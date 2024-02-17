package domain

import "fmt"

type ProjectNameObject struct {
	value string
}

func NewProjectNameObject(name string) (*ProjectNameObject, error) {
	if name == "" {
		return nil, fmt.Errorf("project name is required, but got '%s'", name)
	}
	if len(name) > 100 {
		return nil, fmt.Errorf("project name cannot be longer than 100 characters, but got '%s'", name)
	}
	return &ProjectNameObject{value: name}, nil
}

func (o *ProjectNameObject) Value() string {
	return o.value
}
