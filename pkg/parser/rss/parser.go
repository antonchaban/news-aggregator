package rss

import (
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/mmcdole/gofeed"
	"net/url"
	"os"
)

// Parser is a struct that implements the ParsingAlgorithm interface
type Parser struct{}

// ParseFeed parses the given URL and returns a slice of articles.
func (r *Parser) ParseFeed(url url.URL) ([]model.Article, error) {
	parser := gofeed.NewParser()
	feed, err := parser.ParseURL(url.String())
	if err != nil {
		return nil, err
	}
	return r.parseFeed(feed, url), nil
}

// ParseFile parses the given file and returns a slice of articles.
func (r *Parser) ParseFile(f *os.File) ([]model.Article, error) {
	parser := gofeed.NewParser()
	feed, err := parser.Parse(f)
	if err != nil {
		return nil, err
	}
	return r.parseFeed(feed, url.URL{}), nil
}

// parseFeed is a helper method that processes the parsed feed and returns articles.
func (r *Parser) parseFeed(feed *gofeed.Feed, feedUrl url.URL) []model.Article {
	articles := make([]model.Article, 0)
	for _, item := range feed.Items {
		article := model.Article{
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Description,
			Source: model.Source{
				Name: feed.Title,
				Link: feedUrl.String()},
		}
		if item.PublishedParsed != nil {
			article.PubDate = *item.PublishedParsed
		}
		articles = append(articles, article)
	}
	return articles
}
