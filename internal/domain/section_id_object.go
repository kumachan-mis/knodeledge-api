package domain

import "fmt"

type SectionIdObject struct {
	value string
}

func NewSectionIdObject(sectionId string) (*SectionIdObject, error) {
	if sectionId == "" {
		return nil, fmt.Errorf("section id is required, but got '%v'", sectionId)
	}
	return &SectionIdObject{value: sectionId}, nil
}

func (o *SectionIdObject) Value() string {
	return o.value
}

func (o *SectionIdObject) Len() int {
	return len(o.value)
}
