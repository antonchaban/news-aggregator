package web

import (
	"github.com/antonchaban/news-aggregator/pkg/service"
	"github.com/antonchaban/news-aggregator/pkg/storage"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	articleService service.ArticleService
	srcStorage     storage.SourceStorage
}

func (h *Handler) SrcStorage() storage.SourceStorage {
	return h.srcStorage
}

func (h *Handler) ArticleService() service.ArticleService {
	return h.articleService
}

func NewHandler(asvc service.ArticleService, ss storage.SourceStorage) *Handler {
	h := &Handler{articleService: asvc,
		srcStorage: ss}
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
		sources.GET("/:id", h.getSrcById)
		sources.POST("", h.createSource)
	}
	return router
}
