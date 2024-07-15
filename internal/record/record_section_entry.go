package record

import "time"

type SectionEntry struct {
	Id        string
	Name      string
	UserId    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
