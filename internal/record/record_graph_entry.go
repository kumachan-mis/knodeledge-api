package record

import "time"

type GraphEntry struct {
	Name      string
	Paragraph string
	UserId    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
