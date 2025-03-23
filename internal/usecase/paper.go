package usecase

import (
	"github.com/kumachan-mis/knodeledge-api/internal/domain"
	"github.com/kumachan-mis/knodeledge-api/internal/openapi"
	"github.com/kumachan-mis/knodeledge-api/internal/service"
)

//go:generate mockgen -source=$GOFILE -destination=../../mock/$GOPACKAGE/mock_$GOFILE -package=$GOPACKAGE

type PaperUseCase interface {
	FindPaper(request openapi.PaperFindRequest) (
		*openapi.PaperFindResponse, *Error[openapi.PaperFindErrorResponse])
	UpdatePaper(request openapi.PaperUpdateRequest) (
		*openapi.PaperUpdateResponse, *Error[openapi.PaperUpdateErrorResponse])
}

type paperUseCase struct {
	service service.PaperService
}

func NewPaperUseCase(service service.PaperService) PaperUseCase {
	return paperUseCase{service: service}
}

func (uc paperUseCase) FindPaper(req openapi.PaperFindRequest) (
	*openapi.PaperFindResponse, *Error[openapi.PaperFindErrorResponse]) {
	userId, userIdErr := domain.NewUserIdObject(req.UserId)
	projectId, projectIdErr := domain.NewProjectIdObject(req.ProjectId)
	chapterId, chapterIdErr := domain.NewChapterIdObject(req.ChapterId)

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
		return nil, NewModelBasedError(
			DomainValidationError,
			openapi.PaperFindErrorResponse{
				UserId:    userIdMsg,
				ProjectId: projectIdMsg,
				ChapterId: chapterIdMsg,
			},
		)
	}

	entity, sErr := uc.service.FindPaper(*userId, *projectId, *chapterId)

	if sErr != nil && sErr.Code() == service.NotFoundError {
		return nil, NewMessageBasedError[openapi.PaperFindErrorResponse](
			NotFoundError,
			sErr.Unwrap().Error(),
		)
	}
	if sErr != nil {
		return nil, NewMessageBasedError[openapi.PaperFindErrorResponse](
			InternalErrorPanic,
			sErr.Unwrap().Error(),
		)
	}

	return &openapi.PaperFindResponse{
		Paper: openapi.Paper{
			Id:      entity.Id().Value(),
			Content: entity.Content().Value(),
		},
	}, nil
}

func (uc paperUseCase) UpdatePaper(req openapi.PaperUpdateRequest) (
	*openapi.PaperUpdateResponse, *Error[openapi.PaperUpdateErrorResponse]) {
	userId, userIdErr := domain.NewUserIdObject(req.User.Id)
	projectId, projectIdErr := domain.NewProjectIdObject(req.Project.Id)
	paperId, paperIdErr := domain.NewPaperIdObject(req.Paper.Id)
	paperContent, paperContentErr := domain.NewPaperContentObject(req.Paper.Content)

	userIdMsg := ""
	if userIdErr != nil {
		userIdMsg = userIdErr.Error()
	}
	projectIdMsg := ""
	if projectIdErr != nil {
		projectIdMsg = projectIdErr.Error()
	}
	paperIdMsg := ""
	if paperIdErr != nil {
		paperIdMsg = paperIdErr.Error()
	}
	paperContentMsg := ""
	if paperContentErr != nil {
		paperContentMsg = paperContentErr.Error()
	}

	if userIdErr != nil || projectIdErr != nil || paperIdErr != nil || paperContentErr != nil {
		return nil, NewModelBasedError(
			DomainValidationError,
			openapi.PaperUpdateErrorResponse{
				User:    openapi.UserOnlyIdError{Id: userIdMsg},
				Project: openapi.ProjectOnlyIdError{Id: projectIdMsg},
				Paper:   openapi.PaperError{Id: paperIdMsg, Content: paperContentMsg},
			},
		)
	}

	paper := domain.NewPaperWithoutAutofieldEntity(*paperContent)

	entity, sErr := uc.service.UpdatePaper(*userId, *projectId, *paperId, *paper)

	if sErr != nil && sErr.Code() == service.NotFoundError {
		return nil, NewMessageBasedError[openapi.PaperUpdateErrorResponse](
			NotFoundError,
			sErr.Unwrap().Error(),
		)
	}
	if sErr != nil {
		return nil, NewMessageBasedError[openapi.PaperUpdateErrorResponse](
			InternalErrorPanic,
			sErr.Unwrap().Error(),
		)
	}

	return &openapi.PaperUpdateResponse{
		Paper: openapi.Paper{
			Id:      entity.Id().Value(),
			Content: entity.Content().Value(),
		},
	}, nil
}
