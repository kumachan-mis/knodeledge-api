package repository

import (
	"cloud.google.com/go/firestore"
	"github.com/kumachan-mis/knodeledge-api/internal/db"
	"github.com/kumachan-mis/knodeledge-api/internal/record"
)

//go:generate mockgen -source=$GOFILE -destination=../../mock/$GOPACKAGE/mock_$GOFILE -package=$GOPACKAGE

const HELLO_WORLD_COLLECTION = "hello_world"

type HelloWorldRepository interface {
	FetchHelloWorld(name string) (string, *record.HelloWorldEntry, error)
	CreateHelloWorld(entry record.HelloWorldEntry) (string, error)
}

type helloWorldRepository struct {
	client firestore.Client
}

func NewHelloWorldRepository(client firestore.Client) HelloWorldRepository {
	return helloWorldRepository{client: client}
}

func (r helloWorldRepository) FetchHelloWorld(name string) (string, *record.HelloWorldEntry, error) {
	iter := r.client.Collection(HELLO_WORLD_COLLECTION).
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

func (r helloWorldRepository) CreateHelloWorld(entry record.HelloWorldEntry) (string, error) {
	ref, _, err := r.client.
		Collection(HELLO_WORLD_COLLECTION).
		Add(db.FirestoreContext(), entry)
	if err != nil {
		return "", err
	}
	return ref.ID, nil
}
