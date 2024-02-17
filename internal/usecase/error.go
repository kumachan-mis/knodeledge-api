package usecase

import (
	"encoding/json"
	"fmt"
)

type ErrorCode string
type ErrorMessage string

const (
	InvalidArgumentError ErrorCode = "invalid argument"
	InternalError        ErrorCode = "internal error"
)

func NewMessageBasedError[ErrorModel any](code ErrorCode, message string) *Error[ErrorModel] {
	return &Error[ErrorModel]{code: code, message: ErrorMessage(message), model: nil}
}

func NewModelBasedError[ErrorModel any](code ErrorCode, model ErrorModel) *Error[ErrorModel] {
	return &Error[ErrorModel]{code: code, message: "", model: &model}
}

type Error[ErrorModel any] struct {
	code    ErrorCode
	message ErrorMessage
	model   *ErrorModel
}

func (e Error[ErrorModel]) Code() ErrorCode {
	return e.code
}

func (e Error[ErrorModel]) Model() *ErrorModel {
	return e.model
}

func (e Error[ErrorModel]) Error() string {
	if e.model == nil {
		return fmt.Sprintf("%s: %s", e.code, e.message)
	}

	bytes, err := json.Marshal(e.model)
	if err != nil {
		return fmt.Sprintf("%s: %s", e.code, e.message)
	}
	return fmt.Sprintf("%s: %s", e.code, string(bytes))
}
