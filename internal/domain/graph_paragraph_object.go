package domain

import "fmt"

type GraphParagraphObject struct {
	value string
}

func NewGraphParagraphObject(paragraph string) (*GraphParagraphObject, error) {
	paragraphLen := len(paragraph)
	if paragraphLen > 40000 {
		return nil, fmt.Errorf("graph paragraph must be less than or equal to 40000 bytes, but got %v bytes", paragraphLen)
	}
	return &GraphParagraphObject{value: paragraph}, nil
}

func (o GraphParagraphObject) Value() string {
	return o.value
}
