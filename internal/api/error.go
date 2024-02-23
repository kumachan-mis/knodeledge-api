package api

import (
	"github.com/kumachan-mis/knodeledge-api/internal/usecase"
	"github.com/sirupsen/logrus"
)

func JsonBindErrorToResponseMessage(err error) string {
	return "invalid request format"
}

func JsonBindErrorToResponseModel[ErrorModel any](err error, defaultModel ErrorModel) ErrorModel {
	return defaultModel
}

func UseCaseErrorToResponseMessage[ErrorModel any](err *usecase.Error[ErrorModel]) string {
	switch err.Code() {
	case usecase.InvalidArgumentError:
		return "invalid request value"
	default:
		logrus.WithError(err).Error("internal error")
		return "internal error"
	}
}

func UseCaseErrorToResponseModel[ErrorModel any](err *usecase.Error[ErrorModel], defaultModel ErrorModel) ErrorModel {
	if err.Model() == nil {
		return defaultModel
	}
	return *err.Model()
}
