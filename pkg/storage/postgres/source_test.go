package postgres

import (
	"database/sql"
	"testing"

	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestPostgresSrcStorage_GetAll(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	storage := NewSrc(db)

	rows := sqlmock.NewRows([]string{"id", "name", "link", "short_name"}).
		AddRow(1, "source1", "link1", "short1").
		AddRow(2, "source2", "link2", "short2")

	mock.ExpectQuery(`SELECT id, name, link, short_name FROM sources`).
		WillReturnRows(rows)

	sources, err := storage.GetAll()
	assert.NoError(t, err)
	assert.Len(t, sources, 2)

	expectedSources := []model.Source{
		{Id: 1, Name: "source1", Link: "link1", ShortName: "short1"},
		{Id: 2, Name: "source2", Link: "link2", ShortName: "short2"},
	}

	assert.Equal(t, expectedSources, sources)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresSrcStorage_Save(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	storage := NewSrc(db)

	mock.ExpectQuery(`INSERT INTO sources \(name, link, short_name\) VALUES \(\$1, \$2\, \$3\) RETURNING id`).
		WithArgs("source1", "link1", "short1").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	source := model.Source{Name: "source1", Link: "link1", ShortName: "short1"}
	savedSource, err := storage.Save(source)
	assert.NoError(t, err)
	assert.Equal(t, 1, savedSource.Id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresSrcStorage_SaveAll(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	storage := NewSrc(db)

	mock.ExpectQuery(`INSERT INTO sources \(name, link, short_name\) VALUES \(\$1, \$2, \$3\) RETURNING id`).
		WithArgs("source1", "link1", "short1").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectQuery(`INSERT INTO sources \(name, link, short_name\) VALUES \(\$1, \$2, \$3\) RETURNING id`).
		WithArgs("source2", "link2", "short2").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

	sources := []model.Source{
		{Name: "source1", Link: "link1", ShortName: "short1"},
		{Name: "source2", Link: "link2", ShortName: "short2"},
	}

	err = storage.SaveAll(sources)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresSrcStorage_Delete(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	storage := NewSrc(db)

	mock.ExpectExec(`DELETE FROM sources WHERE id = \$1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = storage.Delete(1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresSrcStorage_GetByID(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	storage := NewSrc(db)

	rows := sqlmock.NewRows([]string{"id", "name", "link"}).
		AddRow(1, "source1", "link1")

	mock.ExpectQuery(`SELECT id, name, link FROM sources WHERE id = \$1`).
		WithArgs(1).
		WillReturnRows(rows)

	source, err := storage.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, model.Source{Id: 1, Name: "source1", Link: "link1"}, source)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresSrcStorage_Update(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	storage := NewSrc(db)

	mock.ExpectExec(`UPDATE sources SET name = \$1, link = \$2 WHERE id = \$3`).
		WithArgs("updated source", "updated link", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	source := model.Source{Name: "updated source", Link: "updated link"}
	updatedSource, err := storage.Update(1, source)
	assert.NoError(t, err)
	assert.Equal(t, 1, updatedSource.Id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresSrcStorage_GetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	storage := NewSrc(db)

	mock.ExpectQuery(`SELECT id, name, link FROM sources WHERE id = \$1`).
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	_, err = storage.GetByID(1)
	assert.Error(t, err)
	assert.Equal(t, "source with id 1 not found", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}
