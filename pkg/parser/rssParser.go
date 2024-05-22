package parser

import (
	"github.com/mmcdole/gofeed"
	"news-aggregator/pkg/model"
	"os"
)

type RssParser struct{}

func (r *RssParser) parseFile(f *os.File) ([]model.Article, error) {
	parser := gofeed.NewParser()
	feed, _ := parser.Parse(f)
	articles := make([]model.Article, 0)
	for _, item := range feed.Items {
		date := item.PublishedParsed
		article := model.Article{
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Description,
			Source:      feed.Title,
			PubDate:     *date,
		}
		articles = append(articles, article)
	}
	return articles, nil
}
