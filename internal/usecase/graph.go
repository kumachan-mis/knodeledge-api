package usecase

import (
	"github.com/kumachan-mis/knodeledge-api/internal/domain"
	"github.com/kumachan-mis/knodeledge-api/internal/openapi"
	"github.com/kumachan-mis/knodeledge-api/internal/service"
)

//go:generate mockgen -source=$GOFILE -destination=../../mock/$GOPACKAGE/mock_$GOFILE -package=$GOPACKAGE

type GraphUseCase interface {
	FindGraph(request openapi.GraphFindRequest) (
		*openapi.GraphFindResponse, *Error[openapi.GraphFindErrorResponse])
	UpdateGraph(request openapi.GraphUpdateRequest) (
		*openapi.GraphUpdateResponse, *Error[openapi.GraphUpdateErrorResponse])
	DeleteGraph(request openapi.GraphDeleteRequest) *Error[openapi.GraphDeleteErrorResponse]
	SectionalizeGraph(request openapi.GraphSectionalizeRequest) (
		*openapi.GraphSectionalizeResponse, *Error[openapi.GraphSectionalizeErrorResponse])
}

type graphUseCase struct {
	service service.GraphService
}

func NewGraphUseCase(service service.GraphService) GraphUseCase {
	return graphUseCase{service: service}
}

func (uc graphUseCase) FindGraph(req openapi.GraphFindRequest) (
	*openapi.GraphFindResponse, *Error[openapi.GraphFindErrorResponse]) {
	userId, userIdErr := domain.NewUserIdObject(req.UserId)
	projectId, projectIdErr := domain.NewProjectIdObject(req.ProjectId)
	chapterId, chapterIdErr := domain.NewChapterIdObject(req.ChapterId)
	sectionId, sectionIdErr := domain.NewSectionIdObject(req.SectionId)

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
			openapi.GraphFindErrorResponse{
				UserId:    userIdMsg,
				ProjectId: projectIdMsg,
				ChapterId: chapterIdMsg,
				SectionId: sectionIdMsg,
			},
		)
	}

	entity, sErr := uc.service.FindGraph(*userId, *projectId, *chapterId, *sectionId)

	if sErr != nil && sErr.Code() == service.NotFoundError {
		return nil, NewMessageBasedError[openapi.GraphFindErrorResponse](
			NotFoundError,
			sErr.Unwrap().Error(),
		)
	}
	if sErr != nil {
		return nil, NewMessageBasedError[openapi.GraphFindErrorResponse](
			InternalErrorPanic,
			sErr.Unwrap().Error(),
		)
	}

	return &openapi.GraphFindResponse{
		Graph: openapi.Graph{
			Id:        entity.Id().Value(),
			Name:      entity.Name().Value(),
			Paragraph: entity.Paragraph().Value(),
			Children:  uc.childrenEntityToModel(entity.Children()),
		},
	}, nil
}

func (uc graphUseCase) UpdateGraph(req openapi.GraphUpdateRequest) (
	*openapi.GraphUpdateResponse, *Error[openapi.GraphUpdateErrorResponse]) {
	userId, userIdErr := domain.NewUserIdObject(req.User.Id)
	projectId, projectIdErr := domain.NewProjectIdObject(req.Project.Id)
	chapterId, chapterIdErr := domain.NewChapterIdObject(req.Chapter.Id)
	graphId, graphIdErr := domain.NewGraphIdObject(req.Graph.Id)
	graphParagraph, graphParagraphErr := domain.NewGraphParagraphObject(req.Graph.Paragraph)
	graphChildren, graphChildrenErr, graphChildrenOk := uc.childrenModelToEntity(req.Graph.Children)

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

	if userIdErr != nil || projectIdErr != nil || chapterIdErr != nil ||
		graphIdErr != nil || graphParagraphErr != nil || !graphChildrenOk {
		return nil, NewModelBasedError(
			DomainValidationError,
			openapi.GraphUpdateErrorResponse{
				User:    openapi.UserOnlyIdError{Id: userIdMsg},
				Project: openapi.ProjectOnlyIdError{Id: projectIdMsg},
				Chapter: openapi.ChapterOnlyIdError{Id: chapterIdMsg},
				Graph: openapi.GraphContentError{
					Id:        graphIdMsg,
					Paragraph: graphParagraphMsg,
					Children:  *graphChildrenErr,
				},
			},
		)
	}

	graph := domain.NewGraphContentEntity(*graphParagraph, *graphChildren)

	entity, uErr := uc.service.UpdateGraphContent(*userId, *projectId, *chapterId, *graphId, *graph)

	if uErr != nil && uErr.Code() == service.NotFoundError {
		return nil, NewMessageBasedError[openapi.GraphUpdateErrorResponse](
			NotFoundError,
			uErr.Unwrap().Error(),
		)
	}
	if uErr != nil {
		return nil, NewMessageBasedError[openapi.GraphUpdateErrorResponse](
			InternalErrorPanic,
			uErr.Unwrap().Error(),
		)
	}

	return &openapi.GraphUpdateResponse{
		Graph: openapi.Graph{
			Id:        entity.Id().Value(),
			Name:      entity.Name().Value(),
			Paragraph: entity.Paragraph().Value(),
			Children:  uc.childrenEntityToModel(entity.Children()),
		},
	}, nil
}

func (uc graphUseCase) DeleteGraph(req openapi.GraphDeleteRequest) *Error[openapi.GraphDeleteErrorResponse] {
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
		return NewModelBasedError(
			DomainValidationError,
			openapi.GraphDeleteErrorResponse{
				User:    openapi.UserOnlyIdError{Id: userIdMsg},
				Project: openapi.ProjectOnlyIdError{Id: projectIdMsg},
				Chapter: openapi.ChapterOnlyIdError{Id: chapterIdMsg},
				Section: openapi.SectionOnlyIdError{Id: sectionIdMsg},
			},
		)
	}

	sErr := uc.service.DeleteGraph(*userId, *projectId, *chapterId, *sectionId)

	if sErr != nil && sErr.Code() == service.NotFoundError {
		return NewMessageBasedError[openapi.GraphDeleteErrorResponse](
			NotFoundError,
			sErr.Unwrap().Error(),
		)
	}
	if sErr != nil {
		return NewMessageBasedError[openapi.GraphDeleteErrorResponse](
			InternalErrorPanic,
			sErr.Unwrap().Error(),
		)
	}

	return nil
}

func (uc graphUseCase) SectionalizeGraph(req openapi.GraphSectionalizeRequest) (
	*openapi.GraphSectionalizeResponse, *Error[openapi.GraphSectionalizeErrorResponse]) {
	userId, userIdErr := domain.NewUserIdObject(req.User.Id)
	projectId, projectIdErr := domain.NewProjectIdObject(req.Project.Id)
	chapterId, chapterIdErr := domain.NewChapterIdObject(req.Chapter.Id)
	sections, sectionsErr, sectionsOk := uc.sectiionsModelToEntity(req.Sections)

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

	if userIdErr != nil || projectIdErr != nil || chapterIdErr != nil || !sectionsOk {
		return nil, NewModelBasedError(
			DomainValidationError,
			openapi.GraphSectionalizeErrorResponse{
				User:     openapi.UserOnlyIdError{Id: userIdMsg},
				Project:  openapi.ProjectOnlyIdError{Id: projectIdMsg},
				Chapter:  openapi.ChapterOnlyIdError{Id: chapterIdMsg},
				Sections: *sectionsErr,
			},
		)
	}

	entities, sErr := uc.service.SectionalizeIntoGraphs(*userId, *projectId, *chapterId, *sections)
	if sErr != nil && sErr.Code() == service.InvalidArgumentError {
		return nil, NewMessageBasedError[openapi.GraphSectionalizeErrorResponse](
			InvalidArgumentError,
			sErr.Unwrap().Error(),
		)
	}
	if sErr != nil && sErr.Code() == service.NotFoundError {
		return nil, NewMessageBasedError[openapi.GraphSectionalizeErrorResponse](
			NotFoundError,
			sErr.Unwrap().Error(),
		)
	}
	if sErr != nil {
		return nil, NewMessageBasedError[openapi.GraphSectionalizeErrorResponse](
			InternalErrorPanic,
			sErr.Unwrap().Error(),
		)
	}

	graphs := make([]openapi.Graph, len(entities))
	for i, entity := range entities {
		graphs[i] = openapi.Graph{
			Id:        entity.Id().Value(),
			Name:      entity.Name().Value(),
			Paragraph: entity.Paragraph().Value(),
			Children:  uc.childrenEntityToModel(entity.Children()),
		}
	}

	return &openapi.GraphSectionalizeResponse{Graphs: graphs}, nil
}

func (uc graphUseCase) childrenEntityToModel(entity *domain.GraphChildrenEntity) []openapi.GraphChild {
	children := make([]openapi.GraphChild, len(entity.Value()))
	for i, child := range entity.Value() {
		children[i] = openapi.GraphChild{
			Name:        child.Name().Value(),
			Relation:    child.Relation().Value(),
			Description: child.Description().Value(),
			Children:    uc.childrenEntityToModel(child.Children()),
		}
	}
	return children
}

func (uc graphUseCase) childrenModelToEntity(children []openapi.GraphChild) (
	*domain.GraphChildrenEntity, *openapi.GraphChildrenError, bool) {
	childItems := make([]domain.GraphChildEntity, len(children))
	childItemErrors := make([]openapi.GraphChildError, len(children))

	childItemErrorExists := false
	for i, child := range children {
		name, nameErr := domain.NewGraphNameObject(child.Name)
		if nameErr != nil {
			childItemErrors[i].Name = nameErr.Error()
			childItemErrorExists = true
		}
		relation, relationErr := domain.NewGraphRelationObject(child.Relation)
		if relationErr != nil {
			childItemErrors[i].Relation = relationErr.Error()
			childItemErrorExists = true
		}
		desc, descErr := domain.NewGraphDescriptionObject(child.Description)
		if descErr != nil {
			childItemErrors[i].Description = descErr.Error()
			childItemErrorExists = true
		}
		children, childrenErr, childrenOk := uc.childrenModelToEntity(child.Children)
		childItemErrors[i].Children = *childrenErr
		if !childrenOk {
			childItemErrorExists = true
		}
		if nameErr == nil && relationErr == nil && descErr == nil && childrenOk {
			childItems[i] = *domain.NewGraphChildEntity(*name, *relation, *desc, *children)
		}
	}

	childrenErrorMessage := ""
	entity, err := domain.NewGraphChildrenEntity(childItems)
	if err != nil {
		childrenErrorMessage = err.Error()
	}

	ok := childrenErrorMessage == "" && !childItemErrorExists
	return entity, &openapi.GraphChildrenError{Message: childrenErrorMessage, Items: childItemErrors}, ok
}

func (uc graphUseCase) sectiionsModelToEntity(sections []openapi.SectionWithoutAutofield) (
	*domain.SectionWithoutAutofieldEntityList, *openapi.SectionWithoutAutofieldListError, bool) {
	sectionItems := make([]domain.SectionWithoutAutofieldEntity, len(sections))
	sectionItemErrors := make([]openapi.SectionWithoutAutofieldError, len(sections))

	sectionItemErrorExists := false
	for i, section := range sections {
		sectionName, sectionNameErr := domain.NewSectionNameObject(section.Name)
		if sectionNameErr != nil {
			sectionItemErrors[i].Name = sectionNameErr.Error()
			sectionItemErrorExists = true
		}
		sectionContent, sectionContentErr := domain.NewSectionContentObject(section.Content)
		if sectionContentErr != nil {
			sectionItemErrors[i].Content = sectionContentErr.Error()
			sectionItemErrorExists = true
		}
		if sectionNameErr == nil && sectionContentErr == nil {
			sectionItems[i] = *domain.NewSectionWithoutAutofieldEntity(*sectionName, *sectionContent)
		}
	}

	sectionsErrorMessage := ""
	entity, err := domain.NewSectionWithoutAutofieldEntityList(sectionItems)
	if err != nil {
		sectionsErrorMessage = err.Error()
	}

	ok := sectionsErrorMessage == "" && !sectionItemErrorExists
	return entity, &openapi.SectionWithoutAutofieldListError{Message: sectionsErrorMessage, Items: sectionItemErrors}, ok
}
