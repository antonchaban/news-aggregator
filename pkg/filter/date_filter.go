package filter

import (
	"fmt"
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/sirupsen/logrus"
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
	logrus.WithField("event_id", "date_range_filter_start").Info("Starting DateRangeFilter")
	var startDateObj, endDateObj time.Time
	var err error

	if f.StartDate != "" {
		startDateObj, err = time.Parse("2006-01-02", f.StartDate)
		if err != nil {
			logrus.WithField("event_id", "parse_start_date_error").Errorf("Failed to parse start date: %v", err)
			return nil, fmt.Errorf("failed to parse start date: %v", err)
		}
	} else {
		startDateObj = time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	if f.EndDate != "" {
		endDateObj, err = time.Parse("2006-01-02", f.EndDate)
		if err != nil {
			logrus.WithField("event_id", "parse_end_date_error").Errorf("Failed to parse end date: %v", err)
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

func (h *DateRangeFilter) BuildFilterQuery(f Filters, query string) (string, []interface{}) {
	if f.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", f.StartDate)
		if err != nil {
			logrus.WithField("event_id", "parse_start_date_error").Errorf("Failed to parse start date: %v", err)
			return query, nil
		}
		query += fmt.Sprintf(" AND pub_date >= '%s'", startDate.Format("2006-01-02"))
	}

	if f.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", f.EndDate)
		if err != nil {
			logrus.WithField("event_id", "parse_end_date_error").Errorf("Failed to parse end date: %v", err)
			return query, nil
		}
		query += fmt.Sprintf(" AND pub_date <= '%s'", endDate.Format("2006-01-02"))
	}

	if h.next != nil {
		return h.next.BuildFilterQuery(f, query)
	}
	return query, nil
}
