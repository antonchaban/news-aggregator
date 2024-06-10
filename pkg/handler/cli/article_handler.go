package cli

import (
	"errors"
	"fmt"
	"github.com/Masterminds/sprig"
	"log"
	"news-aggregator/pkg/model"
	"news-aggregator/pkg/parser"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

// loadData loads articles from the specified files and saves them to the Service.
func (h *Handler) loadData() error {
	execDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory: %v", err)
	}
	log.Printf("Current working directory: %s\n", execDir)

	var articles []model.Article
	files := []string{
		filepath.Join(execDir, "../../../data/abcnews-international-category-19-05-24.xml"),
		filepath.Join(execDir, "../../../data/bbc-world-category-19-05-24.xml"),
		filepath.Join(execDir, "../../../data/washingtontimes-world-category-19-05-24.xml"),
		filepath.Join(execDir, "../../../data/nbc-news.json"),
		filepath.Join(execDir, "../../../data/usatoday-world-news.html"),
	}
	for _, file := range files {
		go func(file string) {
			parsedArticles, err := parser.ParseArticlesFromFile(file)
			if err != nil {
				log.Fatalf("Error parsing articles from file: %v", err)
			} else {
				articles = append(articles, parsedArticles...)
			}
		}(file)
	}

	if err != nil {
		return errors.New("error parsing articles from files")
	}

	h.Service.SaveAll(articles)
	return nil
}

// filterArticles filters the provided articles based on the provided sources, keywords, and date range.
// It returns the filtered articles.
func (h *Handler) filterArticles(sources, keywords, dateStart, dateEnd string) []model.Article {
	var filteredArticles []model.Article
	var err error

	if sources != "" {
		sourceList := strings.Split(sources, ",")
		for _, source := range sourceList {
			sourceArticles, err := h.Service.GetBySource(strings.TrimSpace(source))
			if err != nil {
				log.Fatalf("Error fetching articles by source: %v", err)
			}
			filteredArticles = append(filteredArticles, sourceArticles...)
		}
	} else {
		filteredArticles, err = h.Service.GetAll()
		if err != nil {
			errors.New("error fetching all articles")
		}
	}

	if keywords != "" {
		keywordList := strings.Split(keywords, ",")
		var keywordFilteredArticles []model.Article
		for _, keyword := range keywordList {
			keywordArticles, err := h.Service.GetByKeyword(strings.TrimSpace(keyword))
			if err != nil {
				log.Fatalf("Error fetching articles by keyword: %v", err)
			}
			keywordFilteredArticles = append(keywordFilteredArticles, keywordArticles...)
		}
		filteredArticles = intersect(filteredArticles, keywordFilteredArticles)
	}

	if dateStart != "" || dateEnd != "" {
		dateRangeArticles, err := h.Service.GetByDateInRange(dateStart, dateEnd)
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

// printArticles prints the provided articles to the console using template.
func (h *Handler) printArticles(articles []model.Article, sources, keywords, dateStart, dateEnd string) {
	tmplPath := getTemplatePath()
	tmpl, err := template.New("article.tmpl").Funcs(sprig.FuncMap()).Funcs(template.FuncMap{
		"groupBy": groupBy,
		"gt1": func(s string) bool {
			return len(s) > 1
		},
		"highlightKeywords": func(text string) string {
			if keywords == "" {
				return text
			}
			return highlightKeywords(text, keywords)
		},
	}).ParseFiles(tmplPath)
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	data := struct {
		Articles  []model.Article
		Sources   string
		Keywords  string
		DateStart string
		DateEnd   string
	}{
		Articles:  articles,
		Sources:   sources,
		Keywords:  keywords,
		DateStart: dateStart,
		DateEnd:   dateEnd,
	}

	err = tmpl.ExecuteTemplate(os.Stdout, "page", data)
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}
}

// getTemplatePath returns the path to the template file.
func getTemplatePath() string {
	possiblePaths := []string{
		"../../templates/article.tmpl",
		"templates/article.tmpl",
		"../../../templates/article.tmpl",
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	log.Fatalf("Template file not found in any of the expected paths")
	return ""
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

// groupBy groups the provided articles by the specified field.
func groupBy(articles []model.Article, field string) map[string][]model.Article {
	groupedArticles := make(map[string][]model.Article)
	for _, article := range articles {
		var key string
		switch field {
		case "Source":
			key = article.Source
		default:
			key = ""
		}
		groupedArticles[key] = append(groupedArticles[key], article)
	}
	return groupedArticles
}

// sortArticles sorts the provided articles based on the specified sort order.
func (h *Handler) sortArticles(articles []model.Article, sortOrder string) []model.Article {
	sort.Slice(articles, func(i, j int) bool {
		if sortOrder == "ASC" {
			return articles[i].PubDate.Before(articles[j].PubDate)
		}
		return articles[i].PubDate.After(articles[j].PubDate)
	})
	return articles
}

// highlightKeywords highlights with bold text the specified keywords in the provided text.
func highlightKeywords(text string, keywords string) string {
	keywordList := strings.Split(keywords, ",")
	for _, keyword := range keywordList {
		if len(keyword) > 0 {
			highlighted := fmt.Sprintf("\033[1m%s\033[0m", strings.TrimSpace(keyword))
			text = strings.ReplaceAll(text, strings.TrimSpace(keyword), highlighted)
		}
	}
	return text
}
