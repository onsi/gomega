package gomega

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

type asyncActualType uint

const (
	asyncActualTypeEventually asyncActualType = iota
	asyncActualTypeConsistently
)

type asyncActual struct {
	asyncType       asyncActualType
	actualInput     interface{}
	timeoutInterval time.Duration
	pollingInterval time.Duration
	fail            OmegaFailHandler
	offset          int
}

func newAsyncActual(asyncType asyncActualType, actualInput interface{}, fail OmegaFailHandler, timeoutInterval time.Duration, pollingInterval time.Duration, offset int) *asyncActual {
	actualType := reflect.TypeOf(actualInput)
	if actualType.Kind() == reflect.Func {
		if actualType.NumIn() != 0 || actualType.NumOut() == 0 {
			panic("Expected a function with no arguments and one or more return values.")
		}
	}

	return &asyncActual{
		asyncType:       asyncType,
		actualInput:     actualInput,
		fail:            fail,
		timeoutInterval: timeoutInterval,
		pollingInterval: pollingInterval,
		offset:          offset,
	}
}

func (actual *asyncActual) Should(matcher interface{}, optionalDescription ...interface{}) bool {
	return actual.match(shimIfNecessary(matcher), true, optionalDescription...)
}

func (actual *asyncActual) ShouldNot(matcher interface{}, optionalDescription ...interface{}) bool {
	return actual.match(shimIfNecessary(matcher), false, optionalDescription...)
}

func (actual *asyncActual) buildDescription(optionalDescription ...interface{}) string {
	switch len(optionalDescription) {
	case 0:
		return ""
	default:
		return fmt.Sprintf(optionalDescription[0].(string), optionalDescription[1:]...) + "\n"
	}
}

func (actual *asyncActual) actualInputIsAFunction() bool {
	actualType := reflect.TypeOf(actual.actualInput)
	return actualType.Kind() == reflect.Func && actualType.NumIn() == 0 && actualType.NumOut() > 0
}

func (actual *asyncActual) pollActual() (interface{}, error) {
	if actual.actualInputIsAFunction() {
		values := reflect.ValueOf(actual.actualInput).Call([]reflect.Value{})

		extras := []interface{}{}
		for _, value := range values[1:] {
			extras = append(extras, value.Interface())
		}

		success, message := vetExtras(extras)

		if !success {
			return nil, errors.New(message)
		}

		return values[0].Interface(), nil
	}

	return actual.actualInput, nil
}

type oracleMatcher interface {
	MatchMayChangeInTheFuture(actual interface{}) bool
}

func (actual *asyncActual) matcherMayChange(matcher OmegaMatcher, value interface{}) bool {
	if actual.actualInputIsAFunction() {
		return true
	}

	oracleMatcher, ok := matcher.(oracleMatcher)
	if !ok {
		return true
	}

	return oracleMatcher.MatchMayChangeInTheFuture(value)
}

func (actual *asyncActual) match(matcher OmegaMatcher, desiredMatch bool, optionalDescription ...interface{}) bool {
	timer := time.Now()
	timeout := time.After(actual.timeoutInterval)

	description := actual.buildDescription(optionalDescription...)

	var matches bool
	var err error
	mayChange := true
	value, err := actual.pollActual()
	if err == nil {
		mayChange = actual.matcherMayChange(matcher, value)
		matches, err = matcher.Match(value)
	}

	fail := func(preamble string) {
		errMsg := ""
		message := ""
		if err != nil {
			errMsg = "Error: " + err.Error()
		} else {
			if desiredMatch {
				message = matcher.FailureMessage(value)
			} else {
				message = matcher.NegatedFailureMessage(value)
			}
		}
		actual.fail(fmt.Sprintf("%s after %.3fs.\n%s%s%s", preamble, time.Since(timer).Seconds(), description, message, errMsg), 3+actual.offset)
	}

	if actual.asyncType == asyncActualTypeEventually {
		for {
			if err == nil && matches == desiredMatch {
				return true
			}

			if !mayChange {
				fail("No future change is possible.  Bailing out early")
				return false
			}

			select {
			case <-time.After(actual.pollingInterval):
				value, err = actual.pollActual()
				if err == nil {
					mayChange = actual.matcherMayChange(matcher, value)
					matches, err = matcher.Match(value)
				}
			case <-timeout:
				fail("Timed out")
				return false
			}
		}
	} else if actual.asyncType == asyncActualTypeConsistently {
		for {
			if !(err == nil && matches == desiredMatch) {
				fail("Failed")
				return false
			}

			if !mayChange {
				return true
			}

			select {
			case <-time.After(actual.pollingInterval):
				value, err = actual.pollActual()
				if err == nil {
					mayChange = actual.matcherMayChange(matcher, value)
					matches, err = matcher.Match(value)
				}
			case <-timeout:
				return true
			}
		}
	}

	return false
}
