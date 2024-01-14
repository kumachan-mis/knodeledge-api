package domain

import "fmt"

type MessageObject struct {
	value string
}

func NewMessageObject(message string) (*MessageObject, error) {
	if message == "" {
		return nil, fmt.Errorf("message is empty")
	}

	return &MessageObject{value: message}, nil
}

func (o *MessageObject) Value() string {
	return o.value
}
