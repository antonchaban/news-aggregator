package filter

import (
	"fmt"
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeywordFilter_Filter(t *testing.T) {
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
			name: "Articles with matching keywords",
			args: args{
				articles: []model.Article{
					{
						Id:          1,
						Title:       "Golang is great",
						Description: "Go is an open-source programming language",
					},
					{
						Id:          2,
						Title:       "Python vs Golang",
						Description: "Comparing Python and Golang for backend development",
					},
					{
						Id:          3,
						Title:       "Learning Java",
						Description: "Java is a popular language",
					},
				},
				f: Filters{
					Keyword: "Golang",
				},
			},
			want: []model.Article{
				{
					Id:          1,
					Title:       "Golang is great",
					Description: "Go is an open-source programming language",
				},
				{
					Id:          2,
					Title:       "Python vs Golang",
					Description: "Comparing Python and Golang for backend development",
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Articles without matching keywords",
			args: args{
				articles: []model.Article{
					{
						Id:          1,
						Title:       "Learning Java",
						Description: "Java is a popular language",
					},
					{
						Id:          2,
						Title:       "Python tutorial",
						Description: "Learn Python programming",
					},
				},
				f: Filters{
					Keyword: "Golang",
				},
			},
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			name: "Multiple keywords",
			args: args{
				articles: []model.Article{
					{
						Id:          1,
						Title:       "Golang is great",
						Description: "Go is an open-source programming language",
					},
					{
						Id:          2,
						Title:       "Python vs Golang",
						Description: "Comparing Python and Golang for backend development",
					},
					{
						Id:          3,
						Title:       "Learning Java",
						Description: "Java is a popular language",
					},
				},
				f: Filters{
					Keyword: "Golang, Java",
				},
			},
			want: []model.Article{
				{
					Id:          1,
					Title:       "Golang is great",
					Description: "Go is an open-source programming language",
				},
				{
					Id:          2,
					Title:       "Python vs Golang",
					Description: "Comparing Python and Golang for backend development",
				},
				{
					Id:          3,
					Title:       "Learning Java",
					Description: "Java is a popular language",
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Empty keyword",
			args: args{
				articles: []model.Article{
					{
						Id:          1,
						Title:       "Golang is great",
						Description: "Go is an open-source programming language",
					},
					{
						Id:          2,
						Title:       "Python vs Golang",
						Description: "Comparing Python and Golang for backend development",
					},
				},
				f: Filters{
					Keyword: "",
				},
			},
			want: []model.Article{
				{
					Id:          1,
					Title:       "Golang is great",
					Description: "Go is an open-source programming language",
				},
				{
					Id:          2,
					Title:       "Python vs Golang",
					Description: "Comparing Python and Golang for backend development",
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &KeywordFilter{
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
