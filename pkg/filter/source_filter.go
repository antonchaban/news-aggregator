package filter

import (
	"strings"

	"news-aggregator/pkg/model"
	"news-aggregator/pkg/service"
)

type SourceFilter struct {
	next ArticleFilter
}

func (h *SourceFilter) SetNext(handler ArticleFilter) ArticleFilter {
	h.next = handler
	return handler
}

func (h *SourceFilter) Filter(svc service.ArticleService, _ []model.Article, f Filters) ([]model.Article, error) {
	var filteredArticles []model.Article
	if f.Source != "" {
		sourceList := strings.Split(f.Source, ",")
		for _, source := range sourceList {
			sourceArticles, err := svc.GetBySource(strings.TrimSpace(source))
			if err != nil {
				return nil, err
			}
			filteredArticles = append(filteredArticles, sourceArticles...)
		}
	} else {
		var err error
		filteredArticles, err = svc.GetAll()
		if err != nil {
			return nil, err
		}
	}
	if h.next != nil {
		return h.next.Filter(svc, filteredArticles, f)
	}
	return filteredArticles, nil
}
