package record

import "time"

type SectionEntry struct {
	Name      string
	Number    int
	UserId    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
