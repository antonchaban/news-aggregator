package filter

import (
	"github.com/antonchaban/news-aggregator/pkg/model"
)

type ArticleFilter interface {
	SetNext(handler ArticleFilter) ArticleFilter
	Filter(articles []model.Article, f Filters) ([]model.Article, error)
	BuildFilterQuery(f Filters) (string, []interface{})
}

type Filters struct {
	Keyword   string
	Source    string
	StartDate string
	EndDate   string
}
