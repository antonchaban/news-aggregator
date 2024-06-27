package parser

import (
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/antonchaban/news-aggregator/pkg/parser/json"
	"github.com/antonchaban/news-aggregator/pkg/parser/rss"
	"net/url"
	"reflect"
	"testing"
	"time"
)

func TestParseArticlesFromFeed(t *testing.T) {
	type args struct {
		urlPath url.URL
	}
	rssURL, _ := url.Parse("http://rss.cnn.com/rss/cnn_topstories.rss")
	tests := []struct {
		name    string
		args    args
		want    []model.Article
		wantErr bool
	}{
		{
			name: "rss",
			args: args{
				urlPath: *rssURL,
			},
			want: []model.Article{
				{
					Title:       "Some on-air claims about Dominion Voting Systems were false, Fox News acknowledges in statement after deal is announced",
					Description: "",
					Link:        "https://www.cnn.com/business/live-news/fox-news-dominion-trial-04-18-23/index.html",
					Source: model.Source{
						Name: "CNN.com - RSS Channel - HP Hero",
						Link: "http://rss.cnn.com/rss/cnn_topstories.rss",
					},
					PubDate: time.Date(2023, 4, 19, 12, 44, 51, 0, time.UTC),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseArticlesFromFeed(tt.args.urlPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseArticlesFromFeed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got[0], tt.want[0]) {

				t.Errorf("ParseArticlesFromFeed() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseArticlesFromFile(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name    string
		args    args
		want    []model.Article
		wantErr bool
	}{
		{
			name: "rss",
			args: args{
				file: "testdata/rss.xml",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseArticlesFromFile(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseArticlesFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseArticlesFromFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createParser(t *testing.T) {
	type args struct {
		format string
	}
	tests := []struct {
		name    string
		args    args
		want    Parser
		wantErr bool
	}{
		{
			name: "rss",
			args: args{
				format: rssFormat,
			},
			want:    &rss.Parser{},
			wantErr: false,
		},
		{
			name: "json",
			args: args{
				format: jsonFormat,
			},
			want:    &json.Parser{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createParser(tt.args.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("createParser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createParser() got = %v, want %v", got, tt.want)
			}
		})
	}
}
