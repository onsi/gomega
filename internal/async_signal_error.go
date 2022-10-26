package internal

import (
	"errors"
	"time"
)

type AsyncSignalErrorType int

const (
	AsyncSignalErrorTypeStopTrying AsyncSignalErrorType = iota
	AsyncSignalErrorTypeTryAgainAfter
)

type StopTryingError interface {
	error
	Wrap(err error) StopTryingError
	Attach(description string, obj any) StopTryingError
	Now()
}

type TryAgainAfterError interface {
	error
	Now()
}

var StopTrying = func(message string) StopTryingError {
	return &AsyncSignalError{
		message:              message,
		asyncSignalErrorType: AsyncSignalErrorTypeStopTrying,
	}
}

var TryAgainAfter = func(duration time.Duration) TryAgainAfterError {
	return &AsyncSignalError{
		duration:             duration,
		asyncSignalErrorType: AsyncSignalErrorTypeTryAgainAfter,
	}
}

type AsyncSignalErrorAttachment struct {
	Description string
	Object      any
}

type AsyncSignalError struct {
	message              string
	wrappedErr           error
	asyncSignalErrorType AsyncSignalErrorType
	duration             time.Duration
	Attachments          []AsyncSignalErrorAttachment
}

func (s *AsyncSignalError) Wrap(err error) StopTryingError {
	s.wrappedErr = err
	return s
}

func (s *AsyncSignalError) Attach(description string, obj any) StopTryingError {
	s.Attachments = append(s.Attachments, AsyncSignalErrorAttachment{description, obj})
	return s
}

func (s *AsyncSignalError) Error() string {
	if s.wrappedErr == nil {
		return s.message
	} else {
		return s.message + ": " + s.wrappedErr.Error()
	}
}

func (s *AsyncSignalError) Unwrap() error {
	if s == nil {
		return nil
	}
	return s.wrappedErr
}

func (s *AsyncSignalError) Now() {
	panic(s)
}

func (s *AsyncSignalError) IsStopTrying() bool {
	return s.asyncSignalErrorType == AsyncSignalErrorTypeStopTrying
}

func (s *AsyncSignalError) IsTryAgainAfter() bool {
	return s.asyncSignalErrorType == AsyncSignalErrorTypeTryAgainAfter
}

func (s *AsyncSignalError) TryAgainDuration() time.Duration {
	return s.duration
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
