package filter

import (
	"news-aggregator/pkg/model"
)

type ArticleFilter interface {
	SetNext(handler ArticleFilter) ArticleFilter
	Filter(articles []model.Article, f Filters) ([]model.Article, error)
}

type Filters struct {
	Keyword   string
	Source    string
	StartDate string
	EndDate   string
}
