package domain

import "fmt"

type SectionContentObject struct {
	value string
}

func NewSectionContentObject(content string) (*SectionContentObject, error) {
	contentLen := len(content)
	if contentLen > 40000 {
		return nil, fmt.Errorf("section content must be less than or equal to 40000 bytes, but got %v bytes", contentLen)
	}
	return &SectionContentObject{value: content}, nil
}

func (o SectionContentObject) Value() string {
	return o.value
}
