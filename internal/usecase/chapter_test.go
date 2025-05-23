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

func TestListChaptersValidEntity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mock_service.NewMockChapterService(ctrl)

	id, err := domain.NewChapterIdObject("1000000000000001")
	assert.Nil(t, err)
	name, err := domain.NewChapterNameObject("Chapter 1")
	assert.Nil(t, err)
	number, err := domain.NewChapterNumberObject(1)
	assert.Nil(t, err)
	createdAt, err := domain.NewCreatedAtObject(testutil.Date())
	assert.Nil(t, err)
	updatedAt, err := domain.NewUpdatedAtObject(testutil.Date())
	assert.Nil(t, err)

	sectionId, err := domain.NewSectionIdObject("2000000000000001")
	assert.Nil(t, err)
	sectionName, err := domain.NewSectionNameObject("Section 1")
	assert.Nil(t, err)
	sectionCreatedAt, err := domain.NewCreatedAtObject(testutil.Date())
	assert.Nil(t, err)
	sectionUpdatedAt, err := domain.NewUpdatedAtObject(testutil.Date())
	assert.Nil(t, err)

	section1 := domain.NewSectionOfChapterEntity(*sectionId, *sectionName, *sectionCreatedAt, *sectionUpdatedAt)
	sections := &[]domain.SectionOfChapterEntity{*section1}

	chapter1 := domain.NewChapterEntity(*id, *name, *number, *sections, *createdAt, *updatedAt)

	id, err = domain.NewChapterIdObject("1000000000000002")
	assert.Nil(t, err)
	name, err = domain.NewChapterNameObject("Chapter 2")
	assert.Nil(t, err)
	number, err = domain.NewChapterNumberObject(2)
	assert.Nil(t, err)
	createdAt, err = domain.NewCreatedAtObject(testutil.Date())
	assert.Nil(t, err)
	updatedAt, err = domain.NewUpdatedAtObject(testutil.Date())
	assert.Nil(t, err)

	sectionId, err = domain.NewSectionIdObject("2000000000000002")
	assert.Nil(t, err)
	sectionName, err = domain.NewSectionNameObject("Section 1")
	assert.Nil(t, err)
	sectionCreatedAt, err = domain.NewCreatedAtObject(testutil.Date())
	assert.Nil(t, err)
	sectionUpdatedAt, err = domain.NewUpdatedAtObject(testutil.Date())
	assert.Nil(t, err)

	section1 = domain.NewSectionOfChapterEntity(*sectionId, *sectionName, *sectionCreatedAt, *sectionUpdatedAt)
	sections = &[]domain.SectionOfChapterEntity{*section1}

	chapter2 := domain.NewChapterEntity(*id, *name, *number, *sections, *createdAt, *updatedAt)

	s.EXPECT().
		ListChapters(gomock.Any(), gomock.Any()).
		Do(func(userId domain.UserIdObject, projectId domain.ProjectIdObject) {
			assert.Equal(t, testutil.ReadOnlyUserId(), userId.Value())
			assert.Equal(t, "0000000000000001", projectId.Value())
		}).
		Return([]domain.ChapterEntity{*chapter1, *chapter2}, nil)

	uc := usecase.NewChapterUseCase(s)

	res, ucErr := uc.ListChapters(openapi.ChapterListRequest{
		UserId:    testutil.ReadOnlyUserId(),
		ProjectId: "0000000000000001",
	})

	assert.Nil(t, ucErr)

	assert.Len(t, res.Chapters, 2)

	chapter := res.Chapters[0]
	assert.Equal(t, "1000000000000001", chapter.Id)
	assert.Equal(t, "Chapter 1", chapter.Name)
	assert.Equal(t, int32(1), chapter.Number)
	assert.Len(t, chapter.Sections, 1)

	section := chapter.Sections[0]
	assert.Equal(t, "2000000000000001", section.Id)
	assert.Equal(t, "Section 1", section.Name)

	chapter = res.Chapters[1]
	assert.Equal(t, "1000000000000002", chapter.Id)
	assert.Equal(t, "Chapter 2", chapter.Name)
	assert.Equal(t, int32(2), chapter.Number)
	assert.Len(t, chapter.Sections, 1)

	section = chapter.Sections[0]
	assert.Equal(t, "2000000000000002", section.Id)
	assert.Equal(t, "Section 1", section.Name)

}

func TestListChaptersDomainValidationError(t *testing.T) {
	tt := []struct {
		name      string
		userId    string
		projectId string
		expected  openapi.ChapterListErrorResponse
	}{
		{
			name:      "should return error when user id is empty",
			userId:    "",
			projectId: "0000000000000001",
			expected: openapi.ChapterListErrorResponse{
				UserId:    "user id is required, but got ''",
				ProjectId: "",
			},
		},
		{
			name:      "should return error when project id is empty",
			userId:    testutil.ReadOnlyUserId(),
			projectId: "",
			expected: openapi.ChapterListErrorResponse{
				UserId:    "",
				ProjectId: "project id is required, but got ''",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mock_service.NewMockChapterService(ctrl)

			uc := usecase.NewChapterUseCase(s)

			res, ucErr := uc.ListChapters(openapi.ChapterListRequest{
				UserId:    tc.userId,
				ProjectId: tc.projectId,
			})

			expectedJson, _ := json.Marshal(tc.expected)
			assert.Equal(t, fmt.Sprintf("domain validation error: %s", expectedJson), ucErr.Error())
			assert.Equal(t, usecase.DomainValidationError, ucErr.Code())
			assert.Equal(t, tc.expected, *ucErr.Response())

			assert.Nil(t, res)
		})
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
			name:          "should return error when repository returns not found error",
			errorCode:     service.NotFoundError,
			errorMessage:  "failed to fetch project",
			expectedError: "not found: failed to fetch project",
			expectedCode:  usecase.NotFoundError,
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
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mock_service.NewMockChapterService(ctrl)

			s.EXPECT().
				ListChapters(gomock.Any(), gomock.Any()).
				Return(nil, service.Errorf(tc.errorCode, "%s", tc.errorMessage))

			uc := usecase.NewChapterUseCase(s)

			res, ucErr := uc.ListChapters(openapi.ChapterListRequest{
				UserId:    testutil.ReadOnlyUserId(),
				ProjectId: "0000000000000001",
			})

			assert.Equal(t, tc.expectedError, ucErr.Error())
			assert.Equal(t, tc.expectedCode, ucErr.Code())
			assert.Nil(t, res)
		})
	}
}

func TestCreateChapterValidEntity(t *testing.T) {
	maxLengthChapterName := testutil.RandomString(100)

	tt := []struct {
		name      string
		userId    string
		projectId string
		chapter   openapi.ChapterWithoutAutofield
	}{
		{
			name:      "should create chapter",
			userId:    testutil.ReadOnlyUserId(),
			projectId: "0000000000000001",
			chapter: openapi.ChapterWithoutAutofield{
				Name:   "Chapter 1",
				Number: int32(1),
			},
		},
		{
			name:      "should create chapter with max length name",
			userId:    testutil.ReadOnlyUserId(),
			projectId: "0000000000000001",
			chapter: openapi.ChapterWithoutAutofield{
				Name:   maxLengthChapterName,
				Number: int32(1),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mock_service.NewMockChapterService(ctrl)

			chapterId, err := domain.NewChapterIdObject("1000000000000001")
			assert.Nil(t, err)
			name, err := domain.NewChapterNameObject(tc.chapter.Name)
			assert.Nil(t, err)
			number, err := domain.NewChapterNumberObject(int(tc.chapter.Number))
			assert.Nil(t, err)
			sections := &[]domain.SectionOfChapterEntity{}
			createdAt, err := domain.NewCreatedAtObject(testutil.Date())
			assert.Nil(t, err)
			updatedAt, err := domain.NewUpdatedAtObject(testutil.Date())
			assert.Nil(t, err)

			chapter := domain.NewChapterEntity(*chapterId, *name, *number, *sections, *createdAt, *updatedAt)

			s.EXPECT().
				CreateChapter(gomock.Any(), gomock.Any(), gomock.Any()).
				Do(func(userId domain.UserIdObject, projectId domain.ProjectIdObject, chapter domain.ChapterWithoutAutofieldEntity) {
					assert.Equal(t, tc.userId, userId.Value())
					assert.Equal(t, tc.projectId, projectId.Value())
					assert.Equal(t, tc.chapter.Name, chapter.Name().Value())
					assert.Equal(t, int(tc.chapter.Number), chapter.Number().Value())
				}).
				Return(chapter, nil)

			uc := usecase.NewChapterUseCase(s)

			res, ucErr := uc.CreateChapter(openapi.ChapterCreateRequest{
				User:    openapi.UserOnlyId{Id: tc.userId},
				Project: openapi.ProjectOnlyId{Id: tc.projectId},
				Chapter: tc.chapter,
			})

			assert.Nil(t, ucErr)

			assert.Equal(t, "1000000000000001", res.Chapter.Id)
			assert.Equal(t, tc.chapter.Name, res.Chapter.Name)
			assert.Equal(t, tc.chapter.Number, res.Chapter.Number)
			assert.Len(t, res.Chapter.Sections, 0)
		})
	}
}

func TestCreateChapterDomainValidationError(t *testing.T) {
	tooLongChapterName := testutil.RandomString(101)

	tt := []struct {
		name      string
		userId    string
		projectId string
		chapter   openapi.ChapterWithoutAutofield
		expected  openapi.ChapterCreateErrorResponse
	}{
		{
			name:      "should return error when user id is empty",
			userId:    "",
			projectId: "0000000000000001",
			chapter: openapi.ChapterWithoutAutofield{
				Name:   "Chapter 1",
				Number: int32(1),
			},
			expected: openapi.ChapterCreateErrorResponse{
				User:    openapi.UserOnlyIdError{Id: "user id is required, but got ''"},
				Project: openapi.ProjectOnlyIdError{Id: ""},
				Chapter: openapi.ChapterWithoutAutofieldError{Name: "", Number: ""},
			},
		},
		{
			name:      "should return error when project id is empty",
			userId:    testutil.ReadOnlyUserId(),
			projectId: "",
			chapter: openapi.ChapterWithoutAutofield{
				Name:   "Chapter 1",
				Number: int32(1),
			},
			expected: openapi.ChapterCreateErrorResponse{
				User:    openapi.UserOnlyIdError{Id: ""},
				Project: openapi.ProjectOnlyIdError{Id: "project id is required, but got ''"},
				Chapter: openapi.ChapterWithoutAutofieldError{Name: "", Number: ""},
			},
		},
		{
			name:      "should return error when chapter name is empty",
			userId:    testutil.ReadOnlyUserId(),
			projectId: "0000000000000001",
			chapter: openapi.ChapterWithoutAutofield{
				Name:   "",
				Number: int32(1),
			},
			expected: openapi.ChapterCreateErrorResponse{
				User:    openapi.UserOnlyIdError{Id: ""},
				Project: openapi.ProjectOnlyIdError{Id: ""},
				Chapter: openapi.ChapterWithoutAutofieldError{Name: "chapter name is required, but got ''"},
			},
		},
		{
			name:      "should return error when chapter name is too long",
			userId:    testutil.ReadOnlyUserId(),
			projectId: "0000000000000001",
			chapter: openapi.ChapterWithoutAutofield{
				Name:   tooLongChapterName,
				Number: int32(1),
			},
			expected: openapi.ChapterCreateErrorResponse{
				User:    openapi.UserOnlyIdError{Id: ""},
				Project: openapi.ProjectOnlyIdError{Id: ""},
				Chapter: openapi.ChapterWithoutAutofieldError{
					Name: fmt.Sprintf("chapter name cannot be longer than 100 characters, but got '%v'",
						tooLongChapterName),
				},
			},
		},
		{
			name:      "should return error when chapter number is zero",
			userId:    testutil.ReadOnlyUserId(),
			projectId: "0000000000000001",
			chapter: openapi.ChapterWithoutAutofield{
				Name:   "Chapter 1",
				Number: int32(0),
			},
			expected: openapi.ChapterCreateErrorResponse{
				User:    openapi.UserOnlyIdError{Id: ""},
				Project: openapi.ProjectOnlyIdError{Id: ""},
				Chapter: openapi.ChapterWithoutAutofieldError{Number: "chapter number must be greater than 0, but got 0"},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mock_service.NewMockChapterService(ctrl)

			uc := usecase.NewChapterUseCase(s)

			res, ucErr := uc.CreateChapter(openapi.ChapterCreateRequest{
				User:    openapi.UserOnlyId{Id: tc.userId},
				Project: openapi.ProjectOnlyId{Id: tc.projectId},
				Chapter: tc.chapter,
			})

			expectedJson, _ := json.Marshal(tc.expected)
			assert.Equal(t, fmt.Sprintf("domain validation error: %s", expectedJson), ucErr.Error())
			assert.Equal(t, usecase.DomainValidationError, ucErr.Code())
			assert.Equal(t, tc.expected, *ucErr.Response())

			assert.Nil(t, res)
		})
	}
}

func TestCreateChapterChapterServiceError(t *testing.T) {
	tt := []struct {
		name          string
		errorCode     service.ErrorCode
		errorMessage  string
		expectedError string
		expectedCode  usecase.ErrorCode
	}{
		{
			name:          "should return error when repository returns invalid argument error",
			errorCode:     service.InvalidArgumentError,
			errorMessage:  "chapter number is too large",
			expectedError: "invalid argument: chapter number is too large",
			expectedCode:  usecase.InvalidArgumentError,
		},
		{
			name:          "should return error when repository returns not found error",
			errorCode:     service.NotFoundError,
			errorMessage:  "failed to fetch project",
			expectedError: "not found: failed to fetch project",
			expectedCode:  usecase.NotFoundError,
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
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mock_service.NewMockChapterService(ctrl)

			s.EXPECT().
				CreateChapter(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, service.Errorf(tc.errorCode, "%s", tc.errorMessage))

			uc := usecase.NewChapterUseCase(s)

			res, ucErr := uc.CreateChapter(openapi.ChapterCreateRequest{
				User: openapi.UserOnlyId{
					Id: testutil.ReadOnlyUserId(),
				},
				Project: openapi.ProjectOnlyId{
					Id: "0000000000000001",
				},
				Chapter: openapi.ChapterWithoutAutofield{
					Name:   "Chapter 1",
					Number: int32(1),
				},
			})

			assert.Equal(t, tc.expectedError, ucErr.Error())
			assert.Equal(t, tc.expectedCode, ucErr.Code())
			assert.Nil(t, res)
		})
	}
}

func TestUpdateChapterValidEntity(t *testing.T) {
	maxLengthChapterName := testutil.RandomString(100)

	tt := []struct {
		name      string
		userId    string
		projectId string
		chapterId string
		chapter   openapi.ChapterWithoutAutofield
	}{
		{
			name:      "should update chapter",
			userId:    testutil.ReadOnlyUserId(),
			projectId: "0000000000000001",
			chapterId: "1000000000000001",
			chapter: openapi.ChapterWithoutAutofield{
				Name:   "Chapter 1",
				Number: int32(1),
			},
		},
		{
			name:      "should update chapter with max length name",
			userId:    testutil.ReadOnlyUserId(),
			projectId: "0000000000000001",
			chapterId: "1000000000000001",
			chapter: openapi.ChapterWithoutAutofield{
				Name:   maxLengthChapterName,
				Number: int32(1),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mock_service.NewMockChapterService(ctrl)

			id, err := domain.NewChapterIdObject(tc.chapterId)
			assert.Nil(t, err)
			name, err := domain.NewChapterNameObject(tc.chapter.Name)
			assert.Nil(t, err)
			number, err := domain.NewChapterNumberObject(int(tc.chapter.Number))
			assert.Nil(t, err)
			sections := &[]domain.SectionOfChapterEntity{}
			createdAt, err := domain.NewCreatedAtObject(testutil.Date())
			assert.Nil(t, err)
			updatedAt, err := domain.NewUpdatedAtObject(testutil.Date())
			assert.Nil(t, err)

			chapter := domain.NewChapterEntity(*id, *name, *number, *sections, *createdAt, *updatedAt)

			s.EXPECT().
				UpdateChapter(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Do(func(userId domain.UserIdObject, projectId domain.ProjectIdObject, chapterId domain.ChapterIdObject, chapter domain.ChapterWithoutAutofieldEntity) {
					assert.Equal(t, tc.userId, userId.Value())
					assert.Equal(t, tc.projectId, projectId.Value())
					assert.Equal(t, tc.chapterId, chapterId.Value())
					assert.Equal(t, tc.chapter.Name, chapter.Name().Value())
					assert.Equal(t, int(tc.chapter.Number), chapter.Number().Value())
				}).
				Return(chapter, nil)

			uc := usecase.NewChapterUseCase(s)

			res, ucErr := uc.UpdateChapter(openapi.ChapterUpdateRequest{
				User:    openapi.UserOnlyId{Id: tc.userId},
				Project: openapi.ProjectOnlyId{Id: tc.projectId},
				Chapter: openapi.Chapter{
					Id:     tc.chapterId,
					Name:   tc.chapter.Name,
					Number: tc.chapter.Number,
				},
			})

			assert.Nil(t, ucErr)

			assert.Equal(t, tc.chapterId, res.Chapter.Id)
			assert.Equal(t, tc.chapter.Name, res.Chapter.Name)
			assert.Equal(t, tc.chapter.Number, res.Chapter.Number)
			assert.Len(t, res.Chapter.Sections, 0)
		})
	}
}

func TestUpdateChapterDomainValidationError(t *testing.T) {
	tooLongChapterName := testutil.RandomString(101)

	tt := []struct {
		name      string
		userId    string
		projectId string
		chapterId string
		chapter   openapi.Chapter
		expected  openapi.ChapterUpdateErrorResponse
	}{
		{
			name:      "should return error when user id is empty",
			userId:    "",
			projectId: "0000000000000001",
			chapterId: "1000000000000001",
			chapter: openapi.Chapter{
				Id:     "1000000000000001",
				Name:   "Chapter 1",
				Number: int32(1),
			},
			expected: openapi.ChapterUpdateErrorResponse{
				User:    openapi.UserOnlyIdError{Id: "user id is required, but got ''"},
				Project: openapi.ProjectOnlyIdError{Id: ""},
				Chapter: openapi.ChapterError{Id: "", Name: "", Number: ""},
			},
		},
		{
			name:      "should return error when project id is empty",
			userId:    testutil.ReadOnlyUserId(),
			projectId: "",
			chapterId: "1000000000000001",
			chapter: openapi.Chapter{
				Id:     "1000000000000001",
				Name:   "Chapter 1",
				Number: int32(1),
			},
			expected: openapi.ChapterUpdateErrorResponse{
				User:    openapi.UserOnlyIdError{Id: ""},
				Project: openapi.ProjectOnlyIdError{Id: "project id is required, but got ''"},
				Chapter: openapi.ChapterError{Id: "", Name: "", Number: ""},
			},
		},
		{
			name:      "should return error when chapter id is empty",
			userId:    testutil.ReadOnlyUserId(),
			projectId: "0000000000000001",
			chapterId: "",
			chapter: openapi.Chapter{
				Id:     "",
				Name:   "Chapter 1",
				Number: int32(1),
			},
			expected: openapi.ChapterUpdateErrorResponse{
				User:    openapi.UserOnlyIdError{Id: ""},
				Project: openapi.ProjectOnlyIdError{Id: ""},
				Chapter: openapi.ChapterError{Id: "chapter id is required, but got ''", Name: "", Number: ""},
			},
		},
		{
			name:      "should return error when chapter name is empty",
			userId:    testutil.ReadOnlyUserId(),
			projectId: "0000000000000001",
			chapterId: "1000000000000001",
			chapter: openapi.Chapter{
				Id:     "1000000000000001",
				Name:   "",
				Number: int32(1),
			},
			expected: openapi.ChapterUpdateErrorResponse{
				User:    openapi.UserOnlyIdError{Id: ""},
				Project: openapi.ProjectOnlyIdError{Id: ""},
				Chapter: openapi.ChapterError{Id: "", Name: "chapter name is required, but got ''", Number: ""},
			},
		},
		{
			name:      "should return error when chapter name is too long",
			userId:    testutil.ReadOnlyUserId(),
			projectId: "0000000000000001",
			chapterId: "1000000000000001",
			chapter: openapi.Chapter{
				Id:     "1000000000000001",
				Name:   tooLongChapterName,
				Number: int32(1),
			},
			expected: openapi.ChapterUpdateErrorResponse{
				User:    openapi.UserOnlyIdError{Id: ""},
				Project: openapi.ProjectOnlyIdError{Id: ""},
				Chapter: openapi.ChapterError{
					Id: "",
					Name: fmt.Sprintf("chapter name cannot be longer than 100 characters, but got '%v'",
						tooLongChapterName),
					Number: "",
				},
			},
		},
		{
			name:      "should return error when chapter number is zero",
			userId:    testutil.ReadOnlyUserId(),
			projectId: "0000000000000001",
			chapterId: "1000000000000001",
			chapter: openapi.Chapter{
				Id:     "1000000000000001",
				Name:   "Chapter 1",
				Number: int32(0),
			},
			expected: openapi.ChapterUpdateErrorResponse{
				User:    openapi.UserOnlyIdError{Id: ""},
				Project: openapi.ProjectOnlyIdError{Id: ""},
				Chapter: openapi.ChapterError{Id: "", Name: "", Number: "chapter number must be greater than 0, but got 0"},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mock_service.NewMockChapterService(ctrl)

			uc := usecase.NewChapterUseCase(s)

			res, ucErr := uc.UpdateChapter(openapi.ChapterUpdateRequest{
				User:    openapi.UserOnlyId{Id: tc.userId},
				Project: openapi.ProjectOnlyId{Id: tc.projectId},
				Chapter: tc.chapter,
			})

			expectedJson, _ := json.Marshal(tc.expected)
			assert.Equal(t, fmt.Sprintf("domain validation error: %s", expectedJson), ucErr.Error())
			assert.Equal(t, usecase.DomainValidationError, ucErr.Code())
			assert.Equal(t, tc.expected, *ucErr.Response())

			assert.Nil(t, res)
		})
	}
}

func TestUpdateChapterServiceError(t *testing.T) {
	tt := []struct {
		name          string
		errorCode     service.ErrorCode
		errorMessage  string
		expectedError string
		expectedCode  usecase.ErrorCode
	}{
		{
			name:          "should return error when repository returns not found error",
			errorCode:     service.NotFoundError,
			errorMessage:  "failed to update chapter",
			expectedError: "not found: failed to update chapter",
			expectedCode:  usecase.NotFoundError,
		},
		{
			name:          "should return error when repository returns invalid argument error",
			errorCode:     service.InvalidArgumentError,
			errorMessage:  "chapter number is too large",
			expectedError: "invalid argument: chapter number is too large",
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
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mock_service.NewMockChapterService(ctrl)

			s.EXPECT().
				UpdateChapter(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, service.Errorf(tc.errorCode, "%s", tc.errorMessage))

			uc := usecase.NewChapterUseCase(s)

			res, ucErr := uc.UpdateChapter(openapi.ChapterUpdateRequest{
				User: openapi.UserOnlyId{
					Id: testutil.ReadOnlyUserId(),
				},
				Project: openapi.ProjectOnlyId{
					Id: "0000000000000001",
				},
				Chapter: openapi.Chapter{
					Id:     "1000000000000001",
					Name:   "Chapter 1",
					Number: int32(1),
				},
			})

			assert.Equal(t, tc.expectedError, ucErr.Error())
			assert.Equal(t, tc.expectedCode, ucErr.Code())
			assert.Nil(t, res)
		})
	}
}

func TestDeleteChapterValidEntity(t *testing.T) {
	tt := []struct {
		name      string
		userId    string
		projectId string
		chapterId string
	}{
		{
			name:      "should delete chapter",
			userId:    testutil.ReadOnlyUserId(),
			projectId: "0000000000000001",
			chapterId: "1000000000000001",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mock_service.NewMockChapterService(ctrl)

			s.EXPECT().
				DeleteChapter(gomock.Any(), gomock.Any(), gomock.Any()).
				Do(func(userId domain.UserIdObject, projectId domain.ProjectIdObject, chapterId domain.ChapterIdObject) {
					assert.Equal(t, tc.userId, userId.Value())
					assert.Equal(t, tc.projectId, projectId.Value())
					assert.Equal(t, tc.chapterId, chapterId.Value())
				}).
				Return(nil)

			uc := usecase.NewChapterUseCase(s)

			ucErr := uc.DeleteChapter(openapi.ChapterDeleteRequest{
				User:    openapi.UserOnlyId{Id: tc.userId},
				Project: openapi.ProjectOnlyId{Id: tc.projectId},
				Chapter: openapi.ChapterOnlyId{Id: tc.chapterId},
			})

			assert.Nil(t, ucErr)
		})
	}
}

func TestDeleteChapterDomainValidationError(t *testing.T) {
	tt := []struct {
		name      string
		userId    string
		projectId string
		chapterId string
		expected  openapi.ChapterDeleteErrorResponse
	}{
		{
			name:      "should return error when user id is empty",
			userId:    "",
			projectId: "0000000000000001",
			chapterId: "1000000000000001",
			expected: openapi.ChapterDeleteErrorResponse{
				User:    openapi.UserOnlyIdError{Id: "user id is required, but got ''"},
				Project: openapi.ProjectOnlyIdError{Id: ""},
				Chapter: openapi.ChapterOnlyIdError{Id: ""},
			},
		},
		{
			name:      "should return error when project id is empty",
			userId:    testutil.ReadOnlyUserId(),
			projectId: "",
			chapterId: "1000000000000001",
			expected: openapi.ChapterDeleteErrorResponse{
				User:    openapi.UserOnlyIdError{Id: ""},
				Project: openapi.ProjectOnlyIdError{Id: "project id is required, but got ''"},
				Chapter: openapi.ChapterOnlyIdError{Id: ""},
			},
		},
		{
			name:      "should return error when chapter id is empty",
			userId:    testutil.ReadOnlyUserId(),
			projectId: "0000000000000001",
			chapterId: "",
			expected: openapi.ChapterDeleteErrorResponse{
				User:    openapi.UserOnlyIdError{Id: ""},
				Project: openapi.ProjectOnlyIdError{Id: ""},
				Chapter: openapi.ChapterOnlyIdError{Id: "chapter id is required, but got ''"},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mock_service.NewMockChapterService(ctrl)

			uc := usecase.NewChapterUseCase(s)

			ucErr := uc.DeleteChapter(openapi.ChapterDeleteRequest{
				User:    openapi.UserOnlyId{Id: tc.userId},
				Project: openapi.ProjectOnlyId{Id: tc.projectId},
				Chapter: openapi.ChapterOnlyId{Id: tc.chapterId},
			})

			expectedJson, _ := json.Marshal(tc.expected)
			assert.Equal(t, fmt.Sprintf("domain validation error: %s", expectedJson), ucErr.Error())
			assert.Equal(t, usecase.DomainValidationError, ucErr.Code())
			assert.Equal(t, tc.expected, *ucErr.Response())
		})
	}
}

func TestDeleteChapterServiceError(t *testing.T) {
	tt := []struct {
		name          string
		errorCode     service.ErrorCode
		errorMessage  string
		expectedError string
		expectedCode  usecase.ErrorCode
	}{
		{
			name:          "should return error when repository returns not found error",
			errorCode:     service.NotFoundError,
			errorMessage:  "failed to delete chapter",
			expectedError: "not found: failed to delete chapter",
			expectedCode:  usecase.NotFoundError,
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
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mock_service.NewMockChapterService(ctrl)

			s.EXPECT().
				DeleteChapter(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(service.Errorf(tc.errorCode, "%s", tc.errorMessage))

			uc := usecase.NewChapterUseCase(s)

			ucErr := uc.DeleteChapter(openapi.ChapterDeleteRequest{
				User: openapi.UserOnlyId{
					Id: testutil.ReadOnlyUserId(),
				},
				Project: openapi.ProjectOnlyId{
					Id: "0000000000000001",
				},
				Chapter: openapi.ChapterOnlyId{
					Id: "1000000000000001",
				},
			})

			assert.Equal(t, tc.expectedError, ucErr.Error())
			assert.Equal(t, tc.expectedCode, ucErr.Code())
		})
	}
}
