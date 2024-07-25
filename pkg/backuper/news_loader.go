package backuper

import (
	"encoding/json"
	"errors"
	"github.com/antonchaban/news-aggregator/pkg/handler/web"
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

// Loader is an interface for loading articles and sources from a file.
type Loader interface {
	LoadAllFromFile() ([]model.Article, error)
	LoadSrcsFromFile() ([]model.Source, error)
}

// newsLoader is an implementation of the Loader interface.
type newsLoader struct {
	srcService web.SourceService
}

// NewLoader creates a new Loader instance.
func NewLoader(srcSvc web.SourceService) Loader {
	return &newsLoader{
		srcService: srcSvc,
	}
}

// LoadAllFromFile loads all articles from a JSON file.
// It reads the file specified by the SAVES_DIR environment variable and unmarshals
// its contents into a slice of Article structs.
func (n *newsLoader) LoadAllFromFile() ([]model.Article, error) {
	// Construct the file path
	dataDir := os.Getenv("SAVES_DIR")
	if dataDir == "" {
		return nil, errors.New("SAVES_DIR environment variable is not set")
	}
	filePath := filepath.Join(dataDir, "articles.json")

	// Check if the file exists
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) || fileInfo.Size() == 0 {
		logrus.Warn("No articles found in the backup file")
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	// Read the file contents
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON data
	var articles []model.Article
	err = json.Unmarshal(fileData, &articles)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

// LoadSrcsFromFile loads all sources from a JSON file.
// It reads the file specified by the SAVES_DIR environment variable and unmarshals
// its contents into a slice of Source structs.
func (n *newsLoader) LoadSrcsFromFile() ([]model.Source, error) {
	// Construct the file path
	dataDir := os.Getenv("SAVES_DIR")
	if dataDir == "" {
		return nil, errors.New("SAVES_DIR environment variable is not set")
	}
	filePath := filepath.Join(dataDir, "sources.json")

	// Check if the file exists
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) || fileInfo.Size() == 0 {
		logrus.Warn("No sources found in the backup file")
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	// Read the file contents
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON data
	var sources []model.Source
	err = json.Unmarshal(fileData, &sources)
	if err != nil {
		return nil, err
	}

	return sources, nil
}
