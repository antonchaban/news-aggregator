package backuper

import (
	"encoding/json"
	"fmt"
	"github.com/antonchaban/news-aggregator/pkg/model"
	"os"
	"path/filepath"
)

// Saver is an interface for saving articles to a file.
type Saver interface {
	SaveAllToFile() error
}

// newsSaver is an implementation of the Saver interface.
type newsSaver struct {
	articles []model.Article
}

// NewSaver creates a new Saver instance.
func NewSaver(articles []model.Article) Saver {
	return &newsSaver{
		articles: articles,
	}
}

// SaveAllToFile saves all articles to a JSON file.
// It marshals the articles into JSON and writes them to the file specified by
// the SAVES_DIR environment variable.
func (n newsSaver) SaveAllToFile() error {
	// Convert articles to JSON
	data, err := json.MarshalIndent(n.articles, "", "  ")
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(os.Getenv("SAVES_DIR"), "articles.json"))
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Println(file.Name())
	// Write the JSON data to the file
	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}
