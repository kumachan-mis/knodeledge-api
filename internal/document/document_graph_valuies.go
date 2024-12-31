package document

import "time"

type GraphValues struct {
	Paragraph string             `firestore:"paragraph"`
	Children  []GraphChildValues `firestore:"children"`
	CreatedAt time.Time          `firestore:"createdAt"`
	UpdatedAt time.Time          `firestore:"updatedAt"`
}

type GraphChildValues struct {
	Name       string             `firestore:"name"`
	Relation   string             `firestore:"relation"`
	Descrition string             `firestore:"description"`
	Children   []GraphChildValues `firestore:"children"`
}
