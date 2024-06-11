package filter

import (
	"strings"

	"news-aggregator/pkg/model"
)

type KeywordFilter struct {
	next ArticleFilter
}

func (h *KeywordFilter) SetNext(handler ArticleFilter) ArticleFilter {
	h.next = handler
	return handler
}

func (h *KeywordFilter) Filter(articles []model.Article, f Filters) ([]model.Article, error) {
	if f.Keyword != "" {
		keywordList := strings.Split(f.Keyword, ",")
		var keywordFilteredArticles []model.Article
		for _, article := range articles {
			for _, keyword := range keywordList {
				if strings.Contains(article.Title, strings.TrimSpace(keyword)) || strings.Contains(article.Description, strings.TrimSpace(keyword)) {
					keywordFilteredArticles = append(keywordFilteredArticles, article)
				}
			}
		}
		articles = intersect(articles, keywordFilteredArticles)
	}
	if h.next != nil {
		return h.next.Filter(articles, f)
	}
	return articles, nil
}
