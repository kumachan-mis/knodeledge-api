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
	"github.com/kumachan-mis/knodeledge-api/internal/testutil"
	"github.com/kumachan-mis/knodeledge-api/internal/usecase"
	"github.com/stretchr/testify/assert"
)

func TestPaperFind(t *testing.T) {
	router := setupPaperRouter()

	recorder := httptest.NewRecorder()
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

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var responseBody map[string]any
	err := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
	assert.Nil(t, err)

	content := strings.Join([]string{
		"[** Introduction]",
		"This is an example project of kNODEledge.",
		"",
		"[** Section of Chapter One]",
		"Section of Chapter One. Section of Chapter One. Section of Chapter One. Section of Chapter One. Section of Chapter One. Section of Chapter One.",
		"Section of Chapter One. Section of Chapter One. Section of Chapter One. Section of Chapter One. Section of Chapter One. Section of Chapter One.",
		""}, "\n")

	assert.Equal(t, map[string]any{
		"paper": map[string]any{
			"id":      "CHAPTER_ONE",
			"content": content,
		},
	}, responseBody)
}

func TestPaperFindNotFound(t *testing.T) {
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
			recorder := httptest.NewRecorder()
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
			recorder := httptest.NewRecorder()
			requestBody, _ := json.Marshal(tc.request)
			req, _ := http.NewRequest("POST", "/api/papers/find", strings.NewReader(string(requestBody)))

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

func TestPaperFindInvalidRequestFormat(t *testing.T) {
	router := setupPaperRouter()

	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/papers/find", strings.NewReader(""))

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	var responseBody map[string]any
	err := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
	assert.Nil(t, err)

	assert.Equal(t, map[string]any{
		"message": "invalid request format",
		"user":    map[string]any{},
		"project": map[string]any{},
		"chapter": map[string]any{},
	}, responseBody)
}

func TestPaperFindInternalError(t *testing.T) {
	router := setupPaperRouter()

	recorder := httptest.NewRecorder()
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

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	var responseBody map[string]any
	err := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
	assert.Nil(t, err)

	assert.Equal(t, map[string]any{
		"message": "internal error",
	}, responseBody)
}

func TestPaperUpdate(t *testing.T) {
	router := setupPaperRouter()

	content := strings.Join([]string{
		"## Introduction",
		"This is the introduction of the paper.",
		"",
		"## What is note apps?",
		"Note apps is a web application that allows users to create, read, update, and delete notes.",
		"",
	}, "\n")

	recorder := httptest.NewRecorder()
	requestBody, _ := json.Marshal(map[string]any{
		"user": map[string]any{
			"id": testutil.ModifyOnlyUserId(),
		},
		"project": map[string]any{
			"id": "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_API",
		},
		"paper": map[string]any{
			"id":      "CHAPTER_ONE",
			"content": content,
		},
	})
	req, _ := http.NewRequest("POST", "/api/papers/update", strings.NewReader(string(requestBody)))

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var responseBody map[string]any
	err := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
	assert.Nil(t, err)

	assert.Equal(t, map[string]any{
		"paper": map[string]any{
			"id":      "CHAPTER_ONE",
			"content": content,
		},
	}, responseBody)
}

func TestPaperUpdateNotFound(t *testing.T) {
	tt := []struct {
		name      string
		userId    string
		projectId string
		paperId   string
	}{
		{
			name:      "should return error when project not found",
			userId:    testutil.ModifyOnlyUserId(),
			projectId: "UNKNOWN_PROJECT",
			paperId:   "CHAPTER_ONE",
		},
		{
			name:      "should return not found when user is not author of the project",
			userId:    testutil.ReadOnlyUserId(),
			projectId: "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_API",
			paperId:   "CHAPTER_ONE",
		},
		{
			name:      "should return error when chapter not found",
			userId:    testutil.ModifyOnlyUserId(),
			projectId: "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_API",
			paperId:   "UNKNOWN_CHAPTER",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			router := setupPaperRouter()

			recorder := httptest.NewRecorder()
			requestBody, _ := json.Marshal(map[string]any{
				"user": map[string]any{
					"id": tc.userId,
				},
				"project": map[string]any{
					"id": tc.projectId,
				},
				"paper": map[string]any{
					"id":      tc.paperId,
					"content": "This is paper content.",
				},
			})
			req, _ := http.NewRequest("POST", "/api/papers/update", strings.NewReader(string(requestBody)))

			router.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusNotFound, recorder.Code)

			var responseBody map[string]any
			err := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
			assert.Nil(t, err)

			assert.Equal(t, map[string]any{
				"message": "not found",
				"user":    map[string]any{},
				"project": map[string]any{},
				"paper":   map[string]any{},
			}, responseBody)
		})
	}
}

func TestPaperUpdateDomainValidationError(t *testing.T) {
	tooLongPaperContent := testutil.RandomString(40001)

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
					"id": "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_API",
				},
				"paper": map[string]any{
					"id":      "CHAPTER_ONE",
					"content": "This is paper content.",
				},
			},
			expectedResponse: map[string]any{
				"user": map[string]any{
					"id": "user id is required, but got ''",
				},
				"project": map[string]any{},
				"paper":   map[string]any{},
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
				"paper": map[string]any{
					"id":      "CHAPTER_ONE",
					"content": "This is paper content.",
				},
			},
			expectedResponse: map[string]any{
				"user": map[string]any{},
				"project": map[string]any{
					"id": "project id is required, but got ''",
				},
				"paper": map[string]any{},
			},
		},
		{
			name: "should return error when paper id is empty",
			request: map[string]any{
				"user": map[string]any{
					"id": testutil.ModifyOnlyUserId(),
				},
				"project": map[string]any{
					"id": "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_API",
				},
				"paper": map[string]any{
					"id":      "",
					"content": "This is paper content.",
				},
			},
			expectedResponse: map[string]any{
				"user":    map[string]any{},
				"project": map[string]any{},
				"paper": map[string]any{
					"id": "paper id is required, but got ''",
				},
			},
		},
		{
			name: "should return error when paper content is too long",
			request: map[string]any{
				"user": map[string]any{
					"id": testutil.ModifyOnlyUserId(),
				},
				"project": map[string]any{
					"id": "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_API",
				},
				"paper": map[string]any{
					"id":      "CHAPTER_ONE",
					"content": tooLongPaperContent,
				},
			},
			expectedResponse: map[string]any{
				"user":    map[string]any{},
				"project": map[string]any{},
				"paper": map[string]any{
					"content": "paper content must be less than or equal to 40000 bytes, but got 40001 bytes",
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
				"paper": map[string]any{
					"id": "paper id is required, but got ''",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			router := setupPaperRouter()

			recorder := httptest.NewRecorder()
			requestBody, _ := json.Marshal(tc.request)
			req, _ := http.NewRequest("POST", "/api/papers/update", strings.NewReader(string(requestBody)))

			router.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusBadRequest, recorder.Code)

			var responseBody map[string]any
			err := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
			assert.Nil(t, err)

			assert.Equal(t, map[string]any{
				"message": "invalid request value",
				"user":    tc.expectedResponse["user"],
				"project": tc.expectedResponse["project"],
				"paper":   tc.expectedResponse["paper"],
			}, responseBody)
		})
	}
}

func TestPaperUpdateInvalidRequestFormat(t *testing.T) {
	router := setupPaperRouter()

	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/papers/update", strings.NewReader(""))

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	var responseBody map[string]any
	err := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
	assert.Nil(t, err)

	assert.Equal(t, map[string]any{
		"message": "invalid request format",
		"user":    map[string]any{},
		"project": map[string]any{},
		"paper":   map[string]any{},
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
	router.POST("/api/papers/update", api.HandleUpdate)

	return router
}
