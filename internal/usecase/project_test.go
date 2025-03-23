package usecase_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/kumachan-mis/knodeledge-api/internal/domain"
	"github.com/kumachan-mis/knodeledge-api/internal/openapi"
	"github.com/kumachan-mis/knodeledge-api/internal/service"
	"github.com/kumachan-mis/knodeledge-api/internal/testutil"
	"github.com/kumachan-mis/knodeledge-api/internal/usecase"
	mock_service "github.com/kumachan-mis/knodeledge-api/mock/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestListProjectsValidEntity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mock_service.NewMockProjectService(ctrl)

	id, err := domain.NewProjectIdObject("0000000000000001")
	assert.NoError(t, err)
	name, err := domain.NewProjectNameObject("Project With Description")
	assert.NoError(t, err)
	description, err := domain.NewProjectDescriptionObject("This is a project")
	assert.NoError(t, err)
	createdAt, err := domain.NewCreatedAtObject(testutil.Date())
	assert.NoError(t, err)
	updatedAt, err := domain.NewUpdatedAtObject(testutil.Date())
	assert.NoError(t, err)

	projectWithDesc := domain.NewProjectEntity(*id, *name, *description, *createdAt, *updatedAt)

	id, err = domain.NewProjectIdObject("0000000000000002")
	assert.NoError(t, err)
	name, err = domain.NewProjectNameObject("Project Without Description")
	assert.NoError(t, err)
	description, err = domain.NewProjectDescriptionObject("")
	assert.NoError(t, err)
	createdAt, err = domain.NewCreatedAtObject(testutil.Date().Add(-1 * time.Hour))
	assert.NoError(t, err)
	updatedAt, err = domain.NewUpdatedAtObject(testutil.Date().Add(-1 * time.Hour))
	assert.NoError(t, err)

	projectWithoutDesc := domain.NewProjectEntity(*id, *name, *description, *createdAt, *updatedAt)

	s.EXPECT().
		ListProjects(gomock.Any()).
		Do(func(userId domain.UserIdObject) {
			assert.Equal(t, testutil.ReadOnlyUserId(), userId.Value())
		}).
		Return([]domain.ProjectEntity{*projectWithDesc, *projectWithoutDesc}, nil)

	uc := usecase.NewProjectUseCase(s)

	res, ucErr := uc.ListProjects(openapi.ProjectListRequest{
		UserId: testutil.ReadOnlyUserId(),
	})
	assert.Nil(t, ucErr)

	assert.Len(t, res.Projects, 2)

	project := res.Projects[0]
	assert.Equal(t, "0000000000000001", project.Id)
	assert.Equal(t, "Project With Description", project.Name)
	assert.Equal(t, "This is a project", project.Description)

	project = res.Projects[1]
	assert.Equal(t, "0000000000000002", project.Id)
	assert.Equal(t, "Project Without Description", project.Name)
	assert.Equal(t, "", project.Description)
}

func TestListProjectsDomainValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tt := []struct {
		name     string
		userId   string
		expected openapi.ProjectListErrorResponse
	}{
		{
			name:   "should return error when user id is empty",
			userId: "",
			expected: openapi.ProjectListErrorResponse{
				UserId: "user id is required, but got ''",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			s := mock_service.NewMockProjectService(ctrl)

			uc := usecase.NewProjectUseCase(s)

			res, ucErr := uc.ListProjects(openapi.ProjectListRequest{
				UserId: tc.userId,
			})
			assert.NotNil(t, ucErr)

			expectedJson, _ := json.Marshal(tc.expected)
			assert.Equal(t, fmt.Sprintf("domain validation error: %s", expectedJson), ucErr.Error())
			assert.Equal(t, usecase.DomainValidationError, ucErr.Code())
			assert.Equal(t, tc.expected, *ucErr.Response())

			assert.Nil(t, res)
		})
	}
}

func TestListProjectsServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mock_service.NewMockProjectService(ctrl)

	s.EXPECT().
		ListProjects(gomock.Any()).
		Return(nil, service.Errorf(service.RepositoryFailurePanic, "service error"))

	uc := usecase.NewProjectUseCase(s)

	res, ucErr := uc.ListProjects(openapi.ProjectListRequest{
		UserId: testutil.ReadOnlyUserId(),
	})
	assert.NotNil(t, ucErr)
	assert.Equal(t, "internal error: service error", ucErr.Error())
	assert.Equal(t, usecase.InternalErrorPanic, ucErr.Code())
	assert.Nil(t, ucErr.Response())
	assert.Nil(t, res)
}

func TestFindProjectValidEntity(t *testing.T) {
	tt := []struct {
		name               string
		projectId          string
		projectName        string
		projectDescription string
	}{
		{
			name:               "should find project with description",
			projectId:          "0000000000000001",
			projectName:        "Project With Description",
			projectDescription: "This is a project",
		},
		{
			name:               "should find project without description",
			projectId:          "0000000000000002",
			projectName:        "Project Without Description",
			projectDescription: "",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mock_service.NewMockProjectService(ctrl)

			id, err := domain.NewProjectIdObject(tc.projectId)
			assert.NoError(t, err)
			name, err := domain.NewProjectNameObject(tc.projectName)
			assert.NoError(t, err)
			description, err := domain.NewProjectDescriptionObject(tc.projectDescription)
			assert.NoError(t, err)
			createdAt, err := domain.NewCreatedAtObject(testutil.Date())
			assert.NoError(t, err)
			updatedAt, err := domain.NewUpdatedAtObject(testutil.Date())
			assert.NoError(t, err)

			project := domain.NewProjectEntity(*id, *name, *description, *createdAt, *updatedAt)

			s.EXPECT().
				FindProject(gomock.Any(), gomock.Any()).
				Do(func(userId domain.UserIdObject, projectId domain.ProjectIdObject) {
					assert.Equal(t, testutil.ReadOnlyUserId(), userId.Value())
					assert.Equal(t, tc.projectId, projectId.Value())
				}).
				Return(project, nil)

			uc := usecase.NewProjectUseCase(s)

			res, ucErr := uc.FindProject(openapi.ProjectFindRequest{
				UserId:    testutil.ReadOnlyUserId(),
				ProjectId: tc.projectId,
			})
			assert.Nil(t, ucErr)

			assert.Equal(t, tc.projectId, res.Project.Id)
			assert.Equal(t, tc.projectName, res.Project.Name)
			assert.Equal(t, tc.projectDescription, res.Project.Description)
		})
	}
}

func TestFindProjectDomainValidationError(t *testing.T) {
	tt := []struct {
		name      string
		userId    string
		projectId string
		expected  openapi.ProjectFindErrorResponse
	}{
		{
			name:      "should return error when user id is empty",
			userId:    "",
			projectId: "0000000000000001",
			expected: openapi.ProjectFindErrorResponse{
				UserId:    "user id is required, but got ''",
				ProjectId: "",
			},
		},
		{
			name:      "should return error when project id is empty",
			userId:    testutil.ReadOnlyUserId(),
			projectId: "",
			expected: openapi.ProjectFindErrorResponse{
				UserId:    "",
				ProjectId: "project id is required, but got ''",
			},
		},
		{
			name:      "should return error when all fields are empty",
			userId:    "",
			projectId: "",
			expected: openapi.ProjectFindErrorResponse{
				UserId:    "user id is required, but got ''",
				ProjectId: "project id is required, but got ''",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mock_service.NewMockProjectService(ctrl)

			uc := usecase.NewProjectUseCase(s)

			res, ucErr := uc.FindProject(openapi.ProjectFindRequest{
				UserId:    tc.userId,
				ProjectId: tc.projectId,
			})
			assert.NotNil(t, ucErr)

			expectedJson, _ := json.Marshal(tc.expected)
			assert.Equal(t, fmt.Sprintf("domain validation error: %s", expectedJson), ucErr.Error())
			assert.Equal(t, usecase.DomainValidationError, ucErr.Code())
			assert.Equal(t, tc.expected, *ucErr.Response())

			assert.Nil(t, res)
		})
	}
}

func TestFindProjectServiceError(t *testing.T) {
	tt := []struct {
		name          string
		errorCode     service.ErrorCode
		errorMessage  string
		expectedError string
		expectedCode  usecase.ErrorCode
	}{
		{
			name:          "should return error when project not found",
			errorCode:     service.NotFoundError,
			errorMessage:  "failed to find project",
			expectedError: "not found: failed to find project",
			expectedCode:  usecase.NotFoundError,
		},
		{
			name:          "should return error when repository failure",
			errorCode:     service.RepositoryFailurePanic,
			errorMessage:  "service error",
			expectedError: "internal error: service error",
			expectedCode:  usecase.InternalErrorPanic,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mock_service.NewMockProjectService(ctrl)

			s.EXPECT().
				FindProject(gomock.Any(), gomock.Any()).
				Return(nil, service.Errorf(tc.errorCode, "%s", tc.errorMessage))

			uc := usecase.NewProjectUseCase(s)

			res, ucErr := uc.FindProject(openapi.ProjectFindRequest{
				UserId:    testutil.ReadOnlyUserId(),
				ProjectId: "0000000000000001",
			})
			assert.NotNil(t, ucErr)
			assert.Equal(t, tc.expectedError, ucErr.Error())
			assert.Equal(t, tc.expectedCode, ucErr.Code())
			assert.Nil(t, ucErr.Response())
			assert.Nil(t, res)
		})
	}
}

func TestCreateProjectValidEntity(t *testing.T) {
	maxLengthProjectName := testutil.RandomString(100)
	maxLengthProjectDescription := testutil.RandomString(400)

	tt := []struct {
		name    string
		project openapi.ProjectWithoutAutofield
	}{
		{
			name: "should create project with description",
			project: openapi.ProjectWithoutAutofield{
				Name:        "Project With Description",
				Description: "This is a project",
			},
		},
		{
			name: "should create project without description",
			project: openapi.ProjectWithoutAutofield{
				Name: "Project Without Description",
			},
		},
		{
			name: "should create project with max length name",
			project: openapi.ProjectWithoutAutofield{
				Name: maxLengthProjectName,
			},
		},
		{
			name: "should create project with max length description",
			project: openapi.ProjectWithoutAutofield{
				Name:        "Project With Max Length Description",
				Description: maxLengthProjectDescription,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mock_service.NewMockProjectService(ctrl)

			id, err := domain.NewProjectIdObject("0000000000000001")
			assert.NoError(t, err)
			name, err := domain.NewProjectNameObject(tc.project.Name)
			assert.NoError(t, err)
			description, err := domain.NewProjectDescriptionObject(tc.project.Description)
			assert.NoError(t, err)
			createdAt, err := domain.NewCreatedAtObject(testutil.Date())
			assert.NoError(t, err)
			updatedAt, err := domain.NewUpdatedAtObject(testutil.Date())
			assert.NoError(t, err)

			project := domain.NewProjectEntity(*id, *name, *description, *createdAt, *updatedAt)

			s.EXPECT().
				CreateProject(gomock.Any(), gomock.Any()).
				Do(func(userId domain.UserIdObject, project domain.ProjectWithoutAutofieldEntity) {
					assert.Equal(t, testutil.ModifyOnlyUserId(), userId.Value())
					assert.Equal(t, tc.project.Name, project.Name().Value())
					assert.Equal(t, tc.project.Description, project.Description().Value())
				}).
				Return(project, nil)

			uc := usecase.NewProjectUseCase(s)

			res, ucErr := uc.CreateProject(openapi.ProjectCreateRequest{
				User:    openapi.UserOnlyId{Id: testutil.ModifyOnlyUserId()},
				Project: tc.project,
			})
			assert.Nil(t, ucErr)

			assert.Equal(t, "0000000000000001", res.Project.Id)
			assert.Equal(t, tc.project.Name, res.Project.Name)
			assert.Equal(t, tc.project.Description, res.Project.Description)

		})
	}
}

func TestCreateProjectDomainValidationError(t *testing.T) {
	tooLongProjectName := testutil.RandomString(101)
	tooLongProjectDescription := testutil.RandomString(401)

	tt := []struct {
		name     string
		userId   string
		project  openapi.ProjectWithoutAutofield
		expected openapi.ProjectCreateErrorResponse
	}{
		{
			name:   "should return error when user id is empty",
			userId: "",
			project: openapi.ProjectWithoutAutofield{
				Name:        "Project With Description",
				Description: "This is a project",
			},
			expected: openapi.ProjectCreateErrorResponse{
				User: openapi.UserOnlyIdError{
					Id: "user id is required, but got ''",
				},
			},
		},
		{
			name:   "should return error when project name is empty",
			userId: testutil.ModifyOnlyUserId(),
			project: openapi.ProjectWithoutAutofield{
				Name:        "",
				Description: "This is a project",
			},
			expected: openapi.ProjectCreateErrorResponse{
				Project: openapi.ProjectWithoutAutofieldError{
					Name: "project name is required, but got ''",
				},
			},
		},
		{
			name:   "should return error when project name is too long",
			userId: testutil.ModifyOnlyUserId(),
			project: openapi.ProjectWithoutAutofield{
				Name: tooLongProjectName,
			},
			expected: openapi.ProjectCreateErrorResponse{
				Project: openapi.ProjectWithoutAutofieldError{
					Name: fmt.Sprintf(
						"project name cannot be longer than 100 characters, but got '%v'",
						tooLongProjectName,
					),
				},
			},
		},
		{
			name:   "should return error when project description is too long",
			userId: testutil.ModifyOnlyUserId(),
			project: openapi.ProjectWithoutAutofield{
				Name:        "Project With Description",
				Description: tooLongProjectDescription,
			},
			expected: openapi.ProjectCreateErrorResponse{
				Project: openapi.ProjectWithoutAutofieldError{
					Description: fmt.Sprintf(
						"project description cannot be longer than 400 characters, but got '%v'",
						tooLongProjectDescription,
					),
				},
			},
		},
		{
			name:    "should return error when all fields are empty",
			userId:  "",
			project: openapi.ProjectWithoutAutofield{},
			expected: openapi.ProjectCreateErrorResponse{
				User: openapi.UserOnlyIdError{
					Id: "user id is required, but got ''",
				},
				Project: openapi.ProjectWithoutAutofieldError{
					Name: "project name is required, but got ''",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mock_service.NewMockProjectService(ctrl)

			uc := usecase.NewProjectUseCase(s)

			res, ucErr := uc.CreateProject(openapi.ProjectCreateRequest{
				User:    openapi.UserOnlyId{Id: tc.userId},
				Project: tc.project,
			})
			assert.NotNil(t, ucErr)

			expectedJson, _ := json.Marshal(tc.expected)
			assert.Equal(t, fmt.Sprintf("domain validation error: %s", expectedJson), ucErr.Error())
			assert.Equal(t, usecase.DomainValidationError, ucErr.Code())
			assert.Equal(t, tc.expected, *ucErr.Response())

			assert.Nil(t, res)
		})
	}
}

func TestCreateProjectServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mock_service.NewMockProjectService(ctrl)

	s.EXPECT().
		CreateProject(gomock.Any(), gomock.Any()).
		Return(nil, service.Errorf(service.RepositoryFailurePanic, "service error"))

	uc := usecase.NewProjectUseCase(s)

	res, ucErr := uc.CreateProject(openapi.ProjectCreateRequest{
		User: openapi.UserOnlyId{Id: testutil.ModifyOnlyUserId()},
		Project: openapi.ProjectWithoutAutofield{
			Name: "Project Name",
		},
	})
	assert.NotNil(t, ucErr)
	assert.Equal(t, "internal error: service error", ucErr.Error())
	assert.Equal(t, usecase.InternalErrorPanic, ucErr.Code())
	assert.Nil(t, ucErr.Response())
	assert.Nil(t, res)
}

func TestUpdateProjectValidEntity(t *testing.T) {
	maxLengthProjectName := testutil.RandomString(100)
	maxLengthProjectDescription := testutil.RandomString(400)

	tt := []struct {
		name    string
		project openapi.ProjectWithoutAutofield
	}{
		{
			name: "should update project with description",
			project: openapi.ProjectWithoutAutofield{
				Name:        "Project With Description",
				Description: "This is a project",
			},
		},
		{
			name: "should update project without description",
			project: openapi.ProjectWithoutAutofield{
				Name: "Project Without Description",
			},
		},
		{
			name: "should update project with max length name",
			project: openapi.ProjectWithoutAutofield{
				Name: maxLengthProjectName,
			},
		},
		{
			name: "should update project with max length description",
			project: openapi.ProjectWithoutAutofield{
				Name:        "Project With Max Length Description",
				Description: maxLengthProjectDescription,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mock_service.NewMockProjectService(ctrl)

			id, err := domain.NewProjectIdObject("0000000000000001")
			assert.NoError(t, err)
			name, err := domain.NewProjectNameObject(tc.project.Name)
			assert.NoError(t, err)
			description, err := domain.NewProjectDescriptionObject(tc.project.Description)
			assert.NoError(t, err)
			createdAt, err := domain.NewCreatedAtObject(testutil.Date())
			assert.NoError(t, err)
			updatedAt, err := domain.NewUpdatedAtObject(testutil.Date())
			assert.NoError(t, err)

			project := domain.NewProjectEntity(*id, *name, *description, *createdAt, *updatedAt)

			s.EXPECT().
				UpdateProject(gomock.Any(), gomock.Any(), gomock.Any()).
				Do(func(userId domain.UserIdObject, projectId domain.ProjectIdObject, project domain.ProjectWithoutAutofieldEntity) {
					assert.Equal(t, testutil.ModifyOnlyUserId(), userId.Value())
					assert.Equal(t, "0000000000000001", projectId.Value())
					assert.Equal(t, tc.project.Name, project.Name().Value())
					assert.Equal(t, tc.project.Description, project.Description().Value())
				}).Return(project, nil)

			uc := usecase.NewProjectUseCase(s)

			res, ucErr := uc.UpdateProject(openapi.ProjectUpdateRequest{
				User:    openapi.UserOnlyId{Id: testutil.ModifyOnlyUserId()},
				Project: openapi.Project{Id: "0000000000000001", Name: tc.project.Name, Description: tc.project.Description},
			})
			assert.Nil(t, ucErr)

			assert.Equal(t, "0000000000000001", res.Project.Id)
			assert.Equal(t, tc.project.Name, res.Project.Name)
			assert.Equal(t, tc.project.Description, res.Project.Description)
		})
	}
}

func TestUpdateProjectDomainValidationError(t *testing.T) {
	tooLongProjectName := testutil.RandomString(101)
	tooLongProjectDescription := testutil.RandomString(401)

	tt := []struct {
		name     string
		userId   string
		project  openapi.Project
		expected openapi.ProjectUpdateErrorResponse
	}{
		{
			name:   "should return error when user id is empty",
			userId: "",
			project: openapi.Project{
				Id:          "0000000000000001",
				Name:        "Project With Description",
				Description: "This is a project",
			},
			expected: openapi.ProjectUpdateErrorResponse{
				User: openapi.UserOnlyIdError{
					Id: "user id is required, but got ''",
				},
			},
		},
		{
			name:   "should return error when project id is empty",
			userId: testutil.ModifyOnlyUserId(),
			project: openapi.Project{
				Id:          "",
				Name:        "Project With Description",
				Description: "This is a project",
			},
			expected: openapi.ProjectUpdateErrorResponse{
				Project: openapi.ProjectError{
					Id: "project id is required, but got ''",
				},
			},
		},
		{
			name:   "should return error when project name is empty",
			userId: testutil.ModifyOnlyUserId(),
			project: openapi.Project{
				Id:          "0000000000000001",
				Name:        "",
				Description: "This is a project",
			},
			expected: openapi.ProjectUpdateErrorResponse{
				Project: openapi.ProjectError{
					Name: "project name is required, but got ''",
				},
			},
		},
		{
			name:   "should return error when project name is too long",
			userId: testutil.ModifyOnlyUserId(),
			project: openapi.Project{
				Id:          "0000000000000001",
				Name:        tooLongProjectName,
				Description: "This is a project",
			},
			expected: openapi.ProjectUpdateErrorResponse{
				Project: openapi.ProjectError{
					Name: fmt.Sprintf(
						"project name cannot be longer than 100 characters, but got '%v'",
						tooLongProjectName,
					),
				},
			},
		},
		{
			name:   "should return error when project description is too long",
			userId: testutil.ModifyOnlyUserId(),
			project: openapi.Project{
				Id:          "0000000000000001",
				Name:        "Project With Description",
				Description: tooLongProjectDescription,
			},
			expected: openapi.ProjectUpdateErrorResponse{
				Project: openapi.ProjectError{
					Description: fmt.Sprintf(
						"project description cannot be longer than 400 characters, but got '%v'",
						tooLongProjectDescription,
					),
				},
			},
		},
		{
			name:    "should return error when all fields are empty",
			userId:  "",
			project: openapi.Project{},
			expected: openapi.ProjectUpdateErrorResponse{
				User: openapi.UserOnlyIdError{
					Id: "user id is required, but got ''",
				},
				Project: openapi.ProjectError{
					Id:   "project id is required, but got ''",
					Name: "project name is required, but got ''",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mock_service.NewMockProjectService(ctrl)

			uc := usecase.NewProjectUseCase(s)

			res, ucErr := uc.UpdateProject(openapi.ProjectUpdateRequest{
				User:    openapi.UserOnlyId{Id: tc.userId},
				Project: tc.project,
			})
			assert.NotNil(t, ucErr)

			expectedJson, _ := json.Marshal(tc.expected)
			assert.Equal(t, fmt.Sprintf("domain validation error: %s", expectedJson), ucErr.Error())
			assert.Equal(t, usecase.DomainValidationError, ucErr.Code())
			assert.Equal(t, tc.expected, *ucErr.Response())

			assert.Nil(t, res)
		})
	}
}

func TestUpdateProjectServiceError(t *testing.T) {
	tt := []struct {
		name          string
		errorCode     service.ErrorCode
		errorMessage  string
		expectedError string
		expectedCode  usecase.ErrorCode
	}{
		{
			name:          "should return error when project not found",
			errorCode:     service.NotFoundError,
			errorMessage:  "failed to update project",
			expectedError: "not found: failed to update project",
			expectedCode:  usecase.NotFoundError,
		},
		{
			name:          "should return error when repository failure",
			errorCode:     service.RepositoryFailurePanic,
			errorMessage:  "service error",
			expectedError: "internal error: service error",
			expectedCode:  usecase.InternalErrorPanic,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mock_service.NewMockProjectService(ctrl)

			s.EXPECT().
				UpdateProject(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, service.Errorf(tc.errorCode, "%s", tc.errorMessage))

			uc := usecase.NewProjectUseCase(s)

			res, ucErr := uc.UpdateProject(openapi.ProjectUpdateRequest{
				User: openapi.UserOnlyId{Id: testutil.ModifyOnlyUserId()},
				Project: openapi.Project{
					Id:   "0000000000000001",
					Name: "Project Name",
				},
			})
			assert.NotNil(t, ucErr)
			assert.Equal(t, tc.expectedError, ucErr.Error())
			assert.Equal(t, tc.expectedCode, ucErr.Code())
			assert.Nil(t, ucErr.Response())
			assert.Nil(t, res)
		})
	}
}

func TestDeleteProjectValidEntity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mock_service.NewMockProjectService(ctrl)

	s.EXPECT().
		DeleteProject(gomock.Any(), gomock.Any()).
		Do(func(userId domain.UserIdObject, projectId domain.ProjectIdObject) {
			assert.Equal(t, testutil.ModifyOnlyUserId(), userId.Value())
			assert.Equal(t, "0000000000000001", projectId.Value())
		}).
		Return(nil)

	uc := usecase.NewProjectUseCase(s)

	ucErr := uc.DeleteProject(openapi.ProjectDeleteRequest{
		User:    openapi.UserOnlyId{Id: testutil.ModifyOnlyUserId()},
		Project: openapi.ProjectOnlyId{Id: "0000000000000001"},
	})
	assert.Nil(t, ucErr)
}

func TestDeleteProjectDomainValidationError(t *testing.T) {
	tt := []struct {
		name      string
		userId    string
		projectId string
		expected  openapi.ProjectDeleteErrorResponse
	}{
		{
			name:      "should return error when user id is empty",
			userId:    "",
			projectId: "0000000000000001",
			expected: openapi.ProjectDeleteErrorResponse{
				User:    openapi.UserOnlyIdError{Id: "user id is required, but got ''"},
				Project: openapi.ProjectOnlyIdError{Id: ""},
			},
		},
		{
			name:      "should return error when project id is empty",
			userId:    testutil.ModifyOnlyUserId(),
			projectId: "",
			expected: openapi.ProjectDeleteErrorResponse{
				User:    openapi.UserOnlyIdError{Id: ""},
				Project: openapi.ProjectOnlyIdError{Id: "project id is required, but got ''"},
			},
		},
		{
			name:      "should return error when all fields are empty",
			userId:    "",
			projectId: "",
			expected: openapi.ProjectDeleteErrorResponse{
				User:    openapi.UserOnlyIdError{Id: "user id is required, but got ''"},
				Project: openapi.ProjectOnlyIdError{Id: "project id is required, but got ''"},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mock_service.NewMockProjectService(ctrl)

			uc := usecase.NewProjectUseCase(s)

			ucErr := uc.DeleteProject(openapi.ProjectDeleteRequest{
				User:    openapi.UserOnlyId{Id: tc.userId},
				Project: openapi.ProjectOnlyId{Id: tc.projectId},
			})
			assert.NotNil(t, ucErr)

			expectedJson, _ := json.Marshal(tc.expected)
			assert.Equal(t, fmt.Sprintf("domain validation error: %s", expectedJson), ucErr.Error())
			assert.Equal(t, usecase.DomainValidationError, ucErr.Code())
			assert.Equal(t, tc.expected, *ucErr.Response())
		})
	}
}

func TestDeleteProjectServiceError(t *testing.T) {
	tt := []struct {
		name          string
		errorCode     service.ErrorCode
		errorMessage  string
		expectedError string
		expectedCode  usecase.ErrorCode
	}{
		{
			name:          "should return error when project not found",
			errorCode:     service.NotFoundError,
			errorMessage:  "failed to delete project",
			expectedError: "not found: failed to delete project",
			expectedCode:  usecase.NotFoundError,
		},
		{
			name:          "should return error when repository failure",
			errorCode:     service.RepositoryFailurePanic,
			errorMessage:  "service error",
			expectedError: "internal error: service error",
			expectedCode:  usecase.InternalErrorPanic,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mock_service.NewMockProjectService(ctrl)

			s.EXPECT().
				DeleteProject(gomock.Any(), gomock.Any()).
				Return(service.Errorf(tc.errorCode, "%s", tc.errorMessage))

			uc := usecase.NewProjectUseCase(s)

			ucErr := uc.DeleteProject(openapi.ProjectDeleteRequest{
				User:    openapi.UserOnlyId{Id: testutil.ModifyOnlyUserId()},
				Project: openapi.ProjectOnlyId{Id: "0000000000000001"},
			})
			assert.NotNil(t, ucErr)
			assert.Equal(t, tc.expectedError, ucErr.Error())
			assert.Equal(t, tc.expectedCode, ucErr.Code())
			assert.Nil(t, ucErr.Response())
		})
	}
}
