package matchers

import (
	"fmt"
	"reflect"
)

type AssignableToMatcher struct {
	Expected interface{}
}

func (matcher *AssignableToMatcher) Match(actual interface{}) (success bool, message string, err error) {
	if actual == nil || matcher.Expected == nil {
		return false, "", fmt.Errorf("Refusing to compare <nil> to <nil>.")
	}

	actualType := reflect.TypeOf(actual)
	expectedType := reflect.TypeOf(matcher.Expected)

	if actualType.AssignableTo(expectedType) {
		return true, formatMessage(actual, "not fitting type", matcher.Expected), nil
	} else {
		return false, formatMessage(actual, "fitting type", matcher.Expected), nil
	}
}
