package parser

import (
	"github.com/stretchr/testify/require"
	"news-aggregator/pkg/model"
	"testing"
	"time"
)

func TestParseArticlesFromFiles(t *testing.T) {
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
					Source:      "Sample Feed",
					PubDate:     time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC),
				},
				{
					Title:       "Article 2",
					Link:        "http://example.com/article2",
					Description: "This is the second article.",
					Source:      "Sample Feed",
					PubDate:     time.Date(2006, time.January, 3, 15, 4, 5, 0, time.UTC),
				},
				{
					Title:       "Test Title",
					Link:        "http://testurl.com",
					Description: "Test Description",
					Source:      "Test Source",
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
					Source:      "Sample Feed",
					PubDate:     time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC),
				},
				{
					Title:       "Article 2",
					Link:        "http://example.com/article2",
					Description: "This is the second article.",
					Source:      "Sample Feed",
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
				articles, e := ParseArticlesFromFile(file)
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
