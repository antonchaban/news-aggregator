package web

import (
	"github.com/antonchaban/news-aggregator/pkg/filter"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) getArticlesByFilter(c *gin.Context) {
	f := filter.Filters{
		Keyword:   c.Query("keywords"),
		Source:    c.Query("sources"),
		StartDate: c.Query("date_start"),
		EndDate:   c.Query("date_end"),
	}

	articles, err := h.articleService.GetByFilter(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, articles)
}
