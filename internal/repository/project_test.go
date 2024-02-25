package repository_test

import (
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

	userId := testutil.UserId()
	projects, err := r.FetchUserProjects(userId)

	assert.NoError(t, err)

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
	projects, err := r.FetchUserProjects(userId)

	assert.NoError(t, err)

	assert.Empty(t, projects)
}

func TestFetchUserProjectsInvalidDocument(t *testing.T) {
	tt := []struct {
		name   string
		userId string
	}{
		{
			name:   "should return error when project name is invalid",
			userId: testutil.ErrorUserId(0),
		},
		{
			name:   "should return error when project description is invalid",
			userId: testutil.ErrorUserId(1),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			client := db.FirestoreClient()
			r := repository.NewProjectRepository(*client)

			projects, err := r.FetchUserProjects(tc.userId)

			assert.ErrorContains(t, err, "failed to convert snapshot to entry:")
			assert.Nil(t, projects)
		})
	}
}

func TestInsertProjectValidEntry(t *testing.T) {
	client := db.FirestoreClient()
	r := repository.NewProjectRepository(*client)

	userId := testutil.UserId()
	entry := record.ProjectWithoutAutofieldEntry{
		Name:        "New Project",
		Description: "This is new project",
		UserId:      userId,
	}

	id, project, err := r.InsertProject(entry)
	now := time.Now()

	assert.NoError(t, err)

	assert.NotEmpty(t, id)
	assert.Equal(t, entry.Name, project.Name)
	assert.Equal(t, entry.Description, project.Description)
	assert.Equal(t, entry.UserId, project.UserId)
	assert.Less(t, now.Sub(project.CreatedAt), time.Second)
	assert.Less(t, now.Sub(project.UpdatedAt), time.Second)
}
