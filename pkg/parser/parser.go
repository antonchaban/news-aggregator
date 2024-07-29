package parser

import (
	"fmt"
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/antonchaban/news-aggregator/pkg/parser/html"
	"github.com/antonchaban/news-aggregator/pkg/parser/json"
	"github.com/antonchaban/news-aggregator/pkg/parser/rss"
	"github.com/sirupsen/logrus"
	"net/url"
	"os"
)

// Parser is an interface that defines parsing strategy
type Parser interface {
	ParseFile(f *os.File) ([]model.Article, error)
	ParseFeed(urlPath url.URL) ([]model.Article, error)
}

func ParseArticlesFromFeed(urlPath url.URL) ([]model.Article, error) {
	format, err := DetermineFeedFormat(urlPath)
	parser, err := createParser(format)
	feed, err := parser.ParseFeed(urlPath)
	if err != nil {
		logrus.Errorf("error occurred while parsing feed: %s", err.Error())
	}
	return feed, err
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
			ArticleSelector:     "div.gnt_m.gnt_m_flm > a.gnt_m_flm_a",
			TitleSelector:       "",
			LinkSelector:        "",
			DescriptionSelector: "data-c-br",
			PubDateSelector:     "div.gnt_m_flm_sbt",
			Source:              "USA TODAY",
			DateAttribute:       "data-c-dt",
			TimeFormat: []string{
				"3:04 p.m. ET January 2, 2006",
				"2006-01-02 15:04",
				"Jan 02, 2006",
			},
		}
		return html.NewHtmlParser(config), nil
	default:
		return nil, fmt.Errorf("unsupported file format: %s", format)
	}
}
