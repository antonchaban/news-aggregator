package filter

import (
	"github.com/stretchr/testify/assert"
	"news-aggregator/pkg/model"
	"testing"
)

func Test_intersect(t *testing.T) {
	type args struct {
		a []model.Article
		b []model.Article
	}
	tests := []struct {
		name         string
		args         args
		wantArticles []model.Article
	}{
		{
			name: "No intersection",
			args: args{
				a: []model.Article{{Id: 1}, {Id: 2}},
				b: []model.Article{{Id: 3}, {Id: 4}},
			},
			wantArticles: nil,
		},
		{
			name: "Some intersection",
			args: args{
				a: []model.Article{{Id: 1}, {Id: 2}},
				b: []model.Article{{Id: 2}, {Id: 3}},
			},
			wantArticles: []model.Article{{Id: 2}},
		},
		{
			name: "Complete intersection",
			args: args{
				a: []model.Article{{Id: 1}, {Id: 2}},
				b: []model.Article{{Id: 1}, {Id: 2}},
			},
			wantArticles: []model.Article{{Id: 1}, {Id: 2}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.wantArticles, intersect(tt.args.a, tt.args.b), "intersect(%v, %v)", tt.args.a, tt.args.b)
		})
	}
}
