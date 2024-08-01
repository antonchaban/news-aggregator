package filter

import (
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/reiver/go-porterstemmer"
	"github.com/sirupsen/logrus"
	"strings"
)

const (
	eventKeywordFilterStart       = "keyword_filter_start"
	eventKeywordFilteringComplete = "keyword_filtering_complete"
)

// KeywordFilter filters articles based on keywords in their title or description.
type KeywordFilter struct {
	next ArticleFilter
}

// SetNext sets the next filter in the chain and returns the filter.
func (h *KeywordFilter) SetNext(filter ArticleFilter) ArticleFilter {
	h.next = filter
	return filter
}

// Filter filters articles by keywords in their title or description based on the provided Filters
// using stemming algorithm.
func (h *KeywordFilter) Filter(articles []model.Article, f Filters) ([]model.Article, error) {
	logrus.WithField("event_id", eventKeywordFilterStart).Info("Starting KeywordFilter")

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
		logrus.WithField("filtered_count", len(keywordFilteredArticles)).Info(eventKeywordFilteringComplete)
	}

	if h.next != nil {
		return h.next.Filter(articles, f)
	}
	return articles, nil
}

func (h *KeywordFilter) BuildFilterQuery(f Filters, query string) (string, []interface{}) {
	if f.Keyword != "" {
		keywordList := strings.Split(f.Keyword, ",")
		var keywordConditions []string
		for _, keyword := range keywordList {
			normalizedKeyword := "%" + strings.ToLower(keyword) + "%"
			condition := "(LOWER(title) ILIKE '" + normalizedKeyword + "' OR LOWER(description) ILIKE '" + normalizedKeyword + "')"
			keywordConditions = append(keywordConditions, condition)
		}
		if len(keywordConditions) > 0 {
			query += " AND (" + strings.Join(keywordConditions, " OR ") + ")"
		}
	}
	if h.next != nil {
		return h.next.BuildFilterQuery(f, query)
	}
	return query, nil
}
