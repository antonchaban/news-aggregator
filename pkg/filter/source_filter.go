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

const (
	eventSourceFilterStart       = "source_filter_start"
	eventSourceNotFound          = "source_not_found"
	eventSourceFilteringComplete = "source_filtering_complete"
)

// SourceFilter filters articles based on their source.
type SourceFilter struct {
	next ArticleFilter
}

// SetNext sets the next filter in the chain and returns the filter.
func (h *SourceFilter) SetNext(filter ArticleFilter) ArticleFilter {
	h.next = filter
	return filter
}

// Filter filters articles by their source based on the provided Filters.
func (h *SourceFilter) Filter(articles []model.Article, f Filters) ([]model.Article, error) {
	logrus.WithField("event_id", eventSourceFilterStart).Info("Starting SourceFilter")

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
					logrus.WithField("event_id", eventSourceNotFound).Errorf("Source not found: %s", source)
					return nil, fmt.Errorf("source not found: %s", source)
				}
			}
		}
		articles = filteredArticles
		logrus.WithField("filtered_count", len(filteredArticles)).Info(eventSourceFilteringComplete)
	}

	if h.next != nil {
		return h.next.Filter(articles, f)
	}
	return articles, nil
}

// findOther checks if the source name is not one of the predefined sources
func findOther(sourceName string) bool {
	return sourceName != abcNewsSource && sourceName != bbcNewsSource &&
		sourceName != washingtonTimesSource && sourceName != nbcNewsSource &&
		sourceName != usaTodaySource
}
