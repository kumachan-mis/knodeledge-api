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

func TestListChaptersValidEntity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mock_service.NewMockChapterService(ctrl)

	id, err := domain.NewChapterIdObject("1000000000000001")
	assert.Nil(t, err)
	name, err := domain.NewChapterNameObject("Chapter 1")
	assert.Nil(t, err)
	nextId, err := domain.NewChapterNextIdObject("1000000000000002")
	assert.Nil(t, err)
	createdAt, err := domain.NewCreatedAtObject(testutil.Date())
	assert.Nil(t, err)
	updatedAt, err := domain.NewUpdatedAtObject(testutil.Date())
	assert.Nil(t, err)

	chapter1 := domain.NewChapterEntity(*id, *name, *nextId, *createdAt, *updatedAt)

	id, err = domain.NewChapterIdObject("1000000000000002")
	assert.Nil(t, err)
	name, err = domain.NewChapterNameObject("Chapter 2")
	assert.Nil(t, err)
	nextId, err = domain.NewChapterNextIdObject("")
	assert.Nil(t, err)
	createdAt, err = domain.NewCreatedAtObject(testutil.Date())
	assert.Nil(t, err)
	updatedAt, err = domain.NewUpdatedAtObject(testutil.Date())
	assert.Nil(t, err)

	chapter2 := domain.NewChapterEntity(*id, *name, *nextId, *createdAt, *updatedAt)

	s.EXPECT().
		ListChapters(gomock.Any(), gomock.Any()).
		Do(func(userId domain.UserIdObject, projectId domain.ProjectIdObject) {
			assert.Equal(t, testutil.ReadOnlyUserId(), userId.Value())
			assert.Equal(t, "0000000000000001", projectId.Value())
		}).
		Return([]domain.ChapterEntity{*chapter1, *chapter2}, nil)

	uc := usecase.NewChapterUseCase(s)

	res, ucErr := uc.ListChapters(model.ChapterListRequest{
		User:    model.UserOnlyId{Id: testutil.ReadOnlyUserId()},
		Project: model.ProjectOnlyId{Id: "0000000000000001"},
	})

	assert.Nil(t, ucErr)

	assert.Len(t, res.Chapters, 2)

	chapter := res.Chapters[0]
	assert.Equal(t, "1000000000000001", chapter.Id)
	assert.Equal(t, "Chapter 1", chapter.Name)
	assert.Equal(t, "1000000000000002", chapter.NextId)
	assert.Len(t, chapter.Sections, 0)

	chapter = res.Chapters[1]
	assert.Equal(t, "1000000000000002", chapter.Id)
	assert.Equal(t, "Chapter 2", chapter.Name)
	assert.Equal(t, "", chapter.NextId)
	assert.Len(t, chapter.Sections, 0)
}

func TestListChaptersDomainValidationError(t *testing.T) {
	tt := []struct {
		name      string
		userId    string
		projectId string
		expected  model.ChapterListErrorResponse
	}{
		{
			name:      "empty user id",
			userId:    "",
			projectId: "0000000000000001",
			expected: model.ChapterListErrorResponse{
				User:    model.UserOnlyIdError{Id: "user id is required, but got ''"},
				Project: model.ProjectOnlyIdError{Id: ""},
			},
		},
		{
			name:      "empty project id",
			userId:    testutil.ReadOnlyUserId(),
			projectId: "",
			expected: model.ChapterListErrorResponse{
				User:    model.UserOnlyIdError{Id: ""},
				Project: model.ProjectOnlyIdError{Id: "project id is required, but got ''"},
			},
		},
	}

	for _, tc := range tt {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		s := mock_service.NewMockChapterService(ctrl)

		uc := usecase.NewChapterUseCase(s)

		res, ucErr := uc.ListChapters(model.ChapterListRequest{
			User:    model.UserOnlyId{Id: tc.userId},
			Project: model.ProjectOnlyId{Id: tc.projectId},
		})

		expectedJson, _ := json.Marshal(tc.expected)
		assert.Equal(t, fmt.Sprintf("domain validation error: %s", expectedJson), ucErr.Error())
		assert.Equal(t, usecase.DomainValidationError, ucErr.Code())
		assert.Equal(t, tc.expected, *ucErr.Response())

		assert.Nil(t, res)
	}
}

func TestListChaptersServiceError(t *testing.T) {
	tt := []struct {
		name          string
		errorCode     service.ErrorCode
		errorMessage  string
		expectedError string
		expectedCode  usecase.ErrorCode
	}{
		{
			name:          "should return error when repository returns invalid argument error",
			errorCode:     service.InvalidArgument,
			errorMessage:  "project document does not exist",
			expectedError: "invalid argument: project document does not exist",
			expectedCode:  usecase.InvalidArgumentError,
		},
		{
			name:          "should return error when repository returns failure panic",
			errorCode:     service.RepositoryFailurePanic,
			errorMessage:  "service error",
			expectedError: "internal error: service error",
			expectedCode:  usecase.InternalErrorPanic,
		},
	}

	for _, tc := range tt {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		s := mock_service.NewMockChapterService(ctrl)

		s.EXPECT().
			ListChapters(gomock.Any(), gomock.Any()).
			Return(nil, service.Errorf(tc.errorCode, tc.errorMessage))

		uc := usecase.NewChapterUseCase(s)

		res, ucErr := uc.ListChapters(model.ChapterListRequest{
			User:    model.UserOnlyId{Id: testutil.ReadOnlyUserId()},
			Project: model.ProjectOnlyId{Id: "0000000000000001"},
		})

		assert.Equal(t, tc.expectedError, ucErr.Error())
		assert.Equal(t, tc.expectedCode, ucErr.Code())
		assert.Nil(t, res)
	}
}

func TestCreateChapterValidEntity(t *testing.T) {
	maxLengthChapterName := testutil.RandomString(100)

	tt := []struct {
		name      string
		userId    string
		projectId string
		chapter   model.ChapterWithoutAutofield
	}{
		{
			name:      "valid chapter",
			userId:    testutil.ReadOnlyUserId(),
			projectId: "0000000000000001",
			chapter: model.ChapterWithoutAutofield{
				Name: "Chapter 1",
			},
		},
		{
			name:      "valid chapter with next id",
			userId:    testutil.ReadOnlyUserId(),
			projectId: "0000000000000001",
			chapter: model.ChapterWithoutAutofield{
				Name:   "Chapter 1",
				NextId: "1000000000000002",
			},
		},
		{
			name:      "valid chapter with max length name",
			userId:    testutil.ReadOnlyUserId(),
			projectId: "0000000000000001",
			chapter: model.ChapterWithoutAutofield{
				Name: maxLengthChapterName,
			},
		},
	}

	for _, tc := range tt {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		s := mock_service.NewMockChapterService(ctrl)

		id, err := domain.NewChapterIdObject("1000000000000001")
		assert.Nil(t, err)
		name, err := domain.NewChapterNameObject(tc.chapter.Name)
		assert.Nil(t, err)
		nextId, err := domain.NewChapterNextIdObject(tc.chapter.NextId)
		assert.Nil(t, err)
		createdAt, err := domain.NewCreatedAtObject(testutil.Date())
		assert.Nil(t, err)
		updatedAt, err := domain.NewUpdatedAtObject(testutil.Date())
		assert.Nil(t, err)

		chapter := domain.NewChapterEntity(*id, *name, *nextId, *createdAt, *updatedAt)

		s.EXPECT().
			CreateChapter(gomock.Any(), gomock.Any(), gomock.Any()).
			Do(func(userId domain.UserIdObject, projectId domain.ProjectIdObject, chapter domain.ChapterWithoutAutofieldEntity) {
				assert.Equal(t, tc.userId, userId.Value())
				assert.Equal(t, tc.projectId, projectId.Value())
				assert.Equal(t, tc.chapter.Name, chapter.Name().Value())
				assert.Equal(t, tc.chapter.NextId, chapter.NextId().Value())
			}).
			Return(chapter, nil)

		uc := usecase.NewChapterUseCase(s)

		res, ucErr := uc.CreateChapter(model.ChapterCreateRequest{
			User:    model.UserOnlyId{Id: tc.userId},
			Project: model.ProjectOnlyId{Id: tc.projectId},
			Chapter: tc.chapter,
		})

		assert.Nil(t, ucErr)

		assert.Equal(t, "1000000000000001", res.Chapter.Id)
		assert.Equal(t, tc.chapter.Name, res.Chapter.Name)
		assert.Equal(t, tc.chapter.NextId, res.Chapter.NextId)
		assert.Len(t, res.Chapter.Sections, 0)
	}
}

func TestCreateChapterDomainValidationError(t *testing.T) {
	tooLongChapterName := testutil.RandomString(101)

	tt := []struct {
		name      string
		userId    string
		projectId string
		chapter   model.ChapterWithoutAutofield
		expected  model.ChapterCreateErrorResponse
	}{
		{
			name:      "empty user id",
			userId:    "",
			projectId: "0000000000000001",
			chapter: model.ChapterWithoutAutofield{
				Name: "Chapter 1",
			},
			expected: model.ChapterCreateErrorResponse{
				User:    model.UserOnlyIdError{Id: "user id is required, but got ''"},
				Project: model.ProjectOnlyIdError{Id: ""},
			},
		},
		{
			name:      "empty project id",
			userId:    testutil.ReadOnlyUserId(),
			projectId: "",
			chapter: model.ChapterWithoutAutofield{
				Name: "Chapter 1",
			},
			expected: model.ChapterCreateErrorResponse{
				User:    model.UserOnlyIdError{Id: ""},
				Project: model.ProjectOnlyIdError{Id: "project id is required, but got ''"},
			},
		},
		{
			name:      "empty chapter name",
			userId:    testutil.ReadOnlyUserId(),
			projectId: "0000000000000001",
			chapter: model.ChapterWithoutAutofield{
				Name: "",
			},
			expected: model.ChapterCreateErrorResponse{
				User:    model.UserOnlyIdError{Id: ""},
				Project: model.ProjectOnlyIdError{Id: ""},
				Chapter: model.ChapterWithoutAutofieldError{Name: "chapter name is required, but got ''"},
			},
		},
		{
			name:      "too long chapter name",
			userId:    testutil.ReadOnlyUserId(),
			projectId: "0000000000000001",
			chapter: model.ChapterWithoutAutofield{
				Name: tooLongChapterName,
			},
			expected: model.ChapterCreateErrorResponse{
				User:    model.UserOnlyIdError{Id: ""},
				Project: model.ProjectOnlyIdError{Id: ""},
				Chapter: model.ChapterWithoutAutofieldError{
					Name: fmt.Sprintf("chapter name cannot be longer than 100 characters, but got '%s'", tooLongChapterName),
				},
			},
		},
	}

	for _, tc := range tt {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		s := mock_service.NewMockChapterService(ctrl)

		uc := usecase.NewChapterUseCase(s)

		res, ucErr := uc.CreateChapter(model.ChapterCreateRequest{
			User:    model.UserOnlyId{Id: tc.userId},
			Project: model.ProjectOnlyId{Id: tc.projectId},
			Chapter: tc.chapter,
		})

		expectedJson, _ := json.Marshal(tc.expected)
		assert.Equal(t, fmt.Sprintf("domain validation error: %s", expectedJson), ucErr.Error())
		assert.Equal(t, usecase.DomainValidationError, ucErr.Code())
		assert.Equal(t, tc.expected, *ucErr.Response())

		assert.Nil(t, res)
	}
}

func TestCreateChapterServiceError(t *testing.T) {
	tt := []struct {
		name          string
		errorCode     service.ErrorCode
		errorMessage  string
		expectedError string
		expectedCode  usecase.ErrorCode
	}{
		{
			name:          "should return error when repository returns failure panic",
			errorCode:     service.RepositoryFailurePanic,
			errorMessage:  "service error",
			expectedError: "internal error: service error",
			expectedCode:  usecase.InternalErrorPanic,
		},
		{
			name:          "should return error when repository returns invalid argument error",
			errorCode:     service.InvalidArgument,
			errorMessage:  "id of next chapter does not exist",
			expectedError: "invalid argument: id of next chapter does not exist",
			expectedCode:  usecase.InvalidArgumentError,
		},
	}

	for _, tc := range tt {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		s := mock_service.NewMockChapterService(ctrl)

		s.EXPECT().
			CreateChapter(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, service.Errorf(tc.errorCode, tc.errorMessage))

		uc := usecase.NewChapterUseCase(s)

		res, ucErr := uc.CreateChapter(model.ChapterCreateRequest{
			User:    model.UserOnlyId{Id: testutil.ReadOnlyUserId()},
			Project: model.ProjectOnlyId{Id: "0000000000000001"},
			Chapter: model.ChapterWithoutAutofield{Name: "Chapter 1"},
		})

		assert.Equal(t, tc.expectedError, ucErr.Error())
		assert.Equal(t, tc.expectedCode, ucErr.Code())
		assert.Nil(t, res)
	}
}
