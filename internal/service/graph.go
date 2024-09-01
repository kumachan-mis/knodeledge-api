package service

import (
	"errors"

	"github.com/kumachan-mis/knodeledge-api/internal/domain"
	"github.com/kumachan-mis/knodeledge-api/internal/record"
	"github.com/kumachan-mis/knodeledge-api/internal/repository"
)

//go:generate mockgen -source=$GOFILE -destination=../../mock/$GOPACKAGE/mock_$GOFILE -package=$GOPACKAGE

type GraphService interface {
	SectionalizeIntoGraphs(
		userId domain.UserIdObject,
		projectId domain.ProjectIdObject,
		chapterId domain.ChapterIdObject,
		sections domain.SectionWithoutAutofieldEntityList,
	) ([]domain.GraphEntity, *Error)
}

type graphService struct {
	repository        repository.GraphRepository
	chapterRepository repository.ChapterRepository
}

func NewGraphService(
	repository repository.GraphRepository,
	chapterRepository repository.ChapterRepository,
) GraphService {
	return graphService{repository: repository, chapterRepository: chapterRepository}
}

func (s graphService) SectionalizeIntoGraphs(
	userId domain.UserIdObject,
	projectId domain.ProjectIdObject,
	chapterId domain.ChapterIdObject,
	sections domain.SectionWithoutAutofieldEntityList,
) ([]domain.GraphEntity, *Error) {
	exists, rErr := s.repository.GraphExists(userId.Value(), projectId.Value(), chapterId.Value())
	if rErr != nil && rErr.Code() == repository.NotFoundError {
		return nil, Errorf(NotFoundError, "failed to sectionalize into graphs: %w", rErr.Unwrap())
	}
	if rErr != nil {
		return nil, Errorf(RepositoryFailurePanic, "failed to check graph existence: %w", rErr.Unwrap())
	}
	if exists {
		err := errors.New("graph already exists")
		return nil, Errorf(InvalidArgument, "failed to sectionalize into graphs: %w", err)
	}

	entriesWithoutAutofield := make([]record.GraphWithoutAutofieldEntry, sections.Len())
	for i, section := range sections.Value() {
		entriesWithoutAutofield[i] = record.GraphWithoutAutofieldEntry{
			Name:      section.Name().Value(),
			Paragraph: section.Content().Value(),
		}
	}

	keys, entries, rErr := s.repository.InsertGraphs(
		userId.Value(),
		projectId.Value(),
		chapterId.Value(),
		entriesWithoutAutofield,
	)
	if rErr != nil && rErr.Code() == repository.NotFoundError {
		return nil, Errorf(NotFoundError, "failed to sectionalize into graphs: %w", rErr.Unwrap())
	}
	if rErr != nil {
		return nil, Errorf(RepositoryFailurePanic, "failed to insert graphs: %w", rErr.Unwrap())
	}

	entities := make([]domain.GraphEntity, sections.Len())
	for i, entry := range entries {
		entity, sErr := s.entryToEntity(keys[i], entry)
		if sErr != nil {
			return nil, sErr
		}
		entities[i] = *entity
	}

	sectionWithoutAutofield := make([]record.SectionWithoutAutofieldEntry, sections.Len())
	for i, section := range sections.Value() {
		sectionWithoutAutofield[i] = record.SectionWithoutAutofieldEntry{
			Id:   keys[i],
			Name: section.Name().Value(),
		}
	}

	_, rErr = s.chapterRepository.UpdateChapterSections(
		userId.Value(),
		projectId.Value(),
		chapterId.Value(),
		sectionWithoutAutofield,
	)
	if rErr != nil && rErr.Code() == repository.NotFoundError {
		return nil, Errorf(NotFoundError, "failed to sectionalize into graphs: %w", rErr.Unwrap())
	}
	if rErr != nil {
		return nil, Errorf(RepositoryFailurePanic, "failed to update sections of chapter: %w", rErr.Unwrap())
	}

	return entities, nil
}

func (s graphService) entryToEntity(key string, entry record.GraphEntry) (*domain.GraphEntity, *Error) {
	id, err := domain.NewGraphIdObject(key)
	if err != nil {
		return nil, Errorf(DomainFailurePanic, "failed to convert entry to entity (id): %w", err)
	}
	name, err := domain.NewGraphNameObject(entry.Name)
	if err != nil {
		return nil, Errorf(DomainFailurePanic, "failed to convert entry to entity (name): %w", err)
	}
	paragraph, err := domain.NewGraphParagraphObject(entry.Paragraph)
	if err != nil {
		return nil, Errorf(DomainFailurePanic, "failed to convert entry to entity (paragraph): %w", err)
	}
	createdAt, err := domain.NewCreatedAtObject(entry.CreatedAt)
	if err != nil {
		return nil, Errorf(DomainFailurePanic, "failed to convert entry to entity (createdAt): %w", err)
	}
	updatedAt, err := domain.NewUpdatedAtObject(entry.UpdatedAt)
	if err != nil {
		return nil, Errorf(DomainFailurePanic, "failed to convert entry to entity (updatedAt): %w", err)
	}

	return domain.NewGraphEntity(*id, *name, *paragraph, *createdAt, *updatedAt), nil
}
