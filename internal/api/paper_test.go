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

func TestPaperFind(t *testing.T) {
	router := setupPaperRouter()

	ecorder := httptest.NewRecorder()
	requestBody, _ := json.Marshal(map[string]any{
		"user": map[string]any{
			"id": testutil.ReadOnlyUserId(),
		},
		"project": map[string]any{
			"id": "PROJECT_WITHOUT_DESCRIPTION",
		},
		"chapter": map[string]any{
			"id": "CHAPTER_ONE",
		},
	})
	req, _ := http.NewRequest("POST", "/api/papers/find", strings.NewReader(string(requestBody)))

	router.ServeHTTP(ecorder, req)

	assert.Equal(t, http.StatusOK, ecorder.Code)

	var responseBody map[string]any
	err := json.Unmarshal(ecorder.Body.Bytes(), &responseBody)
	assert.Nil(t, err)

	assert.Equal(t, map[string]any{
		"paper": map[string]any{
			"id":      "CHAPTER_ONE",
			"content": "This is paper content.",
		},
	}, responseBody)
}

func TestPaperFindProjectOrChapterNotFound(t *testing.T) {
	router := setupPaperRouter()

	tt := []struct {
		name    string
		user    string
		project string
		chapter string
	}{
		{
			name:    "should return error when project not found",
			user:    testutil.ReadOnlyUserId(),
			project: "UNKNOWN_PROJECT",
			chapter: "CHAPTER_ONE",
		},
		{
			name:    "should return not found when user is not author of the project",
			user:    testutil.ModifyOnlyUserId(),
			project: "PROJECT_WITH_DESCRIPTION",
			chapter: "CHAPTER_ONE",
		},
		{
			name:    "should return error when chapter not found",
			user:    testutil.ReadOnlyUserId(),
			project: "PROJECT_WITH_DESCRIPTION",
			chapter: "UNKNOWN_CHAPTER",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ecorder := httptest.NewRecorder()
			requestBody, _ := json.Marshal(map[string]any{
				"user": map[string]any{
					"id": tc.user,
				},
				"project": map[string]any{
					"id": tc.project,
				},
				"chapter": map[string]any{
					"id": tc.chapter,
				},
			})
			req, _ := http.NewRequest("POST", "/api/papers/find", strings.NewReader(string(requestBody)))

			router.ServeHTTP(ecorder, req)

			assert.Equal(t, http.StatusNotFound, ecorder.Code)

			var responseBody map[string]any
			err := json.Unmarshal(ecorder.Body.Bytes(), &responseBody)
			assert.Nil(t, err)

			assert.Equal(t, map[string]any{
				"message": "not found",
				"user":    map[string]any{},
				"project": map[string]any{},
				"chapter": map[string]any{},
			}, responseBody)
		})
	}
}

func TestPaperFindDomainValidationError(t *testing.T) {
	router := setupPaperRouter()

	tt := []struct {
		name             string
		request          map[string]any
		expectedResponse map[string]any
	}{
		{
			name: "should return error when user id is empty",
			request: map[string]any{
				"user": map[string]any{
					"id": "",
				},
				"project": map[string]any{
					"id": "PROJECT_WITHOUT_DESCRIPTION",
				},
				"chapter": map[string]any{
					"id": "CHAPTER_ONE",
				},
			},
			expectedResponse: map[string]any{
				"user": map[string]any{
					"id": "user id is required, but got ''",
				},
				"project": map[string]any{},
				"chapter": map[string]any{},
			},
		},
		{
			name: "should return error when project id is empty",
			request: map[string]any{
				"user": map[string]any{
					"id": testutil.ReadOnlyUserId(),
				},
				"project": map[string]any{
					"id": "",
				},
				"chapter": map[string]any{
					"id": "CHAPTER_ONE",
				},
			},
			expectedResponse: map[string]any{
				"user": map[string]any{},
				"project": map[string]any{
					"id": "project id is required, but got ''",
				},
				"chapter": map[string]any{},
			},
		},
		{
			name: "should return error when chapter id is empty",
			request: map[string]any{
				"user": map[string]any{
					"id": testutil.ReadOnlyUserId(),
				},
				"project": map[string]any{
					"id": "PROJECT_WITHOUT_DESCRIPTION",
				},
				"chapter": map[string]any{
					"id": "",
				},
			},
			expectedResponse: map[string]any{
				"user":    map[string]any{},
				"project": map[string]any{},
				"chapter": map[string]any{
					"id": "chapter id is required, but got ''",
				},
			},
		},
		{
			name:    "should return error when empty object is passed",
			request: map[string]any{},
			expectedResponse: map[string]any{
				"user": map[string]any{
					"id": "user id is required, but got ''",
				},
				"project": map[string]any{
					"id": "project id is required, but got ''",
				},
				"chapter": map[string]any{
					"id": "chapter id is required, but got ''",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ecorder := httptest.NewRecorder()
			requestBody, _ := json.Marshal(tc.request)
			req, _ := http.NewRequest("POST", "/api/papers/find", strings.NewReader(string(requestBody)))

			router.ServeHTTP(ecorder, req)

			assert.Equal(t, http.StatusBadRequest, ecorder.Code)

			var responseBody map[string]any
			err := json.Unmarshal(ecorder.Body.Bytes(), &responseBody)
			assert.Nil(t, err)

			assert.Equal(t, map[string]any{
				"message": "invalid request value",
				"user":    tc.expectedResponse["user"],
				"project": tc.expectedResponse["project"],
				"chapter": tc.expectedResponse["chapter"],
			}, responseBody)
		})
	}
}

func TestPaperFindInvalidRequestFormat(t *testing.T) {
	tt := []struct {
		name    string
		request string
	}{
		{
			name:    "should return error when user id is not string",
			request: `{"user": {"id": 1}, "project": {"id": "PROJECT_WITHOUT_DESCRIPTION"}, "chapter": {"id": "CHAPTER_ONE"}}`,
		},
		{
			name:    "should return error when project id is not string",
			request: `{"user": {"id": "USER_READ_ONLY"}, "project": {"id": 1}, "chapter": {"id": "CHAPTER_ONE"}}`,
		},
		{
			name:    "should return error when chapter id is not string",
			request: `{"user": {"id": "USER_READ_ONLY"}, "project": {"id": "PROJECT_WITHOUT_DESCRIPTION"}, "chapter": {"id": 1}}`,
		},
		{
			name:    "should return error when user is not object",
			request: `{"user": "USER_READ_ONLY", "project": {"id": "PROJECT_WITHOUT_DESCRIPTION"}, "chapter": {"id": "CHAPTER_ONE"}}`,
		},
		{
			name:    "should return error when project is not object",
			request: `{"user": {"id": "USER_READ_ONLY"}, "project": "PROJECT_WITHOUT_DESCRIPTION", "chapter": {"id": "CHAPTER_ONE"}}`,
		},
		{
			name:    "should return error when chapter is not object",
			request: `{"user": {"id": "USER_READ_ONLY"}, "project": {"id": "PROJECT_WITHOUT_DESCRIPTION"}, "chapter": "CHAPTER_ONE"}`,
		},
		{
			name:    "should return error when request is not object",
			request: `[]`,
		},
		{
			name:    "should return error when request is empty",
			request: ``,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			router := setupPaperRouter()

			ecorder := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/papers/find", strings.NewReader(tc.request))

			router.ServeHTTP(ecorder, req)

			assert.Equal(t, http.StatusBadRequest, ecorder.Code)

			var responseBody map[string]any
			err := json.Unmarshal(ecorder.Body.Bytes(), &responseBody)
			assert.Nil(t, err)

			assert.Equal(t, map[string]any{
				"message": "invalid request format",
				"user":    map[string]any{},
				"project": map[string]any{},
				"chapter": map[string]any{},
			}, responseBody)
		})
	}
}

func TestPaperFindInternalError(t *testing.T) {
	router := setupPaperRouter()

	ecorder := httptest.NewRecorder()
	requestBody, _ := json.Marshal(map[string]any{
		"user": map[string]any{
			"id": testutil.ErrorUserId(6),
		},
		"project": map[string]any{
			"id": "PROJECT_WITH_INVALID_PAPER_CONTENT",
		},
		"chapter": map[string]any{
			"id": "CHAPTER_WITH_INVALID_PAPER_CONTENT",
		},
	})
	req, _ := http.NewRequest("POST", "/api/papers/find", strings.NewReader(string(requestBody)))

	router.ServeHTTP(ecorder, req)

	assert.Equal(t, http.StatusInternalServerError, ecorder.Code)

	var responseBody map[string]any
	err := json.Unmarshal(ecorder.Body.Bytes(), &responseBody)
	assert.Nil(t, err)

	assert.Equal(t, map[string]any{
		"message": "internal error",
	}, responseBody)
}

func setupPaperRouter() *gin.Engine {
	router := gin.Default()

	client := db.FirestoreClient()
	r := repository.NewPaperRepository(*client)
	s := service.NewPaperService(r)
	uc := usecase.NewPaperUseCase(s)
	api := api.NewPaperApi(uc)

	router.POST("/api/papers/find", api.HandleFind)

	return router
}
