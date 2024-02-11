package domain

import "fmt"

type ProjectDescriptionObject struct {
	value string
}

func NewProjectDescriptionObject(description string) (*ProjectDescriptionObject, error) {
	if len(description) > 400 {
		return nil, fmt.Errorf("project description cannot be longer than 400 characters")
	}
	return &ProjectDescriptionObject{value: description}, nil
}

func (o *ProjectDescriptionObject) Value() string {
	return o.value
}
