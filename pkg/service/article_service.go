package service

import (
	"errors"
	"github.com/antonchaban/news-aggregator/pkg/filter"
	"github.com/antonchaban/news-aggregator/pkg/handler/cli"
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/antonchaban/news-aggregator/pkg/parser"
	"log"
	"os"
	"path/filepath"
)

//go:generate mockgen -destination=mocks/mock_article.go -package=mocks news-aggregator/pkg/storage ArticleStorage

// ArticleStorage is an interface that defines the methods for interacting with the article storage.
//
// Key Responsibilities:
// 1. Retrieve all stored articles.
// 2. Save a single article to the storage.
// 3. Save multiple articles to the storage.
// 4. Delete an article from the storage by its ID.
// 5. Retrieve articles that match a given keyword.
// 6. Retrieve articles from a specific source.
// 7. Retrieve articles within a specified date range.
//
// Expected Behaviors or Guarantees:
// - The GetAll method should return all articles available in the storage or an error if the operation fails.
// - The Save method should store the provided article and return the saved article (with updated fields such as ID) or an error if the operation fails.
// - The SaveAll method should store all provided articles and return an error if the operation fails for any reason.
// - The Delete method should remove the article with the specified ID from the storage and return an error if the operation fails.
// - The GetByKeyword method should return all articles that contain the specified keyword in their content or metadata, or an error if the operation fails.
// - The GetBySource method should return all articles from the specified source, or an error if the operation fails.
// - The GetByDateInRange method should return all articles within the specified date range, or an error if the operation fails.
//
// Common Errors or Exceptions and Handling:
// - `error`: This general error can occur in any of the methods. It should be handled by logging the error and returning an appropriate message to the user or retrying the operation if possible.
// - Data validation errors: Validate input data before attempting to save or retrieve articles. For example, check if the `id` in Delete method is a valid positive integer.
//
// Known Limitations or Restrictions:
// - The Delete method does not specify what happens if the article does not exist. It should be clarified whether it returns an error or silently succeeds.
// - The methods do not define any constraints on the size or format of the articles being saved.
//
// Usage Guidelines or Best Practices:
// - Use transactions where necessary to ensure data consistency, for example when saving multiple articles with SaveAll.
// - Ensure proper error handling and logging to facilitate debugging and monitoring of storage operations.
// - Validate input data thoroughly before performing any operations to prevent injection attacks or corrupt data.
type ArticleStorage interface {
	GetAll() ([]model.Article, error)
	Save(article model.Article) (model.Article, error)
	SaveAll(articles []model.Article) error
	Delete(id int) error
}

type articleService struct {
	articleStorage ArticleStorage
}

func New(articleRepo ArticleStorage) cli.ArticleService {
	return &articleService{articleStorage: articleRepo}
}

func (a *articleService) LoadDataFromFiles() error {
	files, err := getFilesInDir()
	if err != nil {
		return err
	}
	var articles []model.Article
	for _, file := range files {
		parsedArticles, err := parser.ParseArticlesFromFile(file)
		if err != nil {
			return err
		}
		articles = append(articles, parsedArticles...)

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
