package rss

import (
	"news-aggregator/pkg/model"
	"os"
	"reflect"
	"testing"
)

func TestParser_ParseFile(t *testing.T) {
	type args struct {
		f *os.File
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
