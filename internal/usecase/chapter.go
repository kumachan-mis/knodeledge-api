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
			DomainValidationError,
			model.ChapterListErrorResponse{
				User:    model.UserOnlyIdError{Id: uidMsg},
				Project: model.ProjectOnlyIdError{Id: pidMsg},
			},
		)
	}

	entities, sErr := uc.service.ListChapters(*uid, *pid)
	if sErr != nil && sErr.Code() == service.InvalidArgument {
		return nil, NewMessageBasedError[model.ChapterListErrorResponse](
			InvalidArgumentError,
			sErr.Unwrap().Error(),
		)
	}
	if sErr != nil {
		return nil, NewMessageBasedError[model.ChapterListErrorResponse](
			InternalErrorPanic,
			sErr.Unwrap().Error(),
		)
	}

	chapters := make([]model.ChapterWithSections, len(entities))
	i := 0
	for _, entity := range entities {
		chapters[i] = model.ChapterWithSections{
			Id:       entity.Id().Value(),
			Name:     entity.Name().Value(),
			Number:   int32(entity.Number().Value()),
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
	cname, cnameErr := domain.NewChapterNameObject(req.Chapter.Name)
	cnumber, cnumberErr := domain.NewChapterNumberObject(int(req.Chapter.Number))

	uidMsg := ""
	if uidErr != nil {
		uidMsg = uidErr.Error()
	}
	pidMsg := ""
	if pidErr != nil {
		pidMsg = pidErr.Error()
	}
	cnameMsg := ""
	if cnameErr != nil {
		cnameMsg = cnameErr.Error()
	}
	cnumberMsg := ""
	if cnumberErr != nil {
		cnumberMsg = cnumberErr.Error()
	}

	if uidErr != nil || pidErr != nil || cnameErr != nil || cnumberErr != nil {
		return nil, NewModelBasedError(
			DomainValidationError,
			model.ChapterCreateErrorResponse{
				User: model.UserOnlyIdError{
					Id: uidMsg,
				},
				Project: model.ProjectOnlyIdError{
					Id: pidMsg,
				},
				Chapter: model.ChapterWithoutAutofieldError{
					Name:   cnameMsg,
					Number: cnumberMsg,
				},
			},
		)
	}

	chapter := domain.NewChapterWithoutAutofieldEntity(*cname, *cnumber)

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
		Chapter: model.ChapterWithSections{
			Id:       entity.Id().Value(),
			Name:     entity.Name().Value(),
			Number:   int32(entity.Number().Value()),
			Sections: []model.Section{},
		},
	}, nil
}
