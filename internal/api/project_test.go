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
				"user":    map[string]any{"id": testutil.UserId()},
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
			assert.Equal(t, responseBody, map[string]any{
				"project": projectWithId,
			})
		})
	}
}

func TestProjectCreateInvalidArgument(t *testing.T) {
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
				"user": map[string]any{"id": testutil.UserId()},
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
						"project name cannot be longer than 100 characters, but got '%s'",
						tooLongProjectName,
					),
					"description": fmt.Sprintf(
						"project description cannot be longer than 400 characters, but got '%s'",
						tooLongProjectDescription,
					),
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
	tt := []struct {
		name    string
		request string
	}{
		{
			name:    "should return error when user id is not string",
			request: `{ "user": { "id": 123 }, "project": { "name": "New Project" } }`,
		},
		{
			name:    "should return error when project name is not string",
			request: `{ "user": { "id": "user-id" }, "project": { "name": 123 } }`,
		},
		{
			name:    "should return error when user is not object",
			request: `{ "user": 123, "project": { "name": "New Project" } }`,
		},
		{
			name:    "should return error when project is not object",
			request: `{ "user": { "id": "user-id" }, "project": "New Project" }`,
		},
		{
			name:    "should return error when request body is invalid JSON",
			request: `{ "user": { "id": "user-id", "project": { "name": "New Project" }`,
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
			req, _ := http.NewRequest("POST", "/api/projects/create", strings.NewReader(tc.request))

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

func setupProjectRouter() *gin.Engine {
	router := gin.Default()

	client := db.FirestoreClient()
	r := repository.NewProjectRepository(*client)
	s := service.NewProjectService(r)
	uc := usecase.NewProjectUseCase(s)
	a := api.NewProjectApi(uc)

	router.POST("/api/projects/list", a.HandleList)
	router.POST("/api/projects/create", a.HandleCreate)
	return router
}
