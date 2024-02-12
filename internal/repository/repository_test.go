package repository_test

import (
	"fmt"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/kumachan-mis/knodeledge-api/internal/db"
)

func UserId() string {
	return "auth0|65a3d656ca600978b0f9501b"
}

func ErrorUserId(i int) string {
	return fmt.Sprintf("error|%024d", i)
}

func UnknownUserId() string {
	return "unknown|000000000000000000000000"
}

func TestMain(m *testing.M) {
	db.InitDatabaseClient(firestore.DetectProjectID)
	m.Run()
	db.FinalizeDatabaseClient()
}
