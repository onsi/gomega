package matchers

import (
	"reflect"
)

type FitsTypeMatcher struct {
	Expected interface{}
}

func (matcher *FitsTypeMatcher) Match(actual interface{}) (success bool, message string, err error) {
	actualType := reflect.TypeOf(actual)
	expectedType := reflect.TypeOf(matcher.Expected)

	if actualType.AssignableTo(expectedType) {
		return true, formatMessage(actual, "not fitting type", matcher.Expected), nil
	} else {
		return false, formatMessage(actual, "fitting type", matcher.Expected), nil
	}
}
