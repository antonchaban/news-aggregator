package service

import (
	"errors"
	"github.com/antonchaban/news-aggregator/pkg/handler/web"
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/antonchaban/news-aggregator/pkg/parser"
	"github.com/sirupsen/logrus"
	"log"
	"net/url"
	"os"
	"path/filepath"
)

//go:generate mockgen -destination=mocks/mock_source.go -package=mocks github.com/antonchaban/news-aggregator/pkg/storage SourceStorage

// SourceStorage is an interface that defines the methods for interacting with the source storage.
// Expected Behaviors or Guarantees:
// - The GetAll method should return all sources available in the storage or an error if the operation fails.
//
// - The Save method should store the provided source and return the saved source (with updated fields such as ID) or an error if the operation fails.
//
// - The SaveAll method should store all provided sources and return an error if the operation fails for any reason.
//
// - The Delete method should remove the source with the specified ID from the storage and return an error if the operation fails.
//
// - The GetByID method should return the source with the specified ID or an error if the source does not exist or the operation fails.
//
// Common Errors or Exceptions and Handling:
// - `error`: This general error can occur in any of the methods. It should be handled by logging the error and returning an appropriate message to the user or retrying the operation if possible.
//
// - Data validation errors: Validate input data before attempting to save or retrieve sources. For example, check if the `id` in Delete or GetByID methods is a valid positive integer.
//
// Known Limitations or Restrictions:
// - The Delete method does not specify what happens if the source does not exist. It should be clarified whether it returns an error or silently succeeds.
// - The methods do not define any constraints on the size or format of the sources being saved.
//
// Usage Guidelines or Best Practices:
// - Use transactions where necessary to ensure data consistency, for example, when saving multiple sources with SaveAll.
// - Ensure proper error handling and logging to facilitate debugging and monitoring of storage operations.
// - Validate input data thoroughly before performing any operations to prevent injection attacks or corrupt data.
type SourceStorage interface {
	GetAll() ([]model.Source, error)
	Save(src model.Source) (model.Source, error)
	SaveAll(sources []model.Source) error
	Delete(id int) error
	GetByID(id int) (model.Source, error)
	Update(id int, src model.Source) (model.Source, error)
}

// sourceService is the implementation of the SourceService interface.
type sourceService struct {
	articleStorage ArticleStorage
	srcStorage     SourceStorage
}

// NewSourceService creates a new SourceService with the given article and source repositories.
func NewSourceService(articleRepo ArticleStorage, srcRepo SourceStorage) web.SourceService {
	return &sourceService{articleStorage: articleRepo, srcStorage: srcRepo}
}

// GetAll returns all sources from the database.
func (s *sourceService) GetAll() ([]model.Source, error) {
	return s.srcStorage.GetAll()
}

// UpdateSource updates the source with the given ID in the database.
func (s *sourceService) UpdateSource(id int, source model.Source) (model.Source, error) {
	return s.srcStorage.Update(id, source)
}

// DeleteSource removes the source with the given ID from the database.
func (s *sourceService) DeleteSource(id int) error {
	err := s.articleStorage.DeleteBySourceID(id)
	if err != nil {
		return err
	}
	return s.srcStorage.Delete(id)
}

// AddSource adds a new source to the database.
func (s *sourceService) AddSource(source model.Source) (model.Source, error) {
	save, err := s.srcStorage.Save(source)
	if err != nil {
		return model.Source{}, err
	}
	return save, nil
}

// FetchFromAllSources fetches articles from all sources.
func (s *sourceService) FetchFromAllSources() error {
	allSrcs, err := s.srcStorage.GetAll()
	if err != nil {
		return err
	}

	for _, src := range allSrcs {
		urlParsed, err := url.Parse(src.Link)
		if err != nil {
			return err
		}
		articles, err := parser.ParseArticlesFromFeed(*urlParsed)
		if err != nil {
			return err
		}
		for i := range articles {
			articles[i].Source = src
		}
		err = s.articleStorage.SaveAll(articles)
		if err != nil {
			logrus.Printf("Error saving articles: %v", err)
			continue
		}
	}
	return nil
}

// FetchSourceByID fetches articles from the source with the given ID.
func (s *sourceService) FetchSourceByID(id int) ([]model.Article, error) {
	src, err := s.srcStorage.GetByID(id)
	if err != nil {
		return nil, err
	}
	urlParsed, err := url.Parse(src.Link)
	if err != nil {
		return nil, err
	}
	articles, err := parser.ParseArticlesFromFeed(*urlParsed)
	if err != nil {
		return nil, err
	}
	err = s.articleStorage.SaveAll(articles)
	if err != nil {
		return nil, err
	}
	return articles, nil
}

// LoadDataFromFiles loads articles from files.
func (s *sourceService) LoadDataFromFiles() ([]model.Article, error) {
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
