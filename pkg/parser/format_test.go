package parser

import (
	"net/url"
	"testing"
)

func TestDetermineFileFormat(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name       string
		args       args
		wantFormat string
	}{
		{
			name: "rss",
			args: args{
				filename: "file.xml",
			},
			wantFormat: rssFormat,
		},
		{
			name: "json",
			args: args{
				filename: "file.json",
			},
			wantFormat: jsonFormat,
		},
		{
			name: "html",
			args: args{
				filename: "file.html",
			},
			wantFormat: htmlFormat,
		},
		{
			name: "unknown",
			args: args{
				filename: "file.txt",
			},
			wantFormat: unknownFormat,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotFormat := DetermineFileFormat(tt.args.filename); gotFormat != tt.wantFormat {
				t.Errorf("DetermineFileFormat() = %v, want %v", gotFormat, tt.wantFormat)
			}
		})
	}
}

func TestDetermineFeedFormat(t *testing.T) {
	type args struct {
		urlPath url.URL
	}
	rssURL, _ := url.Parse("http://rss.cnn.com/rss/cnn_topstories.rss")
	tests := []struct {
		name       string
		args       args
		wantFormat string
		wantErr    bool
	}{
		{
			name: "html",
			args: args{
				urlPath: url.URL{
					Scheme: "http",
					Host:   "example.com",
				},
			},
			wantFormat: htmlFormat,
			wantErr:    false,
		},
		{
			name: "rss",
			args: args{
				urlPath: *rssURL,
			},
			wantFormat: rssFormat,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFormat, err := DetermineFeedFormat(tt.args.urlPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("DetermineFeedFormat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotFormat != tt.wantFormat {
				t.Errorf("DetermineFeedFormat() gotFormat = %v, want %v", gotFormat, tt.wantFormat)
			}
		})
	}
}
