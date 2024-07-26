package backuper

import (
	"encoding/json"
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestNewSaver(t *testing.T) {
	articles := []model.Article{
		{Title: "Test Article"},
	}
	sources := []model.Source{
		{Name: "Test Source"},
	}
	saver := NewSaver(articles, sources)
	assert.NotNil(t, saver)
}

func TestSaveAllToFile(t *testing.T) {
	articles := []model.Article{
		{Title: "Test Article"},
	}
	sources := []model.Source{
		{Name: "Test Source"},
	}
	saver := NewSaver(articles, sources)
	os.Setenv("SAVES_DIR", "./testdata")
	defer os.Unsetenv("SAVES_DIR")

	t.Run("should save articles to file", func(t *testing.T) {
		err := saver.SaveAllToFile()
		assert.Nil(t, err)

		filePath := filepath.Join("./testdata", "articles.json")
		fileData, err := os.ReadFile(filePath)
		assert.Nil(t, err)

		var savedArticles []model.Article
		err = json.Unmarshal(fileData, &savedArticles)
		assert.Nil(t, err)
		assert.Equal(t, articles, savedArticles)
	})
}

func TestSaveSrcsToFile(t *testing.T) {
	articles := []model.Article{
		{Title: "Test Article"},
	}
	sources := []model.Source{
		{Name: "Test Source"},
	}
	saver := NewSaver(articles, sources)
	os.Setenv("SAVES_DIR", "./testdata")
	defer os.Unsetenv("SAVES_DIR")

	t.Run("should save sources to file", func(t *testing.T) {
		err := saver.SaveSrcsToFile()
		assert.Nil(t, err)

		filePath := filepath.Join("./testdata", "sources.json")
		fileData, err := os.ReadFile(filePath)
		assert.Nil(t, err)

		var savedSources []model.Source
		err = json.Unmarshal(fileData, &savedSources)
		assert.Nil(t, err)
		assert.Equal(t, sources, savedSources)
	})
}
