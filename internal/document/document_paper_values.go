package document

import "time"

type PaperValues struct {
	Content   string    `firestore:"content"`
	CreatedAt time.Time `firestore:"createdAt"`
	UpdatedAt time.Time `firestore:"updatedAt"`
}
