package usecase

import (
	"github.com/kumachan-mis/knodeledge-api/internal/domain"
	"github.com/kumachan-mis/knodeledge-api/internal/openapi"
	"github.com/kumachan-mis/knodeledge-api/internal/service"
)

//go:generate mockgen -source=$GOFILE -destination=../../mock/$GOPACKAGE/mock_$GOFILE -package=$GOPACKAGE

type ChapterUseCase interface {
	ListChapters(req openapi.ChapterListRequest) (
		*openapi.ChapterListResponse, *Error[openapi.ChapterListErrorResponse])
	CreateChapter(req openapi.ChapterCreateRequest) (
		*openapi.ChapterCreateResponse, *Error[openapi.ChapterCreateErrorResponse])
	UpdateChapter(req openapi.ChapterUpdateRequest) (
		*openapi.ChapterUpdateResponse, *Error[openapi.ChapterUpdateErrorResponse])
	DeleteChapter(req openapi.ChapterDeleteRequest) *Error[openapi.ChapterDeleteErrorResponse]
}

type chapterUseCase struct {
	service service.ChapterService
}

func NewChapterUseCase(service service.ChapterService) ChapterUseCase {
	return chapterUseCase{service: service}
}

func (uc chapterUseCase) ListChapters(req openapi.ChapterListRequest) (
	*openapi.ChapterListResponse, *Error[openapi.ChapterListErrorResponse]) {
	userId, userIdErr := domain.NewUserIdObject(req.UserId)
	projectId, projectIdErr := domain.NewProjectIdObject(req.ProjectId)

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
			openapi.ChapterListErrorResponse{
				UserId:    userIdMsg,
				ProjectId: projectIdMsg,
			},
		)
	}

	entities, sErr := uc.service.ListChapters(*userId, *projectId)
	if sErr != nil && sErr.Code() == service.NotFoundError {
		return nil, NewMessageBasedError[openapi.ChapterListErrorResponse](
			NotFoundError,
			sErr.Unwrap().Error(),
		)
	}
	if sErr != nil {
		return nil, NewMessageBasedError[openapi.ChapterListErrorResponse](
			InternalErrorPanic,
			sErr.Unwrap().Error(),
		)
	}

	chapters := make([]openapi.ChapterWithSections, len(entities))
	for i, entity := range entities {
		sections := make([]openapi.SectionOfChapter, len(entity.Sections()))
		for j, section := range entity.Sections() {
			sections[j] = openapi.SectionOfChapter{
				Id:   section.Id().Value(),
				Name: section.Name().Value(),
			}
		}

		chapters[i] = openapi.ChapterWithSections{
			Id:       entity.Id().Value(),
			Name:     entity.Name().Value(),
			Number:   int32(entity.Number().Value()),
			Sections: sections,
		}
	}

	return &openapi.ChapterListResponse{Chapters: chapters}, nil
}

func (uc chapterUseCase) CreateChapter(req openapi.ChapterCreateRequest) (
	*openapi.ChapterCreateResponse, *Error[openapi.ChapterCreateErrorResponse]) {
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
			openapi.ChapterCreateErrorResponse{
				User: openapi.UserOnlyIdError{
					Id: userIdMsg,
				},
				Project: openapi.ProjectOnlyIdError{
					Id: projectIdMsg,
				},
				Chapter: openapi.ChapterWithoutAutofieldError{
					Name:   chapterNameMsg,
					Number: chapterNumberMsg,
				},
			},
		)
	}

	chapter := domain.NewChapterWithoutAutofieldEntity(*chapterName, *chapterNumber)

	chapterEntity, sErr := uc.service.CreateChapter(*userId, *projectId, *chapter)
	if sErr != nil && sErr.Code() == service.InvalidArgumentError {
		return nil, NewMessageBasedError[openapi.ChapterCreateErrorResponse](
			InvalidArgumentError,
			sErr.Unwrap().Error(),
		)
	}
	if sErr != nil && sErr.Code() == service.NotFoundError {
		return nil, NewMessageBasedError[openapi.ChapterCreateErrorResponse](
			NotFoundError,
			sErr.Unwrap().Error(),
		)
	}
	if sErr != nil {
		return nil, NewMessageBasedError[openapi.ChapterCreateErrorResponse](
			InternalErrorPanic,
			sErr.Unwrap().Error(),
		)
	}

	return &openapi.ChapterCreateResponse{
		Chapter: openapi.ChapterWithSections{
			Id:       chapterEntity.Id().Value(),
			Name:     chapterEntity.Name().Value(),
			Number:   int32(chapterEntity.Number().Value()),
			Sections: []openapi.SectionOfChapter{},
		},
	}, nil
}

func (uc chapterUseCase) UpdateChapter(req openapi.ChapterUpdateRequest) (
	*openapi.ChapterUpdateResponse, *Error[openapi.ChapterUpdateErrorResponse]) {
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
			openapi.ChapterUpdateErrorResponse{
				User: openapi.UserOnlyIdError{
					Id: userIdMsg,
				},
				Project: openapi.ProjectOnlyIdError{
					Id: projectIdMsg,
				},
				Chapter: openapi.ChapterError{
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
		return nil, NewMessageBasedError[openapi.ChapterUpdateErrorResponse](
			InvalidArgumentError,
			sErr.Unwrap().Error(),
		)
	}
	if sErr != nil && sErr.Code() == service.NotFoundError {
		return nil, NewMessageBasedError[openapi.ChapterUpdateErrorResponse](
			NotFoundError,
			sErr.Unwrap().Error(),
		)
	}
	if sErr != nil {
		return nil, NewMessageBasedError[openapi.ChapterUpdateErrorResponse](
			InternalErrorPanic,
			sErr.Unwrap().Error(),
		)
	}

	return &openapi.ChapterUpdateResponse{
		Chapter: openapi.ChapterWithSections{
			Id:       entity.Id().Value(),
			Name:     entity.Name().Value(),
			Number:   int32(entity.Number().Value()),
			Sections: []openapi.SectionOfChapter{},
		},
	}, nil
}

func (uc chapterUseCase) DeleteChapter(req openapi.ChapterDeleteRequest) *Error[openapi.ChapterDeleteErrorResponse] {
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
			openapi.ChapterDeleteErrorResponse{
				User: openapi.UserOnlyIdError{
					Id: userIdMsg,
				},
				Project: openapi.ProjectOnlyIdError{
					Id: projectIdMsg,
				},
				Chapter: openapi.ChapterOnlyIdError{
					Id: chapterIdMsg,
				},
			},
		)
	}

	sErr := uc.service.DeleteChapter(*userId, *projectId, *chapterId)
	if sErr != nil && sErr.Code() == service.NotFoundError {
		return NewMessageBasedError[openapi.ChapterDeleteErrorResponse](
			NotFoundError,
			sErr.Unwrap().Error(),
		)
	}
	if sErr != nil {
		return NewMessageBasedError[openapi.ChapterDeleteErrorResponse](
			InternalErrorPanic,
			sErr.Unwrap().Error(),
		)
	}

	return nil
}
