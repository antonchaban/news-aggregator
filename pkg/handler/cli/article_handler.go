package cli

import (
	"fmt"
	"log"
	"news-aggregator/pkg/model"
	"news-aggregator/pkg/parser"
	"os"
	"path/filepath"
	"strings"
)

// loadData loads articles from the specified files and saves them to the service.
func (h *Handler) loadData() []model.Article {
	execDir, err := os.Executable()
	var articles []model.Article
	if err != nil {
		log.Fatalf("Error getting executable directory: %v", err)
	}
	execDir = filepath.Dir(execDir)

	// Define file paths relative to the executable directory
	files := []string{
		filepath.Join(execDir, "../../data/abcnews-international-category-19-05-24.xml"),
		filepath.Join(execDir, "../../data/bbc-world-category-19-05-24.xml"),
		filepath.Join(execDir, "../../data/washingtontimes-world-category-19-05-24.xml"),
		filepath.Join(execDir, "../../data/nbc-news.json"),
		filepath.Join(execDir, "../../data/usatoday-world-news.html"),
	}
	articles, err = parser.ParseArticlesFromFiles(files)
	if err != nil {
		log.Fatalf("Error loading articles from files: %v", err)
	}

	h.service.SaveAll(articles)

	allArticles, err := h.service.GetAll()
	if err != nil {
		log.Fatalf("Error fetching all articles: %v", err)
	}
	return allArticles
}

// filterArticles filters the provided articles based on the provided sources, keywords, and date range.
// It returns the filtered articles.
func (h *Handler) filterArticles(articles []model.Article, sources, keywords, dateStart, dateEnd string) []model.Article {
	var filteredArticles []model.Article

	if sources != "" {
		sourceList := strings.Split(sources, ",")
		for _, source := range sourceList {
			sourceArticles, err := h.service.GetBySource(strings.TrimSpace(source))
			if err != nil {
				log.Fatalf("Error fetching articles by source: %v", err)
			}
			filteredArticles = append(filteredArticles, sourceArticles...)
		}
	} else {
		filteredArticles = articles
	}

	if keywords != "" {
		keywordList := strings.Split(keywords, ",")
		var keywordFilteredArticles []model.Article
		for _, keyword := range keywordList {
			keywordArticles, err := h.service.GetByKeyword(strings.TrimSpace(keyword))
			if err != nil {
				log.Fatalf("Error fetching articles by keyword: %v", err)
			}
			keywordFilteredArticles = append(keywordFilteredArticles, keywordArticles...)
		}
		filteredArticles = intersect(filteredArticles, keywordFilteredArticles)
	}

	if dateStart != "" || dateEnd != "" {
		dateRangeArticles, err := h.service.GetByDateInRange(dateStart, dateEnd)
		if err != nil {
			log.Fatalf("Error fetching articles by date range: %v", err)
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
	return result
}

func (h *Handler) printArticles(articles []model.Article) {
	for _, article := range articles {
		fmt.Println(article)
	}
}

// intersect returns the common elements between two slices of articles.
func intersect(a, b []model.Article) (articles []model.Article) {
	articleMap := make(map[int]model.Article)
	for _, article := range a {
		articleMap[article.Id] = article
	}

	var intersection []model.Article
	for _, article := range b {
		if _, found := articleMap[article.Id]; found {
			intersection = append(intersection, article)
		}
	}
	return intersection
}
