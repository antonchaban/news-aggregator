package filter

import (
	"news-aggregator/pkg/model"
)

type ArticleFilter interface {
	SetNext(handler ArticleFilter) ArticleFilter
	Filter(articles []model.Article, f Filters) ([]model.Article, error)
}

type Filters struct {
	Keyword   string
	Source    string
	StartDate string
	EndDate   string
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
