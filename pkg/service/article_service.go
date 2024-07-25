package service

import (
	"errors"
	"github.com/antonchaban/news-aggregator/pkg/filter"
	"github.com/antonchaban/news-aggregator/pkg/handler/web"
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -destination=mocks/mock_article.go -package=mocks github.com/antonchaban/news-aggregator/pkg/storage ArticleStorage

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
}

type articleService struct {
	articleStorage ArticleStorage
}

func New(articleRepo ArticleStorage) web.ArticleService {
	return &articleService{articleStorage: articleRepo}
}

// SaveAll saves multiple articles to the database.
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
	logrus.WithField("event_id", eventGetByFilterStart).Info("Fetching articles with filter")

	// Fetch all articles initially
	articles, err := a.GetAll()
	if err != nil {
		logrus.WithField("event_id", eventGetAllArticlesError).Error("Error fetching all articles", err)
		return nil, err
	}
	logrus.WithField("event_id", eventAllArticlesFetched).Info("All articles fetched successfully")

	// Create filter handlers
	sourceFilter := &filter.SourceFilter{}
	keywordFilter := &filter.KeywordFilter{}
	dateRangeFilter := &filter.DateRangeFilter{}

	logrus.WithField("event_id", eventFiltersCreated).Info("Filter handlers created")

	// Create the chain
	sourceFilter.SetNext(keywordFilter).SetNext(dateRangeFilter)
	logrus.WithField("event_id", eventFiltersChained).Info("Filters chained together")

	// Start filtering
	filteredArticles, err := sourceFilter.Filter(articles, f)
	if err != nil {
		logrus.WithField("event_id", eventFilteringError).Error("Error during filtering", err)
		return nil, err
	}
	logrus.WithField("event_id", eventFilteringComplete).Info("Filtering completed successfully")

	return filteredArticles, nil
}
