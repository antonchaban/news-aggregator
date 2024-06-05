package storage

import (
	"news-aggregator/pkg/model"
	"time"
)

// ArticleStorage is an interface that defines the methods for interacting with the article storage.
type ArticleStorage interface {
	GetAll() ([]model.Article, error)
	GetById(id int) (model.Article, error)
	Create(article model.Article) (model.Article, error)
	SaveAll(articles []model.Article) error
	Delete(id int) error
	GetByKeyword(keyword string) ([]model.Article, error)
	GetBySource(source string) ([]model.Article, error)
	GetByDateInRange(startDate, endDate time.Time) ([]model.Article, error)
}
