package web

import (
	"github.com/antonchaban/news-aggregator/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	articleService service.ArticleService
	srcService     service.SourceService
}

func (h *Handler) SrcService() service.SourceService {
	return h.srcService
}

func (h *Handler) ArticleService() service.ArticleService {
	return h.articleService
}

func NewHandler(asvc service.ArticleService, ss service.SourceService) *Handler {
	h := &Handler{articleService: asvc,
		srcService: ss}
	return h
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	articles := router.Group("/articles")
	{
		articles.GET("", h.getArticlesByFilter)
	}
	sources := router.Group("/sources")
	{
		sources.GET("/:id", h.fetchSrcById)
		sources.POST("", h.createSource)
		sources.DELETE("/:id", h.deleteSource)
	}
	return router
}
