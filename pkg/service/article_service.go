package service

import (
	"errors"
	"news-aggregator/pkg/filter"
	"news-aggregator/pkg/model"
	"news-aggregator/pkg/parser"
	"news-aggregator/pkg/storage"
)

//go:generate mockgen -destination=../service/mocks/mock_article_service.go -package=mocks news-aggregator/pkg/service ArticleService

// ArticleService is an interface that defines the methods for interacting with the article storage.
type ArticleService interface {
	GetAll() ([]model.Article, error)
	Create(article model.Article) (model.Article, error)
	Delete(id int) error
	SaveAll(articles []model.Article) error
	GetByFilter(f filter.Filters) ([]model.Article, error)
	LoadDataFromFiles(files []string) error
}

type articleService struct {
	articleStorage storage.ArticleStorage
}

func New(articleRepo storage.ArticleStorage) ArticleService {
	return &articleService{articleStorage: articleRepo}
}

func (a *articleService) LoadDataFromFiles(files []string) error {
	articles, err := parser.ParseArticlesFromFiles(files)
	if err != nil {
		return errors.New("error parsing articles from files")
	}
	return a.SaveAll(articles)
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
	// Fetch all articles initially
	articles, err := a.GetAll()
	if err != nil {
		return nil, err
	}

	// Create filter handlers
	sourceFilter := &filter.SourceFilter{}
	keywordFilter := &filter.KeywordFilter{}
	dateRangeFilter := &filter.DateRangeFilter{}

	// Create the chain
	sourceFilter.SetNext(keywordFilter).SetNext(dateRangeFilter)

	// Start filtering
	return sourceFilter.Filter(articles, f)
}
