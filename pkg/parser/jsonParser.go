package parser

import (
	"encoding/json"
	"io"
	"news-aggregator/pkg/model"
	"os"
)

type JsonParser struct {
}

type JsonFeed struct {
	Status       string `json:"status"`
	TotalResults int    `json:"totalResults"`
	Articles     []struct {
		Source struct {
			Name string `json:"name"`
		} `json:"source"`
		Author      string `json:"author"`
		Title       string `json:"title"`
		Description string `json:"description"`
		URL         string `json:"url"`
		PublishedAt string `json:"publishedAt"`
	} `json:"articles"`
}

func (j *JsonParser) parseFile(f *os.File) ([]model.Article, error) {
	bytes, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var feed JsonFeed
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
			Source:      item.Source.Name,
			PubDate:     item.PublishedAt,
		}
		articles = append(articles, article)
	}

	return articles, nil
}
