package internal

import (
	"errors"
	"time"
	"fmt"
)

type AsyncSignalErrorType int

const (
	AsyncSignalErrorTypeStopTrying AsyncSignalErrorType = iota
	AsyncSignalErrorTypeTryAgainAfter
)

type AsyncSignalError interface {
	error
	Wrap(err error) AsyncSignalError
	Attach(description string, obj any) AsyncSignalError
	Now()
}


var StopTrying = func(message string) AsyncSignalError {
	return &AsyncSignalErrorImpl{
		message:              message,
		asyncSignalErrorType: AsyncSignalErrorTypeStopTrying,
	}
}

var TryAgainAfter = func(duration time.Duration) AsyncSignalError {
	return &AsyncSignalErrorImpl{
		message: fmt.Sprintf("told to try again after %s", duration),
		duration:             duration,
		asyncSignalErrorType: AsyncSignalErrorTypeTryAgainAfter,
	}
}

type AsyncSignalErrorAttachment struct {
	Description string
	Object      any
}

type AsyncSignalErrorImpl struct {
	message              string
	wrappedErr           error
	asyncSignalErrorType AsyncSignalErrorType
	duration             time.Duration
	Attachments          []AsyncSignalErrorAttachment
}

func (s *AsyncSignalErrorImpl) Wrap(err error) AsyncSignalError {
	s.wrappedErr = err
	return s
}

func (s *AsyncSignalErrorImpl) Attach(description string, obj any) AsyncSignalError {
	s.Attachments = append(s.Attachments, AsyncSignalErrorAttachment{description, obj})
	return s
}

func (s *AsyncSignalErrorImpl) Error() string {
	if s.wrappedErr == nil {
		return s.message
	} else {
		return s.message + ": " + s.wrappedErr.Error()
	}
}

func (s *AsyncSignalErrorImpl) Unwrap() error {
	if s == nil {
		return nil
	}
	return s.wrappedErr
}

func (s *AsyncSignalErrorImpl) Now() {
	panic(s)
}

func (s *AsyncSignalErrorImpl) IsStopTrying() bool {
	return s.asyncSignalErrorType == AsyncSignalErrorTypeStopTrying
}

func (s *AsyncSignalErrorImpl) IsTryAgainAfter() bool {
	return s.asyncSignalErrorType == AsyncSignalErrorTypeTryAgainAfter
}

func (s *AsyncSignalErrorImpl) TryAgainDuration() time.Duration {
	return s.duration
}

func AsAsyncSignalError(actual interface{}) (*AsyncSignalErrorImpl, bool) {
	if actual == nil {
		return nil, false
	}
	if actualErr, ok := actual.(error); ok {
		var target *AsyncSignalErrorImpl
		if errors.As(actualErr, &target) {
			return target, true
		} else {
			return nil, false
		}
	}

	return nil, false
}
