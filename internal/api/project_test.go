package api_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kumachan-mis/knodeledge-api/internal/api"
	"github.com/kumachan-mis/knodeledge-api/internal/db"
	"github.com/kumachan-mis/knodeledge-api/internal/repository"
	"github.com/kumachan-mis/knodeledge-api/internal/service"
	"github.com/kumachan-mis/knodeledge-api/internal/usecase"
	"github.com/kumachan-mis/knodeledge-api/test/testutil"
	"github.com/stretchr/testify/assert"
)

func TestProjectList(t *testing.T) {
	router := setupProjectRouter()

	recorder := httptest.NewRecorder()
	requestBody, _ := json.Marshal(map[string]any{
		"user": map[string]any{
			"id": testutil.UserId(),
		},
	})
	req, _ := http.NewRequest("POST", "/api/projects/list", strings.NewReader(string(requestBody)))

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var responseBody map[string]any
	assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &responseBody))
	assert.Equal(t, responseBody, map[string]any{
		"projects": []any{
			map[string]any{
				"id":   "PROJECT_WITHOUT_DESCRIPTION",
				"name": "No Description Project",
			},
			map[string]any{
				"id":          "PROJECT_WITH_DESCRIPTION",
				"name":        "Described Project",
				"description": "This is project description",
			},
		},
	})
}

func TestProjectListInvalidArgument(t *testing.T) {
	router := setupProjectRouter()

	recorder := httptest.NewRecorder()
	requestBody, _ := json.Marshal(map[string]any{
		"user": map[string]any{
			"id": "",
		},
	})
	req, _ := http.NewRequest("POST", "/api/projects/list", strings.NewReader(string(requestBody)))

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	var responseBody map[string]any
	assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &responseBody))
	assert.Equal(t, map[string]any{
		"message": "invalid request value",
		"user": map[string]any{
			"id": "user id is required, but got ''",
		},
	}, responseBody)
}

func TestProjectListInvalidRequestFormat(t *testing.T) {
	tt := []struct {
		name    string
		request string
	}{
		{
			name:    "should return error when user id is not string",
			request: `{"user": {"id":123}}`,
		},
		{
			name:    "should return error when user is not object",
			request: `{"user": 123}`,
		},
		{
			name:    "should return error when request body is invalid JSON",
			request: fmt.Sprintf(`{"user": {"id": "%s"`, testutil.UserId()),
		},
		{
			name:    "should return error when request body is empty",
			request: "",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			router := setupProjectRouter()

			recorder := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/projects/list", strings.NewReader(tc.request))

			router.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusBadRequest, recorder.Code)

			var responseBody map[string]any
			assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &responseBody))
			assert.Equal(t, map[string]any{
				"message": "invalid request format",
				"user":    map[string]any{},
			}, responseBody)
		})
	}
}

func TestProjectListInternalError(t *testing.T) {
	router := setupProjectRouter()

	recorder := httptest.NewRecorder()
	requestBody, _ := json.Marshal(map[string]any{
		"user": map[string]any{
			"id": testutil.ErrorUserId(0),
		},
	})
	req, _ := http.NewRequest("POST", "/api/projects/list", strings.NewReader(string(requestBody)))

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	var responseBody map[string]any
	assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &responseBody))
	assert.Equal(t, map[string]any{
		"message": "internal error",
	}, responseBody)
}

func setupProjectRouter() *gin.Engine {
	router := gin.Default()

	client := db.FirestoreClient()
	r := repository.NewProjectRepository(*client)
	s := service.NewProjectService(r)
	uc := usecase.NewProjectUseCase(s)
	a := api.NewProjectApi(uc)

	router.POST("/api/projects/list", a.HandleList)
	return router
}
