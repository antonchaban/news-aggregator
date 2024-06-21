package web

import (
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) getSrcById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	sources, err := h.SrcStorage().GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sources)
}

func (h *Handler) createSource(c *gin.Context) {
	var input model.Source
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	sources, err := h.SrcStorage().Save(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, sources)
}
