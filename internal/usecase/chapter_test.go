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
	number, err := domain.NewChapterNumberObject(1)
	assert.Nil(t, err)
	createdAt, err := domain.NewCreatedAtObject(testutil.Date())
	assert.Nil(t, err)
	updatedAt, err := domain.NewUpdatedAtObject(testutil.Date())
	assert.Nil(t, err)

	chapter1 := domain.NewChapterEntity(*id, *name, *number, *createdAt, *updatedAt)

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

	chapter2 := domain.NewChapterEntity(*id, *name, *number, *createdAt, *updatedAt)

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
	assert.Equal(t, int32(1), chapter.Number)
	assert.Len(t, chapter.Sections, 0)

	chapter = res.Chapters[1]
	assert.Equal(t, "1000000000000002", chapter.Id)
	assert.Equal(t, "Chapter 2", chapter.Name)
	assert.Equal(t, int32(2), chapter.Number)
	assert.Len(t, chapter.Sections, 0)
}

func TestListChaptersInvalidArgument(t *testing.T) {
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
		assert.Equal(t, fmt.Sprintf("invalid argument: %s", expectedJson), ucErr.Error())
		assert.Equal(t, usecase.InvalidArgumentError, ucErr.Code())
		assert.Equal(t, tc.expected, *ucErr.Response())

		assert.Nil(t, res)
	}
}

func TestListChaptersServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mock_service.NewMockChapterService(ctrl)

	s.EXPECT().
		ListChapters(gomock.Any(), gomock.Any()).
		Return(nil, service.Errorf(service.RepositoryFailurePanic, "service error"))

	uc := usecase.NewChapterUseCase(s)

	res, ucErr := uc.ListChapters(model.ChapterListRequest{
		User:    model.UserOnlyId{Id: testutil.ReadOnlyUserId()},
		Project: model.ProjectOnlyId{Id: "0000000000000001"},
	})

	assert.Equal(t, "internal error: service error", ucErr.Error())
	assert.Equal(t, usecase.InternalErrorPanic, ucErr.Code())
	assert.Nil(t, res)
}
