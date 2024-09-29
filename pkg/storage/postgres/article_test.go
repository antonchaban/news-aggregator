package postgres

import (
	"testing"
	"time"

	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestPostgresArticleStorage_GetAll(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	storage := New(db)

	rows := sqlmock.NewRows([]string{"id", "title", "description", "link", "pub_date", "source_id", "source_name", "source_link", "source_short_name"}).
		AddRow(1, "title1", "description1", "link1", time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC), 1, "source1", "source_link1", "short1").
		AddRow(2, "title2", "description2", "link2", time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC), 2, "source2", "source_link2", "short2")

	mock.ExpectQuery("SELECT a.id, a.title, a.description, a.link, a.pub_date, s.id AS source_id, s.name AS source_name, s.link AS source_link, s.short_name AS source_short_name FROM articles a JOIN sources s ON a.source_id = s.id").
		WillReturnRows(rows)

	articles, err := storage.GetAll()
	assert.NoError(t, err)
	assert.Len(t, articles, 2)

	expectedArticles := []model.Article{
		{
			Id:          1,
			Title:       "title1",
			Description: "description1",
			Link:        "link1",
			PubDate:     time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC),
			Source: model.Source{
				Id:        1,
				Name:      "source1",
				Link:      "source_link1",
				ShortName: "short1",
			},
		},
		{
			Id:          2,
			Title:       "title2",
			Description: "description2",
			Link:        "link2",
			PubDate:     time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
			Source: model.Source{
				Id:        2,
				Name:      "source2",
				Link:      "source_link2",
				ShortName: "short2",
			},
		},
	}

	for i, article := range articles {
		assert.Equal(t, expectedArticles[i].Id, article.Id)
		assert.Equal(t, expectedArticles[i].Title, article.Title)
		assert.Equal(t, expectedArticles[i].Description, article.Description)
		assert.Equal(t, expectedArticles[i].Link, article.Link)
		assert.WithinDuration(t, expectedArticles[i].PubDate, article.PubDate, time.Second)
		assert.Equal(t, expectedArticles[i].Source.Id, article.Source.Id)
		assert.Equal(t, expectedArticles[i].Source.Name, article.Source.Name)
		assert.Equal(t, expectedArticles[i].Source.Link, article.Source.Link)
		assert.Equal(t, expectedArticles[i].Source.ShortName, article.Source.ShortName)
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresArticleStorage_Save(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	storage := New(db)

	mock.ExpectQuery("INSERT INTO articles").
		WithArgs("title1", "description1", "link1", 1, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	article := model.Article{
		Title:       "title1",
		Description: "description1",
		Link:        "link1",
		Source:      model.Source{Id: 1},
		PubDate:     time.Now(),
	}

	savedArticle, err := storage.Save(article)
	assert.NoError(t, err)
	assert.Equal(t, 1, savedArticle.Id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresArticleStorage_SaveAll(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	storage := New(db)

	mock.ExpectQuery("INSERT INTO articles").
		WithArgs("title1", "description1", "link1", 1, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectQuery("INSERT INTO articles").
		WithArgs("title2", "description2", "link2", 2, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

	articles := []model.Article{
		{Title: "title1", Description: "description1", Link: "link1", Source: model.Source{Id: 1}, PubDate: time.Now()},
		{Title: "title2", Description: "description2", Link: "link2", Source: model.Source{Id: 2}, PubDate: time.Now()},
	}

	err = storage.SaveAll(articles)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresArticleStorage_Delete(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	storage := New(db)

	mock.ExpectExec("DELETE FROM articles WHERE id = \\$1").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = storage.Delete(1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresArticleStorage_DeleteBySourceID(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	storage := New(db)

	mock.ExpectExec("DELETE FROM articles WHERE source_id = \\$1").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = storage.DeleteBySourceID(1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
