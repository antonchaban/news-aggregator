package service

import (
	"news-aggregator/pkg/model"
	"news-aggregator/pkg/storage"
	"reflect"
	"testing"
)

func TestArticleService_Create(t *testing.T) {
	type fields struct {
		articleStorage storage.Article
	}
	type args struct {
		article model.Article
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.Article
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &ArticleService{
				articleStorage: tt.fields.articleStorage,
			}
			got, err := a.Create(tt.args.article)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArticleService_Delete(t *testing.T) {
	type fields struct {
		articleStorage storage.Article
	}
	type args struct {
		id int
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
			a := &ArticleService{
				articleStorage: tt.fields.articleStorage,
			}
			if err := a.Delete(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestArticleService_GetAll(t *testing.T) {
	type fields struct {
		articleStorage storage.Article
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
			a := &ArticleService{
				articleStorage: tt.fields.articleStorage,
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

func TestArticleService_GetByDateInRange(t *testing.T) {
	type fields struct {
		articleStorage storage.Article
	}
	type args struct {
		startDate string
		endDate   string
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
			a := &ArticleService{
				articleStorage: tt.fields.articleStorage,
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

func TestArticleService_GetById(t *testing.T) {
	type fields struct {
		articleStorage storage.Article
	}
	type args struct {
		id int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.Article
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &ArticleService{
				articleStorage: tt.fields.articleStorage,
			}
			got, err := a.GetById(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetById() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArticleService_GetByKeyword(t *testing.T) {
	type fields struct {
		articleStorage storage.Article
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
			a := &ArticleService{
				articleStorage: tt.fields.articleStorage,
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

func TestArticleService_GetBySource(t *testing.T) {
	type fields struct {
		articleStorage storage.Article
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
			a := &ArticleService{
				articleStorage: tt.fields.articleStorage,
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

func TestArticleService_SaveAll(t *testing.T) {
	type fields struct {
		articleStorage storage.Article
	}
	type args struct {
		articles []model.Article
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &ArticleService{
				articleStorage: tt.fields.articleStorage,
			}
			a.SaveAll(tt.args.articles)
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		articleRepo storage.Article
	}
	tests := []struct {
		name string
		args args
		want *ArticleService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.articleRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
