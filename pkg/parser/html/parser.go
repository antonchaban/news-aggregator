package html

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/antonchaban/news-aggregator/pkg/model"
	"os"
	"strings"
	"time"
)

// Parser is a struct that contains the configuration for parsing HTML feeds
// and implements the Parser interface.
type Parser struct {
	config FeedConfig
}

// FeedConfig is a struct that contains the configuration for parsing HTML feeds.
type FeedConfig struct {
	ArticleSelector     string
	TitleSelector       string
	LinkSelector        string
	DescriptionSelector string
	PubDateSelector     string
	Source              string
	DateAttribute       string
	TimeFormat          []string
}

// NewHtmlParser creates a new HtmlParser with the given configuration.
func NewHtmlParser(config FeedConfig) *Parser {
	return &Parser{config: config}
}

// ParseFile parses the given file and returns a slice of articles.
func (h *Parser) ParseFile(f *os.File) ([]model.Article, error) {
	var articles []model.Article
	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		return nil, err
	}
	doc.Find(h.config.ArticleSelector).Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Text())
		url, _ := s.Attr("href")
		description := strings.TrimSpace(s.AttrOr(h.config.DescriptionSelector, ""))
		date := strings.TrimSpace(s.Find(h.config.PubDateSelector).AttrOr(h.config.DateAttribute, ""))
		parsedDate, _ := parseDate(date, h.config.TimeFormat)
		article := model.Article{
			Title:       title,
			Link:        "https://www.usatoday.com" + url,
			PubDate:     parsedDate,
			Source:      h.config.Source,
			Description: description,
		}
		if article.Title != "" || article.Description != "" {
			articles = append(articles, article)
		}
	})

	return articles, nil
}

// parseDate parses the given date string using the provided time formats.
func parseDate(date string, timeFormats []string) (parsedDate time.Time, err error) {
	for _, format := range timeFormats {
		parsedTime, err := time.Parse(format, date)

		if err == nil {
			return parsedTime, nil
		}
	}
	return time.Now().UTC(), errors.New(fmt.Sprintf("error parsing date: %s", date))
}
