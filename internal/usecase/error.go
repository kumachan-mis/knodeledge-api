package usecase

import (
	"encoding/json"
	"fmt"
)

type ErrorCode string
type ErrorMessage string

const (
	DomainValidationError ErrorCode = "domain validation error"
	InvalidArgumentError  ErrorCode = "invalid argument"
	NotFoundError         ErrorCode = "not found"
	InternalErrorPanic    ErrorCode = "internal error"
)

func NewMessageBasedError[ErrorResponse any](code ErrorCode, message string) *Error[ErrorResponse] {
	return &Error[ErrorResponse]{code: code, message: ErrorMessage(message), response: nil}
}

func NewModelBasedError[ErrorResponse any](code ErrorCode, response ErrorResponse) *Error[ErrorResponse] {
	return &Error[ErrorResponse]{code: code, message: "", response: &response}
}

type Error[ErrorResponse any] struct {
	code     ErrorCode
	message  ErrorMessage
	response *ErrorResponse
}

func (e Error[ErrorResponse]) Code() ErrorCode {
	return e.code
}

func (e Error[ErrorResponse]) Message() ErrorMessage {
	return e.message
}

func (e Error[ErrorResponse]) Response() *ErrorResponse {
	return e.response
}

func (e Error[ErrorResponse]) Error() string {
	if e.response == nil {
		return fmt.Sprintf("%s: %s", e.code, e.message)
	}

	bytes, err := json.Marshal(e.response)
	if err != nil {
		return fmt.Sprintf("%s: %s", e.code, e.message)
	}
	return fmt.Sprintf("%s: %s", e.code, string(bytes))
}
