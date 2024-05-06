package domain

type PaperEntity struct {
	id        PaperIdObject
	content   PaperContentObject
	createdAt CreatedAtObject
	updatedAt UpdatedAtObject
}

func NewPaperEntity(
	id PaperIdObject,
	content PaperContentObject,
	createdAt CreatedAtObject,
	updatedAt UpdatedAtObject,
) *PaperEntity {
	return &PaperEntity{
		id:        id,
		content:   content,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (e *PaperEntity) Id() *PaperIdObject {
	return &e.id
}

func (e *PaperEntity) Content() *PaperContentObject {
	return &e.content
}

func (e *PaperEntity) CreatedAt() *CreatedAtObject {
	return &e.createdAt
}

func (e *PaperEntity) UpdatedAt() *UpdatedAtObject {
	return &e.updatedAt
}
