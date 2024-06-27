package json

import (
	"github.com/antonchaban/news-aggregator/pkg/model"
	"net/url"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestParser_ParseFile(t *testing.T) {
	tests := []struct {
		name      string
		fileName  string
		want      []model.Article
		wantErr   bool
		setupFunc func() (*os.File, func(), error)
	}{
		{
			name:     "Valid JSON File",
			fileName: "valid.json",
			want: []model.Article{
				{
					Title:       "Test Title",
					Link:        "http://testurl.com",
					Description: "Test Description",
					Source:      model.Source{Name: "Test Source"},
					PubDate:     time.Date(2023, 6, 4, 12, 0, 0, 0, time.UTC),
				},
			},
			wantErr: false,
			setupFunc: func() (*os.File, func(), error) {
				file, err := os.Open("testdata/valid.json")
				return file, func() { file.Close() }, err
			},
		},
		{
			name:     "Invalid JSON File",
			fileName: "invalid.json",
			want:     nil,
			wantErr:  true,
			setupFunc: func() (*os.File, func(), error) {
				file, err := os.Open("pkg/parser/json/testdata/invalid.json")
				return file, func() { file.Close() }, err
			},
		},
		{
			name:     "Empty JSON File",
			fileName: "empty.json",
			want:     []model.Article{},
			wantErr:  false,
			setupFunc: func() (*os.File, func(), error) {
				file, err := os.Open("testdata/empty.json")
				return file, func() { file.Close() }, err
			},
		},
		{
			name:     "File Read Error",
			fileName: "nonexistent.json",
			want:     nil,
			wantErr:  true,
			setupFunc: func() (*os.File, func(), error) {
				return nil, func() {}, os.ErrNotExist
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, cleanup, err := tt.setupFunc()
			if err != nil {
				if !tt.wantErr {
					t.Fatalf("setupFunc() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			defer cleanup()

			j := &Parser{}
			got, err := j.ParseFile(f)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_ParseFeed(t *testing.T) {
	type args struct {
		url url.URL
	}
	tests := []struct {
		name    string
		args    args
		want    []model.Article
		wantErr bool
	}{
		{
			name:    "Not Implemented",
			args:    args{},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Parser{}
			got, err := j.ParseFeed(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFeed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseFeed() got = %v, want %v", got, tt.want)
			}
		})
	}
}
