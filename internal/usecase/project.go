package usecase

import (
	"github.com/kumachan-mis/knodeledge-api/internal/domain"
	"github.com/kumachan-mis/knodeledge-api/internal/model"
	"github.com/kumachan-mis/knodeledge-api/internal/service"
)

//go:generate mockgen -source=$GOFILE -destination=../../mock/$GOPACKAGE/mock_$GOFILE -package=$GOPACKAGE

type ProjectUseCase interface {
	ListProjects(req model.ProjectListRequest) (
		*model.ProjectListResponse, *Error[model.ProjectListErrorResponse])
	FindProject(req model.ProjectFindRequest) (
		*model.ProjectFindResponse, *Error[model.ProjectFindErrorResponse])
	CreateProject(req model.ProjectCreateRequest) (
		*model.ProjectCreateResponse, *Error[model.ProjectCreateErrorResponse])
	UpdateProject(req model.ProjectUpdateRequest) (
		*model.ProjectUpdateResponse, *Error[model.ProjectUpdateErrorResponse])
}

type projectUseCase struct {
	service service.ProjectService
}

func NewProjectUseCase(service service.ProjectService) ProjectUseCase {
	return projectUseCase{service: service}
}

func (uc projectUseCase) ListProjects(req model.ProjectListRequest) (
	*model.ProjectListResponse, *Error[model.ProjectListErrorResponse]) {
	userId, userIdErr := domain.NewUserIdObject(req.User.Id)
	if userIdErr != nil {
		return nil, NewModelBasedError(
			DomainValidationError,
			model.ProjectListErrorResponse{User: model.UserOnlyIdError{Id: userIdErr.Error()}},
		)
	}

	entities, sErr := uc.service.ListProjects(*userId)
	if sErr != nil {
		return nil, NewMessageBasedError[model.ProjectListErrorResponse](
			InternalErrorPanic,
			sErr.Unwrap().Error(),
		)
	}

	projects := make([]model.Project, len(entities))
	for i, entity := range entities {
		projects[i] = model.Project{
			Id:          entity.Id().Value(),
			Name:        entity.Name().Value(),
			Description: entity.Description().Value(),
		}
	}
	return &model.ProjectListResponse{Projects: projects}, nil
}

func (uc projectUseCase) FindProject(req model.ProjectFindRequest) (
	*model.ProjectFindResponse, *Error[model.ProjectFindErrorResponse]) {
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
			model.ProjectFindErrorResponse{
				User:    model.UserOnlyIdError{Id: userIdMsg},
				Project: model.ProjectOnlyIdError{Id: projectIdMsg},
			},
		)
	}

	entity, sErr := uc.service.FindProject(*userId, *projectId)
	if sErr != nil && sErr.Code() == service.NotFoundError {
		return nil, NewMessageBasedError[model.ProjectFindErrorResponse](
			NotFoundError,
			sErr.Unwrap().Error(),
		)
	}
	if sErr != nil {
		return nil, NewMessageBasedError[model.ProjectFindErrorResponse](
			InternalErrorPanic,
			sErr.Unwrap().Error(),
		)
	}

	return &model.ProjectFindResponse{
		Project: model.Project{
			Id:          entity.Id().Value(),
			Name:        entity.Name().Value(),
			Description: entity.Description().Value(),
		},
	}, nil
}

func (uc projectUseCase) CreateProject(req model.ProjectCreateRequest) (
	*model.ProjectCreateResponse, *Error[model.ProjectCreateErrorResponse]) {
	userId, userIdErr := domain.NewUserIdObject(req.User.Id)
	projectName, projectNameErr := domain.NewProjectNameObject(req.Project.Name)
	projectDesc, projectDescErr := domain.NewProjectDescriptionObject(req.Project.Description)

	userIdMsg := ""
	if userIdErr != nil {
		userIdMsg = userIdErr.Error()
	}
	projectNameMsg := ""
	if projectNameErr != nil {
		projectNameMsg = projectNameErr.Error()
	}
	projectDescMsg := ""
	if projectDescErr != nil {
		projectDescMsg = projectDescErr.Error()
	}

	if userIdErr != nil || projectNameErr != nil || projectDescErr != nil {
		return nil, NewModelBasedError(
			DomainValidationError,
			model.ProjectCreateErrorResponse{
				User: model.UserOnlyIdError{
					Id: userIdMsg,
				},
				Project: model.ProjectWithoutAutofieldError{
					Name:        projectNameMsg,
					Description: projectDescMsg,
				},
			},
		)
	}

	project := domain.NewProjectWithoutAutofieldEntity(*projectName, *projectDesc)

	entity, sErr := uc.service.CreateProject(*userId, *project)
	if sErr != nil {
		return nil, NewMessageBasedError[model.ProjectCreateErrorResponse](
			InternalErrorPanic,
			sErr.Unwrap().Error(),
		)
	}

	return &model.ProjectCreateResponse{
		Project: model.Project{
			Id:          entity.Id().Value(),
			Name:        entity.Name().Value(),
			Description: entity.Description().Value(),
		},
	}, nil
}

func (uc projectUseCase) UpdateProject(req model.ProjectUpdateRequest) (
	*model.ProjectUpdateResponse, *Error[model.ProjectUpdateErrorResponse]) {
	userId, userIdErr := domain.NewUserIdObject(req.User.Id)
	projectId, projectIdErr := domain.NewProjectIdObject(req.Project.Id)
	projectName, projectNameErr := domain.NewProjectNameObject(req.Project.Name)
	projectDesc, projectDescErr := domain.NewProjectDescriptionObject(req.Project.Description)

	userIdMsg := ""
	if userIdErr != nil {
		userIdMsg = userIdErr.Error()
	}
	projectIdMsg := ""
	if projectIdErr != nil {
		projectIdMsg = projectIdErr.Error()
	}
	projectNameMsg := ""
	if projectNameErr != nil {
		projectNameMsg = projectNameErr.Error()
	}
	projectDescMsg := ""
	if projectDescErr != nil {
		projectDescMsg = projectDescErr.Error()
	}

	if userIdErr != nil || projectIdErr != nil || projectNameErr != nil || projectDescErr != nil {
		return nil, NewModelBasedError(
			DomainValidationError,
			model.ProjectUpdateErrorResponse{
				User: model.UserOnlyIdError{
					Id: userIdMsg,
				},
				Project: model.ProjectError{
					Id:          projectIdMsg,
					Name:        projectNameMsg,
					Description: projectDescMsg,
				},
			},
		)
	}

	project := domain.NewProjectWithoutAutofieldEntity(*projectName, *projectDesc)

	entity, sErr := uc.service.UpdateProject(*userId, *projectId, *project)
	if sErr != nil && sErr.Code() == service.NotFoundError {
		return nil, NewMessageBasedError[model.ProjectUpdateErrorResponse](
			NotFoundError,
			sErr.Unwrap().Error(),
		)
	}
	if sErr != nil {
		return nil, NewMessageBasedError[model.ProjectUpdateErrorResponse](
			InternalErrorPanic,
			sErr.Unwrap().Error(),
		)
	}

	return &model.ProjectUpdateResponse{
		Project: model.Project{
			Id:          entity.Id().Value(),
			Name:        entity.Name().Value(),
			Description: entity.Description().Value(),
		},
	}, nil
}
