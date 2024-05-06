package domain

import "fmt"

type PaperContentObject struct {
	value string
}

func NewPaperContentObject(content string) (*PaperContentObject, error) {
	contentLen := len(content)
	if contentLen > 40000 {
		return nil, fmt.Errorf("paper content must be less than or equal to 40000 bytes, but got %v bytes", contentLen)
	}
	return &PaperContentObject{value: content}, nil
}
