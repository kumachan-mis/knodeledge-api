package api

import (
	"github.com/kumachan-mis/knodeledge-api/internal/usecase"
	"github.com/sirupsen/logrus"
)

func JsonBindErrorToMessage(err error) string {
	return "invalid request format"
}

func UseCaseErrorToMessage[ErrorResponse any](err *usecase.Error[ErrorResponse]) string {
	switch err.Code() {
	case usecase.InvalidArgumentError:
		return "invalid request value"
	default:
		logrus.WithError(err).Error("internal error")
		return "internal error"
	}
}

func UseCaseErrorToResponse[ErrorResponse any](err *usecase.Error[ErrorResponse]) ErrorResponse {
	return *err.Response()
}
