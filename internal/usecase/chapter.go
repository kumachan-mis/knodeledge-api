package usecase

import (
	"github.com/kumachan-mis/knodeledge-api/internal/domain"
	"github.com/kumachan-mis/knodeledge-api/internal/model"
	"github.com/kumachan-mis/knodeledge-api/internal/service"
)

//go:generate mockgen -source=$GOFILE -destination=../../mock/$GOPACKAGE/mock_$GOFILE -package=$GOPACKAGE

type ChapterUseCase interface {
	ListChapters(req model.ChapterListRequest) (
		*model.ChapterListResponse, *Error[model.ChapterListErrorResponse])
	CreateChapter(req model.ChapterCreateRequest) (
		*model.ChapterCreateResponse, *Error[model.ChapterCreateErrorResponse])
}

type chapterUseCase struct {
	service service.ChapterService
}

func NewChapterUseCase(service service.ChapterService) ChapterUseCase {
	return chapterUseCase{service: service}
}

func (uc chapterUseCase) ListChapters(req model.ChapterListRequest) (
	*model.ChapterListResponse, *Error[model.ChapterListErrorResponse]) {
	uid, uidErr := domain.NewUserIdObject(req.User.Id)
	pid, pidErr := domain.NewProjectIdObject(req.Project.Id)

	uidMsg := ""
	if uidErr != nil {
		uidMsg = uidErr.Error()
	}
	pidMsg := ""
	if pidErr != nil {
		pidMsg = pidErr.Error()
	}

	if uidErr != nil || pidErr != nil {
		return nil, NewModelBasedError(
			InvalidArgumentError,
			model.ChapterListErrorResponse{
				User:    model.UserOnlyIdError{Id: uidMsg},
				Project: model.ProjectOnlyIdError{Id: pidMsg},
			},
		)
	}

	entities, sErr := uc.service.ListChapters(*uid, *pid)
	if sErr != nil {
		return nil, NewMessageBasedError[model.ChapterListErrorResponse](
			InternalErrorPanic,
			sErr.Unwrap().Error(),
		)
	}

	chapters := make([]model.Chapter, len(entities))
	i := 0
	for _, entity := range entities {
		chapters[i] = model.Chapter{
			Id:       entity.Id().Value(),
			Name:     entity.Name().Value(),
			NextId:   entity.NextId().Value(),
			Sections: []model.Section{},
		}
		i++
	}

	return &model.ChapterListResponse{Chapters: chapters}, nil
}

func (uc chapterUseCase) CreateChapter(req model.ChapterCreateRequest) (
	*model.ChapterCreateResponse, *Error[model.ChapterCreateErrorResponse]) {
	uid, uidErr := domain.NewUserIdObject(req.User.Id)
	pid, pidErr := domain.NewProjectIdObject(req.Project.Id)
	name, nameErr := domain.NewChapterNameObject(req.Chapter.Name)
	nextId, nextIdErr := domain.NewChapterNextIdObject(req.Chapter.NextId)

	uidMsg := ""
	if uidErr != nil {
		uidMsg = uidErr.Error()
	}
	pidMsg := ""
	if pidErr != nil {
		pidMsg = pidErr.Error()
	}
	nameMsg := ""
	if nameErr != nil {
		nameMsg = nameErr.Error()
	}
	nextIdMsg := ""
	if nextIdErr != nil {
		nextIdMsg = nextIdErr.Error()
	}

	if uidErr != nil || pidErr != nil || nameErr != nil || nextIdErr != nil {
		return nil, NewModelBasedError(
			InvalidArgumentError,
			model.ChapterCreateErrorResponse{
				User: model.UserOnlyIdError{
					Id: uidMsg,
				},
				Project: model.ProjectOnlyIdError{
					Id: pidMsg,
				},
				Chapter: model.ChapterWithoutAutofieldError{
					Name:   nameMsg,
					NextId: nextIdMsg,
				},
			},
		)
	}

	chapter := domain.NewChapterWithoutAutofieldEntity(*name, *nextId)

	entity, sErr := uc.service.CreateChapter(*uid, *pid, *chapter)
	if sErr != nil && sErr.Code() == service.InvalidArgument {
		return nil, NewMessageBasedError[model.ChapterCreateErrorResponse](
			InvalidArgumentError,
			sErr.Unwrap().Error(),
		)
	}
	if sErr != nil {
		return nil, NewMessageBasedError[model.ChapterCreateErrorResponse](
			InternalErrorPanic,
			sErr.Unwrap().Error(),
		)
	}

	return &model.ChapterCreateResponse{
		Chapter: model.Chapter{
			Id:       entity.Id().Value(),
			Name:     entity.Name().Value(),
			NextId:   entity.NextId().Value(),
			Sections: []model.Section{},
		},
	}, nil
}
