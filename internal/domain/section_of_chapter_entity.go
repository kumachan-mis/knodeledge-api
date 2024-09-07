package domain

type SectionOfChapterEntity struct {
	id        SectionIdObject
	name      SectionNameObject
	createdAt CreatedAtObject
	updatedAt UpdatedAtObject
}

func NewSectionOfChapterEntity(
	id SectionIdObject,
	name SectionNameObject,
	createdAt CreatedAtObject,
	updatedAt UpdatedAtObject,
) *SectionOfChapterEntity {
	return &SectionOfChapterEntity{
		id:        id,
		name:      name,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (e *SectionOfChapterEntity) Id() *SectionIdObject {
	return &e.id
}

func (e *SectionOfChapterEntity) Name() *SectionNameObject {
	return &e.name
}

func (e *SectionOfChapterEntity) CreatedAt() *CreatedAtObject {
	return &e.createdAt
}

func (e *SectionOfChapterEntity) UpdatedAt() *UpdatedAtObject {
	return &e.updatedAt
}
