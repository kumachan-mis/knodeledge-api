package repository_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/kumachan-mis/knodeledge-api/internal/db"
	"github.com/kumachan-mis/knodeledge-api/internal/record"
	"github.com/kumachan-mis/knodeledge-api/internal/repository"
	"github.com/kumachan-mis/knodeledge-api/test/testutil"
	"github.com/stretchr/testify/assert"
)

func TestInsertPaperValidEntry(t *testing.T) {
	client := db.FirestoreClient()
	r := repository.NewPaperRepository(*client)

	projectId := "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_REPOSITORY"
	chapterId := "CHAPTER_ONE"
	content := strings.Join([]string{
		"## Introduction",
		"This is the introduction of the paper.",
		"",
		"## What is note apps?",
		"Note apps is a web application that allows users to create, read, update, and delete notes.",
		"",
	}, "\n")

	id, entry, err := r.InsertPaper(projectId, chapterId, record.PaperWithoutAutofieldEntry{
		Content: content,
		UserId:  testutil.ModifyOnlyUserId(),
	})
	now := time.Now()

	assert.Nil(t, err)

	assert.Equal(t, "CHAPTER_ONE", id)
	assert.Equal(t, content, entry.Content)
	assert.Equal(t, testutil.ModifyOnlyUserId(), entry.UserId)
	assert.Less(t, now.Sub(entry.CreatedAt), time.Second)
	assert.Less(t, now.Sub(entry.UpdatedAt), time.Second)
}

func TestInsertPaperProjectOrChapterNotFound(t *testing.T) {
	tt := []struct {
		name          string
		userId        string
		projectId     string
		chapterId     string
		expectedError string
	}{
		{
			name:          "should return error when project not found",
			userId:        testutil.ModifyOnlyUserId(),
			projectId:     "UNKNOWN_PROJECT",
			chapterId:     "CHAPTER_ONE",
			expectedError: "project not found",
		},
		{
			name:          "should return not found when user is not author of the project",
			userId:        testutil.ReadOnlyUserId(),
			projectId:     "PROJECT_WITH_DESCRIPTION_TO_UPDATE_FROM_REPOSITORY",
			chapterId:     "CHAPTER_ONE",
			expectedError: "project not found",
		},
		{
			name:          "should return error when chapter not found",
			userId:        testutil.ModifyOnlyUserId(),
			projectId:     "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_REPOSITORY",
			chapterId:     "UNKNOWN_CHAPTER",
			expectedError: "chapter not found",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			client := db.FirestoreClient()
			r := repository.NewPaperRepository(*client)

			id, createdPaper, rErr := r.InsertPaper(tc.projectId, tc.chapterId, record.PaperWithoutAutofieldEntry{
				Content: "content",
				UserId:  tc.userId,
			})

			assert.NotNil(t, rErr)

			assert.Empty(t, id)
			assert.Equal(t, repository.NotFoundError, rErr.Code())
			assert.Equal(t, fmt.Sprintf("not found: %s", tc.expectedError), rErr.Error())
			assert.Nil(t, createdPaper)
		})
	}
}
