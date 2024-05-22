package main

import (
	"fmt"
	"news-aggregator/pkg/model"
	"news-aggregator/pkg/parser"
	"news-aggregator/pkg/repository"
	"news-aggregator/pkg/service"
	"os"
)

func main() {
	files := []string{
		"data/abcnews-international-category-19-05-24.xml",
		"data/bbc-world-category-19-05-24.xml",
		"data/washingtontimes-world-category-19-05-24.xml",
		"data/nbc-news.json",
		"data/usatoday-world-news.html",
	}

	var articles []model.Article
	db := repository.NewArticleInMemory(articles)
	svc := service.NewArticleService(db)
	context := &parser.Context{}

	for _, filePath := range files {
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Println("Error opening file:", err)
			continue
		}

		defer file.Close()

		format := parser.DetermineFileFormat(filePath)
		switch format {
		case "rss":
			context.SetParser(&parser.RssParser{})
		case "json":
			context.SetParser(&parser.JsonParser{})
		case "html":
			config := parser.HtmlFeedConfig{
				ArticleSelector:     "a.gnt_m_flm_a",
				LinkSelector:        "",
				DescriptionSelector: "data-c-br",
				PubDateSelector:     "div.gnt_m_flm_sbt",
				Source:              "USA TODAY",
				DateAttribute:       "data-c-dt",
			}

			htmlParser := parser.NewHtmlParser(config)
			context.SetParser(htmlParser)
		default:
			fmt.Println("Unsupported file format:", filePath)
			continue
		}

		articles, err := context.Parse(file)
		if err != nil {
			fmt.Println("Error parser file with", format, "parser:", err)
			continue
		}

		for _, item := range articles {
			article := model.Article{
				Id:          len(db.Articles) + 1,
				Title:       item.Title,
				Link:        item.Link,
				Description: item.Description,
				Source:      item.Source,
				PubDate:     item.PubDate,
			}
			_, err := svc.Create(article)
			if err != nil {
				fmt.Println("Error creating article:", err)
				continue
			}
		}
	}

	articlesInDb, _ := svc.GetAll()
	for _, article := range articlesInDb {
		fmt.Println(article)
	}
}
