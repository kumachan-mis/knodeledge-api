package usecase

import (
	"github.com/kumachan-mis/knodeledge-api/internal/domain"
	"github.com/kumachan-mis/knodeledge-api/internal/model"
	"github.com/kumachan-mis/knodeledge-api/internal/service"
)

//go:generate mockgen -source=$GOFILE -destination=../../mock/$GOPACKAGE/mock_$GOFILE -package=$GOPACKAGE

type GraphUseCase interface {
	SectionalizeGraph(request model.GraphSectionalizeRequest) (
		*model.GraphSectionalizeResponse, *Error[model.GraphSectionalizeErrorResponse])
}

type graphUseCase struct {
	service service.GraphService
}

func NewGraphUseCase(service service.GraphService) GraphUseCase {
	return graphUseCase{service: service}
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
