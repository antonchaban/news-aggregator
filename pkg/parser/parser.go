package parser

import (
	"fmt"
	"news-aggregator/pkg/model"
	"news-aggregator/pkg/parser/html"
	"news-aggregator/pkg/parser/json"
	"news-aggregator/pkg/parser/rss"
	"news-aggregator/pkg/parser/strategy"
	"news-aggregator/pkg/service"
	"os"
)

func LoadArticlesFromFiles(files []string, svc *service.ArticleService) error {
	context := &strategy.Context{}

	for _, filePath := range files {
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Println("Error opening file:", err)
			continue
		}

		defer file.Close()

		format := DetermineFileFormat(filePath)
		switch format {
		case "rss":
			context.SetParser(&rss.Parser{})
		case "json":
			context.SetParser(&json.Parser{})
		case "html":
			config := html.FeedConfig{
				ArticleSelector:     "a.gnt_m_flm_a",
				LinkSelector:        "",
				DescriptionSelector: "data-c-br",
				PubDateSelector:     "div.gnt_m_flm_sbt",
				Source:              "USA TODAY",
				DateAttribute:       "data-c-dt",
				TimeFormat: []string{
					"12:59 p.m. ET May 19 2006",
					"Jan 02, 2006",
				},
			}

			htmlParser := html.NewHtmlParser(config)
			context.SetParser(htmlParser)
		default:
			fmt.Println("Unsupported file format:", filePath)
			continue
		}

		parsedArticles, err := context.Parse(file)
		if err != nil {
			fmt.Println("Error parsing file with", format, "parser:", err)
			continue
		}

		for _, item := range parsedArticles {
			article := model.Article{
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

	return nil
}
