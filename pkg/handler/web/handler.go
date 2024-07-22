package web

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Handler represents the handler with article and source services.
type Handler struct {
	articleService ArticleService
	srcService     SourceService
}

// SrcService returns the source service.
func (h *Handler) SrcService() SourceService {
	return h.srcService
}

// ArticleService returns the article service.
func (h *Handler) ArticleService() ArticleService {
	return h.articleService
}

// NewHandler creates a new Handler instance.
func NewHandler(asvc ArticleService, ss SourceService) *Handler {
	h := &Handler{articleService: asvc,
		srcService: ss}
	return h
}

// InitRoutes initializes the routes for the HTTP server.
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	articles := router.Group("/articles")
	{
		articles.GET("", h.getArticlesByFilter)
	}
	sources := router.Group("/sources")
	{
		sources.GET("/:id", h.fetchSrcById)
		sources.POST("", h.createSource)
		sources.DELETE("/:id", h.deleteSource)
		sources.PUT("/:id", h.updateSource)
	}
	return router
}
