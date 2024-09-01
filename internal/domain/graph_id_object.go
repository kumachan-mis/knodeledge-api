package domain

import "fmt"

type GraphIdObject struct {
	value string
}

func NewGraphIdObject(graphId string) (*GraphIdObject, error) {
	if graphId == "" {
		return nil, fmt.Errorf("graph id is required, but got '%v'", graphId)
	}
	return &GraphIdObject{value: graphId}, nil
}

func (o *GraphIdObject) Value() string {
	return o.value
}
