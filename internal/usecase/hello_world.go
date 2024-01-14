package usecase

import (
	"github.com/kumachan-mis/knodeledge-api/internal/domain"
	"github.com/kumachan-mis/knodeledge-api/internal/service"
)

func HelloWorldUseCase(name string) (string, error) {
	n := domain.NewNameObject(name)

	m, err := service.SearchHelloWorld(n)
	if err != nil {
		m, err = service.LogHelloWorld(n)
	}

	if err != nil {
		return "", err
	}
	return m.Value(), nil
}
