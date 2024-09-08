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

func TestGraphFind(t *testing.T) {
	router := setupGraphRouter()

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
		"section": map[string]any{
			"id": "SECTION_ONE",
		},
	})
	req, _ := http.NewRequest("POST", "/api/graphs/find", strings.NewReader(string(requestBody)))

	router.ServeHTTP(ecorder, req)

	assert.Equal(t, http.StatusOK, ecorder.Code)

	var responseBody map[string]any
	err := json.Unmarshal(ecorder.Body.Bytes(), &responseBody)
	assert.Nil(t, err)

	assert.Equal(t, map[string]any{
		"graph": map[string]any{
			"id":        "SECTION_ONE",
			"name":      "Introduction",
			"paragraph": "This is an example project of kNODEledge.",
		},
	}, responseBody)
}

func TestGraphFindProjectOrChapterOrSectionNotFound(t *testing.T) {
	router := setupGraphRouter()

	tt := []struct {
		name    string
		user    string
		project string
		chapter string
		section string
	}{
		{
			name:    "should return error when project not found",
			user:    testutil.ReadOnlyUserId(),
			project: "UNKNOWN_PROJECT",
			chapter: "CHAPTER_ONE",
			section: "SECTION_ONE",
		},
		{
			name:    "should return not found when user is not author of the project",
			user:    testutil.ModifyOnlyUserId(),
			project: "PROJECT_WITH_DESCRIPTION",
			chapter: "CHAPTER_ONE",
			section: "SECTION_ONE",
		},
		{
			name:    "should return error when chapter not found",
			user:    testutil.ReadOnlyUserId(),
			project: "PROJECT_WITH_DESCRIPTION",
			chapter: "UNKNOWN_CHAPTER",
			section: "SECTION_ONE",
		},
		{
			name:    "should return error when section not found",
			user:    testutil.ReadOnlyUserId(),
			project: "PROJECT_WITH_DESCRIPTION",
			chapter: "CHAPTER_ONE",
			section: "UNKNOWN_SECTION",
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
				"section": map[string]any{
					"id": tc.section,
				},
			})
			req, _ := http.NewRequest("POST", "/api/graphs/find", strings.NewReader(string(requestBody)))

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
				"section": map[string]any{},
			}, responseBody)
		})
	}
}

func TestGraphFindDomainValidationError(t *testing.T) {
	router := setupGraphRouter()

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
				"section": map[string]any{
					"id": "SECTION_ONE",
				},
			},
			expectedResponse: map[string]any{
				"user": map[string]any{
					"id": "user id is required, but got ''",
				},
				"project": map[string]any{},
				"chapter": map[string]any{},
				"section": map[string]any{},
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
				"section": map[string]any{
					"id": "SECTION_ONE",
				},
			},
			expectedResponse: map[string]any{
				"user": map[string]any{},
				"project": map[string]any{
					"id": "project id is required, but got ''",
				},
				"chapter": map[string]any{},
				"section": map[string]any{},
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
				"section": map[string]any{
					"id": "SECTION_ONE",
				},
			},
			expectedResponse: map[string]any{
				"user":    map[string]any{},
				"project": map[string]any{},
				"chapter": map[string]any{
					"id": "chapter id is required, but got ''",
				},
				"section": map[string]any{},
			},
		},
		{
			name: "should return error when section id is empty",
			request: map[string]any{
				"user": map[string]any{
					"id": testutil.ReadOnlyUserId(),
				},
				"project": map[string]any{
					"id": "PROJECT_WITHOUT_DESCRIPTION",
				},
				"chapter": map[string]any{
					"id": "CHAPTER_ONE",
				},
				"section": map[string]any{
					"id": "",
				},
			},
			expectedResponse: map[string]any{
				"user":    map[string]any{},
				"project": map[string]any{},
				"chapter": map[string]any{},
				"section": map[string]any{
					"id": "section id is required, but got ''",
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
				"section": map[string]any{
					"id": "section id is required, but got ''",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ecorder := httptest.NewRecorder()
			requestBody, _ := json.Marshal(tc.request)
			req, _ := http.NewRequest("POST", "/api/graphs/find", strings.NewReader(string(requestBody)))

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
				"section": tc.expectedResponse["section"],
			}, responseBody)
		})
	}
}

func TestGraphFindInvalidRequestFormat(t *testing.T) {
	router := setupGraphRouter()

	ecorder := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/graphs/find", strings.NewReader(""))

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
		"section": map[string]any{},
	}, responseBody)
}

func TestGraphFindInternalError(t *testing.T) {
	router := setupGraphRouter()

	ecorder := httptest.NewRecorder()
	requestBody, _ := json.Marshal(map[string]any{
		"user": map[string]any{
			"id": testutil.ErrorUserId(7),
		},
		"project": map[string]any{
			"id": "PROJECT_WITH_INVALID_GRAPH_PARAGRAPH",
		},
		"chapter": map[string]any{
			"id": "CHAPTER_WITH_INVALID_GRAPH_PARAGRAPH",
		},
		"section": map[string]any{
			"id": "SECTION_WITH_INVALID_GRAPH_PARAGRAPH",
		},
	})
	req, _ := http.NewRequest("POST", "/api/graphs/find", strings.NewReader(string(requestBody)))

	router.ServeHTTP(ecorder, req)

	assert.Equal(t, http.StatusInternalServerError, ecorder.Code)

	var responseBody map[string]any
	err := json.Unmarshal(ecorder.Body.Bytes(), &responseBody)
	assert.Nil(t, err)

	assert.Equal(t, map[string]any{
		"message": "internal error",
	}, responseBody)
}

func TestGraphSectionalize(t *testing.T) {
	maxLengthSectionName := testutil.RandomString(100)
	maxLengthSectionContent := testutil.RandomString(40000)

	router := setupGraphRouter()

	ecorder := httptest.NewRecorder()
	requestBody, _ := json.Marshal(map[string]any{
		"user": map[string]any{
			"id": testutil.ModifyOnlyUserId(),
		},
		"project": map[string]any{
			"id": "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_API",
		},
		"chapter": map[string]any{
			"id": "CHAPTER_TWO",
		},
		"sections": []map[string]any{
			{
				"name":    "Section One",
				"content": "Content of Section One",
			},
			{
				"name":    maxLengthSectionName,
				"content": maxLengthSectionContent,
			},
		},
	})
	req, _ := http.NewRequest("POST", "/api/graphs/sectionalize", strings.NewReader(string(requestBody)))

	router.ServeHTTP(ecorder, req)

	assert.Equal(t, http.StatusCreated, ecorder.Code)

	var responseBody map[string]any
	err := json.Unmarshal(ecorder.Body.Bytes(), &responseBody)
	assert.Nil(t, err)

	//graphId is generated by firestore and it's not predictable
	graphId1 := responseBody["graphs"].([]any)[0].(map[string]any)["id"]
	assert.NotEmpty(t, graphId1)
	graphId2 := responseBody["graphs"].([]any)[1].(map[string]any)["id"]
	assert.NotEmpty(t, graphId2)

	assert.Equal(t, map[string]any{
		"graphs": []any{
			map[string]any{
				"id":        graphId1,
				"name":      "Section One",
				"paragraph": "Content of Section One",
			},
			map[string]any{
				"id":        graphId2,
				"name":      maxLengthSectionName,
				"paragraph": maxLengthSectionContent,
			},
		},
	}, responseBody)
}

func TestGraphSectionalizeProjectOrChapterNotFound(t *testing.T) {
	tt := []struct {
		name      string
		userId    string
		projectId string
		chapterId string
	}{
		{
			name:      "should return error when project not found",
			userId:    testutil.ModifyOnlyUserId(),
			projectId: "UNKNOWN_PROJECT",
			chapterId: "CHAPTER_ONE",
		},
		{
			name:      "should return not found when user is not author of the project",
			userId:    testutil.ReadOnlyUserId(),
			projectId: "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_API",
			chapterId: "CHAPTER_ONE",
		},
		{
			name:      "should return error when chapter not found",
			userId:    testutil.ModifyOnlyUserId(),
			projectId: "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_API",
			chapterId: "UNKNOWN_CHAPTER",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			router := setupGraphRouter()

			ecorder := httptest.NewRecorder()
			requestBody, _ := json.Marshal(map[string]any{
				"user": map[string]any{
					"id": tc.userId,
				},
				"project": map[string]any{
					"id": tc.projectId,
				},
				"chapter": map[string]any{
					"id": tc.chapterId,
				},
				"sections": []map[string]any{
					{
						"name":    "Section One",
						"content": "Content of Section One",
					},
				},
			})
			req, _ := http.NewRequest("POST", "/api/graphs/sectionalize", strings.NewReader(string(requestBody)))

			router.ServeHTTP(ecorder, req)

			assert.Equal(t, http.StatusNotFound, ecorder.Code)

			var responseBody map[string]any
			err := json.Unmarshal(ecorder.Body.Bytes(), &responseBody)

			assert.Nil(t, err)

			assert.Equal(t, map[string]any{
				"message":  "not found",
				"user":     map[string]any{},
				"project":  map[string]any{},
				"chapter":  map[string]any{},
				"sections": map[string]any{},
			}, responseBody)
		})
	}
}

func TestGraphSectionalizeGraphAlreadyExists(t *testing.T) {
	router := setupGraphRouter()

	ecorder := httptest.NewRecorder()
	requestBody, _ := json.Marshal(map[string]any{
		"user": map[string]any{
			"id": testutil.ModifyOnlyUserId(),
		},
		"project": map[string]any{
			"id": "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_API",
		},
		"chapter": map[string]any{
			"id": "CHAPTER_ONE",
		},
		"sections": []map[string]any{
			{
				"name":    "Section One",
				"content": "Content of Section One",
			},
		},
	})
	req, _ := http.NewRequest("POST", "/api/graphs/sectionalize", strings.NewReader(string(requestBody)))

	router.ServeHTTP(ecorder, req)

	assert.Equal(t, http.StatusBadRequest, ecorder.Code)

	var responseBody map[string]any
	err := json.Unmarshal(ecorder.Body.Bytes(), &responseBody)

	assert.Nil(t, err)

	assert.Equal(t, map[string]any{
		"message":  "invalid request value: failed to sectionalize into graphs: graph already exists",
		"user":     map[string]any{},
		"project":  map[string]any{},
		"chapter":  map[string]any{},
		"sections": map[string]any{},
	}, responseBody)
}

func TestGraphSectionalizeDomainValidationError(t *testing.T) {
	tooLongSectionName := testutil.RandomString(101)
	tooLongSectionContent := testutil.RandomString(40001)
	tooLongSections := make([]any, 21)
	for i := 0; i < 21; i++ {
		tooLongSections[i] = map[string]any{
			"name":    "Section",
			"content": "Content of Section",
		}
	}
	tooLongSectionErrors := make([]any, 21)
	for i := 0; i < 21; i++ {
		tooLongSectionErrors[i] = map[string]any{}
	}

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
					"id": "chapter id is required, but got ''",
				},
				"sections": map[string]any{
					"message": "sections are required, but got []",
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
					"id": "CHAPTER_ONE",
				},
				"sections": []any{
					map[string]any{
						"name":    "Section One",
						"content": "Content of Section One",
					},
				},
			},
			expectedResponse: map[string]any{
				"user": map[string]any{
					"id": "user id is required, but got ''",
				},
				"project": map[string]any{},
				"chapter": map[string]any{},
				"sections": map[string]any{
					"items": []any{
						map[string]any{},
					},
				},
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
					"id": "CHAPTER_ONE",
				},
				"sections": []any{
					map[string]any{
						"name":    "Section One",
						"content": "Content of Section One",
					},
				},
			},
			expectedResponse: map[string]any{
				"user": map[string]any{},
				"project": map[string]any{
					"id": "project id is required, but got ''",
				},
				"chapter": map[string]any{},
				"sections": map[string]any{
					"items": []any{
						map[string]any{},
					},
				},
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
					"id": "",
				},
				"sections": []any{
					map[string]any{
						"name":    "Section One",
						"content": "Content of Section One",
					},
				},
			},
			expectedResponse: map[string]any{
				"user":    map[string]any{},
				"project": map[string]any{},
				"chapter": map[string]any{
					"id": "chapter id is required, but got ''",
				},
				"sections": map[string]any{
					"items": []any{
						map[string]any{},
					},
				},
			},
		},
		{
			name: "should return error when sections are empty",
			request: map[string]any{
				"user": map[string]any{
					"id": testutil.ModifyOnlyUserId(),
				},
				"project": map[string]any{
					"id": "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_API",
				},
				"chapter": map[string]any{
					"id": "CHAPTER_ONE",
				},
				"sections": []any{},
			},
			expectedResponse: map[string]any{
				"user":    map[string]any{},
				"project": map[string]any{},
				"chapter": map[string]any{},
				"sections": map[string]any{
					"message": "sections are required, but got []",
				},
			},
		},
		{
			name: "should return error when sections are too many",
			request: map[string]any{
				"user": map[string]any{
					"id": testutil.ModifyOnlyUserId(),
				},
				"project": map[string]any{
					"id": "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_API",
				},
				"chapter": map[string]any{
					"id": "CHAPTER_ONE",
				},
				"sections": tooLongSections,
			},
			expectedResponse: map[string]any{
				"user":    map[string]any{},
				"project": map[string]any{},
				"chapter": map[string]any{},
				"sections": map[string]any{
					"message": "sections length must be less than or equal to 20, but got 21",
					"items":   tooLongSectionErrors,
				},
			},
		},
		{
			name: "should return error when section name is empty",
			request: map[string]any{
				"user": map[string]any{
					"id": testutil.ModifyOnlyUserId(),
				},
				"project": map[string]any{
					"id": "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_API",
				},
				"chapter": map[string]any{
					"id": "CHAPTER_ONE",
				},
				"sections": []map[string]any{
					{
						"name":    "",
						"content": "Content of Section",
					},
				},
			},
			expectedResponse: map[string]any{
				"user":    map[string]any{},
				"project": map[string]any{},
				"chapter": map[string]any{},
				"sections": map[string]any{
					"items": []any{
						map[string]any{
							"name": "section name is required, but got ''",
						},
					},
				},
			},
		},
		{
			name: "should return error when section name is too long",
			request: map[string]any{
				"user": map[string]any{
					"id": testutil.ModifyOnlyUserId(),
				},
				"project": map[string]any{
					"id": "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_API",
				},
				"chapter": map[string]any{
					"id": "CHAPTER_ONE",
				},
				"sections": []any{
					map[string]any{
						"name":    tooLongSectionName,
						"content": "Content of Section",
					},
					map[string]any{
						"name":    "Section",
						"content": "Content of Section",
					},
				},
			},
			expectedResponse: map[string]any{
				"user":    map[string]any{},
				"project": map[string]any{},
				"chapter": map[string]any{},
				"sections": map[string]any{
					"items": []any{
						map[string]any{
							"name": fmt.Sprintf("section name cannot be longer than 100 characters, but got '%v'",
								tooLongSectionName),
						},
						map[string]any{},
					},
				},
			},
		},
		{
			name: "should return error when section content is too long",
			request: map[string]any{
				"user": map[string]any{
					"id": testutil.ModifyOnlyUserId(),
				},
				"project": map[string]any{
					"id": "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_API",
				},
				"chapter": map[string]any{
					"id": "CHAPTER_ONE",
				},
				"sections": []any{
					map[string]any{
						"name":    "Section One",
						"content": "Content of Section",
					},
					map[string]any{
						"name":    "Section Two",
						"content": tooLongSectionContent,
					},
				},
			},
			expectedResponse: map[string]any{
				"user":    map[string]any{},
				"project": map[string]any{},
				"chapter": map[string]any{},
				"sections": map[string]any{
					"items": []any{
						map[string]any{},
						map[string]any{
							"content": "section content must be less than or equal to 40000 bytes, but got 40001 bytes",
						},
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			router := setupGraphRouter()

			recorder := httptest.NewRecorder()
			requestBody, _ := json.Marshal(tc.request)
			req, _ := http.NewRequest("POST", "/api/graphs/sectionalize", strings.NewReader(string(requestBody)))

			router.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusBadRequest, recorder.Code)

			var responseBody map[string]any
			assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &responseBody))
			assert.Equal(t, map[string]any{
				"message":  "invalid request value",
				"user":     tc.expectedResponse["user"],
				"project":  tc.expectedResponse["project"],
				"chapter":  tc.expectedResponse["chapter"],
				"sections": tc.expectedResponse["sections"],
			}, responseBody)
		})
	}
}

func TestGraphSectionalizeInvalidRequestFormat(t *testing.T) {
	router := setupGraphRouter()

	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/graphs/sectionalize", strings.NewReader(""))

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	var responseBody map[string]any
	assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &responseBody))
	assert.Equal(t, map[string]any{
		"message":  "invalid request format",
		"user":     map[string]any{},
		"project":  map[string]any{},
		"chapter":  map[string]any{},
		"sections": map[string]any{},
	}, responseBody)
}

func setupGraphRouter() *gin.Engine {
	router := gin.Default()
	client := db.FirestoreClient()
	r := repository.NewGraphRepository(*client)
	cr := repository.NewChapterRepository(*client)
	s := service.NewGraphService(r, cr)

	uc := usecase.NewGraphUseCase(s)
	api := api.NewGraphApi(uc)

	router.POST("/api/graphs/find", api.HandleFind)
	router.POST("/api/graphs/update", api.HandleUpdate)
	router.POST("/api/graphs/sectionalize", api.HandleSectionalize)

	return router
}
