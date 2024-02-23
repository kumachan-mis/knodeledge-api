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
	createdAt, err := domain.NewCreatedAtObject(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	assert.NoError(t, err)
	updatedAt, err := domain.NewUpdatedAtObject(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	assert.NoError(t, err)
	projectWithDesc, err := domain.NewProjectEntity(*id, *name, *description, *createdAt, *updatedAt)
	assert.NoError(t, err)

	id, err = domain.NewProjectIdObject("0000000000000002")
	assert.NoError(t, err)
	name, err = domain.NewProjectNameObject("Project Without Description")
	assert.NoError(t, err)
	description, err = domain.NewProjectDescriptionObject("")
	assert.NoError(t, err)
	createdAt, err = domain.NewCreatedAtObject(time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC))
	assert.NoError(t, err)
	updatedAt, err = domain.NewUpdatedAtObject(time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC))
	assert.NoError(t, err)
	projectWithoutDesc, err := domain.NewProjectEntity(*id, *name, *description, *createdAt, *updatedAt)
	assert.NoError(t, err)

	s.EXPECT().
		ListProjects(gomock.Any()).
		Do(func(userId domain.UserIdObject) {
			assert.Equal(t, testutil.UserId(), userId.Value())
		}).
		Return([]domain.ProjectEntity{
			*projectWithDesc,
			*projectWithoutDesc,
		}, nil)

	uc := usecase.NewProjectUseCase(s)

	projects, err := uc.ListProjects(model.User{Id: testutil.UserId()})
	assert.Nil(t, err)

	assert.Len(t, projects, 2)

	project := projects[0]
	assert.Equal(t, "0000000000000001", project.Id)
	assert.Equal(t, "Project With Description", project.Name)
	assert.Equal(t, "This is a project", project.Description)

	project = projects[1]
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
		expected model.UserError
	}{
		{
			name:     "empty user id",
			userId:   "",
			expected: model.UserError{Id: "user id is required, but got ''"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			s := service.NewMockProjectService(ctrl)

			uc := usecase.NewProjectUseCase(s)

			projects, ucErr := uc.ListProjects(model.User{Id: tc.userId})
			assert.Error(t, ucErr)

			expectedJson, _ := json.Marshal(tc.expected)
			assert.Equal(t, fmt.Sprintf("invalid argument: %s", expectedJson), ucErr.Error())
			assert.Equal(t, usecase.ErrorCode("invalid argument"), ucErr.Code())
			assert.Equal(t, tc.expected, *ucErr.Model())

			assert.Nil(t, projects)
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

	projects, ucErr := uc.ListProjects(model.User{Id: testutil.UserId()})
	assert.Error(t, ucErr)
	assert.Equal(t, usecase.ErrorCode("internal error"), ucErr.Code())
	assert.Nil(t, ucErr.Model())
	assert.Equal(t, "internal error: service error", ucErr.Error())
	assert.Nil(t, projects)
}
