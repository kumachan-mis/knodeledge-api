package db

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
)

var FirestoreContext context.Context
var FirestoreClient *firestore.Client

func InitDatabase() error {
	FirestoreContext = context.Background()
	client, err := firestore.NewClient(FirestoreContext, os.Getenv("FIREBASE_PROJECT_ID"))
	if err != nil {
		return err
	}
	FirestoreClient = client
	return nil
}
