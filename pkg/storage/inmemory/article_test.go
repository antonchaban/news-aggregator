package inmemory

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"news-aggregator/pkg/model"
	"news-aggregator/pkg/storage/mocks"
	"reflect"
	"testing"
	"time"
)

func TestArticleInMemory_Create(t *testing.T) {
	tests := []struct {
		name           string
		newArticle     model.Article
		createdArticle model.Article
	}{
		{
			name: "Create new article successfully",
			newArticle: model.Article{
				Title:       "New Article",
				Description: "New Description",
				Link:        "http://newlink.com",
				Source:      "New Source",
				PubDate:     time.Now(),
			},
			createdArticle: model.Article{
				Id:          1,
				Title:       "New Article",
				Description: "New Description",
				Link:        "http://newlink.com",
				Source:      "New Source",
				PubDate:     time.Now(),
			},
		},
		{
			name: "Create another new article successfully",
			newArticle: model.Article{
				Title:       "Another New Article",
				Description: "Another New Description",
				Link:        "http://anotherlink.com",
				Source:      "Another New Source",
				PubDate:     time.Now(),
			},
			createdArticle: model.Article{
				Id:          2,
				Title:       "Another New Article",
				Description: "Another New Description",
				Link:        "http://anotherlink.com",
				Source:      "Another New Source",
				PubDate:     time.Now(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockArticleStorage := mocks.NewMockArticleStorage(ctrl)

			mockArticleStorage.EXPECT().Create(tt.newArticle).Return(tt.createdArticle, nil)

			article, err := mockArticleStorage.Create(tt.newArticle)

			assert.NoError(t, err)
			assert.Equal(t, tt.createdArticle, article)
		})
	}
}

func TestArticleInMemory_Delete(t *testing.T) {
	tests := []struct {
		name      string
		id        int
		expectErr bool
	}{
		{
			name:      "Delete article successfully",
			id:        1,
			expectErr: false,
		},
		{
			name:      "Delete non-existing article",
			id:        2,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockArticleStorage := mocks.NewMockArticleStorage(ctrl)
			if tt.expectErr {
				mockArticleStorage.EXPECT().Delete(tt.id).Return(fmt.Errorf("article not found"))
			} else {
				mockArticleStorage.EXPECT().Delete(tt.id).Return(nil)
			}

			err := mockArticleStorage.Delete(tt.id)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestArticleInMemory_GetAll(t *testing.T) {
	type fields struct {
		Articles []model.Article
		nextID   int
	}
	tests := []struct {
		name    string
		fields  fields
		want    []model.Article
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &ArticleInMemory{
				Articles: tt.fields.Articles,
				nextID:   tt.fields.nextID,
			}
			got, err := a.GetAll()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAll() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArticleInMemory_GetByDateInRange(t *testing.T) {
	type fields struct {
		Articles []model.Article
		nextID   int
	}
	type args struct {
		startDate time.Time
		endDate   time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []model.Article
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &ArticleInMemory{
				Articles: tt.fields.Articles,
				nextID:   tt.fields.nextID,
			}
			got, err := a.GetByDateInRange(tt.args.startDate, tt.args.endDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByDateInRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByDateInRange() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArticleInMemory_GetByKeyword(t *testing.T) {
	type fields struct {
		Articles []model.Article
		nextID   int
	}
	type args struct {
		keyword string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []model.Article
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &ArticleInMemory{
				Articles: tt.fields.Articles,
				nextID:   tt.fields.nextID,
			}
			got, err := a.GetByKeyword(tt.args.keyword)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByKeyword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByKeyword() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArticleInMemory_GetBySource(t *testing.T) {
	type fields struct {
		Articles []model.Article
		nextID   int
	}
	type args struct {
		source string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []model.Article
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &ArticleInMemory{
				Articles: tt.fields.Articles,
				nextID:   tt.fields.nextID,
			}
			got, err := a.GetBySource(tt.args.source)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBySource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBySource() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArticleInMemory_SaveAll(t *testing.T) {
	type fields struct {
		Articles []model.Article
		nextID   int
	}
	type args struct {
		articles []model.Article
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &ArticleInMemory{
				Articles: tt.fields.Articles,
				nextID:   tt.fields.nextID,
			}
			if err := a.SaveAll(tt.args.articles); (err != nil) != tt.wantErr {
				t.Errorf("SaveAll() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *ArticleInMemory
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
