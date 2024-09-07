package domain

import "fmt"

type GraphNameObject struct {
	value string
}

func NewGraphNameObject(name string) (*GraphNameObject, error) {
	if name == "" {
		return nil, fmt.Errorf("graph name is required, but got '%v'", name)
	}
	if len(name) > 100 {
		return nil, fmt.Errorf("graph name cannot be longer than 100 characters, but got '%v'", name)
	}
	return &GraphNameObject{value: name}, nil
}

func (o *GraphNameObject) Value() string {
	return o.value
}
