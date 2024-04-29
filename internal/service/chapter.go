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
}

type chapterService struct {
	repository repository.ChapterRepository
}

func NewChapterService(repository repository.ChapterRepository) ChapterService {
	return chapterService{repository: repository}
}

func (s chapterService) ListChapters(
	userId domain.UserIdObject,
	projectId domain.ProjectIdObject,
) ([]domain.ChapterEntity, *Error) {
	entries, rErr := s.repository.FetchProjectChapters(userId.Value(), projectId.Value())
	if rErr != nil && rErr.Code() == repository.NotFoundError {
		return []domain.ChapterEntity{}, nil
	}

	if rErr != nil {
		return nil, Errorf(RepositoryFailurePanic, "failed to fetch project chapters: %w", rErr.Unwrap())
	}

	chapters := []domain.ChapterEntity{}
	for key, entry := range entries {
		chapter, err := s.entryToEntity(key, entry)
		if err != nil {
			return nil, err
		}

		chapters = append(chapters, *chapter)
	}

	order := s.entriesOrder(entries)
	sort.Slice(chapters, func(i, j int) bool {
		return order[chapters[i].Id().Value()] < order[chapters[j].Id().Value()]
	})

	return chapters, nil
}

func (s chapterService) CreateChapter(
	userId domain.UserIdObject,
	projectId domain.ProjectIdObject,
	chapter domain.ChapterWithoutAutofieldEntity,
) (*domain.ChapterEntity, *Error) {
	entryyWithoutAutofield := record.ChapterWithoutAutofieldEntry{
		Name:   chapter.Name().Value(),
		NextId: chapter.NextId().Value(),
		UserId: userId.Value(),
	}

	key, entry, rErr := s.repository.InsertChapter(projectId.Value(), entryyWithoutAutofield)
	if rErr != nil && rErr.Code() == repository.InvalidArgument {
		return nil, Errorf(InvalidArgument, "failed to create chapter: %w", rErr.Unwrap())
	}
	if rErr != nil {
		return nil, Errorf(RepositoryFailurePanic, "failed to create chapter: %w", rErr.Unwrap())
	}

	return s.entryToEntity(key, *entry)
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
	nextId, err := domain.NewChapterNextIdObject(entry.NextId)
	if err != nil {
		return nil, Errorf(DomainFailurePanic, "failed to convert entry to entity (nextId): %w", err)
	}
	createdAt, err := domain.NewCreatedAtObject(entry.CreatedAt)
	if err != nil {
		return nil, Errorf(DomainFailurePanic, "failed to convert entry to entity (createdAt): %w", err)
	}
	updatedAt, err := domain.NewUpdatedAtObject(entry.UpdatedAt)
	if err != nil {
		return nil, Errorf(DomainFailurePanic, "failed to convert entry to entity (updatedAt): %w", err)
	}

	return domain.NewChapterEntity(*id, *name, *nextId, *createdAt, *updatedAt), nil
}

func (s chapterService) entriesOrder(entries map[string]record.ChapterEntry) map[string]int {
	prevs := make(map[string]string)
	for key, entry := range entries {
		prevs[entry.NextId] = key
	}

	order := make(map[string]int)

	id := ""
	for i := len(entries); i > 0; i-- {
		order[id] = i
		id = prevs[id]
	}

	return order
}
