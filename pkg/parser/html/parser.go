package html

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/antonchaban/news-aggregator/pkg/model"
	"net/http"
	"net/url"
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

func (h *Parser) ParseFeed(url url.URL) ([]model.Article, error) {
	resp, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch URL: %s", url.String())
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	return h.parseDocument(doc), nil
}

// ParseFile parses the given file and returns a slice of articles.
func (h *Parser) ParseFile(f *os.File) ([]model.Article, error) {
	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		return nil, err
	}
	return h.parseDocument(doc), nil
}

// parseDocument parses the goquery document and returns a slice of articles.
func (h *Parser) parseDocument(doc *goquery.Document) []model.Article {
	var articles []model.Article
	doc.Find(h.config.ArticleSelector).Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Contents().Not("svg").Text())
		link, _ := s.Attr("href")
		description := strings.TrimSpace(s.AttrOr(h.config.DescriptionSelector, ""))
		date := strings.TrimSpace(s.Find(h.config.PubDateSelector).AttrOr(h.config.DateAttribute, ""))
		parsedDate, _ := parseDate(date, h.config.TimeFormat)
		article := model.Article{
			Title:       title,
			Link:        resolveLink(link, "https://www.usatoday.com"),
			PubDate:     parsedDate,
			Source:      model.Source{Name: h.config.Source}, // todo think about how to add links too
			Description: description,
		}
		if article.Title != "" || article.Description != "" {
			articles = append(articles, article)
		}
	})

	return articles
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

// resolveLink resolves relative links to absolute using the base URL
func resolveLink(link, baseURL string) string {
	if strings.HasPrefix(link, "http") {
		return link
	}
	return baseURL + link
}
