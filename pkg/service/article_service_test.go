package service

import (
	"github.com/antonchaban/news-aggregator/pkg/filter"
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/antonchaban/news-aggregator/pkg/storage/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestArticleService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStorage := mocks.NewMockArticleStorage(ctrl)
	article := model.Article{
		Title:       "Test Title",
		Description: "Test Description",
		Link:        "http://test.com",
		Source:      model.Source{Id: 1, Link: "http://test.com", Name: "Test Source"},
	}
	mockStorage.EXPECT().Save(article).Return(model.Article{
		Id:          1,
		Title:       "Test Title",
		Description: "Test Description",
		Link:        "http://test.com",
		Source:      model.Source{Id: 1, Link: "http://test.com", Name: "Test Source"},
	}, nil)
	a := New(mockStorage)
	createdArticle, err := a.Create(article)
	if err != nil {
		t.Errorf("Save() error = %v", err)
		return
	}
	assert.Equal(t, 1, createdArticle.Id)
}

func TestArticleService_Delete(t *testing.T) {
	type fields struct {
		articleStorage *mocks.MockArticleStorage
	}
	type args struct {
		id int
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				id: 1,
			},
			prepare: func(f *fields) {
				f.articleStorage.EXPECT().Delete(1).Return(nil)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				articleStorage: mocks.NewMockArticleStorage(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}
			a := New(f.articleStorage)
			if err := a.Delete(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestArticleService_GetAll(t *testing.T) {
	type fields struct {
		articleStorage *mocks.MockArticleStorage
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		wantErr bool
	}{
		{
			name: "success",
			prepare: func(f *fields) {
				f.articleStorage.EXPECT().GetAll().Return([]model.Article{
					{
						Id:          1,
						Title:       "Test Title",
						Description: "Test Description",
						Link:        "http://test.com",
						Source:      model.Source{Id: 1, Link: "http://test.com", Name: "Test Source"},
					},
					{
						Id:          2,
						Title:       "Test Title 2",
						Description: "Test Description 2",
						Link:        "http://test2.com",
						Source:      model.Source{Id: 2, Link: "http://test2.com", Name: "Test Source 2"},
					},
				}, nil)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				articleStorage: mocks.NewMockArticleStorage(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}
			a := New(f.articleStorage)
			articles, _ := a.GetAll()
			assert.Equal(t, 2, len(articles))
		})
	}
}

func TestArticleService_SaveAll(t *testing.T) {
	type fields struct {
		articleStorage *mocks.MockArticleStorage
	}
	type args struct {
		articles []model.Article
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		wantErr bool
	}{
		{
			name: "success",
			prepare: func(f *fields) {
				f.articleStorage.EXPECT().SaveAll(gomock.Any()).Return(nil)
			},
			args: args{
				articles: []model.Article{
					{
						Title:       "Test Title",
						Description: "Test Description",
						Link:        "http://test.com",
						Source:      model.Source{Id: 1, Link: "http://test.com", Name: "Test Source"},
					},
					{
						Title:       "Test Title 2",
						Description: "Test Description 2",
						Link:        "http://test2.com",
						Source:      model.Source{Id: 2, Link: "http://test2.com", Name: "Test Source 2"},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				articleStorage: mocks.NewMockArticleStorage(ctrl),
			}

			if tt.prepare != nil {
				tt.prepare(&f)
			}

			a := New(f.articleStorage)
			err := a.SaveAll(tt.args.articles)

			if (err != nil) && tt.wantErr {
				assert.Equal(t, "save failed", err.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestArticleService_GetByFilter(t *testing.T) {
	type fields struct {
		articleStorage *mocks.MockArticleStorage
	}
	type args struct {
		f filter.Filters
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		wantErr bool
	}{
		{
			name: "success",
			prepare: func(f *fields) {
				f.articleStorage.EXPECT().GetAll().Return([]model.Article{
					{
						Id:          1,
						Title:       "Test Title",
						Description: "Test Description",
						Link:        "http://test.com",
						Source:      model.Source{Id: 1, Link: "http://test.com", Name: "Test Source"},
					},
					{
						Id:          2,
						Title:       "Test Title 2",
						Description: "Test Description 2",
						Link:        "http://test2.com",
						Source:      model.Source{Id: 2, Link: "http://test2.com", Name: "Test Source 2"},
					},
				}, nil)
			},
			args: args{
				f: filter.Filters{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				articleStorage: mocks.NewMockArticleStorage(ctrl),
			}

			if tt.prepare != nil {
				tt.prepare(&f)
			}

			a := New(f.articleStorage)
			articles, err := a.GetByFilter(tt.args.f)

			if (err != nil) && tt.wantErr {
				assert.Equal(t, "failed to fetch articles", err.Error())
			} else {
				assert.Nil(t, err)
				assert.Equal(t, 2, len(articles))
			}
		})
	}

}

func TestNew(t *testing.T) {
	tests := []struct {
		name         string
		articleRepo  ArticleStorage
		expectedRepo *articleService
	}{
		{
			name:         "Valid ArticleRepo",
			articleRepo:  mocks.NewMockArticleStorage(gomock.NewController(t)),
			expectedRepo: &articleService{articleStorage: mocks.NewMockArticleStorage(gomock.NewController(t))},
		},
		{
			name:         "Nil ArticleRepo",
			articleRepo:  nil,
			expectedRepo: &articleService{articleStorage: nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualRepo := New(tt.articleRepo)
			assert.Equal(t, tt.expectedRepo, actualRepo)
		})
	}
}
