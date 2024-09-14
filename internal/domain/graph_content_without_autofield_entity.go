package domain

type GraphContentWithoutAutofieldEntity struct {
	paragraph GraphParagraphObject
}

func NewGraphContentWithoutAutofieldEntity(
	paragraph GraphParagraphObject,
) *GraphContentWithoutAutofieldEntity {
	return &GraphContentWithoutAutofieldEntity{
		paragraph: paragraph,
	}
}

func (e *GraphContentWithoutAutofieldEntity) Paragraph() *GraphParagraphObject {
	return &e.paragraph
}
