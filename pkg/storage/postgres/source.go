package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/antonchaban/news-aggregator/pkg/storage"
	"github.com/jmoiron/sqlx"
)

type postgresSrcStorage struct {
	db *sqlx.DB
}

func NewSrc(db *sqlx.DB) storage.SourceStorage {
	return &postgresSrcStorage{db: db}
}

func (psrc *postgresSrcStorage) GetAll() ([]model.Source, error) {
	var sources []model.Source
	query := fmt.Sprintf(`SELECT id, name, link FROM sources`)
	err := psrc.db.Select(&sources, query)
	if err != nil {
		return nil, err

	}
	return sources, nil
	/*	rows, err := psrc.db.Query(context.Background(), "SELECT id, name, link FROM sources")
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var sources []model.Source
		for rows.Next() {
			var src model.Source
			err = rows.Scan(&src.Id, &src.Name, &src.Link)
			if err != nil {
				return nil, err
			}
			sources = append(sources, src)
		}
		return sources, nil*/
}

func (psrc *postgresSrcStorage) Save(src model.Source) (model.Source, error) {
	var id int
	createQuery := fmt.Sprintf(`INSERT INTO sources (name, link) VALUES ($1, $2) RETURNING id`)
	err := psrc.db.QueryRow(createQuery, src.Name, src.Link).Scan(&id)
	if err != nil {
		return model.Source{}, err
	}
	src.Id = id
	return src, nil
	/*	var id int
		err := psrc.db.QueryRow(context.Background(),
			"INSERT INTO sources (name, link) VALUES ($1, $2) RETURNING id", src.Name, src.Link).Scan(&id)
		if err != nil {
			return model.Source{}, err
		}
		src.Id = id
		return src, nil*/
}

func (psrc *postgresSrcStorage) SaveAll(sources []model.Source) error {
	for _, src := range sources {
		_, err := psrc.Save(src)
		if err != nil {
			return err
		}
	}
	return nil

	/*	tx, err := psrc.db.Begin(context.Background())
		if err != nil {
			return err
		}
		defer tx.Rollback(context.Background())

		for _, src := range sources {
			_, err := tx.Exec(context.Background(), "INSERT INTO sources (name, link) VALUES ($1, $2)", src.Name, src.Link)
			if err != nil {
				return err
			}
		}

		err = tx.Commit(context.Background())
		if err != nil {
			return err
		}
		return nil*/
}

func (psrc *postgresSrcStorage) Delete(id int) error {
	query := fmt.Sprintf(`DELETE FROM sources WHERE id = $1`)
	_, err := psrc.db.Exec(query, id)
	return err

	/*	_, err := psrc.db.Exec(context.Background(), "DELETE FROM sources WHERE id = $1", id)
		return err*/
}

func (psrc *postgresSrcStorage) GetByID(id int) (model.Source, error) {
	var src model.Source
	query := fmt.Sprintf(`SELECT id, name, link FROM sources WHERE id = $1`)
	err := psrc.db.Get(&src, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Source{}, fmt.Errorf("source with id %d not found", id)
		}
		return model.Source{}, err
	}
	return src, nil

	/*	var src model.Source
		err := psrc.db.QueryRow(context.Background(), "SELECT id, name, link FROM sources WHERE id = $1", id).Scan(&src.Id, &src.Name, &src.Link)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return model.Source{}, fmt.Errorf("source with id %d not found", id)
			}
			return model.Source{}, err
		}
		return src, nil*/
}

func (psrc *postgresSrcStorage) Update(id int, src model.Source) (model.Source, error) {
	query := fmt.Sprintf(`UPDATE sources SET name = $1, link = $2 WHERE id = $3`)
	_, err := psrc.db.Exec(query, src.Name, src.Link, id)
	if err != nil {
		return model.Source{}, err
	}
	src.Id = id
	return src, nil
	/*	_, err := psrc.db.Exec(context.Background(), "UPDATE sources SET name = $1, link = $2 WHERE id = $3", src.Name, src.Link, id)
		if err != nil {
			return model.Source{}, err
		}
		src.Id = id
		return src, nil*/
}
