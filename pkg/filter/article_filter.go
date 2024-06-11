package filter

import (
	"news-aggregator/pkg/model"
	"news-aggregator/pkg/service"
)

type ArticleFilter interface {
	SetNext(handler ArticleFilter) ArticleFilter
	Filter(svc service.ArticleService, articles []model.Article, f Filters) ([]model.Article, error)
}

type Filters struct {
	Keyword   string
	Source    string
	StartDate string
	EndDate   string
}
