package repository

import (
	"cloud.google.com/go/firestore"
	"github.com/kumachan-mis/knodeledge-api/internal/db"
	"github.com/kumachan-mis/knodeledge-api/internal/document"
	"github.com/kumachan-mis/knodeledge-api/internal/record"
)

//go:generate mockgen -source=$GOFILE -destination=../../mock/$GOPACKAGE/mock_$GOFILE -package=$GOPACKAGE

const GraphCollection = "graphs"

type GraphRepository interface {
	GraphExists(
		userId string,
		projectId string,
		chapterId string,
	) (bool, *Error)
	InsertGraphs(
		userId string,
		projectId string,
		chapterId string,
		entries []record.GraphWithoutAutofieldEntry,
	) ([]string, []record.GraphEntry, *Error)
}

type graphRepository struct {
	client            firestore.Client
	chapterRepository ChapterRepository
}

func NewGraphRepository(client firestore.Client) GraphRepository {
	chapterRepository := NewChapterRepository(client)
	return graphRepository{client: client, chapterRepository: chapterRepository}
}

func (r graphRepository) GraphExists(
	userId string,
	projectId string,
	chapterId string,
) (bool, *Error) {
	_, rErr := r.chapterRepository.FetchChapter(userId, projectId, chapterId)
	if rErr != nil {
		return false, rErr
	}

	snapshots, err := r.client.Collection(ProjectCollection).
		Doc(projectId).
		Collection(ChapterCollection).
		Doc(chapterId).
		Collection(GraphCollection).
		Limit(1).
		Documents(db.FirestoreContext()).
		GetAll()

	if err != nil {
		return false, Errorf(ReadFailurePanic, "failed to fetch graphs: %v", err)
	}

	return len(snapshots) > 0, nil
}

func (r graphRepository) InsertGraphs(
	userId string,
	projectId string,
	chapterId string,
	entries []record.GraphWithoutAutofieldEntry,
) ([]string, []record.GraphEntry, *Error) {
	_, rErr := r.chapterRepository.FetchChapter(userId, projectId, chapterId)
	if rErr != nil {
		return nil, nil, rErr
	}

	collRef := r.client.Collection(ProjectCollection).
		Doc(projectId).
		Collection(ChapterCollection).
		Doc(chapterId).
		Collection(GraphCollection)

	bw := r.client.BulkWriter(db.FirestoreContext())

	ids := make([]string, len(entries))
	docRefs := make([]*firestore.DocumentRef, len(entries))
	for i, entry := range entries {
		docRef := collRef.NewDoc()

		_, err := bw.Create(docRef, map[string]any{
			"paragraph": entry.Paragraph,
			"createdAt": firestore.ServerTimestamp,
			"updatedAt": firestore.ServerTimestamp,
		})
		if err != nil {
			return nil, nil, Errorf(WriteFailurePanic, "failed to insert graph: %v", err)
		}

		ids[i] = docRef.ID
		docRefs[i] = docRef
	}

	bw.End()

	snapshots, err := r.client.GetAll(db.FirestoreContext(), docRefs)
	if err != nil {
		return nil, nil, Errorf(ReadFailurePanic, "failed to fetch created graphs: %v", err)
	}

	res := make([]record.GraphEntry, len(snapshots))
	for i, snapshot := range snapshots {
		var values document.GraphValues
		err = snapshot.DataTo(&values)

		if err != nil {
			return nil, nil, Errorf(ReadFailurePanic, "failed to convert snapshot to values: %v", err)
		}

		res[i] = *r.valuesToEntry(values, entries[i].Name, userId)
	}

	return ids, res, nil
}

func (r graphRepository) valuesToEntry(
	values document.GraphValues,
	name string,
	userId string,
) *record.GraphEntry {
	return &record.GraphEntry{
		Paragraph: values.Paragraph,
		Name:      name,
		UserId:    userId,
		CreatedAt: values.CreatedAt,
		UpdatedAt: values.UpdatedAt,
	}
}
