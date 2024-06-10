package parser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"news-aggregator/pkg/model"
	"os"
	"strings"
	"time"
)

// HtmlParser is a struct that contains the configuration for parsing HTML feeds
// and implements the ParsingAlgorithm interface.
type HtmlParser struct {
	config HtmlFeedConfig
}

// HtmlFeedConfig is a struct that contains the configuration for parsing HTML feeds.
type HtmlFeedConfig struct {
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
func NewHtmlParser(config HtmlFeedConfig) *HtmlParser {
	return &HtmlParser{config: config}
}

// parseFile parses the given file and returns a slice of articles.
func (h *HtmlParser) parseFile(f *os.File) ([]model.Article, error) {
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
		parsedDate, err := parseDate(date, h.config.TimeFormat)
		if err != nil {
			//fmt.Println("Setting current date because of error parsing date:", err)
		}
		article := model.Article{
			Title:       title,
			Link:        url,
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
	return time.Now(), fmt.Errorf("no matching format for date: %s", date)
}
