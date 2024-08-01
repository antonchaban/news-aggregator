package web

import (
	service_mocks "github.com/antonchaban/news-aggregator/pkg/handler/web/mocks"
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestHandler_fetchSrcById(t *testing.T) {
	type mockBehavior func(r *service_mocks.MockSourceService, id int)
	tests := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedCode         int
		expectedResponseBody string
		inputID              string
	}{
		{
			name: "OK",
			mockBehavior: func(r *service_mocks.MockSourceService, id int) {
				r.EXPECT().FetchSourceByID(id).Return([]model.Article{
					{
						Id:          1,
						Title:       "Title1",
						Link:        "Link1",
						Description: "Description1",
						Source: model.Source{
							Name: "CNN",
						},
					},
					{
						Id:          2,
						Title:       "Title2",
						Link:        "Link2",
						Description: "Description2",
						Source: model.Source{
							Name: "CNN",
						},
					},
				}, nil)
			},
			expectedCode:         200,
			expectedResponseBody: `[{"Id":1,"Title":"Title1","Description":"Description1","Link":"Link1","Source":{"id":0,"name":"CNN","link":""},"PubDate":"0001-01-01T00:00:00Z"},{"Id":2,"Title":"Title2","Description":"Description2","Link":"Link2","Source":{"id":0,"name":"CNN","link":""},"PubDate":"0001-01-01T00:00:00Z"}]`,
			inputID:              "1",
		},
		{
			name: "BadRequest",
			mockBehavior: func(r *service_mocks.MockSourceService, id int) {
				// no expectations
			},
			expectedCode:         400,
			expectedResponseBody: `{"message":"strconv.Atoi: parsing \"abc\": invalid syntax"}`,
			inputID:              "abc",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			srcSvc := service_mocks.NewMockSourceService(c)
			id, _ := strconv.Atoi(test.inputID)
			test.mockBehavior(srcSvc, id)

			// Init Endpoint
			r := gin.New()
			r.GET("/source/:id", NewHandler(nil, srcSvc).fetchSrcById)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/source/"+test.inputID, nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_createSource(t *testing.T) {
	type mockBehavior func(r *service_mocks.MockSourceService, src model.Source)
	tests := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedCode         int
		expectedResponseBody string
		inputBody            string
	}{
		{
			name: "OK",
			mockBehavior: func(r *service_mocks.MockSourceService, src model.Source) {
				r.EXPECT().AddSource(src).Return(src, nil)
			},
			expectedCode:         200,
			expectedResponseBody: `{"id":1,"name":"CNN","link":"http://cnn.com"}`,
			inputBody:            `{"id":1,"name":"CNN","link":"http://cnn.com"}`,
		},
		{
			name:                 "BadRequest",
			mockBehavior:         func(r *service_mocks.MockSourceService, src model.Source) {},
			expectedCode:         400,
			expectedResponseBody: `{"message":"EOF"}`,
			inputBody:            ``,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			srcSvc := service_mocks.NewMockSourceService(c)
			var src model.Source
			if test.inputBody != "" {
				src = model.Source{Id: 1, Name: "CNN", Link: "http://cnn.com"}
			}
			test.mockBehavior(srcSvc, src)

			// Init Endpoint
			r := gin.New()
			r.POST("/source", NewHandler(nil, srcSvc).createSource)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/source", strings.NewReader(test.inputBody))
			req.Header.Set("Content-Type", "application/json")

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedCode)
			assert.JSONEq(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_deleteSource(t *testing.T) {
	type mockBehavior func(r *service_mocks.MockSourceService, id int)
	tests := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedCode         int
		expectedResponseBody string
		inputID              string
	}{
		{
			name: "OK",
			mockBehavior: func(r *service_mocks.MockSourceService, id int) {
				r.EXPECT().DeleteSource(id).Return(nil)
			},
			expectedCode:         200,
			expectedResponseBody: `{"message":"source deleted"}`,
			inputID:              "1",
		},
		{
			name: "BadRequest",
			mockBehavior: func(r *service_mocks.MockSourceService, id int) {
				// no expectations
			},
			expectedCode:         400,
			expectedResponseBody: `{"message":"strconv.Atoi: parsing \"abc\": invalid syntax"}`,
			inputID:              "abc",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			srcSvc := service_mocks.NewMockSourceService(c)
			id, _ := strconv.Atoi(test.inputID)
			test.mockBehavior(srcSvc, id)

			// Init Endpoint
			r := gin.New()
			r.DELETE("/source/:id", NewHandler(nil, srcSvc).deleteSource)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/source/"+test.inputID, nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_updateSource(t *testing.T) {
	type mockBehavior func(r *service_mocks.MockSourceService, id int, src model.Source)
	tests := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedCode         int
		expectedResponseBody string
		inputBody            string
		inputID              string
	}{
		{
			name: "OK",
			mockBehavior: func(r *service_mocks.MockSourceService, id int, src model.Source) {
				r.EXPECT().UpdateSource(id, src).Return(src, nil)
			},
			expectedCode:         200,
			expectedResponseBody: `{"id":1,"name":"CNN","link":"http://cnn.com"}`,
			inputBody:            `{"id":1,"name":"CNN","link":"http://cnn.com"}`,
			inputID:              "1",
		},
		{
			name:                 "BadRequest",
			mockBehavior:         func(r *service_mocks.MockSourceService, id int, src model.Source) {},
			expectedCode:         400,
			expectedResponseBody: `{"message":"EOF"}`,
			inputBody:            ``,
			inputID:              "1",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			srcSvc := service_mocks.NewMockSourceService(c)
			var src model.Source
			if test.inputBody != "" {
				src = model.Source{Id: 1, Name: "CNN", Link: "http://cnn.com"}
			}
			id, _ := strconv.Atoi(test.inputID)
			test.mockBehavior(srcSvc, id, src)

			// Init Endpoint
			r := gin.New()
			r.PUT("/source/:id", NewHandler(nil, srcSvc).updateSource)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/source/"+test.inputID, strings.NewReader(test.inputBody))
			req.Header.Set("Content-Type", "application/json")

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedCode)
			assert.JSONEq(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}
