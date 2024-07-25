package rss

import (
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/mmcdole/gofeed"
	"github.com/sirupsen/logrus"
	"net/url"
	"os"
)

const (
	eventParseRssFeedStart         = "parse_rss_feed_start"
	eventParseRssUrlError          = "parse_rss_url_error"
	eventParseRssFeedSuccess       = "parse_rss_feed_success"
	eventParseRssFileStart         = "parse_rss_file_start"
	eventParseRssFileError         = "parse_rss_file_error"
	eventParseRssFileSuccess       = "parse_rss_file_success"
	eventParseRssFeedItems         = "parse_rss_feed_items"
	eventParseRssFeedItemsComplete = "parse_rss_feed_items_complete"
)

// Parser is a struct that implements the ParsingAlgorithm interface
type Parser struct{}

// ParseFeed parses the given URL and returns a slice of articles.
func (r *Parser) ParseFeed(url url.URL) ([]model.Article, error) {
	logrus.WithField("event_id", eventParseRssFeedStart).Infof("Starting to parse feed from URL: %s", url.String())

	parser := gofeed.NewParser()
	feed, err := parser.ParseURL(url.String())
	if err != nil {
		logrus.WithField("event_id", eventParseRssUrlError).Errorf("Error parsing URL: %s", err.Error())
		return nil, err
	}

	logrus.WithField("event_id", eventParseRssFeedSuccess).Infof("Successfully parsed feed from URL: %s", url.String())
	return r.parseFeed(feed, url), nil
}

// ParseFile parses the given file and returns a slice of articles.
func (r *Parser) ParseFile(f *os.File) ([]model.Article, error) {
	logrus.WithField("event_id", eventParseRssFileStart).Infof("Starting to parse file: %s", f.Name())

	parser := gofeed.NewParser()
	feed, err := parser.Parse(f)
	if err != nil {
		logrus.WithField("event_id", eventParseRssFileError).Errorf("Error parsing file: %s", err.Error())
		return nil, err
	}

	logrus.WithField("event_id", eventParseRssFileSuccess).Infof("Successfully parsed file: %s", f.Name())
	return r.parseFeed(feed, url.URL{}), nil
}

// parseFeed is a helper method that processes the parsed feed and returns articles.
func (r *Parser) parseFeed(feed *gofeed.Feed, feedUrl url.URL) []model.Article {
	logrus.WithField("event_id", eventParseRssFeedItems).Info("Processing feed items")
	articles := make([]model.Article, 0)
	for _, item := range feed.Items {
		article := model.Article{
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Description,
			Source: model.Source{
				Name: feed.Title,
				Link: feedUrl.String(),
			},
		}
		if item.PublishedParsed != nil {
			article.PubDate = *item.PublishedParsed
		}
		articles = append(articles, article)
	}
	logrus.WithField("event_id", eventParseRssFeedItemsComplete).Infof("Completed processing feed items, found %d articles", len(articles))
	return articles
}
