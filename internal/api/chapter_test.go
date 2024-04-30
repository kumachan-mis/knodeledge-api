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
	assert.Equal(t, map[string]any{
		"chapters": []any{
			map[string]any{
				"id":       "CHAPTER_ONE",
				"name":     "Chapter One",
				"nextId":   "CHAPTER_TWO",
				"sections": []any{},
			},
			map[string]any{
				"id":       "CHAPTER_TWO",
				"name":     "Chapter Two",
				"nextId":   "",
				"sections": []any{},
			},
		},
	}, responseBody)
}

func TestChapterListEmpty(t *testing.T) {
	router := setupChapterRouter()

	recorder := httptest.NewRecorder()
	requestBody, _ := json.Marshal(map[string]any{
		"user": map[string]any{
			"id": testutil.ReadOnlyUserId(),
		},
		"project": map[string]any{
			"id": "PROJECT_WITH_DESCRIPTION",
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
}

func TestChapterListProjectNotFound(t *testing.T) {
	tt := []struct {
		name    string
		request map[string]any
	}{
		{
			name: "should return invalid argument when project id is not found",
			request: map[string]any{
				"user": map[string]any{
					"id": testutil.ReadOnlyUserId(),
				},
				"project": map[string]any{
					"id": "UNKNOWN_PROJECT",
				},
			},
		},
		{
			name: "should return invalid argument when user is not author of the project",
			request: map[string]any{
				"user": map[string]any{
					"id": testutil.ModifyOnlyUserId(),
				},
				"project": map[string]any{
					"id": "PROJECT_WITHOUT_DESCRIPTION",
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
				"message": "invalid request value: " +
					"failed to list chapters: project document does not exist",
				"user":    map[string]any{},
				"project": map[string]any{},
			}, responseBody)
		})
	}
}

func TestChapterListDomainValidationError(t *testing.T) {
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
			"id": testutil.ErrorUserId(2),
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

func TestChapterCreate(t *testing.T) {
	router := setupChapterRouter()

	recorder := httptest.NewRecorder()
	requestBody, _ := json.Marshal(map[string]any{
		"user": map[string]any{
			"id": testutil.ModifyOnlyUserId(),
		},
		"project": map[string]any{
			"id": "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_API",
		},
		"chapter": map[string]any{
			"name": "Chapter One",
		},
	})
	req, _ := http.NewRequest("POST", "/api/chapters/create", strings.NewReader(string(requestBody)))

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusCreated, recorder.Code)

	var responseBody map[string]any
	assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &responseBody))

	//chapterId is generated by firestore and it's not predictable
	chapterId := responseBody["chapter"].(map[string]any)["id"]
	assert.NotEmpty(t, chapterId)

	assert.Equal(t, map[string]any{
		"chapter": map[string]any{
			"id":       chapterId,
			"name":     "Chapter One",
			"nextId":   "",
			"sections": []any{},
		},
	}, responseBody)
}

func TestChapterCreateProjectNotFound(t *testing.T) {
	tt := []struct {
		name    string
		request map[string]any
	}{
		{
			name: "should return invalid argument when project id is not found",
			request: map[string]any{
				"user": map[string]any{
					"id": testutil.ModifyOnlyUserId(),
				},
				"project": map[string]any{
					"id": "UNKNOWN_PROJECT",
				},
				"chapter": map[string]any{
					"name": "Chapter One",
				},
			},
		},
		{
			name: "should return invalid argument when user is not author of the project",
			request: map[string]any{
				"user": map[string]any{
					"id": testutil.ReadOnlyUserId(),
				},
				"project": map[string]any{
					"id": "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_API",
				},
				"chapter": map[string]any{
					"name": "Chapter One",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			router := setupChapterRouter()

			recorder := httptest.NewRecorder()
			requestBody, _ := json.Marshal(tc.request)
			req, _ := http.NewRequest("POST", "/api/chapters/create", strings.NewReader(string(requestBody)))

			router.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusBadRequest, recorder.Code)

			var responseBody map[string]any
			assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &responseBody))
			assert.Equal(t, map[string]any{
				"message": "invalid request value: " +
					"failed to create chapter: project document does not exist",
				"user":    map[string]any{},
				"project": map[string]any{},
				"chapter": map[string]any{},
			}, responseBody)
		})
	}
}

func TestChapterCreateNextChapterNotFound(t *testing.T) {
	router := setupChapterRouter()

	recorder := httptest.NewRecorder()
	requestBody, _ := json.Marshal(map[string]any{
		"user": map[string]any{
			"id": testutil.ModifyOnlyUserId(),
		},
		"project": map[string]any{
			"id": "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_API",
		},
		"chapter": map[string]any{
			"name":   "Chapter One",
			"nextId": "UNKNOWN_CHAPTER",
		},
	})
	req, _ := http.NewRequest("POST", "/api/chapters/create", strings.NewReader(string(requestBody)))

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	var responseBody map[string]any
	assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &responseBody))
	assert.Equal(t, map[string]any{
		"message": "invalid request value: " +
			"failed to create chapter: id of next chapter does not exist",
		"user":    map[string]any{},
		"project": map[string]any{},
		"chapter": map[string]any{},
	}, responseBody)
}

func TestChapterCreateDomainValidationError(t *testing.T) {
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
				"chapter": map[string]any{
					"name": "chapter name is required, but got ''",
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
					"id": "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_API",
				},
				"chapter": map[string]any{
					"name": "Chapter One",
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
			name: "Empty project id",
			request: map[string]any{
				"user": map[string]any{
					"id": testutil.ModifyOnlyUserId(),
				},
				"project": map[string]any{
					"id": "",
				},
				"chapter": map[string]any{
					"name": "Chapter One",
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
			name: "Empty chapter name",
			request: map[string]any{
				"user": map[string]any{
					"id": testutil.ModifyOnlyUserId(),
				},
				"project": map[string]any{
					"id": "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_API",
				},
				"chapter": map[string]any{
					"name": "",
				},
			},
			expectedResponse: map[string]any{
				"user":    map[string]any{},
				"project": map[string]any{},
				"chapter": map[string]any{
					"name": "chapter name is required, but got ''",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			router := setupChapterRouter()

			recorder := httptest.NewRecorder()
			requestBody, _ := json.Marshal(tc.request)
			req, _ := http.NewRequest("POST", "/api/chapters/create", strings.NewReader(string(requestBody)))

			router.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusBadRequest, recorder.Code)

			var responseBody map[string]any
			assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &responseBody))
			assert.Equal(t, map[string]any{
				"message": "invalid request value",
				"user":    tc.expectedResponse["user"],
				"project": tc.expectedResponse["project"],
				"chapter": tc.expectedResponse["chapter"],
			}, responseBody)
		})
	}
}

func TestChapterCreateInvalidRequestFormat(t *testing.T) {
	tt := []struct {
		name    string
		request string
	}{
		{
			name:    "should return error when user id is not string",
			request: `{"user": {"id":123}, "project": {"id": "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_API"}, "chapter": {"name": "Chapter One"}}`,
		},
		{
			name:    "should return error when project id is not string",
			request: `{"user": {"id": "user-id"}, "project": {"id": 123}, "chapter": {"name": "Chapter One"}}`,
		},
		{
			name:    "should return error when chapter name is not string",
			request: `{"user": {"id": "user-id"}, "project": {"id": "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_API"}, "chapter": {"name": 123}}`,
		},
		{
			name:    "should return error when user is not object",
			request: `{"user": 123, "project": {"id": "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_API"}, "chapter": {"name": "Chapter One"}}`,
		},
		{
			name:    "should return error when project is not object",
			request: `{"user": {"id": "user-id"}, "project": "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_API", "chapter": {"name": "Chapter One"}}`,
		},
		{
			name:    "should return error when chapter is not object",
			request: `{"user": {"id": "user-id"}, "project": {"id": "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_API"}, "chapter": "Chapter One"}`,
		},
		{
			name:    "should return error when request body is invalid JSON",
			request: `{"user": {"id": "user-id", "project": {"id": "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_API"}, "chapter": {"name": "Chapter One"}`,
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
			req, _ := http.NewRequest("POST", "/api/chapters/create", strings.NewReader(tc.request))

			router.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusBadRequest, recorder.Code)

			var responseBody map[string]any
			assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &responseBody))
			assert.Equal(t, map[string]any{
				"message": "invalid request format",
				"user":    map[string]any{},
				"project": map[string]any{},
				"chapter": map[string]any{},
			}, responseBody)
		})
	}
}

func setupChapterRouter() *gin.Engine {
	router := gin.Default()

	client := db.FirestoreClient()
	r := repository.NewChapterRepository(*client)
	s := service.NewChapterService(r)
	uc := usecase.NewChapterUseCase(s)
	api := api.NewChapterApi(uc)

	router.POST("/api/chapters/list", api.HandleList)
	router.POST("/api/chapters/create", api.HandleCreate)

	return router
}
