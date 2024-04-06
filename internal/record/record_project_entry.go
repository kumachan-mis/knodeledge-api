package record

import "time"

type ProjectEntry struct {
	Name        string
	Description string
	UserId      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
