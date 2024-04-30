package document

import "time"

type ChapterValues struct {
	Name      string    `firestore:"name"`
	CreatedAt time.Time `firestore:"createdAt"`
	UpdatedAt time.Time `firestore:"updatedAt"`
}
