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

func TestFetchProjectChaptersInvalidArgument(t *testing.T) {
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
			expectedError: "project document does not exist",
		},
		{
			name:          "should return error whenwhen user is not author of the project",
			userId:        testutil.ModifyOnlyUserId(),
			projectId:     "PROJECT_WITHOUT_DESCRIPTION",
			expectedError: "project document does not exist",
		},
	}

	for _, tc := range tt {

		client := db.FirestoreClient()
		r := repository.NewChapterRepository(*client)

		chapters, rErr := r.FetchProjectChapters(tc.userId, tc.projectId)

		assert.NotNil(t, rErr)
		assert.Equal(t, repository.InvalidArgument, rErr.Code())
		assert.Equal(t, fmt.Sprintf("invalid argument: %s", tc.expectedError), rErr.Error())
		assert.Nil(t, chapters)
	}
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
			name:      "should return error when chapter ids are invalid",
			userId:    testutil.ErrorUserId(3),
			projectId: "PROJECT_WITH_INVALID_CHAPTER_IDS",
			expectedError: "failed to convert snapshot to values: document.ProjectWithChapterIdsValues.chapterIds: " +
				"firestore: cannot set type []string to string",
		},
		{
			name:          "should return error when chapter ids have excessive elements",
			userId:        testutil.ErrorUserId(4),
			projectId:     "PROJECT_WITH_TOO_MANY_CHAPTER_IDS",
			expectedError: "failed to convert values to entry: document.ProjectWithChapterIdsValues.chapterIds have excessive elements",
		},
		{
			name:          "should return error when chapter ids have deficient elements",
			userId:        testutil.ErrorUserId(5),
			projectId:     "PROJECT_WITH_TOO_FEW_CHAPTER_IDS",
			expectedError: "failed to convert values to entry: document.ProjectWithChapterIdsValues.chapterIds have deficient elements",
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

func TestInsertChapterValidEntry(t *testing.T) {
	client := db.FirestoreClient()
	r := repository.NewChapterRepository(*client)

	projectId := "PROJECT_WITH_DESCRIPTION_TO_UPDATE_FROM_REPOSITORY"

	entry := record.ChapterWithoutAutofieldEntry{
		Name:   "Chapter One",
		Number: 1,
		UserId: testutil.ModifyOnlyUserId(),
	}

	id1, createdChapter1, rErr := r.InsertChapter(projectId, entry)
	now := time.Now()

	assert.Nil(t, rErr)

	assert.NotEmpty(t, id1)
	assert.Equal(t, "Chapter One", createdChapter1.Name)
	assert.Equal(t, 1, createdChapter1.Number)
	assert.Equal(t, testutil.ModifyOnlyUserId(), createdChapter1.UserId)
	assert.Less(t, now.Sub(createdChapter1.CreatedAt), time.Second)
	assert.Less(t, now.Sub(createdChapter1.UpdatedAt), time.Second)

	entry = record.ChapterWithoutAutofieldEntry{
		Name:   "Chapter Three",
		Number: 2,
		UserId: testutil.ModifyOnlyUserId(),
	}

	id3, createdChapter3, rErr := r.InsertChapter(projectId, entry)
	now = time.Now()

	assert.Nil(t, rErr)

	assert.NotEmpty(t, id3)
	assert.Equal(t, "Chapter Three", createdChapter3.Name)
	assert.Equal(t, 2, createdChapter3.Number)
	assert.Equal(t, testutil.ModifyOnlyUserId(), createdChapter3.UserId)
	assert.Less(t, now.Sub(createdChapter3.CreatedAt), time.Second)
	assert.Less(t, now.Sub(createdChapter3.UpdatedAt), time.Second)

	entry = record.ChapterWithoutAutofieldEntry{
		Name:   "Chapter Two",
		Number: 2,
		UserId: testutil.ModifyOnlyUserId(),
	}

	id2, createdChapter2, rErr := r.InsertChapter(projectId, entry)
	now = time.Now()

	assert.Nil(t, rErr)

	assert.NotEmpty(t, id2)
	assert.Equal(t, "Chapter Two", createdChapter2.Name)
	assert.Equal(t, 2, createdChapter2.Number)
	assert.Equal(t, testutil.ModifyOnlyUserId(), createdChapter2.UserId)
	assert.Less(t, now.Sub(createdChapter2.CreatedAt), time.Second)
	assert.Less(t, now.Sub(createdChapter2.UpdatedAt), time.Second)

	entry = record.ChapterWithoutAutofieldEntry{
		Name:   "Chapter Zero",
		Number: 1,
		UserId: testutil.ModifyOnlyUserId(),
	}

	id0, createdChapter0, rErr := r.InsertChapter(projectId, entry)

	assert.Nil(t, rErr)

	assert.NotEmpty(t, id0)
	assert.Equal(t, "Chapter Zero", createdChapter0.Name)
	assert.Equal(t, 1, createdChapter0.Number)
	assert.Equal(t, testutil.ModifyOnlyUserId(), createdChapter0.UserId)
	assert.Less(t, now.Sub(createdChapter0.CreatedAt), time.Second)
	assert.Less(t, now.Sub(createdChapter0.UpdatedAt), time.Second)

	chapters, rErr := r.FetchProjectChapters(testutil.ModifyOnlyUserId(), projectId)

	assert.Nil(t, rErr)

	assert.Equal(t, map[string]record.ChapterEntry{
		id0: {
			Name:      "Chapter Zero",
			Number:    1,
			UserId:    testutil.ModifyOnlyUserId(),
			CreatedAt: createdChapter0.CreatedAt,
			UpdatedAt: createdChapter0.UpdatedAt,
		},
		id1: {
			Name:      "Chapter One",
			Number:    2,
			UserId:    testutil.ModifyOnlyUserId(),
			CreatedAt: createdChapter1.CreatedAt,
			UpdatedAt: createdChapter1.UpdatedAt,
		},
		id2: {
			Name:      "Chapter Two",
			Number:    3,
			UserId:    testutil.ModifyOnlyUserId(),
			CreatedAt: createdChapter2.CreatedAt,
			UpdatedAt: createdChapter2.UpdatedAt,
		},
		id3: {
			Name:      "Chapter Three",
			Number:    4,
			UserId:    testutil.ModifyOnlyUserId(),
			CreatedAt: createdChapter3.CreatedAt,
			UpdatedAt: createdChapter3.UpdatedAt,
		},
	}, chapters)
}

func TestInsertChapterInvalidArgument(t *testing.T) {
	tt := []struct {
		name          string
		projectId     string
		entry         record.ChapterWithoutAutofieldEntry
		expectedError string
	}{
		{
			name:      "should return error when chapter number is too large",
			projectId: "PROJECT_WITH_DESCRIPTION_TO_UPDATE_FROM_REPOSITORY",
			entry: record.ChapterWithoutAutofieldEntry{
				Name:   "Chapter Ninety-Nine",
				Number: 99,
				UserId: testutil.ModifyOnlyUserId(),
			},
			expectedError: "chapter number is too large",
		},
		{
			name:      "should return error when project is not found",
			projectId: "UNKNOWN_PROJECT",
			entry: record.ChapterWithoutAutofieldEntry{
				Name:   "Chapter One",
				Number: 1,
				UserId: testutil.ModifyOnlyUserId(),
			},
			expectedError: "project document does not exist",
		},
		{
			name:      "should return error when user is not author of the project",
			projectId: "PROJECT_WITH_DESCRIPTION_TO_UPDATE_FROM_REPOSITORY",
			entry: record.ChapterWithoutAutofieldEntry{
				Name:   "Chapter One",
				Number: 1,
				UserId: testutil.ReadOnlyUserId(),
			},
			expectedError: "project document does not exist",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			client := db.FirestoreClient()
			r := repository.NewChapterRepository(*client)

			id, createdChapter, rErr := r.InsertChapter(tc.projectId, tc.entry)

			assert.NotNil(t, rErr)

			assert.Empty(t, id)
			assert.Equal(t, repository.InvalidArgument, rErr.Code())
			assert.Equal(t, fmt.Sprintf("invalid argument: %s", tc.expectedError), rErr.Error())
			assert.Nil(t, createdChapter)
		})
	}
}

func TestUpdateChapterValidEntry(t *testing.T) {
	client := db.FirestoreClient()
	r := repository.NewChapterRepository(*client)

	projectId := "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_REPOSITORY"

	entry := record.ChapterWithoutAutofieldEntry{
		Name:   "Chapter One",
		Number: 1,
		UserId: testutil.ModifyOnlyUserId(),
	}

	updatedChapter2, rErr := r.UpdateChapter(projectId, "CHAPTER_TWO", entry)
	now := time.Now()

	assert.Nil(t, rErr)

	assert.Equal(t, "Chapter One", updatedChapter2.Name)
	assert.Equal(t, 1, updatedChapter2.Number)
	assert.Equal(t, testutil.ModifyOnlyUserId(), updatedChapter2.UserId)
	assert.Equal(t, testutil.Date(), updatedChapter2.CreatedAt)
	assert.Less(t, now.Sub(updatedChapter2.UpdatedAt), time.Second)

	chapters, err := r.FetchProjectChapters(testutil.ModifyOnlyUserId(), projectId)

	assert.Nil(t, err)

	assert.Equal(t, map[string]record.ChapterEntry{
		"CHAPTER_TWO": {
			Name:      "Chapter One",
			Number:    1,
			UserId:    testutil.ModifyOnlyUserId(),
			CreatedAt: testutil.Date(),
			UpdatedAt: updatedChapter2.UpdatedAt,
		},
		"CHAPTER_ONE": {
			Name:      "Chapter One",
			Number:    2,
			UserId:    testutil.ModifyOnlyUserId(),
			CreatedAt: testutil.Date(),
			UpdatedAt: testutil.Date(),
		},
	}, chapters)

	entry = record.ChapterWithoutAutofieldEntry{
		Name:   "Chapter Two",
		Number: 2,
		UserId: testutil.ModifyOnlyUserId(),
	}

	updatedChapter1, rErr := r.UpdateChapter(projectId, "CHAPTER_ONE", entry)

	assert.Nil(t, rErr)

	assert.Equal(t, "Chapter Two", updatedChapter1.Name)
	assert.Equal(t, 2, updatedChapter1.Number)
	assert.Equal(t, testutil.ModifyOnlyUserId(), updatedChapter1.UserId)
	assert.Equal(t, testutil.Date(), updatedChapter1.CreatedAt)
	assert.Less(t, now.Sub(updatedChapter1.UpdatedAt), time.Second)

	chapters, err = r.FetchProjectChapters(testutil.ModifyOnlyUserId(), projectId)

	assert.Nil(t, err)

	assert.Equal(t, map[string]record.ChapterEntry{
		"CHAPTER_TWO": {
			Name:      "Chapter One",
			Number:    1,
			UserId:    testutil.ModifyOnlyUserId(),
			CreatedAt: testutil.Date(),
			UpdatedAt: updatedChapter2.UpdatedAt,
		},
		"CHAPTER_ONE": {
			Name:      "Chapter Two",
			Number:    2,
			UserId:    testutil.ModifyOnlyUserId(),
			CreatedAt: testutil.Date(),
			UpdatedAt: updatedChapter1.UpdatedAt,
		},
	}, chapters)
}

func TestUpdateChapterNoDocument(t *testing.T) {
	client := db.FirestoreClient()
	r := repository.NewChapterRepository(*client)

	updatedChapter, rErr := r.UpdateChapter(
		"PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_REPOSITORY",
		"UNKNOWN_CHAPTER",
		record.ChapterWithoutAutofieldEntry{
			Name:   "Chapter One",
			Number: 1,
			UserId: testutil.ModifyOnlyUserId(),
		})

	assert.NotNil(t, rErr)

	assert.Nil(t, updatedChapter)
	assert.Equal(t, repository.NotFoundError, rErr.Code())
	assert.Equal(t, "not found: failed to update chapter", rErr.Error())
}

func TestUpdateChapterInvalidArgument(t *testing.T) {
	tt := []struct {
		name          string
		projectId     string
		chapterId     string
		entry         record.ChapterWithoutAutofieldEntry
		expectedError string
	}{
		{
			name:      "should return error when chapter number is too large",
			projectId: "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_REPOSITORY",
			chapterId: "CHAPTER_ONE",
			entry: record.ChapterWithoutAutofieldEntry{
				Name:   "Chapter Ninety-Nine",
				Number: 99,
				UserId: testutil.ModifyOnlyUserId(),
			},
			expectedError: "chapter number is too large",
		},
		{
			name:      "should return error when project is not found",
			projectId: "UNKNOWN_PROJECT",
			chapterId: "CHAPTER_ONE",
			entry: record.ChapterWithoutAutofieldEntry{
				Name:   "Chapter One",
				Number: 1,
				UserId: testutil.ModifyOnlyUserId(),
			},
			expectedError: "project document does not exist",
		},
		{
			name:      "should return error when user is not author of the project",
			projectId: "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_REPOSITORY",
			chapterId: "CHAPTER_ONE",
			entry: record.ChapterWithoutAutofieldEntry{
				Name:   "Chapter One",
				Number: 1,
				UserId: testutil.ReadOnlyUserId(),
			},
			expectedError: "project document does not exist",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			client := db.FirestoreClient()
			r := repository.NewChapterRepository(*client)

			updatedChapter, rErr := r.UpdateChapter(tc.projectId, tc.chapterId, tc.entry)

			assert.NotNil(t, rErr)

			assert.Nil(t, updatedChapter)
			assert.Equal(t, repository.InvalidArgument, rErr.Code())
			assert.Equal(t, fmt.Sprintf("invalid argument: %s", tc.expectedError), rErr.Error())
		})
	}
}

func TestUpdateChapterInvalidDocument(t *testing.T) {
	tt := []struct {
		name          string
		userId        string
		projectId     string
		chapterId     string
		expectedError string
	}{
		{
			name:      "should return error when chapter ids are invalid",
			userId:    testutil.ErrorUserId(3),
			projectId: "PROJECT_WITH_INVALID_CHAPTER_IDS",
			chapterId: "CHAPTER",
			expectedError: "failed to convert snapshot to values: document.ProjectWithChapterIdsValues.chapterIds: " +
				"firestore: cannot set type []string to string",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			client := db.FirestoreClient()
			r := repository.NewChapterRepository(*client)

			updatedChapter, rErr := r.UpdateChapter(
				tc.projectId,
				tc.chapterId,
				record.ChapterWithoutAutofieldEntry{
					Name:   "Updated Chapter",
					Number: 1,
					UserId: tc.userId,
				})

			assert.NotNil(t, rErr)

			assert.Nil(t, updatedChapter)
			assert.Equal(t, repository.ReadFailurePanic, rErr.Code())
			assert.Equal(t, fmt.Sprintf("read failure: %s", tc.expectedError), rErr.Error())
		})
	}
}
