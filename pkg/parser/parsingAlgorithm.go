package parser

import (
	"news-aggregator/pkg/model"
	"os"
)

// ParsingAlgorithm is an interface that defines parsing strategy
type ParsingAlgorithm interface {
	parseFile(f *os.File) ([]model.Article, error)
}
