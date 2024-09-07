package domain

type SectionWithoutAutofieldEntity struct {
	name    SectionNameObject
	content SectionContentObject
}

func NewSectionWithoutAutofieldEntity(
	name SectionNameObject,
	content SectionContentObject,
) *SectionWithoutAutofieldEntity {
	return &SectionWithoutAutofieldEntity{
		name:    name,
		content: content,
	}
}

func (e *SectionWithoutAutofieldEntity) Name() *SectionNameObject {
	return &e.name
}

func (e *SectionWithoutAutofieldEntity) Content() *SectionContentObject {
	return &e.content
}
