package service_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/kumachan-mis/knodeledge-api/internal/domain"
	"github.com/kumachan-mis/knodeledge-api/internal/record"
	"github.com/kumachan-mis/knodeledge-api/internal/repository"
	"github.com/kumachan-mis/knodeledge-api/internal/service"
	mock_repository "github.com/kumachan-mis/knodeledge-api/mock/repository"
	"github.com/kumachan-mis/knodeledge-api/test/testutil"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestListProjectsValidEntry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	maxLengthProjectName := testutil.RandomString(100)
	maxLengthProjectDescription := testutil.RandomString(400)

	r := mock_repository.NewMockProjectRepository(ctrl)
	r.EXPECT().
		FetchUserProjects(testutil.ReadOnlyUserId()).
		Return(map[string]record.ProjectEntry{
			"0000000000000003": {
				Name:        maxLengthProjectName,
				Description: maxLengthProjectDescription,
				UserId:      testutil.ReadOnlyUserId(),
				CreatedAt:   testutil.Date().Add(-3 * time.Hour),
				UpdatedAt:   testutil.Date().Add(-3 * time.Hour),
			},
			"0000000000000002": {
				Name:      "Second Project",
				UserId:    testutil.ReadOnlyUserId(),
				CreatedAt: testutil.Date().Add(-2 * time.Hour),
				UpdatedAt: testutil.Date().Add(-2 * time.Hour),
			},
			"0000000000000001": {
				Name:        "First Project",
				Description: "This is my first project",
				UserId:      testutil.ReadOnlyUserId(),
				CreatedAt:   testutil.Date().Add(-1 * time.Hour),
				UpdatedAt:   testutil.Date().Add(-1 * time.Hour),
			},
		}, nil)

	s := service.NewProjectService(r)

	userId, err := domain.NewUserIdObject(testutil.ReadOnlyUserId())
	assert.NoError(t, err)

	projects, sErr := s.ListProjects(*userId)
	assert.Nil(t, sErr)

	assert.Len(t, projects, 3)

	project := projects[0]
	assert.Equal(t, "0000000000000001", project.Id().Value())
	assert.Equal(t, "First Project", project.Name().Value())
	assert.Equal(t, "This is my first project", project.Description().Value())
	assert.Equal(t, testutil.Date().Add(-1*time.Hour), project.CreatedAt().Value())
	assert.Equal(t, testutil.Date().Add(-1*time.Hour), project.UpdatedAt().Value())
	assert.True(t, project.AuthoredBy(userId))

	project = projects[1]
	assert.Equal(t, "0000000000000002", project.Id().Value())
	assert.Equal(t, "Second Project", project.Name().Value())
	assert.Equal(t, "", project.Description().Value())
	assert.Equal(t, testutil.Date().Add(-2*time.Hour), project.CreatedAt().Value())
	assert.Equal(t, testutil.Date().Add(-2*time.Hour), project.UpdatedAt().Value())
	assert.True(t, project.AuthoredBy(userId))

	project = projects[2]
	assert.Equal(t, "0000000000000003", project.Id().Value())
	assert.Equal(t, maxLengthProjectName, project.Name().Value())
	assert.Equal(t, maxLengthProjectDescription, project.Description().Value())
	assert.Equal(t, testutil.Date().Add(-3*time.Hour), project.CreatedAt().Value())
	assert.Equal(t, testutil.Date().Add(-3*time.Hour), project.UpdatedAt().Value())
	assert.True(t, project.AuthoredBy(userId))
}

func TestListProjectsNoEntry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := mock_repository.NewMockProjectRepository(ctrl)
	r.EXPECT().
		FetchUserProjects(testutil.ReadOnlyUserId()).
		Return(map[string]record.ProjectEntry{}, nil)

	s := service.NewProjectService(r)

	userId, err := domain.NewUserIdObject(testutil.ReadOnlyUserId())
	assert.NoError(t, err)

	projects, sErr := s.ListProjects(*userId)
	assert.Nil(t, sErr)

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
				UserId:    testutil.ReadOnlyUserId(),
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
				UserId:    testutil.ReadOnlyUserId(),
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
				UserId:    testutil.ReadOnlyUserId(),
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
				UserId:      testutil.ReadOnlyUserId(),
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

			r := mock_repository.NewMockProjectRepository(ctrl)
			r.EXPECT().
				FetchUserProjects(testutil.ReadOnlyUserId()).
				Return(map[string]record.ProjectEntry{
					tc.projectId: tc.project,
				}, nil)

			s := service.NewProjectService(r)

			userId, err := domain.NewUserIdObject(testutil.ReadOnlyUserId())
			assert.NoError(t, err)

			projects, sErr := s.ListProjects(*userId)
			assert.NotNil(t, sErr)
			assert.Equal(t, service.DomainFailurePanic, sErr.Code())
			assert.Equal(t, fmt.Sprintf("domain failure: %s", tc.expectedError), sErr.Error())
			assert.Nil(t, projects)
		})
	}
}

func TestListProjectsRepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := mock_repository.NewMockProjectRepository(ctrl)
	r.EXPECT().
		FetchUserProjects(testutil.ReadOnlyUserId()).
		Return(nil, repository.Errorf(repository.ReadFailurePanic, "repository error"))

	s := service.NewProjectService(r)

	userId, err := domain.NewUserIdObject(testutil.ReadOnlyUserId())
	assert.NoError(t, err)

	projects, err := s.ListProjects(*userId)
	assert.NotNil(t, err)
	assert.Equal(t, "repository failure: failed to fetch user projects: repository error", err.Error())
	assert.Nil(t, projects)
}

func TestFindProjectValidEntry(t *testing.T) {
	maxLengthProjectName := testutil.RandomString(100)
	maxLengthProjectDescription := testutil.RandomString(400)

	tt := []struct {
		name      string
		projectId string
		project   record.ProjectEntry
	}{
		{
			name:      "should return project with valid entry",
			projectId: "0000000000000001",
			project: record.ProjectEntry{
				Name:        "New Project",
				Description: "This is new project",
				UserId:      testutil.ReadOnlyUserId(),
				CreatedAt:   testutil.Date(),
				UpdatedAt:   testutil.Date(),
			},
		},
		{
			name:      "should return project with max-length valid entry",
			projectId: "0000000000000002",
			project: record.ProjectEntry{
				Name:        maxLengthProjectName,
				Description: maxLengthProjectDescription,
				UserId:      testutil.ReadOnlyUserId(),
				CreatedAt:   testutil.Date(),
				UpdatedAt:   testutil.Date(),
			},
		},
	}

	for _, tc := range tt {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		r := mock_repository.NewMockProjectRepository(ctrl)
		r.EXPECT().
			FetchProject(tc.projectId).
			Return(&tc.project, nil)

		s := service.NewProjectService(r)

		userId, err := domain.NewUserIdObject(tc.project.UserId)
		assert.NoError(t, err)

		projectId, err := domain.NewProjectIdObject(tc.projectId)
		assert.NoError(t, err)

		project, sErr := s.FindProject(*userId, *projectId)
		assert.Nil(t, sErr)

		assert.Equal(t, tc.projectId, project.Id().Value())
		assert.Equal(t, tc.project.Name, project.Name().Value())
		assert.Equal(t, tc.project.Description, project.Description().Value())
		assert.Equal(t, tc.project.CreatedAt, project.CreatedAt().Value())
		assert.Equal(t, tc.project.UpdatedAt, project.UpdatedAt().Value())
		assert.True(t, project.AuthoredBy(userId))
	}
}

func TestFindProjectNoEntry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := mock_repository.NewMockProjectRepository(ctrl)
	r.EXPECT().
		FetchProject("0000000000000001").
		Return(nil, repository.Errorf(repository.NotFoundError, "failed to get project"))

	s := service.NewProjectService(r)

	userId, err := domain.NewUserIdObject(testutil.ReadOnlyUserId())
	assert.NoError(t, err)

	projectId, err := domain.NewProjectIdObject("0000000000000001")
	assert.NoError(t, err)

	project, sErr := s.FindProject(*userId, *projectId)
	assert.NotNil(t, sErr)
	assert.Equal(t, service.NotFoundError, sErr.Code())
	assert.Equal(t, "not found: failed to find project", sErr.Error())
	assert.Nil(t, project)
}

func TestFindProjectUnauthoredEntry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := mock_repository.NewMockProjectRepository(ctrl)
	r.EXPECT().
		FetchProject("0000000000000001").
		Return(&record.ProjectEntry{
			Name:      "Project",
			UserId:    testutil.ModifyOnlyUserId(),
			CreatedAt: testutil.Date(),
			UpdatedAt: testutil.Date(),
		}, nil)

	s := service.NewProjectService(r)

	userId, err := domain.NewUserIdObject(testutil.ReadOnlyUserId())
	assert.NoError(t, err)

	projectId, err := domain.NewProjectIdObject("0000000000000001")
	assert.NoError(t, err)

	project, sErr := s.FindProject(*userId, *projectId)
	assert.NotNil(t, sErr)
	assert.Equal(t, service.NotFoundError, sErr.Code())
	assert.Equal(t, "not found: failed to find project", sErr.Error())
	assert.Nil(t, project)
}

func TestFindProjectInvalidEntry(t *testing.T) {
	tooLongProjectName := testutil.RandomString(101)
	tooLongProjectDescription := testutil.RandomString(401)

	tt := []struct {
		name          string
		projectId     string
		project       record.ProjectEntry
		expectedError string
	}{
		{
			name:      "should return error when project name is empty",
			projectId: "0000000000000001",
			project: record.ProjectEntry{
				Name:      "",
				UserId:    testutil.ReadOnlyUserId(),
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
				UserId:    testutil.ReadOnlyUserId(),
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
				UserId:      testutil.ReadOnlyUserId(),
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

			r := mock_repository.NewMockProjectRepository(ctrl)
			r.EXPECT().
				FetchProject(tc.projectId).
				Return(&tc.project, nil)

			s := service.NewProjectService(r)

			userId, err := domain.NewUserIdObject(tc.project.UserId)
			assert.NoError(t, err)

			projectId, err := domain.NewProjectIdObject(tc.projectId)
			assert.NoError(t, err)

			project, sErr := s.FindProject(*userId, *projectId)
			assert.NotNil(t, sErr)
			assert.Equal(t, service.DomainFailurePanic, sErr.Code())
			assert.Equal(t, fmt.Sprintf("domain failure: %s", tc.expectedError), sErr.Error())
			assert.Nil(t, project)
		})
	}
}

func TestFindProjectRepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := mock_repository.NewMockProjectRepository(ctrl)
	r.EXPECT().
		FetchProject("0000000000000001").
		Return(nil, repository.Errorf(repository.ReadFailurePanic, "repository error"))

	s := service.NewProjectService(r)

	userId, err := domain.NewUserIdObject(testutil.ReadOnlyUserId())
	assert.NoError(t, err)

	projectId, err := domain.NewProjectIdObject("0000000000000001")
	assert.NoError(t, err)

	project, sErr := s.FindProject(*userId, *projectId)
	assert.NotNil(t, sErr)
	assert.Equal(t, service.RepositoryFailurePanic, sErr.Code())
	assert.Equal(t, "repository failure: failed to fetch project: repository error", sErr.Error())
	assert.Nil(t, project)
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
				UserId:      testutil.ModifyOnlyUserId(),
			},
		},
		{
			name:      "should return project with max-length valid entry",
			projectId: "0000000000000002",
			project: record.ProjectWithoutAutofieldEntry{
				Name:        maxLengthProjectName,
				Description: maxLengthProjectDescription,
				UserId:      testutil.ModifyOnlyUserId(),
			},
		},
	}

	for _, tc := range tt {

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		r := mock_repository.NewMockProjectRepository(ctrl)
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

		userId, err := domain.NewUserIdObject(testutil.ModifyOnlyUserId())
		assert.NoError(t, err)

		name, err := domain.NewProjectNameObject(tc.project.Name)
		assert.NoError(t, err)
		description, err := domain.NewProjectDescriptionObject(tc.project.Description)
		assert.NoError(t, err)

		project := domain.NewProjectWithoutAutofieldEntity(*name, *description)

		createdProject, sErr := s.CreateProject(*userId, *project)
		assert.Nil(t, sErr)

		assert.Equal(t, tc.projectId, createdProject.Id().Value())
		assert.Equal(t, tc.project.Name, createdProject.Name().Value())
		assert.Equal(t, tc.project.Description, createdProject.Description().Value())
		assert.Equal(t, testutil.Date(), createdProject.CreatedAt().Value())
		assert.Equal(t, testutil.Date(), createdProject.UpdatedAt().Value())
		assert.True(t, createdProject.AuthoredBy(userId))
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
				UserId:      testutil.ModifyOnlyUserId(),
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
				UserId:      testutil.ModifyOnlyUserId(),
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
				UserId:      testutil.ModifyOnlyUserId(),
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

			r := mock_repository.NewMockProjectRepository(ctrl)
			r.EXPECT().
				InsertProject(record.ProjectWithoutAutofieldEntry{
					Name:        "New Project",
					Description: "This is new project",
					UserId:      testutil.ModifyOnlyUserId(),
				}).
				Return("0000000000000001", &tc.createdProject, nil)

			s := service.NewProjectService(r)

			userId, err := domain.NewUserIdObject(testutil.ModifyOnlyUserId())
			assert.NoError(t, err)

			name, err := domain.NewProjectNameObject("New Project")
			assert.NoError(t, err)
			description, err := domain.NewProjectDescriptionObject("This is new project")
			assert.NoError(t, err)

			project := domain.NewProjectWithoutAutofieldEntity(*name, *description)

			createdProject, sErr := s.CreateProject(*userId, *project)
			assert.NotNil(t, sErr)
			assert.Equal(t, service.DomainFailurePanic, sErr.Code())
			assert.Equal(t, fmt.Sprintf("domain failure: %s", tc.expectedError), sErr.Error())
			assert.Nil(t, createdProject)
		})
	}
}

func TestCreateProjectRepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := mock_repository.NewMockProjectRepository(ctrl)
	r.EXPECT().
		InsertProject(gomock.Any()).
		Return("", nil, repository.Errorf(repository.WriteFailurePanic, "repository error"))

	s := service.NewProjectService(r)

	userId, err := domain.NewUserIdObject(testutil.ModifyOnlyUserId())
	assert.NoError(t, err)

	name, err := domain.NewProjectNameObject("New Project")
	assert.NoError(t, err)
	description, err := domain.NewProjectDescriptionObject("This is new project")
	assert.NoError(t, err)

	project := domain.NewProjectWithoutAutofieldEntity(*name, *description)

	createdProject, sErr := s.CreateProject(*userId, *project)
	assert.NotNil(t, sErr)
	assert.Equal(t, service.RepositoryFailurePanic, sErr.Code())
	assert.Equal(t, "repository failure: failed to insert project: repository error", sErr.Error())
	assert.Nil(t, createdProject)
}

func TestEditProjectValidEntry(t *testing.T) {
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
				Name:        "Updated Project",
				Description: "This is updated project",
				UserId:      testutil.ModifyOnlyUserId(),
			},
		},
		{
			name:      "should return project with max-length valid entry",
			projectId: "0000000000000002",
			project: record.ProjectWithoutAutofieldEntry{
				Name:        maxLengthProjectName,
				Description: maxLengthProjectDescription,
				UserId:      testutil.ModifyOnlyUserId(),
			},
		},
	}

	for _, tc := range tt {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		r := mock_repository.NewMockProjectRepository(ctrl)
		r.EXPECT().
			UpdateProject(tc.projectId, tc.project).
			Return(&record.ProjectEntry{
				Name:        tc.project.Name,
				Description: tc.project.Description,
				UserId:      tc.project.UserId,
				CreatedAt:   testutil.Date(),
				UpdatedAt:   testutil.Date(),
			}, nil)

		s := service.NewProjectService(r)

		userId, err := domain.NewUserIdObject(testutil.ModifyOnlyUserId())
		assert.NoError(t, err)

		projectId, err := domain.NewProjectIdObject(tc.projectId)
		assert.NoError(t, err)

		name, err := domain.NewProjectNameObject(tc.project.Name)
		assert.NoError(t, err)
		description, err := domain.NewProjectDescriptionObject(tc.project.Description)
		assert.NoError(t, err)

		project := domain.NewProjectWithoutAutofieldEntity(*name, *description)

		editedProject, sErr := s.UpdateProject(*userId, *projectId, *project)
		assert.Nil(t, sErr)

		assert.Equal(t, tc.projectId, editedProject.Id().Value())
		assert.Equal(t, tc.project.Name, editedProject.Name().Value())
		assert.Equal(t, tc.project.Description, editedProject.Description().Value())
		assert.Equal(t, testutil.Date(), editedProject.CreatedAt().Value())
		assert.Equal(t, testutil.Date(), editedProject.UpdatedAt().Value())
		assert.True(t, editedProject.AuthoredBy(userId))
	}
}

func TestEditProjectNoEntry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := mock_repository.NewMockProjectRepository(ctrl)
	r.EXPECT().
		UpdateProject("0000000000000001", record.ProjectWithoutAutofieldEntry{
			Name:        "Updated Project",
			Description: "This is updated project",
			UserId:      testutil.ModifyOnlyUserId(),
		}).
		Return(nil, repository.Errorf(repository.NotFoundError, "failed to update project"))

	s := service.NewProjectService(r)

	userId, err := domain.NewUserIdObject(testutil.ModifyOnlyUserId())
	assert.NoError(t, err)

	projectId, err := domain.NewProjectIdObject("0000000000000001")
	assert.NoError(t, err)

	name, err := domain.NewProjectNameObject("Updated Project")
	assert.NoError(t, err)
	description, err := domain.NewProjectDescriptionObject("This is updated project")
	assert.NoError(t, err)

	project := domain.NewProjectWithoutAutofieldEntity(*name, *description)

	editedProject, sErr := s.UpdateProject(*userId, *projectId, *project)
	assert.NotNil(t, sErr)
	assert.Equal(t, service.NotFoundError, sErr.Code())
	assert.Equal(t, "not found: failed to update project", sErr.Error())
	assert.Nil(t, editedProject)
}

func TestEditProjectUnauthoredEntry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := mock_repository.NewMockProjectRepository(ctrl)
	r.EXPECT().
		UpdateProject("0000000000000001", record.ProjectWithoutAutofieldEntry{
			Name:        "Updated Project",
			Description: "This is updated project",
			UserId:      testutil.ModifyOnlyUserId(),
		}).
		Return(&record.ProjectEntry{
			Name:      "Updated Project",
			UserId:    testutil.ReadOnlyUserId(),
			CreatedAt: testutil.Date(),
			UpdatedAt: testutil.Date(),
		}, nil)

	s := service.NewProjectService(r)

	userId, err := domain.NewUserIdObject(testutil.ModifyOnlyUserId())
	assert.NoError(t, err)

	projectId, err := domain.NewProjectIdObject("0000000000000001")
	assert.NoError(t, err)

	name, err := domain.NewProjectNameObject("Updated Project")
	assert.NoError(t, err)
	description, err := domain.NewProjectDescriptionObject("This is updated project")
	assert.NoError(t, err)

	project := domain.NewProjectWithoutAutofieldEntity(*name, *description)

	editedProject, sErr := s.UpdateProject(*userId, *projectId, *project)
	assert.NotNil(t, sErr)
	assert.Equal(t, service.NotFoundError, sErr.Code())
	assert.Equal(t, "not found: failed to update project", sErr.Error())
	assert.Nil(t, editedProject)
}

func TestEditProjectInvalidEditedEntry(t *testing.T) {
	tooLongProjectName := testutil.RandomString(101)
	tooLongProjectDescription := testutil.RandomString(401)

	tt := []struct {
		name          string
		editedProject record.ProjectEntry
		expectedError string
	}{
		{
			name: "should return error when project name is empty",
			editedProject: record.ProjectEntry{
				Name:        "",
				Description: "This is updated project",
				UserId:      testutil.ModifyOnlyUserId(),
				CreatedAt:   testutil.Date(),
				UpdatedAt:   testutil.Date(),
			},
			expectedError: "failed to convert entry to entity (name): project name is required, but got ''",
		},
		{
			name: "should return error when project name is too long",
			editedProject: record.ProjectEntry{
				Name:        tooLongProjectName,
				Description: "This is updated project",
				UserId:      testutil.ModifyOnlyUserId(),
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
			editedProject: record.ProjectEntry{
				Name:        "Updated Project",
				Description: tooLongProjectDescription,
				UserId:      testutil.ModifyOnlyUserId(),
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

			r := mock_repository.NewMockProjectRepository(ctrl)
			r.EXPECT().
				UpdateProject("0000000000000001", record.ProjectWithoutAutofieldEntry{
					Name:        "Updated Project",
					Description: "This is updated project",
					UserId:      testutil.ModifyOnlyUserId(),
				}).
				Return(&tc.editedProject, nil)

			s := service.NewProjectService(r)

			userId, err := domain.NewUserIdObject(testutil.ModifyOnlyUserId())
			assert.NoError(t, err)

			projectId, err := domain.NewProjectIdObject("0000000000000001")
			assert.NoError(t, err)

			name, err := domain.NewProjectNameObject("Updated Project")
			assert.NoError(t, err)
			description, err := domain.NewProjectDescriptionObject("This is updated project")
			assert.NoError(t, err)

			project := domain.NewProjectWithoutAutofieldEntity(*name, *description)

			editedProject, sErr := s.UpdateProject(*userId, *projectId, *project)
			assert.NotNil(t, sErr)
			assert.Equal(t, service.DomainFailurePanic, sErr.Code())
			assert.Equal(t, fmt.Sprintf("domain failure: %s", tc.expectedError), sErr.Error())
			assert.Nil(t, editedProject)
		})
	}
}

func TestEditProjectRepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := mock_repository.NewMockProjectRepository(ctrl)
	r.EXPECT().
		UpdateProject("0000000000000001", record.ProjectWithoutAutofieldEntry{
			Name:        "Updated Project",
			Description: "This is updated project",
			UserId:      testutil.ModifyOnlyUserId(),
		}).
		Return(nil, repository.Errorf(repository.ReadFailurePanic, "repository error"))

	s := service.NewProjectService(r)

	userId, err := domain.NewUserIdObject(testutil.ModifyOnlyUserId())
	assert.NoError(t, err)

	projectId, err := domain.NewProjectIdObject("0000000000000001")
	assert.NoError(t, err)

	name, err := domain.NewProjectNameObject("Updated Project")
	assert.NoError(t, err)
	description, err := domain.NewProjectDescriptionObject("This is updated project")
	assert.NoError(t, err)

	project := domain.NewProjectWithoutAutofieldEntity(*name, *description)

	editedProject, sErr := s.UpdateProject(*userId, *projectId, *project)
	assert.NotNil(t, sErr)
	assert.Equal(t, service.RepositoryFailurePanic, sErr.Code())
	assert.Equal(t, "repository failure: failed to update project: repository error", sErr.Error())
	assert.Nil(t, editedProject)
}
