package service_test

import (
	"testing"

	"github.com/kumachan-mis/knodeledge-api/internal/domain"
	"github.com/kumachan-mis/knodeledge-api/internal/record"
	"github.com/kumachan-mis/knodeledge-api/internal/service"
	"github.com/kumachan-mis/knodeledge-api/mock/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSearchHelloWorldFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	entry := record.HelloWorldEntry{Name: "SearchHelloWorld Test", Message: "Hello, SearchHelloWorld Test!"}
	r := repository.NewMockHelloWorldRepository(ctrl)
	r.EXPECT().
		FetchHelloWorld(entry.Name).
		Return("SEARCH_HELLO_WORLD_TEST_DOC", &entry, nil)

	s := service.NewHelloWorldService(r)
	m, err := s.SearchHelloWorld(domain.NewNameObject("SearchHelloWorld Test"))

	assert.NoError(t, err)
	assert.Equal(t, "Hello, SearchHelloWorld Test!", m.Value())
}

func TestSearchHelloWorldNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := repository.NewMockHelloWorldRepository(ctrl)
	r.EXPECT().
		FetchHelloWorld("SearchHelloWorld Test").
		Return("", nil, assert.AnError)

	s := service.NewHelloWorldService(r)
	m, err := s.SearchHelloWorld(domain.NewNameObject("SearchHelloWorld Test"))

	assert.Error(t, err)
	assert.Nil(t, m)
}

func TestLogHelloWorldSucceeded(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := repository.NewMockHelloWorldRepository(ctrl)
	r.EXPECT().
		CreateHelloWorld(gomock.Any()).
		Do(func(entry record.HelloWorldEntry) {
			assert.Equal(t, "LogHelloWorld Test", entry.Name)
			assert.Equal(t, "Hello, LogHelloWorld Test!", entry.Message)
		}).
		Return("LOG_HELLO_WORLD_TEST_DOC", nil)

	s := service.NewHelloWorldService(r)
	m, err := s.LogHelloWorld(domain.NewNameObject("LogHelloWorld Test"))

	assert.NoError(t, err)
	assert.Equal(t, "Hello, LogHelloWorld Test!", m.Value())
}

func TestLogHelloWorldFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := repository.NewMockHelloWorldRepository(ctrl)
	r.EXPECT().
		CreateHelloWorld(gomock.Any()).
		Do(func(entry record.HelloWorldEntry) {
			assert.Equal(t, "LogHelloWorld Test", entry.Name)
			assert.Equal(t, "Hello, LogHelloWorld Test!", entry.Message)
		}).
		Return("", assert.AnError)

	s := service.NewHelloWorldService(r)
	m, err := s.LogHelloWorld(domain.NewNameObject("LogHelloWorld Test"))

	assert.Error(t, err)
	assert.Nil(t, m)
}
