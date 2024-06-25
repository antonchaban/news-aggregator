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

type Loader interface {
	LoadAllFromFile() ([]model.Article, error)
}

type newsLoader struct {
	srcService service.SourceService
}

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

func NewLoader(srcSvc service.SourceService) Loader {
	return &newsLoader{
		srcService: srcSvc,
	}
}
