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
}

type projectUseCase struct {
	service service.ProjectService
}

func NewProjectUseCase(service service.ProjectService) ProjectUseCase {
	return projectUseCase{service: service}
}

func (uc projectUseCase) ListProjects(req model.ProjectListRequest) (
	*model.ProjectListResponse, *Error[model.ProjectListErrorResponse]) {
	uid, err := domain.NewUserIdObject(req.User.Id)
	if err != nil {
		return nil, NewModelBasedError(
			InvalidArgumentError,
			model.ProjectListErrorResponse{User: model.UserError{Id: err.Error()}},
		)
	}

	entities, err := uc.service.ListProjects(*uid)
	if err != nil {
		return nil, NewMessageBasedError[model.ProjectListErrorResponse](
			InternalError,
			err.Error(),
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
