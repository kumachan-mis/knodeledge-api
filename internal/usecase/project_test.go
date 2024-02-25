package usecase_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/kumachan-mis/knodeledge-api/internal/domain"
	"github.com/kumachan-mis/knodeledge-api/internal/model"
	"github.com/kumachan-mis/knodeledge-api/internal/usecase"
	"github.com/kumachan-mis/knodeledge-api/mock/service"
	"github.com/kumachan-mis/knodeledge-api/test/testutil"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestListProjectsValidEntity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := service.NewMockProjectService(ctrl)

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
			assert.Equal(t, testutil.UserId(), userId.Value())
		}).
		Return([]domain.ProjectEntity{*projectWithDesc, *projectWithoutDesc}, nil)

	uc := usecase.NewProjectUseCase(s)

	res, err := uc.ListProjects(model.ProjectListRequest{
		User: model.User{Id: testutil.UserId()},
	})
	assert.Nil(t, err)

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

func TestListProjectsInvalidArgument(t *testing.T) {
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
				User: model.UserError{Id: "user id is required, but got ''"},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			s := service.NewMockProjectService(ctrl)

			uc := usecase.NewProjectUseCase(s)

			res, ucErr := uc.ListProjects(model.ProjectListRequest{
				User: model.User{Id: tc.userId},
			})
			assert.Error(t, ucErr)

			expectedJson, _ := json.Marshal(tc.expected)
			assert.Equal(t, fmt.Sprintf("invalid argument: %s", expectedJson), ucErr.Error())
			assert.Equal(t, usecase.ErrorCode("invalid argument"), ucErr.Code())
			assert.Equal(t, tc.expected, *ucErr.Response())

			assert.Nil(t, res)
		})
	}
}

func TestListProjectsServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := service.NewMockProjectService(ctrl)

	s.EXPECT().
		ListProjects(gomock.Any()).
		Return(nil, fmt.Errorf("service error"))

	uc := usecase.NewProjectUseCase(s)

	res, ucErr := uc.ListProjects(model.ProjectListRequest{
		User: model.User{Id: testutil.UserId()},
	})
	assert.Error(t, ucErr)
	assert.Equal(t, usecase.ErrorCode("internal error"), ucErr.Code())
	assert.Nil(t, ucErr.Response())
	assert.Equal(t, "internal error: service error", ucErr.Error())
	assert.Nil(t, res)
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

			s := service.NewMockProjectService(ctrl)

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

			projectWithDesc := domain.NewProjectEntity(*id, *name, *description, *createdAt, *updatedAt)

			s.EXPECT().
				CreateProject(gomock.Any(), gomock.Any()).
				Do(func(userId domain.UserIdObject, project domain.ProjectWithoutAutofieldEntity) {
					assert.Equal(t, testutil.UserId(), userId.Value())
					assert.Equal(t, tc.project.Name, project.Name().Value())
					assert.Equal(t, tc.project.Description, project.Description().Value())
				}).
				Return(projectWithDesc, nil)

			uc := usecase.NewProjectUseCase(s)

			res, err := uc.CreateProject(model.ProjectCreateRequest{
				User:    model.User{Id: testutil.UserId()},
				Project: tc.project,
			})
			assert.Nil(t, err)

			assert.Equal(t, "0000000000000001", res.Project.Id)
			assert.Equal(t, tc.project.Name, res.Project.Name)
			assert.Equal(t, tc.project.Description, res.Project.Description)

		})
	}
}

func TestCreateProjectInvalidArgument(t *testing.T) {
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
				User: model.UserError{
					Id: "user id is required, but got ''",
				},
			},
		},
		{
			name:   "should return error when project name is empty",
			userId: testutil.UserId(),
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
			userId: testutil.UserId(),
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
			userId: testutil.UserId(),
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
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := service.NewMockProjectService(ctrl)

			uc := usecase.NewProjectUseCase(s)

			res, ucErr := uc.CreateProject(model.ProjectCreateRequest{
				User:    model.User{Id: tc.userId},
				Project: tc.project,
			})
			assert.Error(t, ucErr)

			expectedJson, _ := json.Marshal(tc.expected)
			assert.Equal(t, fmt.Sprintf("invalid argument: %s", expectedJson), ucErr.Error())
			assert.Equal(t, usecase.ErrorCode("invalid argument"), ucErr.Code())
			assert.Equal(t, tc.expected, *ucErr.Response())

			assert.Nil(t, res)
		})
	}
}

func TestCreateProjectServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := service.NewMockProjectService(ctrl)

	s.EXPECT().
		CreateProject(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("service error"))

	uc := usecase.NewProjectUseCase(s)

	res, ucErr := uc.CreateProject(model.ProjectCreateRequest{
		User: model.User{Id: testutil.UserId()},
		Project: model.ProjectWithoutAutofield{
			Name: "Project Name",
		},
	})
	assert.Error(t, ucErr)
	assert.Equal(t, usecase.ErrorCode("internal error"), ucErr.Code())
	assert.Nil(t, ucErr.Response())
	assert.Equal(t, "internal error: service error", ucErr.Error())
	assert.Nil(t, res)
}
