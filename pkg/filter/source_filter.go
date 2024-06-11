package filter

import (
	"fmt"
	"strings"

	"news-aggregator/pkg/model"
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
		sourceList := strings.Split(f.Source, ",")
		var filteredArticles []model.Article
		for _, article := range articles {
			for _, source := range sourceList {
				switch source {
				case "abcnews":
					if article.Source == abcNewsSource {
						filteredArticles = append(filteredArticles, article)
					}
				case "bbc":
					if article.Source == bbcNewsSource {
						filteredArticles = append(filteredArticles, article)

					}
				case "washingtontimes":
					if article.Source == washingtonTimesSource {
						filteredArticles = append(filteredArticles, article)
					}
				case "nbc":
					if article.Source == nbcNewsSource {
						filteredArticles = append(filteredArticles, article)
					}
				case "usatoday":
					if article.Source == usaTodaySource {
						filteredArticles = append(filteredArticles, article)
					}
				default:
					return nil, fmt.Errorf("source not found")
				}
			}
			articles = filteredArticles
		}

	}
	if h.next != nil {
		return h.next.Filter(articles, f)
	}
	return articles, nil
}
