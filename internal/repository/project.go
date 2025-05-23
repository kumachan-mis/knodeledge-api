package repository

import (
	"cloud.google.com/go/firestore"
	"github.com/kumachan-mis/knodeledge-api/internal/db"
	"github.com/kumachan-mis/knodeledge-api/internal/document"
	"github.com/kumachan-mis/knodeledge-api/internal/record"
)

//go:generate mockgen -source=$GOFILE -destination=../../mock/$GOPACKAGE/mock_$GOFILE -package=$GOPACKAGE

const ProjectCollection = "projects"

type ProjectRepository interface {
	FetchProjects(
		userId string,
	) (map[string]record.ProjectEntry, *Error)
	FetchProject(
		userId string,
		projectId string,
	) (*record.ProjectEntry, *Error)
	InsertProject(
		userId string,
		entry record.ProjectWithoutAutofieldEntry,
	) (string, *record.ProjectEntry, *Error)
	UpdateProject(
		userId string,
		projectId string,
		entry record.ProjectWithoutAutofieldEntry,
	) (*record.ProjectEntry, *Error)
	DeleteProject(
		userId string,
		projectId string,
	) *Error
}

type projectRepository struct {
	client firestore.Client
}

func NewProjectRepository(client firestore.Client) ProjectRepository {
	return projectRepository{client: client}
}

func (r projectRepository) FetchProjects(
	userId string,
) (map[string]record.ProjectEntry, *Error) {
	iter := r.client.Collection(ProjectCollection).
		Where("userId", "==", userId).
		Documents(db.FirestoreContext())

	entries := make(map[string]record.ProjectEntry)

	for {
		snapshot, err := iter.Next()
		if err != nil {
			break
		}

		var values document.ProjectValues
		err = snapshot.DataTo(&values)
		if err != nil {
			return nil, Errorf(ReadFailurePanic, "failed to convert snapshot to values: %w", err)
		}

		entries[snapshot.Ref.ID] = *r.valuesToEntry(values)
	}

	return entries, nil
}

func (r projectRepository) FetchProject(
	userId string,
	projectId string,
) (*record.ProjectEntry, *Error) {
	snapshot, err := r.client.Collection(ProjectCollection).
		Doc(projectId).
		Get(db.FirestoreContext())
	if err != nil {
		return nil, Errorf(NotFoundError, "failed to fetch project")
	}

	var values document.ProjectValues
	err = snapshot.DataTo(&values)
	if err != nil {
		return nil, Errorf(ReadFailurePanic, "failed to convert snapshot to values: %w", err)
	}

	if values.UserId != userId {
		return nil, Errorf(NotFoundError, "failed to fetch project")
	}

	return r.valuesToEntry(values), nil
}

func (r projectRepository) InsertProject(
	userId string,
	entry record.ProjectWithoutAutofieldEntry,
) (string, *record.ProjectEntry, *Error) {
	ref, _, err := r.client.Collection(ProjectCollection).
		Add(db.FirestoreContext(), map[string]any{
			"name":        entry.Name,
			"description": entry.Description,
			"userId":      userId,
			"createdAt":   firestore.ServerTimestamp,
			"updatedAt":   firestore.ServerTimestamp,
		})
	if err != nil {
		return "", nil, Errorf(WriteFailurePanic, "failed to insert project: %w", err)
	}

	snapshot, err := ref.Get(db.FirestoreContext())
	if err != nil {
		return "", nil, Errorf(ReadFailurePanic, "failed to fetch inserted project: %w", err)
	}

	var values document.ProjectValues
	err = snapshot.DataTo(&values)
	if err != nil {
		return "", nil, Errorf(ReadFailurePanic, "failed to convert snapshot to values: %w", err)
	}

	return ref.ID, r.valuesToEntry(values), nil
}

func (r projectRepository) UpdateProject(
	userId string,
	projectId string,
	entry record.ProjectWithoutAutofieldEntry,
) (*record.ProjectEntry, *Error) {
	ref := r.client.Collection(ProjectCollection).
		Doc(projectId)

	snapshotToBeUpdated, err := ref.Get(db.FirestoreContext())
	if err != nil {
		return nil, Errorf(NotFoundError, "failed to update project")
	}

	var valuesToBeUpdated document.ProjectValues
	err = snapshotToBeUpdated.DataTo(&valuesToBeUpdated)
	if err != nil {
		return nil, Errorf(ReadFailurePanic, "failed to convert snapshot to values: %w", err)
	}

	if valuesToBeUpdated.UserId != userId {
		return nil, Errorf(NotFoundError, "failed to update project")
	}

	_, err = ref.Update(db.FirestoreContext(), []firestore.Update{
		{Path: "name", Value: entry.Name},
		{Path: "description", Value: entry.Description},
		{Path: "updatedAt", Value: firestore.ServerTimestamp},
	})
	if err != nil {
		return nil, Errorf(WriteFailurePanic, "failed to update project: %w", err)
	}

	snapshot, err := ref.Get(db.FirestoreContext())
	if err != nil {
		return nil, Errorf(ReadFailurePanic, "failed to fetch updated project: %w", err)
	}

	var values document.ProjectValues
	err = snapshot.DataTo(&values)
	if err != nil {
		return nil, Errorf(ReadFailurePanic, "failed to convert snapshot to values: %w", err)
	}

	return r.valuesToEntry(values), nil
}

func (r projectRepository) DeleteProject(
	userId string,
	projectId string,
) *Error {
	ref := r.client.Collection(ProjectCollection).
		Doc(projectId)

	snapshot, err := ref.Get(db.FirestoreContext())
	if err != nil {
		return Errorf(NotFoundError, "failed to delete project")
	}

	var values document.ProjectValues
	err = snapshot.DataTo(&values)
	if err != nil {
		return Errorf(ReadFailurePanic, "failed to convert snapshot to values: %w", err)
	}

	if values.UserId != userId {
		return Errorf(NotFoundError, "failed to delete project")
	}

	_, err = ref.Delete(db.FirestoreContext())
	if err != nil {
		return Errorf(WriteFailurePanic, "failed to delete project: %w", err)
	}

	return nil
}

func (r projectRepository) valuesToEntry(
	values document.ProjectValues,
) *record.ProjectEntry {
	return &record.ProjectEntry{
		Name:        values.Name,
		Description: values.Description,
		UserId:      values.UserId,
		CreatedAt:   values.CreatedAt,
		UpdatedAt:   values.UpdatedAt,
	}
}
