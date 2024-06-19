package filter

import (
	"fmt"
	"github.com/antonchaban/news-aggregator/pkg/model"
	"time"
)

type DateRangeFilter struct {
	next ArticleFilter
}

func (h *DateRangeFilter) SetNext(handler ArticleFilter) ArticleFilter {
	h.next = handler
	return handler
}

func (h *DateRangeFilter) Filter(articles []model.Article, f Filters) ([]model.Article, error) {
	var startDateObj, endDateObj time.Time
	var err error

	if f.StartDate != "" {
		startDateObj, err = time.Parse("2006-01-02", f.StartDate)
		if err != nil {
			return nil, fmt.Errorf("failed to parse start date: %v", err)
		}
	} else {
		startDateObj = time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)

	}

	if f.EndDate != "" {
		endDateObj, err = time.Parse("2006-01-02", f.EndDate)
		if err != nil {
			return nil, fmt.Errorf("failed to parse end date: %v", err)
		}
	} else {
		endDateObj = time.Now()
	}

	// Filter articles by date range
	var dateRangeFilteredArticles []model.Article
	for _, article := range articles {
		if (f.StartDate == "" || article.PubDate.After(startDateObj) || article.PubDate.Equal(startDateObj)) &&
			(f.EndDate == "" || article.PubDate.Before(endDateObj) || article.PubDate.Equal(endDateObj)) {
			dateRangeFilteredArticles = append(dateRangeFilteredArticles, article)
		}
	}

	articles = dateRangeFilteredArticles

	if h.next != nil {
		return h.next.Filter(articles, f)
	}
	return articles, nil
}
