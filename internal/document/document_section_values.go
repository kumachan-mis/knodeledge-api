package document

import "time"

type SectionValues struct {
	Name      string    `firestore:"name"`
	CreatedAt time.Time `firestore:"createdAt"`
	UpdatedAt time.Time `firestore:"updatedAt"`
}
