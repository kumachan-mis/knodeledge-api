package record

import "time"

type ChapterEntry struct {
	Name      string    `firestore:"name"`
	Number    int       `firestore:"number"`
	CreatedAt time.Time `firestore:"createdAt"`
	UpdatedAt time.Time `firestore:"updatedAt"`
}
