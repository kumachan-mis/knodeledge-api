package service

import (
	"github.com/kumachan-mis/knodeledge-api/internal/domain"
	"github.com/kumachan-mis/knodeledge-api/internal/record"
	"github.com/kumachan-mis/knodeledge-api/internal/repository"
)

//go:generate mockgen -source=$GOFILE -destination=../../mock/$GOPACKAGE/mock_$GOFILE -package=$GOPACKAGE

type PaperService interface {
	FindPaper(
		userId domain.UserIdObject,
		projectId domain.ProjectIdObject,
		chapterId domain.ChapterIdObject,
	) (*domain.PaperEntity, *Error)
	CreatePaper(
		userId domain.UserIdObject,
		projectId domain.ProjectIdObject,
		chapterId domain.ChapterIdObject,
		paper domain.PaperWithoutAutofieldEntity,
	) (*domain.PaperEntity, *Error)
}

type paperService struct {
	repository repository.PaperRepository
}

func NewPaperService(repository repository.PaperRepository) PaperService {
	return paperService{repository: repository}
}

func (s paperService) FindPaper(
	userId domain.UserIdObject,
	projectId domain.ProjectIdObject,
	chapterId domain.ChapterIdObject,
) (*domain.PaperEntity, *Error) {
	entry, rErr := s.repository.FetchPaper(userId.Value(), projectId.Value(), chapterId.Value())
	if rErr != nil && rErr.Code() == repository.NotFoundError {
		return nil, Errorf(NotFoundError, "failed to find paper: %w", rErr.Unwrap())
	}
	if rErr != nil {
		return nil, Errorf(RepositoryFailurePanic, "failed to fetch paper: %w", rErr.Unwrap())
	}

	return s.entryToEntity(chapterId.Value(), *entry)
}

func (s paperService) CreatePaper(
	userId domain.UserIdObject,
	projectId domain.ProjectIdObject,
	chapterId domain.ChapterIdObject,
	paper domain.PaperWithoutAutofieldEntity,
) (*domain.PaperEntity, *Error) {
	entryWithoutAutofield := record.PaperWithoutAutofieldEntry{
		Content: paper.Content().Value(),
		UserId:  userId.Value(),
	}

	key, entry, rErr := s.repository.InsertPaper(projectId.Value(), chapterId.Value(), entryWithoutAutofield)
	if rErr != nil && rErr.Code() == repository.NotFoundError {
		return nil, Errorf(NotFoundError, "failed to create paper: %w", rErr.Unwrap())
	}
	if rErr != nil {
		return nil, Errorf(RepositoryFailurePanic, "failed to insert paper: %w", rErr.Unwrap())
	}

	return s.entryToEntity(key, *entry)
}

func (s paperService) entryToEntity(key string, entry record.PaperEntry) (*domain.PaperEntity, *Error) {
	id, err := domain.NewPaperIdObject(key)
	if err != nil {
		return nil, Errorf(DomainFailurePanic, "failed to convert entry to entity (id): %w", err)
	}
	content, err := domain.NewPaperContentObject(entry.Content)
	if err != nil {
		return nil, Errorf(DomainFailurePanic, "failed to convert entry to entity (content): %w", err)
	}
	createdAt, err := domain.NewCreatedAtObject(entry.CreatedAt)
	if err != nil {
		return nil, Errorf(DomainFailurePanic, "failed to convert entry to entity (createdAt): %w", err)
	}
	updatedAt, err := domain.NewUpdatedAtObject(entry.UpdatedAt)
	if err != nil {
		return nil, Errorf(DomainFailurePanic, "failed to convert entry to entity (updatedAt): %w", err)
	}

	return domain.NewPaperEntity(*id, *content, *createdAt, *updatedAt), nil
}
