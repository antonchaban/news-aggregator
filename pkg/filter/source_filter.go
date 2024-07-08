package filter

import (
	"fmt"
	"strings"

	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/sirupsen/logrus"
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
	logrus.WithField("event_id", "source_filter_start").Info("Starting SourceFilter")

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
				if source == "other" {
					if findOther(article.Source.Name) {
						filteredArticles = append(filteredArticles, article)
						break
					}
				} else if sourceName, ok := sourceMap[source]; ok {
					if article.Source.Name == sourceName {
						filteredArticles = append(filteredArticles, article)
						break
					}
				} else {
					logrus.WithField("event_id", "source_not_found").Errorf("Source not found: %s", source)
					return nil, fmt.Errorf("source not found: %s", source)
				}
			}
		}
		articles = filteredArticles
		logrus.WithField("filtered_count", len(filteredArticles)).Info("Source filtering complete")
	}

	if h.next != nil {
		return h.next.Filter(articles, f)
	}
	return articles, nil
}

func (h *SourceFilter) BuildFilterQuery(f Filters, query string) (string, []interface{}) {
	if f.Source != "" {
		sourceMap := map[string]string{
			"abcnews":         abcNewsSource,
			"bbc":             bbcNewsSource,
			"washingtontimes": washingtonTimesSource,
			"nbc":             nbcNewsSource,
			"usatoday":        usaTodaySource,
		}

		sourceList := strings.Split(f.Source, ",")
		var sourceConditions []string
		for _, source := range sourceList {
			if source == "other" {
				condition := fmt.Sprintf("s.name NOT IN ('%s', '%s', '%s', '%s', '%s')",
					abcNewsSource, bbcNewsSource, washingtonTimesSource, nbcNewsSource, usaTodaySource)
				sourceConditions = append(sourceConditions, condition)
			} else if sourceName, ok := sourceMap[source]; ok {
				condition := fmt.Sprintf("s.name = '%s'", sourceName)
				sourceConditions = append(sourceConditions, condition)
			} else {
				logrus.WithField("event_id", "source_not_found").Errorf("Source not found: %s", source)
				return query, nil
			}
		}
		if len(sourceConditions) > 0 {
			query += " AND (" + strings.Join(sourceConditions, " OR ") + ")"
		}
	}
	if h.next != nil {
		return h.next.BuildFilterQuery(f, query)
	}
	return query, nil
}

// findOther checks if the source name is not one of the predefined sources
func findOther(sourceName string) bool {
	return sourceName != abcNewsSource && sourceName != bbcNewsSource &&
		sourceName != washingtonTimesSource && sourceName != nbcNewsSource &&
		sourceName != usaTodaySource
}
