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

func (h *KeywordFilter) BuildFilterQuery(f Filters, query string) (string, []interface{}) {
	var args []interface{}
	if f.Keyword != "" {
		keywordList := strings.Split(f.Keyword, ",")
		var keywordConditions []string
		for _, keyword := range keywordList {
			normalizedKeyword := "%" + strings.ToLower(keyword) + "%"
			stemmedKeyword := "%" + porterstemmer.StemString(strings.ToLower(keyword)) + "%"
			condition := "(LOWER(articles.title) ILIKE ? OR LOWER(articles.description) ILIKE ? OR " +
				"LOWER(articles.title) ILIKE ? OR LOWER(articles.description) ILIKE ?)"
			keywordConditions = append(keywordConditions, condition)
			args = append(args, normalizedKeyword, normalizedKeyword, stemmedKeyword, stemmedKeyword)
		}
		if len(keywordConditions) > 0 {
			query += " AND (" + strings.Join(keywordConditions, " OR ") + ")"
		}
	}
	if h.next != nil {
		return h.next.BuildFilterQuery(f, query)
	}
	return query, args
}
