package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/antonchaban/news-aggregator/pkg/service"
	"github.com/jmoiron/sqlx"
)

type postgresSrcStorage struct {
	db *sqlx.DB
}

func NewSrc(db *sqlx.DB) service.SourceStorage {
	return &postgresSrcStorage{db: db}
}

func (psrc *postgresSrcStorage) GetAll() ([]model.Source, error) {
	var sources []model.Source
	query := `SELECT id, name, link, short_name FROM sources`
	err := psrc.db.Select(&sources, query)
	if err != nil {
		return nil, err

	}
	return sources, nil
}

func (psrc *postgresSrcStorage) Save(src model.Source) (model.Source, error) {
	var id int
	createQuery := `INSERT INTO sources (name, link, short_name) VALUES ($1, $2, $3) RETURNING id`
	err := psrc.db.QueryRow(createQuery, src.Name, src.Link, src.ShortName).Scan(&id)
	if err != nil {
		return model.Source{}, err
	}
	src.Id = id
	return src, nil
}

func (psrc *postgresSrcStorage) SaveAll(sources []model.Source) error {
	for _, src := range sources {
		_, err := psrc.Save(src)
		if err != nil {
			return err
		}
	}
	return nil
}

func (psrc *postgresSrcStorage) Delete(id int) error {
	query := `DELETE FROM sources WHERE id = $1`
	_, err := psrc.db.Exec(query, id)
	return err
}

func (psrc *postgresSrcStorage) GetByID(id int) (model.Source, error) {
	var src model.Source
	query := `SELECT id, name, link FROM sources WHERE id = $1`
	err := psrc.db.Get(&src, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Source{}, fmt.Errorf("source with id %d not found", id)
		}
		return model.Source{}, err
	}
	return src, nil
}

func (psrc *postgresSrcStorage) GetByShortName(shortName string) (model.Source, error) {
	var src model.Source
	query := `SELECT id, name, link, short_name FROM sources WHERE short_name = $1`
	err := psrc.db.Get(&src, query, shortName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Source{}, fmt.Errorf("source with short name %s not found", shortName)
		}
		return model.Source{}, err
	}
	return src, nil
}

func (psrc *postgresSrcStorage) Update(id int, src model.Source) (model.Source, error) {
	query := `UPDATE sources SET name = $1, link = $2, short_name = $3 WHERE id = $4`
	_, err := psrc.db.Exec(query, src.Name, src.Link, src.ShortName, id)
	if err != nil {
		return model.Source{}, err
	}
	src.Id = id
	return src, nil
}
