package json

import (
	"encoding/json"
	"github.com/antonchaban/news-aggregator/pkg/model"
	"io"
	"net/url"
	"os"
	"time"
)

// Parser is a struct that implements the Parser interface
type Parser struct {
}

func (j *Parser) ParseFeed(url url.URL) ([]model.Article, error) {
	//TODO implement me
	panic("implement me")
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
	bytes, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var feed Feed
	err = json.Unmarshal(bytes, &feed)
	if err != nil {
		return nil, err
	}

	articles := make([]model.Article, 0)
	for _, item := range feed.Articles {
		article := model.Article{
			Title:       item.Title,
			Link:        item.URL,
			Description: item.Description,
			Source:      model.Source{Name: item.Source.Name}, //todo think about how to add links too
			PubDate:     item.PublishedAt,
		}
		articles = append(articles, article)
	}

	return articles, nil
}
