package service

import (
	"errors"
	"github.com/antonchaban/news-aggregator/pkg/filter"
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/antonchaban/news-aggregator/pkg/parser"
	"github.com/antonchaban/news-aggregator/pkg/storage"
	"log"
	"net/url"
	"os"
	"path/filepath"
)

//go:generate mockgen -destination=../service/mocks/mock_article_service.go -package=mocks news-aggregator/pkg/service ArticleService

// ArticleService is an interface that defines the methods for interacting with the article storage.
type ArticleService interface {
	GetAll() ([]model.Article, error)
	Create(article model.Article) (model.Article, error)
	Delete(id int) error
	SaveAll(articles []model.Article) error
	GetByFilter(f filter.Filters) ([]model.Article, error)
	LoadDataFromFiles() ([]model.Article, error)
	LoadFromFeed(urlPath string) ([]model.Article, error)
}

type articleService struct {
	articleStorage storage.ArticleStorage
}

func New(articleRepo storage.ArticleStorage) ArticleService {
	return &articleService{articleStorage: articleRepo}
}

func (a *articleService) LoadFromFeed(urlPath string) ([]model.Article, error) {
	urlParsed, err := url.Parse(urlPath)
	if err != nil {
		return nil, err
	}
	// Get all articles from the feed
	articles, err := parser.ParseArticlesFromFeed(*urlParsed)
	if err != nil {
		return nil, err
	}

	// Save all articles to the database
	err = a.SaveAll(articles)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func (a *articleService) LoadDataFromFiles() ([]model.Article, error) {
	files, err := getFilesInDir()
	if err != nil {
		return nil, err
	}
	var articles []model.Article
	for _, file := range files {
		parsedArticles, err := parser.ParseArticlesFromFile(file)
		if err != nil {
			return nil, err
		}
		articles = append(articles, parsedArticles...)

	}

	return articles, nil
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

func getFilesInDir() ([]string, error) {
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		return nil, errors.New("environment variable DATA_DIR not set")
	}

	// Get all files in the data directory
	files, err := filepath.Glob(filepath.Join(dataDir, "*"))
	if err != nil {
		log.Fatalf("Error reading files from directory: %v", err)
		return nil, err
	}
	return files, nil
}
