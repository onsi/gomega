package matchers

import (
	"fmt"
	"reflect"
)

type EqualMatcher struct {
	Expected interface{}
}

func (matcher *EqualMatcher) Match(actual interface{}) (success bool, message string, err error) {
	if actual == nil && matcher.Expected == nil {
		return false, "", fmt.Errorf("Refusing to compare <nil> to <nil>.")
	}
	if reflect.DeepEqual(actual, matcher.Expected) {
		return true, formatMessage(actual, "not to equal", matcher.Expected), nil
	} else {
		return false, formatMessage(actual, "to equal", matcher.Expected), nil
	}
}
