package usecase_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/kumachan-mis/knodeledge-api/internal/domain"
	"github.com/kumachan-mis/knodeledge-api/internal/model"
	"github.com/kumachan-mis/knodeledge-api/internal/service"
	"github.com/kumachan-mis/knodeledge-api/internal/usecase"
	mock_service "github.com/kumachan-mis/knodeledge-api/mock/service"
	"github.com/kumachan-mis/knodeledge-api/test/testutil"
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

	res, ucErr := uc.FindPaper(model.PaperFindRequest{
		User:    model.UserOnlyId{Id: testutil.ReadOnlyUserId()},
		Project: model.ProjectOnlyId{Id: "0000000000000001"},
		Chapter: model.ChapterOnlyId{Id: "1000000000000001"},
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
		expected  model.PaperFindErrorResponse
	}{
		{
			name:      "should return error when user id is empty",
			userId:    "",
			projectId: "0000000000000001",
			chapterId: "1000000000000001",
			expected: model.PaperFindErrorResponse{
				User:    model.UserOnlyIdError{Id: "user id is required, but got ''"},
				Project: model.ProjectOnlyIdError{Id: ""},
				Chapter: model.ChapterOnlyIdError{Id: ""},
			},
		},
		{
			name:      "should return error when project id is empty",
			userId:    testutil.ReadOnlyUserId(),
			projectId: "",
			chapterId: "1000000000000001",
			expected: model.PaperFindErrorResponse{
				User:    model.UserOnlyIdError{Id: ""},
				Project: model.ProjectOnlyIdError{Id: "project id is required, but got ''"},
				Chapter: model.ChapterOnlyIdError{Id: ""},
			},
		},
		{
			name:      "should return error when chapter id is empty",
			userId:    testutil.ReadOnlyUserId(),
			projectId: "0000000000000001",
			chapterId: "",
			expected: model.PaperFindErrorResponse{
				User:    model.UserOnlyIdError{Id: ""},
				Project: model.ProjectOnlyIdError{Id: ""},
				Chapter: model.ChapterOnlyIdError{Id: "chapter id is required, but got ''"},
			},
		},
		{
			name:      "should return error when all fields are empty",
			userId:    "",
			projectId: "",
			chapterId: "",
			expected: model.PaperFindErrorResponse{
				User:    model.UserOnlyIdError{Id: "user id is required, but got ''"},
				Project: model.ProjectOnlyIdError{Id: "project id is required, but got ''"},
				Chapter: model.ChapterOnlyIdError{Id: "chapter id is required, but got ''"},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mock_service.NewMockPaperService(ctrl)

			uc := usecase.NewPaperUseCase(s)

			res, ucErr := uc.FindPaper(model.PaperFindRequest{
				User:    model.UserOnlyId{Id: tc.userId},
				Project: model.ProjectOnlyId{Id: tc.projectId},
				Chapter: model.ChapterOnlyId{Id: tc.chapterId},
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
				Return(nil, service.Errorf(tc.errorCode, tc.errorMessage))

			res, ucErr := uc.FindPaper(model.PaperFindRequest{
				User:    model.UserOnlyId{Id: testutil.ReadOnlyUserId()},
				Project: model.ProjectOnlyId{Id: "0000000000000001"},
				Chapter: model.ChapterOnlyId{Id: "1000000000000001"},
			})

			assert.Nil(t, res)
			assert.Equal(t, tc.expectedError, ucErr.Error())
			assert.Equal(t, tc.expectedCode, ucErr.Code())
		})
	}
}
