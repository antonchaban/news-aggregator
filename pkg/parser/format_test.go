package parser

import "testing"

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
			wantFormat: "rss",
		},
		{
			name: "json",
			args: args{
				filename: "file.json",
			},
			wantFormat: "json",
		},
		{
			name: "html",
			args: args{
				filename: "file.html",
			},
			wantFormat: "html",
		},
		{
			name: "unknown",
			args: args{
				filename: "file.txt",
			},
			wantFormat: "unknown",
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
