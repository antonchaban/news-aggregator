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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotFormat := DetermineFileFormat(tt.args.filename); gotFormat != tt.wantFormat {
				t.Errorf("DetermineFileFormat() = %v, want %v", gotFormat, tt.wantFormat)
			}
		})
	}
}
