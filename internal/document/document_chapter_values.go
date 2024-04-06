package document

import "time"

type ChapterValues struct {
	Name      string    `firestore:"name"`
	Number    int       `firestore:"number"`
	CreatedAt time.Time `firestore:"createdAt"`
	UpdatedAt time.Time `firestore:"updatedAt"`
}
