package main

import (
	"flag"
	"fmt"
	"log"
	"news-aggregator/pkg/model"
	"news-aggregator/pkg/parser"
	"news-aggregator/pkg/repository"
	"news-aggregator/pkg/service"
	"strings"
)

func main() {
	// Init CLI arguments
	help := flag.Bool("help", false, "Show all available arguments and their descriptions.")
	sources := flag.String("sources", "", "Select the desired news sources to get the news from. Supported sources: abcnews, bbc, washingtontimes, nbc, usatoday")
	keywords := flag.String("keywords", "", "Specify the keywords to filter the news by.")
	dateStart := flag.String("date-start", "", "Specify the start date to filter the news by (format: YYYY-MM-DD).")
	dateEnd := flag.String("date-end", "", "Specify the end date to filter the news by (format: YYYY-MM-DD).")

	flag.Parse()

	if *help {
		flag.Usage()
		return
	}
	var articles []model.Article
	db := repository.NewArticleInMemory(articles)
	svc := service.NewArticleService(db)

	files := []string{
		"../../data/abcnews-international-category-19-05-24.xml",
		"../../data/bbc-news-category-19-05-24.xml",
		"../../data/washingtontimes-world-category-19-05-24.xml",
		"../../data/nbc-news.json",
		"../../data/usatoday-world-news.html",
	}

	err := parser.LoadArticlesFromFiles(files, svc)
	if err != nil {
		log.Fatalf("Error loading articles from files: %v", err)
	}

	var filteredArticles []model.Article
	if *sources != "" {
		sourceList := strings.Split(*sources, ",")
		for _, source := range sourceList {
			sourceArticles, err := svc.GetBySource(strings.TrimSpace(source))
			if err != nil {
				log.Fatalf("Error fetching articles by source: %v", err)
			}
			filteredArticles = append(filteredArticles, sourceArticles...)
		}
	} else {
		filteredArticles, err = svc.GetAll()
		if err != nil {
			log.Fatalf("Error fetching all articles: %v", err)
		}
	}

	// Filter by keywords
	if *keywords != "" {
		keywordList := strings.Split(*keywords, ",")
		var keywordFilteredArticles []model.Article
		for _, keyword := range keywordList {
			keywordArticles, err := svc.GetByKeyword(strings.TrimSpace(keyword))
			if err != nil {
				log.Fatalf("Error fetching articles by keyword: %v", err)
			}
			keywordFilteredArticles = append(keywordFilteredArticles, keywordArticles...)
		}
		filteredArticles = intersect(filteredArticles, keywordFilteredArticles)
	}

	if *dateStart != "" || *dateEnd != "" {
		dateRangeArticles, err := svc.GetByDateInRange(*dateStart, *dateEnd)
		if err != nil {
			log.Fatalf("Error fetching articles by date range: %v", err)
		}
		filteredArticles = intersect(filteredArticles, dateRangeArticles)
	}

	uniqueArticles := make(map[int]model.Article)
	for _, article := range filteredArticles {
		uniqueArticles[article.Id] = article
	}

	for _, article := range uniqueArticles {
		fmt.Println(article)
	}
}

// intersect returns the common elements between two slices of articles.
func intersect(a, b []model.Article) []model.Article {
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
