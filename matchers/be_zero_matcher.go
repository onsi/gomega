package matchers

import (
	"reflect"
	"github.com/onsi/gomega/format"
)

type BeZeroMatcher struct {
}

func (matcher *BeZeroMatcher) Match(actual interface{}) (success bool, message string, err error) {
	if actual == nil {
		return true, format.Message(actual, "not to be zero-valued"), nil
	}
	zeroValue := reflect.Zero(reflect.TypeOf(actual)).Interface()
	if reflect.DeepEqual(zeroValue, actual) {
		return true, format.Message(actual, "not to be zero-valued"), nil
	} else {
		return false, format.Message(actual, "to be zero-valued"), nil
	}
}
