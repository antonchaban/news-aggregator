package parser

import (
	"news-aggregator/pkg/model"
	"reflect"
	"testing"
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseArticlesFromFiles(tt.args.files)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseArticlesFromFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseArticlesFromFiles() got = %v, want %v", got, tt.want)
			}
		})
	}
}
