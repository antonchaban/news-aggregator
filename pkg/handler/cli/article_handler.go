package cli

import (
	"fmt"
	"github.com/Masterminds/sprig"
	"log"
	"news-aggregator/pkg/filter"
	"news-aggregator/pkg/model"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

// filterArticles filters the provided articles based on the provided sources, keywords, and date range.
// It returns the filtered articles.
func (h *cliHandler) filterArticles(f filter.Filters) ([]model.Article, error) {
	articles, err := h.service.GetByFilter(f)
	if err != nil {
		log.Fatalf("Error filtering articles: %v", err)
		return nil, err
	}
	return articles, nil
}

// printArticles prints the provided articles to the console using template.
func (h *cliHandler) printArticles(articles []model.Article, filters filter.Filters) {
	tmplPath, err := getTemplatePath()
	if err != nil {
		log.Fatalf("Error getting template path: %v", err)
		return
	}
	tmpl, err := template.New("article.tmpl").Funcs(sprig.FuncMap()).Funcs(template.FuncMap{
		"groupBy": groupBy,
		"gt1": func(s string) bool {
			return len(s) > 1
		},
		"highlightKeywords": func(text string) string {
			if filters.Keyword == "" {
				return text
			}
			return highlightKeywords(text, filters.Keyword)
		},
	}).ParseFiles(tmplPath)
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	data := struct {
		Articles []model.Article
		Filters  filter.Filters
	}{
		Articles: articles,
		Filters:  filters,
	}

	err = tmpl.ExecuteTemplate(os.Stdout, "page", data)
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}
}

// getTemplatePath returns the path to the template file.
func getTemplatePath() (string, error) {
	execDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory: %v", err)
		return "", err
	}

	tmplPath := filepath.Join(execDir, "../../../templates", "article.tmpl")

	if _, err := os.Stat(tmplPath); os.IsNotExist(err) {
		log.Fatalf("Template file not found: %s", tmplPath)
		return "", err
	}
	return tmplPath, nil
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
func (h *cliHandler) sortArticles(articles []model.Article, sortOrder string) []model.Article {
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
