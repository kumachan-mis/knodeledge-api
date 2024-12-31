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

	entity, sErr := uc.service.FindGraph(*userId, *projectId, *chapterId, *sectionId)

	if sErr != nil && sErr.Code() == service.NotFoundError {
		return nil, NewMessageBasedError[model.GraphFindErrorResponse](
			NotFoundError,
			sErr.Unwrap().Error(),
		)
	}
	if sErr != nil {
		return nil, NewMessageBasedError[model.GraphFindErrorResponse](
			InternalErrorPanic,
			sErr.Unwrap().Error(),
		)
	}

	return &model.GraphFindResponse{
		Graph: model.Graph{
			Id:        entity.Id().Value(),
			Name:      entity.Name().Value(),
			Paragraph: entity.Paragraph().Value(),
			Children:  uc.childrenEntityToModel(entity.Children()),
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
			model.GraphUpdateErrorResponse{
				User:    model.UserOnlyIdError{Id: userIdMsg},
				Project: model.ProjectOnlyIdError{Id: projectIdMsg},
				Chapter: model.ChapterOnlyIdError{Id: chapterIdMsg},
				Graph: model.GraphContentError{
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
			Children:  uc.childrenEntityToModel(entity.Children()),
		},
	}, nil
}

func (uc graphUseCase) SectionalizeGraph(req model.GraphSectionalizeRequest) (
	*model.GraphSectionalizeResponse, *Error[model.GraphSectionalizeErrorResponse]) {
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
			model.GraphSectionalizeErrorResponse{
				User:     model.UserOnlyIdError{Id: userIdMsg},
				Project:  model.ProjectOnlyIdError{Id: projectIdMsg},
				Chapter:  model.ChapterOnlyIdError{Id: chapterIdMsg},
				Sections: *sectionsErr,
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
			Children:  uc.childrenEntityToModel(entity.Children()),
		}
	}

	return &model.GraphSectionalizeResponse{Graphs: graphs}, nil
}

func (uc graphUseCase) childrenEntityToModel(entity *domain.GraphChildrenEntity) []model.GraphChild {
	children := make([]model.GraphChild, len(entity.Value()))
	for i, child := range entity.Value() {
		children[i] = model.GraphChild{
			Name:        child.Name().Value(),
			Relation:    child.Relation().Value(),
			Description: child.Description().Value(),
			Children:    uc.childrenEntityToModel(child.Children()),
		}
	}
	return children
}

func (uc graphUseCase) childrenModelToEntity(children []model.GraphChild) (
	*domain.GraphChildrenEntity, *model.GraphChildrenError, bool) {
	childItems := make([]domain.GraphChildEntity, len(children))
	childItemErrors := make([]model.GraphChildError, len(children))

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
	return entity, &model.GraphChildrenError{Message: childrenErrorMessage, Items: childItemErrors}, ok
}

func (uc graphUseCase) sectiionsModelToEntity(sections []model.SectionWithoutAutofield) (
	*domain.SectionWithoutAutofieldEntityList, *model.SectionWithoutAutofieldListError, bool) {
	sectionItems := make([]domain.SectionWithoutAutofieldEntity, len(sections))
	sectionItemErrors := make([]model.SectionWithoutAutofieldError, len(sections))

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
	return entity, &model.SectionWithoutAutofieldListError{Message: sectionsErrorMessage, Items: sectionItemErrors}, ok
}
