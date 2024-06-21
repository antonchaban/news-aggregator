package web

import (
	"github.com/antonchaban/news-aggregator/pkg/backuper"
	"github.com/antonchaban/news-aggregator/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	articleService service.ArticleService
}

func NewHandler(asvc service.ArticleService) *Handler {
	h := &Handler{articleService: asvc}
	articles, err := backuper.NewLoader(asvc).LoadAllFromFile()
	for _, feed := range getSupportedFeeds() {
		articlesFeed, err := backuper.NewLoader(asvc).UpdateFromFeed(feed)
		if err != nil {
			logrus.Fatalf("error occurred while updating articles from feed: %s", err.Error())
		}
		articles = append(articles, articlesFeed...)

	}
	if err != nil {
		return nil
	}
	err = asvc.SaveAll(articles)
	return h
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	articles := router.Group("/articles")
	{
		articles.GET("", h.getArticlesByFilter)
	}
	return router
}

func getSupportedFeeds() []string {
	return []string{
		"https://feeds.bbci.co.uk/news/rss.xml",
		"https://abcnews.go.com/abcnews/internationalheadlines",
		"https://www.washingtontimes.com/rss/headlines/news/world/",
	}
}
