package repository

import (
	"fmt"
)

type ErrorCode string

const (
	InValidArgument   ErrorCode = "invalid argument"
	NotFoundError     ErrorCode = "not found"
	ReadFailurePanic  ErrorCode = "read failure"
	WriteFailurePanic ErrorCode = "write failure"
)

type Error struct {
	code ErrorCode
	err  error
}

func Errorf(code ErrorCode, format string, args ...any) *Error {
	return &Error{code: code, err: fmt.Errorf(format, args...)}
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", e.code, e.err.Error())
}

func (e Error) Unwrap() error {
	return e.err
}

func (e Error) Code() ErrorCode {
	return e.code
}
