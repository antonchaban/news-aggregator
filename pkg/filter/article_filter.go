package filter

import (
	"github.com/antonchaban/news-aggregator/pkg/model"
)

// ArticleFilter interface defines methods for setting the next handler in the chain and for filtering articles.
type ArticleFilter interface {
	// SetNext sets the next handler in the chain and returns the handler.
	SetNext(handler ArticleFilter) ArticleFilter
	// Filter filters the articles based on the given filters and returns the filtered articles or an error.
	Filter(articles []model.Article, f Filters) ([]model.Article, error)
	BuildFilterQuery(f Filters, query string) (string, []interface{})
}

// Filters struct contains the filtering criteria.
type Filters struct {
	Keyword   string
	Source    string
	StartDate string
	EndDate   string
	UseDB     bool
}
