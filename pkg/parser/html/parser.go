package html

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	eventParseHtmlFeedStart     = "parse_html_feed_start"
	eventHttpGetError           = "http_get_error"
	eventHttpStatusError        = "http_status_error"
	eventParseHtmlDocumentError = "parse_html_document_error"
	eventParseHtmlFeedSuccess   = "parse_html_feed_success"
	eventParseHtmlFileStart     = "parse_html_file_start"
	eventParseHtmlFileSuccess   = "parse_html_file_success"
	eventParseHtmlDocumentStart = "parse_html_document_start"
	eventParseHtmlDocumentEnd   = "parse_html_document_end"
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
	logrus.WithField("event_id", eventParseHtmlFeedStart).Infof("Starting to parse feed from URL: %s", url.String())

	resp, err := http.Get(url.String())
	if err != nil {
		logrus.WithField("event_id", eventHttpGetError).Errorf("Error fetching URL: %s", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logrus.WithField("event_id", eventHttpStatusError).Errorf("Failed to fetch URL: %s, Status Code: %d", url.String(), resp.StatusCode)
		return nil, fmt.Errorf("failed to fetch URL: %s", url.String())
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logrus.WithField("event_id", eventParseHtmlDocumentError).Errorf("Error parsing document from URL: %s", err.Error())
		return nil, err
	}

	logrus.WithField("event_id", eventParseHtmlFeedSuccess).Info("Successfully fetched and parsed feed")
	return h.parseDocument(doc), nil
}

// ParseFile parses the given file and returns a slice of articles.
func (h *Parser) ParseFile(f *os.File) ([]model.Article, error) {
	logrus.WithField("event_id", eventParseHtmlFileStart).Infof("Starting to parse file: %s", f.Name())

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		logrus.WithField("event_id", eventParseHtmlDocumentError).Errorf("Error parsing document from file: %s", err.Error())
		return nil, err
	}

	logrus.WithField("event_id", eventParseHtmlFileSuccess).Infof("Successfully parsed file: %s", f.Name())
	return h.parseDocument(doc), nil
}

// parseDocument parses the goquery document and returns a slice of articles.
func (h *Parser) parseDocument(doc *goquery.Document) []model.Article {
	logrus.WithField("event_id", eventParseHtmlDocumentStart).Info("Starting to parse document")
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
			Source:      model.Source{Name: h.config.Source},
			Description: description,
		}
		if article.Title != "" || article.Description != "" {
			articles = append(articles, article)
		}
	})

	logrus.WithField("event_id", eventParseHtmlDocumentEnd).Infof("Document parsing completed, found %d articles", len(articles))
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
	return time.Now().UTC(), fmt.Errorf("failed to parse date: %s", date)
}

// resolveLink resolves relative links to absolute using the base URL.
func resolveLink(link, baseURL string) string {
	if strings.HasPrefix(link, "http") {
		return link
	}
	return baseURL + link
}
