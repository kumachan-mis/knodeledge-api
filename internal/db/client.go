package db

import (
	"context"

	"cloud.google.com/go/firestore"
)

var firestoreClient *firestore.Client
var firestoreContext context.Context

func InitDatabaseClient(projectId string) error {
	firestoreContext = context.Background()
	client, err := firestore.NewClient(firestoreContext, projectId)
	if err != nil {
		return err
	}
	firestoreClient = client
	return nil
}

func FirestoreClient() *firestore.Client {
	return firestoreClient
}

func FirestoreContext() context.Context {
	return firestoreContext
}
