package web

import (
	"github.com/antonchaban/news-aggregator/pkg/filter"
	"github.com/antonchaban/news-aggregator/pkg/model"
	service_mocks "github.com/antonchaban/news-aggregator/pkg/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestHandler_getArticlesByFilter(t *testing.T) {
	type mockBehavior func(r *service_mocks.MockArticleService, filters filter.Filters)
	tests := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedCode         int
		expectedResponseBody string
		inputQuery           map[string]string
	}{
		{
			name: "OK",
			mockBehavior: func(r *service_mocks.MockArticleService, filters filter.Filters) {
				r.EXPECT().GetByFilter(filters).Return([]model.Article{
					{
						Id:          1,
						Title:       "Title",
						Link:        "Link",
						Description: "Description",
						Source: model.Source{
							Name: "CNN",
						},
					},
				}, nil)
			},
			expectedCode:         200,
			expectedResponseBody: `[{"Id":1,"Title":"Title","Description":"Description","Link":"Link","Source":{"id":0,"name":"CNN","link":""},"PubDate":"0001-01-01T00:00:00Z"}]`,
			inputQuery: map[string]string{
				"sources": "other",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			artSvc := service_mocks.NewMockArticleService(c)
			sSvc := service_mocks.NewMockSourceService(c)
			for range test.inputQuery {
				f := filter.Filters{
					Keyword:   test.inputQuery["keywords"],
					Source:    test.inputQuery["sources"],
					StartDate: test.inputQuery["date_start"],
					EndDate:   test.inputQuery["date_end"],
				}
				test.mockBehavior(artSvc, f)
			}

			// Init Endpoint
			r := gin.New()
			r.POST("/articles", NewHandler(artSvc, sSvc).getArticlesByFilter)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/articles?sources=other", nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}
