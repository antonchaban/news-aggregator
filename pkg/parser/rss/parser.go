package rss

import (
	"github.com/mmcdole/gofeed"
	"news-aggregator/pkg/model"
	"os"
)

// Parser is a struct that implements the ParsingAlgorithm interface
type Parser struct{}

// ParseFile parses the given file and returns a slice of articles.
func (r *Parser) ParseFile(f *os.File) ([]model.Article, error) {
	parser := gofeed.NewParser()
	feed, err := parser.Parse(f)
	if err != nil {
		return nil, err
	}
	articles := make([]model.Article, 0)
	for _, item := range feed.Items {
		article := model.Article{
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Description,
			Source:      feed.Title,
		}
		if item.PublishedParsed != nil {
			article.PubDate = *item.PublishedParsed
		}
		articles = append(articles, article)
	}
	return articles, nil
}
