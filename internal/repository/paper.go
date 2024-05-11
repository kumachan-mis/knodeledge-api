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
		projectId string,
		chapterId string,
		entry record.PaperWithoutAutofieldEntry,
	) (string, *record.PaperEntry, *Error)
}

type paperRepository struct {
	client firestore.Client
}

func NewPaperRepository(client firestore.Client) PaperRepository {
	return paperRepository{client: client}
}

func (r paperRepository) FetchPaper(
	userId string,
	projectId string,
	chapterId string,
) (*record.PaperEntry, *Error) {
	_, rErr := r.chapterValues(userId, projectId, chapterId)
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
	projectId string,
	chapterId string,
	entry record.PaperWithoutAutofieldEntry,
) (string, *record.PaperEntry, *Error) {
	_, rErr := r.chapterValues(entry.UserId, projectId, chapterId)
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

	return chapterId, r.valuesToEntry(values, entry.UserId), nil
}

func (r paperRepository) chapterValues(userId string, projectId string, chapterId string) (*document.ChapterValues, *Error) {
	projectSnapshot, err := r.client.Collection(ProjectCollection).
		Doc(projectId).
		Get(db.FirestoreContext())
	if err != nil {
		return nil, Errorf(NotFoundError, "project not found")
	}

	var projectValues document.ProjectValues
	err = projectSnapshot.DataTo(&projectValues)
	if err != nil {
		return nil, Errorf(ReadFailurePanic, "failed to convert snapshot to values: %w", err)
	}

	if projectValues.UserId != userId {
		return nil, Errorf(NotFoundError, "project not found")
	}

	chapterSnapshot, err := r.client.Collection(ProjectCollection).
		Doc(projectId).
		Collection(ChapterCollection).
		Doc(chapterId).
		Get(db.FirestoreContext())
	if err != nil {
		return nil, Errorf(NotFoundError, "chapter not found")
	}

	var chapterValues document.ChapterValues
	err = chapterSnapshot.DataTo(&chapterValues)
	if err != nil {
		return nil, Errorf(ReadFailurePanic, "failed to convert snapshot to values: %w", err)
	}

	return &chapterValues, nil
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
