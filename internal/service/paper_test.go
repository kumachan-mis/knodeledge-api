package service_test

import (
	"fmt"
	"testing"

	"github.com/kumachan-mis/knodeledge-api/internal/domain"
	"github.com/kumachan-mis/knodeledge-api/internal/record"
	"github.com/kumachan-mis/knodeledge-api/internal/repository"
	"github.com/kumachan-mis/knodeledge-api/internal/service"
	mock_repository "github.com/kumachan-mis/knodeledge-api/mock/repository"
	"github.com/kumachan-mis/knodeledge-api/test/testutil"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestFindPaperValidEntry(t *testing.T) {
	maxLengthPaperContent := testutil.RandomString(40000)

	tt := []struct {
		name  string
		entry record.PaperEntry
	}{
		{
			name: "should return paper with valid entry",
			entry: record.PaperEntry{
				Content:   "valid content",
				UserId:    testutil.ReadOnlyUserId(),
				CreatedAt: testutil.Date(),
				UpdatedAt: testutil.Date(),
			},
		},
		{
			name: "should return paper with max-length valid entry",
			entry: record.PaperEntry{
				Content:   maxLengthPaperContent,
				UserId:    testutil.ReadOnlyUserId(),
				CreatedAt: testutil.Date(),
				UpdatedAt: testutil.Date(),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := mock_repository.NewMockPaperRepository(ctrl)
			r.EXPECT().
				FetchPaper("USER", "PROJECT", "CHAPTER").
				Return(&tc.entry, nil)

			s := service.NewPaperService(r)

			userId, err := domain.NewUserIdObject("USER")
			assert.Nil(t, err)
			projectId, err := domain.NewProjectIdObject("PROJECT")
			assert.Nil(t, err)
			chapterId, err := domain.NewChapterIdObject("CHAPTER")
			assert.Nil(t, err)

			paper, sErr := s.FindPaper(*userId, *projectId, *chapterId)
			assert.Nil(t, sErr)

			assert.Equal(t, tc.entry.Content, paper.Content().Value())
			assert.Equal(t, tc.entry.CreatedAt, paper.CreatedAt().Value())
			assert.Equal(t, tc.entry.UpdatedAt, paper.UpdatedAt().Value())
		})
	}
}

func TestFindPaperInvalidEntry(t *testing.T) {
	tooLongPaperContent := testutil.RandomString(40001)

	tt := []struct {
		name          string
		entry         record.PaperEntry
		expectedError string
	}{
		{
			name: "should return error when content is too long",
			entry: record.PaperEntry{
				Content:   tooLongPaperContent,
				UserId:    testutil.ReadOnlyUserId(),
				CreatedAt: testutil.Date(),
				UpdatedAt: testutil.Date(),
			},
			expectedError: "failed to convert entry to entity (content): " +
				"paper content must be less than or equal to 40000 bytes, but got 40001 bytes",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := mock_repository.NewMockPaperRepository(ctrl)
			r.EXPECT().
				FetchPaper("USER", "PROJECT", "CHAPTER").
				Return(&tc.entry, nil)

			s := service.NewPaperService(r)

			userId, err := domain.NewUserIdObject("USER")
			assert.Nil(t, err)
			projectId, err := domain.NewProjectIdObject("PROJECT")
			assert.Nil(t, err)
			chapterId, err := domain.NewChapterIdObject("CHAPTER")
			assert.Nil(t, err)

			paper, sErr := s.FindPaper(*userId, *projectId, *chapterId)
			assert.NotNil(t, sErr)
			assert.Equal(t, service.DomainFailurePanic, sErr.Code())
			assert.Equal(t, fmt.Sprintf("domain failure: %v", tc.expectedError), sErr.Error())
			assert.Nil(t, paper)
		})
	}
}

func TestFindPaperRepositoryError(t *testing.T) {
	tt := []struct {
		name          string
		errorCode     repository.ErrorCode
		errorMessage  string
		expectedError string
		expectedCode  service.ErrorCode
	}{
		{
			name:          "should return error when repository returns not found error",
			errorCode:     repository.NotFoundError,
			errorMessage:  "paper not found",
			expectedError: "failed to find paper: paper not found",
			expectedCode:  service.NotFoundError,
		},
		{
			name:          "should return error when repository returns read failure error",
			errorCode:     repository.ReadFailurePanic,
			errorMessage:  "repository error",
			expectedError: "failed to fetch paper: repository error",
			expectedCode:  service.RepositoryFailurePanic,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := mock_repository.NewMockPaperRepository(ctrl)
			r.EXPECT().
				FetchPaper("USER", "PROJECT", "CHAPTER").
				Return(nil, repository.Errorf(tc.errorCode, tc.errorMessage))

			s := service.NewPaperService(r)

			userId, err := domain.NewUserIdObject("USER")
			assert.Nil(t, err)
			projectId, err := domain.NewProjectIdObject("PROJECT")
			assert.Nil(t, err)
			chapterId, err := domain.NewChapterIdObject("CHAPTER")
			assert.Nil(t, err)

			paper, sErr := s.FindPaper(*userId, *projectId, *chapterId)
			assert.NotNil(t, sErr)
			assert.Equal(t, tc.expectedCode, sErr.Code())
			assert.Equal(t, fmt.Sprintf("%v: %v", tc.expectedCode, tc.expectedError), sErr.Error())
			assert.Nil(t, paper)
		})
	}
}

func TestCreatePaperValidEntry(t *testing.T) {
	maxLengthPaperContent := testutil.RandomString(40000)

	tt := []struct {
		name  string
		paper record.PaperWithoutAutofieldEntry
	}{
		{
			name: "should return paper with valid entry",
			paper: record.PaperWithoutAutofieldEntry{
				Content: "valid content",
				UserId:  testutil.ModifyOnlyUserId(),
			},
		},
		{

			name: "should return paper with max-length valid entry",
			paper: record.PaperWithoutAutofieldEntry{
				Content: maxLengthPaperContent,
				UserId:  testutil.ModifyOnlyUserId(),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := mock_repository.NewMockPaperRepository(ctrl)
			r.EXPECT().
				InsertPaper("PROJECT", "CHAPTER", tc.paper).
				Return("CHAPTER", &record.PaperEntry{
					Content:   tc.paper.Content,
					UserId:    tc.paper.UserId,
					CreatedAt: testutil.Date(),
					UpdatedAt: testutil.Date(),
				}, nil)

			s := service.NewPaperService(r)

			userId, err := domain.NewUserIdObject(tc.paper.UserId)
			assert.Nil(t, err)
			projectId, err := domain.NewProjectIdObject("PROJECT")
			assert.Nil(t, err)
			chapterId, err := domain.NewChapterIdObject("CHAPTER")
			assert.Nil(t, err)

			content, err := domain.NewPaperContentObject(tc.paper.Content)
			assert.Nil(t, err)

			paper := domain.NewPaperWithoutAutofieldEntity(*content)

			createdPaper, sErr := s.CreatePaper(*userId, *projectId, *chapterId, *paper)
			assert.Nil(t, sErr)

			assert.Equal(t, "CHAPTER", createdPaper.Id().Value())
			assert.Equal(t, tc.paper.Content, createdPaper.Content().Value())
			assert.Equal(t, testutil.Date(), createdPaper.CreatedAt().Value())
			assert.Equal(t, testutil.Date(), createdPaper.UpdatedAt().Value())
		})
	}
}

func TestCreatePaperInvalidCreatedEntry(t *testing.T) {
	tooLongPaperContent := testutil.RandomString(40001)

	tt := []struct {
		name          string
		createdPaper  record.PaperEntry
		expectedError string
	}{
		{
			name: "should return error when content is too long",
			createdPaper: record.PaperEntry{
				Content:   tooLongPaperContent,
				UserId:    testutil.ModifyOnlyUserId(),
				CreatedAt: testutil.Date(),
				UpdatedAt: testutil.Date(),
			},
			expectedError: "failed to convert entry to entity (content): " +
				"paper content must be less than or equal to 40000 bytes, but got 40001 bytes",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := mock_repository.NewMockPaperRepository(ctrl)
			r.EXPECT().
				InsertPaper("PROJECT", "CHAPTER", gomock.Any()).
				Return("CHAPTER", &tc.createdPaper, nil)

			s := service.NewPaperService(r)

			userId, err := domain.NewUserIdObject(tc.createdPaper.UserId)
			assert.Nil(t, err)
			projectId, err := domain.NewProjectIdObject("PROJECT")
			assert.Nil(t, err)
			chapterId, err := domain.NewChapterIdObject("CHAPTER")
			assert.Nil(t, err)

			content, err := domain.NewPaperContentObject("content")
			assert.Nil(t, err)

			paper := domain.NewPaperWithoutAutofieldEntity(*content)

			createdPaper, sErr := s.CreatePaper(*userId, *projectId, *chapterId, *paper)
			assert.NotNil(t, sErr)
			assert.Equal(t, service.DomainFailurePanic, sErr.Code())
			assert.Equal(t, fmt.Sprintf("domain failure: %v", tc.expectedError), sErr.Error())
			assert.Nil(t, createdPaper)
		})
	}
}

func TestCreatePaperRepositoryError(t *testing.T) {
	tt := []struct {
		name          string
		errorCode     repository.ErrorCode
		errorMessage  string
		expectedError string
		expectedCode  service.ErrorCode
	}{
		{
			name:          "should return error when repository returns not found error",
			errorCode:     repository.NotFoundError,
			errorMessage:  "project not found",
			expectedError: "failed to create paper: project not found",
			expectedCode:  service.NotFoundError,
		},
		{
			name:          "should return error when repository returns write failure error",
			errorCode:     repository.WriteFailurePanic,
			errorMessage:  "repository error",
			expectedError: "failed to insert paper: repository error",
			expectedCode:  service.RepositoryFailurePanic,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := mock_repository.NewMockPaperRepository(ctrl)
			r.EXPECT().
				InsertPaper("PROJECT", "CHAPTER", gomock.Any()).
				Return("", nil, repository.Errorf(tc.errorCode, tc.errorMessage))

			s := service.NewPaperService(r)

			userId, err := domain.NewUserIdObject(testutil.ModifyOnlyUserId())
			assert.Nil(t, err)
			projectId, err := domain.NewProjectIdObject("PROJECT")
			assert.Nil(t, err)
			chapterId, err := domain.NewChapterIdObject("CHAPTER")
			assert.Nil(t, err)

			content, err := domain.NewPaperContentObject("content")
			assert.Nil(t, err)

			paper := domain.NewPaperWithoutAutofieldEntity(*content)

			createdPaper, sErr := s.CreatePaper(*userId, *projectId, *chapterId, *paper)
			assert.NotNil(t, sErr)
			assert.Equal(t, tc.expectedCode, sErr.Code())
			assert.Equal(t, fmt.Sprintf("%v: %v", tc.expectedCode, tc.expectedError), sErr.Error())
			assert.Nil(t, createdPaper)
		})
	}
}
