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

func TestGraphExistsValidEntry(t *testing.T) {
	tt := []struct {
		name      string
		userId    string
		projectId string
		chapterId string
		expected  bool
	}{
		{
			name:      "should return true when graph exists",
			userId:    testutil.ReadOnlyUserId(),
			projectId: "PROJECT_WITHOUT_DESCRIPTION",
			chapterId: "CHAPTER_ONE",
			expected:  true,
		},
		{
			name:      "should return false when graph not exists",
			userId:    testutil.ReadOnlyUserId(),
			projectId: "PROJECT_WITHOUT_DESCRIPTION",
			chapterId: "CHAPTER_TWO",
			expected:  false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			client := db.FirestoreClient()
			r := repository.NewGraphRepository(*client)

			exists, rErr := r.GraphExists(tc.userId, tc.projectId, tc.chapterId)

			assert.Nil(t, rErr)

			assert.Equal(t, tc.expected, exists)
		})
	}
}

func TestGraphExistsProjectOrChapterNotFound(t *testing.T) {
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
			expectedError: "failed to fetch project",
		},
		{
			name:          "should return not found when user is not author of the project",
			userId:        testutil.ReadOnlyUserId(),
			projectId:     "PROJECT_WITH_DESCRIPTION_TO_UPDATE_FROM_REPOSITORY",
			chapterId:     "CHAPTER_ONE",
			expectedError: "failed to fetch project",
		},
		{
			name:          "should return error when chapter not found",
			userId:        testutil.ModifyOnlyUserId(),
			projectId:     "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_REPOSITORY",
			chapterId:     "UNKNOWN_CHAPTER",
			expectedError: "failed to fetch chapter",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			client := db.FirestoreClient()
			r := repository.NewGraphRepository(*client)

			exists, rErr := r.GraphExists(tc.userId, tc.projectId, tc.chapterId)

			assert.NotNil(t, rErr)

			assert.False(t, exists)
			assert.Equal(t, repository.NotFoundError, rErr.Code())
			assert.Equal(t, fmt.Sprintf("not found: %s", tc.expectedError), rErr.Error())
		})
	}
}

func TestInsertGraphsValidEntry(t *testing.T) {
	client := db.FirestoreClient()
	r := repository.NewGraphRepository(*client)

	userId := testutil.ModifyOnlyUserId()
	projectId := "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_REPOSITORY"
	chapterId := "CHAPTER_ONE"

	paragraph1 := strings.Join([]string{
		"## Introduction",
		"This is the introduction of the paper.",
	}, "\n")
	paragraph2 := strings.Join([]string{
		"## What is note apps?",
		"Note apps is a web application that allows users to create, read, update, and delete notes.",
	}, "\n")
	entries := []record.GraphWithoutAutofieldEntry{
		{
			Name:      "Introduction",
			Paragraph: paragraph1,
		},
		{
			Name:      "What is note apps?",
			Paragraph: paragraph2,
		},
	}

	ids, createdEntries, rErr := r.InsertGraphs(userId, projectId, chapterId, entries)
	now := time.Now()

	assert.Nil(t, rErr)

	assert.Len(t, ids, 2)
	assert.Len(t, createdEntries, 2)

	id := ids[0]
	createdEntry := createdEntries[0]
	assert.NotEmpty(t, id)
	assert.Equal(t, "Introduction", createdEntry.Name)
	assert.Equal(t, paragraph1, createdEntry.Paragraph)
	assert.Equal(t, testutil.ModifyOnlyUserId(), createdEntry.UserId)
	assert.Less(t, now.Sub(createdEntry.CreatedAt), time.Second)
	assert.Less(t, now.Sub(createdEntry.UpdatedAt), time.Second)

	id = ids[1]
	createdEntry = createdEntries[1]
	assert.NotEmpty(t, id)
	assert.Equal(t, "What is note apps?", createdEntry.Name)
	assert.Equal(t, paragraph2, createdEntry.Paragraph)
	assert.Equal(t, testutil.ModifyOnlyUserId(), createdEntry.UserId)
	assert.Less(t, now.Sub(createdEntry.CreatedAt), time.Second)
	assert.Less(t, now.Sub(createdEntry.UpdatedAt), time.Second)
}

func TestInsertGraphProjectOrChapterNotFound(t *testing.T) {
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
			expectedError: "failed to fetch project",
		},
		{
			name:          "should return not found when user is not author of the project",
			userId:        testutil.ReadOnlyUserId(),
			projectId:     "PROJECT_WITH_DESCRIPTION_TO_UPDATE_FROM_REPOSITORY",
			chapterId:     "CHAPTER_ONE",
			expectedError: "failed to fetch project",
		},
		{
			name:          "should return error when chapter not found",
			userId:        testutil.ModifyOnlyUserId(),
			projectId:     "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_REPOSITORY",
			chapterId:     "UNKNOWN_CHAPTER",
			expectedError: "failed to fetch chapter",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			client := db.FirestoreClient()
			r := repository.NewGraphRepository(*client)

			id, createdPaper, rErr := r.InsertGraphs(tc.userId, tc.projectId, tc.chapterId, []record.GraphWithoutAutofieldEntry{
				{
					Name:      "Section Name",
					Paragraph: "paragraph",
				},
			})

			assert.NotNil(t, rErr)

			assert.Empty(t, id)
			assert.Equal(t, repository.NotFoundError, rErr.Code())
			assert.Equal(t, fmt.Sprintf("not found: %s", tc.expectedError), rErr.Error())
			assert.Nil(t, createdPaper)
		})
	}
}