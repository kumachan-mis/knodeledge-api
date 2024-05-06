package domain

type PaperWithoutAutofieldEntity struct {
	content PaperContentObject
}

func NewPaperWithoutAutofieldEntity(content PaperContentObject) *PaperWithoutAutofieldEntity {
	return &PaperWithoutAutofieldEntity{content: content}
}

func (e *PaperWithoutAutofieldEntity) Content() *PaperContentObject {
	return &e.content
}
