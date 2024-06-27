package json

import (
	"encoding/json"
	"errors"
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/sirupsen/logrus"
	"io"
	"net/url"
	"os"
	"time"
)

// Parser is a struct that implements the Parser interface
type Parser struct{}

func (j *Parser) ParseFeed(url url.URL) ([]model.Article, error) {
	return nil, errors.New("not implemented")
}

// Feed is a struct that represents the JSON feed
type Feed struct {
	Status       string `json:"status"`
	TotalResults int    `json:"totalResults"`
	Articles     []struct {
		Source struct {
			Name string `json:"name"`
		} `json:"source"`
		Author      string    `json:"author"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		URL         string    `json:"url"`
		PublishedAt time.Time `json:"publishedAt"`
	} `json:"articles"`
}

// ParseFile parses the given file and returns a slice of articles.
func (j *Parser) ParseFile(f *os.File) ([]model.Article, error) {
	logrus.WithField("event_id", "parse_json_file_start").Infof("Starting to parse file: %s", f.Name())

	bytes, err := io.ReadAll(f)
	if err != nil {
		logrus.WithField("event_id", "read_file_error").Errorf("Error reading file: %s", err.Error())
		return nil, err
	}

	var feed Feed
	err = json.Unmarshal(bytes, &feed)
	if err != nil {
		logrus.WithField("event_id", "json_unmarshal_error").Errorf("Error unmarshalling JSON: %s", err.Error())
		return nil, err
	}

	articles := make([]model.Article, 0)
	for _, item := range feed.Articles {
		article := model.Article{
			Title:       item.Title,
			Link:        item.URL,
			Description: item.Description,
			Source:      model.Source{Name: item.Source.Name},
			PubDate:     item.PublishedAt,
		}
		articles = append(articles, article)
	}

	logrus.WithField("event_id", "parse_json_file_success").Infof("File parsing completed, found %d articles", len(articles))
	return articles, nil
}
