package document

import "time"

type ChapterValues struct {
	Name      string    `firestore:"name"`
	NextId    string    `firestore:"nextId"`
	CreatedAt time.Time `firestore:"createdAt"`
	UpdatedAt time.Time `firestore:"updatedAt"`
}
