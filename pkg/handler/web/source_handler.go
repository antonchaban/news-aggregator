package web

import (
	"github.com/antonchaban/news-aggregator/pkg/model"
	_ "github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//go:generate mockgen -destination=../service/mocks/mock_source_service.go -package=mocks github.com/antonchaban/news-aggregator/pkg/service SourceService

// SourceService represents the service for sources.
type SourceService interface {
	FetchFromAllSources() error
	FetchSourceByID(id int) ([]model.Article, error)
	LoadDataFromFiles() ([]model.Article, error)
	AddSource(source model.Source) (model.Source, error)
	DeleteSource(id int) error
	UpdateSource(id int, source model.Source) (model.Source, error)
	GetAll() ([]model.Source, error)
}

// @Summary Fetch source by ID
// @Description Immediately fetches news from source by ID
// @Tags sources
// @ID fetch-source-by-id
// @Accept json
// @Produce json
// @Param id path int true "Source ID"
// @Success 200 {object} []model.Article
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /sources/{id} [get]
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

// @Summary Create a new source
// @Description Create a new source
// @Tags sources
// @ID create-source
// @Accept json
// @Produce json
// @Param source body model.Source true "Source object"
// @Success 200 {object} model.Source
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /sources [post]
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

// @Summary Delete source by ID
// @Description Delete source and all associated articles by ID
// @Tags sources
// @ID delete-source-by-id
// @Accept json
// @Produce json
// @Param id path int true "Source ID"
// @Success 200 {object} errorResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /sources/{id} [delete]
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

// @Summary Update source by ID
// @Description Update source by ID
// @Tags sources
// @ID update-source-by-id
// @Accept json
// @Produce json
// @Param id path int true "Source ID"
// @Param source body model.Source true "Source object"
// @Success 200 {object} model.Source
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /sources/{id} [put]
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
