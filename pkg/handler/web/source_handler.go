package web

import (
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) fetchSrcById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	fetchedArticles, err := h.SrcService().FetchSourceByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, fetchedArticles)
}

func (h *Handler) createSource(c *gin.Context) {
	var input model.Source
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	sources, err := h.SrcService().AddSource(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, sources)
}
