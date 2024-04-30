package api

import (
	"fmt"

	"github.com/kumachan-mis/knodeledge-api/internal/usecase"
	"github.com/sirupsen/logrus"
)

func JsonBindErrorToMessage(err error) string {
	return "invalid request format"
}

func UseCaseErrorToMessage[ErrorResponse any](err *usecase.Error[ErrorResponse]) string {
	switch err.Code() {
	case usecase.DomainValidationError:
		return "invalid request value"
	case usecase.InvalidArgumentError:
		return fmt.Sprintf("invalid request value: %s", err.Message())
	case usecase.NotFoundError:
		return "not found"
	default:
		logrus.WithError(err).Error("internal error")
		return "internal error"
	}
}

func UseCaseErrorToResponse[ErrorResponse any](err *usecase.Error[ErrorResponse]) ErrorResponse {
	return *err.Response()
}
