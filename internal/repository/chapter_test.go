package repository_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/kumachan-mis/knodeledge-api/internal/db"
	"github.com/kumachan-mis/knodeledge-api/internal/record"
	"github.com/kumachan-mis/knodeledge-api/internal/repository"
	"github.com/kumachan-mis/knodeledge-api/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestFetchChaptersValidDocument(t *testing.T) {
	client := db.FirestoreClient()
	r := repository.NewChapterRepository(*client)

	chapters, rErr := r.FetchChapters(testutil.ReadOnlyUserId(), "PROJECT_WITHOUT_DESCRIPTION")

	assert.Nil(t, rErr)

	assert.Len(t, chapters, 2)

	chapter := chapters["CHAPTER_ONE"]
	assert.Equal(t, "Chapter One", chapter.Name)
	assert.Equal(t, 1, chapter.Number)
	assert.Len(t, chapter.Sections, 2)
	assert.Equal(t, "SECTION_ONE", chapter.Sections[0].Id)
	assert.Equal(t, "Introduction", chapter.Sections[0].Name)
	assert.Equal(t, "SECTION_TWO", chapter.Sections[1].Id)
	assert.Equal(t, "Section of Chapter One", chapter.Sections[1].Name)
	assert.Equal(t, testutil.Date(), chapter.CreatedAt)
	assert.Equal(t, testutil.Date(), chapter.UpdatedAt)

	chapter = chapters["CHAPTER_TWO"]
	assert.Equal(t, "Chapter Two", chapter.Name)
	assert.Equal(t, 2, chapter.Number)
	assert.Len(t, chapter.Sections, 0)
	assert.Equal(t, testutil.Date(), chapter.CreatedAt)
	assert.Equal(t, testutil.Date(), chapter.UpdatedAt)
}

func TestFetchChaptersNoDocument(t *testing.T) {
	client := db.FirestoreClient()
	r := repository.NewChapterRepository(*client)

	chapters, rErr := r.FetchChapters(testutil.ReadOnlyUserId(), "PROJECT_WITH_DESCRIPTION")

	assert.Nil(t, rErr)

	assert.Empty(t, chapters)
}

func TestFetchChaptersNotFound(t *testing.T) {
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
			expectedError: "failed to fetch project",
		},
		{
			name:          "should return error when user is not author of the project",
			userId:        testutil.ModifyOnlyUserId(),
			projectId:     "PROJECT_WITHOUT_DESCRIPTION",
			expectedError: "failed to fetch project",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			client := db.FirestoreClient()
			r := repository.NewChapterRepository(*client)

			chapters, rErr := r.FetchChapters(tc.userId, tc.projectId)

			assert.NotNil(t, rErr)
			assert.Equal(t, repository.NotFoundError, rErr.Code())
			assert.Equal(t, fmt.Sprintf("not found: %v", tc.expectedError), rErr.Error())
			assert.Nil(t, chapters)
		})
	}
}

func TestFetchChaptersInvalidDocument(t *testing.T) {
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
			expectedError: "failed to convert snapshot to values: document.ProjectValues.chapterIds: " +
				"firestore: cannot set type []string to string",
		},
		{
			name:          "should return error when chapter ids have excessive elements",
			userId:        testutil.ErrorUserId(4),
			projectId:     "PROJECT_WITH_TOO_MANY_CHAPTER_IDS",
			expectedError: "failed to convert values to entry: document.ProjectValues.chapterIds have excessive elements",
		},
		{
			name:          "should return error when chapter ids have insufficient elements",
			userId:        testutil.ErrorUserId(5),
			projectId:     "PROJECT_WITH_TOO_FEW_CHAPTER_IDS",
			expectedError: "failed to convert values to entry: document.ProjectValues.chapterIds have insufficient elements",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			client := db.FirestoreClient()
			r := repository.NewChapterRepository(*client)

			chapters, rErr := r.FetchChapters(tc.userId, tc.projectId)

			assert.NotNil(t, rErr)
			assert.Equal(t, repository.ReadFailurePanic, rErr.Code())
			assert.Equal(t, fmt.Sprintf("read failure: %v", tc.expectedError), rErr.Error())
			assert.Nil(t, chapters)
		})
	}
}

func TestFetchChapterValidDocument(t *testing.T) {
	client := db.FirestoreClient()
	r := repository.NewChapterRepository(*client)

	chapter, rErr := r.FetchChapter(testutil.ReadOnlyUserId(), "PROJECT_WITHOUT_DESCRIPTION", "CHAPTER_ONE")

	assert.Nil(t, rErr)

	assert.Equal(t, "Chapter One", chapter.Name)
	assert.Equal(t, 1, chapter.Number)
	assert.Len(t, chapter.Sections, 2)
	assert.Equal(t, "SECTION_ONE", chapter.Sections[0].Id)
	assert.Equal(t, "Introduction", chapter.Sections[0].Name)
	assert.Equal(t, "SECTION_TWO", chapter.Sections[1].Id)
	assert.Equal(t, "Section of Chapter One", chapter.Sections[1].Name)
	assert.Equal(t, testutil.Date(), chapter.CreatedAt)
	assert.Equal(t, testutil.Date(), chapter.UpdatedAt)
}

func TestFetchChapterNotFound(t *testing.T) {
	tt := []struct {
		name          string
		userId        string
		projectId     string
		chapterId     string
		expectedError string
	}{
		{
			name:          "should return error when project is not found",
			userId:        testutil.ReadOnlyUserId(),
			projectId:     "UNKNOWN_PROJECT",
			chapterId:     "UNKNOWN_CHAPTER",
			expectedError: "failed to fetch project",
		},
		{
			name:          "should return error when user is not author of the project",
			userId:        testutil.ModifyOnlyUserId(),
			projectId:     "PROJECT_WITHOUT_DESCRIPTION",
			chapterId:     "CHAPTER_ONE",
			expectedError: "failed to fetch project",
		},
		{
			name:          "should return error when chapter is not found",
			userId:        testutil.ReadOnlyUserId(),
			projectId:     "PROJECT_WITHOUT_DESCRIPTION",
			chapterId:     "UNKNOWN_CHAPTER",
			expectedError: "failed to fetch chapter",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			client := db.FirestoreClient()
			r := repository.NewChapterRepository(*client)

			chapters, rErr := r.FetchChapter(tc.userId, tc.projectId, tc.chapterId)

			assert.NotNil(t, rErr)
			assert.Equal(t, repository.NotFoundError, rErr.Code())
			assert.Equal(t, fmt.Sprintf("not found: %v", tc.expectedError), rErr.Error())
			assert.Nil(t, chapters)
		})
	}
}

func TestFetchChapterInvalidDocument(t *testing.T) {
	tt := []struct {
		name          string
		userId        string
		projectId     string
		chapterId     string
		expectedError string
	}{
		{
			name:      "should return error when chapter name is invalid",
			userId:    testutil.ErrorUserId(2),
			projectId: "PROJECT_WITH_INVALID_CHAPTER_NAME",
			chapterId: "CHAPTER_WITH_INVALID_NAME",
			expectedError: "failed to convert snapshot to values: document.ChapterValues.name: " +
				"firestore: cannot set type string to bool",
		},
		{
			name:      "should return error when chapter ids are invalid",
			userId:    testutil.ErrorUserId(3),
			projectId: "PROJECT_WITH_INVALID_CHAPTER_IDS",
			chapterId: "CHAPTER",
			expectedError: "failed to convert snapshot to values: document.ProjectValues.chapterIds: " +
				"firestore: cannot set type []string to string",
		},
		{
			name:          "should return error when chapter ids have excessive elements",
			userId:        testutil.ErrorUserId(4),
			projectId:     "PROJECT_WITH_TOO_MANY_CHAPTER_IDS",
			chapterId:     "UNKNOWN_CHAPTER",
			expectedError: "failed to convert values to entry: document.ProjectValues.chapterIds have excessive elements",
		},
		{
			name:          "should return error when chapter ids have insufficient elements",
			userId:        testutil.ErrorUserId(5),
			projectId:     "PROJECT_WITH_TOO_FEW_CHAPTER_IDS",
			chapterId:     "CHAPTER_UNKNOWN",
			expectedError: "failed to convert values to entry: document.ProjectValues.chapterIds have insufficient elements",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			client := db.FirestoreClient()
			r := repository.NewChapterRepository(*client)

			chapters, rErr := r.FetchChapter(tc.userId, tc.projectId, tc.chapterId)

			assert.NotNil(t, rErr)
			assert.Equal(t, repository.ReadFailurePanic, rErr.Code())
			assert.Equal(t, fmt.Sprintf("read failure: %v", tc.expectedError), rErr.Error())
			assert.Nil(t, chapters)
		})
	}
}

func TestInsertChapterValidEntry(t *testing.T) {
	client := db.FirestoreClient()
	r := repository.NewChapterRepository(*client)

	userId := testutil.ModifyOnlyUserId()
	projectId := "PROJECT_WITH_DESCRIPTION_TO_UPDATE_FROM_REPOSITORY"

	entry := record.ChapterWithoutAutofieldEntry{
		Name:   "Chapter One",
		Number: 1,
	}

	id1, createdChapter1, rErr := r.InsertChapter(userId, projectId, entry)
	now := time.Now()

	assert.Nil(t, rErr)

	assert.NotEmpty(t, id1)
	assert.Equal(t, "Chapter One", createdChapter1.Name)
	assert.Equal(t, 1, createdChapter1.Number)
	assert.Equal(t, []record.SectionEntry{}, createdChapter1.Sections)
	assert.Equal(t, testutil.ModifyOnlyUserId(), createdChapter1.UserId)
	assert.Less(t, now.Sub(createdChapter1.CreatedAt), time.Second)
	assert.Less(t, now.Sub(createdChapter1.UpdatedAt), time.Second)

	entry = record.ChapterWithoutAutofieldEntry{
		Name:   "Chapter Three",
		Number: 2,
	}

	id3, createdChapter3, rErr := r.InsertChapter(userId, projectId, entry)
	now = time.Now()

	assert.Nil(t, rErr)

	assert.NotEmpty(t, id3)
	assert.Equal(t, "Chapter Three", createdChapter3.Name)
	assert.Equal(t, 2, createdChapter3.Number)
	assert.Equal(t, []record.SectionEntry{}, createdChapter1.Sections)
	assert.Equal(t, testutil.ModifyOnlyUserId(), createdChapter3.UserId)
	assert.Less(t, now.Sub(createdChapter3.CreatedAt), time.Second)
	assert.Less(t, now.Sub(createdChapter3.UpdatedAt), time.Second)

	entry = record.ChapterWithoutAutofieldEntry{
		Name:   "Chapter Two",
		Number: 2,
	}

	id2, createdChapter2, rErr := r.InsertChapter(userId, projectId, entry)
	now = time.Now()

	assert.Nil(t, rErr)

	assert.NotEmpty(t, id2)
	assert.Equal(t, "Chapter Two", createdChapter2.Name)
	assert.Equal(t, 2, createdChapter2.Number)
	assert.Equal(t, []record.SectionEntry{}, createdChapter1.Sections)
	assert.Equal(t, testutil.ModifyOnlyUserId(), createdChapter2.UserId)
	assert.Less(t, now.Sub(createdChapter2.CreatedAt), time.Second)
	assert.Less(t, now.Sub(createdChapter2.UpdatedAt), time.Second)

	entry = record.ChapterWithoutAutofieldEntry{
		Name:   "Chapter Zero",
		Number: 1,
	}

	id0, createdChapter0, rErr := r.InsertChapter(userId, projectId, entry)

	assert.Nil(t, rErr)

	assert.NotEmpty(t, id0)
	assert.Equal(t, "Chapter Zero", createdChapter0.Name)
	assert.Equal(t, 1, createdChapter0.Number)
	assert.Equal(t, []record.SectionEntry{}, createdChapter1.Sections)
	assert.Equal(t, testutil.ModifyOnlyUserId(), createdChapter0.UserId)
	assert.Less(t, now.Sub(createdChapter0.CreatedAt), time.Second)
	assert.Less(t, now.Sub(createdChapter0.UpdatedAt), time.Second)

	chapters, rErr := r.FetchChapters(testutil.ModifyOnlyUserId(), projectId)

	assert.Nil(t, rErr)

	assert.Equal(t, map[string]record.ChapterEntry{
		id0: {
			Name:      "Chapter Zero",
			Number:    1,
			Sections:  []record.SectionEntry{},
			UserId:    testutil.ModifyOnlyUserId(),
			CreatedAt: createdChapter0.CreatedAt,
			UpdatedAt: createdChapter0.UpdatedAt,
		},
		id1: {
			Name:      "Chapter One",
			Number:    2,
			Sections:  []record.SectionEntry{},
			UserId:    testutil.ModifyOnlyUserId(),
			CreatedAt: createdChapter1.CreatedAt,
			UpdatedAt: createdChapter1.UpdatedAt,
		},
		id2: {
			Name:      "Chapter Two",
			Number:    3,
			Sections:  []record.SectionEntry{},
			UserId:    testutil.ModifyOnlyUserId(),
			CreatedAt: createdChapter2.CreatedAt,
			UpdatedAt: createdChapter2.UpdatedAt,
		},
		id3: {
			Name:      "Chapter Three",
			Number:    4,
			Sections:  []record.SectionEntry{},
			UserId:    testutil.ModifyOnlyUserId(),
			CreatedAt: createdChapter3.CreatedAt,
			UpdatedAt: createdChapter3.UpdatedAt,
		},
	}, chapters)
}

func TestInsertChapterNotFound(t *testing.T) {
	tt := []struct {
		name          string
		userId        string
		projectId     string
		expectedError string
	}{
		{
			name:          "should return error when project is not found",
			userId:        testutil.ModifyOnlyUserId(),
			projectId:     "UNKNOWN_PROJECT",
			expectedError: "failed to fetch project",
		},
		{
			name:          "should return not found when user is not author of the project",
			userId:        testutil.ReadOnlyUserId(),
			projectId:     "PROJECT_WITH_DESCRIPTION_TO_UPDATE_FROM_REPOSITORY",
			expectedError: "failed to fetch project",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			client := db.FirestoreClient()
			r := repository.NewChapterRepository(*client)

			id, createdChapter, rErr := r.InsertChapter(tc.userId, tc.projectId, record.ChapterWithoutAutofieldEntry{
				Name:   "Chapter One",
				Number: 1,
			})

			assert.NotNil(t, rErr)

			assert.Empty(t, id)
			assert.Equal(t, repository.NotFoundError, rErr.Code())
			assert.Equal(t, fmt.Sprintf("not found: %v", tc.expectedError), rErr.Error())
			assert.Nil(t, createdChapter)
		})
	}
}

func TestInsertChapterInvalidArgument(t *testing.T) {
	tt := []struct {
		name          string
		userId        string
		projectId     string
		entry         record.ChapterWithoutAutofieldEntry
		expectedError string
	}{
		{
			name:      "should return error when chapter number is too large",
			userId:    testutil.ModifyOnlyUserId(),
			projectId: "PROJECT_WITH_DESCRIPTION_TO_UPDATE_FROM_REPOSITORY",
			entry: record.ChapterWithoutAutofieldEntry{
				Name:   "Chapter Ninety-Nine",
				Number: 99,
			},
			expectedError: "chapter number is too large",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			client := db.FirestoreClient()
			r := repository.NewChapterRepository(*client)

			id, createdChapter, rErr := r.InsertChapter(tc.userId, tc.projectId, tc.entry)

			assert.NotNil(t, rErr)

			assert.Empty(t, id)
			assert.Equal(t, repository.InvalidArgument, rErr.Code())
			assert.Equal(t, fmt.Sprintf("invalid argument: %v", tc.expectedError), rErr.Error())
			assert.Nil(t, createdChapter)
		})
	}
}

func TestUpdateChapterValidEntry(t *testing.T) {
	client := db.FirestoreClient()
	r := repository.NewChapterRepository(*client)

	userId := testutil.ModifyOnlyUserId()
	projectId := "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_REPOSITORY"

	entry := record.ChapterWithoutAutofieldEntry{
		Name:   "Chapter One",
		Number: 1,
	}

	updatedChapter, rErr := r.UpdateChapter(userId, projectId, "CHAPTER_TWO", entry)
	now := time.Now()

	assert.Nil(t, rErr)

	assert.Equal(t, "Chapter One", updatedChapter.Name)
	assert.Equal(t, 1, updatedChapter.Number)
	assert.Equal(t, []record.SectionEntry{}, updatedChapter.Sections)
	assert.Equal(t, testutil.ModifyOnlyUserId(), updatedChapter.UserId)
	assert.Equal(t, testutil.Date(), updatedChapter.CreatedAt)
	assert.Less(t, now.Sub(updatedChapter.UpdatedAt), time.Second)

	chapters, err := r.FetchChapters(testutil.ModifyOnlyUserId(), projectId)

	assert.Nil(t, err)

	assert.Equal(t, map[string]record.ChapterEntry{
		"CHAPTER_TWO": {
			Name:      "Chapter One",
			Number:    1,
			Sections:  []record.SectionEntry{},
			UserId:    testutil.ModifyOnlyUserId(),
			CreatedAt: testutil.Date(),
			UpdatedAt: updatedChapter.UpdatedAt,
		},
		"CHAPTER_ONE": {
			Name:   "Chapter One",
			Number: 2,
			Sections: []record.SectionEntry{
				{
					Id:        "SECTION_ONE",
					Name:      "Introduction",
					UserId:    testutil.ModifyOnlyUserId(),
					CreatedAt: testutil.Date(),
					UpdatedAt: testutil.Date(),
				},
				{
					Id:        "SECTION_TWO",
					Name:      "Section of Chapter One",
					UserId:    testutil.ModifyOnlyUserId(),
					CreatedAt: testutil.Date(),
					UpdatedAt: testutil.Date(),
				},
			},
			UserId:    testutil.ModifyOnlyUserId(),
			CreatedAt: testutil.Date(),
			UpdatedAt: testutil.Date(),
		},
	}, chapters)

	entry = record.ChapterWithoutAutofieldEntry{
		Name:   "Chapter Two",
		Number: 2,
	}

	updatedChapter, rErr = r.UpdateChapter(userId, projectId, "CHAPTER_TWO", entry)

	assert.Nil(t, rErr)

	assert.Equal(t, "Chapter Two", updatedChapter.Name)
	assert.Equal(t, 2, updatedChapter.Number)
	assert.Equal(t, []record.SectionEntry{}, updatedChapter.Sections)
	assert.Equal(t, testutil.ModifyOnlyUserId(), updatedChapter.UserId)
	assert.Equal(t, testutil.Date(), updatedChapter.CreatedAt)
	assert.Less(t, now.Sub(updatedChapter.UpdatedAt), time.Second)

	chapters, err = r.FetchChapters(testutil.ModifyOnlyUserId(), projectId)

	assert.Nil(t, err)

	assert.Equal(t, map[string]record.ChapterEntry{
		"CHAPTER_ONE": {
			Name:   "Chapter One",
			Number: 1,
			Sections: []record.SectionEntry{
				{
					Id:        "SECTION_ONE",
					Name:      "Introduction",
					UserId:    testutil.ModifyOnlyUserId(),
					CreatedAt: testutil.Date(),
					UpdatedAt: testutil.Date(),
				},
				{
					Id:        "SECTION_TWO",
					Name:      "Section of Chapter One",
					UserId:    testutil.ModifyOnlyUserId(),
					CreatedAt: testutil.Date(),
					UpdatedAt: testutil.Date(),
				},
			},
			UserId:    testutil.ModifyOnlyUserId(),
			CreatedAt: testutil.Date(),
			UpdatedAt: testutil.Date(),
		},
		"CHAPTER_TWO": {
			Name:      "Chapter Two",
			Number:    2,
			Sections:  []record.SectionEntry{},
			UserId:    testutil.ModifyOnlyUserId(),
			CreatedAt: testutil.Date(),
			UpdatedAt: updatedChapter.UpdatedAt,
		},
	}, chapters)
}

func TestUpdateChapterNotFound(t *testing.T) {
	tt := []struct {
		name          string
		userId        string
		projectId     string
		chapterId     string
		expectedError string
	}{
		{
			name:          "should return error when project is not found",
			userId:        testutil.ModifyOnlyUserId(),
			projectId:     "UNKNOWN_PROJECT",
			chapterId:     "CHAPTER_ONE",
			expectedError: "failed to fetch project",
		},
		{
			name:          "should return not found when user is not author of the project",
			userId:        testutil.ReadOnlyUserId(),
			projectId:     "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_REPOSITORY",
			chapterId:     "CHAPTER_ONE",
			expectedError: "failed to fetch project",
		},
		{
			name:          "should return error when chapter is not found",
			userId:        testutil.ModifyOnlyUserId(),
			projectId:     "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_REPOSITORY",
			chapterId:     "UNKNOWN_CHAPTER",
			expectedError: "failed to update chapter",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			client := db.FirestoreClient()
			r := repository.NewChapterRepository(*client)

			updatedChapter, rErr := r.UpdateChapter(
				tc.userId,
				tc.projectId,
				tc.chapterId,
				record.ChapterWithoutAutofieldEntry{
					Name:   "Chapter One",
					Number: 1,
				})

			assert.NotNil(t, rErr)

			assert.Nil(t, updatedChapter)
			assert.Equal(t, repository.NotFoundError, rErr.Code())
			assert.Equal(t, fmt.Sprintf("not found: %v", tc.expectedError), rErr.Error())
		})
	}
}

func TestUpdateChapterInvalidArgument(t *testing.T) {
	tt := []struct {
		name          string
		userId        string
		projectId     string
		chapterId     string
		entry         record.ChapterWithoutAutofieldEntry
		expectedError string
	}{
		{
			name:      "should return error when chapter number is too large",
			userId:    testutil.ModifyOnlyUserId(),
			projectId: "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_REPOSITORY",
			chapterId: "CHAPTER_ONE",
			entry: record.ChapterWithoutAutofieldEntry{
				Name:   "Chapter Ninety-Nine",
				Number: 99,
			},
			expectedError: "chapter number is too large",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			client := db.FirestoreClient()
			r := repository.NewChapterRepository(*client)

			updatedChapter, rErr := r.UpdateChapter(tc.userId, tc.projectId, tc.chapterId, tc.entry)

			assert.NotNil(t, rErr)

			assert.Nil(t, updatedChapter)
			assert.Equal(t, repository.InvalidArgument, rErr.Code())
			assert.Equal(t, fmt.Sprintf("invalid argument: %v", tc.expectedError), rErr.Error())
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
			expectedError: "failed to convert snapshot to values: document.ProjectValues.chapterIds: " +
				"firestore: cannot set type []string to string",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			client := db.FirestoreClient()
			r := repository.NewChapterRepository(*client)

			updatedChapter, rErr := r.UpdateChapter(
				tc.userId,
				tc.projectId,
				tc.chapterId,
				record.ChapterWithoutAutofieldEntry{
					Name:   "Updated Chapter",
					Number: 1,
				})

			assert.NotNil(t, rErr)

			assert.Nil(t, updatedChapter)
			assert.Equal(t, repository.ReadFailurePanic, rErr.Code())
			assert.Equal(t, fmt.Sprintf("read failure: %v", tc.expectedError), rErr.Error())
		})
	}
}

func TestUpdateChapterSectionsValidEntry(t *testing.T) {
	client := db.FirestoreClient()
	r := repository.NewChapterRepository(*client)

	userId := testutil.ModifyOnlyUserId()
	projectId := "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_REPOSITORY"
	chapterId := "CHAPTER_ONE"

	updatedSections, rErr := r.UpdateChapterSections(userId, projectId, chapterId, []record.SectionWithoutAutofieldEntry{})
	now := time.Now()

	assert.Nil(t, rErr)

	assert.Len(t, updatedSections, 0)

	updatedChapter, err := r.FetchChapter(testutil.ModifyOnlyUserId(), projectId, chapterId)

	assert.Nil(t, err)

	assert.Equal(t, record.ChapterEntry{
		Name:      "Chapter One",
		Number:    1,
		Sections:  []record.SectionEntry{},
		UserId:    testutil.ModifyOnlyUserId(),
		CreatedAt: testutil.Date(),
		UpdatedAt: updatedChapter.UpdatedAt,
	}, *updatedChapter)
	assert.Less(t, now.Sub(updatedChapter.UpdatedAt), time.Second)

	updatedSections, rErr = r.UpdateChapterSections(userId, projectId, chapterId, []record.SectionWithoutAutofieldEntry{
		{
			Id:   "SECTION_ONE",
			Name: "Section One",
		},
		{
			Id:   "SECTION_TWO",
			Name: "Section Two",
		},
	})

	assert.Nil(t, rErr)

	assert.Len(t, updatedSections, 2)

	section := updatedSections[0]
	assert.Equal(t, "SECTION_ONE", section.Id)
	assert.Equal(t, "Section One", section.Name)
	assert.Equal(t, testutil.ModifyOnlyUserId(), section.UserId)
	assert.Equal(t, section.CreatedAt, testutil.Date())
	assert.Less(t, now.Sub(section.UpdatedAt), time.Second)

	section = updatedSections[1]
	assert.Equal(t, "SECTION_TWO", section.Id)
	assert.Equal(t, "Section Two", section.Name)
	assert.Equal(t, testutil.ModifyOnlyUserId(), section.UserId)
	assert.Equal(t, section.CreatedAt, testutil.Date())
	assert.Less(t, now.Sub(section.UpdatedAt), time.Second)

	updatedChapter, err = r.FetchChapter(testutil.ModifyOnlyUserId(), projectId, chapterId)

	assert.Nil(t, err)

	assert.Equal(t, record.ChapterEntry{
		Name:   "Chapter One",
		Number: 1,
		Sections: []record.SectionEntry{
			{
				Id:        "SECTION_ONE",
				Name:      "Section One",
				UserId:    testutil.ModifyOnlyUserId(),
				CreatedAt: testutil.Date(),
				UpdatedAt: updatedChapter.UpdatedAt,
			},
			{
				Id:        "SECTION_TWO",
				Name:      "Section Two",
				UserId:    testutil.ModifyOnlyUserId(),
				CreatedAt: testutil.Date(),
				UpdatedAt: updatedChapter.UpdatedAt,
			},
		},
		UserId:    testutil.ModifyOnlyUserId(),
		CreatedAt: testutil.Date(),
		UpdatedAt: updatedChapter.UpdatedAt,
	}, *updatedChapter)
	assert.Less(t, now.Sub(updatedChapter.UpdatedAt), time.Second)
}

func TestUpdateChaptersNotFound(t *testing.T) {
	tt := []struct {
		name          string
		userId        string
		projectId     string
		chapterId     string
		expectedError string
	}{
		{
			name:          "should return error when project is not found",
			userId:        testutil.ModifyOnlyUserId(),
			projectId:     "UNKNOWN_PROJECT",
			chapterId:     "CHAPTER_ONE",
			expectedError: "failed to fetch project",
		},
		{
			name:          "should return not found when user is not author of the project",
			userId:        testutil.ReadOnlyUserId(),
			projectId:     "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_REPOSITORY",
			chapterId:     "CHAPTER_ONE",
			expectedError: "failed to fetch project",
		},
		{
			name:          "should return error when chapter is not found",
			userId:        testutil.ModifyOnlyUserId(),
			projectId:     "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_REPOSITORY",
			chapterId:     "UNKNOWN_CHAPTER",
			expectedError: "failed to update sections of chapter",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			client := db.FirestoreClient()
			r := repository.NewChapterRepository(*client)

			updatedChapter, rErr := r.UpdateChapterSections(
				tc.userId,
				tc.projectId,
				tc.chapterId,
				[]record.SectionWithoutAutofieldEntry{},
			)

			assert.NotNil(t, rErr)

			assert.Nil(t, updatedChapter)
			assert.Equal(t, repository.NotFoundError, rErr.Code())
			assert.Equal(t, fmt.Sprintf("not found: %v", tc.expectedError), rErr.Error())
		})
	}
}

func TestUpdateChapterSectionsInvalidDocument(t *testing.T) {
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
			expectedError: "failed to convert snapshot to values: document.ProjectValues.chapterIds: " +
				"firestore: cannot set type []string to string",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			client := db.FirestoreClient()
			r := repository.NewChapterRepository(*client)

			updatedChapter, rErr := r.UpdateChapterSections(
				tc.userId,
				tc.projectId,
				tc.chapterId,
				[]record.SectionWithoutAutofieldEntry{},
			)

			assert.NotNil(t, rErr)

			assert.Nil(t, updatedChapter)
			assert.Equal(t, repository.ReadFailurePanic, rErr.Code())
			assert.Equal(t, fmt.Sprintf("read failure: %v", tc.expectedError), rErr.Error())
		})
	}
}
