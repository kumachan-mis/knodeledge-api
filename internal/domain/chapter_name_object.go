package domain

import "fmt"

type ChapterNameObject struct {
	value string
}

func NewChapterNameObject(name string) (*ChapterNameObject, error) {
	if name == "" {
		return nil, fmt.Errorf("chapter name is required, but got '%v'", name)
	}
	if len(name) > 100 {
		return nil, fmt.Errorf("chapter name cannot be longer than 100 characters, but got '%v'", name)
	}
	return &ChapterNameObject{value: name}, nil
}

func (o *ChapterNameObject) Value() string {
	return o.value
}
