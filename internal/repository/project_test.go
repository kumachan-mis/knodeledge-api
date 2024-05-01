package repository_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/kumachan-mis/knodeledge-api/internal/db"
	"github.com/kumachan-mis/knodeledge-api/internal/record"
	"github.com/kumachan-mis/knodeledge-api/internal/repository"
	"github.com/kumachan-mis/knodeledge-api/test/testutil"
	"github.com/stretchr/testify/assert"
)

func TestFetchProjectsValidDocument(t *testing.T) {
	client := db.FirestoreClient()
	r := repository.NewProjectRepository(*client)

	userId := testutil.ReadOnlyUserId()
	projects, rErr := r.FetchProjects(userId)

	assert.Nil(t, rErr)

	assert.Len(t, projects, 2)

	project := projects["PROJECT_WITHOUT_DESCRIPTION"]
	assert.Equal(t, "No Description Project", project.Name)
	assert.Equal(t, "", project.Description)
	assert.Equal(t, userId, project.UserId)
	assert.Equal(t, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), project.CreatedAt)
	assert.Equal(t, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), project.UpdatedAt)

	project = projects["PROJECT_WITH_DESCRIPTION"]
	assert.Equal(t, "Described Project", project.Name)
	assert.Equal(t, "This is project description", project.Description)
	assert.Equal(t, userId, project.UserId)
	assert.Equal(t, time.Date(2023, 12, 31, 23, 0, 0, 0, time.UTC), project.CreatedAt)
	assert.Equal(t, time.Date(2023, 12, 31, 23, 0, 0, 0, time.UTC), project.UpdatedAt)
}

func TestFetchProjectsNoDocument(t *testing.T) {
	client := db.FirestoreClient()
	r := repository.NewProjectRepository(*client)

	userId := testutil.UnknownUserId()
	projects, rErr := r.FetchProjects(userId)

	assert.Nil(t, rErr)

	assert.Empty(t, projects)
}

func TestFetchProjectsInvalidDocument(t *testing.T) {
	tt := []struct {
		name          string
		userId        string
		expectedError string
	}{
		{
			name:   "should return error when project name is invalid",
			userId: testutil.ErrorUserId(0),
			expectedError: "failed to convert snapshot to values: document.ProjectValues.name: " +
				"firestore: cannot set type string to int",
		},
		{
			name:   "should return error when project description is invalid",
			userId: testutil.ErrorUserId(1),
			expectedError: "failed to convert snapshot to values: document.ProjectValues.description: " +
				"firestore: cannot set type string to bool",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			client := db.FirestoreClient()
			r := repository.NewProjectRepository(*client)

			projects, rErr := r.FetchProjects(tc.userId)

			assert.NotNil(t, rErr)
			assert.Equal(t, repository.ReadFailurePanic, rErr.Code())
			assert.Equal(t, fmt.Sprintf("read failure: %s", tc.expectedError), rErr.Error())
			assert.Nil(t, projects)
		})
	}
}

func TestFetchProjectValidDocument(t *testing.T) {
	tt := []struct {
		name      string
		projectId string
		expected  record.ProjectEntry
	}{
		{
			name:      "should return project without description",
			projectId: "PROJECT_WITHOUT_DESCRIPTION",
			expected: record.ProjectEntry{
				Name:        "No Description Project",
				Description: "",
				UserId:      testutil.ReadOnlyUserId(),
				CreatedAt:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name:      "should return project with description",
			projectId: "PROJECT_WITH_DESCRIPTION",
			expected: record.ProjectEntry{
				Name:        "Described Project",
				Description: "This is project description",
				UserId:      testutil.ReadOnlyUserId(),
				CreatedAt:   time.Date(2023, 12, 31, 23, 0, 0, 0, time.UTC),
				UpdatedAt:   time.Date(2023, 12, 31, 23, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			client := db.FirestoreClient()
			r := repository.NewProjectRepository(*client)

			project, rErr := r.FetchProject(testutil.ReadOnlyUserId(), tc.projectId)

			assert.Nil(t, rErr)

			assert.Equal(t, tc.expected.Name, project.Name)
			assert.Equal(t, tc.expected.Description, project.Description)
			assert.Equal(t, tc.expected.UserId, project.UserId)
			assert.Equal(t, tc.expected.CreatedAt, project.CreatedAt)
			assert.Equal(t, tc.expected.UpdatedAt, project.UpdatedAt)
		})
	}
}

func TestFetchProjectNotFound(t *testing.T) {
	tt := []struct {
		name          string
		userId        string
		projectId     string
		expectedError string
	}{
		{
			name:          "should return error when project is not found",
			userId:        testutil.ReadOnlyUserId(),
			projectId:     "UNKNOWN_PROJECT",
			expectedError: "failed to get project",
		},
		{
			name:          "should return error when user is not author of the project",
			userId:        testutil.ModifyOnlyUserId(),
			projectId:     "PROJECT_WITHOUT_DESCRIPTION",
			expectedError: "failed to get project",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			client := db.FirestoreClient()
			r := repository.NewProjectRepository(*client)

			project, rErr := r.FetchProject(tc.userId, tc.projectId)

			assert.NotNil(t, rErr)
			assert.Equal(t, repository.NotFoundError, rErr.Code())
			assert.Equal(t, fmt.Sprintf("not found: %s", tc.expectedError), rErr.Error())
			assert.Nil(t, project)
		})
	}
}

func TestFetchProjectInvalidDocument(t *testing.T) {
	tt := []struct {
		name          string
		userId        string
		projectId     string
		expectedError string
	}{
		{
			name:      "should return error when project name is invalid",
			userId:    testutil.ErrorUserId(0),
			projectId: "PROJECT_WITH_NAME_ERROR",
			expectedError: "failed to convert snapshot to values: document.ProjectValues.name: " +
				"firestore: cannot set type string to int",
		},
		{
			name:      "should return error when project description is invalid",
			userId:    testutil.ErrorUserId(1),
			projectId: "PROJECT_WITH_DESCRIPTION_ERROR",
			expectedError: "failed to convert snapshot to values: document.ProjectValues.description: " +
				"firestore: cannot set type string to bool",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			client := db.FirestoreClient()
			r := repository.NewProjectRepository(*client)

			project, rErr := r.FetchProject(tc.userId, tc.projectId)

			assert.NotNil(t, rErr)
			assert.Equal(t, repository.ReadFailurePanic, rErr.Code())
			assert.Equal(t, fmt.Sprintf("read failure: %s", tc.expectedError), rErr.Error())
			assert.Nil(t, project)
		})
	}
}

func TestInsertProjectValidEntry(t *testing.T) {
	tt := []struct {
		name    string
		project record.ProjectWithoutAutofieldEntry
	}{
		{
			name: "should insert project without description",
			project: record.ProjectWithoutAutofieldEntry{
				Name:   "New Project",
				UserId: testutil.ModifyOnlyUserId(),
			},
		},
		{
			name: "should insert project with description",
			project: record.ProjectWithoutAutofieldEntry{
				Name:        "New Project",
				Description: "This is new project",
				UserId:      testutil.ModifyOnlyUserId(),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			client := db.FirestoreClient()
			r := repository.NewProjectRepository(*client)

			id, createdProject, rErr := r.InsertProject(tc.project)
			now := time.Now()

			assert.Nil(t, rErr)

			assert.NotEmpty(t, id)
			assert.Equal(t, tc.project.Name, createdProject.Name)
			assert.Equal(t, tc.project.Description, createdProject.Description)
			assert.Equal(t, tc.project.UserId, createdProject.UserId)
			assert.Less(t, now.Sub(createdProject.CreatedAt), time.Second)
			assert.Less(t, now.Sub(createdProject.UpdatedAt), time.Second)
		})
	}
}

func TestUpdateProjectValidEntry(t *testing.T) {
	tt := []struct {
		name      string
		projectId string
		project   record.ProjectWithoutAutofieldEntry
	}{
		{
			name:      "should update project without description",
			projectId: "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_REPOSITORY",
			project: record.ProjectWithoutAutofieldEntry{
				Name:   "Updated Project",
				UserId: testutil.ModifyOnlyUserId(),
			},
		},
		{
			name:      "should update project with description",
			projectId: "PROJECT_WITH_DESCRIPTION_TO_UPDATE_FROM_REPOSITORY",
			project: record.ProjectWithoutAutofieldEntry{
				Name:        "Updated Project",
				Description: "This is updated project",
				UserId:      testutil.ModifyOnlyUserId(),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			client := db.FirestoreClient()
			r := repository.NewProjectRepository(*client)

			updatedProject, rErr := r.UpdateProject(tc.projectId, tc.project)
			now := time.Now()

			assert.Nil(t, rErr)

			assert.Equal(t, tc.project.Name, updatedProject.Name)
			assert.Equal(t, tc.project.Description, updatedProject.Description)
			assert.Equal(t, tc.project.UserId, updatedProject.UserId)
			assert.Equal(t, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), updatedProject.CreatedAt)
			assert.Less(t, now.Sub(updatedProject.UpdatedAt), time.Second)
		})
	}
}

func TestUpdateProjectNotFound(t *testing.T) {
	tt := []struct {
		name          string
		userId        string
		projectId     string
		expectedError string
	}{
		{
			name:          "should return error when project is not found",
			userId:        testutil.ReadOnlyUserId(),
			projectId:     "UNKNOWN_PROJECT",
			expectedError: "failed to update project",
		},
		{
			name:          "should return error when user is not author of the project",
			userId:        testutil.ModifyOnlyUserId(),
			projectId:     "PROJECT_WITHOUT_DESCRIPTION",
			expectedError: "failed to update project",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			client := db.FirestoreClient()
			r := repository.NewProjectRepository(*client)

			project, rErr := r.UpdateProject(tc.projectId, record.ProjectWithoutAutofieldEntry{
				Name:        "Updated Project",
				Description: "This is updated project",
				UserId:      tc.userId,
			})

			assert.NotNil(t, rErr)
			assert.Equal(t, repository.NotFoundError, rErr.Code())
			assert.Equal(t, fmt.Sprintf("not found: %s", tc.expectedError), rErr.Error())
			assert.Nil(t, project)
		})
	}
}

func TestUpdateProjectInvalidDocument(t *testing.T) {
	tt := []struct {
		name          string
		userId        string
		projectId     string
		expectedError string
	}{
		{
			name:      "should return error when project name is invalid",
			userId:    testutil.ErrorUserId(0),
			projectId: "PROJECT_WITH_NAME_ERROR",
			expectedError: "failed to convert snapshot to values: document.ProjectValues.name: " +
				"firestore: cannot set type string to int",
		},
		{
			name:      "should return error when project description is invalid",
			userId:    testutil.ErrorUserId(1),
			projectId: "PROJECT_WITH_DESCRIPTION_ERROR",
			expectedError: "failed to convert snapshot to values: document.ProjectValues.description: " +
				"firestore: cannot set type string to bool",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			client := db.FirestoreClient()
			r := repository.NewProjectRepository(*client)

			project, rErr := r.UpdateProject(tc.projectId, record.ProjectWithoutAutofieldEntry{
				Name:        "Updated Project",
				Description: "This is updated project",
				UserId:      tc.userId,
			})

			assert.NotNil(t, rErr)
			assert.Equal(t, repository.ReadFailurePanic, rErr.Code())
			assert.Equal(t, fmt.Sprintf("read failure: %s", tc.expectedError), rErr.Error())
			assert.Nil(t, project)
		})
	}
}
