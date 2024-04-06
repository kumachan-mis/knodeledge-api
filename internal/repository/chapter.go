package repository

import (
	"cloud.google.com/go/firestore"
	"github.com/kumachan-mis/knodeledge-api/internal/db"
	"github.com/kumachan-mis/knodeledge-api/internal/record"
)

//go:generate mockgen -source=$GOFILE -destination=../../mock/$GOPACKAGE/mock_$GOFILE -package=$GOPACKAGE

const ChapterCollection = "chapters"

type ChapterRepository interface {
	FetchProjectChapters(projectId string) (map[string]record.ChapterEntry, *Error)
}

type chapterRepository struct {
	client firestore.Client
}

func NewChapterRepository(client firestore.Client) ChapterRepository {
	return chapterRepository{client: client}
}

func (r chapterRepository) FetchProjectChapters(projectId string) (map[string]record.ChapterEntry, *Error) {
	docref := r.client.Collection(ProjectCollection).
		Doc(projectId)

	if _, err := docref.Get(db.FirestoreContext()); err != nil {
		return nil, Errorf(NotFoundError, "parent collection not found")
	}

	iter := docref.
		Collection(ChapterCollection).
		Documents(db.FirestoreContext())

	entries := make(map[string]record.ChapterEntry)

	for {
		snapshot, err := iter.Next()
		if err != nil {
			break
		}

		var entry record.ChapterEntry
		err = snapshot.DataTo(&entry)
		if err != nil {
			return nil, Errorf(ReadFailurePanic, "failed to convert snapshot to entry: %w", err)
		}

		entries[snapshot.Ref.ID] = entry
	}

	return entries, nil
}
