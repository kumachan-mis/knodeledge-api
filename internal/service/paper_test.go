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
				FetchPaper(testutil.ReadOnlyUserId(), "0000000000000001", "1000000000000001").
				Return(&tc.entry, nil)

			s := service.NewPaperService(r)

			userId, err := domain.NewUserIdObject(testutil.ReadOnlyUserId())
			assert.Nil(t, err)
			projectId, err := domain.NewProjectIdObject("0000000000000001")
			assert.Nil(t, err)
			chapterId, err := domain.NewChapterIdObject("1000000000000001")
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
				FetchPaper(testutil.ReadOnlyUserId(), "0000000000000001", "1000000000000001").
				Return(&tc.entry, nil)

			s := service.NewPaperService(r)

			userId, err := domain.NewUserIdObject(testutil.ReadOnlyUserId())
			assert.Nil(t, err)
			projectId, err := domain.NewProjectIdObject("0000000000000001")
			assert.Nil(t, err)
			chapterId, err := domain.NewChapterIdObject("1000000000000001")
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
				FetchPaper(testutil.ReadOnlyUserId(), "0000000000000001", "1000000000000001").
				Return(nil, repository.Errorf(tc.errorCode, tc.errorMessage))

			s := service.NewPaperService(r)

			userId, err := domain.NewUserIdObject(testutil.ReadOnlyUserId())
			assert.Nil(t, err)
			projectId, err := domain.NewProjectIdObject("0000000000000001")
			assert.Nil(t, err)
			chapterId, err := domain.NewChapterIdObject("1000000000000001")
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
			},
		},
		{

			name: "should return paper with max-length valid entry",
			paper: record.PaperWithoutAutofieldEntry{
				Content: maxLengthPaperContent,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := mock_repository.NewMockPaperRepository(ctrl)
			r.EXPECT().
				InsertPaper(testutil.ModifyOnlyUserId(), "0000000000000001", "1000000000000001", tc.paper).
				Return("1000000000000001", &record.PaperEntry{
					Content:   tc.paper.Content,
					CreatedAt: testutil.Date(),
					UpdatedAt: testutil.Date(),
				}, nil)

			s := service.NewPaperService(r)

			userId, err := domain.NewUserIdObject(testutil.ModifyOnlyUserId())
			assert.Nil(t, err)
			projectId, err := domain.NewProjectIdObject("0000000000000001")
			assert.Nil(t, err)
			chapterId, err := domain.NewChapterIdObject("1000000000000001")
			assert.Nil(t, err)

			content, err := domain.NewPaperContentObject(tc.paper.Content)
			assert.Nil(t, err)

			paper := domain.NewPaperWithoutAutofieldEntity(*content)

			createdPaper, sErr := s.CreatePaper(*userId, *projectId, *chapterId, *paper)
			assert.Nil(t, sErr)

			assert.Equal(t, "1000000000000001", createdPaper.Id().Value())
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
				InsertPaper(testutil.ModifyOnlyUserId(), "0000000000000001", "1000000000000001", gomock.Any()).
				Return("1000000000000001", &tc.createdPaper, nil)

			s := service.NewPaperService(r)

			userId, err := domain.NewUserIdObject(tc.createdPaper.UserId)
			assert.Nil(t, err)
			projectId, err := domain.NewProjectIdObject("0000000000000001")
			assert.Nil(t, err)
			chapterId, err := domain.NewChapterIdObject("1000000000000001")
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
			errorMessage:  "failed to fetch project",
			expectedError: "failed to create paper: failed to fetch project",
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
				InsertPaper(testutil.ModifyOnlyUserId(), "0000000000000001", "1000000000000001", gomock.Any()).
				Return("", nil, repository.Errorf(tc.errorCode, tc.errorMessage))

			s := service.NewPaperService(r)

			userId, err := domain.NewUserIdObject(testutil.ModifyOnlyUserId())
			assert.Nil(t, err)
			projectId, err := domain.NewProjectIdObject("0000000000000001")
			assert.Nil(t, err)
			chapterId, err := domain.NewChapterIdObject("1000000000000001")
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

func TestUpdatePaperValidEntry(t *testing.T) {
	maxLengthPaperContent := testutil.RandomString(40000)

	tt := []struct {
		name  string
		paper record.PaperWithoutAutofieldEntry
	}{
		{
			name: "should return paper with valid entry",
			paper: record.PaperWithoutAutofieldEntry{
				Content: "valid content",
			},
		},
		{

			name: "should return paper with max-length valid entry",
			paper: record.PaperWithoutAutofieldEntry{
				Content: maxLengthPaperContent,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := mock_repository.NewMockPaperRepository(ctrl)
			r.EXPECT().
				UpdatePaper(testutil.ModifyOnlyUserId(), "0000000000000001", "1000000000000001", tc.paper).
				Return(&record.PaperEntry{
					Content:   tc.paper.Content,
					CreatedAt: testutil.Date(),
					UpdatedAt: testutil.Date(),
				}, nil)

			s := service.NewPaperService(r)

			userId, err := domain.NewUserIdObject(testutil.ModifyOnlyUserId())
			assert.Nil(t, err)
			projectId, err := domain.NewProjectIdObject("0000000000000001")
			assert.Nil(t, err)
			paperId, err := domain.NewPaperIdObject("1000000000000001")
			assert.Nil(t, err)

			content, err := domain.NewPaperContentObject(tc.paper.Content)
			assert.Nil(t, err)

			paper := domain.NewPaperWithoutAutofieldEntity(*content)

			updatedPaper, sErr := s.UpdatePaper(*userId, *projectId, *paperId, *paper)
			assert.Nil(t, sErr)

			assert.Equal(t, "1000000000000001", updatedPaper.Id().Value())
			assert.Equal(t, tc.paper.Content, updatedPaper.Content().Value())
			assert.Equal(t, testutil.Date(), updatedPaper.CreatedAt().Value())
			assert.Equal(t, testutil.Date(), updatedPaper.UpdatedAt().Value())
		})
	}
}

func TestUpdatePaperInvalidUpdatedEntry(t *testing.T) {
	tooLongPaperContent := testutil.RandomString(40001)

	tt := []struct {
		name          string
		updatedPaper  record.PaperEntry
		expectedError string
	}{
		{
			name: "should return error when content is too long",
			updatedPaper: record.PaperEntry{
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
				UpdatePaper(testutil.ModifyOnlyUserId(), "0000000000000001", "1000000000000001", gomock.Any()).
				Return(&tc.updatedPaper, nil)

			s := service.NewPaperService(r)

			userId, err := domain.NewUserIdObject(tc.updatedPaper.UserId)
			assert.Nil(t, err)
			projectId, err := domain.NewProjectIdObject("0000000000000001")
			assert.Nil(t, err)
			paperId, err := domain.NewPaperIdObject("1000000000000001")
			assert.Nil(t, err)

			content, err := domain.NewPaperContentObject("content")
			assert.Nil(t, err)

			paper := domain.NewPaperWithoutAutofieldEntity(*content)

			updatedPaper, sErr := s.UpdatePaper(*userId, *projectId, *paperId, *paper)
			assert.NotNil(t, sErr)
			assert.Equal(t, service.DomainFailurePanic, sErr.Code())
			assert.Equal(t, fmt.Sprintf("domain failure: %v", tc.expectedError), sErr.Error())
			assert.Nil(t, updatedPaper)
		})
	}
}

func TestUpdatePaperRepositoryError(t *testing.T) {
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
			expectedError: "failed to update paper: paper not found",
			expectedCode:  service.NotFoundError,
		},
		{
			name:          "should return error when repository returns write failure error",
			errorCode:     repository.WriteFailurePanic,
			errorMessage:  "repository error",
			expectedError: "failed to update paper: repository error",
			expectedCode:  service.RepositoryFailurePanic,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := mock_repository.NewMockPaperRepository(ctrl)
			r.EXPECT().
				UpdatePaper(testutil.ModifyOnlyUserId(), "0000000000000001", "1000000000000001", gomock.Any()).
				Return(nil, repository.Errorf(tc.errorCode, tc.errorMessage))

			s := service.NewPaperService(r)

			userId, err := domain.NewUserIdObject(testutil.ModifyOnlyUserId())
			assert.Nil(t, err)
			projectId, err := domain.NewProjectIdObject("0000000000000001")
			assert.Nil(t, err)
			paperId, err := domain.NewPaperIdObject("1000000000000001")
			assert.Nil(t, err)

			content, err := domain.NewPaperContentObject("content")
			assert.Nil(t, err)

			paper := domain.NewPaperWithoutAutofieldEntity(*content)

			updatedPaper, sErr := s.UpdatePaper(*userId, *projectId, *paperId, *paper)
			assert.NotNil(t, sErr)
			assert.Equal(t, tc.expectedCode, sErr.Code())
			assert.Equal(t, fmt.Sprintf("%v: %v", tc.expectedCode, tc.expectedError), sErr.Error())
			assert.Nil(t, updatedPaper)
		})
	}
}
