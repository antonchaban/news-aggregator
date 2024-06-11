package cli

import (
	"github.com/golang/mock/gomock"
	"news-aggregator/pkg/filter"
	"news-aggregator/pkg/model"
	"news-aggregator/pkg/service/mocks"
	"testing"
	"time"
)

func TestHandler_filterArticles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockArticleService := mocks.NewMockArticleService(ctrl)
	handler := cliHandler{
		service: mockArticleService,
	}

	pubDate1 := time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC)
	pubDate2 := time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC)
	articles := []model.Article{
		{Id: 1, Source: "abcnews", Title: "Title 1", Description: "Description 1", Link: "http://link1.com", PubDate: pubDate1},
		{Id: 2, Source: "bbc", Title: "Title 2", Description: "Description 2", Link: "http://link2.com", PubDate: pubDate2},
	}

	mockArticleService.EXPECT().GetByFilter(filter.Filters{Source: "abcnews"}).Return([]model.Article{articles[0]}, nil).Times(1)

	f := filter.Filters{
		Source: "abcnews",
	}

	filtered, err := handler.filterArticles(f)
	if err != nil {
		t.Errorf("Expected to filter articles, but got an error")
	}
	if len(filtered) != 1 {
		t.Errorf("Expected 1 article but got %d", len(filtered))
	}
}

func TestHandler_loadData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockArticleService := mocks.NewMockArticleService(ctrl)
	handler := cliHandler{
		service: mockArticleService,
	}

	mockArticleService.EXPECT().SaveAll(gomock.Any()).Return(nil).Times(1)
	err := handler.loadData()

	if err != nil {
		t.Errorf("Expected to load articles, but got none")
	}
}
