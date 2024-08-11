package domain

type SectionEntity struct {
	id        SectionIdObject
	name      SectionNameObject
	createdAt CreatedAtObject
	updatedAt UpdatedAtObject
}

func NewSectionEntity(
	id SectionIdObject,
	name SectionNameObject,
	createdAt CreatedAtObject,
	updatedAt UpdatedAtObject,
) *SectionEntity {
	return &SectionEntity{
		id:        id,
		name:      name,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (e *SectionEntity) Id() *SectionIdObject {
	return &e.id
}

func (e *SectionEntity) Name() *SectionNameObject {
	return &e.name
}

func (e *SectionEntity) CreatedAt() *CreatedAtObject {
	return &e.createdAt
}

func (e *SectionEntity) UpdatedAt() *UpdatedAtObject {
	return &e.updatedAt
}
