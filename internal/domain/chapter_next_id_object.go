package domain

type ChapterNextIdObject struct {
	value string
}

func NewChapterNextIdObject(chapterId string) (*ChapterNextIdObject, error) {
	return &ChapterNextIdObject{value: chapterId}, nil
}

func (o *ChapterNextIdObject) Value() string {
	return o.value
}
