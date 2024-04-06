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

func TestFetchUserProjectsValidDocument(t *testing.T) {
	client := db.FirestoreClient()
	r := repository.NewProjectRepository(*client)

	userId := testutil.ReadOnlyUserId()
	projects, rErr := r.FetchUserProjects(userId)

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

func TestFetchUserProjectsNoDocument(t *testing.T) {
	client := db.FirestoreClient()
	r := repository.NewProjectRepository(*client)

	userId := testutil.UnknownUserId()
	projects, rErr := r.FetchUserProjects(userId)

	assert.Nil(t, rErr)

	assert.Empty(t, projects)
}

func TestFetchUserProjectsInvalidDocument(t *testing.T) {
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

			projects, rErr := r.FetchUserProjects(tc.userId)

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

			project, rErr := r.FetchProject(tc.projectId)

			assert.Nil(t, rErr)

			assert.Equal(t, tc.expected.Name, project.Name)
			assert.Equal(t, tc.expected.Description, project.Description)
			assert.Equal(t, tc.expected.UserId, project.UserId)
			assert.Equal(t, tc.expected.CreatedAt, project.CreatedAt)
			assert.Equal(t, tc.expected.UpdatedAt, project.UpdatedAt)
		})
	}
}

func TestFetchProjectNoDocument(t *testing.T) {
	client := db.FirestoreClient()
	r := repository.NewProjectRepository(*client)

	projectId := "UNKNOWN_PROJECT"
	project, rErr := r.FetchProject(projectId)

	assert.NotNil(t, rErr)
	assert.Equal(t, repository.NotFoundError, rErr.Code())
	assert.Equal(t, "not found: failed to get project", rErr.Error())
	assert.Nil(t, project)
}

func TestFetchProjectInvalidDocument(t *testing.T) {
	tt := []struct {
		name          string
		projectId     string
		expectedError string
	}{
		{
			name:      "should return error when project name is invalid",
			projectId: "PROJECT_WITH_NAME_ERROR",
			expectedError: "failed to convert snapshot to values: document.ProjectValues.name: " +
				"firestore: cannot set type string to int",
		},
		{
			name:      "should return error when project description is invalid",
			projectId: "PROJECT_WITH_DESCRIPTION_ERROR",
			expectedError: "failed to convert snapshot to values: document.ProjectValues.description: " +
				"firestore: cannot set type string to bool",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			client := db.FirestoreClient()
			r := repository.NewProjectRepository(*client)

			project, rErr := r.FetchProject(tc.projectId)

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
		name    string
		id      string
		project record.ProjectWithoutAutofieldEntry
	}{
		{
			name: "should update project without description",
			id:   "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_REPOSITORY",
			project: record.ProjectWithoutAutofieldEntry{
				Name:   "Updated Project",
				UserId: testutil.ModifyOnlyUserId(),
			},
		},
		{
			name: "should update project with description",
			id:   "PROJECT_WITH_DESCRIPTION_TO_UPDATE_FROM_REPOSITORY",
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

			updatedProject, rErr := r.UpdateProject(tc.id, tc.project)
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

func TestUpdateProjectNoDocument(t *testing.T) {
	client := db.FirestoreClient()
	r := repository.NewProjectRepository(*client)

	projectId := "UNKNOWN_PROJECT"
	entry := record.ProjectWithoutAutofieldEntry{
		Name:        "Updated Project",
		Description: "This is updated project",
		UserId:      testutil.ModifyOnlyUserId(),
	}

	project, rErr := r.UpdateProject(projectId, entry)

	assert.NotNil(t, rErr)
	assert.Equal(t, repository.WriteFailurePanic, rErr.Code())
	assert.Equal(t, "write failure: failed to update project", rErr.Error())
	assert.Nil(t, project)
}
