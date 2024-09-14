package record

import "time"

type GraphContentEntry struct {
	Paragraph string
	UserId    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
