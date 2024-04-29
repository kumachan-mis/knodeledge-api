package record

import "time"

type ChapterEntry struct {
	Name      string
	NextId    string
	UserId    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
