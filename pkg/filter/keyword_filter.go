package filter

import (
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/reiver/go-porterstemmer"
	"github.com/sirupsen/logrus"
	"strings"
)

type KeywordFilter struct {
	next ArticleFilter
}

func (h *KeywordFilter) SetNext(handler ArticleFilter) ArticleFilter {
	h.next = handler
	return handler
}

func (h *KeywordFilter) Filter(articles []model.Article, f Filters) ([]model.Article, error) {
	logrus.WithField("event_id", "keyword_filter_start").Info("Starting KeywordFilter")

	if f.Keyword != "" {
		keywordList := strings.Split(f.Keyword, ",")
		var keywordFilteredArticles []model.Article
		for _, article := range articles {
			for _, keyword := range keywordList {
				normalizedTitle := strings.ToLower(article.Title)
				normalizedDesc := strings.ToLower(article.Description)
				stemmedKeyword := porterstemmer.StemString(keyword)
				stemmedTitle := porterstemmer.StemString(normalizedTitle)
				stemmedDesc := porterstemmer.StemString(normalizedDesc)
				if strings.Contains(stemmedTitle, stemmedKeyword) ||
					strings.Contains(stemmedDesc, stemmedKeyword) {
					keywordFilteredArticles = append(keywordFilteredArticles, article)
					break // Avoid adding the same article multiple times for different keywords
				}
			}
		}
		articles = keywordFilteredArticles
		logrus.WithField("filtered_count", len(keywordFilteredArticles)).Info("Keyword filtering complete")
	}

	if h.next != nil {
		return h.next.Filter(articles, f)
	}
	return articles, nil
}

func (h *KeywordFilter) BuildFilterQuery(f Filters) (string, []interface{}) {
	return "", nil // todo implement
}
