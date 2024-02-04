package usecase

import (
	"github.com/kumachan-mis/knodeledge-api/internal/domain"
	"github.com/kumachan-mis/knodeledge-api/internal/service"
)

//go:generate mockgen -source=$GOFILE -destination=../../mock/$GOPACKAGE/mock_$GOFILE -package=$GOPACKAGE

type HelloWorldUseCase interface {
	UseHelloWorld(name string) (string, error)
}

type helloWorldUseCase struct {
	service service.HelloWorldService
}

func NewHelloWorldUseCase(service service.HelloWorldService) HelloWorldUseCase {
	return helloWorldUseCase{service: service}
}

func (uc helloWorldUseCase) UseHelloWorld(name string) (string, error) {
	n := domain.NewNameObject(name)

	m, err := uc.service.SearchHelloWorld(n)
	if err != nil {
		m, err = uc.service.LogHelloWorld(n)
	}

	if err != nil {
		return "", err
	}
	return m.Value(), nil
}
