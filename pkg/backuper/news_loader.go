package backuper

import (
	"encoding/json"
	"errors"
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/antonchaban/news-aggregator/pkg/service"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

// Loader is an interface for loading articles from a file.
type Loader interface {
	LoadAllFromFile() ([]model.Article, error)
}

// newsLoader is an implementation of the Loader interface.
type newsLoader struct {
	srcService service.SourceService
}

// NewLoader creates a new Loader instance.
func NewLoader(srcSvc service.SourceService) Loader {
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
