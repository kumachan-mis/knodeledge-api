package domain

import "fmt"

type GraphChildrenEntity struct {
	children []GraphChildEntity
	len      int
}

func NewGraphChildrenEntity(children []GraphChildEntity) (*GraphChildrenEntity, error) {
	names := make(map[string]struct{}, len(children)+1)
	for _, child := range children {
		if _, ok := names[child.Name().Value()]; ok {
			return nil, fmt.Errorf("names of children must be unique, but got '%v' duplicated", child.Name().Value())
		}
		names[child.Name().Value()] = struct{}{}
	}

	return &GraphChildrenEntity{children: children, len: len(children)}, nil
}

func (e *GraphChildrenEntity) Value() []GraphChildEntity {
	return e.children
}

func (e *GraphChildrenEntity) Len() int {
	return e.len
}
