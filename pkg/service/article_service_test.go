package service

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"news-aggregator/pkg/model"
	"news-aggregator/pkg/storage"
	"news-aggregator/pkg/storage/mocks"
	"testing"
	"time"
)

func TestArticleService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStorage := mocks.NewMockArticleStorage(ctrl)
	article := model.Article{
		Title:       "Test Title",
		Description: "Test Description",
		Link:        "http://test.com",
		Source:      "Test Source",
	}
	mockStorage.EXPECT().Create(article).Return(model.Article{
		Id:          1,
		Title:       "Test Title",
		Description: "Test Description",
		Link:        "http://test.com",
		Source:      "Test Source",
	}, nil)
	a := &ArticleService{
		articleStorage: mockStorage,
	}
	createdArticle, err := a.Create(article)
	if err != nil {
		t.Errorf("Create() error = %v", err)
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
			a := &ArticleService{
				articleStorage: f.articleStorage,
			}
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
						Source:      "Test Source",
					},
					{
						Id:          2,
						Title:       "Test Title 2",
						Description: "Test Description 2",
						Link:        "http://test2.com",
						Source:      "Test Source 2",
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
			a := &ArticleService{
				articleStorage: f.articleStorage,
			}
			articles, _ := a.GetAll()
			assert.Equal(t, 2, len(articles))
		})
	}
}

func TestArticleService_GetByDateInRange(t *testing.T) {
	type fields struct {
		articleStorage *mocks.MockArticleStorage
	}
	type args struct {
		startDate string
		endDate   string
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		want    []model.Article
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				startDate: "2021-01-01",
				endDate:   "2021-01-03",
			},
			prepare: func(f *fields) {
				f.articleStorage.EXPECT().GetByDateInRange(gomock.Any(), gomock.Any()).Return([]model.Article{
					{
						Id:          1,
						Title:       "Test Title",
						Description: "Test Description",
						Link:        "http://test.com",
						Source:      "Test Source",
						PubDate:     time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
					},
					{
						Id:          2,
						Title:       "Test Title 2",
						Description: "Test Description 2",
						Link:        "http://test2.com",
						Source:      "Test Source 2",
						PubDate:     time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
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
			a := &ArticleService{
				articleStorage: f.articleStorage,
			}
			articles, _ := a.GetByDateInRange(tt.args.startDate, tt.args.endDate)
			assert.Equal(t, 2, len(articles))
		})
	}
}

func TestArticleService_GetByKeyword(t *testing.T) {
	type fields struct {
		articleStorage *mocks.MockArticleStorage
	}
	type args struct {
		keyword string
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		want    []model.Article
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				keyword: "test",
			},
			prepare: func(f *fields) {
				f.articleStorage.EXPECT().GetByKeyword("test").Return([]model.Article{
					{
						Id:          1,
						Title:       "Test Title",
						Description: "Test Description",
						Link:        "http://test.com",
						Source:      "Test Source",
					},
					{
						Id:          2,
						Title:       "Test Title 2",
						Description: "Test Description 2",
						Link:        "http://test2.com",
						Source:      "Test Source 2",
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
			a := &ArticleService{
				articleStorage: f.articleStorage,
			}
			articles, _ := a.GetByKeyword(tt.args.keyword)
			assert.Equal(t, 2, len(articles))
		})
	}
}

func TestArticleService_GetBySource(t *testing.T) {
	type fields struct {
		articleStorage *mocks.MockArticleStorage
	}
	type args struct {
		source string
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		want    []model.Article
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				source: "bbc",
			},
			prepare: func(f *fields) {
				f.articleStorage.EXPECT().GetBySource("BBC News").Return([]model.Article{
					{
						Id:          1,
						Title:       "Test Title",
						Description: "Test Description",
						Link:        "http://test.com",
						Source:      "BBC News",
					},
					{
						Id:          2,
						Title:       "Test Title 2",
						Description: "Test Description 2",
						Link:        "http://test2.com",
						Source:      "BBC News",
					},
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "wrong source",
			args: args{
				source: "test",
			},
			prepare: func(f *fields) {

			},
			wantErr: true,
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
			a := &ArticleService{
				articleStorage: f.articleStorage,
			}
			articles, err := a.GetBySource(tt.args.source)
			if tt.wantErr {
				assert.Equal(t, "source not found", err.Error())
			} else {
				assert.Equal(t, 2, len(articles))
			}
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
						Source:      "Test Source",
					},
					{
						Title:       "Test Title 2",
						Description: "Test Description 2",
						Link:        "http://test2.com",
						Source:      "Test Source 2",
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

			a := &ArticleService{
				articleStorage: f.articleStorage,
			}
			err := a.SaveAll(tt.args.articles)

			if (err != nil) && tt.wantErr {
				assert.Equal(t, "save failed", err.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name         string
		articleRepo  storage.ArticleStorage
		expectedRepo *ArticleService
	}{
		{
			name:         "Valid ArticleRepo",
			articleRepo:  mocks.NewMockArticleStorage(gomock.NewController(t)),
			expectedRepo: &ArticleService{articleStorage: mocks.NewMockArticleStorage(gomock.NewController(t))},
		},
		{
			name:         "Nil ArticleRepo",
			articleRepo:  nil,
			expectedRepo: &ArticleService{articleStorage: nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualRepo := New(tt.articleRepo)
			assert.Equal(t, tt.expectedRepo, actualRepo)
		})
	}
}
