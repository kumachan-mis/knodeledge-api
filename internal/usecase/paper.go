package usecase

import (
	"github.com/kumachan-mis/knodeledge-api/internal/domain"
	"github.com/kumachan-mis/knodeledge-api/internal/model"
	"github.com/kumachan-mis/knodeledge-api/internal/service"
)

//go:generate mockgen -source=$GOFILE -destination=../../mock/$GOPACKAGE/mock_$GOFILE -package=$GOPACKAGE

type PaperUseCase interface {
	FindPaper(request model.PaperFindRequest) (
		*model.PaperFindResponse, *Error[model.PaperFindErrorResponse])
	UpdatePaper(request model.PaperUpdateRequest) (
		*model.PaperUpdateResponse, *Error[model.PaperUpdateErrorResponse])
}

type paperUseCase struct {
	service service.PaperService
}

func NewPaperUseCase(service service.PaperService) PaperUseCase {
	return paperUseCase{service: service}
}

func (uc paperUseCase) FindPaper(req model.PaperFindRequest) (
	*model.PaperFindResponse, *Error[model.PaperFindErrorResponse]) {
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
		return nil, NewModelBasedError(
			DomainValidationError,
			model.PaperFindErrorResponse{
				User:    model.UserOnlyIdError{Id: userIdMsg},
				Project: model.ProjectOnlyIdError{Id: projectIdMsg},
				Chapter: model.ChapterOnlyIdError{Id: chapterIdMsg},
			},
		)
	}

	entity, sErr := uc.service.FindPaper(*userId, *projectId, *chapterId)

	if sErr != nil && sErr.Code() == service.NotFoundError {
		return nil, NewMessageBasedError[model.PaperFindErrorResponse](
			NotFoundError,
			sErr.Unwrap().Error(),
		)
	}
	if sErr != nil {
		return nil, NewMessageBasedError[model.PaperFindErrorResponse](
			InternalErrorPanic,
			sErr.Unwrap().Error(),
		)
	}

	return &model.PaperFindResponse{
		Paper: model.Paper{
			Id:      entity.Id().Value(),
			Content: entity.Content().Value(),
		},
	}, nil
}

func (uc paperUseCase) UpdatePaper(req model.PaperUpdateRequest) (
	*model.PaperUpdateResponse, *Error[model.PaperUpdateErrorResponse]) {
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
			model.PaperUpdateErrorResponse{
				User:    model.UserOnlyIdError{Id: userIdMsg},
				Project: model.ProjectOnlyIdError{Id: projectIdMsg},
				Paper:   model.PaperError{Id: paperIdMsg, Content: paperContentMsg},
			},
		)
	}

	paper := domain.NewPaperWithoutAutofieldEntity(*paperContent)

	entity, sErr := uc.service.UpdatePaper(*userId, *projectId, *paperId, *paper)

	if sErr != nil && sErr.Code() == service.NotFoundError {
		return nil, NewMessageBasedError[model.PaperUpdateErrorResponse](
			NotFoundError,
			sErr.Unwrap().Error(),
		)
	}
	if sErr != nil {
		return nil, NewMessageBasedError[model.PaperUpdateErrorResponse](
			InternalErrorPanic,
			sErr.Unwrap().Error(),
		)
	}

	return &model.PaperUpdateResponse{
		Paper: model.Paper{
			Id:      entity.Id().Value(),
			Content: entity.Content().Value(),
		},
	}, nil
}
