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
	projectWithDesc, err := domain.NewProjectEntity(*id, *name, *description, *createdAt, *updatedAt)
	assert.NoError(t, err)

	id, err = domain.NewProjectIdObject("0000000000000002")
	assert.NoError(t, err)
	name, err = domain.NewProjectNameObject("Project Without Description")
	assert.NoError(t, err)
	description, err = domain.NewProjectDescriptionObject("")
	assert.NoError(t, err)
	createdAt, err = domain.NewCreatedAtObject(testutil.Date().Add(1 * time.Hour))
	assert.NoError(t, err)
	updatedAt, err = domain.NewUpdatedAtObject(testutil.Date().Add(1 * time.Hour))
	assert.NoError(t, err)
	projectWithoutDesc, err := domain.NewProjectEntity(*id, *name, *description, *createdAt, *updatedAt)
	assert.NoError(t, err)

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

func TestListProjectsInvalidUserId(t *testing.T) {
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
