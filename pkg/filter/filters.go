package filter

import (
	"news-aggregator/pkg/model"
	"news-aggregator/pkg/service"
	"strings"
)

type Filters struct {
	Keyword   string
	Source    string
	StartDate string
	EndDate   string
}

func FilterArticles(svc service.ArticleService, f Filters) ([]model.Article, error) {
	var filteredArticles []model.Article
	var err error

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
		filteredArticles, err = svc.GetAll()
		if err != nil {
			return nil, err
		}
	}

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
		filteredArticles = intersect(filteredArticles, keywordFilteredArticles)
	}

	if f.StartDate != "" || f.EndDate != "" {
		dateRangeArticles, err := svc.GetByDateInRange(f.StartDate, f.EndDate)
		if err != nil {
			return nil, err
		}
		filteredArticles = intersect(filteredArticles, dateRangeArticles)
	}

	uniqueArticles := make(map[int]model.Article)
	for _, article := range filteredArticles {
		uniqueArticles[article.Id] = article
	}

	var result []model.Article
	for _, article := range uniqueArticles {
		result = append(result, article)
	}
	return result, nil
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
