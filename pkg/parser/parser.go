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

// ParseArticlesFromFile Parse function takes a file and returns a slice of parsed articles
func ParseArticlesFromFile(file string) ([]model.Article, error) {
	f, err := os.Open(file)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)

	format := DetermineFileFormat(file)
	parser, err := createParser(format)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	articles, err := parser.ParseFile(f)
	if err != nil {
		fmt.Println("Error parsing file with", format, "parser:", err)
		return nil, err
	}

	return articles, nil
}

func createParser(format string) (Parser, error) {
	switch format {
	case rssFormat:
		return &rss.Parser{}, nil
	case jsonFormat:
		return &json.Parser{}, nil
	case htmlFormat:
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
		return html.NewHtmlParser(config), nil
	default:
		return nil, fmt.Errorf("unsupported file format: %s", format)
	}
}
