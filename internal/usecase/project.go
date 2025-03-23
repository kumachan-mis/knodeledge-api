package usecase

import (
	"github.com/kumachan-mis/knodeledge-api/internal/domain"
	"github.com/kumachan-mis/knodeledge-api/internal/openapi"
	"github.com/kumachan-mis/knodeledge-api/internal/service"
)

//go:generate mockgen -source=$GOFILE -destination=../../mock/$GOPACKAGE/mock_$GOFILE -package=$GOPACKAGE

type ProjectUseCase interface {
	ListProjects(req openapi.ProjectListRequest) (
		*openapi.ProjectListResponse, *Error[openapi.ProjectListErrorResponse])
	FindProject(req openapi.ProjectFindRequest) (
		*openapi.ProjectFindResponse, *Error[openapi.ProjectFindErrorResponse])
	CreateProject(req openapi.ProjectCreateRequest) (
		*openapi.ProjectCreateResponse, *Error[openapi.ProjectCreateErrorResponse])
	UpdateProject(req openapi.ProjectUpdateRequest) (
		*openapi.ProjectUpdateResponse, *Error[openapi.ProjectUpdateErrorResponse])
	DeleteProject(req openapi.ProjectDeleteRequest) *Error[openapi.ProjectDeleteErrorResponse]
}

type projectUseCase struct {
	service service.ProjectService
}

func NewProjectUseCase(service service.ProjectService) ProjectUseCase {
	return projectUseCase{service: service}
}

func (uc projectUseCase) ListProjects(req openapi.ProjectListRequest) (
	*openapi.ProjectListResponse, *Error[openapi.ProjectListErrorResponse]) {
	userId, userIdErr := domain.NewUserIdObject(req.UserId)
	if userIdErr != nil {
		return nil, NewModelBasedError(
			DomainValidationError,
			openapi.ProjectListErrorResponse{UserId: userIdErr.Error()},
		)
	}

	entities, sErr := uc.service.ListProjects(*userId)
	if sErr != nil {
		return nil, NewMessageBasedError[openapi.ProjectListErrorResponse](
			InternalErrorPanic,
			sErr.Unwrap().Error(),
		)
	}

	projects := make([]openapi.Project, len(entities))
	for i, entity := range entities {
		projects[i] = openapi.Project{
			Id:          entity.Id().Value(),
			Name:        entity.Name().Value(),
			Description: entity.Description().Value(),
		}
	}
	return &openapi.ProjectListResponse{Projects: projects}, nil
}

func (uc projectUseCase) FindProject(req openapi.ProjectFindRequest) (
	*openapi.ProjectFindResponse, *Error[openapi.ProjectFindErrorResponse]) {
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
			openapi.ProjectFindErrorResponse{
				UserId:    userIdMsg,
				ProjectId: projectIdMsg,
			},
		)
	}

	entity, sErr := uc.service.FindProject(*userId, *projectId)
	if sErr != nil && sErr.Code() == service.NotFoundError {
		return nil, NewMessageBasedError[openapi.ProjectFindErrorResponse](
			NotFoundError,
			sErr.Unwrap().Error(),
		)
	}
	if sErr != nil {
		return nil, NewMessageBasedError[openapi.ProjectFindErrorResponse](
			InternalErrorPanic,
			sErr.Unwrap().Error(),
		)
	}

	return &openapi.ProjectFindResponse{
		Project: openapi.Project{
			Id:          entity.Id().Value(),
			Name:        entity.Name().Value(),
			Description: entity.Description().Value(),
		},
	}, nil
}

func (uc projectUseCase) CreateProject(req openapi.ProjectCreateRequest) (
	*openapi.ProjectCreateResponse, *Error[openapi.ProjectCreateErrorResponse]) {
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
			openapi.ProjectCreateErrorResponse{
				User: openapi.UserOnlyIdError{
					Id: userIdMsg,
				},
				Project: openapi.ProjectWithoutAutofieldError{
					Name:        projectNameMsg,
					Description: projectDescMsg,
				},
			},
		)
	}

	project := domain.NewProjectWithoutAutofieldEntity(*projectName, *projectDesc)

	entity, sErr := uc.service.CreateProject(*userId, *project)
	if sErr != nil {
		return nil, NewMessageBasedError[openapi.ProjectCreateErrorResponse](
			InternalErrorPanic,
			sErr.Unwrap().Error(),
		)
	}

	return &openapi.ProjectCreateResponse{
		Project: openapi.Project{
			Id:          entity.Id().Value(),
			Name:        entity.Name().Value(),
			Description: entity.Description().Value(),
		},
	}, nil
}

func (uc projectUseCase) UpdateProject(req openapi.ProjectUpdateRequest) (
	*openapi.ProjectUpdateResponse, *Error[openapi.ProjectUpdateErrorResponse]) {
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
			openapi.ProjectUpdateErrorResponse{
				User: openapi.UserOnlyIdError{
					Id: userIdMsg,
				},
				Project: openapi.ProjectError{
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
		return nil, NewMessageBasedError[openapi.ProjectUpdateErrorResponse](
			NotFoundError,
			sErr.Unwrap().Error(),
		)
	}
	if sErr != nil {
		return nil, NewMessageBasedError[openapi.ProjectUpdateErrorResponse](
			InternalErrorPanic,
			sErr.Unwrap().Error(),
		)
	}

	return &openapi.ProjectUpdateResponse{
		Project: openapi.Project{
			Id:          entity.Id().Value(),
			Name:        entity.Name().Value(),
			Description: entity.Description().Value(),
		},
	}, nil
}

func (uc projectUseCase) DeleteProject(req openapi.ProjectDeleteRequest) *Error[openapi.ProjectDeleteErrorResponse] {
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
		return NewModelBasedError(
			DomainValidationError,
			openapi.ProjectDeleteErrorResponse{
				User:    openapi.UserOnlyIdError{Id: userIdMsg},
				Project: openapi.ProjectOnlyIdError{Id: projectIdMsg},
			},
		)
	}

	sErr := uc.service.DeleteProject(*userId, *projectId)
	if sErr != nil && sErr.Code() == service.NotFoundError {
		return NewMessageBasedError[openapi.ProjectDeleteErrorResponse](
			NotFoundError,
			sErr.Unwrap().Error(),
		)
	}
	if sErr != nil {
		return NewMessageBasedError[openapi.ProjectDeleteErrorResponse](
			InternalErrorPanic,
			sErr.Unwrap().Error(),
		)
	}

	return nil
}
