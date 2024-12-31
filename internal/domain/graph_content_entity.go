package domain

type GraphContentEntity struct {
	paragraph GraphParagraphObject
	children  GraphChildrenEntity
}

func NewGraphContentEntity(
	paragraph GraphParagraphObject,
	children GraphChildrenEntity,
) *GraphContentEntity {
	return &GraphContentEntity{
		paragraph: paragraph,
		children:  children,
	}
}

func (e *GraphContentEntity) Paragraph() *GraphParagraphObject {
	return &e.paragraph
}

func (e *GraphContentEntity) Children() *GraphChildrenEntity {
	return &e.children
}
