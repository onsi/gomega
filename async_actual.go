package gomega

import (
	"fmt"
	"reflect"
	"time"
)

type asyncActual struct {
	actualInput     interface{}
	timeoutInterval time.Duration
	pollingInterval time.Duration
	fail            OmegaFailHandler
}

func newAsyncActual(actualInput interface{}, fail OmegaFailHandler, timeoutInterval time.Duration, pollingInterval time.Duration) *asyncActual {
	actualType := reflect.TypeOf(actualInput)
	if actualType.Kind() == reflect.Func {
		if actualType.NumIn() != 0 || actualType.NumOut() != 1 {
			panic("Expected a function with no arguments and one return value.")
		}
	}

	return &asyncActual{
		actualInput:     actualInput,
		fail:            fail,
		timeoutInterval: timeoutInterval,
		pollingInterval: pollingInterval,
	}
}

func (actual *asyncActual) Should(matcher OmegaMatcher, optionalDescription ...interface{}) bool {
	return actual.match(matcher, true, optionalDescription...)
}

func (actual *asyncActual) ShouldNot(matcher OmegaMatcher, optionalDescription ...interface{}) bool {
	return actual.match(matcher, false, optionalDescription...)
}

func (actual *asyncActual) buildDescription(optionalDescription ...interface{}) string {
	switch len(optionalDescription) {
	case 0:
		return ""
	default:
		return fmt.Sprintf(optionalDescription[0].(string), optionalDescription[1:]...) + "\n"
	}
}
func (actual *asyncActual) pollActual() interface{} {
	actualType := reflect.TypeOf(actual.actualInput)

	if actualType.Kind() == reflect.Func && actualType.NumIn() == 0 && actualType.NumOut() == 1 {
		return reflect.ValueOf(actual.actualInput).Call([]reflect.Value{})[0].Interface()
	}

	return actual.actualInput
}

func (actual *asyncActual) match(matcher OmegaMatcher, desiredMatch bool, optionalDescription ...interface{}) bool {
	timer := time.Now()
	timeout := time.After(actual.timeoutInterval)

	description := actual.buildDescription(optionalDescription...)
	matches, message, err := matcher.Match(actual.pollActual())

	for {
		if err == nil && matches == desiredMatch {
			return true
		}

		select {
		case <-time.After(actual.pollingInterval):
			matches, message, err = matcher.Match(actual.pollActual())
		case <-timeout:
			errMsg := ""
			if err != nil {
				errMsg = "Error: " + err.Error()
			}
			actual.fail(fmt.Sprintf("Timed out after %.3fs.\n%s%s%s", time.Since(timer).Seconds(), description, message, errMsg), 2)
			return false
		}
	}
}
