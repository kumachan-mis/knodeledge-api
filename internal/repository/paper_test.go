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

func TestFetchPaperValidEntry(t *testing.T) {
	client := db.FirestoreClient()
	r := repository.NewPaperRepository(*client)

	projectId := "PROJECT_WITHOUT_DESCRIPTION"
	chapterId := "CHAPTER_ONE"
	content := strings.Join([]string{
		"[** Introduction]",
		"This is an example project of kNODEledge.",
		"",
		"[** Section of Chapter One]",
		"Section of Chapter One. Section of Chapter One. Section of Chapter One. Section of Chapter One. Section of Chapter One. Section of Chapter One.",
		"Section of Chapter One. Section of Chapter One. Section of Chapter One. Section of Chapter One. Section of Chapter One. Section of Chapter One.",
		""}, "\n")

	entry, err := r.FetchPaper(testutil.ReadOnlyUserId(), projectId, chapterId)

	assert.Nil(t, err)

	assert.Equal(t, content, entry.Content)
	assert.Equal(t, testutil.ReadOnlyUserId(), entry.UserId)
	assert.Equal(t, testutil.Date(), entry.CreatedAt)
	assert.Equal(t, testutil.Date(), entry.UpdatedAt)
}

func TestFetchPaperProjectOrChapterNotFound(t *testing.T) {
	tt := []struct {
		name          string
		userId        string
		projectId     string
		chapterId     string
		expectedError string
	}{
		{
			name:          "should return error when project not found",
			userId:        testutil.ReadOnlyUserId(),
			projectId:     "UNKNOWN_PROJECT",
			chapterId:     "CHAPTER_ONE",
			expectedError: "failed to fetch project",
		},
		{
			name:          "should return not found when user is not author of the project",
			userId:        testutil.ModifyOnlyUserId(),
			projectId:     "PROJECT_WITH_DESCRIPTION",
			chapterId:     "CHAPTER_ONE",
			expectedError: "failed to fetch project",
		},
		{
			name:          "should return error when chapter not found",
			userId:        testutil.ReadOnlyUserId(),
			projectId:     "PROJECT_WITH_DESCRIPTION",
			chapterId:     "UNKNOWN_CHAPTER",
			expectedError: "failed to fetch chapter",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			client := db.FirestoreClient()
			r := repository.NewPaperRepository(*client)

			entry, rErr := r.FetchPaper(tc.userId, tc.projectId, tc.chapterId)

			assert.NotNil(t, rErr)

			assert.Nil(t, entry)
			assert.Equal(t, repository.NotFoundError, rErr.Code())
			assert.Equal(t, fmt.Sprintf("not found: %s", tc.expectedError), rErr.Error())
		})
	}
}

func TestFetchPaperInvalidDocument(t *testing.T) {
	tt := []struct {
		name          string
		userId        string
		projectId     string
		chapterId     string
		expectedError string
	}{
		{
			name:      "should return error when paper content is invalid",
			userId:    testutil.ErrorUserId(6),
			projectId: "PROJECT_WITH_INVALID_PAPER_CONTENT",
			chapterId: "CHAPTER_WITH_INVALID_PAPER_CONTENT",
			expectedError: "failed to convert snapshot to values: document.PaperValues.content: " +
				"firestore: cannot set type string to array",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			client := db.FirestoreClient()
			r := repository.NewPaperRepository(*client)

			entry, rErr := r.FetchPaper(tc.userId, tc.projectId, tc.chapterId)

			assert.NotNil(t, rErr)

			assert.Nil(t, entry)
			assert.Equal(t, repository.ReadFailurePanic, rErr.Code())
			assert.Equal(t, fmt.Sprintf("read failure: %v", tc.expectedError), rErr.Error())
		})
	}
}

func TestInsertPaperValidEntry(t *testing.T) {
	client := db.FirestoreClient()
	r := repository.NewPaperRepository(*client)

	userId := testutil.ModifyOnlyUserId()
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

	id, entry, err := r.InsertPaper(userId, projectId, chapterId, record.PaperWithoutAutofieldEntry{
		Content: content,
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
			r := repository.NewPaperRepository(*client)

			id, createdPaper, rErr := r.InsertPaper(tc.userId, tc.projectId, tc.chapterId, record.PaperWithoutAutofieldEntry{
				Content: "content",
			})

			assert.NotNil(t, rErr)

			assert.Empty(t, id)
			assert.Equal(t, repository.NotFoundError, rErr.Code())
			assert.Equal(t, fmt.Sprintf("not found: %s", tc.expectedError), rErr.Error())
			assert.Nil(t, createdPaper)
		})
	}
}

func TestUpdatePaterValidEntry(t *testing.T) {
	client := db.FirestoreClient()
	r := repository.NewPaperRepository(*client)

	userId := testutil.ModifyOnlyUserId()
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

	entry, err := r.UpdatePaper(userId, projectId, chapterId, record.PaperWithoutAutofieldEntry{
		Content: content,
	})
	now := time.Now()

	assert.Nil(t, err)

	assert.Equal(t, content, entry.Content)
	assert.Equal(t, testutil.ModifyOnlyUserId(), entry.UserId)
	assert.Less(t, now.Sub(entry.CreatedAt), time.Second)
	assert.Less(t, now.Sub(entry.UpdatedAt), time.Second)
}

func TestUpdatePaperProjectOrChapterNotFound(t *testing.T) {
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
			r := repository.NewPaperRepository(*client)

			entry, rErr := r.UpdatePaper(tc.userId, tc.projectId, tc.chapterId, record.PaperWithoutAutofieldEntry{
				Content: "content",
			})

			assert.NotNil(t, rErr)

			assert.Nil(t, entry)
			assert.Equal(t, repository.NotFoundError, rErr.Code())
			assert.Equal(t, fmt.Sprintf("not found: %s", tc.expectedError), rErr.Error())
		})
	}
}
