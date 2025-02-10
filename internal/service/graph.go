package service

import (
	"errors"

	"github.com/kumachan-mis/knodeledge-api/internal/domain"
	"github.com/kumachan-mis/knodeledge-api/internal/record"
	"github.com/kumachan-mis/knodeledge-api/internal/repository"
)

//go:generate mockgen -source=$GOFILE -destination=../../mock/$GOPACKAGE/mock_$GOFILE -package=$GOPACKAGE

type GraphService interface {
	FindGraph(
		userId domain.UserIdObject,
		projectId domain.ProjectIdObject,
		chapterId domain.ChapterIdObject,
		sectionId domain.SectionIdObject,
	) (*domain.GraphEntity, *Error)
	UpdateGraphContent(
		userId domain.UserIdObject,
		projectId domain.ProjectIdObject,
		chapterId domain.ChapterIdObject,
		graphId domain.GraphIdObject,
		graph domain.GraphContentEntity,
	) (*domain.GraphEntity, *Error)
	DeleteGraph(
		userId domain.UserIdObject,
		projectId domain.ProjectIdObject,
		chapterId domain.ChapterIdObject,
		sectionId domain.SectionIdObject,
	) *Error
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

func (s graphService) FindGraph(
	userId domain.UserIdObject,
	projectId domain.ProjectIdObject,
	chapterId domain.ChapterIdObject,
	sectionId domain.SectionIdObject,
) (*domain.GraphEntity, *Error) {
	entry, rErr := s.repository.FetchGraph(userId.Value(), projectId.Value(), chapterId.Value(), sectionId.Value())
	if rErr != nil && rErr.Code() == repository.NotFoundError {
		return nil, Errorf(NotFoundError, "failed to find graph: %w", rErr.Unwrap())
	}
	if rErr != nil {
		return nil, Errorf(RepositoryFailurePanic, "failed to fetch graph: %w", rErr.Unwrap())
	}

	return s.entryToEntity(sectionId.Value(), *entry)
}

func (s graphService) UpdateGraphContent(
	userId domain.UserIdObject,
	projectId domain.ProjectIdObject,
	chapterId domain.ChapterIdObject,
	graphId domain.GraphIdObject,
	graph domain.GraphContentEntity,
) (*domain.GraphEntity, *Error) {
	entryWithoutAutofield := record.GraphContentEntry{
		Paragraph: graph.Paragraph().Value(),
		Children:  s.childrenEntityToEntry(*graph.Children()),
	}

	entry, rErr := s.repository.UpdateGraphContent(
		userId.Value(),
		projectId.Value(),
		chapterId.Value(),
		graphId.Value(),
		entryWithoutAutofield,
	)
	if rErr != nil && rErr.Code() == repository.NotFoundError {
		return nil, Errorf(NotFoundError, "failed to update graph content: %w", rErr.Unwrap())
	}
	if rErr != nil {
		return nil, Errorf(RepositoryFailurePanic, "failed to update graph content: %w", rErr.Unwrap())
	}

	return s.entryToEntity(graphId.Value(), *entry)
}

func (s graphService) DeleteGraph(
	userId domain.UserIdObject,
	projectId domain.ProjectIdObject,
	chapterId domain.ChapterIdObject,
	sectionId domain.SectionIdObject,
) *Error {
	chapter, rErr := s.chapterRepository.FetchChapter(userId.Value(), projectId.Value(), chapterId.Value())
	if rErr != nil && rErr.Code() == repository.NotFoundError {
		return Errorf(NotFoundError, "failed to delete graph: %w", rErr.Unwrap())
	}
	if rErr != nil {
		return Errorf(RepositoryFailurePanic, "failed to delete graph: %w", rErr.Unwrap())
	}

	rErr = s.repository.DeleteGraph(userId.Value(), projectId.Value(), chapterId.Value(), sectionId.Value())
	if rErr != nil && rErr.Code() == repository.NotFoundError {
		return Errorf(NotFoundError, "failed to delete graph: %w", rErr.Unwrap())
	}
	if rErr != nil {
		return Errorf(RepositoryFailurePanic, "failed to delete graph: %w", rErr.Unwrap())
	}

	sectionWithoutAutofield := make([]record.SectionWithoutAutofieldEntry, 0, len(chapter.Sections)-1)
	for _, section := range chapter.Sections {
		if section.Id == sectionId.Value() {
			continue
		}
		sectionWithoutAutofield = append(sectionWithoutAutofield, record.SectionWithoutAutofieldEntry{
			Id:   section.Id,
			Name: section.Name,
		})
	}

	_, rErr = s.chapterRepository.UpdateChapterSections(
		userId.Value(),
		projectId.Value(),
		chapterId.Value(),
		sectionWithoutAutofield,
	)
	if rErr != nil && rErr.Code() == repository.NotFoundError {
		return Errorf(NotFoundError, "failed to delete graph: %w", rErr.Unwrap())
	}
	if rErr != nil {
		return Errorf(RepositoryFailurePanic, "failed to delete graph: %w", rErr.Unwrap())
	}

	return nil
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
			Children:  []record.GraphChildEntry{},
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
	children, sErr := s.childrenEntryToEntity(entry.Children)
	if sErr != nil {
		return nil, Errorf(DomainFailurePanic, "failed to convert entry to entity (children): %w", sErr.Unwrap())
	}
	createdAt, err := domain.NewCreatedAtObject(entry.CreatedAt)
	if err != nil {
		return nil, Errorf(DomainFailurePanic, "failed to convert entry to entity (createdAt): %w", err)
	}
	updatedAt, err := domain.NewUpdatedAtObject(entry.UpdatedAt)
	if err != nil {
		return nil, Errorf(DomainFailurePanic, "failed to convert entry to entity (updatedAt): %w", err)
	}

	return domain.NewGraphEntity(*id, *name, *paragraph, *children, *createdAt, *updatedAt), nil
}

func (s graphService) childrenEntryToEntity(entry []record.GraphChildEntry) (*domain.GraphChildrenEntity, *Error) {
	entities := make([]domain.GraphChildEntity, len(entry))
	for i, child := range entry {
		name, err := domain.NewGraphNameObject(child.Name)
		if err != nil {
			return nil, Errorf(DomainFailurePanic, "failed to convert child entry to entity (name): %w", err)
		}
		relation, err := domain.NewGraphRelationObject(child.Relation)
		if err != nil {
			return nil, Errorf(DomainFailurePanic, "failed to convert child entry to entity (relation): %w", err)
		}
		description, err := domain.NewGraphDescriptionObject(child.Description)
		if err != nil {
			return nil, Errorf(DomainFailurePanic, "failed to convert child entry to entity (description): %w", err)
		}
		children, err := s.childrenEntryToEntity(child.Children)
		entities[i] = *domain.NewGraphChildEntity(*name, *relation, *description, *children)
	}

	entity, err := domain.NewGraphChildrenEntity(entities)
	if err != nil {
		return nil, Errorf(DomainFailurePanic, "%w", err)
	}
	return entity, nil
}

func (s graphService) childrenEntityToEntry(entity domain.GraphChildrenEntity) []record.GraphChildEntry {
	entries := make([]record.GraphChildEntry, entity.Len())
	for i, child := range entity.Value() {
		entries[i] = record.GraphChildEntry{
			Name:        child.Name().Value(),
			Relation:    child.Relation().Value(),
			Description: child.Description().Value(),
			Children:    s.childrenEntityToEntry(*child.Children()),
		}
	}
	return entries
}
