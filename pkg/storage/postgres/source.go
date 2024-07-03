package postgres

import (
	"context"
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/antonchaban/news-aggregator/pkg/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresSrcStorage struct {
	db *pgxpool.Pool
}

func NewSrc(db *pgxpool.Pool) storage.SourceStorage {
	return &postgresSrcStorage{db: db}
}

func (psrc *postgresSrcStorage) GetAll() ([]model.Source, error) {
	rows, err := psrc.db.Query(context.Background(), "SELECT id, name, link FROM sources")
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
	return sources, nil
}

func (psrc *postgresSrcStorage) Save(src model.Source) (model.Source, error) {
	var id int
	err := psrc.db.QueryRow(context.Background(),
		"INSERT INTO sources (name, link) VALUES ($1, $2) RETURNING id", src.Name, src.Link).Scan(&id)
	if err != nil {
		return model.Source{}, err
	}
	src.Id = id
	return src, nil
}

func (psrc *postgresSrcStorage) SaveAll(sources []model.Source) error {
	//TODO implement me
	panic("implement me")
}

func (psrc *postgresSrcStorage) Delete(id int) error {
	//TODO implement me
	panic("implement me")
}

func (psrc *postgresSrcStorage) GetByID(id int) (model.Source, error) {
	//TODO implement me
	panic("implement me")
}

func (psrc *postgresSrcStorage) Update(id int, src model.Source) (model.Source, error) {
	//TODO implement me
	panic("implement me")
}
