package parser

import (
	"github.com/mmcdole/gofeed"
	"os"
)

type HtmlParser struct{}

func (h *HtmlParser) parseFile(f *os.File) (*gofeed.Feed, error) {
	parser := gofeed.NewParser()
	feed, _ := parser.Parse(f)
	return feed, nil
} // Todo: Implement the parseFile method for the HtmlParser struct
