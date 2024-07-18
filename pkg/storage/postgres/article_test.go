package postgres

/*import (
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/antonchaban/news-aggregator/pkg/storage"
	"github.com/jackc/pgx/v5/pgxpool"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		db *pgxpool.Pool
	}
	tests := []struct {
		name string
		args args
		want storage.ArticleStorage
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_postgresArticleStorage_Delete(t *testing.T) {
	type fields struct {
		db *pgxpool.Pool
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
			pa := &postgresArticleStorage{
				db: tt.fields.db,
			}
			if err := pa.Delete(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_postgresArticleStorage_DeleteBySourceID(t *testing.T) {
	type fields struct {
		db *pgxpool.Pool
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
			pa := &postgresArticleStorage{
				db: tt.fields.db,
			}
			if err := pa.DeleteBySourceID(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("DeleteBySourceID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_postgresArticleStorage_GetAll(t *testing.T) {
	type fields struct {
		db *pgxpool.Pool
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
			pa := &postgresArticleStorage{
				db: tt.fields.db,
			}
			got, err := pa.GetAll()
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

func Test_postgresArticleStorage_GetByFilter(t *testing.T) {
	type fields struct {
		db *pgxpool.Pool
	}
	type args struct {
		query string
		args  []interface{}
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
			pa := &postgresArticleStorage{
				db: tt.fields.db,
			}
			got, err := pa.GetByFilter(tt.args.query, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByFilter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByFilter() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_postgresArticleStorage_Save(t *testing.T) {
	type fields struct {
		db *pgxpool.Pool
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
			pa := &postgresArticleStorage{
				db: tt.fields.db,
			}
			got, err := pa.Save(tt.args.article)
			if (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Save() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_postgresArticleStorage_SaveAll(t *testing.T) {
	type fields struct {
		db *pgxpool.Pool
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
			pa := &postgresArticleStorage{
				db: tt.fields.db,
			}
			if err := pa.SaveAll(tt.args.articles); (err != nil) != tt.wantErr {
				t.Errorf("SaveAll() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
*/
