package repository_test

import (
	"testing"

	"github.com/kumachan-mis/knodeledge-api/internal/db"
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

	project = projects["PROJECT_WITH_DESCRIPTION"]
	assert.Equal(t, "Described Project", project.Name)
	assert.Equal(t, "This is project description", project.Description)
	assert.Equal(t, userId, project.UserId)
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
