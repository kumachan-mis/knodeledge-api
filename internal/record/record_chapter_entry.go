package record

import "time"

type ChapterEntry struct {
	Name      string    `firestore:"name"`
	Number    int       `firestore:"number"`
	UserId    string    `firestore:"userId"`
	CreatedAt time.Time `firestore:"createdAt"`
	UpdatedAt time.Time `firestore:"updatedAt"`
}
