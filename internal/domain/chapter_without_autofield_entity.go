package domain

type ChapterWithoutAutofieldEntity struct {
	name   ChapterNameObject
	nextId ChapterNextIdObject
}

func NewChapterWithoutAutofieldEntity(
	name ChapterNameObject,
	nextId ChapterNextIdObject,

) *ChapterWithoutAutofieldEntity {
	return &ChapterWithoutAutofieldEntity{
		name:   name,
		nextId: nextId,
	}
}

func (e *ChapterWithoutAutofieldEntity) Name() *ChapterNameObject {
	return &e.name
}

func (e *ChapterWithoutAutofieldEntity) NextId() *ChapterNextIdObject {
	return &e.nextId
}
