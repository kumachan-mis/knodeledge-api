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

func TestFetchProjectChaptersNoProject(t *testing.T) {
	client := db.FirestoreClient()
	r := repository.NewChapterRepository(*client)

	chapters, rErr := r.FetchProjectChapters(testutil.ReadOnlyUserId(), "UNKNOWN_PROJECT")

	assert.NotNil(t, rErr)
	assert.Equal(t, repository.InvalidArgument, rErr.Code())
	assert.Equal(t, "invalid argument: project document does not exist", rErr.Error())
	assert.Nil(t, chapters)
}

func TestFetchProjectChaptersUnauthorizedProject(t *testing.T) {
	client := db.FirestoreClient()
	r := repository.NewChapterRepository(*client)

	chapters, rErr := r.FetchProjectChapters(testutil.ModifyOnlyUserId(), "PROJECT_WITHOUT_DESCRIPTION")

	assert.NotNil(t, rErr)
	assert.Equal(t, repository.InvalidArgument, rErr.Code())
	assert.Equal(t, "invalid argument: project document does not exist", rErr.Error())
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

	projectId := "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_REPOSITORY"

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

func TestInsertProjectChaptersTooLargeChapterNumber(t *testing.T) {
	client := db.FirestoreClient()
	r := repository.NewChapterRepository(*client)

	projectId := "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_REPOSITORY"

	entry := record.ChapterWithoutAutofieldEntry{
		Name:   "Chapter Ninety-Nine",
		Number: 99,
		UserId: testutil.ModifyOnlyUserId(),
	}

	id, createdChapter, rErr := r.InsertChapter(projectId, entry)

	assert.NotNil(t, rErr)

	assert.Empty(t, id)
	assert.Equal(t, repository.InvalidArgument, rErr.Code())
	assert.Equal(t, "invalid argument: chapter number is too large", rErr.Error())
	assert.Nil(t, createdChapter)
}

func TestInsertProjectChaptersNoProject(t *testing.T) {
	client := db.FirestoreClient()
	r := repository.NewChapterRepository(*client)

	projectId := "UNKNOWN_PROJECT"

	entry := record.ChapterWithoutAutofieldEntry{
		Name:   "Chapter One",
		Number: 1,
		UserId: testutil.ModifyOnlyUserId(),
	}

	id, createdChapter, rErr := r.InsertChapter(projectId, entry)

	assert.NotNil(t, rErr)

	assert.Empty(t, id)
	assert.Equal(t, repository.InvalidArgument, rErr.Code())
	assert.Equal(t, "invalid argument: project document does not exist", rErr.Error())
	assert.Nil(t, createdChapter)
}

func TestInsertProjectChaptersUnauthorizedProject(t *testing.T) {
	client := db.FirestoreClient()
	r := repository.NewChapterRepository(*client)

	projectId := "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_REPOSITORY"

	entry := record.ChapterWithoutAutofieldEntry{
		Name:   "Chapter One",
		Number: 1,
		UserId: testutil.ReadOnlyUserId(),
	}

	id, createdChapter, rErr := r.InsertChapter(projectId, entry)

	assert.NotNil(t, rErr)

	assert.Empty(t, id)
	assert.Equal(t, repository.InvalidArgument, rErr.Code())
	assert.Equal(t, "invalid argument: project document does not exist", rErr.Error())
	assert.Nil(t, createdChapter)
}
