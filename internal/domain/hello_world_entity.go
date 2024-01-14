package domain

import "fmt"

type HelloWorldEntity struct {
	name NameObject
}

func NewHelloWorldEntity(name NameObject) HelloWorldEntity {
	return HelloWorldEntity{name: name}
}

func (e *HelloWorldEntity) Name() NameObject {
	return e.name
}

func (e *HelloWorldEntity) Message() (*MessageObject, error) {
	if e.name.IsGuest() {
		return NewMessageObject("Hello World!")
	}

	return NewMessageObject(fmt.Sprintf("Hello, %s!", e.name.Value()))
}
