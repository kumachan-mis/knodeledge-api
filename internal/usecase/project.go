package usecase

import (
	"github.com/kumachan-mis/knodeledge-api/internal/domain"
	"github.com/kumachan-mis/knodeledge-api/internal/model"
	"github.com/kumachan-mis/knodeledge-api/internal/service"
)

//go:generate mockgen -source=$GOFILE -destination=../../mock/$GOPACKAGE/mock_$GOFILE -package=$GOPACKAGE

type ProjectUseCase interface {
	ListProjects(user model.User) ([]model.Project, *Error[model.UserError])
}

type projectUseCase struct {
	service service.ProjectService
}

func NewProjectUseCase(service service.ProjectService) ProjectUseCase {
	return projectUseCase{service: service}
}

func (uc projectUseCase) ListProjects(user model.User) ([]model.Project, *Error[model.UserError]) {
	uid, err := domain.NewUserIdObject(user.Id)
	if err != nil {
		return nil, NewModelBasedError(InvalidArgumentError, model.UserError{Id: err.Error()})
	}

	entities, err := uc.service.ListProjects(*uid)
	if err != nil {
		return nil, NewMessageBasedError[model.UserError](InternalError, err.Error())
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
	return projects, nil
}
