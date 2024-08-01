package backuper

/*import (
	"encoding/json"
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/antonchaban/news-aggregator/pkg/service/mocks"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestNewLoader(t *testing.T) {
	mockSrcService := new(mocks.MockSourceService)
	loader := NewLoader(mockSrcService)
	assert.NotNil(t, loader)
}

func TestLoadAllFromFile(t *testing.T) {
	mockSrcService := new(mocks.MockSourceService)
	loader := NewLoader(mockSrcService)
	os.Setenv("SAVES_DIR", "./testdata")
	defer os.Unsetenv("SAVES_DIR")

	t.Run("should return error when SAVES_DIR is not set", func(t *testing.T) {
		os.Unsetenv("SAVES_DIR")
		articles, err := loader.LoadAllFromFile()
		assert.Nil(t, articles)
		assert.EqualError(t, err, "SAVES_DIR environment variable is not set")
		os.Setenv("SAVES_DIR", "./testdata")
	})

	t.Run("should return error when json is invalid", func(t *testing.T) {
		os.Setenv("SAVES_DIR", "./testdata")
		invalidFilePath := filepath.Join("./testdata", "articles.json")
		os.WriteFile(invalidFilePath, []byte(`invalid json`), 0644)

		articles, err := loader.LoadAllFromFile()
		assert.Nil(t, articles)
		assert.NotNil(t, err)
	})

	t.Run("should load articles successfully", func(t *testing.T) {
		os.Setenv("SAVES_DIR", "./testdata")
		validArticles := []model.Article{
			{Title: "Test Article"},
		}
		fileData, _ := json.Marshal(validArticles)
		validFilePath := filepath.Join("./testdata", "articles.json")
		os.WriteFile(validFilePath, fileData, 0644)

		articles, err := loader.LoadAllFromFile()
		assert.NotNil(t, articles)
		assert.Nil(t, err)
		assert.Equal(t, validArticles, articles)
	})
}

func TestLoadSrcsFromFile(t *testing.T) {
	mockSrcService := new(mocks.MockSourceService)
	loader := NewLoader(mockSrcService)
	os.Setenv("SAVES_DIR", "./testdata")
	defer os.Unsetenv("SAVES_DIR")

	t.Run("should return error when SAVES_DIR is not set", func(t *testing.T) {
		os.Unsetenv("SAVES_DIR")
		sources, err := loader.LoadSrcsFromFile()
		assert.Nil(t, sources)
		assert.EqualError(t, err, "SAVES_DIR environment variable is not set")
		os.Setenv("SAVES_DIR", "./testdata")
	})

	t.Run("should return error when json is invalid", func(t *testing.T) {
		os.Setenv("SAVES_DIR", "./testdata")
		invalidFilePath := filepath.Join("./testdata", "sources.json")
		os.WriteFile(invalidFilePath, []byte(`invalid json`), 0644)

		sources, err := loader.LoadSrcsFromFile()
		assert.Nil(t, sources)
		assert.NotNil(t, err)
	})

	t.Run("should load sources successfully", func(t *testing.T) {
		os.Setenv("SAVES_DIR", "./testdata")
		validSources := []model.Source{
			{Name: "Test Source"},
		}
		fileData, _ := json.Marshal(validSources)
		validFilePath := filepath.Join("./testdata", "sources.json")
		os.WriteFile(validFilePath, fileData, 0644)

		sources, err := loader.LoadSrcsFromFile()
		assert.NotNil(t, sources)
		assert.Nil(t, err)
		assert.Equal(t, validSources, sources)
	})
}
*/
