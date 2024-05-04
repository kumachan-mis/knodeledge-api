package record

import "time"

type PaperEntry struct {
	Content   string
	UserId    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
