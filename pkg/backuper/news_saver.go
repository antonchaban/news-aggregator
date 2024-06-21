package backuper

import (
	"encoding/json"
	"fmt"
	"github.com/antonchaban/news-aggregator/pkg/model"
	"os"
	"path/filepath"
)

type Saver interface {
	SaveAllToFile() error
}

type newsSaver struct {
	articles []model.Article
}

func NewSaver(articles []model.Article) Saver {
	return &newsSaver{
		articles: articles,
	}
}

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
