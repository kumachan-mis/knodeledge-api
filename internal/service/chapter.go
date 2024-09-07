package service

import (
	"sort"

	"github.com/kumachan-mis/knodeledge-api/internal/domain"
	"github.com/kumachan-mis/knodeledge-api/internal/record"
	"github.com/kumachan-mis/knodeledge-api/internal/repository"
)

//go:generate mockgen -source=$GOFILE -destination=../../mock/$GOPACKAGE/mock_$GOFILE -package=$GOPACKAGE

type ChapterService interface {
	ListChapters(
		userId domain.UserIdObject,
		projectId domain.ProjectIdObject,
	) ([]domain.ChapterEntity, *Error)
	CreateChapter(
		userId domain.UserIdObject,
		projectId domain.ProjectIdObject,
		chapter domain.ChapterWithoutAutofieldEntity,
	) (*domain.ChapterEntity, *Error)
	UpdateChapter(
		userId domain.UserIdObject,
		projectId domain.ProjectIdObject,
		chapterId domain.ChapterIdObject,
		chapter domain.ChapterWithoutAutofieldEntity,
	) (*domain.ChapterEntity, *Error)
}

type chapterService struct {
	repository      repository.ChapterRepository
	paperRepository repository.PaperRepository
}

func NewChapterService(
	repository repository.ChapterRepository,
	paperRepository repository.PaperRepository,
) ChapterService {
	return chapterService{repository: repository, paperRepository: paperRepository}
}

func (s chapterService) ListChapters(
	userId domain.UserIdObject,
	projectId domain.ProjectIdObject,
) ([]domain.ChapterEntity, *Error) {
	entries, rErr := s.repository.FetchChapters(userId.Value(), projectId.Value())
	if rErr != nil && rErr.Code() == repository.NotFoundError {
		return nil, Errorf(NotFoundError, "failed to list chapters: %w", rErr.Unwrap())
	}
	if rErr != nil {
		return nil, Errorf(RepositoryFailurePanic, "failed to fetch chapters: %w", rErr.Unwrap())
	}

	chapters := []domain.ChapterEntity{}
	for key, entry := range entries {
		chapter, err := s.entryToEntity(key, entry)
		if err != nil {
			return nil, err
		}

		chapters = append(chapters, *chapter)
	}

	sort.Slice(chapters, func(i, j int) bool {
		return chapters[i].Number().Value() < chapters[j].Number().Value()
	})

	return chapters, nil
}

func (s chapterService) CreateChapter(
	userId domain.UserIdObject,
	projectId domain.ProjectIdObject,
	chapter domain.ChapterWithoutAutofieldEntity,
) (*domain.ChapterEntity, *Error) {
	entryWithoutAutofield := record.ChapterWithoutAutofieldEntry{
		Name:   chapter.Name().Value(),
		Number: chapter.Number().Value(),
	}

	key, entry, rErr := s.repository.InsertChapter(userId.Value(), projectId.Value(), entryWithoutAutofield)
	if rErr != nil && rErr.Code() == repository.InvalidArgument {
		return nil, Errorf(InvalidArgument, "failed to create chapter: %w", rErr.Unwrap())
	}
	if rErr != nil && rErr.Code() == repository.NotFoundError {
		return nil, Errorf(NotFoundError, "failed to create chapter: %w", rErr.Unwrap())
	}
	if rErr != nil {
		return nil, Errorf(RepositoryFailurePanic, "failed to create chapter: %w", rErr.Unwrap())
	}

	paperEntryWithoutAutofield := record.PaperWithoutAutofieldEntry{
		Content: "",
	}

	_, _, rErr = s.paperRepository.InsertPaper(userId.Value(), projectId.Value(), key, paperEntryWithoutAutofield)
	if rErr != nil && rErr.Code() == repository.NotFoundError {
		return nil, Errorf(NotFoundError, "failed to create initial paper: %w", rErr.Unwrap())
	}
	if rErr != nil {
		return nil, Errorf(RepositoryFailurePanic, "failed to create initial paper: %w", rErr.Unwrap())
	}

	return s.entryToEntity(key, *entry)
}

func (s chapterService) UpdateChapter(
	userId domain.UserIdObject,
	projectId domain.ProjectIdObject,
	chapterId domain.ChapterIdObject,
	chapter domain.ChapterWithoutAutofieldEntity,
) (*domain.ChapterEntity, *Error) {
	entryWithoutAutofield := record.ChapterWithoutAutofieldEntry{
		Name:   chapter.Name().Value(),
		Number: chapter.Number().Value(),
	}

	entry, rErr := s.repository.UpdateChapter(
		userId.Value(),
		projectId.Value(),
		chapterId.Value(),
		entryWithoutAutofield,
	)
	if rErr != nil && rErr.Code() == repository.InvalidArgument {
		return nil, Errorf(InvalidArgument, "failed to update chapter: %w", rErr.Unwrap())
	}
	if rErr != nil && rErr.Code() == repository.NotFoundError {
		return nil, Errorf(NotFoundError, "failed to update chapter: %w", rErr.Unwrap())
	}
	if rErr != nil {
		return nil, Errorf(RepositoryFailurePanic, "failed to update chapter: %w", rErr.Unwrap())
	}

	return s.entryToEntity(chapterId.Value(), *entry)
}

func (s chapterService) entryToEntity(key string, entry record.ChapterEntry) (*domain.ChapterEntity, *Error) {
	id, err := domain.NewChapterIdObject(key)
	if err != nil {
		return nil, Errorf(DomainFailurePanic, "failed to convert entry to entity (id): %w", err)
	}
	name, err := domain.NewChapterNameObject(entry.Name)
	if err != nil {
		return nil, Errorf(DomainFailurePanic, "failed to convert entry to entity (name): %w", err)
	}
	number, err := domain.NewChapterNumberObject(entry.Number)
	if err != nil {
		return nil, Errorf(DomainFailurePanic, "failed to convert entry to entity (number): %w", err)
	}
	sections := make([]domain.SectionOfChapterEntity, len(entry.Sections))
	for i, sectionEntry := range entry.Sections {
		section, err := s.sectionEntryToEntity(sectionEntry)
		if err != nil {
			return nil, Errorf(DomainFailurePanic, "failed to convert entry to entity (section): %w", err.Unwrap())
		}
		sections[i] = *section
	}

	createdAt, err := domain.NewCreatedAtObject(entry.CreatedAt)
	if err != nil {
		return nil, Errorf(DomainFailurePanic, "failed to convert entry to entity (createdAt): %w", err)
	}
	updatedAt, err := domain.NewUpdatedAtObject(entry.UpdatedAt)
	if err != nil {
		return nil, Errorf(DomainFailurePanic, "failed to convert entry to entity (updatedAt): %w", err)
	}

	return domain.NewChapterEntity(*id, *name, *number, sections, *createdAt, *updatedAt), nil
}

func (s chapterService) sectionEntryToEntity(entry record.SectionEntry) (*domain.SectionOfChapterEntity, *Error) {
	id, err := domain.NewSectionIdObject(entry.Id)
	if err != nil {
		return nil, Errorf(DomainFailurePanic, "failed to convert entry to entity (id): %w", err)
	}
	name, err := domain.NewSectionNameObject(entry.Name)
	if err != nil {
		return nil, Errorf(DomainFailurePanic, "failed to convert entry to entity (name): %w", err)
	}
	createdAt, err := domain.NewCreatedAtObject(entry.CreatedAt)
	if err != nil {
		return nil, Errorf(DomainFailurePanic, "failed to convert entry to entity (createdAt): %w", err)
	}
	updatedAt, err := domain.NewUpdatedAtObject(entry.UpdatedAt)
	if err != nil {
		return nil, Errorf(DomainFailurePanic, "failed to convert entry to entity (updatedAt): %w", err)
	}

	return domain.NewSectionOfChapterEntity(*id, *name, *createdAt, *updatedAt), nil
}
