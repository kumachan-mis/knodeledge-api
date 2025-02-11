package repository

import (
	"cloud.google.com/go/firestore"
	"github.com/kumachan-mis/knodeledge-api/internal/db"
	"github.com/kumachan-mis/knodeledge-api/internal/document"
	"github.com/kumachan-mis/knodeledge-api/internal/record"
)

//go:generate mockgen -source=$GOFILE -destination=../../mock/$GOPACKAGE/mock_$GOFILE -package=$GOPACKAGE

const PaperCollection = "papers"

type PaperRepository interface {
	FetchPaper(
		userId string,
		projectId string,
		chapterId string,
	) (*record.PaperEntry, *Error)
	InsertPaper(
		userId string,
		projectId string,
		chapterId string,
		entry record.PaperWithoutAutofieldEntry,
	) (string, *record.PaperEntry, *Error)
	UpdatePaper(
		userId string,
		projectId string,
		chapterId string,
		entry record.PaperWithoutAutofieldEntry,
	) (*record.PaperEntry, *Error)
	DeletePaper(
		userId string,
		projectId string,
		chapterId string,
	) *Error
}

type paperRepository struct {
	client            firestore.Client
	chapterRepository ChapterRepository
}

func NewPaperRepository(client firestore.Client) PaperRepository {
	chapterRepository := NewChapterRepository(client)
	return paperRepository{client: client, chapterRepository: chapterRepository}
}

func (r paperRepository) FetchPaper(
	userId string,
	projectId string,
	chapterId string,
) (*record.PaperEntry, *Error) {
	_, rErr := r.chapterRepository.FetchChapter(userId, projectId, chapterId)
	if rErr != nil {
		return nil, rErr
	}

	snapshot, err := r.client.Collection(ProjectCollection).
		Doc(projectId).
		Collection(PaperCollection).
		Doc(chapterId).
		Get(db.FirestoreContext())

	if err != nil {
		return nil, Errorf(ReadFailurePanic, "failed to fetch paper: %w", err)
	}

	var values document.PaperValues
	err = snapshot.DataTo(&values)
	if err != nil {
		return nil, Errorf(ReadFailurePanic, "failed to convert snapshot to values: %w", err)
	}

	return r.valuesToEntry(values, userId), nil
}

func (r paperRepository) InsertPaper(
	userId string,
	projectId string,
	chapterId string,
	entry record.PaperWithoutAutofieldEntry,
) (string, *record.PaperEntry, *Error) {
	_, rErr := r.chapterRepository.FetchChapter(userId, projectId, chapterId)
	if rErr != nil {
		return "", nil, rErr
	}

	ref := r.client.Collection(ProjectCollection).
		Doc(projectId).
		Collection(PaperCollection).
		Doc(chapterId)

	_, err := ref.Set(db.FirestoreContext(), map[string]any{
		"content":   entry.Content,
		"createdAt": firestore.ServerTimestamp,
		"updatedAt": firestore.ServerTimestamp,
	})

	if err != nil {
		return "", nil, Errorf(WriteFailurePanic, "failed to insert paper: %w", err)
	}

	snapshpt, err := ref.Get(db.FirestoreContext())
	if err != nil {
		return "", nil, Errorf(ReadFailurePanic, "failed to fetch inserted paper: %w", err)
	}

	var values document.PaperValues
	err = snapshpt.DataTo(&values)
	if err != nil {
		return "", nil, Errorf(ReadFailurePanic, "failed to convert snapshot to values: %w", err)
	}

	return chapterId, r.valuesToEntry(values, userId), nil
}

func (r paperRepository) UpdatePaper(
	userId string,
	projectId string,
	chapterId string,
	entry record.PaperWithoutAutofieldEntry,
) (*record.PaperEntry, *Error) {
	_, rErr := r.chapterRepository.FetchChapter(userId, projectId, chapterId)
	if rErr != nil {
		return nil, rErr
	}

	ref := r.client.Collection(ProjectCollection).
		Doc(projectId).
		Collection(PaperCollection).
		Doc(chapterId)

	_, err := ref.Set(db.FirestoreContext(), map[string]any{
		"content":   entry.Content,
		"updatedAt": firestore.ServerTimestamp,
	}, firestore.MergeAll)

	if err != nil {
		return nil, Errorf(WriteFailurePanic, "failed to update paper: %w", err)
	}

	snapshot, err := ref.Get(db.FirestoreContext())
	if err != nil {
		return nil, Errorf(ReadFailurePanic, "failed to fetch updated paper: %w", err)
	}

	var values document.PaperValues
	err = snapshot.DataTo(&values)
	if err != nil {
		return nil, Errorf(ReadFailurePanic, "failed to convert snapshot to values: %w", err)
	}

	return r.valuesToEntry(values, userId), nil
}

func (r paperRepository) DeletePaper(
	userId string,
	projectId string,
	chapterId string,
) *Error {
	_, rErr := r.chapterRepository.FetchChapter(userId, projectId, chapterId)
	if rErr != nil {
		return rErr
	}

	_, err := r.client.Collection(ProjectCollection).
		Doc(projectId).
		Collection(PaperCollection).
		Doc(chapterId).
		Delete(db.FirestoreContext())

	if err != nil {
		return Errorf(WriteFailurePanic, "failed to delete paper: %w", err)
	}

	return nil
}

func (r paperRepository) valuesToEntry(
	values document.PaperValues,
	userId string,
) *record.PaperEntry {
	return &record.PaperEntry{
		Content:   values.Content,
		UserId:    userId,
		CreatedAt: values.CreatedAt,
		UpdatedAt: values.UpdatedAt,
	}
}
