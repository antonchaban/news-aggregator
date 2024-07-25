package filter

import (
	"fmt"
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	eventDateRangeFilterStart = "date_range_filter_start"
	eventParseStartDateError  = "parse_start_date_error"
	eventParseEndDateError    = "parse_end_date_error"
)

// DateRangeFilter filters articles based on their publication date.
type DateRangeFilter struct {
	next ArticleFilter
}

// SetNext sets the next filter in the chain and returns the filter.
func (h *DateRangeFilter) SetNext(filter ArticleFilter) ArticleFilter {
	h.next = filter
	return filter
}

// Filter filters articles by their publication date based on the provided Filters.
func (h *DateRangeFilter) Filter(articles []model.Article, f Filters) ([]model.Article, error) {
	logrus.WithField("event_id", eventDateRangeFilterStart).Info("Starting DateRangeFilter")
	var startDateObj, endDateObj time.Time
	var err error

	if f.StartDate != "" {
		startDateObj, err = time.Parse("2006-01-02", f.StartDate)
		if err != nil {
			logrus.WithField("event_id", eventParseStartDateError).Errorf("Failed to parse start date: %v", err)
			return nil, fmt.Errorf("failed to parse start date: %v", err)
		}
	} else {
		startDateObj = time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	if f.EndDate != "" {
		endDateObj, err = time.Parse("2006-01-02", f.EndDate)
		if err != nil {
			logrus.WithField("event_id", eventParseEndDateError).Errorf("Failed to parse end date: %v", err)
			return nil, fmt.Errorf("failed to parse end date: %v", err)
		}
	} else {
		endDateObj = time.Now()
	}

	logrus.WithFields(logrus.Fields{
		"start_date": startDateObj,
		"end_date":   endDateObj,
	}).Info("Filtering articles by date range")

	var dateRangeFilteredArticles []model.Article
	for _, article := range articles {
		if (f.StartDate == "" || article.PubDate.After(startDateObj) || article.PubDate.Equal(startDateObj)) &&
			(f.EndDate == "" || article.PubDate.Before(endDateObj) || article.PubDate.Equal(endDateObj)) {
			dateRangeFilteredArticles = append(dateRangeFilteredArticles, article)
		}
	}

	logrus.WithField("filtered_count", len(dateRangeFilteredArticles)).Info("Date range filtering complete")

	articles = dateRangeFilteredArticles

	if h.next != nil {
		return h.next.Filter(articles, f)
	}
	return articles, nil
}
