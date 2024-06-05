package storage

import (
	"news-aggregator/pkg/model"
	"time"
)

//go:generate mockgen -destination=mocks/mock_article.go -package=mocks news-aggregator/pkg/storage Article

// Article is an interface that defines the methods for interacting with the article storage.
type Article interface {
	GetAll() ([]model.Article, error)
	Create(article model.Article) (model.Article, error)
	SaveAll(articles []model.Article) error
	Delete(id int) error
	GetByKeyword(keyword string) ([]model.Article, error)
	GetBySource(source string) ([]model.Article, error)
	GetByDateInRange(startDate, endDate time.Time) ([]model.Article, error)
}
