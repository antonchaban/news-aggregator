package service

import (
	"errors"
	"github.com/antonchaban/news-aggregator/pkg/filter"
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/antonchaban/news-aggregator/pkg/storage"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -destination=../service/mocks/mock_article_service.go -package=mocks news-aggregator/pkg/service ArticleService

// ArticleService is an interface that defines the methods for interacting with the article storage.
type ArticleService interface {
	GetAll() ([]model.Article, error)
	Create(article model.Article) (model.Article, error)
	Delete(id int) error
	SaveAll(articles []model.Article) error
	GetByFilter(f filter.Filters) ([]model.Article, error)
}

type articleService struct {
	articleStorage storage.ArticleStorage
}

func New(articleRepo storage.ArticleStorage) ArticleService {
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

	// Fetch all articles initially
	articles, err := a.GetAll()
	if err != nil {
		logrus.WithField("event_id", "get_all_articles_error").Error("Error fetching all articles", err)
		return nil, err
	}
	logrus.WithField("event_id", "all_articles_fetched").Info("All articles fetched successfully")

	// Create filter handlers
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
