package service_test

import (
	"fmt"
	"testing"

	"github.com/kumachan-mis/knodeledge-api/internal/domain"
	"github.com/kumachan-mis/knodeledge-api/internal/record"
	"github.com/kumachan-mis/knodeledge-api/internal/service"
	"github.com/kumachan-mis/knodeledge-api/mock/repository"
	"github.com/kumachan-mis/knodeledge-api/test/testutil"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestListProjectsValidEntry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	maxLengthProjectName := testutil.RandomString(100)
	maxLengthProjectDescription := testutil.RandomString(400)

	r := repository.NewMockProjectRepository(ctrl)
	r.EXPECT().
		FetchUserProjects(testutil.UserId()).
		Return(map[string]record.ProjectEntry{
			"0000000000000001": {
				Name:        "First Project",
				Description: "This is my first project",
				UserId:      testutil.UserId(),
			},
			"0000000000000002": {
				Name:   "Second Project",
				UserId: testutil.UserId(),
			},
			"0000000000000003": {
				Name:        maxLengthProjectName,
				Description: maxLengthProjectDescription,
				UserId:      testutil.UserId(),
			},
		}, nil)

	s := service.NewProjectService(r)

	userId, err := domain.NewUserIdObject(testutil.UserId())
	assert.NoError(t, err)

	projects, err := s.ListProjects(*userId)
	assert.NoError(t, err)

	assert.Len(t, projects, 3)

	project := projects[0]
	assert.Equal(t, "0000000000000001", project.Id().Value())
	assert.Equal(t, "First Project", project.Name().Value())
	assert.Equal(t, "This is my first project", project.Description().Value())

	project = projects[1]
	assert.Equal(t, "0000000000000002", project.Id().Value())
	assert.Equal(t, "Second Project", project.Name().Value())
	assert.Equal(t, "", project.Description().Value())

	project = projects[2]
	assert.Equal(t, "0000000000000003", project.Id().Value())
	assert.Equal(t, maxLengthProjectName, project.Name().Value())
	assert.Equal(t, maxLengthProjectDescription, project.Description().Value())

}

func TestListProjectsNoEntry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := repository.NewMockProjectRepository(ctrl)
	r.EXPECT().
		FetchUserProjects(testutil.UserId()).
		Return(map[string]record.ProjectEntry{}, nil)

	s := service.NewProjectService(r)

	userId, err := domain.NewUserIdObject(testutil.UserId())
	assert.NoError(t, err)

	projects, err := s.ListProjects(*userId)
	assert.NoError(t, err)

	assert.Empty(t, projects)
}

func TestListProjectsInvalidEntry(t *testing.T) {
	tooLongProjectName := testutil.RandomString(101)
	tooLongProjectDescription := testutil.RandomString(401)

	tt := []struct {
		name          string
		projectId     string
		project       record.ProjectEntry
		expectedError string
	}{
		{
			name:      "should return error when project id is empty",
			projectId: "",
			project: record.ProjectEntry{
				Name:   "Project",
				UserId: testutil.UserId(),
			},
			expectedError: "failed to convert entry to entity (id): project id is required, but got ''",
		},
		{
			name:      "should return error when project name is empty",
			projectId: "0000000000000001",
			project: record.ProjectEntry{
				Name:   "",
				UserId: testutil.UserId(),
			},
			expectedError: "failed to convert entry to entity (name): project name is required, but got ''",
		},
		{
			name:      "should return error when project name is too long",
			projectId: "0000000000000001",
			project: record.ProjectEntry{
				Name:   tooLongProjectName,
				UserId: testutil.UserId(),
			},
			expectedError: "failed to convert entry to entity (name): " +
				fmt.Sprintf(
					"project name cannot be longer than 100 characters, but got '%s'",
					tooLongProjectName,
				),
		},
		{
			name:      "should return error when project description is too long",
			projectId: "0000000000000001",
			project: record.ProjectEntry{
				Name:        "Project",
				Description: tooLongProjectDescription,
				UserId:      testutil.UserId(),
			},
			expectedError: "failed to convert entry to entity (description): " +
				fmt.Sprintf(
					"project description cannot be longer than 400 characters, but got '%s'",
					tooLongProjectDescription,
				),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := repository.NewMockProjectRepository(ctrl)
			r.EXPECT().
				FetchUserProjects(testutil.UserId()).
				Return(map[string]record.ProjectEntry{
					tc.projectId: tc.project,
				}, nil)

			s := service.NewProjectService(r)

			userId, err := domain.NewUserIdObject(testutil.UserId())
			assert.NoError(t, err)

			projects, err := s.ListProjects(*userId)
			assert.ErrorContains(t, err, tc.expectedError)
			assert.Nil(t, projects)
		})
	}
}

func TestListProjectsRepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := repository.NewMockProjectRepository(ctrl)
	r.EXPECT().
		FetchUserProjects(testutil.UserId()).
		Return(nil, fmt.Errorf("repository error"))

	s := service.NewProjectService(r)

	userId, err := domain.NewUserIdObject(testutil.UserId())
	assert.NoError(t, err)

	projects, err := s.ListProjects(*userId)
	assert.ErrorContains(t, err, "failed to fetch user projects: repository error")
	assert.Nil(t, projects)
}
