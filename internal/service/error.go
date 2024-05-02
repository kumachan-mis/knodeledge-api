package service

import (
	"fmt"
)

type ErrorCode string

const (
	InvalidArgument        ErrorCode = "invalid argument"
	NotFoundError          ErrorCode = "not found"
	DomainFailurePanic     ErrorCode = "domain failure"
	RepositoryFailurePanic ErrorCode = "repository failure"
)

type Error struct {
	code ErrorCode
	err  error
}

func Errorf(code ErrorCode, format string, args ...any) *Error {
	return &Error{code: code, err: fmt.Errorf(format, args...)}
}

func (e Error) Error() string {
	return fmt.Sprintf("%v: %v", e.code, e.err.Error())
}

func (e Error) Unwrap() error {
	return e.err
}

func (e Error) Code() ErrorCode {
	return e.code
}
