package repository

import (
	"cloud.google.com/go/firestore"
	"github.com/kumachan-mis/knodeledge-api/internal/db"
	"github.com/kumachan-mis/knodeledge-api/internal/record"
)

//go:generate mockgen -source=$GOFILE -destination=../../mock/$GOPACKAGE/mock_$GOFILE -package=$GOPACKAGE

const ProjectCollection = "projects"

type ProjectRepository interface {
	FetchUserProjects(userId string) (map[string]record.ProjectEntry, *Error)
	FetchProject(projectId string) (*record.ProjectEntry, *Error)
	InsertProject(entry record.ProjectWithoutAutofieldEntry) (string, *record.ProjectEntry, *Error)
	UpdateProject(projectId string, entry record.ProjectWithoutAutofieldEntry) (*record.ProjectEntry, *Error)
}

type projectRepository struct {
	client firestore.Client
}

func NewProjectRepository(client firestore.Client) ProjectRepository {
	return projectRepository{client: client}
}

func (r projectRepository) FetchUserProjects(userId string) (map[string]record.ProjectEntry, *Error) {
	iter := r.client.Collection(ProjectCollection).
		Where("userId", "==", userId).
		Documents(db.FirestoreContext())

	entries := make(map[string]record.ProjectEntry)

	for {
		snapshot, err := iter.Next()
		if err != nil {
			break
		}

		var entry record.ProjectEntry
		err = snapshot.DataTo(&entry)
		if err != nil {
			return nil, Errorf(ReadFailurePanic, "failed to convert snapshot to entry: %w", err)
		}

		entries[snapshot.Ref.ID] = entry
	}

	return entries, nil
}

func (r projectRepository) FetchProject(projectId string) (*record.ProjectEntry, *Error) {
	snapshot, err := r.client.Collection(ProjectCollection).
		Doc(projectId).
		Get(db.FirestoreContext())
	if err != nil {
		return nil, Errorf(NotFoundError, "failed to get project")
	}

	var entry record.ProjectEntry
	err = snapshot.DataTo(&entry)
	if err != nil {
		return nil, Errorf(ReadFailurePanic, "failed to convert snapshot to entry: %w", err)
	}

	return &entry, nil
}

func (r projectRepository) InsertProject(entry record.ProjectWithoutAutofieldEntry) (string, *record.ProjectEntry, *Error) {
	ref, _, err := r.client.Collection(ProjectCollection).
		Add(db.FirestoreContext(), map[string]any{
			"name":        entry.Name,
			"description": entry.Description,
			"userId":      entry.UserId,
			"createdAt":   firestore.ServerTimestamp,
			"updatedAt":   firestore.ServerTimestamp,
		})
	if err != nil {
		return "", nil, Errorf(WriteFailurePanic, "failed to insert project")
	}

	snapshot, err := ref.Get(db.FirestoreContext())
	if err != nil {
		return "", nil, Errorf(ReadFailurePanic, "failed to get inserted project")
	}

	var entryWithAutofield record.ProjectEntry
	err = snapshot.DataTo(&entryWithAutofield)
	if err != nil {
		return "", nil, Errorf(ReadFailurePanic, "failed to convert snapshot to entry: %w", err)
	}

	return ref.ID, &entryWithAutofield, nil
}

func (r projectRepository) UpdateProject(projectId string, entry record.ProjectWithoutAutofieldEntry) (*record.ProjectEntry, *Error) {
	_, err := r.client.Collection(ProjectCollection).
		Doc(projectId).
		Update(db.FirestoreContext(), []firestore.Update{
			{Path: "name", Value: entry.Name},
			{Path: "description", Value: entry.Description},
			{Path: "userId", Value: entry.UserId},
			{Path: "updatedAt", Value: firestore.ServerTimestamp},
		})
	if err != nil {
		return nil, Errorf(WriteFailurePanic, "failed to update project")
	}

	snapshot, err := r.client.Collection(ProjectCollection).
		Doc(projectId).
		Get(db.FirestoreContext())
	if err != nil {
		return nil, Errorf(ReadFailurePanic, "failed to get updated project")
	}

	var entryWithAutofield record.ProjectEntry
	err = snapshot.DataTo(&entryWithAutofield)
	if err != nil {
		return nil, Errorf(ReadFailurePanic, "failed to convert snapshot to entry: %w", err)
	}

	return &entryWithAutofield, nil
}
