package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHelloWorldHandler(t *testing.T) {
	testCases := []struct {
		name     string
		request  string
		expected string
	}{
		{
			name:     "should return Hello, Kumachan! when name is Kumachan",
			request:  `{"name": "Kumachan"}`,
			expected: `{"message":"Hello, Kumachan!"}`,
		},
		{
			name:     "should return Hello World! when name is empty",
			request:  `{"name": ""}`,
			expected: `{"message":"Hello World!"}`,
		},
		{
			name:     "should return Hello World! when name is null",
			request:  `{"name": null}`,
			expected: `{"message":"Hello World!"}`,
		},
		{
			name:     "should return Hello World! when name is not specified",
			request:  `{}`,
			expected: `{"message":"Hello World!"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := setupRouter()

			recorder := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/hello-world", strings.NewReader(tc.request))

			router.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusOK, recorder.Code)
			assert.Equal(t, tc.expected, recorder.Body.String())
		})
	}
}

func TestHelloWorldHandlerError(t *testing.T) {
	testCases := []struct {
		name    string
		request io.Reader
	}{
		{
			name:    "should return error when name is not string",
			request: strings.NewReader(`{"name": 123}`),
		},
		{
			name:    "should return error when request body is invalid JSON",
			request: strings.NewReader(`{"name": "Kumachan"`),
		},
		{
			name:    "should return error when request body is nil",
			request: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := setupRouter()

			recorder := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/hello-world", nil)

			router.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusBadRequest, recorder.Code)
			assert.NotEmpty(t, recorder.Body.String())
		})
	}
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/api/hello-world", HelloWorldHandler)
	return router
}
