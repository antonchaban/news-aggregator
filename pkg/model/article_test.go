package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestArticle_String(t *testing.T) {
	type fields struct {
		Id          int
		Title       string
		Description string
		Link        string
		Source      string
		PubDate     time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Test with all fields",
			fields: fields{
				Id:          1,
				Title:       "Test Title",
				Description: "Test Description",
				Link:        "http://example.com",
				Source:      "Example Source",
				PubDate:     time.Date(2023, time.June, 4, 12, 0, 0, 0, time.UTC),
			},
			want: "ID: 1\nTitle: Test Title\nDate: 2023-06-04 12:00:00 +0000 UTC\nDescription: Test Description\nLink: http://example.com\nSource: Example Source\n",
		},
		{
			name: "Test with empty fields",
			fields: fields{
				Id:          0,
				Title:       "",
				Description: "",
				Link:        "",
				Source:      "",
				PubDate:     time.Time{},
			},
			want: "ID: 0\nTitle: \nDate: 0001-01-01 00:00:00 +0000 UTC\nDescription: \nLink: \nSource: \n",
		},
		{
			name: "Test with partial fields",
			fields: fields{
				Id:          2,
				Title:       "Partial Title",
				Description: "",
				Link:        "http://partial.com",
				Source:      "",
				PubDate:     time.Date(2023, time.December, 31, 23, 59, 59, 0, time.UTC),
			},
			want: "ID: 2\nTitle: Partial Title\nDate: 2023-12-31 23:59:59 +0000 UTC\nDescription: \nLink: http://partial.com\nSource: \n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Article{
				Id:          tt.fields.Id,
				Title:       tt.fields.Title,
				Description: tt.fields.Description,
				Link:        tt.fields.Link,
				Source:      tt.fields.Source,
				PubDate:     tt.fields.PubDate,
			}
			assert.Equal(t, tt.want, a.String())
		})
	}
}
