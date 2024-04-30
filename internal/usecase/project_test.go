package usecase_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/kumachan-mis/knodeledge-api/internal/domain"
	"github.com/kumachan-mis/knodeledge-api/internal/model"
	"github.com/kumachan-mis/knodeledge-api/internal/service"
	"github.com/kumachan-mis/knodeledge-api/internal/usecase"
	mock_service "github.com/kumachan-mis/knodeledge-api/mock/service"
	"github.com/kumachan-mis/knodeledge-api/test/testutil"
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

	res, ucErr := uc.ListProjects(model.ProjectListRequest{
		User: model.UserOnlyId{Id: testutil.ReadOnlyUserId()},
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
		expected model.ProjectListErrorResponse
	}{
		{
			name:   "empty user id",
			userId: "",
			expected: model.ProjectListErrorResponse{
				User: model.UserOnlyIdError{Id: "user id is required, but got ''"},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			s := mock_service.NewMockProjectService(ctrl)

			uc := usecase.NewProjectUseCase(s)

			res, ucErr := uc.ListProjects(model.ProjectListRequest{
				User: model.UserOnlyId{Id: tc.userId},
			})
			assert.Error(t, ucErr)

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

	res, ucErr := uc.ListProjects(model.ProjectListRequest{
		User: model.UserOnlyId{Id: testutil.ReadOnlyUserId()},
	})
	assert.Error(t, ucErr)
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

		res, ucErr := uc.FindProject(model.ProjectFindRequest{
			User:    model.UserOnlyId{Id: testutil.ReadOnlyUserId()},
			Project: model.ProjectOnlyId{Id: tc.projectId},
		})
		assert.Nil(t, ucErr)

		assert.Equal(t, tc.projectId, res.Project.Id)
		assert.Equal(t, tc.projectName, res.Project.Name)
		assert.Equal(t, tc.projectDescription, res.Project.Description)
	}
}

func TestFindProjectDomainValidationError(t *testing.T) {
	tt := []struct {
		name      string
		userId    string
		projectId string
		expected  model.ProjectFindErrorResponse
	}{
		{
			name:      "empty user id",
			userId:    "",
			projectId: "0000000000000001",
			expected: model.ProjectFindErrorResponse{
				User:    model.UserOnlyIdError{Id: "user id is required, but got ''"},
				Project: model.ProjectOnlyIdError{Id: ""},
			},
		},
		{
			name:      "empty project id",
			userId:    testutil.ReadOnlyUserId(),
			projectId: "",
			expected: model.ProjectFindErrorResponse{
				User:    model.UserOnlyIdError{Id: ""},
				Project: model.ProjectOnlyIdError{Id: "project id is required, but got ''"},
			},
		},
		{
			name:      "empty all fields",
			userId:    "",
			projectId: "",
			expected: model.ProjectFindErrorResponse{
				User:    model.UserOnlyIdError{Id: "user id is required, but got ''"},
				Project: model.ProjectOnlyIdError{Id: "project id is required, but got ''"},
			},
		},
	}

	for _, tc := range tt {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		s := mock_service.NewMockProjectService(ctrl)

		uc := usecase.NewProjectUseCase(s)

		res, ucErr := uc.FindProject(model.ProjectFindRequest{
			User:    model.UserOnlyId{Id: tc.userId},
			Project: model.ProjectOnlyId{Id: tc.projectId},
		})
		assert.Error(t, ucErr)

		expectedJson, _ := json.Marshal(tc.expected)
		assert.Equal(t, fmt.Sprintf("domain validation error: %s", expectedJson), ucErr.Error())
		assert.Equal(t, usecase.DomainValidationError, ucErr.Code())
		assert.Equal(t, tc.expected, *ucErr.Response())

		assert.Nil(t, res)
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
			name:          "not found",
			errorCode:     service.NotFoundError,
			errorMessage:  "failed to find project",
			expectedError: "not found: failed to find project",
			expectedCode:  usecase.NotFoundError,
		},
		{
			name:          "internal error",
			errorCode:     service.RepositoryFailurePanic,
			errorMessage:  "service error",
			expectedError: "internal error: service error",
			expectedCode:  usecase.InternalErrorPanic,
		},
	}

	for _, tc := range tt {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		s := mock_service.NewMockProjectService(ctrl)

		s.EXPECT().
			FindProject(gomock.Any(), gomock.Any()).
			Return(nil, service.Errorf(tc.errorCode, tc.errorMessage))

		uc := usecase.NewProjectUseCase(s)

		res, ucErr := uc.FindProject(model.ProjectFindRequest{
			User:    model.UserOnlyId{Id: testutil.ReadOnlyUserId()},
			Project: model.ProjectOnlyId{Id: "0000000000000001"},
		})
		assert.Error(t, ucErr)
		assert.Equal(t, tc.expectedError, ucErr.Error())
		assert.Equal(t, tc.expectedCode, ucErr.Code())
		assert.Nil(t, ucErr.Response())
		assert.Nil(t, res)
	}
}

func TestCreateProjectValidEntity(t *testing.T) {
	maxLengthProjectName := testutil.RandomString(100)
	maxLengthProjectDescription := testutil.RandomString(400)

	tt := []struct {
		name    string
		project model.ProjectWithoutAutofield
	}{
		{
			name: "should create project with description",
			project: model.ProjectWithoutAutofield{
				Name:        "Project With Description",
				Description: "This is a project",
			},
		},
		{
			name: "should create project without description",
			project: model.ProjectWithoutAutofield{
				Name: "Project Without Description",
			},
		},
		{
			name: "should create project with max length name",
			project: model.ProjectWithoutAutofield{
				Name: maxLengthProjectName,
			},
		},
		{
			name: "should create project with max length description",
			project: model.ProjectWithoutAutofield{
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

			res, ucErr := uc.CreateProject(model.ProjectCreateRequest{
				User:    model.UserOnlyId{Id: testutil.ModifyOnlyUserId()},
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
		project  model.ProjectWithoutAutofield
		expected model.ProjectCreateErrorResponse
	}{
		{
			name:   "should return error when user id is empty",
			userId: "",
			project: model.ProjectWithoutAutofield{
				Name:        "Project With Description",
				Description: "This is a project",
			},
			expected: model.ProjectCreateErrorResponse{
				User: model.UserOnlyIdError{
					Id: "user id is required, but got ''",
				},
			},
		},
		{
			name:   "should return error when project name is empty",
			userId: testutil.ModifyOnlyUserId(),
			project: model.ProjectWithoutAutofield{
				Name:        "",
				Description: "This is a project",
			},
			expected: model.ProjectCreateErrorResponse{
				Project: model.ProjectWithoutAutofieldError{
					Name: "project name is required, but got ''",
				},
			},
		},
		{
			name:   "should return error when project name is too long",
			userId: testutil.ModifyOnlyUserId(),
			project: model.ProjectWithoutAutofield{
				Name: tooLongProjectName,
			},
			expected: model.ProjectCreateErrorResponse{
				Project: model.ProjectWithoutAutofieldError{
					Name: fmt.Sprintf(
						"project name cannot be longer than 100 characters, but got '%s'",
						tooLongProjectName,
					),
				},
			},
		},
		{
			name:   "should return error when project description is too long",
			userId: testutil.ModifyOnlyUserId(),
			project: model.ProjectWithoutAutofield{
				Name:        "Project With Description",
				Description: tooLongProjectDescription,
			},
			expected: model.ProjectCreateErrorResponse{
				Project: model.ProjectWithoutAutofieldError{
					Description: fmt.Sprintf(
						"project description cannot be longer than 400 characters, but got '%s'",
						tooLongProjectDescription,
					),
				},
			},
		},
		{
			name:    "should return error when all fields are empty",
			userId:  "",
			project: model.ProjectWithoutAutofield{},
			expected: model.ProjectCreateErrorResponse{
				User: model.UserOnlyIdError{
					Id: "user id is required, but got ''",
				},
				Project: model.ProjectWithoutAutofieldError{
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

			res, ucErr := uc.CreateProject(model.ProjectCreateRequest{
				User:    model.UserOnlyId{Id: tc.userId},
				Project: tc.project,
			})
			assert.Error(t, ucErr)

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

	res, ucErr := uc.CreateProject(model.ProjectCreateRequest{
		User: model.UserOnlyId{Id: testutil.ModifyOnlyUserId()},
		Project: model.ProjectWithoutAutofield{
			Name: "Project Name",
		},
	})
	assert.Error(t, ucErr)
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
		project model.ProjectWithoutAutofield
	}{
		{
			name: "should update project with description",
			project: model.ProjectWithoutAutofield{
				Name:        "Project With Description",
				Description: "This is a project",
			},
		},
		{
			name: "should update project without description",
			project: model.ProjectWithoutAutofield{
				Name: "Project Without Description",
			},
		},
		{
			name: "should update project with max length name",
			project: model.ProjectWithoutAutofield{
				Name: maxLengthProjectName,
			},
		},
		{
			name: "should update project with max length description",
			project: model.ProjectWithoutAutofield{
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

			res, ucErr := uc.UpdateProject(model.ProjectUpdateRequest{
				User:    model.UserOnlyId{Id: testutil.ModifyOnlyUserId()},
				Project: model.Project{Id: "0000000000000001", Name: tc.project.Name, Description: tc.project.Description},
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
		project  model.Project
		expected model.ProjectUpdateErrorResponse
	}{
		{
			name:   "should return error when user id is empty",
			userId: "",
			project: model.Project{
				Id:          "0000000000000001",
				Name:        "Project With Description",
				Description: "This is a project",
			},
			expected: model.ProjectUpdateErrorResponse{
				User: model.UserOnlyIdError{
					Id: "user id is required, but got ''",
				},
			},
		},
		{
			name:   "should return error when project id is empty",
			userId: testutil.ModifyOnlyUserId(),
			project: model.Project{
				Id:          "",
				Name:        "Project With Description",
				Description: "This is a project",
			},
			expected: model.ProjectUpdateErrorResponse{
				Project: model.ProjectError{
					Id: "project id is required, but got ''",
				},
			},
		},
		{
			name:   "should return error when project name is empty",
			userId: testutil.ModifyOnlyUserId(),
			project: model.Project{
				Id:          "0000000000000001",
				Name:        "",
				Description: "This is a project",
			},
			expected: model.ProjectUpdateErrorResponse{
				Project: model.ProjectError{
					Name: "project name is required, but got ''",
				},
			},
		},
		{
			name:   "should return error when project name is too long",
			userId: testutil.ModifyOnlyUserId(),
			project: model.Project{
				Id:          "0000000000000001",
				Name:        tooLongProjectName,
				Description: "This is a project",
			},
			expected: model.ProjectUpdateErrorResponse{
				Project: model.ProjectError{
					Name: fmt.Sprintf(
						"project name cannot be longer than 100 characters, but got '%s'",
						tooLongProjectName,
					),
				},
			},
		},
		{
			name:   "should return error when project description is too long",
			userId: testutil.ModifyOnlyUserId(),
			project: model.Project{
				Id:          "0000000000000001",
				Name:        "Project With Description",
				Description: tooLongProjectDescription,
			},
			expected: model.ProjectUpdateErrorResponse{
				Project: model.ProjectError{
					Description: fmt.Sprintf(
						"project description cannot be longer than 400 characters, but got '%s'",
						tooLongProjectDescription,
					),
				},
			},
		},
		{
			name:    "should return error when all fields are empty",
			userId:  "",
			project: model.Project{},
			expected: model.ProjectUpdateErrorResponse{
				User: model.UserOnlyIdError{
					Id: "user id is required, but got ''",
				},
				Project: model.ProjectError{
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

			res, ucErr := uc.UpdateProject(model.ProjectUpdateRequest{
				User:    model.UserOnlyId{Id: tc.userId},
				Project: tc.project,
			})
			assert.Error(t, ucErr)

			expectedJson, _ := json.Marshal(tc.expected)
			assert.Equal(t, fmt.Sprintf("domain validation error: %s", expectedJson), ucErr.Error())
			assert.Equal(t, usecase.DomainValidationError, ucErr.Code())
			assert.Equal(t, tc.expected, *ucErr.Response())

			assert.Nil(t, res)
		})
	}
}

func TestUpdateProjectServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mock_service.NewMockProjectService(ctrl)

	s.EXPECT().
		UpdateProject(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, service.Errorf(service.RepositoryFailurePanic, "service error"))

	uc := usecase.NewProjectUseCase(s)

	res, ucErr := uc.UpdateProject(model.ProjectUpdateRequest{
		User: model.UserOnlyId{Id: testutil.ModifyOnlyUserId()},
		Project: model.Project{
			Id:   "0000000000000001",
			Name: "Project Name",
		},
	})
	assert.Error(t, ucErr)
	assert.Equal(t, "internal error: service error", ucErr.Error())
	assert.Equal(t, usecase.InternalErrorPanic, ucErr.Code())
	assert.Nil(t, ucErr.Response())
	assert.Nil(t, res)
}
