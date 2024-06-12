package parser

import (
	"fmt"
	"news-aggregator/pkg/model"
	"news-aggregator/pkg/parser/html"
	"news-aggregator/pkg/parser/json"
	"news-aggregator/pkg/parser/rss"
	"os"
)

// Parser is an interface that defines parsing strategy
type Parser interface {
	ParseFile(f *os.File) ([]model.Article, error)
}

func ParseArticlesFromFiles(files []string) ([]model.Article, error) {
	var parsedArticles []model.Article
	var parser Parser
	for _, filePath := range files {
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Println("Error opening file:", err)
			continue
		}

		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				panic(err)
			}
		}(file)

		format := DetermineFileFormat(filePath)

		switch format {
		case "rss":
			parser = &rss.Parser{}
		case "json":
			parser = &json.Parser{}
		case "html":
			config := html.FeedConfig{
				ArticleSelector:     "a.gnt_m_flm_a",
				LinkSelector:        "",
				DescriptionSelector: "data-c-br",
				PubDateSelector:     "div.gnt_m_flm_sbt",
				Source:              "USA TODAY",
				DateAttribute:       "data-c-dt",
				TimeFormat: []string{
					"2006-01-02 15:04",
					"Jan 02, 2006",
				},
			}
			parser = html.NewHtmlParser(config)
		default:
			fmt.Println("Unsupported file format:", filePath)
			continue
		}

		articles, err := parser.ParseFile(file)
		parsedArticles = append(parsedArticles, articles...)
		if err != nil {
			fmt.Println("Error parsing file with", format, "parser:", err)
			continue
		}

	}

	return parsedArticles, nil
}
