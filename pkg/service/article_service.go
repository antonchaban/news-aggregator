package service

import (
	"errors"
	"github.com/antonchaban/news-aggregator/pkg/filter"
	"github.com/antonchaban/news-aggregator/pkg/handler/web"
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/sirupsen/logrus"
	"os"
)

const (
	eventGetByFilterStart    = "get_by_filter_start"
	eventGetAllArticlesError = "get_all_articles_error"
	eventAllArticlesFetched  = "all_articles_fetched"
	eventFiltersCreated      = "filters_created"
	eventFiltersChained      = "filters_chained"
	eventFilteringError      = "filtering_error"
	eventFilteringComplete   = "filtering_complete"
)

// ArticleStorage is an interface that defines the methods for interacting with the article storage.
type ArticleStorage interface {
	GetAll() ([]model.Article, error)
	Save(article model.Article) (model.Article, error)
	SaveAll(articles []model.Article) error
	Delete(id int) error
	DeleteBySourceID(id int) error
	GetByFilter(query string, args []interface{}) ([]model.Article, error)
}

type articleService struct {
	articleStorage ArticleStorage
}

func New(articleRepo ArticleStorage) web.ArticleService {
	return &articleService{articleStorage: articleRepo}
}

func (a *articleService) SaveAll(articles []model.Article) error {
	err := a.articleStorage.SaveAll(articles)
	if err != nil {
		return errors.New("failed to save articles")
	}
	return err
}

// GetAll returns all articles in the database.
func (a *articleService) GetAll() ([]model.Article, error) {
	return a.articleStorage.GetAll()
}

// Create adds a new article to the database.
func (a *articleService) Create(article model.Article) (model.Article, error) {
	return a.articleStorage.Save(article)
}

// Delete removes the article with the given ID from the database.
func (a *articleService) Delete(id int) error {
	return a.articleStorage.Delete(id)
}

// GetByFilter returns all articles that match the given filters.
func (a *articleService) GetByFilter(f filter.Filters) ([]model.Article, error) {
	logrus.WithField("event_id", "get_by_filter_start").Info("Fetching articles with filter")

	if os.Getenv("STORAGE_TYPE") == "postgres" {
		articles, err := a.getByFilterDB(f)
		if err != nil {
			if err.Error() == "GetByFilter operation is not supported in in-memory storage" {
				logrus.WithField("event_id", "fallback_to_inmemory").Warn("Falling back to in-memory filtering")
				return a.getByFilterInMemory(f)
			}
			return nil, err
		}
		return articles, nil
	}
	return a.getByFilterInMemory(f)
}

func (a *articleService) getByFilterInMemory(f filter.Filters) ([]model.Article, error) {
	articles, err := a.articleStorage.GetAll()
	if err != nil {
		logrus.WithField("event_id", "get_all_articles_error").Error("Error fetching all articles", err)
		return nil, err
	}
	logrus.WithField("event_id", "all_articles_fetched").Info("All articles fetched successfully")
	sourceFilter := &filter.SourceFilter{}
	keywordFilter := &filter.KeywordFilter{}
	dateRangeFilter := &filter.DateRangeFilter{}

	logrus.WithField("event_id", "filters_created").Info("Filter handlers created")

	// Create the chain
	sourceFilter.SetNext(keywordFilter).SetNext(dateRangeFilter)
	logrus.WithField("event_id", "filters_chained").Info("Filters chained together")

	// Start filtering
	filteredArticles, err := sourceFilter.Filter(articles, f)
	if err != nil {
		logrus.WithField("event_id", "filtering_error").Error("Error during filtering", err)
		return nil, err
	}
	logrus.WithField("event_id", "filtering_complete").Info("Filtering completed successfully")

	return filteredArticles, nil
}

func (a *articleService) getByFilterDB(f filter.Filters) ([]model.Article, error) {
	baseQuery := `
		SELECT a.id, a.title, a.description, a.link, a.pub_date,
		       s.id, s.name, s.link
		FROM articles a
		JOIN sources s ON a.source_id = s.id
		WHERE 1=1
	`

	sourceFilter := &filter.SourceFilter{}
	keywordFilter := &filter.KeywordFilter{}
	dateRangeFilter := &filter.DateRangeFilter{}

	logrus.WithField("event_id", "filters_created").Info("Filter handlers created")

	sourceFilter.SetNext(keywordFilter).SetNext(dateRangeFilter)
	logrus.WithField("event_id", "filters_chained").Info("Filters chained together")

	query, args := sourceFilter.BuildFilterQuery(f, baseQuery)

	articles, err := a.articleStorage.GetByFilter(query, args)
	if err != nil {
		logrus.WithField("event_id", "query_error").Error("Error executing query", err)
		return nil, err
	}

	logrus.WithField("event_id", "filtering_complete").Info("Filtering completed successfully")
	return articles, nil
}
