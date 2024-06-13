package filter

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"news-aggregator/pkg/model"
	"testing"
	"time"
)

func TestDateRangeFilter_Filter(t *testing.T) {
	type fields struct {
		next ArticleFilter
	}
	type args struct {
		articles []model.Article
		f        Filters
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []model.Article
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Test case 1: Filter articles by date range",
			args: args{
				articles: []model.Article{
					{
						Id:      1,
						Title:   "Article 1",
						Source:  "Source 1",
						PubDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						Link:    "http://example.com/article1",
					},
					{
						Id:      2,
						Title:   "Article 2",
						Source:  "Source 2",
						PubDate: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
						Link:    "http://example.com/article2",
					},
					{
						Id:      3,
						Title:   "Article 3",
						Source:  "Source 3",
						PubDate: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
						Link:    "http://example.com/article3",
					},
				},
				f: Filters{
					StartDate: "2021-01-02",
					EndDate:   "2021-01-03",
				},
			},
			want: []model.Article{
				{
					Id:      2,
					Title:   "Article 2",
					Source:  "Source 2",
					PubDate: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
					Link:    "http://example.com/article2",
				},
				{
					Id:      3,
					Title:   "Article 3",
					Source:  "Source 3",
					PubDate: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
					Link:    "http://example.com/article3",
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Test case 2: Articles outside date range",
			args: args{
				articles: []model.Article{
					{
						Id:      1,
						Title:   "Article 1",
						Source:  "Source 1",
						PubDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						Link:    "http://example.com/article1",
					},
					{
						Id:      4,
						Title:   "Article 4",
						Source:  "Source 4",
						PubDate: time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC),
						Link:    "http://example.com/article4",
					},
				},
				f: Filters{
					StartDate: "2021-01-02",
					EndDate:   "2021-01-03",
				},
			},
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			name: "Test case 3: Empty start and end dates",
			args: args{
				articles: []model.Article{
					{
						Id:      1,
						Title:   "Article 1",
						Source:  "Source 1",
						PubDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						Link:    "http://example.com/article1",
					},
					{
						Id:      2,
						Title:   "Article 2",
						Source:  "Source 2",
						PubDate: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
						Link:    "http://example.com/article2",
					},
				},
				f: Filters{
					StartDate: "",
					EndDate:   "",
				},
			},
			want: []model.Article{
				{
					Id:      1,
					Title:   "Article 1",
					Source:  "Source 1",
					PubDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					Link:    "http://example.com/article1",
				},
				{
					Id:      2,
					Title:   "Article 2",
					Source:  "Source 2",
					PubDate: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
					Link:    "http://example.com/article2",
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &DateRangeFilter{
				next: tt.fields.next,
			}
			got, err := h.Filter(tt.args.articles, tt.args.f)
			if !tt.wantErr(t, err, fmt.Sprintf("Filter(%v, %v)", tt.args.articles, tt.args.f)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Filter(%v, %v)", tt.args.articles, tt.args.f)
		})
	}
}
