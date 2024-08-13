package service

import (
	"sort"

	"github.com/kumachan-mis/knodeledge-api/internal/domain"
	"github.com/kumachan-mis/knodeledge-api/internal/record"
	"github.com/kumachan-mis/knodeledge-api/internal/repository"
)

//go:generate mockgen -source=$GOFILE -destination=../../mock/$GOPACKAGE/mock_$GOFILE -package=$GOPACKAGE

type ProjectService interface {
	ListProjects(
		userId domain.UserIdObject,
	) ([]domain.ProjectEntity, *Error)
	FindProject(
		userId domain.UserIdObject,
		projectId domain.ProjectIdObject,
	) (*domain.ProjectEntity, *Error)
	CreateProject(
		userId domain.UserIdObject,
		project domain.ProjectWithoutAutofieldEntity,
	) (*domain.ProjectEntity, *Error)
	UpdateProject(
		userId domain.UserIdObject,
		projectId domain.ProjectIdObject,
		project domain.ProjectWithoutAutofieldEntity,
	) (*domain.ProjectEntity, *Error)
}

type projectService struct {
	repository repository.ProjectRepository
}

func NewProjectService(repository repository.ProjectRepository) ProjectService {
	return projectService{repository: repository}
}

func (s projectService) ListProjects(
	userId domain.UserIdObject,
) ([]domain.ProjectEntity, *Error) {
	entries, rErr := s.repository.FetchProjects(userId.Value())
	if rErr != nil {
		return nil, Errorf(RepositoryFailurePanic, "failed to fetch user projects: %w", rErr.Unwrap())
	}

	projects := []domain.ProjectEntity{}
	for key, entry := range entries {
		project, err := s.entryToEntity(key, entry)
		if err != nil {
			return nil, err
		}

		projects = append(projects, *project)
	}

	sort.Slice(projects, func(i, j int) bool {
		ikey := projects[i].UpdatedAt().Value()
		jkey := projects[j].UpdatedAt().Value()
		return ikey.After(jkey)
	})

	return projects, nil
}

func (s projectService) FindProject(
	userId domain.UserIdObject,
	projectId domain.ProjectIdObject,
) (*domain.ProjectEntity, *Error) {
	entry, rErr := s.repository.FetchProject(userId.Value(), projectId.Value())
	if rErr != nil && rErr.Code() == repository.NotFoundError {
		return nil, Errorf(NotFoundError, "failed to find project: %w", rErr.Unwrap())
	}
	if rErr != nil {
		return nil, Errorf(RepositoryFailurePanic, "failed to fetch project: %w", rErr.Unwrap())
	}

	entity, err := s.entryToEntity(projectId.Value(), *entry)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (s projectService) CreateProject(
	userId domain.UserIdObject,
	project domain.ProjectWithoutAutofieldEntity,
) (*domain.ProjectEntity, *Error) {
	entryWithoutAutofield := record.ProjectWithoutAutofieldEntry{
		Name:        project.Name().Value(),
		Description: project.Description().Value(),
	}

	key, entry, rErr := s.repository.InsertProject(userId.Value(), entryWithoutAutofield)
	if rErr != nil {
		return nil, Errorf(RepositoryFailurePanic, "failed to insert project: %w", rErr.Unwrap())
	}

	return s.entryToEntity(key, *entry)
}

func (s projectService) UpdateProject(
	userId domain.UserIdObject,
	projectId domain.ProjectIdObject,
	project domain.ProjectWithoutAutofieldEntity,
) (*domain.ProjectEntity, *Error) {
	entryWithoutAutofield := record.ProjectWithoutAutofieldEntry{
		Name:        project.Name().Value(),
		Description: project.Description().Value(),
	}

	entry, rErr := s.repository.UpdateProject(userId.Value(), projectId.Value(), entryWithoutAutofield)
	if rErr != nil && rErr.Code() == repository.NotFoundError {
		return nil, Errorf(NotFoundError, "failed to update project: %w", rErr.Unwrap())
	}
	if rErr != nil {
		return nil, Errorf(RepositoryFailurePanic, "failed to update project: %w", rErr.Unwrap())
	}

	return s.entryToEntity(projectId.Value(), *entry)
}

func (s projectService) entryToEntity(key string, entry record.ProjectEntry) (*domain.ProjectEntity, *Error) {
	id, err := domain.NewProjectIdObject(key)
	if err != nil {
		return nil, Errorf(DomainFailurePanic, "failed to convert entry to entity (id): %w", err)
	}
	name, err := domain.NewProjectNameObject(entry.Name)
	if err != nil {
		return nil, Errorf(DomainFailurePanic, "failed to convert entry to entity (name): %w", err)
	}
	description, err := domain.NewProjectDescriptionObject(entry.Description)
	if err != nil {
		return nil, Errorf(DomainFailurePanic, "failed to convert entry to entity (description): %w", err)
	}
	createdAt, err := domain.NewCreatedAtObject(entry.CreatedAt)
	if err != nil {
		return nil, Errorf(DomainFailurePanic, "failed to convert entry to entity (createdAt): %w", err)
	}
	updatedAt, err := domain.NewUpdatedAtObject(entry.UpdatedAt)
	if err != nil {
		return nil, Errorf(DomainFailurePanic, "failed to convert entry to entity (updatedAt): %w", err)
	}

	return domain.NewProjectEntity(*id, *name, *description, *createdAt, *updatedAt), nil
}
