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

func TestProjectList(t *testing.T) {
	router := setupProjectRouter()

	recorder := httptest.NewRecorder()
	requestBody, _ := json.Marshal(map[string]any{
		"user": map[string]any{
			"id": testutil.ReadOnlyUserId(),
		},
	})
	req, _ := http.NewRequest("POST", "/api/projects/list", strings.NewReader(string(requestBody)))

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var responseBody map[string]any
	assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &responseBody))
	assert.Equal(t, map[string]any{
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
	}, responseBody)
}

func TestProjectListDomainValidationError(t *testing.T) {
	tt := []struct {
		name             string
		request          map[string]any
		expectedResponse map[string]any
	}{
		{
			name: "should return error when user id is empty",
			request: map[string]any{
				"user": map[string]any{"id": ""},
			},
			expectedResponse: map[string]any{
				"user": map[string]any{
					"id": "user id is required, but got ''",
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
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			router := setupProjectRouter()

			recorder := httptest.NewRecorder()
			requestBody, _ := json.Marshal(tc.request)
			req, _ := http.NewRequest("POST", "/api/projects/list", strings.NewReader(string(requestBody)))

			router.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusBadRequest, recorder.Code)

			var responseBody map[string]any
			assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &responseBody))
			assert.Equal(t, map[string]any{
				"message": "invalid request value",
				"user":    tc.expectedResponse["user"],
			}, responseBody)
		})
	}
}

func TestProjectListInvalidRequestFormat(t *testing.T) {
	router := setupProjectRouter()

	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/projects/list", strings.NewReader(""))

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	var responseBody map[string]any
	assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &responseBody))
	assert.Equal(t, map[string]any{
		"message": "invalid request format",
		"user":    map[string]any{},
	}, responseBody)
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

func TestProjectFind(t *testing.T) {
	tt := []struct {
		name             string
		request          map[string]any
		expectedResponse map[string]any
	}{
		{
			name: "should find project without description",
			request: map[string]any{
				"user":    map[string]any{"id": testutil.ReadOnlyUserId()},
				"project": map[string]any{"id": "PROJECT_WITHOUT_DESCRIPTION"},
			},
			expectedResponse: map[string]any{
				"project": map[string]any{
					"id":   "PROJECT_WITHOUT_DESCRIPTION",
					"name": "No Description Project",
				},
			},
		},
		{
			name: "should find project with description",
			request: map[string]any{
				"user":    map[string]any{"id": testutil.ReadOnlyUserId()},
				"project": map[string]any{"id": "PROJECT_WITH_DESCRIPTION"},
			},
			expectedResponse: map[string]any{
				"project": map[string]any{
					"id":          "PROJECT_WITH_DESCRIPTION",
					"name":        "Described Project",
					"description": "This is project description",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			router := setupProjectRouter()

			recorder := httptest.NewRecorder()
			requestBody, _ := json.Marshal(tc.request)
			req, _ := http.NewRequest("POST", "/api/projects/find", strings.NewReader(string(requestBody)))

			router.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusOK, recorder.Code)

			var responseBody map[string]any
			assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &responseBody))

			assert.Equal(t, tc.expectedResponse, responseBody)
		})
	}
}

func TestProjectFindNotFound(t *testing.T) {
	tt := []struct {
		name    string
		request map[string]any
	}{
		{
			name: "should return not found when project is not found",
			request: map[string]any{
				"user":    map[string]any{"id": testutil.ReadOnlyUserId()},
				"project": map[string]any{"id": "NOT_FOUND_PROJECT"},
			},
		},
		{
			name: "should return not found when user is not author of the project",
			request: map[string]any{
				"user":    map[string]any{"id": testutil.ModifyOnlyUserId()},
				"project": map[string]any{"id": "PROJECT_WITHOUT_DESCRIPTION"},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			router := setupProjectRouter()

			recorder := httptest.NewRecorder()
			requestBody, _ := json.Marshal(tc.request)
			req, _ := http.NewRequest("POST", "/api/projects/find", strings.NewReader(string(requestBody)))

			router.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusNotFound, recorder.Code)

			var responseBody map[string]any
			assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &responseBody))
			assert.Equal(t, map[string]any{
				"message": "not found",
				"user":    map[string]any{},
				"project": map[string]any{},
			}, responseBody)
		})
	}
}

func TestProjectFindDomainValidationError(t *testing.T) {
	tt := []struct {
		name             string
		request          map[string]any
		expectedResponse map[string]any
	}{
		{
			name: "should return error when user id is empty",
			request: map[string]any{
				"user":    map[string]any{"id": ""},
				"project": map[string]any{"id": "PROJECT_WITHOUT_DESCRIPTION"},
			},
			expectedResponse: map[string]any{
				"user": map[string]any{
					"id": "user id is required, but got ''",
				},
				"project": map[string]any{},
			},
		},
		{
			name: "should return error when project id is empty",
			request: map[string]any{
				"user":    map[string]any{"id": testutil.ReadOnlyUserId()},
				"project": map[string]any{"id": ""},
			},
			expectedResponse: map[string]any{
				"user": map[string]any{},
				"project": map[string]any{
					"id": "project id is required, but got ''",
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
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			router := setupProjectRouter()

			recorder := httptest.NewRecorder()
			requestBody, _ := json.Marshal(tc.request)
			req, _ := http.NewRequest("POST", "/api/projects/find", strings.NewReader(string(requestBody)))

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

func TestProjectFindInvalidRequestFormat(t *testing.T) {
	router := setupProjectRouter()

	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/projects/find", strings.NewReader(""))

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	var responseBody map[string]any
	assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &responseBody))
	assert.Equal(t, map[string]any{
		"message": "invalid request format",
		"user":    map[string]any{},
		"project": map[string]any{},
	}, responseBody)
}

func TestProjectFindInternalError(t *testing.T) {
	router := setupProjectRouter()

	recorder := httptest.NewRecorder()
	requestBody, _ := json.Marshal(map[string]any{
		"user": map[string]any{
			"id": testutil.ErrorUserId(0),
		},
		"project": map[string]any{
			"id": "PROJECT_WITH_INVALID_NAME",
		},
	})
	req, _ := http.NewRequest("POST", "/api/projects/find", strings.NewReader(string(requestBody)))

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	var responseBody map[string]any
	assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &responseBody))
	assert.Equal(t, map[string]any{
		"message": "internal error",
	}, responseBody)
}

func TestProjectCreate(t *testing.T) {
	maxLengthProjectName := testutil.RandomString(100)
	maxLengthProjectDescription := testutil.RandomString(400)

	tt := []struct {
		name    string
		project map[string]any
	}{
		{
			name: "should create project without description",
			project: map[string]any{
				"name": "New Project",
			},
		},
		{
			name: "should create project with description",
			project: map[string]any{
				"name":        "New Project",
				"description": "This is project description",
			},
		},
		{
			name: "should create project with maximum length properties",
			project: map[string]any{
				"name":        maxLengthProjectName,
				"description": maxLengthProjectDescription,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			router := setupProjectRouter()

			recorder := httptest.NewRecorder()
			requestBody, _ := json.Marshal(map[string]any{
				"user":    map[string]any{"id": testutil.ModifyOnlyUserId()},
				"project": tc.project,
			})
			req, _ := http.NewRequest("POST", "/api/projects/create", strings.NewReader(string(requestBody)))

			router.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusCreated, recorder.Code)

			var responseBody map[string]any
			assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &responseBody))

			//projectId is generated by firestore and it's not predictable
			projectId := responseBody["project"].(map[string]any)["id"]
			assert.NotEmpty(t, projectId)

			projectWithId := tc.project
			projectWithId["id"] = projectId
			assert.Equal(t, map[string]any{
				"project": projectWithId,
			}, responseBody)
		})
	}
}

func TestProjectCreateDomainValidationError(t *testing.T) {
	tooLongProjectName := testutil.RandomString(101)
	tooLongProjectDescription := testutil.RandomString(401)

	tt := []struct {
		name             string
		request          map[string]any
		expectedResponse map[string]any
	}{
		{
			name: "should return error when user id is empty",
			request: map[string]any{
				"user": map[string]any{"id": ""},
				"project": map[string]any{
					"name": "New Project",
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
			name: "should return error when project name is empty",
			request: map[string]any{
				"user": map[string]any{"id": testutil.ModifyOnlyUserId()},
				"project": map[string]any{
					"name": "",
				},
			},
			expectedResponse: map[string]any{
				"user": map[string]any{},
				"project": map[string]any{
					"name": "project name is required, but got ''",
				},
			},
		},
		{
			name: "should return error when project properties are too long",
			request: map[string]any{
				"user": map[string]any{"id": ""},
				"project": map[string]any{
					"name":        tooLongProjectName,
					"description": tooLongProjectDescription,
				},
			},
			expectedResponse: map[string]any{
				"user": map[string]any{
					"id": "user id is required, but got ''",
				},
				"project": map[string]any{
					"name": fmt.Sprintf(
						"project name cannot be longer than 100 characters, but got '%v'",
						tooLongProjectName,
					),
					"description": fmt.Sprintf(
						"project description cannot be longer than 400 characters, but got '%v'",
						tooLongProjectDescription,
					),
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
					"name": "project name is required, but got ''",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			router := setupProjectRouter()

			recorder := httptest.NewRecorder()
			requestBody, _ := json.Marshal(tc.request)
			req, _ := http.NewRequest("POST", "/api/projects/create", strings.NewReader(string(requestBody)))

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

func TestProjectCreateInvalidRequestFormat(t *testing.T) {
	router := setupProjectRouter()

	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/projects/create", strings.NewReader(""))

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	var responseBody map[string]any
	assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &responseBody))
	assert.Equal(t, map[string]any{
		"message": "invalid request format",
		"user":    map[string]any{},
		"project": map[string]any{},
	}, responseBody)
}

func TestProjectUpdate(t *testing.T) {
	tt := []struct {
		name    string
		project map[string]any
	}{
		{
			name: "should update project name",
			project: map[string]any{
				"id":   "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_API",
				"name": "Updated Project",
			},
		},
		{
			name: "should update project name and description",
			project: map[string]any{
				"id":          "PROJECT_WITH_DESCRIPTION_TO_UPDATE_FROM_API",
				"name":        "Updated Project",
				"description": "Updated project description",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			router := setupProjectRouter()

			recorder := httptest.NewRecorder()
			requestBody, _ := json.Marshal(map[string]any{
				"user":    map[string]any{"id": testutil.ModifyOnlyUserId()},
				"project": tc.project,
			})
			req, _ := http.NewRequest("POST", "/api/projects/update", strings.NewReader(string(requestBody)))

			router.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusOK, recorder.Code)

			var responseBody map[string]any
			assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &responseBody))

			assert.Equal(t, map[string]any{
				"project": tc.project,
			}, responseBody)
		})
	}
}

func TestProjectUpdateNotFound(t *testing.T) {
	tt := []struct {
		name    string
		project map[string]any
	}{
		{
			name: "should return not found when project is not found",
			project: map[string]any{
				"id":   "NOT_FOUND_PROJECT",
				"name": "Updated Project",
			},
		},
		{
			name: "should return not found when user is not author of the project",
			project: map[string]any{
				"id":   "PROJECT_TO_UPDATE_WITHOUT_DESCRIPTION",
				"name": "Updated Project",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			router := setupProjectRouter()

			recorder := httptest.NewRecorder()
			requestBody, _ := json.Marshal(map[string]any{
				"user":    map[string]any{"id": testutil.ReadOnlyUserId()},
				"project": tc.project,
			})
			req, _ := http.NewRequest("POST", "/api/projects/update", strings.NewReader(string(requestBody)))

			router.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusNotFound, recorder.Code)

			var responseBody map[string]any
			assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &responseBody))
			assert.Equal(t, map[string]any{
				"message": "not found",
				"user":    map[string]any{},
				"project": map[string]any{},
			}, responseBody)
		})
	}
}

func TestProjectUpdateDomainValidationError(t *testing.T) {
	tooLongProjectName := testutil.RandomString(101)
	tooLongProjectDescription := testutil.RandomString(401)

	tt := []struct {
		name             string
		request          map[string]any
		expectedResponse map[string]any
	}{
		{
			name: "should return error when user id is empty",
			request: map[string]any{
				"user": map[string]any{"id": ""},
				"project": map[string]any{
					"id":   "PROJECT_TO_UPDATE_WITHOUT_DESCRIPTION",
					"name": "Updated Project",
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
			name: "should return error when project id is empty",
			request: map[string]any{
				"user": map[string]any{"id": testutil.ModifyOnlyUserId()},
				"project": map[string]any{
					"id":   "",
					"name": "Updated Project",
				},
			},
			expectedResponse: map[string]any{
				"user": map[string]any{},
				"project": map[string]any{
					"id": "project id is required, but got ''",
				},
			},
		},
		{
			name: "should return error when project name is empty",
			request: map[string]any{
				"user": map[string]any{"id": testutil.ModifyOnlyUserId()},
				"project": map[string]any{
					"id":   "PROJECT_TO_UPDATE_WITHOUT_DESCRIPTION",
					"name": "",
				},
			},
			expectedResponse: map[string]any{
				"user": map[string]any{},
				"project": map[string]any{
					"name": "project name is required, but got ''",
				},
			},
		},
		{
			name: "should return error when project properties are too long",
			request: map[string]any{
				"user": map[string]any{"id": ""},
				"project": map[string]any{
					"id":          "PROJECT_TO_UPDATE_WITHOUT_DESCRIPTION",
					"name":        tooLongProjectName,
					"description": tooLongProjectDescription,
				},
			},
			expectedResponse: map[string]any{
				"user": map[string]any{
					"id": "user id is required, but got ''",
				},
				"project": map[string]any{
					"name": fmt.Sprintf(
						"project name cannot be longer than 100 characters, but got '%v'",
						tooLongProjectName,
					),
					"description": fmt.Sprintf(
						"project description cannot be longer than 400 characters, but got '%v'",
						tooLongProjectDescription,
					),
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
					"id":   "project id is required, but got ''",
					"name": "project name is required, but got ''",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			router := setupProjectRouter()

			recorder := httptest.NewRecorder()
			requestBody, _ := json.Marshal(tc.request)
			req, _ := http.NewRequest("POST", "/api/projects/update", strings.NewReader(string(requestBody)))

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

func TestProjectUpdateInvalidRequestFormat(t *testing.T) {
	router := setupProjectRouter()

	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/projects/update", strings.NewReader(""))

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	var responseBody map[string]any
	assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &responseBody))
	assert.Equal(t, map[string]any{
		"message": "invalid request format",
		"user":    map[string]any{},
		"project": map[string]any{},
	}, responseBody)
}

func TestProjectUpdateInternalError(t *testing.T) {
	router := setupProjectRouter()

	recorder := httptest.NewRecorder()
	requestBody, _ := json.Marshal(map[string]any{
		"user": map[string]any{
			"id": testutil.ErrorUserId(0),
		},
		"project": map[string]any{
			"id":   "PROJECT_WITH_INVALID_NAME",
			"name": "Updated Project",
		},
	})
	req, _ := http.NewRequest("POST", "/api/projects/update", strings.NewReader(string(requestBody)))

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	var responseBody map[string]any
	assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &responseBody))
	assert.Equal(t, map[string]any{
		"message": "internal error",
	}, responseBody)
}

func TestProjectDelete(t *testing.T) {
	router := setupProjectRouter()

	recorder := httptest.NewRecorder()
	requestBody, _ := json.Marshal(map[string]any{
		"user": map[string]any{
			"id": testutil.ModifyOnlyUserId(),
		},
		"project": map[string]any{
			"id": "PROJECT_WITH_DESCRIPTION_TO_DELETE_FROM_API",
		},
	})
	req, _ := http.NewRequest("POST", "/api/projects/delete", strings.NewReader(string(requestBody)))

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusNoContent, recorder.Code)
}

func TestProjectDeleteNotFound(t *testing.T) {
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
			},
		},
		{
			name: "should return not found when user is not author of the project",
			request: map[string]any{
				"user": map[string]any{
					"id": testutil.ReadOnlyUserId(),
				},
				"project": map[string]any{
					"id": "PROJECT_WITH_DESCRIPTION_TO_DELETE_FROM_API",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			router := setupProjectRouter()

			recorder := httptest.NewRecorder()
			requestBody, _ := json.Marshal(tc.request)
			req, _ := http.NewRequest("POST", "/api/projects/delete", strings.NewReader(string(requestBody)))

			router.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusNotFound, recorder.Code)

			var responseBody map[string]any
			assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &responseBody))
			assert.Equal(t, map[string]any{
				"message": "not found",
				"user":    map[string]any{},
				"project": map[string]any{},
			}, responseBody)
		})
	}
}

func TestProjectDeleteDomainValidationError(t *testing.T) {
	tt := []struct {
		name             string
		request          map[string]any
		expectedResponse map[string]any
	}{
		{
			name: "should return error when user id is empty",
			request: map[string]any{
				"user":    map[string]any{"id": ""},
				"project": map[string]any{"id": "PROJECT_WITHOUT_DESCRIPTION"},
			},
			expectedResponse: map[string]any{
				"user": map[string]any{
					"id": "user id is required, but got ''",
				},
				"project": map[string]any{},
			},
		},
		{
			name: "should return error when project id is empty",
			request: map[string]any{
				"user":    map[string]any{"id": testutil.ReadOnlyUserId()},
				"project": map[string]any{"id": ""},
			},
			expectedResponse: map[string]any{
				"user": map[string]any{},
				"project": map[string]any{
					"id": "project id is required, but got ''",
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
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			router := setupProjectRouter()

			recorder := httptest.NewRecorder()
			requestBody, _ := json.Marshal(tc.request)
			req, _ := http.NewRequest("POST", "/api/projects/delete", strings.NewReader(string(requestBody)))

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

func TestProjectDeleteInvalidRequestFormat(t *testing.T) {
	router := setupProjectRouter()

	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/projects/delete", strings.NewReader(""))

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestProjectDeleteInternalError(t *testing.T) {
	router := setupProjectRouter()

	recorder := httptest.NewRecorder()
	requestBody, _ := json.Marshal(map[string]any{
		"user": map[string]any{
			"id": testutil.ErrorUserId(0),
		},
		"project": map[string]any{
			"id": "PROJECT_WITH_INVALID_NAME",
		},
	})
	req, _ := http.NewRequest("POST", "/api/projects/delete", strings.NewReader(string(requestBody)))

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
	router.POST("/api/projects/create", a.HandleCreate)
	router.POST("/api/projects/find", a.HandleFind)
	router.POST("/api/projects/update", a.HandleUpdate)
	router.POST("/api/projects/delete", a.HandleDelete)
	return router
}
