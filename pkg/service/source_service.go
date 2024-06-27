package service

import (
	"errors"
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/antonchaban/news-aggregator/pkg/parser"
	"github.com/antonchaban/news-aggregator/pkg/storage"
	"log"
	"net/url"
	"os"
	"path/filepath"
)

//go:generate mockgen -destination=../service/mocks/mock_source_service.go -package=mocks github.com/antonchaban/news-aggregator/pkg/service SourceService

// SourceService represents the service for sources.
type SourceService interface {
	FetchFromAllSources() error
	FetchSourceByID(id int) ([]model.Article, error)
	LoadDataFromFiles() ([]model.Article, error)
	AddSource(source model.Source) (model.Source, error)
	DeleteSource(id int) error
}

// sourceService is the implementation of the SourceService interface.
type sourceService struct {
	articleStorage storage.ArticleStorage
	srcStorage     storage.SourceStorage
}

// NewSourceService creates a new SourceService with the given article and source repositories.
func NewSourceService(articleRepo storage.ArticleStorage, srcRepo storage.SourceStorage) SourceService {
	return &sourceService{articleStorage: articleRepo, srcStorage: srcRepo}
}

// DeleteSource removes the source with the given ID from the database.
func (s *sourceService) DeleteSource(id int) error {
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
			return err
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
