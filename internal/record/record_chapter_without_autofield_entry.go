package record

type ChapterWithoutAutofieldEntry struct {
	Name     string
	Number   int
	Sections []SectionWithoutAutofieldEntry
	UserId   string
}

type SectionWithoutAutofieldEntry struct {
	Name string
}
