package repository_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/kumachan-mis/knodeledge-api/internal/db"
	"github.com/kumachan-mis/knodeledge-api/internal/record"
	"github.com/kumachan-mis/knodeledge-api/internal/repository"
	"github.com/kumachan-mis/knodeledge-api/internal/testutil"
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

func TestGraphExistsNotFound(t *testing.T) {
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

func TestFetchGraphValidEntry(t *testing.T) {

	userId := testutil.ReadOnlyUserId()
	projectId := "PROJECT_WITHOUT_DESCRIPTION"
	chapterId := "CHAPTER_ONE"
	sectionId := "SECTION_ONE"

	client := db.FirestoreClient()
	r := repository.NewGraphRepository(*client)

	entry, rErr := r.FetchGraph(userId, projectId, chapterId, sectionId)

	assert.Nil(t, rErr)

	assert.Equal(t, record.GraphEntry{
		Name:      "Introduction",
		Paragraph: "This is an example project of kNODEledge.",
		Children: []record.GraphChildEntry{
			{
				Name:        "Background",
				Relation:    "part of",
				Description: "This is background part.",
				Children: []record.GraphChildEntry{
					{
						Name:        "IT in Education",
						Relation:    "one of",
						Description: "This is IT in Education part.",
						Children:    []record.GraphChildEntry{},
					},
				},
			},
			{
				Name:        "Motivation",
				Relation:    "part of",
				Description: "This is motivation part.",
				Children:    []record.GraphChildEntry{},
			},
			{
				Name:        "Literature Review",
				Relation:    "part of",
				Description: "This is literature review part.",
				Children:    []record.GraphChildEntry{},
			},
		},
		UserId:    testutil.ReadOnlyUserId(),
		CreatedAt: testutil.Date(),
		UpdatedAt: testutil.Date(),
	}, *entry)
}

func TestFetchGraphNotFound(t *testing.T) {
	tt := []struct {
		name          string
		userId        string
		projectId     string
		chapterId     string
		sectionId     string
		expectedError string
	}{
		{
			name:          "should return error when project not found",
			userId:        testutil.ReadOnlyUserId(),
			projectId:     "UNKNOWN_PROJECT",
			chapterId:     "CHAPTER_ONE",
			sectionId:     "SECTION_ONE",
			expectedError: "failed to fetch project",
		},
		{
			name:          "should return not found when user is not author of the project",
			userId:        testutil.ModifyOnlyUserId(),
			projectId:     "PROJECT_WITHOUT_DESCRIPTION",
			chapterId:     "CHAPTER_ONE",
			sectionId:     "SECTION_ONE",
			expectedError: "failed to fetch project",
		},
		{
			name:          "should return error when chapter not found",
			userId:        testutil.ReadOnlyUserId(),
			projectId:     "PROJECT_WITHOUT_DESCRIPTION",
			chapterId:     "UNKNOWN_CHAPTER",
			sectionId:     "SECTION_ONE",
			expectedError: "failed to fetch chapter",
		},
		{
			name:          "should return error when section not found",
			userId:        testutil.ReadOnlyUserId(),
			projectId:     "PROJECT_WITHOUT_DESCRIPTION",
			chapterId:     "CHAPTER_ONE",
			sectionId:     "UNKNOWN_SECTION",
			expectedError: "failed to fetch graph",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			client := db.FirestoreClient()
			r := repository.NewGraphRepository(*client)

			entry, rErr := r.FetchGraph(tc.userId, tc.projectId, tc.chapterId, tc.sectionId)

			assert.NotNil(t, rErr)

			assert.Nil(t, entry)
			assert.Equal(t, repository.NotFoundError, rErr.Code())
			assert.Equal(t, fmt.Sprintf("not found: %s", tc.expectedError), rErr.Error())
		})
	}
}

func TestFetchGraphInvalidDocument(t *testing.T) {
	tt := []struct {
		name          string
		userId        string
		projectId     string
		chapterId     string
		sectionId     string
		expectedError string
	}{
		{
			name:      "should return error when graph paragraph is invalid",
			userId:    testutil.ErrorUserId(7),
			projectId: "PROJECT_WITH_INVALID_GRAPH_PARAGRAPH",
			chapterId: "CHAPTER_WITH_INVALID_GRAPH_PARAGRAPH",
			sectionId: "SECTION_WITH_INVALID_GRAPH_PARAGRAPH",
			expectedError: "failed to convert snapshot to values: document.GraphValues.paragraph: " +
				"firestore: cannot set type string to array",
		},
		{
			name:          "should return error when sections have excessive elements",
			userId:        testutil.ErrorUserId(8),
			projectId:     "PROJECT_WITH_TOO_MANY_SECTIONS",
			chapterId:     "CHAPTER_WITH_TOO_MANY_SECTIONS",
			sectionId:     "UNKNOWN_SECTION",
			expectedError: "failed to convert values to entry: document.ChapterValues.sections have excessive elements",
		},
		{
			name:          "should return error when sections have insufficient elements",
			userId:        testutil.ErrorUserId(9),
			projectId:     "PROJECT_WITH_TOO_FEW_SECTIONS",
			chapterId:     "CHAPTER_WITH_TOO_FEW_SECTIONS",
			sectionId:     "SECTION_UNKNOWN",
			expectedError: "failed to convert values to entry: document.ChapterValues.sections have insufficient elements",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			client := db.FirestoreClient()
			r := repository.NewGraphRepository(*client)

			entry, rErr := r.FetchGraph(tc.userId, tc.projectId, tc.chapterId, tc.sectionId)

			assert.NotNil(t, rErr)

			assert.Nil(t, entry)
			assert.Equal(t, repository.ReadFailurePanic, rErr.Code())
			assert.Equal(t, fmt.Sprintf("read failure: %v", tc.expectedError), rErr.Error())
		})
	}
}

func TestInsertGraphsValidEntry(t *testing.T) {
	client := db.FirestoreClient()
	r := repository.NewGraphRepository(*client)

	userId := testutil.ModifyOnlyUserId()
	projectId := "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_REPOSITORY"
	chapterId := "CHAPTER_TWO"

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
			Children:  []record.GraphChildEntry{},
		},
		{
			Name:      "What is note apps?",
			Paragraph: paragraph2,
			Children:  []record.GraphChildEntry{},
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
	assert.Equal(t, []record.GraphChildEntry{}, createdEntry.Children)
	assert.Equal(t, testutil.ModifyOnlyUserId(), createdEntry.UserId)
	assert.Less(t, now.Sub(createdEntry.CreatedAt), time.Second)
	assert.Less(t, now.Sub(createdEntry.UpdatedAt), time.Second)

	id = ids[1]
	createdEntry = createdEntries[1]
	assert.NotEmpty(t, id)
	assert.Equal(t, "What is note apps?", createdEntry.Name)
	assert.Equal(t, paragraph2, createdEntry.Paragraph)
	assert.Equal(t, []record.GraphChildEntry{}, createdEntry.Children)
	assert.Equal(t, testutil.ModifyOnlyUserId(), createdEntry.UserId)
	assert.Less(t, now.Sub(createdEntry.CreatedAt), time.Second)
	assert.Less(t, now.Sub(createdEntry.UpdatedAt), time.Second)
}

func TestInsertGraphNotFound(t *testing.T) {
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

			id, createdPaper, rErr := r.InsertGraphs(tc.userId, tc.projectId, tc.chapterId,
				[]record.GraphWithoutAutofieldEntry{
					{
						Name:      "Section Name",
						Paragraph: "paragraph",
						Children:  []record.GraphChildEntry{},
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

func TestUpdateGraphContentValidEntry(t *testing.T) {
	userId := testutil.ModifyOnlyUserId()
	projectId := "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_REPOSITORY"
	chapterId := "CHAPTER_ONE"
	sectionId := "SECTION_ONE"

	client := db.FirestoreClient()
	r := repository.NewGraphRepository(*client)

	paragraph := "This is the introduction of the paper."
	children := []record.GraphChildEntry{
		{
			Name:        "Background",
			Relation:    "part of",
			Description: "This is background part.",
			Children: []record.GraphChildEntry{
				{
					Name:        "IT in Education",
					Relation:    "one of",
					Description: "This is IT in Education part.",
					Children:    []record.GraphChildEntry{},
				},
				{
					Name:        "IT in Business",
					Relation:    "one of",
					Description: "This is IT in Business part.",
					Children:    []record.GraphChildEntry{},
				},
			},
		},
		{
			Name:        "Motivation",
			Relation:    "one of",
			Description: "This is motivation part.",
			Children:    []record.GraphChildEntry{},
		},
		{
			Name:        "Lterature Review",
			Relation:    "one of",
			Description: "This is literature review part.",
			Children:    []record.GraphChildEntry{},
		},
	}

	updatedEntry, rErr := r.UpdateGraphContent(userId, projectId, chapterId, sectionId,
		record.GraphContentEntry{Paragraph: paragraph, Children: children})
	now := time.Now()

	assert.Nil(t, rErr)

	assert.Equal(t, "Introduction", updatedEntry.Name)
	assert.Equal(t, paragraph, updatedEntry.Paragraph)
	assert.Equal(t, children, updatedEntry.Children)
	assert.Equal(t, testutil.ModifyOnlyUserId(), updatedEntry.UserId)
	assert.Equal(t, testutil.Date(), updatedEntry.CreatedAt)
	assert.Less(t, now.Sub(updatedEntry.UpdatedAt), time.Second)

}

func TestUpdateGraphContentNotFound(t *testing.T) {
	tt := []struct {
		name          string
		userId        string
		projectId     string
		chapterId     string
		sectionId     string
		expectedError string
	}{
		{
			name:          "should return error when project not found",
			userId:        testutil.ModifyOnlyUserId(),
			projectId:     "UNKNOWN_PROJECT",
			chapterId:     "CHAPTER_ONE",
			sectionId:     "SECTION_ONE",
			expectedError: "failed to fetch project",
		},
		{
			name:          "should return not found when user is not author of the project",
			userId:        testutil.ReadOnlyUserId(),
			projectId:     "PROJECT_WITH_DESCRIPTION_TO_UPDATE_FROM_REPOSITORY",
			chapterId:     "CHAPTER_ONE",
			sectionId:     "SECTION_ONE",
			expectedError: "failed to fetch project",
		},
		{
			name:          "should return error when chapter not found",
			userId:        testutil.ModifyOnlyUserId(),
			projectId:     "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_REPOSITORY",
			chapterId:     "UNKNOWN_CHAPTER",
			sectionId:     "SECTION_ONE",
			expectedError: "failed to fetch chapter",
		},
		{
			name:          "should return error when section not found",
			userId:        testutil.ModifyOnlyUserId(),
			projectId:     "PROJECT_WITHOUT_DESCRIPTION_TO_UPDATE_FROM_REPOSITORY",
			chapterId:     "CHAPTER_ONE",
			sectionId:     "UNKNOWN_SECTION",
			expectedError: "failed to fetch graph",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			client := db.FirestoreClient()
			r := repository.NewGraphRepository(*client)

			entry, rErr := r.UpdateGraphContent(tc.userId, tc.projectId, tc.chapterId, tc.sectionId, record.GraphContentEntry{
				Paragraph: "content",
			})

			assert.NotNil(t, rErr)

			assert.Nil(t, entry)
			assert.Equal(t, repository.NotFoundError, rErr.Code())
			assert.Equal(t, fmt.Sprintf("not found: %s", tc.expectedError), rErr.Error())
		})
	}
}
