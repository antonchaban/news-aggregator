package web

import (
	"github.com/antonchaban/news-aggregator/pkg/backuper"
	"github.com/antonchaban/news-aggregator/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	articleService service.ArticleService
}

func NewHandler(asvc service.ArticleService) *Handler {
	h := &Handler{articleService: asvc}
	articles, err := backuper.NewLoader(asvc).LoadAllFromFile()
	//https://abcnews.go.com/abcnews/internationalheadlines
	articlesFeed, err := backuper.NewLoader(asvc).UpdateFromFeed("https://feeds.bbci.co.uk/news/rss.xml")
	articles = append(articles, articlesFeed...)
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
