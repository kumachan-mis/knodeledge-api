package service_test

import (
	"fmt"
	"testing"
	"time"

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
			"0000000000000003": {
				Name:        maxLengthProjectName,
				Description: maxLengthProjectDescription,
				UserId:      testutil.UserId(),
				CreatedAt:   testutil.Date().Add(-3 * time.Hour),
				UpdatedAt:   testutil.Date().Add(-3 * time.Hour),
			},
			"0000000000000002": {
				Name:      "Second Project",
				UserId:    testutil.UserId(),
				CreatedAt: testutil.Date().Add(-2 * time.Hour),
				UpdatedAt: testutil.Date().Add(-2 * time.Hour),
			},
			"0000000000000001": {
				Name:        "First Project",
				Description: "This is my first project",
				UserId:      testutil.UserId(),
				CreatedAt:   testutil.Date().Add(-1 * time.Hour),
				UpdatedAt:   testutil.Date().Add(-1 * time.Hour),
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
	assert.Equal(t, testutil.Date().Add(-1*time.Hour), project.CreatedAt().Value())
	assert.Equal(t, testutil.Date().Add(-1*time.Hour), project.UpdatedAt().Value())

	project = projects[1]
	assert.Equal(t, "0000000000000002", project.Id().Value())
	assert.Equal(t, "Second Project", project.Name().Value())
	assert.Equal(t, "", project.Description().Value())
	assert.Equal(t, testutil.Date().Add(-2*time.Hour), project.CreatedAt().Value())
	assert.Equal(t, testutil.Date().Add(-2*time.Hour), project.UpdatedAt().Value())

	project = projects[2]
	assert.Equal(t, "0000000000000003", project.Id().Value())
	assert.Equal(t, maxLengthProjectName, project.Name().Value())
	assert.Equal(t, maxLengthProjectDescription, project.Description().Value())
	assert.Equal(t, testutil.Date().Add(-3*time.Hour), project.CreatedAt().Value())
	assert.Equal(t, testutil.Date().Add(-3*time.Hour), project.UpdatedAt().Value())
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
				Name:      "Project",
				UserId:    testutil.UserId(),
				CreatedAt: testutil.Date(),
				UpdatedAt: testutil.Date(),
			},
			expectedError: "failed to convert entry to entity (id): project id is required, but got ''",
		},
		{
			name:      "should return error when project name is empty",
			projectId: "0000000000000001",
			project: record.ProjectEntry{
				Name:      "",
				UserId:    testutil.UserId(),
				CreatedAt: testutil.Date(),
				UpdatedAt: testutil.Date(),
			},
			expectedError: "failed to convert entry to entity (name): project name is required, but got ''",
		},
		{
			name:      "should return error when project name is too long",
			projectId: "0000000000000001",
			project: record.ProjectEntry{
				Name:      tooLongProjectName,
				UserId:    testutil.UserId(),
				CreatedAt: testutil.Date(),
				UpdatedAt: testutil.Date(),
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
				CreatedAt:   testutil.Date(),
				UpdatedAt:   testutil.Date(),
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

func TestCreateProjectValidEntry(t *testing.T) {
	maxLengthProjectName := testutil.RandomString(100)
	maxLengthProjectDescription := testutil.RandomString(400)

	tt := []struct {
		name      string
		projectId string
		project   record.ProjectWithoutAutofieldEntry
	}{
		{
			name:      "should return project with valid entry",
			projectId: "0000000000000001",
			project: record.ProjectWithoutAutofieldEntry{
				Name:        "New Project",
				Description: "This is new project",
				UserId:      testutil.UserId(),
			},
		},
		{
			name:      "should return project with max-length valid entry",
			projectId: "0000000000000002",
			project: record.ProjectWithoutAutofieldEntry{
				Name:        maxLengthProjectName,
				Description: maxLengthProjectDescription,
				UserId:      testutil.UserId(),
			},
		},
	}

	for _, tc := range tt {

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		r := repository.NewMockProjectRepository(ctrl)
		r.EXPECT().
			InsertProject(tc.project).
			Return(tc.projectId, &record.ProjectEntry{
				Name:        tc.project.Name,
				Description: tc.project.Description,
				UserId:      tc.project.UserId,
				CreatedAt:   testutil.Date(),
				UpdatedAt:   testutil.Date(),
			}, nil)

		s := service.NewProjectService(r)

		userId, err := domain.NewUserIdObject(testutil.UserId())
		assert.NoError(t, err)

		name, err := domain.NewProjectNameObject(tc.project.Name)
		assert.NoError(t, err)
		description, err := domain.NewProjectDescriptionObject(tc.project.Description)
		assert.NoError(t, err)

		project := domain.NewProjectWithoutAutofieldEntity(*name, *description)

		createdProject, err := s.CreateProject(*userId, *project)
		assert.NoError(t, err)

		assert.Equal(t, tc.projectId, createdProject.Id().Value())
		assert.Equal(t, tc.project.Name, createdProject.Name().Value())
		assert.Equal(t, tc.project.Description, createdProject.Description().Value())
		assert.Equal(t, testutil.Date(), createdProject.CreatedAt().Value())
		assert.Equal(t, testutil.Date(), createdProject.UpdatedAt().Value())
	}
}

func TestCreateProjectInvalidCreatedEntry(t *testing.T) {
	tooLongProjectName := testutil.RandomString(101)
	tooLongProjectDescription := testutil.RandomString(401)

	tt := []struct {
		name           string
		createdProject record.ProjectEntry
		expectedError  string
	}{
		{
			name: "should return error when project name is empty",
			createdProject: record.ProjectEntry{
				Name:        "",
				Description: "This is new project",
				UserId:      testutil.UserId(),
				CreatedAt:   testutil.Date(),
				UpdatedAt:   testutil.Date(),
			},
			expectedError: "failed to convert entry to entity (name): project name is required, but got ''",
		},
		{
			name: "should return error when project name is too long",
			createdProject: record.ProjectEntry{
				Name:        tooLongProjectName,
				Description: "This is new project",
				UserId:      testutil.UserId(),
				CreatedAt:   testutil.Date(),
				UpdatedAt:   testutil.Date(),
			},
			expectedError: "failed to convert entry to entity (name): " +
				fmt.Sprintf(
					"project name cannot be longer than 100 characters, but got '%s'",
					tooLongProjectName,
				),
		},
		{
			name: "should return error when project description is too long",
			createdProject: record.ProjectEntry{
				Name:        "New Project",
				Description: tooLongProjectDescription,
				UserId:      testutil.UserId(),
				CreatedAt:   testutil.Date(),
				UpdatedAt:   testutil.Date(),
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
				InsertProject(record.ProjectWithoutAutofieldEntry{
					Name:        "New Project",
					Description: "This is new project",
					UserId:      testutil.UserId(),
				}).
				Return("0000000000000001", &tc.createdProject, nil)

			s := service.NewProjectService(r)

			userId, err := domain.NewUserIdObject(testutil.UserId())
			assert.NoError(t, err)

			name, err := domain.NewProjectNameObject("New Project")
			assert.NoError(t, err)
			description, err := domain.NewProjectDescriptionObject("This is new project")
			assert.NoError(t, err)

			project := domain.NewProjectWithoutAutofieldEntity(*name, *description)

			createdProject, err := s.CreateProject(*userId, *project)
			assert.ErrorContains(t, err, tc.expectedError)
			assert.Nil(t, createdProject)
		})
	}
}

func TestCreateProjectRepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := repository.NewMockProjectRepository(ctrl)
	r.EXPECT().
		InsertProject(gomock.Any()).
		Return("", nil, fmt.Errorf("repository error"))

	s := service.NewProjectService(r)

	userId, err := domain.NewUserIdObject(testutil.UserId())
	assert.NoError(t, err)

	name, err := domain.NewProjectNameObject("New Project")
	assert.NoError(t, err)
	description, err := domain.NewProjectDescriptionObject("This is new project")
	assert.NoError(t, err)

	project := domain.NewProjectWithoutAutofieldEntity(*name, *description)

	createdProject, err := s.CreateProject(*userId, *project)
	assert.ErrorContains(t, err, "failed to insert project: repository error")
	assert.Nil(t, createdProject)
}
