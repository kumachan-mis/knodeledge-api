package domain

import "fmt"

type GraphRelationObject struct {
	value string
}

func NewGraphRelationObject(relation string) (*GraphRelationObject, error) {
	if len(relation) > 100 {
		return nil, fmt.Errorf("graph relation cannot be longer than 100 characters, but got '%v'", relation)
	}
	return &GraphRelationObject{value: relation}, nil
}

func (o *GraphRelationObject) Value() string {
	return o.value
}
