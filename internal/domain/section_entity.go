package domain

type SectionObject struct {
	id        SectionIdObject
	name      SectionNameObject
	createdAt CreatedAtObject
	updatedAt UpdatedAtObject
}

func NewSectionObject(
	id SectionIdObject,
	name SectionNameObject,
	createdAt CreatedAtObject,
	updatedAt UpdatedAtObject,
) *SectionObject {
	return &SectionObject{
		id:        id,
		name:      name,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (e *SectionObject) Id() *SectionIdObject {
	return &e.id
}

func (e *SectionObject) Name() *SectionNameObject {
	return &e.name
}

func (e *SectionObject) CreatedAt() *CreatedAtObject {
	return &e.createdAt
}

func (e *SectionObject) UpdatedAt() *UpdatedAtObject {
	return &e.updatedAt
}
