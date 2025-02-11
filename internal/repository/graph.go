package repository

import (
	"errors"

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
	FetchGraph(
		userId string,
		projectId string,
		chapterId string,
		sectionId string,
	) (*record.GraphEntry, *Error)
	InsertGraphs(
		userId string,
		projectId string,
		chapterId string,
		entries []record.GraphWithoutAutofieldEntry,
	) ([]string, []record.GraphEntry, *Error)
	UpdateGraphContent(
		userId string,
		projectId string,
		chapterId string,
		sectionId string,
		entry record.GraphContentEntry,
	) (*record.GraphEntry, *Error)
	DeleteGraph(
		userId string,
		projectId string,
		chapterId string,
		sectionId string,
	) *Error
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
		return false, Errorf(ReadFailurePanic, "failed to fetch graphs: %w", err)
	}

	return len(snapshots) > 0, nil
}

func (r graphRepository) FetchGraph(
	userId string,
	projectId string,
	chapterId string,
	sectionId string,
) (*record.GraphEntry, *Error) {
	chapter, rErr := r.chapterRepository.FetchChapter(userId, projectId, chapterId)
	if rErr != nil {
		return nil, rErr
	}

	var section *record.SectionEntry
	for _, s := range chapter.Sections {
		if s.Id == sectionId {
			section = &s
			break
		}
	}

	snapshot, err := r.client.Collection(ProjectCollection).
		Doc(projectId).
		Collection(ChapterCollection).
		Doc(chapterId).
		Collection(GraphCollection).
		Doc(sectionId).
		Get(db.FirestoreContext())

	if err != nil && section == nil {
		return nil, Errorf(NotFoundError, "failed to fetch graph")
	}

	if err != nil && section != nil {
		err := errors.New("document.ChapterValues.sections have excessive elements")
		return nil, Errorf(ReadFailurePanic, "failed to convert values to entry: %w", err)
	} else if err == nil && section == nil {
		err := errors.New("document.ChapterValues.sections have insufficient elements")
		return nil, Errorf(ReadFailurePanic, "failed to convert values to entry: %w", err)
	}

	var values document.GraphValues
	err = snapshot.DataTo(&values)
	if err != nil {
		return nil, Errorf(ReadFailurePanic, "failed to convert snapshot to values: %v", err)
	}

	return r.valuesToEntry(values, section.Name, userId), nil
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
			"children":  r.childrenEntryToValues(entry.Children),
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

func (r graphRepository) UpdateGraphContent(
	userId string,
	projectId string,
	chapterId string,
	sectionId string,
	entry record.GraphContentEntry,
) (*record.GraphEntry, *Error) {
	chapter, rErr := r.chapterRepository.FetchChapter(userId, projectId, chapterId)
	if rErr != nil {
		return nil, rErr
	}

	var section *record.SectionEntry
	for _, s := range chapter.Sections {
		if s.Id == sectionId {
			section = &s
			break
		}
	}
	if section == nil {
		return nil, Errorf(NotFoundError, "failed to fetch graph")
	}

	ref := r.client.Collection(ProjectCollection).
		Doc(projectId).
		Collection(ChapterCollection).
		Doc(chapterId).
		Collection(GraphCollection).
		Doc(sectionId)

	_, err := ref.Set(db.FirestoreContext(), map[string]any{
		"paragraph": entry.Paragraph,
		"children":  r.childrenEntryToValues(entry.Children),
		"updatedAt": firestore.ServerTimestamp,
	}, firestore.MergeAll)

	if err != nil {
		return nil, Errorf(WriteFailurePanic, "failed to update graph: %v", err)
	}

	snapshot, err := ref.Get(db.FirestoreContext())
	if err != nil {
		return nil, Errorf(ReadFailurePanic, "failed to fetch updated graph: %w", err)
	}

	var values document.GraphValues
	err = snapshot.DataTo(&values)
	if err != nil {
		return nil, Errorf(ReadFailurePanic, "failed to convert snapshot to values: %w", err)
	}

	return r.valuesToEntry(values, section.Name, userId), nil
}

func (r graphRepository) DeleteGraph(
	userId string,
	projectId string,
	chapterId string,
	sectionId string,
) *Error {
	_, rErr := r.chapterRepository.FetchChapter(userId, projectId, chapterId)
	if rErr != nil {
		return rErr
	}

	_, err := r.client.Collection(ProjectCollection).
		Doc(projectId).
		Collection(ChapterCollection).
		Doc(chapterId).
		Collection(GraphCollection).
		Doc(sectionId).
		Delete(db.FirestoreContext())

	if err != nil {
		return Errorf(WriteFailurePanic, "failed to delete graph: %v", err)
	}

	return nil
}

func (r graphRepository) valuesToEntry(
	values document.GraphValues,
	name string,
	userId string,
) *record.GraphEntry {
	return &record.GraphEntry{
		Name:      name,
		Paragraph: values.Paragraph,
		Children:  r.childrenValuesToEntry(values.Children),
		UserId:    userId,
		CreatedAt: values.CreatedAt,
		UpdatedAt: values.UpdatedAt,
	}
}

func (r graphRepository) childrenValuesToEntry(
	values []document.GraphChildValues,
) []record.GraphChildEntry {
	entries := make([]record.GraphChildEntry, len(values))
	for i, value := range values {
		entries[i] = record.GraphChildEntry{
			Name:        value.Name,
			Relation:    value.Relation,
			Description: value.Descrition,
			Children:    r.childrenValuesToEntry(value.Children),
		}
	}
	return entries
}

func (r graphRepository) childrenEntryToValues(
	children []record.GraphChildEntry,
) []document.GraphChildValues {
	values := make([]document.GraphChildValues, len(children))
	for i, child := range children {
		values[i] = document.GraphChildValues{
			Name:       child.Name,
			Relation:   child.Relation,
			Descrition: child.Description,
			Children:   r.childrenEntryToValues(child.Children),
		}
	}
	return values
}
