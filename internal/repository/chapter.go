package repository

import (
	"cloud.google.com/go/firestore"
	"github.com/kumachan-mis/knodeledge-api/internal/db"
	"github.com/kumachan-mis/knodeledge-api/internal/document"
	"github.com/kumachan-mis/knodeledge-api/internal/record"
)

//go:generate mockgen -source=$GOFILE -destination=../../mock/$GOPACKAGE/mock_$GOFILE -package=$GOPACKAGE

const ChapterCollection = "chapters"

type ChapterRepository interface {
	FetchProjectChapters(userId string, projectId string) (map[string]record.ChapterEntry, *Error)
}

type chapterRepository struct {
	client firestore.Client
}

func NewChapterRepository(client firestore.Client) ChapterRepository {
	return chapterRepository{client: client}
}

func (r chapterRepository) FetchProjectChapters(userId string, projectId string) (map[string]record.ChapterEntry, *Error) {
	ref, err := r.projectDocumentRef(userId, projectId)
	if err != nil {
		return nil, err
	}

	iter := ref.
		Collection(ChapterCollection).
		Documents(db.FirestoreContext())

	entries := make(map[string]record.ChapterEntry)

	for {
		snapshot, err := iter.Next()
		if err != nil {
			break
		}

		var values document.ChapterValues
		err = snapshot.DataTo(&values)
		if err != nil {
			return nil, Errorf(ReadFailurePanic, "failed to convert snapshot to values: %w", err)
		}

		entries[snapshot.Ref.ID] = *r.valuesToEntry(values, userId)
	}

	return entries, nil
}

func (r chapterRepository) projectDocumentRef(userId string, projectId string) (*firestore.DocumentRef, *Error) {
	ref := r.client.Collection(ProjectCollection).
		Doc(projectId)

	snapshot, err := ref.Get(db.FirestoreContext())
	if err != nil {
		return nil, Errorf(NotFoundError, "project document not found")
	}

	var projectValues document.ProjectValues
	err = snapshot.DataTo(&projectValues)
	if err != nil {
		return nil, Errorf(ReadFailurePanic, "failed to convert snapshot to values: %w", err)
	}

	if projectValues.UserId != userId {
		return nil, Errorf(NotFoundError, "project document not found")
	}

	return ref, nil
}

func (r chapterRepository) valuesToEntry(values document.ChapterValues, userId string) *record.ChapterEntry {
	return &record.ChapterEntry{
		Name:      values.Name,
		NextId:    values.NextId,
		UserId:    userId,
		CreatedAt: values.CreatedAt,
		UpdatedAt: values.UpdatedAt,
	}
}
