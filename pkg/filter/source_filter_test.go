package filter

import (
	"fmt"
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSourceFilter_Filter(t *testing.T) {
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
			name: "Articles with matching sources",
			args: args{
				articles: []model.Article{
					{
						Id:     1,
						Title:  "Article 1",
						Source: abcNewsSource,
					},
					{
						Id:     2,
						Title:  "Article 2",
						Source: bbcNewsSource,
					},
					{
						Id:     3,
						Title:  "Article 3",
						Source: "Some other source",
					},
				},
				f: Filters{
					Source: "abcnews,bbc",
				},
			},
			want: []model.Article{
				{
					Id:     1,
					Title:  "Article 1",
					Source: abcNewsSource,
				},
				{
					Id:     2,
					Title:  "Article 2",
					Source: bbcNewsSource,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Articles without matching sources",
			args: args{
				articles: []model.Article{
					{
						Id:     1,
						Title:  "Article 1",
						Source: "Some other source",
					},
					{
						Id:     2,
						Title:  "Article 2",
						Source: "Another source",
					},
				},
				f: Filters{
					Source: "abcnews,bbc",
				},
			},
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			name: "Multiple sources",
			args: args{
				articles: []model.Article{
					{
						Id:     1,
						Title:  "Article 1",
						Source: abcNewsSource,
					},
					{
						Id:     2,
						Title:  "Article 2",
						Source: bbcNewsSource,
					},
					{
						Id:     3,
						Title:  "Article 3",
						Source: usaTodaySource,
					},
				},
				f: Filters{
					Source: "abcnews,bbc,usatoday",
				},
			},
			want: []model.Article{
				{
					Id:     1,
					Title:  "Article 1",
					Source: abcNewsSource,
				},
				{
					Id:     2,
					Title:  "Article 2",
					Source: bbcNewsSource,
				},
				{
					Id:     3,
					Title:  "Article 3",
					Source: usaTodaySource,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Empty source",
			args: args{
				articles: []model.Article{
					{
						Id:     1,
						Title:  "Article 1",
						Source: abcNewsSource,
					},
					{
						Id:     2,
						Title:  "Article 2",
						Source: bbcNewsSource,
					},
				},
				f: Filters{
					Source: "",
				},
			},
			want: []model.Article{
				{
					Id:     1,
					Title:  "Article 1",
					Source: abcNewsSource,
				},
				{
					Id:     2,
					Title:  "Article 2",
					Source: bbcNewsSource,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Invalid source",
			args: args{
				articles: []model.Article{
					{
						Id:     1,
						Title:  "Article 1",
						Source: abcNewsSource,
					},
				},
				f: Filters{
					Source: "invalidsource",
				},
			},
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &SourceFilter{
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
