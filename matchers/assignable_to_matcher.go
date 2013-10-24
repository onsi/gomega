package matchers

import (
	"fmt"
	"reflect"
)

type AssignableToTypeOfMatcher struct {
	Expected interface{}
}

func (matcher *AssignableToTypeOfMatcher) Match(actual interface{}) (success bool, message string, err error) {
	if actual == nil || matcher.Expected == nil {
		return false, "", fmt.Errorf("Refusing to compare <nil> to <nil>.")
	}

	actualType := reflect.TypeOf(actual)
	expectedType := reflect.TypeOf(matcher.Expected)

	if actualType.AssignableTo(expectedType) {
		return true, formatMessage(actual, fmt.Sprintf("not to be assignable to the type: %T", matcher.Expected)), nil
	} else {
		return false, formatMessage(actual, fmt.Sprintf("to be assignable to the type: %T", matcher.Expected)), nil
	}
}
