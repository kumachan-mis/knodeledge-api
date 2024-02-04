package usecase_test

import (
	"testing"

	"github.com/kumachan-mis/knodeledge-api/internal/domain"
	"github.com/kumachan-mis/knodeledge-api/internal/usecase"
	"github.com/kumachan-mis/knodeledge-api/mock/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUseHelloWorldAlreadyGreeted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m, err := domain.NewMessageObject("Hello, UseHelloWorld Test!")
	assert.NoError(t, err)

	s := service.NewMockHelloWorldService(ctrl)
	s.EXPECT().
		SearchHelloWorld(gomock.Any()).
		Do(func(n domain.NameObject) {
			assert.Equal(t, "UseHelloWorld Test", n.Value())
		}).
		Return(m, nil)

	uc := usecase.NewHelloWorldUseCase(s)

	message, err := uc.UseHelloWorld("UseHelloWorld Test")

	assert.NoError(t, err)
	assert.Equal(t, "Hello, UseHelloWorld Test!", message)
}

func TestUseHelloWorldFitstGreeting(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m, err := domain.NewMessageObject("Hello, UseHelloWorld Test!")
	assert.NoError(t, err)

	s := service.NewMockHelloWorldService(ctrl)
	s.EXPECT().
		SearchHelloWorld(gomock.Any()).
		Do(func(n domain.NameObject) {
			assert.Equal(t, "UseHelloWorld Test", n.Value())
		}).
		Return(nil, assert.AnError)
	s.EXPECT().
		LogHelloWorld(gomock.Any()).
		Do(func(n domain.NameObject) {
			assert.Equal(t, "UseHelloWorld Test", n.Value())
		}).
		Return(m, nil)

	uc := usecase.NewHelloWorldUseCase(s)

	message, err := uc.UseHelloWorld("UseHelloWorld Test")

	assert.NoError(t, err)
	assert.Equal(t, "Hello, UseHelloWorld Test!", message)
}

func TestUseHelloWorldGreeingLogFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := service.NewMockHelloWorldService(ctrl)
	s.EXPECT().
		SearchHelloWorld(gomock.Any()).
		Do(func(n domain.NameObject) {
			assert.Equal(t, "UseHelloWorld Test", n.Value())
		}).
		Return(nil, assert.AnError)
	s.EXPECT().
		LogHelloWorld(gomock.Any()).
		Do(func(n domain.NameObject) {
			assert.Equal(t, "UseHelloWorld Test", n.Value())
		}).
		Return(nil, assert.AnError)

	uc := usecase.NewHelloWorldUseCase(s)

	message, err := uc.UseHelloWorld("UseHelloWorld Test")

	assert.Error(t, err)
	assert.Empty(t, message)
}
