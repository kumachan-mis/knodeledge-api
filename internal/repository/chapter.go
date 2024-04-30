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
	InsertChapter(projectId string, entry record.ChapterWithoutAutofieldEntry) (string, *record.ChapterEntry, *Error)
}

type chapterRepository struct {
	client firestore.Client
}

func NewChapterRepository(client firestore.Client) ChapterRepository {
	return chapterRepository{client: client}
}

func (r chapterRepository) FetchProjectChapters(userId string, projectId string) (map[string]record.ChapterEntry, *Error) {
	ref, rErr := r.projectDocumentRef(userId, projectId)
	if rErr != nil {
		return nil, rErr
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

func (r chapterRepository) InsertChapter(projectId string, entry record.ChapterWithoutAutofieldEntry) (string, *record.ChapterEntry, *Error) {
	ref, rErr := r.projectDocumentRef(entry.UserId, projectId)
	if rErr != nil {
		return "", nil, rErr
	}

	if entry.NextId != "" {
		_, err := ref.Collection(ChapterCollection).
			Doc(entry.NextId).
			Get(db.FirestoreContext())
		if err != nil {
			return "", nil, Errorf(InvalidArgument, "id of next chapter does not exist")
		}
	}

	prevSnapshot, _ := ref.Collection(ChapterCollection).
		Where("nextId", "==", entry.NextId).
		Limit(1).
		Documents(db.FirestoreContext()).
		Next()

	ref, _, err := ref.Collection(ChapterCollection).
		Add(db.FirestoreContext(), map[string]any{
			"name":      entry.Name,
			"nextId":    entry.NextId,
			"createdAt": firestore.ServerTimestamp,
			"updatedAt": firestore.ServerTimestamp,
		})
	if err != nil {
		return "", nil, Errorf(WriteFailurePanic, "failed to insert chapter: %w", err)
	}

	if prevSnapshot != nil {
		_, err = prevSnapshot.Ref.Update(db.FirestoreContext(), []firestore.Update{
			{Path: "nextId", Value: ref.ID},
		})
		if err != nil {
			return "", nil, Errorf(WriteFailurePanic, "failed to update previous chapter: %w", err)
		}
	}

	snapshot, err := ref.Get(db.FirestoreContext())
	if err != nil {
		return "", nil, Errorf(WriteFailurePanic, "failed to get inserted chapter: %w", err)
	}

	var values document.ChapterValues
	err = snapshot.DataTo(&values)
	if err != nil {
		return "", nil, Errorf(ReadFailurePanic, "failed to convert snapshot to values: %w", err)
	}

	return ref.ID, r.valuesToEntry(values, entry.UserId), nil
}

func (r chapterRepository) projectDocumentRef(userId string, projectId string) (*firestore.DocumentRef, *Error) {
	ref := r.client.Collection(ProjectCollection).
		Doc(projectId)

	snapshot, err := ref.Get(db.FirestoreContext())
	if err != nil {
		return nil, Errorf(InvalidArgument, "project document does not exist")
	}

	var projectValues document.ProjectValues
	err = snapshot.DataTo(&projectValues)
	if err != nil {
		return nil, Errorf(ReadFailurePanic, "failed to convert snapshot to values: %w", err)
	}

	if projectValues.UserId != userId {
		return nil, Errorf(InvalidArgument, "project document does not exist")
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
