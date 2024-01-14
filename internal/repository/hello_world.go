package repository

import (
	"github.com/kumachan-mis/knodeledge-api/internal/db"
	"github.com/kumachan-mis/knodeledge-api/internal/record"
)

const HELLO_WORLD_COLLECTION = "hello_world"

func FetchHelloWorld(name string) (string, *record.HelloWorldEntry, error) {
	iter := db.FirestoreClient().Collection(HELLO_WORLD_COLLECTION).
		Where("name", "==", name).
		Limit(1).
		Documents(db.FirestoreContext())

	snapshot, err := iter.Next()
	if err != nil {
		return "", nil, err
	}

	var entry record.HelloWorldEntry
	err = snapshot.DataTo(&entry)
	if err != nil {
		return "", nil, err
	}

	return snapshot.Ref.ID, &entry, nil
}

func CreateHelloWorld(entry record.HelloWorldEntry) (string, error) {
	ref, _, err := db.FirestoreClient().
		Collection(HELLO_WORLD_COLLECTION).
		Add(db.FirestoreContext(), entry)
	if err != nil {
		return "", err
	}
	return ref.ID, nil
}
