package service

import (
	"github.com/mmcdole/gofeed"
	"os"
)

type Parser interface {
	ParseFile(filePath string) (gofeed.Feed, error)
}

type RssParser struct {
}

func (p *RssParser) ParseFile(filePath string) (gofeed.Feed, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return gofeed.Feed{}, err
	}
	defer file.Close()

	parser := gofeed.NewParser()
	feed, _ := parser.Parse(file)
	return *feed, nil
}
