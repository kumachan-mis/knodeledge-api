package api_test

import (
	"encoding/json"
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

func TestChapterList(t *testing.T) {
	router := setupChapterRouter()

	recorder := httptest.NewRecorder()
	requestBody, _ := json.Marshal(map[string]any{
		"user": map[string]any{
			"id": testutil.ReadOnlyUserId(),
		},
		"project": map[string]any{
			"id": "PROJECT_WITHOUT_DESCRIPTION",
		},
	})
	req, _ := http.NewRequest("POST", "/api/chapters/list", strings.NewReader(string(requestBody)))

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var responseBody map[string]any
	assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &responseBody))
	assert.Equal(t, responseBody, map[string]any{
		"chapters": []any{
			map[string]any{
				"id":       "CHAPTER_ONE",
				"name":     "Chapter One",
				"number":   float64(1), // json.Unmarshal() converts numbers to float64
				"sections": []any{},
			},
			map[string]any{
				"id":       "CHAPTER_TWO",
				"name":     "Chapter Two",
				"number":   float64(2), // json.Unmarshal() converts numbers to float64
				"sections": []any{},
			},
		},
	})
}

func TestChapterListEmpty(t *testing.T) {
	tt := []struct {
		name      string
		userId    string
		projectId string
	}{
		{
			name:      "Project without chapters",
			userId:    testutil.ReadOnlyUserId(),
			projectId: "PROJECT_WITH_DESCRIPTION",
		},
		{
			name:      "Project not authorized",
			userId:    testutil.ModifyOnlyUserId(),
			projectId: "PROJECT_WITH_DESCRIPTION",
		},
		{
			name:      "Project not found",
			userId:    testutil.ReadOnlyUserId(),
			projectId: "UNKNOWN_PROJECT",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			router := setupChapterRouter()

			recorder := httptest.NewRecorder()
			requestBody, _ := json.Marshal(map[string]any{
				"user": map[string]any{
					"id": tc.userId,
				},
				"project": map[string]any{
					"id": tc.projectId,
				},
			})
			req, _ := http.NewRequest("POST", "/api/chapters/list", strings.NewReader(string(requestBody)))

			router.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusOK, recorder.Code)

			var responseBody map[string]any
			assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &responseBody))
			assert.Equal(t, responseBody, map[string]any{
				"chapters": []any{},
			})
		})
	}
}

func TestChapterListInvalidArgument(t *testing.T) {
	tt := []struct {
		name             string
		request          map[string]any
		expectedResponse map[string]any
	}{
		{
			name:    "Empty request",
			request: map[string]any{},
			expectedResponse: map[string]any{
				"user": map[string]any{
					"id": "user id is required, but got ''",
				},
				"project": map[string]any{
					"id": "project id is required, but got ''",
				},
			},
		},
		{
			name: "Empty user id",
			request: map[string]any{
				"user": map[string]any{
					"id": "",
				},
				"project": map[string]any{
					"id": "PROJECT_WITHOUT_DESCRIPTION",
				},
			},
			expectedResponse: map[string]any{
				"user": map[string]any{
					"id": "user id is required, but got ''",
				},
				"project": map[string]any{},
			},
		},
		{
			name: "Empty project id",
			request: map[string]any{
				"user": map[string]any{
					"id": testutil.ReadOnlyUserId(),
				},
				"project": map[string]any{
					"id": "",
				},
			},
			expectedResponse: map[string]any{
				"user": map[string]any{},
				"project": map[string]any{
					"id": "project id is required, but got ''",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			router := setupChapterRouter()

			recorder := httptest.NewRecorder()
			requestBody, _ := json.Marshal(tc.request)
			req, _ := http.NewRequest("POST", "/api/chapters/list", strings.NewReader(string(requestBody)))

			router.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusBadRequest, recorder.Code)

			var responseBody map[string]any
			assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &responseBody))
			assert.Equal(t, map[string]any{
				"message": "invalid request value",
				"user":    tc.expectedResponse["user"],
				"project": tc.expectedResponse["project"],
			}, responseBody)
		})
	}
}

func TestChapterListInvalidRequestFormat(t *testing.T) {
	tt := []struct {
		name    string
		request string
	}{
		{
			name:    "should return error when user id is not string",
			request: `{"user": {"id":123}, "project": {"id": "PROJECT_WITHOUT_DESCRIPTION"}}`,
		},
		{
			name:    "should return error when project id is not string",
			request: `{"user": {"id": "user-id"}, "project": {"id": 123}}`,
		},
		{
			name:    "should return error when user is not object",
			request: `{"user": 123, "project": {"id": "PROJECT_WITHOUT_DESCRIPTION"}}`,
		},
		{
			name:    "should return error when project is not object",
			request: `{"user": {"id": "user-id"}, "project": "PROJECT_WITHOUT_DESCRIPTION"}`,
		},
		{
			name:    "should return error when request body is invalid JSON",
			request: `{"user": {"id": "user-id", "project": {"id": "PROJECT_WITHOUT_DESCRIPTION"}`,
		},
		{
			name:    "should return error when request body is empty",
			request: "",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			router := setupChapterRouter()

			recorder := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/chapters/list", strings.NewReader(tc.request))

			router.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusBadRequest, recorder.Code)

			var responseBody map[string]any
			assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &responseBody))
			assert.Equal(t, map[string]any{
				"message": "invalid request format",
				"user":    map[string]any{},
				"project": map[string]any{},
			}, responseBody)
		})
	}
}

func TestChapterListInternalError(t *testing.T) {
	router := setupChapterRouter()

	recorder := httptest.NewRecorder()
	requestBody, _ := json.Marshal(map[string]any{
		"user": map[string]any{
			"id": testutil.ReadOnlyUserId(),
		},
		"project": map[string]any{
			"id": "PROJECT_WITH_INVALID_CHAPTER_NAME",
		},
	})
	req, _ := http.NewRequest("POST", "/api/chapters/list", strings.NewReader(string(requestBody)))

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	var responseBody map[string]any
	assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &responseBody))
	assert.Equal(t, map[string]any{
		"message": "internal error",
	}, responseBody)
}

func setupChapterRouter() *gin.Engine {
	router := gin.Default()

	client := db.FirestoreClient()
	r := repository.NewChapterRepository(*client)
	s := service.NewChapterService(r)
	uc := usecase.NewChapterUseCase(s)
	api := api.NewChapterApi(uc)

	router.POST("/api/chapters/list", api.HandleList)

	return router
}
