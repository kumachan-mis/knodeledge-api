package domain

type SectionWithoutAutofieldEntity struct {
	id   SectionIdObject
	name SectionNameObject
}

func NewSectionWithoutAutofieldEntity(
	id SectionIdObject,
	name SectionNameObject,
) *SectionWithoutAutofieldEntity {
	return &SectionWithoutAutofieldEntity{
		id:   id,
		name: name,
	}
}

func (e *SectionWithoutAutofieldEntity) Id() *SectionIdObject {
	return &e.id
}

func (e *SectionWithoutAutofieldEntity) Name() *SectionNameObject {
	return &e.name
}
