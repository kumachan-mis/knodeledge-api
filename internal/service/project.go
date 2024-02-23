package service

import (
	"fmt"
	"sort"

	"github.com/kumachan-mis/knodeledge-api/internal/domain"
	"github.com/kumachan-mis/knodeledge-api/internal/repository"
)

//go:generate mockgen -source=$GOFILE -destination=../../mock/$GOPACKAGE/mock_$GOFILE -package=$GOPACKAGE

type ProjectService interface {
	ListProjects(userId domain.UserIdObject) ([]domain.ProjectEntity, error)
}

type projectService struct {
	repository repository.ProjectRepository
}

func NewProjectService(repository repository.ProjectRepository) ProjectService {
	return projectService{repository: repository}
}

func (s projectService) ListProjects(userId domain.UserIdObject) ([]domain.ProjectEntity, error) {
	entries, err := s.repository.FetchUserProjects(userId.Value())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user projects: %w", err)
	}

	projects := make([]domain.ProjectEntity, len(entries))
	i := 0
	for key, entry := range entries {
		id, err := domain.NewProjectIdObject(key)
		if err != nil {
			return nil, fmt.Errorf("failed to convert entry to entity (id): %w", err)
		}
		name, err := domain.NewProjectNameObject(entry.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to convert entry to entity (name): %w", err)
		}
		description, err := domain.NewProjectDescriptionObject(entry.Description)
		if err != nil {
			return nil, fmt.Errorf("failed to convert entry to entity (description): %w", err)
		}
		createdAt, err := domain.NewCreatedAtObject(entry.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to convert entry to entity (createdAt): %w", err)
		}
		updatedAt, err := domain.NewUpdatedAtObject(entry.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to convert entry to entity (updatedAt): %w", err)
		}
		project, err := domain.NewProjectEntity(*id, *name, *description, *createdAt, *updatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to convert entry to entity: %w", err)
		}
		projects[i] = *project
		i++
	}

	sort.Slice(projects, func(i, j int) bool {
		ikey := projects[i].UpdatedAt().Value()
		jkey := projects[j].UpdatedAt().Value()
		return ikey.Before(jkey)
	})

	return projects, nil
}
