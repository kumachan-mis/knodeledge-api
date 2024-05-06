package domain

import "fmt"

type PaperIdObject struct {
	value string
}

func NewPaperIdObject(chapterId string) (*PaperIdObject, error) {
	if chapterId == "" {
		return nil, fmt.Errorf("paper id is required, but got '%v'", chapterId)
	}
	return &PaperIdObject{value: chapterId}, nil
}

func (o *PaperIdObject) Value() string {
	return o.value
}
