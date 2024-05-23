package parser

import (
	"news-aggregator/pkg/model"
	"os"
)

type ParsingAlgorithm interface {
	parseFile(f *os.File) ([]model.Article, error)
}
