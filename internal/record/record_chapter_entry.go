package record

import "time"

type ChapterEntry struct {
	Name      string
	Number    int
	Sections  []SectionEntry
	UserId    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SectionEntry struct {
	Id        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
