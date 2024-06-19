package filter

import (
	"fmt"
	"strings"

	"github.com/antonchaban/news-aggregator/pkg/model"
)

const (
	abcNewsSource         = "ABC News: International"
	bbcNewsSource         = "BBC News"
	washingtonTimesSource = "The Washington Times stories: World"
	nbcNewsSource         = "NBC News"
	usaTodaySource        = "USA TODAY"
)

type SourceFilter struct {
	next ArticleFilter
}

func (h *SourceFilter) SetNext(handler ArticleFilter) ArticleFilter {
	h.next = handler
	return handler
}

func (h *SourceFilter) Filter(articles []model.Article, f Filters) ([]model.Article, error) {
	if f.Source != "" {
		sourceMap := map[string]string{
			"abcnews":         abcNewsSource,
			"bbc":             bbcNewsSource,
			"washingtontimes": washingtonTimesSource,
			"nbc":             nbcNewsSource,
			"usatoday":        usaTodaySource,
		}

		sourceList := strings.Split(f.Source, ",")
		var filteredArticles []model.Article
		for _, article := range articles {
			for _, source := range sourceList {
				if sourceName, ok := sourceMap[source]; ok {
					if article.Source == sourceName {
						filteredArticles = append(filteredArticles, article)
						break
					}
				} else {
					return nil, fmt.Errorf("source not found: %s", source)
				}
			}
		}
		articles = filteredArticles
	}

	if h.next != nil {
		return h.next.Filter(articles, f)
	}
	return articles, nil
}
