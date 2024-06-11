package filter

import (
	"news-aggregator/pkg/model"
	"news-aggregator/pkg/service"
)

func FilterArticles(svc service.ArticleService, f Filters) ([]model.Article, error) {
	var articles []model.Article

	sourceFilter := &SourceFilter{}
	keywordFilter := &KeywordFilter{}
	dateRangeFilter := &DateRangeFilter{}

	// Create the chain
	sourceFilter.SetNext(keywordFilter).SetNext(dateRangeFilter)

	// Start filtering
	return sourceFilter.Filter(svc, articles, f)
}

// intersect function to get the intersection of two slices of articles
func intersect(slice1, slice2 []model.Article) []model.Article {
	articleMap := make(map[int]model.Article)
	for _, article := range slice1 {
		articleMap[article.Id] = article
	}

	var intersection []model.Article
	for _, article := range slice2 {
		if _, exists := articleMap[article.Id]; exists {
			intersection = append(intersection, article)
		}
	}
	return intersection
}
