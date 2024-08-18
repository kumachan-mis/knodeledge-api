package document

import "time"

type GraphValues struct {
	Paragraph string    `firestore:"paragraph"`
	CreatedAt time.Time `firestore:"createdAt"`
	UpdatedAt time.Time `firestore:"updatedAt"`
}
