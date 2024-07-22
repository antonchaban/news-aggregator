package service

import (
	"errors"
	"fmt"
	"github.com/antonchaban/news-aggregator/pkg/parser"
	"github.com/antonchaban/news-aggregator/pkg/storage"
	"github.com/stretchr/testify/require"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/antonchaban/news-aggregator/pkg/storage/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewSourceService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockArticleStorage := mocks.NewMockArticleStorage(ctrl)
	mockSourceStorage := mocks.NewMockSourceStorage(ctrl)

	tests := []struct {
		name string
		args struct {
			articleRepo storage.ArticleStorage
			srcRepo     storage.SourceStorage
		}
		want SourceService
	}{
		{
			name: "initialize source service",
			args: struct {
				articleRepo storage.ArticleStorage
				srcRepo     storage.SourceStorage
			}{
				articleRepo: mockArticleStorage,
				srcRepo:     mockSourceStorage,
			},
			want: &sourceService{articleStorage: mockArticleStorage, srcStorage: mockSourceStorage},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewSourceService(tt.args.articleRepo, tt.args.srcRepo), "NewSourceService(%v, %v)", tt.args.articleRepo, tt.args.srcRepo)
		})
	}
}

func Test_getFilesInDir(t *testing.T) {
	tests := []struct {
		name    string
		envVar  string
		want    []string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:   "get files successfully",
			envVar: "testdata",
			want: []string{"testdata/empty.json", "testdata/empty_rss.xml", "testdata/file1.txt",
				"testdata/file2.txt", "testdata/invalid.json", "testdata/invalid_rss.xml", "testdata/json.json",
				"testdata/rss.xml"},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
		{
			name:   "environment variable not set",
			envVar: "",
			want:   nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err) && assert.Equal(t, "environment variable DATA_DIR not set", err.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envVar != "" {
				os.Setenv("DATA_DIR", tt.envVar)
				defer os.Unsetenv("DATA_DIR")
			} else {
				os.Unsetenv("DATA_DIR")
			}

			got, err := getFilesInDir()
			if !tt.wantErr(t, err, fmt.Sprintf("getFilesInDir()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "getFilesInDir()")
		})
	}
}

func Test_sourceService_AddSource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSourceStorage := mocks.NewMockSourceStorage(ctrl)

	tests := []struct {
		name    string
		source  model.Source
		setup   func()
		want    model.Source
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:   "add source successfully",
			source: model.Source{Link: "http://example.com"},
			setup: func() {
				mockSourceStorage.EXPECT().Save(gomock.Any()).Return(model.Source{Id: 1, Link: "http://example.com"}, nil)
			},
			want: model.Source{Id: 1, Link: "http://example.com"},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
		{
			name:   "add duplicate source",
			source: model.Source{Link: "http://example.com"},
			setup: func() {
				mockSourceStorage.EXPECT().Save(gomock.Any()).Return(model.Source{}, errors.New("source already exists"))
			},
			want: model.Source{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err) && assert.Equal(t, "source already exists", err.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := &sourceService{
				srcStorage: mockSourceStorage,
			}
			got, err := s.AddSource(tt.source)
			if !tt.wantErr(t, err, fmt.Sprintf("AddSource(%v)", tt.source)) {
				return
			}
			assert.Equalf(t, tt.want, got, "AddSource(%v)", tt.source)
		})
	}
}

func Test_sourceService_DeleteSource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSourceStorage := mocks.NewMockSourceStorage(ctrl)
	mockArticleStorage := mocks.NewMockArticleStorage(ctrl)

	tests := []struct {
		name    string
		id      int
		setup   func()
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "delete source successfully",
			id:   1,
			setup: func() {
				mockArticleStorage.EXPECT().DeleteBySourceID(1).Return(nil)
				mockSourceStorage.EXPECT().Delete(1).Return(nil)
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := &sourceService{
				srcStorage:     mockSourceStorage,
				articleStorage: mockArticleStorage,
			}
			tt.wantErr(t, s.DeleteSource(tt.id), fmt.Sprintf("DeleteSource(%v)", tt.id))
		})
	}
}

func Test_sourceService_FetchFromAllSources(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSourceStorage := mocks.NewMockSourceStorage(ctrl)
	mockArticleStorage := mocks.NewMockArticleStorage(ctrl)

	tests := []struct {
		name    string
		setup   func()
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "fetch articles from all sources successfully",
			setup: func() {
				mockSourceStorage.EXPECT().GetAll().Return([]model.Source{
					{Id: 1, Link: "http://rss.cnn.com/rss/cnn_topstories.rss"},
				}, nil)
				mockArticleStorage.EXPECT().SaveAll(gomock.Any()).Return(nil)
				urlParsed, _ := url.Parse("http://rss.cnn.com/rss/cnn_topstories.rss")
				_, err := parser.ParseArticlesFromFeed(*urlParsed)
				if err != nil {
					return
				}
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
		{
			name: "error fetching sources",
			setup: func() {
				mockSourceStorage.EXPECT().GetAll().Return(nil, errors.New("storage error"))
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err) && assert.Equal(t, "storage error", err.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := &sourceService{
				articleStorage: mockArticleStorage,
				srcStorage:     mockSourceStorage,
			}
			tt.wantErr(t, s.FetchFromAllSources(), fmt.Sprintf("FetchFromAllSources()"))
		})
	}
}

func Test_sourceService_LoadDataFromFiles(t *testing.T) {
	type args struct {
		files []string
	}
	tests := []struct {
		name    string
		args    args
		want    []model.Article
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				files: []string{"testdata/rss.xml", "testdata/json.json"},
			},
			want: []model.Article{
				{
					Title:       "Article 1",
					Link:        "http://example.com/article1",
					Description: "This is the first article.",
					Source:      model.Source{Name: "Sample Feed"},
					PubDate:     time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC),
				},
				{
					Title:       "Article 2",
					Link:        "http://example.com/article2",
					Description: "This is the second article.",
					Source:      model.Source{Name: "Sample Feed"},
					PubDate:     time.Date(2006, time.January, 3, 15, 4, 5, 0, time.UTC),
				},
				{
					Title:       "Test Title",
					Link:        "http://testurl.com",
					Description: "Test Description",
					Source:      model.Source{Name: "Test Source"},
					PubDate:     time.Date(2023, 6, 4, 12, 0, 0, 0, time.UTC),
				},
			},
			wantErr: false,
		},
		{
			name: "one invalid file",
			args: args{
				files: []string{"testdata/rss.xml", "testdata/invalid.json"},
			},
			want: []model.Article{
				{
					Title:       "Article 1",
					Link:        "http://example.com/article1",
					Description: "This is the first article.",
					Source:      model.Source{Name: "Sample Feed"},
					PubDate:     time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC),
				},
				{
					Title:       "Article 2",
					Link:        "http://example.com/article2",
					Description: "This is the second article.",
					Source:      model.Source{Name: "Sample Feed"},
					PubDate:     time.Date(2006, time.January, 3, 15, 4, 5, 0, time.UTC),
				},
			},
			wantErr: true,
		},
		{
			name: "all invalid files",
			args: args{
				files: []string{"testdata/invalid_rss.xml", "testdata/invalid.json"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty files",
			args: args{
				files: []string{"testdata/empty_rss.xml", "testdata/empty.json"},
			},
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := require.New(t)
			var got []model.Article
			var err error
			for _, file := range tt.args.files {
				articles, e := parser.ParseArticlesFromFile(file)
				if e != nil {
					err = e
				} else {
					got = append(got, articles...)
				}
			}
			if tt.wantErr {
				assert.Error(err, "ParseArticlesFromFiles() should return an error")
			} else {
				assert.NoError(err, "ParseArticlesFromFiles() should not return an error")
			}
			assert.Equal(tt.want, got, "ParseArticlesFromFiles() returned unexpected result")
		})
	}
}

func Test_sourceService_UpdateSource(t *testing.T) {
	type fields struct {
		articleStorage storage.ArticleStorage
		srcStorage     storage.SourceStorage
	}
	type args struct {
		id     int
		source model.Source
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.Source
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "update source successfully",
			fields: fields{
				articleStorage: nil,
				srcStorage:     mocks.NewMockSourceStorage(gomock.NewController(t)),
			},
			args: args{
				id: 1,
				source: model.Source{
					Id:   1,
					Name: "CNN",
					Link: "http://rss.cnn.com/rss/cnn_topstories.rss",
				},
			},
			want: model.Source{
				Id:   1,
				Name: "CNN",
				Link: "http://rss.cnn.com/rss/cnn_topstories.rss",
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &sourceService{
				articleStorage: tt.fields.articleStorage,
				srcStorage:     tt.fields.srcStorage,
			}
			tt.fields.srcStorage.(*mocks.MockSourceStorage).EXPECT().Update(tt.args.id, tt.args.source).Return(tt.want, nil)
			got, err := s.UpdateSource(tt.args.id, tt.args.source)
			if !tt.wantErr(t, err, fmt.Sprintf("UpdateSource(%v, %v)", tt.args.id, tt.args.source)) {
				return
			}
			assert.Equalf(t, tt.want, got, "UpdateSource(%v, %v)", tt.args.id, tt.args.source)
		})
	}
}
