package web

import (
	"github.com/antonchaban/news-aggregator/pkg/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Handler represents the HTTP handler with article and source services.
type Handler struct {
	articleService service.ArticleService
	srcService     service.SourceService
}

// SrcService returns the source service.
func (h *Handler) SrcService() service.SourceService {
	return h.srcService
}

// ArticleService returns the article service.
func (h *Handler) ArticleService() service.ArticleService {
	return h.articleService
}

// NewHandler creates a new Handler instance.
func NewHandler(asvc service.ArticleService, ss service.SourceService) *Handler {
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
