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
			"children": []any{
				map[string]any{
					"name":        "Background",
					"relation":    "part of",
					"description": "This is background part.",
					"children": []any{
						map[string]any{
							"name":        "IT in Education",
							"relation":    "one of",
							"description": "This is IT in Education part.",
							"children":    []any{},
						},
					},
				},
				map[string]any{
					"name":        "Motivation",
					"relation":    "part of",
					"description": "This is motivation part.",
					"children":    []any{},
				},
				map[string]any{
					"name":        "Literature Review",
					"relation":    "part of",
					"description": "This is literature review part.",
					"children":    []any{},
				},
			},
		},
	}, responseBody)
}

func TestGraphFindNotFound(t *testing.T) {
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

func TestGraphUpdate(t *testing.T) {
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
		"graph": map[string]any{
			"id":        "SECTION_TWO",
			"paragraph": "Updated paragraph content.",
			"children": []any{
				map[string]any{
					"name":        "Background",
					"relation":    "part of",
					"description": "This is background part.",
					"children": []any{
						map[string]any{
							"name":        "IT in Education",
							"relation":    "one of",
							"description": "This is IT in Education part.",
							"children":    []any{},
						},
						map[string]any{
							"name":        "IT in Business",
							"relation":    "one of",
							"description": "This is IT in Business part.",
							"children":    []any{},
						},
					},
				},
				map[string]any{
					"name":        "Motivation",
					"relation":    "part of",
					"description": "This is motivation part.",
					"children":    []any{},
				},
				map[string]any{
					"name":        "Literature Review",
					"relation":    "part of",
					"description": "This is literature review part.",
					"children":    []any{},
				},
			},
		},
	})
	req, _ := http.NewRequest("POST", "/api/graphs/update", strings.NewReader(string(requestBody)))

	router.ServeHTTP(ecorder, req)

	assert.Equal(t, http.StatusOK, ecorder.Code)

	var responseBody map[string]any
	err := json.Unmarshal(ecorder.Body.Bytes(), &responseBody)
	assert.Nil(t, err)

	assert.Equal(t, map[string]any{
		"graph": map[string]any{
			"id":        "SECTION_TWO",
			"name":      "Section of Chapter One",
			"paragraph": "Updated paragraph content.",
			"children": []any{
				map[string]any{
					"name":        "Background",
					"relation":    "part of",
					"description": "This is background part.",
					"children": []any{
						map[string]any{
							"name":        "IT in Education",
							"relation":    "one of",
							"description": "This is IT in Education part.",
							"children":    []any{},
						},
						map[string]any{
							"name":        "IT in Business",
							"relation":    "one of",
							"description": "This is IT in Business part.",
							"children":    []any{},
						},
					},
				},
				map[string]any{
					"name":        "Motivation",
					"relation":    "part of",
					"description": "This is motivation part.",
					"children":    []any{},
				},
				map[string]any{
					"name":        "Literature Review",
					"relation":    "part of",
					"description": "This is literature review part.",
					"children":    []any{},
				},
			},
		},
	}, responseBody)
}

func TestGraphUpdateNotFound(t *testing.T) {

	tt := []struct {
		name      string
		userId    string
		projectId string
		chapterId string
		graphId   string
	}{
		{
			name:      "should return error when project not found",
			userId:    testutil.ModifyOnlyUserId(),
			projectId: "UNKNOWN_PROJECT",
			chapterId: "CHAPTER_ONE",
			graphId:   "SECTION_ONE",
		},
		{
			name:      "should return not found when user is not author of the project",
			userId:    testutil.ReadOnlyUserId(),
			projectId: "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_API",
			chapterId: "CHAPTER_ONE",
			graphId:   "SECTION_ONE",
		},
		{
			name:      "should return error when chapter not found",
			userId:    testutil.ModifyOnlyUserId(),
			projectId: "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_API",
			chapterId: "UNKNOWN_CHAPTER",
			graphId:   "SECTION_ONE",
		},
		{
			name:      "should return error when section not found",
			userId:    testutil.ModifyOnlyUserId(),
			projectId: "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_API",
			chapterId: "CHAPTER_ONE",
			graphId:   "UNKNOWN_SECTION",
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
				"graph": map[string]any{
					"id":        tc.graphId,
					"paragraph": "Updated paragraph content.",
					"children":  []any{},
				},
			})
			req, _ := http.NewRequest("POST", "/api/graphs/update", strings.NewReader(string(requestBody)))

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
				"graph": map[string]any{
					"children": map[string]any{},
				},
			}, responseBody)
		})
	}
}

func TestGraphUpdateDomainValidationError(t *testing.T) {
	tooLongGraphParagraph := testutil.RandomString(40001)

	tooLongChildName := testutil.RandomString(101)
	tooLongRelation := testutil.RandomString(101)
	tooLongDescription := testutil.RandomString(401)

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
				"chapter": map[string]any{
					"id": "CHAPTER_ONE",
				},
				"graph": map[string]any{
					"id":        "SECTION_ONE",
					"paragraph": "Updated paragraph content.",
					"children":  []any{},
				},
			},
			expectedResponse: map[string]any{
				"user": map[string]any{
					"id": "user id is required, but got ''",
				},
				"project": map[string]any{},
				"chapter": map[string]any{},
				"graph": map[string]any{
					"children": map[string]any{},
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
				"graph": map[string]any{
					"id":        "SECTION_ONE",
					"paragraph": "Updated paragraph content.",
					"children":  []any{},
				},
			},
			expectedResponse: map[string]any{
				"user": map[string]any{},
				"project": map[string]any{
					"id": "project id is required, but got ''",
				},
				"chapter": map[string]any{},
				"graph": map[string]any{
					"children": map[string]any{},
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
				"graph": map[string]any{
					"id":        "SECTION_ONE",
					"paragraph": "Updated paragraph content.",
					"children":  []any{},
				},
			},
			expectedResponse: map[string]any{
				"user":    map[string]any{},
				"project": map[string]any{},
				"chapter": map[string]any{
					"id": "chapter id is required, but got ''",
				},
				"graph": map[string]any{
					"children": map[string]any{},
				},
			},
		},
		{
			name: "should return error when graph id is empty",
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
				"graph": map[string]any{
					"id":        "",
					"paragraph": "Updated paragraph content.",
					"children": []any{
						map[string]any{
							"name":        "Child",
							"relation":    "child relation",
							"description": "child description",
							"children":    []any{},
						},
					},
				},
			},
			expectedResponse: map[string]any{
				"user":    map[string]any{},
				"project": map[string]any{},
				"chapter": map[string]any{},
				"graph": map[string]any{
					"id": "graph id is required, but got ''",
					"children": map[string]any{
						"items": []any{
							map[string]any{
								"children": map[string]any{},
							},
						},
					},
				},
			},
		},
		{
			name: "should return error when graph paragraph is too long",
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
				"graph": map[string]any{
					"id":        "SECTION_ONE",
					"paragraph": tooLongGraphParagraph,
					"children":  []any{},
				},
			},
			expectedResponse: map[string]any{
				"user":    map[string]any{},
				"project": map[string]any{},
				"chapter": map[string]any{},
				"graph": map[string]any{
					"paragraph": "graph paragraph must be less than or equal to 40000 bytes, but got 40001 bytes",
					"children":  map[string]any{},
				},
			},
		},
		{
			name: "should return error when child name is too long",
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
				"graph": map[string]any{
					"id":        "SECTION_ONE",
					"paragraph": "Updated paragraph content.",
					"children": []any{
						map[string]any{
							"name":        tooLongChildName,
							"relation":    "child relation",
							"description": "child description",
							"children":    []any{},
						},
					},
				},
			},
			expectedResponse: map[string]any{
				"user":    map[string]any{},
				"project": map[string]any{},
				"chapter": map[string]any{},
				"graph": map[string]any{
					"children": map[string]any{
						"items": []any{
							map[string]any{
								"name": fmt.Sprintf("graph name cannot be longer than 100 characters, but got '%v'",
									tooLongChildName),
								"children": map[string]any{},
							},
						},
					},
				},
			},
		},
		{
			name: "should return error when child relation is too long",
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
				"graph": map[string]any{
					"id":        "SECTION_ONE",
					"paragraph": "Updated paragraph content.",
					"children": []any{
						map[string]any{
							"name":        "Child",
							"relation":    tooLongRelation,
							"description": "child description",
							"children":    []any{},
						},
					},
				},
			},
			expectedResponse: map[string]any{
				"user":    map[string]any{},
				"project": map[string]any{},
				"chapter": map[string]any{},
				"graph": map[string]any{
					"children": map[string]any{
						"items": []any{
							map[string]any{
								"relation": fmt.Sprintf("graph relation cannot be longer than 100 characters, but got '%v'",
									tooLongRelation),
								"children": map[string]any{},
							},
						},
					},
				},
			},
		},
		{
			name: "should return error when child description is too long",
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
				"graph": map[string]any{
					"id":        "SECTION_ONE",
					"paragraph": "Updated paragraph content.",
					"children": []any{
						map[string]any{
							"name":        "Child",
							"relation":    "child relation",
							"description": tooLongDescription,
							"children":    []any{},
						},
					},
				},
			},
			expectedResponse: map[string]any{
				"user":    map[string]any{},
				"project": map[string]any{},
				"chapter": map[string]any{},
				"graph": map[string]any{
					"children": map[string]any{
						"items": []any{
							map[string]any{
								"description": fmt.Sprintf("graph description cannot be longer than 400 characters, but got '%v'",
									tooLongDescription),
								"children": map[string]any{},
							},
						},
					},
				},
			},
		},
		{
			name: "should return error when child names are duplicated",
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
				"graph": map[string]any{
					"id":        "SECTION_ONE",
					"paragraph": "Updated paragraph content.",
					"children": []any{
						map[string]any{
							"name":        "Child",
							"relation":    "child relation",
							"description": "child description",
							"children":    []any{},
						},
						map[string]any{
							"name":        "Child",
							"relation":    "child relation",
							"description": "child description",
							"children":    []any{},
						},
					},
				},
			},
			expectedResponse: map[string]any{
				"user":    map[string]any{},
				"project": map[string]any{},
				"chapter": map[string]any{},
				"graph": map[string]any{
					"children": map[string]any{
						"message": "names of children must be unique, but got 'Child' duplicated",
						"items": []any{
							map[string]any{
								"children": map[string]any{},
							},
							map[string]any{
								"children": map[string]any{},
							},
						},
					},
				},
			},
		},
		{
			name: "should return error when grand child has errors",
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
				"graph": map[string]any{
					"id":        "SECTION_ONE",
					"paragraph": "Updated paragraph content.",
					"children": []any{
						map[string]any{
							"name":        "Child",
							"relation":    "child relation",
							"description": "child description",
							"children": []any{
								map[string]any{
									"name":        tooLongChildName,
									"relation":    tooLongRelation,
									"description": tooLongDescription,
									"children":    []any{},
								},
							},
						},
					},
				},
			},
			expectedResponse: map[string]any{
				"user":    map[string]any{},
				"project": map[string]any{},
				"chapter": map[string]any{},
				"graph": map[string]any{
					"children": map[string]any{
						"items": []any{
							map[string]any{
								"children": map[string]any{
									"items": []any{
										map[string]any{
											"name": fmt.Sprintf("graph name cannot be longer than 100 characters, but got '%v'",
												tooLongChildName),
											"relation": fmt.Sprintf("graph relation cannot be longer than 100 characters, but got '%v'",
												tooLongRelation),
											"description": fmt.Sprintf("graph description cannot be longer than 400 characters, but got '%v'",
												tooLongDescription),
											"children": map[string]any{},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			router := setupGraphRouter()

			ecorder := httptest.NewRecorder()
			requestBody, _ := json.Marshal(tc.request)
			req, _ := http.NewRequest("POST", "/api/graphs/update", strings.NewReader(string(requestBody)))

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
				"graph":   tc.expectedResponse["graph"],
			}, responseBody)
		})
	}
}

func TestGraphUpdateInvalidRequestFormat(t *testing.T) {
	router := setupGraphRouter()

	ecorder := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/graphs/update", strings.NewReader(""))

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
		"graph": map[string]any{
			"children": map[string]any{},
		},
	}, responseBody)
}

func TestGraphDelete(t *testing.T) {
	router := setupGraphRouter()

	ecorder := httptest.NewRecorder()
	requestBody, _ := json.Marshal(map[string]any{
		"user": map[string]any{
			"id": testutil.ModifyOnlyUserId(),
		},
		"project": map[string]any{
			"id": "PROJECT_WITHOUT_DESCRIPTION_TO_DELETE_FROM_API",
		},
		"chapter": map[string]any{
			"id": "CHAPTER_ONE",
		},
		"section": map[string]any{
			"id": "SECTION_ONE",
		},
	})
	req, _ := http.NewRequest("POST", "/api/graphs/delete", strings.NewReader(string(requestBody)))

	router.ServeHTTP(ecorder, req)

	assert.Equal(t, http.StatusOK, ecorder.Code)

	var responseBody map[string]any
	err := json.Unmarshal(ecorder.Body.Bytes(), &responseBody)
	assert.Nil(t, err)

	assert.Equal(t, map[string]any{
		"message": "graph successfully deleted",
	}, responseBody)
}

func TestGraphDeleteNotFound(t *testing.T) {
	tt := []struct {
		name      string
		userId    string
		projectId string
		chapterId string
		graphId   string
	}{
		{
			name:      "should return error when project not found",
			userId:    testutil.ModifyOnlyUserId(),
			projectId: "UNKNOWN_PROJECT",
			chapterId: "CHAPTER_ONE",
			graphId:   "SECTION_ONE",
		},
		{
			name:      "should return not found when user is not author of the project",
			userId:    testutil.ReadOnlyUserId(),
			projectId: "PROJECT_WITHOUT_DESCRIPTION_TO_DELETE_FROM_API",
			chapterId: "CHAPTER_ONE",
			graphId:   "SECTION_ONE",
		},
		{
			name:      "should return error when chapter not found",
			userId:    testutil.ModifyOnlyUserId(),
			projectId: "PROJECT_WITHOUT_DESCRIPTION_TO_DELETE_FROM_API",
			chapterId: "UNKNOWN_CHAPTER",
			graphId:   "SECTION_ONE",
		},
		{
			name:      "should return error when section not found",
			userId:    testutil.ModifyOnlyUserId(),
			projectId: "PROJECT_WITHOUT_DESCRIPTION_TO_DELETE_FROM_API",
			chapterId: "CHAPTER_ONE",
			graphId:   "UNKNOWN_SECTION",
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
				"section": map[string]any{
					"id": tc.graphId,
				},
			})
			req, _ := http.NewRequest("POST", "/api/graphs/delete", strings.NewReader(string(requestBody)))

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

func TestGraphDeleteDomainValidationError(t *testing.T) {
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
			req, _ := http.NewRequest("POST", "/api/graphs/delete", strings.NewReader(string(requestBody)))

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

func TestGraphDeleteInvalidRequestFormat(t *testing.T) {
	router := setupGraphRouter()

	ecorder := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/graphs/delete", strings.NewReader(""))

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
				"children":  []any{},
			},
			map[string]any{
				"id":        graphId2,
				"name":      maxLengthSectionName,
				"paragraph": maxLengthSectionContent,
				"children":  []any{},
			},
		},
	}, responseBody)
}

func TestGraphSectionalizeNotFound(t *testing.T) {
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
	router.POST("/api/graphs/delete", api.HandleDelete)
	router.POST("/api/graphs/sectionalize", api.HandleSectionalize)

	return router
}
