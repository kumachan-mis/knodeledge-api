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
				NextId:    "1000000000000004",
				UserId:    testutil.ReadOnlyUserId(),
				CreatedAt: testutil.Date().Add(-3 * time.Hour),
				UpdatedAt: testutil.Date().Add(-3 * time.Hour),
			},
			"1000000000000001": {
				Name:      "Chapter 1",
				NextId:    "1000000000000002",
				UserId:    testutil.ReadOnlyUserId(),
				CreatedAt: testutil.Date().Add(-2 * time.Hour),
				UpdatedAt: testutil.Date().Add(-2 * time.Hour),
			},
			"1000000000000004": {
				Name:      maxLengthChapterName,
				NextId:    "",
				UserId:    testutil.ReadOnlyUserId(),
				CreatedAt: testutil.Date().Add(-4 * time.Hour),
				UpdatedAt: testutil.Date().Add(-4 * time.Hour),
			},
			"1000000000000002": {
				Name:      "Chapter 2",
				NextId:    "1000000000000003",
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
	assert.Equal(t, "1000000000000002", chapter.NextId().Value())
	assert.Equal(t, testutil.Date().Add(-2*time.Hour), chapter.CreatedAt().Value())
	assert.Equal(t, testutil.Date().Add(-2*time.Hour), chapter.UpdatedAt().Value())

	chapter = chapters[1]
	assert.Equal(t, "1000000000000002", chapter.Id().Value())
	assert.Equal(t, "Chapter 2", chapter.Name().Value())
	assert.Equal(t, "1000000000000003", chapter.NextId().Value())
	assert.Equal(t, testutil.Date().Add(-1*time.Hour), chapter.CreatedAt().Value())
	assert.Equal(t, testutil.Date().Add(-1*time.Hour), chapter.UpdatedAt().Value())

	chapter = chapters[2]
	assert.Equal(t, "1000000000000003", chapter.Id().Value())
	assert.Equal(t, "Chapter 3", chapter.Name().Value())
	assert.Equal(t, "1000000000000004", chapter.NextId().Value())
	assert.Equal(t, testutil.Date().Add(-3*time.Hour), chapter.CreatedAt().Value())
	assert.Equal(t, testutil.Date().Add(-3*time.Hour), chapter.UpdatedAt().Value())

	chapter = chapters[3]
	assert.Equal(t, "1000000000000004", chapter.Id().Value())
	assert.Equal(t, maxLengthChapterName, chapter.Name().Value())
	assert.Equal(t, "", chapter.NextId().Value())
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
				NextId:    "",
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
				NextId:    "",
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
				NextId:    "",
				UserId:    testutil.ReadOnlyUserId(),
				CreatedAt: testutil.Date(),
				UpdatedAt: testutil.Date(),
			},
			expectedError: "failed to convert entry to entity (name): " + fmt.Sprintf(
				"chapter name cannot be longer than 100 characters, but got '%s'",
				tooLongChapterName,
			),
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

func TestCreateChapterValidEntry(t *testing.T) {
	maxLengthChapterName := testutil.RandomString(100)

	tt := []struct {
		name    string
		chapter record.ChapterWithoutAutofieldEntry
	}{
		{
			name: "should return chapter with valid entry",
			chapter: record.ChapterWithoutAutofieldEntry{
				Name:   "Chapter One",
				NextId: "",
				UserId: testutil.ModifyOnlyUserId(),
			},
		},
		{
			name: "should return chapter with max-length valid entry",
			chapter: record.ChapterWithoutAutofieldEntry{
				Name:   maxLengthChapterName,
				NextId: "",
				UserId: testutil.ModifyOnlyUserId(),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := mock_repository.NewMockChapterRepository(ctrl)
			r.EXPECT().
				InsertChapter("0000000000000001", tc.chapter).
				Return("1000000000000001", &record.ChapterEntry{
					Name:      tc.chapter.Name,
					NextId:    tc.chapter.NextId,
					UserId:    tc.chapter.UserId,
					CreatedAt: testutil.Date(),
					UpdatedAt: testutil.Date(),
				}, nil)

			s := service.NewChapterService(r)

			userId, err := domain.NewUserIdObject(tc.chapter.UserId)
			assert.Nil(t, err)
			projectId, err := domain.NewProjectIdObject("0000000000000001")
			assert.Nil(t, err)

			name, err := domain.NewChapterNameObject(tc.chapter.Name)
			assert.Nil(t, err)
			nextId, err := domain.NewChapterNextIdObject(tc.chapter.NextId)
			assert.Nil(t, err)

			chapter := domain.NewChapterWithoutAutofieldEntity(*name, *nextId)

			createdChapter, sErr := s.CreateChapter(*userId, *projectId, *chapter)
			assert.Nil(t, sErr)

			assert.Equal(t, "1000000000000001", createdChapter.Id().Value())
			assert.Equal(t, tc.chapter.Name, createdChapter.Name().Value())
			assert.Equal(t, tc.chapter.NextId, createdChapter.NextId().Value())
			assert.Equal(t, testutil.Date(), createdChapter.CreatedAt().Value())
			assert.Equal(t, testutil.Date(), createdChapter.UpdatedAt().Value())
		})
	}
}

func TestCreateChapterInvalidCreatedEntry(t *testing.T) {
	tooLongChapterName := testutil.RandomString(101)

	tt := []struct {
		name           string
		createdChapter record.ChapterEntry
		expectedError  string
	}{
		{
			name: "should return error when chapter name is empty",
			createdChapter: record.ChapterEntry{
				Name:      "",
				NextId:    "",
				UserId:    testutil.ModifyOnlyUserId(),
				CreatedAt: testutil.Date(),
				UpdatedAt: testutil.Date(),
			},
			expectedError: "failed to convert entry to entity (name): chapter name is required, but got ''",
		},
		{
			name: "should return error when chapter name is too long",
			createdChapter: record.ChapterEntry{
				Name:      tooLongChapterName,
				NextId:    "",
				UserId:    testutil.ModifyOnlyUserId(),
				CreatedAt: testutil.Date(),
				UpdatedAt: testutil.Date(),
			},
			expectedError: "failed to convert entry to entity (name): " + fmt.Sprintf(
				"chapter name cannot be longer than 100 characters, but got '%s'",
				tooLongChapterName,
			),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := mock_repository.NewMockChapterRepository(ctrl)
			r.EXPECT().
				InsertChapter("0000000000000001", record.ChapterWithoutAutofieldEntry{
					Name:   "Chapter One",
					NextId: "",
					UserId: testutil.ModifyOnlyUserId(),
				}).
				Return("1000000000000001", &tc.createdChapter, nil)

			s := service.NewChapterService(r)

			userId, err := domain.NewUserIdObject(testutil.ModifyOnlyUserId())
			assert.Nil(t, err)
			projectId, err := domain.NewProjectIdObject("0000000000000001")
			assert.Nil(t, err)

			name, err := domain.NewChapterNameObject("Chapter One")
			assert.Nil(t, err)
			nextId, err := domain.NewChapterNextIdObject("")
			assert.Nil(t, err)

			chapter := domain.NewChapterWithoutAutofieldEntity(*name, *nextId)

			createdChapter, sErr := s.CreateChapter(*userId, *projectId, *chapter)
			assert.NotNil(t, sErr)
			assert.Equal(t, service.DomainFailurePanic, sErr.Code())

			assert.Equal(t, fmt.Sprintf("domain failure: %s", tc.expectedError), sErr.Error())
			assert.Nil(t, createdChapter)
		})
	}
}

func TestCreateChapterInvalidArgumentError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := mock_repository.NewMockChapterRepository(ctrl)
	r.EXPECT().
		InsertChapter("0000000000000001", record.ChapterWithoutAutofieldEntry{
			Name:   "Chapter One",
			NextId: "UNKNOWN_CHAPTER",
			UserId: testutil.ModifyOnlyUserId(),
		}).
		Return("", nil, repository.Errorf(repository.InvalidArgument, "id of next chapter does not exist"))

	s := service.NewChapterService(r)

	userId, err := domain.NewUserIdObject(testutil.ModifyOnlyUserId())
	assert.Nil(t, err)
	projectId, err := domain.NewProjectIdObject("0000000000000001")
	assert.Nil(t, err)

	name, err := domain.NewChapterNameObject("Chapter One")
	assert.Nil(t, err)
	nextId, err := domain.NewChapterNextIdObject("UNKNOWN_CHAPTER")
	assert.Nil(t, err)

	chapter := domain.NewChapterWithoutAutofieldEntity(*name, *nextId)

	createdChapter, sErr := s.CreateChapter(*userId, *projectId, *chapter)
	assert.NotNil(t, sErr)
	assert.Equal(t, service.InvalidArgument, sErr.Code())
	assert.Equal(t, "invalid argument: failed to create chapter: id of next chapter does not exist", sErr.Error())
	assert.Nil(t, createdChapter)
}

func TestCreateChapterRepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := mock_repository.NewMockChapterRepository(ctrl)
	r.EXPECT().
		InsertChapter("0000000000000001", gomock.Any()).
		Return("", nil, repository.Errorf(repository.WriteFailurePanic, "repository error"))

	s := service.NewChapterService(r)

	userId, err := domain.NewUserIdObject(testutil.ModifyOnlyUserId())
	assert.Nil(t, err)
	projectId, err := domain.NewProjectIdObject("0000000000000001")
	assert.Nil(t, err)

	name, err := domain.NewChapterNameObject("Chapter One")
	assert.Nil(t, err)
	nextId, err := domain.NewChapterNextIdObject("")
	assert.Nil(t, err)

	chapter := domain.NewChapterWithoutAutofieldEntity(*name, *nextId)

	createdChapter, sErr := s.CreateChapter(*userId, *projectId, *chapter)
	assert.NotNil(t, sErr)
	assert.Equal(t, service.RepositoryFailurePanic, sErr.Code())
	assert.Equal(t, "repository failure: failed to create chapter: repository error", sErr.Error())
	assert.Nil(t, createdChapter)
}
