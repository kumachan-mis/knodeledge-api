package domain

type GraphEntity struct {
	id        GraphIdObject
	name      GraphNameObject
	paragraph GraphParagraphObject
	createdAt CreatedAtObject
	updatedAt UpdatedAtObject
}

func NewGraphEntity(
	id GraphIdObject,
	name GraphNameObject,
	paragraph GraphParagraphObject,
	createdAt CreatedAtObject,
	updatedAt UpdatedAtObject,
) *GraphEntity {
	return &GraphEntity{
		id:        id,
		name:      name,
		paragraph: paragraph,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (e *GraphEntity) Id() *GraphIdObject {
	return &e.id
}

func (e *GraphEntity) Name() *GraphNameObject {
	return &e.name
}

func (e *GraphEntity) Paragraph() *GraphParagraphObject {
	return &e.paragraph
}

func (e *GraphEntity) CreatedAt() *CreatedAtObject {
	return &e.createdAt
}

func (e *GraphEntity) UpdatedAt() *UpdatedAtObject {
	return &e.updatedAt
}
