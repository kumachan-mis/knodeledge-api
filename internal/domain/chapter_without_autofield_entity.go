package domain

type ChapterWithoutAutofieldEntity struct {
	name   ChapterNameObject
	number ChapterNumberObject
}

func NewChapterWithoutAutofieldEntity(
	name ChapterNameObject,
	number ChapterNumberObject,

) *ChapterWithoutAutofieldEntity {
	return &ChapterWithoutAutofieldEntity{
		name:   name,
		number: number,
	}
}

func (e *ChapterWithoutAutofieldEntity) Name() *ChapterNameObject {
	return &e.name
}

func (e *ChapterWithoutAutofieldEntity) Number() *ChapterNumberObject {
	return &e.number
}
