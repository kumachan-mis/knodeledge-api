package usecase_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/kumachan-mis/knodeledge-api/internal/domain"
	"github.com/kumachan-mis/knodeledge-api/internal/openapi"
	"github.com/kumachan-mis/knodeledge-api/internal/service"
	"github.com/kumachan-mis/knodeledge-api/internal/testutil"
	"github.com/kumachan-mis/knodeledge-api/internal/usecase"
	mock_service "github.com/kumachan-mis/knodeledge-api/mock/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestFindPaperValidEntity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mock_service.NewMockPaperService(ctrl)

	id, err := domain.NewPaperIdObject("1000000000000001")
	assert.Nil(t, err)
	content, err := domain.NewPaperContentObject("This is paper content")
	assert.Nil(t, err)
	createdAt, err := domain.NewCreatedAtObject(testutil.Date())
	assert.Nil(t, err)
	updatedAt, err := domain.NewUpdatedAtObject(testutil.Date())
	assert.Nil(t, err)

	paper := domain.NewPaperEntity(*id, *content, *createdAt, *updatedAt)

	s.EXPECT().
		FindPaper(gomock.Any(), gomock.Any(), gomock.Any()).
		Do(func(userId domain.UserIdObject, projectId domain.ProjectIdObject, chapterId domain.ChapterIdObject) {
			assert.Equal(t, testutil.ReadOnlyUserId(), userId.Value())
			assert.Equal(t, "0000000000000001", projectId.Value())
			assert.Equal(t, "1000000000000001", chapterId.Value())
		}).
		Return(paper, nil)

	uc := usecase.NewPaperUseCase(s)

	res, ucErr := uc.FindPaper(openapi.PaperFindRequest{
		UserId:    testutil.ReadOnlyUserId(),
		ProjectId: "0000000000000001",
		ChapterId: "1000000000000001",
	})

	assert.Nil(t, ucErr)

	assert.Equal(t, "1000000000000001", res.Paper.Id)
	assert.Equal(t, "This is paper content", res.Paper.Content)
}

func TestFindPaperDomainValidationError(t *testing.T) {
	tt := []struct {
		name      string
		userId    string
		projectId string
		chapterId string
		expected  openapi.PaperFindErrorResponse
	}{
		{
			name:      "should return error when user id is empty",
			userId:    "",
			projectId: "0000000000000001",
			chapterId: "1000000000000001",
			expected: openapi.PaperFindErrorResponse{
				UserId:    "user id is required, but got ''",
				ProjectId: "",
				ChapterId: "",
			},
		},
		{
			name:      "should return error when project id is empty",
			userId:    testutil.ReadOnlyUserId(),
			projectId: "",
			chapterId: "1000000000000001",
			expected: openapi.PaperFindErrorResponse{
				UserId:    "",
				ProjectId: "project id is required, but got ''",
				ChapterId: "",
			},
		},
		{
			name:      "should return error when chapter id is empty",
			userId:    testutil.ReadOnlyUserId(),
			projectId: "0000000000000001",
			chapterId: "",
			expected: openapi.PaperFindErrorResponse{
				UserId:    "",
				ProjectId: "",
				ChapterId: "chapter id is required, but got ''",
			},
		},
		{
			name:      "should return error when all fields are empty",
			userId:    "",
			projectId: "",
			chapterId: "",
			expected: openapi.PaperFindErrorResponse{
				UserId:    "user id is required, but got ''",
				ProjectId: "project id is required, but got ''",
				ChapterId: "chapter id is required, but got ''",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mock_service.NewMockPaperService(ctrl)

			uc := usecase.NewPaperUseCase(s)

			res, ucErr := uc.FindPaper(openapi.PaperFindRequest{
				UserId:    tc.userId,
				ProjectId: tc.projectId,
				ChapterId: tc.chapterId,
			})

			expectedJson, _ := json.Marshal(tc.expected)
			assert.Equal(t, fmt.Sprintf("domain validation error: %s", expectedJson), ucErr.Error())
			assert.Equal(t, usecase.DomainValidationError, ucErr.Code())
			assert.Equal(t, tc.expected, *ucErr.Response())

			assert.Nil(t, res)
		})
	}
}

func TestFindPaperServiceError(t *testing.T) {
	tt := []struct {
		name          string
		errorCode     service.ErrorCode
		errorMessage  string
		expectedError string
		expectedCode  usecase.ErrorCode
	}{
		{
			name:          "should return error when project not found",
			errorCode:     service.NotFoundError,
			errorMessage:  "failed to find project",
			expectedError: "not found: failed to find project",
			expectedCode:  usecase.NotFoundError,
		},
		{
			name:          "should return error when repository failure",
			errorCode:     service.RepositoryFailurePanic,
			errorMessage:  "service error",
			expectedError: "internal error: service error",
			expectedCode:  usecase.InternalErrorPanic,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mock_service.NewMockPaperService(ctrl)

			uc := usecase.NewPaperUseCase(s)

			s.EXPECT().
				FindPaper(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, service.Errorf(tc.errorCode, "%s", tc.errorMessage))

			res, ucErr := uc.FindPaper(openapi.PaperFindRequest{
				UserId:    testutil.ReadOnlyUserId(),
				ProjectId: "0000000000000001",
				ChapterId: "1000000000000001",
			})

			assert.Nil(t, res)
			assert.Equal(t, tc.expectedError, ucErr.Error())
			assert.Equal(t, tc.expectedCode, ucErr.Code())
		})
	}
}

func TestUpdatePaperValidEntity(t *testing.T) {
	maxLengthPaperContent := testutil.RandomString(40000)

	tt := []struct {
		name    string
		content string
	}{
		{
			name:    "should update paper",
			content: "This is paper content.",
		},
		{
			name:    "should update paper with max length content",
			content: maxLengthPaperContent,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mock_service.NewMockPaperService(ctrl)

			id, err := domain.NewPaperIdObject("1000000000000001")
			assert.Nil(t, err)
			content, err := domain.NewPaperContentObject(tc.content)
			assert.Nil(t, err)
			createdAt, err := domain.NewCreatedAtObject(testutil.Date())
			assert.Nil(t, err)
			updatedAt, err := domain.NewUpdatedAtObject(testutil.Date())
			assert.Nil(t, err)

			paper := domain.NewPaperEntity(*id, *content, *createdAt, *updatedAt)

			s.EXPECT().
				UpdatePaper(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Do(func(userId domain.UserIdObject, projectId domain.ProjectIdObject, paperId domain.PaperIdObject, paper domain.PaperWithoutAutofieldEntity) {
					assert.Equal(t, testutil.ReadOnlyUserId(), userId.Value())
					assert.Equal(t, "0000000000000001", projectId.Value())
					assert.Equal(t, "1000000000000001", paperId.Value())
					assert.Equal(t, tc.content, paper.Content().Value())
				}).
				Return(paper, nil)

			uc := usecase.NewPaperUseCase(s)

			res, ucErr := uc.UpdatePaper(openapi.PaperUpdateRequest{
				User:    openapi.UserOnlyId{Id: testutil.ReadOnlyUserId()},
				Project: openapi.ProjectOnlyId{Id: "0000000000000001"},
				Paper: openapi.Paper{
					Id:      "1000000000000001",
					Content: tc.content,
				},
			})

			assert.Nil(t, ucErr)

			assert.Equal(t, "1000000000000001", res.Paper.Id)
			assert.Equal(t, tc.content, res.Paper.Content)
		})
	}
}

func TestUpdatePaperDomainValidationError(t *testing.T) {
	tooLongPaperContent := testutil.RandomString(40001)

	tt := []struct {
		name      string
		userId    string
		projectId string
		paperId   string
		content   string
		expected  openapi.PaperUpdateErrorResponse
	}{
		{
			name:      "should return error when user id is empty",
			userId:    "",
			projectId: "0000000000000001",
			paperId:   "1000000000000001",
			content:   "This is paper content.",
			expected: openapi.PaperUpdateErrorResponse{
				User:    openapi.UserOnlyIdError{Id: "user id is required, but got ''"},
				Project: openapi.ProjectOnlyIdError{Id: ""},
				Paper:   openapi.PaperError{Id: "", Content: ""},
			},
		},
		{
			name:      "should return error when project id is empty",
			userId:    testutil.ReadOnlyUserId(),
			projectId: "",
			paperId:   "1000000000000001",
			content:   "This is paper content.",
			expected: openapi.PaperUpdateErrorResponse{
				User:    openapi.UserOnlyIdError{Id: ""},
				Project: openapi.ProjectOnlyIdError{Id: "project id is required, but got ''"},
				Paper:   openapi.PaperError{Id: "", Content: ""},
			},
		},
		{
			name:      "should return error when paper id is empty",
			userId:    testutil.ReadOnlyUserId(),
			projectId: "0000000000000001",
			paperId:   "",
			content:   "This is paper content.",
			expected: openapi.PaperUpdateErrorResponse{
				User:    openapi.UserOnlyIdError{Id: ""},
				Project: openapi.ProjectOnlyIdError{Id: ""},
				Paper:   openapi.PaperError{Id: "paper id is required, but got ''", Content: ""},
			},
		},
		{
			name:      "should return error when content is too long",
			userId:    testutil.ReadOnlyUserId(),
			projectId: "0000000000000001",
			paperId:   "1000000000000001",
			content:   tooLongPaperContent,
			expected: openapi.PaperUpdateErrorResponse{
				User:    openapi.UserOnlyIdError{Id: ""},
				Project: openapi.ProjectOnlyIdError{Id: ""},
				Paper: openapi.PaperError{
					Id:      "",
					Content: "paper content must be less than or equal to 40000 bytes, but got 40001 bytes",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mock_service.NewMockPaperService(ctrl)

			uc := usecase.NewPaperUseCase(s)

			res, ucErr := uc.UpdatePaper(openapi.PaperUpdateRequest{
				User:    openapi.UserOnlyId{Id: tc.userId},
				Project: openapi.ProjectOnlyId{Id: tc.projectId},
				Paper: openapi.Paper{
					Id:      tc.paperId,
					Content: tc.content,
				},
			})

			expectedJson, _ := json.Marshal(tc.expected)
			assert.Equal(t, fmt.Sprintf("domain validation error: %s", expectedJson), ucErr.Error())
			assert.Equal(t, usecase.DomainValidationError, ucErr.Code())
			assert.Equal(t, tc.expected, *ucErr.Response())

			assert.Nil(t, res)
		})
	}
}

func TestUpdatePaperServiceError(t *testing.T) {
	tt := []struct {
		name          string
		errorCode     service.ErrorCode
		errorMessage  string
		expectedError string
		expectedCode  usecase.ErrorCode
	}{
		{
			name:          "should return error when project not found",
			errorCode:     service.NotFoundError,
			errorMessage:  "failed to find project",
			expectedError: "not found: failed to find project",
			expectedCode:  usecase.NotFoundError,
		},
		{
			name:          "should return error when repository failure",
			errorCode:     service.RepositoryFailurePanic,
			errorMessage:  "service error",
			expectedError: "internal error: service error",
			expectedCode:  usecase.InternalErrorPanic,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mock_service.NewMockPaperService(ctrl)

			uc := usecase.NewPaperUseCase(s)

			s.EXPECT().
				UpdatePaper(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, service.Errorf(tc.errorCode, "%s", tc.errorMessage))

			res, ucErr := uc.UpdatePaper(openapi.PaperUpdateRequest{
				User:    openapi.UserOnlyId{Id: testutil.ReadOnlyUserId()},
				Project: openapi.ProjectOnlyId{Id: "0000000000000001"},
				Paper: openapi.Paper{
					Id:      "1000000000000001",
					Content: "This is paper content.",
				},
			})

			assert.Nil(t, res)
			assert.Equal(t, tc.expectedError, ucErr.Error())
			assert.Equal(t, tc.expectedCode, ucErr.Code())
		})
	}
}
