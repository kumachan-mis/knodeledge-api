package document

import "time"

type ProjectWithChapterIdsValues struct {
	Name        string    `firestore:"name"`
	Description string    `firestore:"description,omitempty"`
	ChapterIds  []string  `firestore:"chapterIds,omitempty"`
	UserId      string    `firestore:"userId"`
	CreatedAt   time.Time `firestore:"createdAt"`
	UpdatedAt   time.Time `firestore:"updatedAt"`
}
