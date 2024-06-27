package rss

import (
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/mmcdole/gofeed"
	"net/url"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestParser_ParseFile(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		want     []model.Article
		wantErr  bool
	}{
		{
			name:     "Valid RSS",
			fileName: "valid_rss.xml",
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
			wantErr: false,
		},
		{
			name:     "Invalid RSS",
			fileName: "invalid_rss.xml",
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "Empty RSS",
			fileName: "empty_rss.xml",
			want:     []model.Article{},
			wantErr:  false,
		},
		{
			name:     "Missing Fields RSS",
			fileName: "missing_fields_rss.xml",
			want: []model.Article{
				{
					Title:  "Article 1",
					Link:   "http://example.com/article1",
					Source: model.Source{Name: "Sample Feed"},
				},
				{
					Title:       "Article 2",
					Link:        "http://example.com/article2",
					Description: "This is the second article.",
					Source:      model.Source{Name: "Sample Feed"},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Parser{}
			f := loadTestData(t, tt.fileName)
			defer f.Close()
			got, err := r.ParseFile(f)
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

func loadTestData(t *testing.T, filename string) *os.File {
	file, err := os.Open("testdata/" + filename)
	if err != nil {
		t.Fatalf("failed to open test data file %s: %v", filename, err)
	}
	return file
}

func TestParser_parseFeed(t *testing.T) {
	type args struct {
		feed    *gofeed.Feed
		feedUrl url.URL
	}
	tests := []struct {
		name string
		args args
		want []model.Article
	}{
		{
			name: "Valid RSS",
			args: args{
				feed: &gofeed.Feed{
					Title: "Sample Feed",
					Items: []*gofeed.Item{
						{
							Title:       "Article 1",
							Link:        "http://example.com/article1",
							Description: "This is the first article.",
						},
						{
							Title:       "Article 2",
							Link:        "http://example.com/article2",
							Description: "This is the second article.",
						},
					},
				},
				feedUrl: url.URL{},
			},
			want: []model.Article{
				{
					Title:       "Article 1",
					Link:        "http://example.com/article1",
					Description: "This is the first article.",
					Source:      model.Source{Name: "Sample Feed"},
				},
				{
					Title:       "Article 2",
					Link:        "http://example.com/article2",
					Description: "This is the second article.",
					Source:      model.Source{Name: "Sample Feed"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Parser{}
			if got := r.parseFeed(tt.args.feed, tt.args.feedUrl); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseFeed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_ParseFile1(t *testing.T) {
	type args struct {
		f *os.File
	}
	tests := []struct {
		name    string
		args    args
		want    []model.Article
		wantErr bool
	}{
		{
			name: "Valid RSS",
			args: args{
				f: loadTestData(t, "valid_rss.xml"),
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
			wantErr: false,
		},
		{
			name: "Invalid RSS",
			args: args{
				f: loadTestData(t, "invalid_rss.xml"),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Empty RSS",
			args: args{
				f: loadTestData(t, "empty_rss.xml"),
			},
			want:    []model.Article{},
			wantErr: false,
		},
		{
			name: "Missing Fields RSS",
			args: args{
				f: loadTestData(t, "missing_fields_rss.xml"),
			},
			want: []model.Article{
				{
					Title:  "Article 1",
					Link:   "http://example.com/article1",
					Source: model.Source{Name: "Sample Feed"},
				},
				{
					Title:       "Article 2",
					Link:        "http://example.com/article2",
					Description: "This is the second article.",
					Source:      model.Source{Name: "Sample Feed"},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Parser{}
			got, err := r.ParseFile(tt.args.f)
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
