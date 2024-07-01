package web

import (
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// fetchSrcById fetches a source by its ID and returns it in the response.
func (h *Handler) fetchSrcById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	fetchedArticles, err := h.SrcService().FetchSourceByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, fetchedArticles)
}

// createSource creates a new source from the request body and returns it in the response.
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

// deleteSource deletes a source by its ID and returns a success message.
func (h *Handler) deleteSource(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = h.SrcService().DeleteSource(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "source deleted"})
}

func (h *Handler) updateSource(c *gin.Context) {
	var input model.Source
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	sources, err := h.srcService.UpdateSource(id, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, sources)
}
