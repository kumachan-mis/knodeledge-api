package service

import (
	"github.com/kumachan-mis/knodeledge-api/internal/domain"
	"github.com/kumachan-mis/knodeledge-api/internal/record"
	"github.com/kumachan-mis/knodeledge-api/internal/repository"
)

type HelloWorldService interface {
	SearchHelloWorld(name domain.NameObject) (*domain.MessageObject, error)
	LogHelloWorld(name domain.NameObject) (*domain.MessageObject, error)
}

type helloWorldService struct {
	repository repository.HelloWorldRepository
}

func NewHelloWorldService(repository repository.HelloWorldRepository) HelloWorldService {
	return helloWorldService{repository: repository}
}

func (s helloWorldService) SearchHelloWorld(name domain.NameObject) (*domain.MessageObject, error) {
	_, entry, err := s.repository.FetchHelloWorld(name.Value())
	if err != nil {
		return nil, err
	}
	return domain.NewMessageObject(entry.Message)
}

func (s helloWorldService) LogHelloWorld(name domain.NameObject) (*domain.MessageObject, error) {
	entity := domain.NewHelloWorldEntity(name)
	message, err := entity.Message()
	if err != nil {
		return nil, err
	}

	entry := record.HelloWorldEntry{Name: name.Value(), Message: message.Value()}
	_, err = s.repository.CreateHelloWorld(entry)
	if err != nil {
		return nil, err
	}

	return message, nil
}
