package filter

import (
	"strings"

	"news-aggregator/pkg/model"
	"news-aggregator/pkg/service"
)

type KeywordFilter struct {
	next ArticleFilter
}

func (h *KeywordFilter) SetNext(handler ArticleFilter) ArticleFilter {
	h.next = handler
	return handler
}

func (h *KeywordFilter) Filter(svc service.ArticleService, articles []model.Article, f Filters) ([]model.Article, error) {
	if f.Keyword != "" {
		keywordList := strings.Split(f.Keyword, ",")
		var keywordFilteredArticles []model.Article
		for _, keyword := range keywordList {
			keywordArticles, err := svc.GetByKeyword(strings.TrimSpace(keyword))
			if err != nil {
				return nil, err
			}
			keywordFilteredArticles = append(keywordFilteredArticles, keywordArticles...)
		}
		articles = intersect(articles, keywordFilteredArticles)
	}
	if h.next != nil {
		return h.next.Filter(svc, articles, f)
	}
	return articles, nil
}
