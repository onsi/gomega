package internal

import (
	"errors"
	"fmt"
)

type StopTryingError interface {
	error
	Now()
}

var StopTrying = func(format string, a ...any) StopTryingError {
	err := fmt.Errorf(format, a...)
	return &AsyncSignalError{
		message:    err.Error(),
		wrappedErr: errors.Unwrap(err),
		stopTrying: true,
	}
}

type AsyncSignalError struct {
	message    string
	wrappedErr error
	stopTrying bool
	viaPanic   bool
}

func (s *AsyncSignalError) Error() string {
	return s.message
}

func (s *AsyncSignalError) Unwrap() error {
	if s == nil {
		return nil
	}
	return s.wrappedErr
}

func (s *AsyncSignalError) Now() {
	s.viaPanic = true
	panic(s)
}

func (s *AsyncSignalError) WasViaPanic() bool {
	return s.viaPanic
}

func (s *AsyncSignalError) StopTrying() bool {
	return s.stopTrying
}

func AsAsyncSignalError(actual interface{}) (*AsyncSignalError, bool) {
	if actual == nil {
		return nil, false
	}
	if actualErr, ok := actual.(error); ok {
		var target *AsyncSignalError
		if errors.As(actualErr, &target) {
			return target, true
		} else {
			return nil, false
		}
	}

	return nil, false
}
