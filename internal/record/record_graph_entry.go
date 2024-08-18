package record

import "time"

type GraphEntry struct {
	Paragraph string
	UserId    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
