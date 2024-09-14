package domain

type GraphContentEntity struct {
	id        GraphIdObject
	paragraph GraphParagraphObject
	createdAt CreatedAtObject
	updatedAt UpdatedAtObject
}

func NewGraphContentEntity(
	id GraphIdObject,
	paragraph GraphParagraphObject,
	createdAt CreatedAtObject,
	updatedAt UpdatedAtObject,
) *GraphContentEntity {
	return &GraphContentEntity{
		id:        id,
		paragraph: paragraph,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (e *GraphContentEntity) Id() *GraphIdObject {
	return &e.id
}

func (e *GraphContentEntity) Paragraph() *GraphParagraphObject {
	return &e.paragraph
}

func (e *GraphContentEntity) CreatedAt() *CreatedAtObject {
	return &e.createdAt
}

func (e *GraphContentEntity) UpdatedAt() *UpdatedAtObject {
	return &e.updatedAt
}
