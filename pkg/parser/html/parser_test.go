package html

import (
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewHtmlParser(t *testing.T) {
	tests := []struct {
		name   string
		config FeedConfig
		want   *Parser
	}{
		{
			name: "Test with valid config",
			config: FeedConfig{
				ArticleSelector:     ".article",
				TitleSelector:       ".title",
				LinkSelector:        ".link",
				DescriptionSelector: ".description",
				PubDateSelector:     ".date",
				Source:              "Test Source",
				DateAttribute:       "datetime",
				TimeFormat:          []string{time.RFC3339},
			},
			want: &Parser{
				config: FeedConfig{
					ArticleSelector:     ".article",
					TitleSelector:       ".title",
					LinkSelector:        ".link",
					DescriptionSelector: ".description",
					PubDateSelector:     ".date",
					Source:              "Test Source",
					DateAttribute:       "datetime",
					TimeFormat:          []string{time.RFC3339},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewHtmlParser(tt.config)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParser_ParseFile(t *testing.T) {
	type fields struct {
		config FeedConfig
	}
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.Article
		wantErr bool
	}{
		{
			name: "Test USA TODAY HTML file",
			fields: fields{
				config: FeedConfig{
					ArticleSelector:     "a.gnt_m_flm_a",
					LinkSelector:        "",
					DescriptionSelector: "data-c-br",
					PubDateSelector:     "div.gnt_m_flm_sbt",
					Source:              "USA TODAY",
					DateAttribute:       "data-c-dt",
					TimeFormat: []string{
						"2006-01-02 15:04",
						"Jan 02, 2006",
					},
				},
			},
			args: args{
				filePath: "../../../data/usatoday-world-news.html", // Ensure this file exists
			},
			want: model.Article{
				Title:       "A Russian escalation. A new front. Is Ukraine losing the war with Russia?",
				Link:        "https://www.usatoday.com/story/news/world/2024/05/19/ukraine-losing-war-with-russia/73730454007/",
				PubDate:     time.Date(2024, 5, 19, 10, 42, 0, 0, time.UTC),
				Source:      "USA TODAY",
				Description: "In recent days, Russia's forces have seized territory near Kharkiv in Ukraine. Is Ukraine losing the war with Russia?",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Parser{
				config: tt.fields.config,
			}
			file, err := os.Open(filepath.Clean(tt.args.filePath))
			if err != nil {
				t.Fatalf("Failed to open file: %v", err)
			}
			defer file.Close()

			got, err := h.ParseFile(file)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got[1])
			}
		})
	}
}

func Test_parseDate(t *testing.T) {
	tests := []struct {
		name           string
		date           string
		timeFormats    []string
		wantParsedDate time.Time
		wantErr        bool
	}{
		{
			name:           "Test valid RFC3339 date",
			date:           "2023-06-04T12:00:00Z",
			timeFormats:    []string{time.RFC3339},
			wantParsedDate: time.Date(2023, time.June, 4, 12, 0, 0, 0, time.UTC),
			wantErr:        false,
		},
		{
			name:           "Test invalid date",
			date:           "invalid date",
			timeFormats:    []string{time.RFC3339},
			wantParsedDate: time.Now().UTC(), // Current time as fallback
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotParsedDate, err := parseDate(tt.date, tt.timeFormats)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantParsedDate, gotParsedDate)
			}
		})
	}
}
