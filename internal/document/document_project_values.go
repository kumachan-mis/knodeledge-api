package document

import "time"

type ProjectValues struct {
	Name        string    `firestore:"name"`
	Description string    `firestore:"description,omitempty"`
	UserId      string    `firestore:"userId"`
	CreatedAt   time.Time `firestore:"createdAt"`
	UpdatedAt   time.Time `firestore:"updatedAt"`
}
