package repository

import (
	"fmt"
	"reflect"

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
	UpdateChapter(projectId string, chapterId string, entry record.ChapterWithoutAutofieldEntry) (*record.ChapterEntry, *Error)
}

type chapterRepository struct {
	client firestore.Client
}

func NewChapterRepository(client firestore.Client) ChapterRepository {
	return chapterRepository{client: client}
}

func (r chapterRepository) FetchProjectChapters(userId string, projectId string) (map[string]record.ChapterEntry, *Error) {
	projectValues, rErr := r.projectValues(userId, projectId)
	if rErr != nil {
		return nil, rErr
	}

	chapterNumbers := make(map[string]int)
	for i, chapterId := range projectValues.ChapterIds {
		chapterNumbers[chapterId] = i + 1
	}

	iter := r.client.Collection(ProjectCollection).
		Doc(projectId).
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

		number, ok := chapterNumbers[snapshot.Ref.ID]
		if !ok {
			err = fmt.Errorf("%v.chapterIds have deficient elements", reflect.TypeOf(*projectValues))
			return nil, Errorf(ReadFailurePanic, "failed to convert values to entry: %w", err)
		}

		entries[snapshot.Ref.ID] = *r.valuesToEntry(values, number, userId)
	}

	if len(entries) != len(projectValues.ChapterIds) {
		err := fmt.Errorf("%v.chapterIds have excessive elements", reflect.TypeOf(*projectValues))
		return nil, Errorf(ReadFailurePanic, "failed to convert values to entry: %w", err)
	}
	return entries, nil
}

func (r chapterRepository) InsertChapter(projectId string, entry record.ChapterWithoutAutofieldEntry) (string, *record.ChapterEntry, *Error) {
	projectValues, rErr := r.projectValues(entry.UserId, projectId)
	if rErr != nil {
		return "", nil, rErr
	}

	if entry.Number > len(projectValues.ChapterIds)+1 {
		return "", nil, Errorf(InvalidArgument, "chapter number is too large")
	}

	ref, _, err := r.client.Collection(ProjectCollection).
		Doc(projectId).
		Collection(ChapterCollection).
		Add(db.FirestoreContext(), map[string]any{
			"name":      entry.Name,
			"createdAt": firestore.ServerTimestamp,
			"updatedAt": firestore.ServerTimestamp,
		})
	if err != nil {
		return "", nil, Errorf(WriteFailurePanic, "failed to insert chapter: %w", err)
	}

	updatedChapterIds := make([]string, len(projectValues.ChapterIds)+1)
	copy(updatedChapterIds[:entry.Number-1], projectValues.ChapterIds[:entry.Number-1])
	updatedChapterIds[entry.Number-1] = ref.ID
	copy(updatedChapterIds[entry.Number:], projectValues.ChapterIds[entry.Number-1:])

	_, err = r.client.Collection(ProjectCollection).
		Doc(projectId).
		Update(db.FirestoreContext(), []firestore.Update{
			{Path: "chapterIds", Value: updatedChapterIds},
		})
	if err != nil {
		return "", nil, Errorf(WriteFailurePanic, "failed to update chapter ids: %w", err)
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

	return ref.ID, r.valuesToEntry(values, entry.Number, entry.UserId), nil
}

func (r chapterRepository) UpdateChapter(projectId string, chapterId string, entry record.ChapterWithoutAutofieldEntry) (*record.ChapterEntry, *Error) {
	projectValues, rErr := r.projectValues(entry.UserId, projectId)
	if rErr != nil {
		return nil, rErr
	}

	if entry.Number > len(projectValues.ChapterIds) {
		return nil, Errorf(InvalidArgument, "chapter number is too large")
	}

	_, err := r.client.Collection(ProjectCollection).
		Doc(projectId).
		Collection(ChapterCollection).
		Doc(chapterId).
		Update(db.FirestoreContext(), []firestore.Update{
			{Path: "name", Value: entry.Name},
			{Path: "updatedAt", Value: firestore.ServerTimestamp},
		})
	if err != nil {
		return nil, Errorf(NotFoundError, "failed to update chapter")
	}

	chapterIdsWithoutUpdated := []string{}
	updatedNumber := 0
	for i, id := range projectValues.ChapterIds {
		if id == chapterId {
			updatedNumber = i + 1
			continue
		}
		chapterIdsWithoutUpdated = append(chapterIdsWithoutUpdated, id)
	}

	if updatedNumber != entry.Number {
		updatedChapterIds := make([]string, len(projectValues.ChapterIds))
		copy(updatedChapterIds[:entry.Number-1], chapterIdsWithoutUpdated[:entry.Number-1])
		updatedChapterIds[entry.Number-1] = chapterId
		copy(updatedChapterIds[entry.Number:], chapterIdsWithoutUpdated[entry.Number-1:])

		_, err = r.client.Collection(ProjectCollection).
			Doc(projectId).
			Update(db.FirestoreContext(), []firestore.Update{
				{Path: "chapterIds", Value: updatedChapterIds},
			})
		if err != nil {
			return nil, Errorf(WriteFailurePanic, "failed to update chapter ids: %w", err)
		}
	}

	snapshot, err := r.client.Collection(ProjectCollection).
		Doc(projectId).
		Collection(ChapterCollection).
		Doc(chapterId).
		Get(db.FirestoreContext())
	if err != nil {
		return nil, Errorf(ReadFailurePanic, "failed to get updated chapter: %w", err)
	}

	var values document.ChapterValues
	err = snapshot.DataTo(&values)
	if err != nil {
		return nil, Errorf(ReadFailurePanic, "failed to convert snapshot to values: %w", err)
	}

	return r.valuesToEntry(values, entry.Number, entry.UserId), nil
}

func (r chapterRepository) projectValues(userId string, projectId string) (*document.ProjectWithChapterIdsValues, *Error) {
	ref := r.client.Collection(ProjectCollection).
		Doc(projectId)

	snapshot, err := ref.Get(db.FirestoreContext())
	if err != nil {
		return nil, Errorf(NotFoundError, "project not found")
	}

	var projectValues document.ProjectWithChapterIdsValues
	err = snapshot.DataTo(&projectValues)
	if err != nil {
		return nil, Errorf(ReadFailurePanic, "failed to convert snapshot to values: %w", err)
	}

	if projectValues.UserId != userId {
		return nil, Errorf(NotFoundError, "project not found")
	}

	if projectValues.ChapterIds == nil {
		projectValues.ChapterIds = []string{}
	}
	return &projectValues, nil
}

func (r chapterRepository) valuesToEntry(values document.ChapterValues, number int, userId string) *record.ChapterEntry {
	return &record.ChapterEntry{
		Name:      values.Name,
		Number:    number,
		UserId:    userId,
		CreatedAt: values.CreatedAt,
		UpdatedAt: values.UpdatedAt,
	}
}
