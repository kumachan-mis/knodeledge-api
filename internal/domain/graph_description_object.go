package domain

import "fmt"

type GraphDescriptionObject struct {
	value string
}

func NewGraphDescriptionObject(description string) (*GraphDescriptionObject, error) {
	if len(description) > 400 {
		return nil, fmt.Errorf("graph description cannot be longer than 400 characters, but got '%v'", description)
	}
	return &GraphDescriptionObject{value: description}, nil
}

func (o *GraphDescriptionObject) Value() string {
	return o.value
}
