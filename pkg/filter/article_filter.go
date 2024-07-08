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
}

// Filters struct contains the filtering criteria.
type Filters struct {
	Keyword   string // Keywords to filter the articles by.
	Source    string // Sources to filter the articles by.
	StartDate string // Start date to filter the articles by.
	EndDate   string // End date to filter the articles by.
}
