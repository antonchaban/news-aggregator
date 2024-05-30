package strategy

import (
	"news-aggregator/pkg/model"
	"os"
)

type Context struct {
	parser ParsingAlgorithm
}

func (pc *Context) SetParser(parser ParsingAlgorithm) {
	pc.parser = parser
}

func (pc *Context) Parse(f *os.File) ([]model.Article, error) {
	return pc.parser.ParseFile(f)
}
