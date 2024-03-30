package domain

import "fmt"

type ChapterNumberObject struct {
	value int
}

func NewChapterNumberObject(number int) (*ChapterNumberObject, error) {
	if number < 1 {
		return nil, fmt.Errorf("chapter number must be greater than 0, but got '%d'", number)
	}
	return &ChapterNumberObject{value: number}, nil
}

func (o *ChapterNumberObject) Value() int {
	return o.value
}
