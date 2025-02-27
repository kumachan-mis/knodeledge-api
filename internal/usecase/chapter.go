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
	UpdateChapter(req model.ChapterUpdateRequest) (
		*model.ChapterUpdateResponse, *Error[model.ChapterUpdateErrorResponse])
	DeleteChapter(req model.ChapterDeleteRequest) *Error[model.ChapterDeleteErrorResponse]
}

type chapterUseCase struct {
	service service.ChapterService
}

func NewChapterUseCase(service service.ChapterService) ChapterUseCase {
	return chapterUseCase{service: service}
}

func (uc chapterUseCase) ListChapters(req model.ChapterListRequest) (
	*model.ChapterListResponse, *Error[model.ChapterListErrorResponse]) {
	userId, userIdErr := domain.NewUserIdObject(req.User.Id)
	projectId, projectIdErr := domain.NewProjectIdObject(req.Project.Id)

	userIdMsg := ""
	if userIdErr != nil {
		userIdMsg = userIdErr.Error()
	}
	projectIdMsg := ""
	if projectIdErr != nil {
		projectIdMsg = projectIdErr.Error()
	}

	if userIdErr != nil || projectIdErr != nil {
		return nil, NewModelBasedError(
			DomainValidationError,
			model.ChapterListErrorResponse{
				User:    model.UserOnlyIdError{Id: userIdMsg},
				Project: model.ProjectOnlyIdError{Id: projectIdMsg},
			},
		)
	}

	entities, sErr := uc.service.ListChapters(*userId, *projectId)
	if sErr != nil && sErr.Code() == service.NotFoundError {
		return nil, NewMessageBasedError[model.ChapterListErrorResponse](
			NotFoundError,
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
	for i, entity := range entities {
		sections := make([]model.SectionOfChapter, len(entity.Sections()))
		for j, section := range entity.Sections() {
			sections[j] = model.SectionOfChapter{
				Id:   section.Id().Value(),
				Name: section.Name().Value(),
			}
		}

		chapters[i] = model.ChapterWithSections{
			Id:       entity.Id().Value(),
			Name:     entity.Name().Value(),
			Number:   int32(entity.Number().Value()),
			Sections: sections,
		}
	}

	return &model.ChapterListResponse{Chapters: chapters}, nil
}

func (uc chapterUseCase) CreateChapter(req model.ChapterCreateRequest) (
	*model.ChapterCreateResponse, *Error[model.ChapterCreateErrorResponse]) {
	userId, userIdErr := domain.NewUserIdObject(req.User.Id)
	projectId, projectIdErr := domain.NewProjectIdObject(req.Project.Id)
	chapterName, chapterNameErr := domain.NewChapterNameObject(req.Chapter.Name)
	chapterNumber, chapterNumberErr := domain.NewChapterNumberObject(int(req.Chapter.Number))

	userIdMsg := ""
	if userIdErr != nil {
		userIdMsg = userIdErr.Error()
	}
	projectIdMsg := ""
	if projectIdErr != nil {
		projectIdMsg = projectIdErr.Error()
	}
	chapterNameMsg := ""
	if chapterNameErr != nil {
		chapterNameMsg = chapterNameErr.Error()
	}
	chapterNumberMsg := ""
	if chapterNumberErr != nil {
		chapterNumberMsg = chapterNumberErr.Error()
	}

	if userIdErr != nil || projectIdErr != nil ||
		chapterNameErr != nil || chapterNumberErr != nil {
		return nil, NewModelBasedError(
			DomainValidationError,
			model.ChapterCreateErrorResponse{
				User: model.UserOnlyIdError{
					Id: userIdMsg,
				},
				Project: model.ProjectOnlyIdError{
					Id: projectIdMsg,
				},
				Chapter: model.ChapterWithoutAutofieldError{
					Name:   chapterNameMsg,
					Number: chapterNumberMsg,
				},
			},
		)
	}

	chapter := domain.NewChapterWithoutAutofieldEntity(*chapterName, *chapterNumber)

	chapterEntity, sErr := uc.service.CreateChapter(*userId, *projectId, *chapter)
	if sErr != nil && sErr.Code() == service.InvalidArgumentError {
		return nil, NewMessageBasedError[model.ChapterCreateErrorResponse](
			InvalidArgumentError,
			sErr.Unwrap().Error(),
		)
	}
	if sErr != nil && sErr.Code() == service.NotFoundError {
		return nil, NewMessageBasedError[model.ChapterCreateErrorResponse](
			NotFoundError,
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
			Id:       chapterEntity.Id().Value(),
			Name:     chapterEntity.Name().Value(),
			Number:   int32(chapterEntity.Number().Value()),
			Sections: []model.SectionOfChapter{},
		},
	}, nil
}

func (uc chapterUseCase) UpdateChapter(req model.ChapterUpdateRequest) (
	*model.ChapterUpdateResponse, *Error[model.ChapterUpdateErrorResponse]) {
	userId, userIdErr := domain.NewUserIdObject(req.User.Id)
	projectId, projectIdErr := domain.NewProjectIdObject(req.Project.Id)
	chapterId, chapterIdErr := domain.NewChapterIdObject(req.Chapter.Id)
	chapterName, chapterNameErr := domain.NewChapterNameObject(req.Chapter.Name)
	chapterNumber, chapterNumberErr := domain.NewChapterNumberObject(int(req.Chapter.Number))

	userIdMsg := ""
	if userIdErr != nil {
		userIdMsg = userIdErr.Error()
	}
	projectIdMsg := ""
	if projectIdErr != nil {
		projectIdMsg = projectIdErr.Error()
	}
	chapterIdMsg := ""
	if chapterIdErr != nil {
		chapterIdMsg = chapterIdErr.Error()
	}
	chapterNameMsg := ""
	if chapterNameErr != nil {
		chapterNameMsg = chapterNameErr.Error()
	}
	chapterNumberMsg := ""
	if chapterNumberErr != nil {
		chapterNumberMsg = chapterNumberErr.Error()
	}

	if userIdErr != nil || projectIdErr != nil || chapterIdErr != nil || chapterNameErr != nil || chapterNumberErr != nil {
		return nil, NewModelBasedError(
			DomainValidationError,
			model.ChapterUpdateErrorResponse{
				User: model.UserOnlyIdError{
					Id: userIdMsg,
				},
				Project: model.ProjectOnlyIdError{
					Id: projectIdMsg,
				},
				Chapter: model.ChapterError{
					Id:     chapterIdMsg,
					Name:   chapterNameMsg,
					Number: chapterNumberMsg,
				},
			},
		)
	}

	chapter := domain.NewChapterWithoutAutofieldEntity(*chapterName, *chapterNumber)

	entity, sErr := uc.service.UpdateChapter(*userId, *projectId, *chapterId, *chapter)
	if sErr != nil && sErr.Code() == service.InvalidArgumentError {
		return nil, NewMessageBasedError[model.ChapterUpdateErrorResponse](
			InvalidArgumentError,
			sErr.Unwrap().Error(),
		)
	}
	if sErr != nil && sErr.Code() == service.NotFoundError {
		return nil, NewMessageBasedError[model.ChapterUpdateErrorResponse](
			NotFoundError,
			sErr.Unwrap().Error(),
		)
	}
	if sErr != nil {
		return nil, NewMessageBasedError[model.ChapterUpdateErrorResponse](
			InternalErrorPanic,
			sErr.Unwrap().Error(),
		)
	}

	return &model.ChapterUpdateResponse{
		Chapter: model.ChapterWithSections{
			Id:       entity.Id().Value(),
			Name:     entity.Name().Value(),
			Number:   int32(entity.Number().Value()),
			Sections: []model.SectionOfChapter{},
		},
	}, nil
}

func (uc chapterUseCase) DeleteChapter(req model.ChapterDeleteRequest) *Error[model.ChapterDeleteErrorResponse] {
	userId, userIdErr := domain.NewUserIdObject(req.User.Id)
	projectId, projectIdErr := domain.NewProjectIdObject(req.Project.Id)
	chapterId, chapterIdErr := domain.NewChapterIdObject(req.Chapter.Id)

	userIdMsg := ""
	if userIdErr != nil {
		userIdMsg = userIdErr.Error()
	}
	projectIdMsg := ""
	if projectIdErr != nil {
		projectIdMsg = projectIdErr.Error()
	}
	chapterIdMsg := ""
	if chapterIdErr != nil {
		chapterIdMsg = chapterIdErr.Error()
	}

	if userIdErr != nil || projectIdErr != nil || chapterIdErr != nil {
		return NewModelBasedError(
			DomainValidationError,
			model.ChapterDeleteErrorResponse{
				User: model.UserOnlyIdError{
					Id: userIdMsg,
				},
				Project: model.ProjectOnlyIdError{
					Id: projectIdMsg,
				},
				Chapter: model.ChapterOnlyIdError{
					Id: chapterIdMsg,
				},
			},
		)
	}

	sErr := uc.service.DeleteChapter(*userId, *projectId, *chapterId)
	if sErr != nil && sErr.Code() == service.NotFoundError {
		return NewMessageBasedError[model.ChapterDeleteErrorResponse](
			NotFoundError,
			sErr.Unwrap().Error(),
		)
	}
	if sErr != nil {
		return NewMessageBasedError[model.ChapterDeleteErrorResponse](
			InternalErrorPanic,
			sErr.Unwrap().Error(),
		)
	}

	return nil
}
