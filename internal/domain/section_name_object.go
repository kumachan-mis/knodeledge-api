package domain

import "fmt"

type SectionNameObject struct {
	value string
}

func NewSectionNameObject(name string) (*SectionNameObject, error) {
	if name == "" {
		return nil, fmt.Errorf("section name is required, but got '%v'", name)
	}
	if len(name) > 100 {
		return nil, fmt.Errorf("section name cannot be longer than 100 characters, but got '%v'", name)
	}
	return &SectionNameObject{value: name}, nil
}

func (o *SectionNameObject) Value() string {
	return o.value
}
