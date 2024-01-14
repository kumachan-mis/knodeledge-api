package db

import (
	"context"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
)

var firestoreClient *firestore.Client
var firestoreContext context.Context

func InitDatabaseClient(projectID string) error {
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: projectID}

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		return err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return err
	}

	firestoreClient = client
	firestoreContext = ctx
	return nil
}

func FinalizeDatabaseClient() error {
	err := firestoreClient.Close()
	firestoreClient = nil
	firestoreContext = nil
	return err
}

func FirestoreClient() *firestore.Client {
	return firestoreClient
}

func FirestoreContext() context.Context {
	return firestoreContext
}
