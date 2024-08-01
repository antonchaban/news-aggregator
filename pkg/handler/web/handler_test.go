package web

import (
	service_mocks "github.com/antonchaban/news-aggregator/pkg/handler/web/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandler_ArticleService(t *testing.T) {
	mockArticleService := new(service_mocks.MockArticleService)
	mockSourceService := new(service_mocks.MockSourceService)

	h := NewHandler(mockArticleService, mockSourceService)

	assert.Equal(t, mockArticleService, h.ArticleService())
}

func TestHandler_InitRoutes(t *testing.T) {
	mockArticleService := new(service_mocks.MockArticleService)
	mockSourceService := new(service_mocks.MockSourceService)

	h := NewHandler(mockArticleService, mockSourceService)
	router := h.InitRoutes()

	assert.NotNil(t, router)

	// Check if the routes are properly set up
	routes := router.Routes()
	expectedRoutes := []string{"/swagger/*any", "/articles", "/sources/:id", "/sources", "/sources/:id", "/sources/:id"}

	for _, route := range expectedRoutes {
		found := false
		for _, r := range routes {
			if r.Path == route {
				found = true
				break
			}
		}
		assert.True(t, found, "Route %s not found", route)
	}
}

func TestHandler_SrcService(t *testing.T) {
	mockArticleService := new(service_mocks.MockArticleService)
	mockSourceService := new(service_mocks.MockSourceService)

	h := NewHandler(mockArticleService, mockSourceService)

	assert.Equal(t, mockSourceService, h.SrcService())
}

func TestNewHandler(t *testing.T) {
	mockArticleService := new(service_mocks.MockArticleService)
	mockSourceService := new(service_mocks.MockSourceService)

	h := NewHandler(mockArticleService, mockSourceService)

	assert.Equal(t, mockArticleService, h.ArticleService())
	assert.Equal(t, mockSourceService, h.SrcService())
}
