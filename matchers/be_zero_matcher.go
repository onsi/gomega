package matchers

import (
	"reflect"
)

type BeZeroMatcher struct {
}

func (matcher *BeZeroMatcher) Match(actual interface{}) (success bool, message string, err error) {
	if actual == nil {
		return true, formatMessage(actual, "not to be zero-valued"), nil
	}
	zeroValue := reflect.Zero(reflect.TypeOf(actual)).Interface()
	if reflect.DeepEqual(zeroValue, actual) {
		return true, formatMessage(actual, "not to be zero-valued"), nil
	} else {
		return false, formatMessage(actual, "to be zero-valued"), nil
	}
}
