package inmemory

import (
	"github.com/stretchr/testify/assert"
	"news-aggregator/pkg/model"
	"testing"
	"time"
)

func TestArticleInMemory_Create(t *testing.T) {
	r := New()

	type args struct {
		article model.Article
	}
	tests := []struct {
		name    string
		args    args
		want    model.Article
		wantErr bool
	}{
		{
			name: "Create article",
			args: args{
				article: model.Article{
					Title:       "Test Article",
					Description: "This is a test article",
					Source:      "Test Source",
					Link:        "http://example.com",
					PubDate:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			want: model.Article{
				Id:          1,
				Title:       "Test Article",
				Description: "This is a test article",
				Source:      "Test Source",
				Link:        "http://example.com",
				PubDate:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name: "Create second article",
			args: args{
				article: model.Article{
					Title:       "Test Article 2",
					Description: "This is a test article 2",
					Source:      "Test Source 2",
					Link:        "http://example.com/2",
					PubDate:     time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},
			want: model.Article{
				Id:          2,
				Title:       "Test Article 2",
				Description: "This is a test article 2",
				Source:      "Test Source 2",
				Link:        "http://example.com/2",
				PubDate:     time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.Create(tt.args.article)
			if tt.wantErr {
				assert.Error(t, err, "Create() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				assert.NoError(t, err, "Create() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, got, "Create() got = %v, want %v", got, tt.want)
		})
	}
}

func TestArticleInMemory_Delete(t *testing.T) {
	r := New()

	articles := []model.Article{
		{
			Id:          1,
			Title:       "Test Article",
			Description: "This is a test article",
			Source:      "Test Source",
			Link:        "http://example.com",
			PubDate:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			Id:          2,
			Title:       "Test Article 2",
			Description: "This is a test article 2",
			Source:      "Test Source 2",
			Link:        "http://example.com/2",
			PubDate:     time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			Id:          3,
			Title:       "Test Article 3",
			Description: "This is a test article 3",
			Source:      "Test Source 3",
			Link:        "http://example.com/3",
			PubDate:     time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
		},
	}

	err := r.SaveAll(articles)
	if err != nil {
		return
	}

	type args struct {
		id int
	}
	tests := []struct {
		name    string
		args    args
		want    []model.Article
		wantErr bool
	}{
		{
			name: "Delete article",
			args: args{
				id: 2,
			},
			want: []model.Article{
				{
					Id:          1,
					Title:       "Test Article",
					Description: "This is a test article",
					Source:      "Test Source",
					Link:        "http://example.com",
					PubDate:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					Id:          3,
					Title:       "Test Article 3",
					Description: "This is a test article 3",
					Source:      "Test Source 3",
					Link:        "http://example.com/3",
					PubDate:     time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
				},
			}, wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := r.Delete(tt.args.id)
			if err != nil {
				return
			}
			got, _ := r.GetAll()
			if tt.wantErr {
				assert.Error(t, err, "Delete() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				assert.NoError(t, err, "Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, got, "Delete() got = %v, want %v", got, tt.want)
		})
	}
}

//func TestArticleInMemory_GetAll(t *testing.T) {
//	type fields struct {
//		Articles []model.Article
//		nextID   int
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		want    []model.Article
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			a := &ArticleInMemory{
//				Articles: tt.fields.Articles,
//				nextID:   tt.fields.nextID,
//			}
//			got, err := a.GetAll()
//			if (err != nil) != tt.wantErr {
//				t.Errorf("GetAll() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("GetAll() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestArticleInMemory_GetByDateInRange(t *testing.T) {
//	type fields struct {
//		Articles []model.Article
//		nextID   int
//	}
//	type args struct {
//		startDate time.Time
//		endDate   time.Time
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    []model.Article
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			a := &ArticleInMemory{
//				Articles: tt.fields.Articles,
//				nextID:   tt.fields.nextID,
//			}
//			got, err := a.GetByDateInRange(tt.args.startDate, tt.args.endDate)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("GetByDateInRange() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("GetByDateInRange() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestArticleInMemory_GetById(t *testing.T) {
//	type fields struct {
//		Articles []model.Article
//		nextID   int
//	}
//	type args struct {
//		id int
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    model.Article
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			a := &ArticleInMemory{
//				Articles: tt.fields.Articles,
//				nextID:   tt.fields.nextID,
//			}
//			got, err := a.GetById(tt.args.id)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("GetById() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("GetById() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestArticleInMemory_GetByKeyword(t *testing.T) {
//	type fields struct {
//		Articles []model.Article
//		nextID   int
//	}
//	type args struct {
//		keyword string
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    []model.Article
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			a := &ArticleInMemory{
//				Articles: tt.fields.Articles,
//				nextID:   tt.fields.nextID,
//			}
//			got, err := a.GetByKeyword(tt.args.keyword)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("GetByKeyword() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("GetByKeyword() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestArticleInMemory_GetBySource(t *testing.T) {
//	type fields struct {
//		Articles []model.Article
//		nextID   int
//	}
//	type args struct {
//		source string
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    []model.Article
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			a := &ArticleInMemory{
//				Articles: tt.fields.Articles,
//				nextID:   tt.fields.nextID,
//			}
//			got, err := a.GetBySource(tt.args.source)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("GetBySource() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("GetBySource() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestArticleInMemory_SaveAll(t *testing.T) {
//	type fields struct {
//		Articles []model.Article
//		nextID   int
//	}
//	type args struct {
//		articles []model.Article
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			a := &ArticleInMemory{
//				Articles: tt.fields.Articles,
//				nextID:   tt.fields.nextID,
//			}
//			if err := a.SaveAll(tt.args.articles); (err != nil) != tt.wantErr {
//				t.Errorf("SaveAll() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func TestNew(t *testing.T) {
//	tests := []struct {
//		name string
//		want *ArticleInMemory
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := New(); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("New() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
