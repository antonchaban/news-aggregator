package postgres

import (
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/antonchaban/news-aggregator/pkg/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresArticleStorage struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) storage.ArticleStorage {
	return &postgresArticleStorage{db: db}
}

func (pa *postgresArticleStorage) GetAll() ([]model.Article, error) {
	//TODO implement me
	panic("implement me")
}

func (pa *postgresArticleStorage) Save(article model.Article) (model.Article, error) {
	//TODO implement me
	panic("implement me")
}

func (pa *postgresArticleStorage) SaveAll(articles []model.Article) error {
	//TODO implement me
	panic("implement me")
}

func (pa *postgresArticleStorage) Delete(id int) error {
	//TODO implement me
	panic("implement me")
}

func (pa *postgresArticleStorage) DeleteBySourceID(id int) error {
	//TODO implement me
	panic("implement me")
}
