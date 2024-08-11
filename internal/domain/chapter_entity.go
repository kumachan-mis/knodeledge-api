package domain

type ChapterEntity struct {
	id        ChapterIdObject
	name      ChapterNameObject
	number    ChapterNumberObject
	sections  []SectionEntity
	createdAt CreatedAtObject
	updatedAt UpdatedAtObject
}

func NewChapterEntity(
	id ChapterIdObject,
	name ChapterNameObject,
	number ChapterNumberObject,
	sections []SectionEntity,
	createdAt CreatedAtObject,
	updatedAt UpdatedAtObject,

) *ChapterEntity {
	return &ChapterEntity{
		id:        id,
		name:      name,
		number:    number,
		sections:  sections,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (e *ChapterEntity) Id() *ChapterIdObject {
	return &e.id
}

func (e *ChapterEntity) Name() *ChapterNameObject {
	return &e.name
}

func (e *ChapterEntity) Number() *ChapterNumberObject {
	return &e.number
}

func (e *ChapterEntity) Sections() []SectionEntity {
	return e.sections
}

func (e *ChapterEntity) CreatedAt() *CreatedAtObject {
	return &e.createdAt
}

func (e *ChapterEntity) UpdatedAt() *UpdatedAtObject {
	return &e.updatedAt
}
