package domain

import "fmt"

type ChapterIdObject struct {
	value string
}

func NewChapterIdObject(chapterId string) (*ChapterIdObject, error) {
	if chapterId == "" {
		return nil, fmt.Errorf("chapter id is required, but got '%s'", chapterId)
	}
	return &ChapterIdObject{value: chapterId}, nil
}

func (o *ChapterIdObject) Value() string {
	return o.value
}
