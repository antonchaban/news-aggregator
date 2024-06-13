package cli

import (
	"news-aggregator/pkg/filter"
	"news-aggregator/pkg/service/mocks"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"news-aggregator/pkg/model"
)

func TestInitCommands_Help(t *testing.T) {
	os.Args = []string{"cmd", "-help"}

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockArticleService := mocks.NewMockArticleService(ctrl)
	NewHandler(mockArticleService)

	w.Close()
	os.Stdout = old

	var buf [1024]byte
	n, _ := r.Read(buf[:])
	output := string(buf[:n])

	assert.Contains(t, output, "Usage of")
	assert.Contains(t, output, "-help")
	assert.Contains(t, output, "-sources string")
	assert.Contains(t, output, "-keywords string")
	assert.Contains(t, output, "-date-start string")
	assert.Contains(t, output, "-date-end string")
}

func TestExecute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockArticleService := mocks.NewMockArticleService(ctrl)

	pubDate1 := time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC)
	pubDate2 := time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC)
	articles := []model.Article{
		{Id: 1, Source: "abcnews", Title: "Title 1", Description: "Description 1", Link: "http://link1.com", PubDate: pubDate1},
		{Id: 2, Source: "bbc", Title: "Title 2", Description: "Description 2", Link: "http://link2.com", PubDate: pubDate2},
	}

	mockArticleService.EXPECT().GetByFilter(filter.Filters{Source: "abcnews", Keyword: "test", StartDate: "2023-01-01", EndDate: "2023-12-31"}).Return([]model.Article{articles[0]}, nil).Times(1)

	files, err := getFilesInDir()
	mockArticleService.EXPECT().LoadDataFromFiles(files).Return(nil).Times(1)
	handler := cliHandler{
		service: mockArticleService,
	}

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f := filter.Filters{
		Source:    "abcnews",
		Keyword:   "test",
		StartDate: "2023-01-01",
		EndDate:   "2023-12-31",
	}

	err = handler.execute(f, "ASC")
	if err != nil {
		return
	}

	w.Close()
	os.Stdout = old

	var buf [1024]byte
	n, _ := r.Read(buf[:])
	output := string(buf[:n])
	assert.Contains(t, output, "Title: Title 1")
	assert.NotContains(t, output, "Title: Title 2")
}
