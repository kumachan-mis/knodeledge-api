package domain

import "fmt"

type PaperIdObject struct {
	value string
}

func NewPaperIdObject(paperId string) (*PaperIdObject, error) {
	if paperId == "" {
		return nil, fmt.Errorf("paper id is required, but got '%v'", paperId)
	}
	return &PaperIdObject{value: paperId}, nil
}

func (o *PaperIdObject) Value() string {
	return o.value
}
