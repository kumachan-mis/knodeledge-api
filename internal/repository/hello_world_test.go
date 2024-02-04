package repository_test

import (
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/kumachan-mis/knodeledge-api/internal/db"
	"github.com/kumachan-mis/knodeledge-api/internal/record"
	"github.com/kumachan-mis/knodeledge-api/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	db.InitDatabaseClient(firestore.DetectProjectID)
	m.Run()
	db.FinalizeDatabaseClient()
}

func TestFetchHelloWorldValidDocument(t *testing.T) {
	client := db.FirestoreClient()
	r := repository.NewHelloWorldRepository(*client)

	id, entity, err := r.FetchHelloWorld("FetchHelloWorld Test")

	assert.Equal(t, "FETCH_HELLO_WORLD_TEST_DOC", id)
	assert.Equal(t, "FetchHelloWorld Test", entity.Name)
	assert.Equal(t, "Hello, FetchHelloWorld Test!", entity.Message)
	assert.NoError(t, err)
}

func TestFetchHelloWorldInvalidDocumment(t *testing.T) {
	client := db.FirestoreClient()
	r := repository.NewHelloWorldRepository(*client)

	id, entity, err := r.FetchHelloWorld("FetchHelloWorld InvalidDocumment")

	assert.Empty(t, id)
	assert.Nil(t, entity)
	assert.Error(t, err)
}

func TestFetchHelloWorldNotFound(t *testing.T) {
	client := db.FirestoreClient()
	r := repository.NewHelloWorldRepository(*client)

	id, entity, err := r.FetchHelloWorld("FetchHelloWorld NotFoundTest")

	assert.Empty(t, id)
	assert.Nil(t, entity)
	assert.Error(t, err)
}

func TestCreateHelloWorldSucceeded(t *testing.T) {
	client := db.FirestoreClient()
	r := repository.NewHelloWorldRepository(*client)

	id, err := r.CreateHelloWorld(record.HelloWorldEntry{
		Name:    "CreateHelloWorld Test",
		Message: "Hello, CreateHelloWorld Test!",
	})

	assert.NotEmpty(t, id)
	assert.NoError(t, err)

	doc, err := client.Collection(repository.HELLO_WORLD_COLLECTION).
		Doc(id).
		Get(db.FirestoreContext())

	assert.NoError(t, err)

	var entity record.HelloWorldEntry
	err = doc.DataTo(&entity)

	assert.NoError(t, err)
	assert.Equal(t, "CreateHelloWorld Test", entity.Name)
	assert.Equal(t, "Hello, CreateHelloWorld Test!", entity.Message)
}
