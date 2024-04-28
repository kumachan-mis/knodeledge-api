package repository_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/kumachan-mis/knodeledge-api/internal/db"
	"github.com/kumachan-mis/knodeledge-api/internal/repository"
	"github.com/kumachan-mis/knodeledge-api/test/testutil"
	"github.com/stretchr/testify/assert"
)

func TestFetchProjectChaptersValidDocument(t *testing.T) {
	client := db.FirestoreClient()
	r := repository.NewChapterRepository(*client)

	chapters, rErr := r.FetchProjectChapters(testutil.ReadOnlyUserId(), "PROJECT_WITHOUT_DESCRIPTION")

	assert.Nil(t, rErr)

	assert.Len(t, chapters, 2)

	chapter := chapters["CHAPTER_ONE"]
	assert.Equal(t, "Chapter One", chapter.Name)
	assert.Equal(t, 1, chapter.Number)
	assert.Equal(t, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), chapter.CreatedAt)
	assert.Equal(t, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), chapter.UpdatedAt)

	chapter = chapters["CHAPTER_TWO"]
	assert.Equal(t, "Chapter Two", chapter.Name)
	assert.Equal(t, 2, chapter.Number)
	assert.Equal(t, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), chapter.CreatedAt)
	assert.Equal(t, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), chapter.UpdatedAt)
}

func TestFetchProjectChaptersNoDocument(t *testing.T) {
	client := db.FirestoreClient()
	r := repository.NewChapterRepository(*client)

	chapters, rErr := r.FetchProjectChapters(testutil.ReadOnlyUserId(), "PROJECT_WITH_DESCRIPTION")

	assert.Nil(t, rErr)

	assert.Empty(t, chapters)
}

func TestFetchProjectChaptersNoProject(t *testing.T) {
	client := db.FirestoreClient()
	r := repository.NewChapterRepository(*client)

	chapters, rErr := r.FetchProjectChapters(testutil.ReadOnlyUserId(), "UNKNOWN_PROJECT")

	assert.NotNil(t, rErr)
	assert.Equal(t, repository.NotFoundError, rErr.Code())
	assert.Equal(t, "not found: parent document not found", rErr.Error())
	assert.Nil(t, chapters)
}

func TestFetchProjectChaptersUnauthorizedProject(t *testing.T) {
	client := db.FirestoreClient()
	r := repository.NewChapterRepository(*client)

	chapters, rErr := r.FetchProjectChapters(testutil.ModifyOnlyUserId(), "PROJECT_WITHOUT_DESCRIPTION")

	assert.NotNil(t, rErr)
	assert.Equal(t, repository.NotFoundError, rErr.Code())
	assert.Equal(t, "not found: parent document not found", rErr.Error())
	assert.Nil(t, chapters)
}

func TestFetchProjectChaptersInvalidDocument(t *testing.T) {
	tt := []struct {
		name          string
		userId        string
		projectId     string
		expectedError string
	}{
		{
			name:      "should return error when chapter name is invalid",
			userId:    testutil.ErrorUserId(2),
			projectId: "PROJECT_WITH_INVALID_CHAPTER_NAME",
			expectedError: "failed to convert snapshot to values: document.ChapterValues.name: " +
				"firestore: cannot set type string to bool",
		},
		{
			name:      "should return error when chapter number is invalid",
			userId:    testutil.ErrorUserId(3),
			projectId: "PROJECT_WITH_INVALID_CHAPTER_NUMBER",
			expectedError: "failed to convert snapshot to values: document.ChapterValues.number: " +
				"firestore: cannot set type int to string",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			client := db.FirestoreClient()
			r := repository.NewChapterRepository(*client)

			chapters, rErr := r.FetchProjectChapters(tc.userId, tc.projectId)

			assert.NotNil(t, rErr)
			assert.Equal(t, repository.ReadFailurePanic, rErr.Code())
			assert.Equal(t, fmt.Sprintf("read failure: %s", tc.expectedError), rErr.Error())
			assert.Nil(t, chapters)
		})
	}
}
