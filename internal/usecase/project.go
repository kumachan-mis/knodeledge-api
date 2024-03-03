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
	CreateProject(req model.ProjectCreateRequest) (
		*model.ProjectCreateResponse, *Error[model.ProjectCreateErrorResponse])
}

type projectUseCase struct {
	service service.ProjectService
}

func NewProjectUseCase(service service.ProjectService) ProjectUseCase {
	return projectUseCase{service: service}
}

func (uc projectUseCase) ListProjects(req model.ProjectListRequest) (
	*model.ProjectListResponse, *Error[model.ProjectListErrorResponse]) {
	uid, uidErr := domain.NewUserIdObject(req.User.Id)
	if uidErr != nil {
		return nil, NewModelBasedError(
			InvalidArgumentError,
			model.ProjectListErrorResponse{User: model.UserOnlyIdError{Id: uidErr.Error()}},
		)
	}

	entities, sErr := uc.service.ListProjects(*uid)
	if sErr != nil {
		return nil, NewMessageBasedError[model.ProjectListErrorResponse](
			InternalErrorPanic,
			sErr.Unwrap().Error(),
		)
	}

	projects := make([]model.Project, len(entities))
	i := 0
	for _, entity := range entities {
		projects[i] = model.Project{
			Id:          entity.Id().Value(),
			Name:        entity.Name().Value(),
			Description: entity.Description().Value(),
		}
		i++
	}
	return &model.ProjectListResponse{Projects: projects}, nil
}

func (uc projectUseCase) CreateProject(req model.ProjectCreateRequest) (
	*model.ProjectCreateResponse, *Error[model.ProjectCreateErrorResponse]) {
	uid, uidErr := domain.NewUserIdObject(req.User.Id)
	pname, pnameErr := domain.NewProjectNameObject(req.Project.Name)
	pdesc, pdescErr := domain.NewProjectDescriptionObject(req.Project.Description)

	uidMsg := ""
	if uidErr != nil {
		uidMsg = uidErr.Error()
	}
	pnameMsg := ""
	if pnameErr != nil {
		pnameMsg = pnameErr.Error()
	}
	pdescMsg := ""
	if pdescErr != nil {
		pdescMsg = pdescErr.Error()
	}

	if uidErr != nil || pnameErr != nil || pdescErr != nil {
		return nil, NewModelBasedError(
			InvalidArgumentError,
			model.ProjectCreateErrorResponse{
				User: model.UserOnlyIdError{
					Id: uidMsg,
				},
				Project: model.ProjectWithoutAutofieldError{
					Name:        pnameMsg,
					Description: pdescMsg,
				},
			},
		)
	}

	project := domain.NewProjectWithoutAutofieldEntity(*pname, *pdesc)

	entity, sErr := uc.service.CreateProject(*uid, *project)
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
