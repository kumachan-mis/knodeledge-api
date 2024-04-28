package service_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/kumachan-mis/knodeledge-api/internal/domain"
	"github.com/kumachan-mis/knodeledge-api/internal/record"
	"github.com/kumachan-mis/knodeledge-api/internal/repository"
	"github.com/kumachan-mis/knodeledge-api/internal/service"
	mock_repository "github.com/kumachan-mis/knodeledge-api/mock/repository"
	"github.com/kumachan-mis/knodeledge-api/test/testutil"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestListChaptersValidEntry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	maxLengthChapterName := testutil.RandomString(100)

	r := mock_repository.NewMockChapterRepository(ctrl)
	r.EXPECT().
		FetchProjectChapters(testutil.ReadOnlyUserId(), "0000000000000001").
		Return(map[string]record.ChapterEntry{
			"1000000000000003": {
				Name:      "Chapter 3",
				Number:    3,
				UserId:    testutil.ReadOnlyUserId(),
				CreatedAt: testutil.Date().Add(-3 * time.Hour),
				UpdatedAt: testutil.Date().Add(-3 * time.Hour),
			},
			"1000000000000001": {
				Name:      "Chapter 1",
				Number:    1,
				UserId:    testutil.ReadOnlyUserId(),
				CreatedAt: testutil.Date().Add(-2 * time.Hour),
				UpdatedAt: testutil.Date().Add(-2 * time.Hour),
			},
			"1000000000000004": {
				Name:      maxLengthChapterName,
				Number:    4,
				UserId:    testutil.ReadOnlyUserId(),
				CreatedAt: testutil.Date().Add(-4 * time.Hour),
				UpdatedAt: testutil.Date().Add(-4 * time.Hour),
			},
			"1000000000000002": {
				Name:      "Chapter 2",
				Number:    2,
				UserId:    testutil.ReadOnlyUserId(),
				CreatedAt: testutil.Date().Add(-1 * time.Hour),
				UpdatedAt: testutil.Date().Add(-1 * time.Hour),
			},
		}, nil)

	s := service.NewChapterService(r)

	userId, err := domain.NewUserIdObject(testutil.ReadOnlyUserId())
	assert.Nil(t, err)

	projectId, err := domain.NewProjectIdObject("0000000000000001")
	assert.Nil(t, err)

	chapters, err := s.ListChapters(*userId, *projectId)
	assert.Nil(t, err)

	assert.Len(t, chapters, 4)

	chapter := chapters[0]
	assert.Equal(t, "1000000000000001", chapter.Id().Value())
	assert.Equal(t, "Chapter 1", chapter.Name().Value())
	assert.Equal(t, 1, chapter.Number().Value())
	assert.Equal(t, testutil.Date().Add(-2*time.Hour), chapter.CreatedAt().Value())
	assert.Equal(t, testutil.Date().Add(-2*time.Hour), chapter.UpdatedAt().Value())

	chapter = chapters[1]
	assert.Equal(t, "1000000000000002", chapter.Id().Value())
	assert.Equal(t, "Chapter 2", chapter.Name().Value())
	assert.Equal(t, 2, chapter.Number().Value())
	assert.Equal(t, testutil.Date().Add(-1*time.Hour), chapter.CreatedAt().Value())
	assert.Equal(t, testutil.Date().Add(-1*time.Hour), chapter.UpdatedAt().Value())

	chapter = chapters[2]
	assert.Equal(t, "1000000000000003", chapter.Id().Value())
	assert.Equal(t, "Chapter 3", chapter.Name().Value())
	assert.Equal(t, 3, chapter.Number().Value())
	assert.Equal(t, testutil.Date().Add(-3*time.Hour), chapter.CreatedAt().Value())
	assert.Equal(t, testutil.Date().Add(-3*time.Hour), chapter.UpdatedAt().Value())

	chapter = chapters[3]
	assert.Equal(t, "1000000000000004", chapter.Id().Value())
	assert.Equal(t, maxLengthChapterName, chapter.Name().Value())
	assert.Equal(t, 4, chapter.Number().Value())
	assert.Equal(t, testutil.Date().Add(-4*time.Hour), chapter.CreatedAt().Value())
	assert.Equal(t, testutil.Date().Add(-4*time.Hour), chapter.UpdatedAt().Value())
}

func TestListChaptersNoEntry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := mock_repository.NewMockChapterRepository(ctrl)
	r.EXPECT().
		FetchProjectChapters(testutil.ReadOnlyUserId(), "0000000000000001").
		Return(map[string]record.ChapterEntry{}, nil)

	s := service.NewChapterService(r)

	userId, err := domain.NewUserIdObject(testutil.ReadOnlyUserId())
	assert.Nil(t, err)

	projectId, err := domain.NewProjectIdObject("0000000000000001")
	assert.Nil(t, err)

	chapters, sErr := s.ListChapters(*userId, *projectId)
	assert.Nil(t, sErr)

	assert.Len(t, chapters, 0)
}

func TestListChaptersInvalidEntry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tooLongChapterName := testutil.RandomString(101)

	tt := []struct {
		name          string
		chapterId     string
		chapter       record.ChapterEntry
		expectedError string
	}{
		{
			name:      "should return error when chapter id is empty",
			chapterId: "",
			chapter: record.ChapterEntry{
				Name:      "Chapter 1",
				Number:    1,
				UserId:    testutil.ReadOnlyUserId(),
				CreatedAt: testutil.Date(),
				UpdatedAt: testutil.Date(),
			},
			expectedError: "failed to convert entry to entity (id): chapter id is required, but got ''",
		},
		{
			name:      "should return error when chapter name is empty",
			chapterId: "1000000000000001",
			chapter: record.ChapterEntry{
				Name:      "",
				Number:    1,
				UserId:    testutil.ReadOnlyUserId(),
				CreatedAt: testutil.Date(),
				UpdatedAt: testutil.Date(),
			},
			expectedError: "failed to convert entry to entity (name): chapter name is required, but got ''",
		},
		{
			name:      "should return error when chapter name is too long",
			chapterId: "1000000000000001",
			chapter: record.ChapterEntry{
				Name:      tooLongChapterName,
				Number:    1,
				UserId:    testutil.ReadOnlyUserId(),
				CreatedAt: testutil.Date(),
				UpdatedAt: testutil.Date(),
			},
			expectedError: "failed to convert entry to entity (name): " + fmt.Sprintf(
				"chapter name cannot be longer than 100 characters, but got '%s'",
				tooLongChapterName,
			),
		},
		{
			name:      "should return error when chapter number is zero",
			chapterId: "1000000000000001",
			chapter: record.ChapterEntry{
				Name:      "Chapter 1",
				Number:    0,
				UserId:    testutil.ReadOnlyUserId(),
				CreatedAt: testutil.Date(),
				UpdatedAt: testutil.Date(),
			},
			expectedError: "failed to convert entry to entity (number): chapter number must be greater than 0, but got 0",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r := mock_repository.NewMockChapterRepository(ctrl)
			r.EXPECT().
				FetchProjectChapters(testutil.ReadOnlyUserId(), "0000000000000001").
				Return(map[string]record.ChapterEntry{
					tc.chapterId: tc.chapter,
				}, nil)

			s := service.NewChapterService(r)

			userId, err := domain.NewUserIdObject(testutil.ReadOnlyUserId())
			assert.Nil(t, err)

			projectId, err := domain.NewProjectIdObject("0000000000000001")
			assert.Nil(t, err)

			chapters, sErr := s.ListChapters(*userId, *projectId)
			assert.NotNil(t, sErr)
			assert.Equal(t, service.DomainFailurePanic, sErr.Code())
			assert.Equal(t, fmt.Sprintf("domain failure: %s", tc.expectedError), sErr.Error())
			assert.Nil(t, chapters)
		})
	}
}

func TestListChaptersRepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := mock_repository.NewMockChapterRepository(ctrl)
	r.EXPECT().
		FetchProjectChapters(testutil.ReadOnlyUserId(), "0000000000000001").
		Return(nil, repository.Errorf(repository.ReadFailurePanic, "repository error"))

	s := service.NewChapterService(r)

	userId, err := domain.NewUserIdObject(testutil.ReadOnlyUserId())
	assert.Nil(t, err)

	projectId, err := domain.NewProjectIdObject("0000000000000001")
	assert.Nil(t, err)

	chapters, sErr := s.ListChapters(*userId, *projectId)
	assert.NotNil(t, sErr)
	assert.Equal(t, service.RepositoryFailurePanic, sErr.Code())
	assert.Equal(t, "repository failure: failed to fetch project chapters: repository error", sErr.Error())
	assert.Nil(t, chapters)
}
