package matchers

import (
	"reflect"
)

type BeNilMatcher struct {
}

func (matcher *BeNilMatcher) Match(actual interface{}) (success bool, message string, err error) {
	if actual == nil {
		return true, formatMessage(actual, "not to be nil"), nil
	} else {
		t := reflect.ValueOf(actual)
		if t.Kind() == reflect.Ptr {
			return t.IsNil(), formatMessage(actual, "to be nil"), nil
		}

		return false, formatMessage(actual, "to be nil"), nil
	}
}
