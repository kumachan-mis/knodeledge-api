package record

import "time"

type ChapterEntry struct {
	Name      string
	Number    int
	UserId    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
