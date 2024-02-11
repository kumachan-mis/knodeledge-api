package domain

import "fmt"

type ProjectIdObject struct {
	value string
}

func NewProjectIdObject(projectId string) (*ProjectIdObject, error) {
	if projectId == "" {
		return nil, fmt.Errorf("project id is required")
	}
	return &ProjectIdObject{value: projectId}, nil
}

func (o *ProjectIdObject) Value() string {
	return o.value
}
