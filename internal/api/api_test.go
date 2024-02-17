package api_test

import (
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/kumachan-mis/knodeledge-api/internal/db"
)

func TestMain(m *testing.M) {
	db.InitDatabaseClient(firestore.DetectProjectID)
	m.Run()
	db.FinalizeDatabaseClient()
}
