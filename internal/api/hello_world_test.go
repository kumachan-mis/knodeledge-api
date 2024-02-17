package api_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kumachan-mis/knodeledge-api/internal/api"
	"github.com/kumachan-mis/knodeledge-api/internal/db"
	"github.com/kumachan-mis/knodeledge-api/internal/model"
	"github.com/kumachan-mis/knodeledge-api/internal/repository"
	"github.com/kumachan-mis/knodeledge-api/internal/service"
	"github.com/kumachan-mis/knodeledge-api/internal/usecase"
	"github.com/stretchr/testify/assert"
)

func TestHandleHelloWorld(t *testing.T) {
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
			router := setupHelloWorldRouter()

			recorder := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/hello-world", strings.NewReader(tc.request))

			router.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusOK, recorder.Code)
			assert.Equal(t, tc.expected, recorder.Body.String())
		})
	}
}

func TestHandleHelloWorldError(t *testing.T) {
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
			name:    "should return error when request body is empty",
			request: strings.NewReader(""),
		},
		{
			name:    "should return error when request body is nil",
			request: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := setupHelloWorldRouter()

			recorder := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/hello-world", tc.request)

			router.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusBadRequest, recorder.Code)

			var response model.ApplicationErrorResponse
			assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))

			assert.NotEmpty(t, response.Message)

		})
	}
}

func setupHelloWorldRouter() *gin.Engine {
	router := gin.Default()

	client := db.FirestoreClient()
	r := repository.NewHelloWorldRepository(*client)
	s := service.NewHelloWorldService(r)
	uc := usecase.NewHelloWorldUseCase(s)
	a := api.NewHelloWorldApi(uc)

	router.POST("/api/hello-world", a.HandleHelloWorld)
	return router
}
