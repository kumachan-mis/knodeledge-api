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
	if rErr != nil && rErr.Code() == repository.InvalidArgument {
		return nil, Errorf(InvalidArgument, "failed to list chapters: %w", rErr.Unwrap())
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
	entryyWithoutAutofield := record.ChapterWithoutAutofieldEntry{
		Name:   chapter.Name().Value(),
		Number: chapter.Number().Value(),
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

func (s chapterService) UpdateChapter(
	userId domain.UserIdObject,
	projectId domain.ProjectIdObject,
	chapterId domain.ChapterIdObject,
	chapter domain.ChapterWithoutAutofieldEntity,
) (*domain.ChapterEntity, *Error) {
	entryyWithoutAutofield := record.ChapterWithoutAutofieldEntry{
		Name:   chapter.Name().Value(),
		Number: chapter.Number().Value(),
		UserId: userId.Value(),
	}

	entry, rErr := s.repository.UpdateChapter(projectId.Value(), chapterId.Value(), entryyWithoutAutofield)
	if rErr != nil && rErr.Code() == repository.InvalidArgument {
		return nil, Errorf(InvalidArgument, "failed to update chapter: %w", rErr.Unwrap())
	}
	if rErr != nil && rErr.Code() == repository.NotFoundError {
		return nil, Errorf(NotFoundError, "failed to update chapter")
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
	createdAt, err := domain.NewCreatedAtObject(entry.CreatedAt)
	if err != nil {
		return nil, Errorf(DomainFailurePanic, "failed to convert entry to entity (createdAt): %w", err)
	}
	updatedAt, err := domain.NewUpdatedAtObject(entry.UpdatedAt)
	if err != nil {
		return nil, Errorf(DomainFailurePanic, "failed to convert entry to entity (updatedAt): %w", err)
	}

	return domain.NewChapterEntity(*id, *name, *number, *createdAt, *updatedAt), nil
}
