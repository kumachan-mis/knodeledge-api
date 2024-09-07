package domain

type SectionOfChapterWithoutAutofieldEntity struct {
	id   SectionIdObject
	name SectionNameObject
}

func NewSectionOfChapterWithoutAutofieldEntity(
	id SectionIdObject,
	name SectionNameObject,
) *SectionOfChapterWithoutAutofieldEntity {
	return &SectionOfChapterWithoutAutofieldEntity{
		id:   id,
		name: name,
	}
}

func (e *SectionOfChapterWithoutAutofieldEntity) Id() *SectionIdObject {
	return &e.id
}

func (e *SectionOfChapterWithoutAutofieldEntity) Name() *SectionNameObject {
	return &e.name
}
