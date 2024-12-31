package record

import "time"

type GraphEntry struct {
	Name      string
	Paragraph string
	Children  []GraphChildEntry
	UserId    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
