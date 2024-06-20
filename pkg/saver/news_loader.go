package saver

import (
	"encoding/json"
	"errors"
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/antonchaban/news-aggregator/pkg/service"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Loader interface {
	LoadAllFromFile() ([]model.Article, error)
}

type newsLoader struct {
	articleService service.ArticleService
}

func (n newsLoader) LoadAllFromFile() ([]model.Article, error) {
	// Construct the file path
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		return nil, errors.New("DATA_DIR environment variable is not set")
	}
	filePath := filepath.Join(dataDir, "articles.json")

	// Check if the file exists
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) || fileInfo.Size() == 0 {
		// If the file does not exist or is empty, call LoadDataFromFiles
		return n.articleService.LoadDataFromFiles()
	} else if err != nil {
		return nil, err
	}

	// Read the file contents
	fileData, err := ioutil.ReadFile(filePath)
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

func NewLoader() Loader {
	return &newsLoader{}
}
