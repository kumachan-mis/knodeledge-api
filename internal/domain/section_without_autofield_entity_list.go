package domain

import "fmt"

type SectionWithoutAutofieldEntityList struct {
	sections []SectionWithoutAutofieldEntity
	len      int
}

func NewSectionWithoutAutofieldEntityList(
	sections []SectionWithoutAutofieldEntity,
) (*SectionWithoutAutofieldEntityList, error) {
	if len(sections) == 0 {
		return nil, fmt.Errorf("sections are required, but got []")
	}
	if len(sections) > 20 {
		return nil, fmt.Errorf("sections length must be less than or equal to 20, but got %v", len(sections))
	}
	return &SectionWithoutAutofieldEntityList{sections: sections, len: len(sections)}, nil
}

func (o *SectionWithoutAutofieldEntityList) Value() []SectionWithoutAutofieldEntity {
	return o.sections
}

func (o *SectionWithoutAutofieldEntityList) Len() int {
	return o.len
}
