package domain

type GraphContentEntity struct {
	paragraph GraphParagraphObject
}

func NewGraphContentEntity(
	paragraph GraphParagraphObject,
) *GraphContentEntity {
	return &GraphContentEntity{
		paragraph: paragraph,
	}
}

func (e *GraphContentEntity) Paragraph() *GraphParagraphObject {
	return &e.paragraph
}
