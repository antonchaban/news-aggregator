package parser

import (
	"github.com/PuerkitoBio/goquery"
	"news-aggregator/pkg/model"
	"os"
	"strings"
)

type HtmlParser struct {
	config HtmlFeedConfig
}

type HtmlFeedConfig struct {
	ArticleSelector     string
	TitleSelector       string
	LinkSelector        string
	DescriptionSelector string
	PubDateSelector     string
	Source              string
	DateAttribute       string
}

func NewHtmlParser(config HtmlFeedConfig) *HtmlParser {
	return &HtmlParser{config: config}
}

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
		date := strings.TrimSpace(s.Parent().Find(h.config.PubDateSelector).AttrOr(h.config.DateAttribute, ""))

		article := model.Article{
			Title:       title,
			Link:        url,
			PubDate:     date,
			Source:      h.config.Source,
			Description: description,
		}
		if article.Title != "" || article.Description != "" {
			articles = append(articles, article)
		}
	})

	return articles, nil
}
