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
	"github.com/kumachan-mis/knodeledge-api/internal/testutil"
	"github.com/kumachan-mis/knodeledge-api/internal/usecase"
	mock_middleware "github.com/kumachan-mis/knodeledge-api/mock/middleware"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestChapterList(t *testing.T) {
	router := setupChapterRouter(t)

	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/chapters/list", nil)
	query := req.URL.Query()
	query.Add("userId", testutil.ReadOnlyUserId())
	query.Add("projectId", "PROJECT_WITHOUT_DESCRIPTION")
	req.URL.RawQuery = query.Encode()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var responseBody map[string]any
	assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &responseBody))
	assert.Equal(t, map[string]any{
		"chapters": []any{
			map[string]any{
				"id":     "CHAPTER_ONE",
				"name":   "Chapter One",
				"number": float64(1), // json.Unmarshal converts number to float64
				"sections": []any{
					map[string]any{
						"id":   "SECTION_ONE",
						"name": "Introduction",
					},
					map[string]any{
						"id":   "SECTION_TWO",
						"name": "Section of Chapter One",
					},
				},
			},
			map[string]any{
				"id":       "CHAPTER_TWO",
				"name":     "Chapter Two",
				"number":   float64(2), // json.Unmarshal converts number to float64
				"sections": []any{},
			},
		},
	}, responseBody)
}

func TestChapterListEmpty(t *testing.T) {
	router := setupChapterRouter(t)

	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/chapters/list", nil)
	query := req.URL.Query()
	query.Add("userId", testutil.ReadOnlyUserId())
	query.Add("projectId", "PROJECT_WITH_DESCRIPTION")
	req.URL.RawQuery = query.Encode()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var responseBody map[string]any
	assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &responseBody))
	assert.Equal(t, responseBody, map[string]any{
		"chapters": []any{},
	})
}

func TestChapterListNotFound(t *testing.T) {
	tt := []struct {
		name  string
		query map[string]string
	}{
		{
			name: "should return not found when project is not found",
			query: map[string]string{
				"userId":    testutil.ReadOnlyUserId(),
				"projectId": "UNKNOWN_PROJECT",
			},
		},
		{
			name: "should return not found when user is not author of the project",
			query: map[string]string{
				"userId":    testutil.ModifyOnlyUserId(),
				"projectId": "PROJECT_WITHOUT_DESCRIPTION",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			router := setupChapterRouter(t)

			recorder := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/chapters/list", nil)
			query := req.URL.Query()
			for key, value := range tc.query {
				query.Add(key, value)
			}
			req.URL.RawQuery = query.Encode()

			router.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusNotFound, recorder.Code)

			var responseBody map[string]any
			assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &responseBody))
			assert.Equal(t, map[string]any{
				"message": "not found",
			}, responseBody)
		})
	}
}

func TestChapterListDomainValidationError(t *testing.T) {
	tt := []struct {
		name             string
		query            map[string]string
		expectedResponse map[string]any
	}{
		{
			name: "should return error when user id is empty",
			query: map[string]string{
				"userId":    "",
				"projectId": "PROJECT_WITHOUT_DESCRIPTION",
			},
			expectedResponse: map[string]any{
				"userId": "user id is required, but got ''",
			},
		},
		{
			name: "should return error when project id is empty",
			query: map[string]string{
				"userId":    testutil.ReadOnlyUserId(),
				"projectId": "",
			},
			expectedResponse: map[string]any{
				"projectId": "project id is required, but got ''",
			},
		},
		{
			name:  "should return error when empty parameter is passed",
			query: map[string]string{},
			expectedResponse: map[string]any{
				"userId":    "user id is required, but got ''",
				"projectId": "project id is required, but got ''",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			router := setupChapterRouter(t)

			recorder := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/chapters/list", nil)
			query := req.URL.Query()
			for key, value := range tc.query {
				query.Add(key, value)
			}
			req.URL.RawQuery = query.Encode()

			router.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusBadRequest, recorder.Code)

			var responseBody map[string]any
			assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &responseBody))

			expectedResponse := tc.expectedResponse
			expectedResponse["message"] = "invalid request value"
			assert.Equal(t, expectedResponse, responseBody)
		})
	}
}

func TestChapterListInternalError(t *testing.T) {
	router := setupChapterRouter(t)

	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/chapters/list", nil)
	query := req.URL.Query()
	query.Add("userId", testutil.ErrorUserId(2))
	query.Add("projectId", "PROJECT_WITH_INVALID_CHAPTER_NAME")
	req.URL.RawQuery = query.Encode()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	var responseBody map[string]any
	assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &responseBody))
	assert.Equal(t, map[string]any{
		"message": "internal error",
	}, responseBody)
}

func TestChapterCreate(t *testing.T) {
	router := setupChapterRouter(t)

	recorder := httptest.NewRecorder()
	requestBody, _ := json.Marshal(map[string]any{
		"user": map[string]any{
			"id": testutil.ModifyOnlyUserId(),
		},
		"project": map[string]any{
			"id": "PROJECT_WITH_DESCRIPTION_TO_UPDATE_FROM_API",
		},
		"chapter": map[string]any{
			"name":   "Chapter One",
			"number": 1,
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
			"number":   float64(1), // json.Unmarshal converts number to float64
			"sections": []any{},
		},
	}, responseBody)
}

func TestChapterCreateNotFound(t *testing.T) {
	tt := []struct {
		name    string
		request map[string]any
	}{
		{
			name: "should return not found when project is not found",
			request: map[string]any{
				"user": map[string]any{
					"id": testutil.ModifyOnlyUserId(),
				},
				"project": map[string]any{
					"id": "UNKNOWN_PROJECT",
				},
				"chapter": map[string]any{
					"name":   "Chapter One",
					"number": 1,
				},
			},
		},
		{
			name: "should return not found when user is not author of the project",
			request: map[string]any{
				"user": map[string]any{
					"id": testutil.ReadOnlyUserId(),
				},
				"project": map[string]any{
					"id": "PROJECT_WITH_DESCRIPTION_TO_UPDATE_FROM_API",
				},
				"chapter": map[string]any{
					"name":   "Chapter One",
					"number": 1,
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			router := setupChapterRouter(t)

			recorder := httptest.NewRecorder()
			requestBody, _ := json.Marshal(tc.request)
			req, _ := http.NewRequest("POST", "/api/chapters/create", strings.NewReader(string(requestBody)))

			router.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusNotFound, recorder.Code)

			var responseBody map[string]any
			assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &responseBody))
			assert.Equal(t, map[string]any{
				"message": "not found",
				"user":    map[string]any{},
				"project": map[string]any{},
				"chapter": map[string]any{},
			}, responseBody)
		})
	}
}

func TestChapterCreateTooLargeChapterNumber(t *testing.T) {
	router := setupChapterRouter(t)

	recorder := httptest.NewRecorder()
	requestBody, _ := json.Marshal(map[string]any{
		"user": map[string]any{
			"id": testutil.ModifyOnlyUserId(),
		},
		"project": map[string]any{
			"id": "PROJECT_WITH_DESCRIPTION_TO_UPDATE_FROM_API",
		},
		"chapter": map[string]any{
			"name":   "Chapter Ninety-Nine",
			"number": 99,
		},
	})
	req, _ := http.NewRequest("POST", "/api/chapters/create", strings.NewReader(string(requestBody)))

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	var responseBody map[string]any
	assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &responseBody))
	assert.Equal(t, map[string]any{
		"message": "invalid request value: " +
			"failed to create chapter: chapter number is too large",
		"user":    map[string]any{},
		"project": map[string]any{},
		"chapter": map[string]any{},
	}, responseBody)
}

func TestChapterCreateDomainValidationError(t *testing.T) {
	tooLongChapterName := testutil.RandomString(101)

	tt := []struct {
		name             string
		request          map[string]any
		expectedResponse map[string]any
	}{
		{
			name:    "should return error when request is empty",
			request: map[string]any{},
			expectedResponse: map[string]any{
				"user": map[string]any{
					"id": "user id is required, but got ''",
				},
				"project": map[string]any{
					"id": "project id is required, but got ''",
				},
				"chapter": map[string]any{
					"name":   "chapter name is required, but got ''",
					"number": "chapter number must be greater than 0, but got 0",
				},
			},
		},
		{
			name: "should return error when user id is empty",
			request: map[string]any{
				"user": map[string]any{
					"id": "",
				},
				"project": map[string]any{
					"id": "PROJECT_WITH_DESCRIPTION_TO_UPDATE_FROM_API",
				},
				"chapter": map[string]any{
					"name":   "Chapter One",
					"number": 1,
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
					"id": testutil.ModifyOnlyUserId(),
				},
				"project": map[string]any{
					"id": "",
				},
				"chapter": map[string]any{
					"name":   "Chapter One",
					"number": 1,
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
			name: "should return error when chapter name is empty",
			request: map[string]any{
				"user": map[string]any{
					"id": testutil.ModifyOnlyUserId(),
				},
				"project": map[string]any{
					"id": "PROJECT_WITH_DESCRIPTION_TO_UPDATE_FROM_API",
				},
				"chapter": map[string]any{
					"name":   "",
					"number": 1,
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
		{
			name: "should return error when chapter name is too long",
			request: map[string]any{
				"user": map[string]any{
					"id": testutil.ModifyOnlyUserId(),
				},
				"project": map[string]any{
					"id": "PROJECT_WITH_DESCRIPTION_TO_UPDATE_FROM_API",
				},
				"chapter": map[string]any{
					"name":   tooLongChapterName,
					"number": 1,
				},
			},
			expectedResponse: map[string]any{
				"user":    map[string]any{},
				"project": map[string]any{},
				"chapter": map[string]any{
					"name": fmt.Sprintf("chapter name cannot be longer than 100 characters, but got '%v'",
						tooLongChapterName),
				},
			},
		},
		{
			name: "should return error when chapter number is zero",
			request: map[string]any{
				"user": map[string]any{
					"id": testutil.ModifyOnlyUserId(),
				},
				"project": map[string]any{
					"id": "PROJECT_WITH_DESCRIPTION_TO_UPDATE_FROM_API",
				},
				"chapter": map[string]any{
					"name":   "Chapter One",
					"number": 0,
				},
			},
			expectedResponse: map[string]any{
				"user":    map[string]any{},
				"project": map[string]any{},
				"chapter": map[string]any{
					"number": "chapter number must be greater than 0, but got 0",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			router := setupChapterRouter(t)

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
	router := setupChapterRouter(t)

	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/chapters/create", strings.NewReader(""))

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
}

func TestChapterUpdate(t *testing.T) {
	router := setupChapterRouter(t)

	recorder := httptest.NewRecorder()
	requestBody, _ := json.Marshal(map[string]any{
		"user": map[string]any{
			"id": testutil.ModifyOnlyUserId(),
		},
		"project": map[string]any{
			"id": "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_API",
		},
		"chapter": map[string]any{
			"id":     "CHAPTER_ONE",
			"name":   "Updated Chapter One",
			"number": 1,
		},
	})
	req, _ := http.NewRequest("POST", "/api/chapters/update", strings.NewReader(string(requestBody)))

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var responseBody map[string]any
	assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &responseBody))

	assert.Equal(t, map[string]any{
		"chapter": map[string]any{
			"id":       "CHAPTER_ONE",
			"name":     "Updated Chapter One",
			"number":   float64(1), // json.Unmarshal converts number to float64
			"sections": []any{},
		},
	}, responseBody)
}

func TestChapterUpdateNotFound(t *testing.T) {
	tt := []struct {
		name    string
		request map[string]any
	}{
		{
			name: "should return not found when project is not found",
			request: map[string]any{
				"user": map[string]any{
					"id": testutil.ModifyOnlyUserId(),
				},
				"project": map[string]any{
					"id": "UNKNOWN_PROJECT",
				},
				"chapter": map[string]any{
					"id":     "CHAPTER_ONE",
					"name":   "Updated Chapter One",
					"number": 1,
				},
			},
		},
		{
			name: "should return not found when user is not author of the project",
			request: map[string]any{
				"user": map[string]any{
					"id": testutil.ReadOnlyUserId(),
				},
				"project": map[string]any{
					"id": "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_API",
				},
				"chapter": map[string]any{
					"id":     "CHAPTER_ONE",
					"name":   "Updated Chapter One",
					"number": 1,
				},
			},
		},
		{
			name: "should return not found when chapter is not found",
			request: map[string]any{
				"user": map[string]any{
					"id": testutil.ModifyOnlyUserId(),
				},
				"project": map[string]any{
					"id": "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_API",
				},
				"chapter": map[string]any{
					"id":     "UNKNOWN_CHAPTER",
					"name":   "Updated Chapter One",
					"number": 1,
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			router := setupChapterRouter(t)

			recorder := httptest.NewRecorder()
			requestBody, _ := json.Marshal(tc.request)
			req, _ := http.NewRequest("POST", "/api/chapters/update", strings.NewReader(string(requestBody)))

			router.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusNotFound, recorder.Code)

			var responseBody map[string]any
			assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &responseBody))
			assert.Equal(t, map[string]any{
				"message": "not found",
				"user":    map[string]any{},
				"project": map[string]any{},
				"chapter": map[string]any{},
			}, responseBody)
		})
	}
}

func TestChapterUpdateTooLargeChapterNumber(t *testing.T) {
	router := setupChapterRouter(t)

	recorder := httptest.NewRecorder()
	requestBody, _ := json.Marshal(map[string]any{
		"user": map[string]any{
			"id": testutil.ModifyOnlyUserId(),
		},
		"project": map[string]any{
			"id": "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_API",
		},
		"chapter": map[string]any{
			"id":     "CHAPTER_ONE",
			"name":   "Updated Chapter One",
			"number": 3,
		},
	})
	req, _ := http.NewRequest("POST", "/api/chapters/update", strings.NewReader(string(requestBody)))

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	var responseBody map[string]any
	assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &responseBody))
	assert.Equal(t, map[string]any{
		"message": "invalid request value: " +
			"failed to update chapter: chapter number is too large",
		"user":    map[string]any{},
		"project": map[string]any{},
		"chapter": map[string]any{},
	}, responseBody)
}

func TestChapterUpdateDomainValidationError(t *testing.T) {
	tt := []struct {
		name             string
		request          map[string]any
		expectedResponse map[string]any
	}{
		{
			name:    "should return error when request is empty",
			request: map[string]any{},
			expectedResponse: map[string]any{
				"user": map[string]any{
					"id": "user id is required, but got ''",
				},
				"project": map[string]any{
					"id": "project id is required, but got ''",
				},
				"chapter": map[string]any{
					"id":     "chapter id is required, but got ''",
					"name":   "chapter name is required, but got ''",
					"number": "chapter number must be greater than 0, but got 0",
				},
			},
		},
		{
			name: "should return error when user id is empty",
			request: map[string]any{
				"user": map[string]any{
					"id": "",
				},
				"project": map[string]any{
					"id": "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_API",
				},
				"chapter": map[string]any{
					"id":     "CHAPTER_ONE",
					"name":   "Updated Chapter One",
					"number": 1,
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
					"id": testutil.ModifyOnlyUserId(),
				},
				"project": map[string]any{
					"id": "",
				},
				"chapter": map[string]any{
					"id":     "CHAPTER_ONE",
					"name":   "Updated Chapter One",
					"number": 1,
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
					"id": testutil.ModifyOnlyUserId(),
				},
				"project": map[string]any{
					"id": "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_API",
				},
				"chapter": map[string]any{
					"id":     "",
					"name":   "Updated Chapter One",
					"number": 1,
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
			name: "should return error when chapter name is empty",
			request: map[string]any{
				"user": map[string]any{
					"id": testutil.ModifyOnlyUserId(),
				},
				"project": map[string]any{
					"id": "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_API",
				},
				"chapter": map[string]any{
					"id":     "CHAPTER_ONE",
					"name":   "",
					"number": 1,
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
		{
			name: "should return error when chapter number is zero",
			request: map[string]any{
				"user": map[string]any{
					"id": testutil.ModifyOnlyUserId(),
				},
				"project": map[string]any{
					"id": "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_API",
				},
				"chapter": map[string]any{
					"id":     "CHAPTER_ONE",
					"name":   "Updated Chapter One",
					"number": 0,
				},
			},
			expectedResponse: map[string]any{
				"user":    map[string]any{},
				"project": map[string]any{},
				"chapter": map[string]any{
					"number": "chapter number must be greater than 0, but got 0",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			router := setupChapterRouter(t)

			recorder := httptest.NewRecorder()
			requestBody, _ := json.Marshal(tc.request)
			req, _ := http.NewRequest("POST", "/api/chapters/update", strings.NewReader(string(requestBody)))

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

func TestChapterUpdateInvalidRequestFormat(t *testing.T) {
	router := setupChapterRouter(t)

	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/chapters/update", strings.NewReader(""))

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
}

func TestChapterDelete(t *testing.T) {
	router := setupChapterRouter(t)

	recorder := httptest.NewRecorder()
	requestBody, _ := json.Marshal(map[string]any{
		"user": map[string]any{
			"id": testutil.ModifyOnlyUserId(),
		},
		"project": map[string]any{
			"id": "PROJECT_WITHOUT_DESCRIPTION_TO_DELETE_FROM_API",
		},
		"chapter": map[string]any{
			"id": "CHAPTER_TWO",
		},
	})
	req, _ := http.NewRequest("POST", "/api/chapters/delete", strings.NewReader(string(requestBody)))

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusNoContent, recorder.Code)
}

func TestChapterDeleteNotFound(t *testing.T) {
	tt := []struct {
		name    string
		request map[string]any
	}{
		{
			name: "should return not found when project is not found",
			request: map[string]any{
				"user": map[string]any{
					"id": testutil.ModifyOnlyUserId(),
				},
				"project": map[string]any{
					"id": "UNKNOWN_PROJECT",
				},
				"chapter": map[string]any{
					"id": "CHAPTER_ONE",
				},
			},
		},
		{
			name: "should return not found when user is not author of the project",
			request: map[string]any{
				"user": map[string]any{
					"id": testutil.ReadOnlyUserId(),
				},
				"project": map[string]any{
					"id": "PROJECT_WITHOUT_DESCRIPTION_TO_DELETE_FROM_API",
				},
				"chapter": map[string]any{
					"id": "CHAPTER_TWO",
				},
			},
		},
		{
			name: "should return not found when chapter is not found",
			request: map[string]any{
				"user": map[string]any{
					"id": testutil.ModifyOnlyUserId(),
				},
				"project": map[string]any{
					"id": "PROJECT_WITHOUT_DESCRIPTION_TO_DELETE_FROM_API",
				},
				"chapter": map[string]any{
					"id": "UNKNOWN_CHAPTER",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			router := setupChapterRouter(t)

			recorder := httptest.NewRecorder()
			requestBody, _ := json.Marshal(tc.request)
			req, _ := http.NewRequest("POST", "/api/chapters/delete", strings.NewReader(string(requestBody)))

			router.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusNotFound, recorder.Code)

			var responseBody map[string]any
			err := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
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

func TestChapterDeleteDomainValidationError(t *testing.T) {
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
			router := setupChapterRouter(t)

			recorder := httptest.NewRecorder()
			requestBody, _ := json.Marshal(tc.request)
			req, _ := http.NewRequest("POST", "/api/chapters/delete", strings.NewReader(string(requestBody)))

			router.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusBadRequest, recorder.Code)

			var responseBody map[string]any
			err := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
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

func setupChapterRouter(t *testing.T) *gin.Engine {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	router := gin.Default()

	client := db.FirestoreClient()
	r := repository.NewChapterRepository(*client)
	pr := repository.NewPaperRepository(*client)
	s := service.NewChapterService(r, pr)

	v := mock_middleware.NewMockUserVerifier(ctrl)
	v.EXPECT().
		Verify(gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	uc := usecase.NewChapterUseCase(s)
	api := api.NewChaptersApi(v, uc)

	router.GET("/api/chapters/list", api.ChaptersList)
	router.POST("/api/chapters/create", api.ChaptersCreate)
	router.POST("/api/chapters/update", api.ChaptersUpdate)
	router.POST("/api/chapters/delete", api.ChaptersDelete)

	return router
}
