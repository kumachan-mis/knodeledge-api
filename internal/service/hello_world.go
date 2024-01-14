package service

import (
	"github.com/kumachan-mis/knodeledge-api/internal/domain"
	"github.com/kumachan-mis/knodeledge-api/internal/record"
	"github.com/kumachan-mis/knodeledge-api/internal/repository"
)

func SearchHelloWorld(name domain.NameObject) (*domain.MessageObject, error) {
	_, entry, err := repository.FetchHelloWorld(name.Value())
	if err != nil {
		return nil, err
	}
	return domain.NewMessageObject(entry.Message)
}

func LogHelloWorld(name domain.NameObject) (*domain.MessageObject, error) {
	entity := domain.NewHelloWorldEntity(name)
	message, err := entity.Message()
	if err != nil {
		return nil, err
	}

	entry := record.HelloWorldEntry{Name: name.Value(), Message: message.Value()}
	_, err = repository.CreateHelloWorld(entry)
	if err != nil {
		return nil, err
	}

	return message, nil
}
