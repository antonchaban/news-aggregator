package filter

import (
	"news-aggregator/pkg/model"
	"news-aggregator/pkg/service"
)

type DateRangeFilter struct {
	next ArticleFilter
}

func (h *DateRangeFilter) SetNext(handler ArticleFilter) ArticleFilter {
	h.next = handler
	return handler
}

func (h *DateRangeFilter) Filter(svc service.ArticleService, articles []model.Article, f Filters) ([]model.Article, error) {
	if f.StartDate != "" || f.EndDate != "" {
		dateRangeArticles, err := svc.GetByDateInRange(f.StartDate, f.EndDate)
		if err != nil {
			return nil, err
		}
		articles = intersect(articles, dateRangeArticles)
	}
	if h.next != nil {
		return h.next.Filter(svc, articles, f)
	}
	return articles, nil
}
