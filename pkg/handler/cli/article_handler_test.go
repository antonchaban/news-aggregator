package cli

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"news-aggregator/pkg/model"
	"news-aggregator/pkg/service/mocks"
	"testing"
	"time"
)

func TestHandler_filterArticles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockArticleService := mocks.NewMockArticle(ctrl)
	handler := Handler{
		Service: mockArticleService,
	}

	pubDate1 := time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC)
	pubDate2 := time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC)
	articles := []model.Article{
		{Id: 1, Source: "abcnews", Title: "Title 1", Description: "Description 1", Link: "http://link1.com", PubDate: pubDate1},
		{Id: 2, Source: "bbc", Title: "Title 2", Description: "Description 2", Link: "http://link2.com", PubDate: pubDate2},
	}

	mockArticleService.EXPECT().GetBySource("abcnews").Return([]model.Article{articles[0]}, nil).Times(1)

	filtered := handler.filterArticles("abcnews", "", "", "")
	if len(filtered) != 1 {
		t.Errorf("Expected 1 article but got %d", len(filtered))
	}
}

func TestHandler_loadData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockArticleService := mocks.NewMockArticle(ctrl)
	handler := Handler{
		Service: mockArticleService,
	}

	mockArticleService.EXPECT().SaveAll(gomock.Any()).Return(nil).Times(1)
	err := handler.loadData()

	if err != nil {
		t.Errorf("Expected to load articles, but got none")
	}
}

func Test_intersect(t *testing.T) {
	type args struct {
		a []model.Article
		b []model.Article
	}
	tests := []struct {
		name         string
		args         args
		wantArticles []model.Article
	}{
		{
			name: "No intersection",
			args: args{
				a: []model.Article{{Id: 1}, {Id: 2}},
				b: []model.Article{{Id: 3}, {Id: 4}},
			},
			wantArticles: nil,
		},
		{
			name: "Some intersection",
			args: args{
				a: []model.Article{{Id: 1}, {Id: 2}},
				b: []model.Article{{Id: 2}, {Id: 3}},
			},
			wantArticles: []model.Article{{Id: 2}},
		},
		{
			name: "Complete intersection",
			args: args{
				a: []model.Article{{Id: 1}, {Id: 2}},
				b: []model.Article{{Id: 1}, {Id: 2}},
			},
			wantArticles: []model.Article{{Id: 1}, {Id: 2}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.wantArticles, intersect(tt.args.a, tt.args.b), "intersect(%v, %v)", tt.args.a, tt.args.b)
		})
	}
}
