package web

import (
	"github.com/antonchaban/news-aggregator/pkg/filter"
	"github.com/antonchaban/news-aggregator/pkg/model"
	_ "github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

//go:generate mockgen -destination=mocks/mock_article_service.go -package=mocks github.com/antonchaban/news-aggregator/pkg/handler/web ArticleService

// ArticleService is an interface that defines the methods for interacting with the article storage.
type ArticleService interface {
	GetAll() ([]model.Article, error)
	Create(article model.Article) (model.Article, error)
	Delete(id int) error
	SaveAll(articles []model.Article) error
	GetByFilter(f filter.Filters) ([]model.Article, error)
}

// @Summary Get articles by filter
// @Description Get articles by filter parameters
// @Tags articles
// @ID get-articles-by-filter
// @Accept json
// @Produce json
// @Param keywords query string false "Keywords to search for"
// @Param sources query string false "Sources to search for"
// @Param date_start query string false "Start date for search"
// @Param date_end query string false "End date for search"
// @Success 200 {object} []model.Article
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /articles [get]
func (h *Handler) getArticlesByFilter(c *gin.Context) {
	f := filter.Filters{
		Keyword:   c.Query("keywords"),
		Source:    c.Query("sources"),
		StartDate: c.Query("date_start"),
		EndDate:   c.Query("date_end"),
	}

	articles, err := h.articleService.GetByFilter(f)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, articles)
}
