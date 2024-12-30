package usecase

import (
	"github.com/kumachan-mis/knodeledge-api/internal/domain"
	"github.com/kumachan-mis/knodeledge-api/internal/model"
	"github.com/kumachan-mis/knodeledge-api/internal/service"
)

//go:generate mockgen -source=$GOFILE -destination=../../mock/$GOPACKAGE/mock_$GOFILE -package=$GOPACKAGE

type GraphUseCase interface {
	FindGraph(request model.GraphFindRequest) (
		*model.GraphFindResponse, *Error[model.GraphFindErrorResponse])
	UpdateGraph(request model.GraphUpdateRequest) (
		*model.GraphUpdateResponse, *Error[model.GraphUpdateErrorResponse])
	SectionalizeGraph(request model.GraphSectionalizeRequest) (
		*model.GraphSectionalizeResponse, *Error[model.GraphSectionalizeErrorResponse])
}

type graphUseCase struct {
	service service.GraphService
}

func NewGraphUseCase(service service.GraphService) GraphUseCase {
	return graphUseCase{service: service}
}

func (uc graphUseCase) FindGraph(req model.GraphFindRequest) (
	*model.GraphFindResponse, *Error[model.GraphFindErrorResponse]) {
	userId, userIdErr := domain.NewUserIdObject(req.User.Id)
	projectId, projectIdErr := domain.NewProjectIdObject(req.Project.Id)
	chapterId, chapterIdErr := domain.NewChapterIdObject(req.Chapter.Id)
	sectionId, sectionIdErr := domain.NewSectionIdObject(req.Section.Id)

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
	sectionIdMsg := ""
	if sectionIdErr != nil {
		sectionIdMsg = sectionIdErr.Error()
	}

	if userIdErr != nil || projectIdErr != nil || chapterIdErr != nil || sectionIdErr != nil {
		return nil, NewModelBasedError(
			DomainValidationError,
			model.GraphFindErrorResponse{
				User:    model.UserOnlyIdError{Id: userIdMsg},
				Project: model.ProjectOnlyIdError{Id: projectIdMsg},
				Chapter: model.ChapterOnlyIdError{Id: chapterIdMsg},
				Section: model.SectionOnlyIdError{Id: sectionIdMsg},
			},
		)
	}

	entity, fErr := uc.service.FindGraph(*userId, *projectId, *chapterId, *sectionId)

	if fErr != nil && fErr.Code() == service.NotFoundError {
		return nil, NewMessageBasedError[model.GraphFindErrorResponse](
			NotFoundError,
			fErr.Unwrap().Error(),
		)
	}
	if fErr != nil {
		return nil, NewMessageBasedError[model.GraphFindErrorResponse](
			InternalErrorPanic,
			fErr.Unwrap().Error(),
		)
	}

	return &model.GraphFindResponse{
		Graph: model.Graph{
			Id:        entity.Id().Value(),
			Name:      entity.Name().Value(),
			Paragraph: entity.Paragraph().Value(),
		},
	}, nil
}

func (uc graphUseCase) UpdateGraph(req model.GraphUpdateRequest) (
	*model.GraphUpdateResponse, *Error[model.GraphUpdateErrorResponse]) {
	userId, userIdErr := domain.NewUserIdObject(req.User.Id)
	projectId, projectIdErr := domain.NewProjectIdObject(req.Project.Id)
	chapterId, chapterIdErr := domain.NewChapterIdObject(req.Chapter.Id)
	graphId, graphIdErr := domain.NewGraphIdObject(req.Graph.Id)
	graphParagraph, graphParagraphErr := domain.NewGraphParagraphObject(req.Graph.Paragraph)

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
	graphIdMsg := ""
	if graphIdErr != nil {
		graphIdMsg = graphIdErr.Error()
	}
	graphParagraphMsg := ""
	if graphParagraphErr != nil {
		graphParagraphMsg = graphParagraphErr.Error()
	}

	if userIdErr != nil || projectIdErr != nil || chapterIdErr != nil || graphIdErr != nil || graphParagraphErr != nil {
		return nil, NewModelBasedError(
			DomainValidationError,
			model.GraphUpdateErrorResponse{
				User:    model.UserOnlyIdError{Id: userIdMsg},
				Project: model.ProjectOnlyIdError{Id: projectIdMsg},
				Chapter: model.ChapterOnlyIdError{Id: chapterIdMsg},
				Graph:   model.GraphContentError{Id: graphIdMsg, Paragraph: graphParagraphMsg},
			},
		)
	}

	graph := domain.NewGraphContentWithoutAutofieldEntity(*graphParagraph)

	entity, uErr := uc.service.UpdateGraphContent(*userId, *projectId, *chapterId, *graphId, *graph)

	if uErr != nil && uErr.Code() == service.NotFoundError {
		return nil, NewMessageBasedError[model.GraphUpdateErrorResponse](
			NotFoundError,
			uErr.Unwrap().Error(),
		)
	}
	if uErr != nil {
		return nil, NewMessageBasedError[model.GraphUpdateErrorResponse](
			InternalErrorPanic,
			uErr.Unwrap().Error(),
		)
	}

	return &model.GraphUpdateResponse{
		Graph: model.Graph{
			Id:        entity.Id().Value(),
			Name:      entity.Name().Value(),
			Paragraph: entity.Paragraph().Value(),
		},
	}, nil
}

func (uc graphUseCase) SectionalizeGraph(req model.GraphSectionalizeRequest) (
	*model.GraphSectionalizeResponse, *Error[model.GraphSectionalizeErrorResponse]) {
	userId, userIdErr := domain.NewUserIdObject(req.User.Id)
	projectId, projectIdErr := domain.NewProjectIdObject(req.Project.Id)
	chapterId, chapterIdErr := domain.NewChapterIdObject(req.Chapter.Id)

	sectionItems := make([]domain.SectionWithoutAutofieldEntity, len(req.Sections))
	sectionItemErrors := make([]model.SectionWithoutAutofieldError, len(req.Sections))

	ectionItemErrorExists := false
	for i, section := range req.Sections {
		sectionName, sectionNameErr := domain.NewSectionNameObject(section.Name)
		if sectionNameErr != nil {
			sectionItemErrors[i].Name = sectionNameErr.Error()
			ectionItemErrorExists = true
		}
		sectionContent, sectionContentErr := domain.NewSectionContentObject(section.Content)
		if sectionContentErr != nil {
			sectionItemErrors[i].Content = sectionContentErr.Error()
			ectionItemErrorExists = true
		}
		if sectionNameErr == nil && sectionContentErr == nil {
			sectionItems[i] = *domain.NewSectionWithoutAutofieldEntity(*sectionName, *sectionContent)
		}
	}
	sections, sectionsErr := domain.NewSectionWithoutAutofieldEntityList(sectionItems)

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
	sectionsMsg := ""
	if sectionsErr != nil {
		sectionsMsg = sectionsErr.Error()
	}

	if userIdErr != nil || projectIdErr != nil || chapterIdErr != nil || ectionItemErrorExists || sectionsErr != nil {
		return nil, NewModelBasedError(
			DomainValidationError,
			model.GraphSectionalizeErrorResponse{
				User:    model.UserOnlyIdError{Id: userIdMsg},
				Project: model.ProjectOnlyIdError{Id: projectIdMsg},
				Chapter: model.ChapterOnlyIdError{Id: chapterIdMsg},
				Sections: model.SectionWithoutAutofieldListError{
					Message: sectionsMsg,
					Items:   sectionItemErrors,
				},
			},
		)
	}

	entities, sErr := uc.service.SectionalizeIntoGraphs(*userId, *projectId, *chapterId, *sections)
	if sErr != nil && sErr.Code() == service.InvalidArgument {
		return nil, NewMessageBasedError[model.GraphSectionalizeErrorResponse](
			InvalidArgumentError,
			sErr.Unwrap().Error(),
		)
	}
	if sErr != nil && sErr.Code() == service.NotFoundError {
		return nil, NewMessageBasedError[model.GraphSectionalizeErrorResponse](
			NotFoundError,
			sErr.Unwrap().Error(),
		)
	}
	if sErr != nil {
		return nil, NewMessageBasedError[model.GraphSectionalizeErrorResponse](
			InternalErrorPanic,
			sErr.Unwrap().Error(),
		)
	}

	graphs := make([]model.Graph, len(entities))
	for i, entity := range entities {
		graphs[i] = model.Graph{
			Id:        entity.Id().Value(),
			Name:      entity.Name().Value(),
			Paragraph: entity.Paragraph().Value(),
		}
	}

	return &model.GraphSectionalizeResponse{Graphs: graphs}, nil
}
