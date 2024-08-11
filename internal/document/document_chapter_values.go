package document

import "time"

type ChapterValues struct {
	Name      string          `firestore:"name"`
	Sections  []SectionValues `firestore:"sections"`
	CreatedAt time.Time       `firestore:"createdAt"`
	UpdatedAt time.Time       `firestore:"updatedAt"`
}

type SectionValues struct {
	Id   string `firestore:"id"`
	Name string `firestore:"name"`
}
