package strategy

import (
	"news-aggregator/pkg/model"
	"os"
)

// ParsingAlgorithm is an interface that defines parsing strategy
type ParsingAlgorithm interface {
	ParseFile(f *os.File) ([]model.Article, error)
}
