package parser

import (
	"github.com/mmcdole/gofeed"
	"os"
)

type JsonParser struct {
}

func (j *JsonParser) parseFile(f *os.File) (*gofeed.Feed, error) {
	parser := gofeed.NewParser()
	feed, err := parser.Parse(f)
	if err != nil {
		return nil, err
	}
	return feed, nil
}
