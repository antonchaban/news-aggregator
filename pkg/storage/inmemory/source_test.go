package inmemory

import (
	"fmt"
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/antonchaban/news-aggregator/pkg/storage"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSrc(t *testing.T) {
	tests := []struct {
		name string
		want storage.SourceStorage
	}{
		{
			name: "initialize in-memory storage",
			want: &memorySourceStorage{
				Sources: []model.Source{},
				nextID:  1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewSrc(), "NewSrc()")
		})
	}
}

func Test_memorySourceStorage_Delete(t *testing.T) {
	type fields struct {
		Sources []model.Source
		nextID  int
	}
	type args struct {
		id int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "delete existing source",
			fields: fields{
				Sources: []model.Source{
					{Id: 1, Link: "http://example.com"},
				},
				nextID: 2,
			},
			args: args{id: 1},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
		{
			name: "delete non-existing source",
			fields: fields{
				Sources: []model.Source{},
				nextID:  1,
			},
			args: args{id: 1},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err) && assert.Equal(t, "source not found", err.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &memorySourceStorage{
				Sources: tt.fields.Sources,
				nextID:  tt.fields.nextID,
			}
			tt.wantErr(t, m.Delete(tt.args.id), fmt.Sprintf("Delete(%v)", tt.args.id))
		})
	}
}

func Test_memorySourceStorage_GetAll(t *testing.T) {
	type fields struct {
		Sources []model.Source
		nextID  int
	}
	tests := []struct {
		name    string
		fields  fields
		want    []model.Source
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "get all sources",
			fields: fields{
				Sources: []model.Source{
					{Id: 1, Link: "http://example.com"},
				},
				nextID: 2,
			},
			want: []model.Source{
				{Id: 1, Link: "http://example.com"},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
		{
			name: "get all sources when empty",
			fields: fields{
				Sources: []model.Source{},
				nextID:  1,
			},
			want: []model.Source{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &memorySourceStorage{
				Sources: tt.fields.Sources,
				nextID:  tt.fields.nextID,
			}
			got, err := m.GetAll()
			if !tt.wantErr(t, err, fmt.Sprintf("GetAll()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetAll()")
		})
	}
}

func Test_memorySourceStorage_GetByID(t *testing.T) {
	type fields struct {
		Sources []model.Source
		nextID  int
	}
	type args struct {
		id int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.Source
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "get existing source by ID",
			fields: fields{
				Sources: []model.Source{
					{Id: 1, Link: "http://example.com"},
				},
				nextID: 2,
			},
			args: args{id: 1},
			want: model.Source{Id: 1, Link: "http://example.com"},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
		{
			name: "get non-existing source by ID",
			fields: fields{
				Sources: []model.Source{},
				nextID:  1,
			},
			args: args{id: 1},
			want: model.Source{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err) && assert.Equal(t, "source not found", err.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &memorySourceStorage{
				Sources: tt.fields.Sources,
				nextID:  tt.fields.nextID,
			}
			got, err := m.GetByID(tt.args.id)
			if !tt.wantErr(t, err, fmt.Sprintf("GetByID(%v)", tt.args.id)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetByID(%v)", tt.args.id)
		})
	}
}

func Test_memorySourceStorage_Save(t *testing.T) {
	type fields struct {
		Sources []model.Source
		nextID  int
	}
	type args struct {
		src model.Source
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.Source
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "save new source",
			fields: fields{
				Sources: []model.Source{},
				nextID:  1,
			},
			args: args{src: model.Source{Link: "http://example.com"}},
			want: model.Source{Id: 1, Link: "http://example.com"},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
		{
			name: "save duplicate source",
			fields: fields{
				Sources: []model.Source{
					{Id: 1, Link: "http://example.com"},
				},
				nextID: 2,
			},
			args: args{src: model.Source{Link: "http://example.com"}},
			want: model.Source{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err) && assert.Equal(t, "source already exists", err.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &memorySourceStorage{
				Sources: tt.fields.Sources,
				nextID:  tt.fields.nextID,
			}
			got, err := m.Save(tt.args.src)
			if !tt.wantErr(t, err, fmt.Sprintf("Save(%v)", tt.args.src)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Save(%v)", tt.args.src)
		})
	}
}

func Test_memorySourceStorage_SaveAll(t *testing.T) {
	type fields struct {
		Sources []model.Source
		nextID  int
	}
	type args struct {
		sources []model.Source
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "save all sources",
			fields: fields{
				Sources: []model.Source{},
				nextID:  1,
			},
			args: args{sources: []model.Source{
				{Link: "http://example.com"},
				{Link: "http://example.org"},
			}},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &memorySourceStorage{
				Sources: tt.fields.Sources,
				nextID:  tt.fields.nextID,
			}
			tt.wantErr(t, m.SaveAll(tt.args.sources), fmt.Sprintf("SaveAll(%v)", tt.args.sources))
		})
	}
}
