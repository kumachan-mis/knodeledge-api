package domain

type ChapterWithoutAutofieldEntity struct {
	name     ChapterNameObject
	number   ChapterNumberObject
	sections []SectionWithoutAutofieldEntity
}

func NewChapterWithoutAutofieldEntity(
	name ChapterNameObject,
	number ChapterNumberObject,
	sections []SectionWithoutAutofieldEntity,
) *ChapterWithoutAutofieldEntity {
	return &ChapterWithoutAutofieldEntity{
		name:     name,
		number:   number,
		sections: sections,
	}
}

func (e *ChapterWithoutAutofieldEntity) Name() *ChapterNameObject {
	return &e.name
}

func (e *ChapterWithoutAutofieldEntity) Number() *ChapterNumberObject {
	return &e.number
}

func (e *ChapterWithoutAutofieldEntity) Sections() []SectionWithoutAutofieldEntity {
	return e.sections
}
