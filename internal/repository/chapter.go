package repository

import (
	"errors"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/kumachan-mis/knodeledge-api/internal/db"
	"github.com/kumachan-mis/knodeledge-api/internal/document"
	"github.com/kumachan-mis/knodeledge-api/internal/record"
)

//go:generate mockgen -source=$GOFILE -destination=../../mock/$GOPACKAGE/mock_$GOFILE -package=$GOPACKAGE

const ChapterCollection = "chapters"

type ChapterRepository interface {
	FetchChapters(
		userId string,
		projectId string,
	) (map[string]record.ChapterEntry, *Error)
	FetchChapter(
		userId string,
		projectId string,
		chapterId string,
	) (*record.ChapterEntry, *Error)
	InsertChapter(
		userId string,
		projectId string,
		entry record.ChapterWithoutAutofieldEntry,
	) (string, *record.ChapterEntry, *Error)
	UpdateChapter(
		userId string,
		projectId string,
		chapterId string,
		entry record.ChapterWithoutAutofieldEntry,
	) (*record.ChapterEntry, *Error)
	UpdateChapterSections(
		userId string,
		projectId string,
		chapterId string,
		entries []record.SectionWithoutAutofieldEntry,
	) ([]record.SectionEntry, *Error)
	DeleteChapter(
		userId string,
		projectId string,
		chapterId string,
	) *Error
}

type chapterRepository struct {
	client firestore.Client
}

func NewChapterRepository(client firestore.Client) ChapterRepository {
	return chapterRepository{client: client}
}

func (r chapterRepository) FetchChapters(
	userId string,
	projectId string,
) (map[string]record.ChapterEntry, *Error) {
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
			err = errors.New("document.ProjectValues.chapterIds have insufficient elements")
			return nil, Errorf(ReadFailurePanic, "failed to convert values to entry: %w", err)
		}

		entries[snapshot.Ref.ID] = *r.valuesToEntry(values, number, userId)
	}

	if len(entries) != len(projectValues.ChapterIds) {
		err := errors.New("document.ProjectValues.chapterIds have excessive elements")
		return nil, Errorf(ReadFailurePanic, "failed to convert values to entry: %w", err)
	}
	return entries, nil
}

func (r chapterRepository) FetchChapter(
	userId string,
	projectId string,
	chapterId string,
) (*record.ChapterEntry, *Error) {
	projectValues, rErr := r.projectValues(userId, projectId)
	if rErr != nil {
		return nil, rErr
	}

	chapterNumbers := make(map[string]int)
	for i, chapterId := range projectValues.ChapterIds {
		chapterNumbers[chapterId] = i + 1
	}
	number, ok := chapterNumbers[chapterId]

	snapshot, err := r.client.Collection(ProjectCollection).
		Doc(projectId).
		Collection(ChapterCollection).
		Doc(chapterId).
		Get(db.FirestoreContext())

	if err != nil && !ok {
		return nil, Errorf(NotFoundError, "failed to fetch chapter")
	}

	if err != nil && ok {
		err := errors.New("document.ProjectValues.chapterIds have excessive elements")
		return nil, Errorf(ReadFailurePanic, "failed to convert values to entry: %w", err)
	} else if err == nil && !ok {
		err := errors.New("document.ProjectValues.chapterIds have insufficient elements")
		return nil, Errorf(ReadFailurePanic, "failed to convert values to entry: %w", err)
	}

	var values document.ChapterValues
	err = snapshot.DataTo(&values)
	if err != nil {
		return nil, Errorf(ReadFailurePanic, "failed to convert snapshot to values: %w", err)
	}

	return r.valuesToEntry(values, number, userId), nil
}

func (r chapterRepository) InsertChapter(
	userId string,
	projectId string,
	entry record.ChapterWithoutAutofieldEntry,
) (string, *record.ChapterEntry, *Error) {
	projectValues, rErr := r.projectValues(userId, projectId)
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
			"sections":  []map[string]any{},
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
		return "", nil, Errorf(WriteFailurePanic, "failed to fetch inserted chapter: %w", err)
	}

	var values document.ChapterValues
	err = snapshot.DataTo(&values)
	if err != nil {
		return "", nil, Errorf(ReadFailurePanic, "failed to convert snapshot to values: %w", err)
	}

	return ref.ID, r.valuesToEntry(values, entry.Number, userId), nil
}

func (r chapterRepository) UpdateChapter(
	userId string,
	projectId string,
	chapterId string,
	entry record.ChapterWithoutAutofieldEntry,
) (*record.ChapterEntry, *Error) {
	projectValues, rErr := r.projectValues(userId, projectId)
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
		return nil, Errorf(ReadFailurePanic, "failed to fetch updated chapter: %w", err)
	}

	var values document.ChapterValues
	err = snapshot.DataTo(&values)
	if err != nil {
		return nil, Errorf(ReadFailurePanic, "failed to convert snapshot to values: %w", err)
	}

	return r.valuesToEntry(values, entry.Number, userId), nil
}

func (r chapterRepository) UpdateChapterSections(
	userId string,
	projectId string,
	chapterId string,
	entries []record.SectionWithoutAutofieldEntry,
) ([]record.SectionEntry, *Error) {
	_, rErr := r.projectValues(userId, projectId)
	if rErr != nil {
		return nil, rErr
	}

	sections := make([]map[string]any, len(entries))
	for i, sectionEntry := range entries {
		sections[i] = map[string]any{
			"id":   sectionEntry.Id,
			"name": sectionEntry.Name,
		}
	}

	_, err := r.client.Collection(ProjectCollection).
		Doc(projectId).
		Collection(ChapterCollection).
		Doc(chapterId).
		Update(db.FirestoreContext(), []firestore.Update{
			{Path: "sections", Value: sections},
			{Path: "updatedAt", Value: firestore.ServerTimestamp},
		})
	if err != nil {
		return nil, Errorf(NotFoundError, "failed to update sections of chapter")
	}

	snapshot, err := r.client.Collection(ProjectCollection).
		Doc(projectId).
		Collection(ChapterCollection).
		Doc(chapterId).
		Get(db.FirestoreContext())
	if err != nil {
		return nil, Errorf(ReadFailurePanic, "failed to fetch updated chapter")
	}

	var values document.ChapterValues
	err = snapshot.DataTo(&values)
	if err != nil {
		return nil, Errorf(ReadFailurePanic, "failed to convert snapshot to values: %w", err)
	}

	sectionEntries := make([]record.SectionEntry, len(values.Sections))
	for i, sectionValues := range values.Sections {
		sectionEntries[i] = *r.sectionValuesToEntry(sectionValues, userId, values.CreatedAt, values.UpdatedAt)
	}

	return sectionEntries, nil
}

func (r chapterRepository) DeleteChapter(
	userId string,
	projectId string,
	chapterId string,
) *Error {
	projectValues, rErr := r.projectValues(userId, projectId)
	if rErr != nil {
		return rErr
	}

	ref := r.client.Collection(ProjectCollection).
		Doc(projectId).
		Collection(ChapterCollection).
		Doc(chapterId)

	if _, err := ref.Get(db.FirestoreContext()); err != nil {
		return Errorf(NotFoundError, "failed to fetch chapter")
	}

	_, err := ref.Delete(db.FirestoreContext())

	if err != nil {
		return Errorf(WriteFailurePanic, "failed to delete chapter: %w", err)
	}

	updatedChapterIds := []string{}
	for _, id := range projectValues.ChapterIds {
		if id == chapterId {
			continue
		}
		updatedChapterIds = append(updatedChapterIds, id)
	}

	_, err = r.client.Collection(ProjectCollection).
		Doc(projectId).
		Update(db.FirestoreContext(), []firestore.Update{
			{Path: "chapterIds", Value: updatedChapterIds},
		})
	if err != nil {
		return Errorf(WriteFailurePanic, "failed to update chapter ids: %w", err)
	}

	return nil
}

func (r chapterRepository) projectValues(userId string, projectId string) (*document.ProjectValues, *Error) {
	ref := r.client.Collection(ProjectCollection).
		Doc(projectId)

	snapshot, err := ref.Get(db.FirestoreContext())
	if err != nil {
		return nil, Errorf(NotFoundError, "failed to fetch project")
	}

	var projectValues document.ProjectValues
	err = snapshot.DataTo(&projectValues)
	if err != nil {
		return nil, Errorf(ReadFailurePanic, "failed to convert snapshot to values: %w", err)
	}

	if projectValues.UserId != userId {
		return nil, Errorf(NotFoundError, "failed to fetch project")
	}

	if projectValues.ChapterIds == nil {
		projectValues.ChapterIds = []string{}
	}
	return &projectValues, nil
}

func (r chapterRepository) valuesToEntry(
	values document.ChapterValues,
	number int,
	userId string,
) *record.ChapterEntry {
	sections := make([]record.SectionEntry, len(values.Sections))
	for i, sectionValues := range values.Sections {
		sections[i] = record.SectionEntry{
			Id:        sectionValues.Id,
			Name:      sectionValues.Name,
			UserId:    userId,
			CreatedAt: values.CreatedAt,
			UpdatedAt: values.UpdatedAt,
		}
	}

	return &record.ChapterEntry{
		Name:      values.Name,
		Number:    number,
		Sections:  sections,
		UserId:    userId,
		CreatedAt: values.CreatedAt,
		UpdatedAt: values.UpdatedAt,
	}
}

func (r chapterRepository) sectionValuesToEntry(
	values document.SectionValues,
	userId string,
	createdAt time.Time,
	updatedAt time.Time,
) *record.SectionEntry {
	return &record.SectionEntry{
		Id:        values.Id,
		Name:      values.Name,
		UserId:    userId,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
