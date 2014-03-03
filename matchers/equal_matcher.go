package matchers

import (
	"fmt"
	"reflect"
	"github.com/onsi/gomega/format"
)

type EqualMatcher struct {
	Expected interface{}
}

func (matcher *EqualMatcher) Match(actual interface{}) (success bool, message string, err error) {
	if actual == nil && matcher.Expected == nil {
		return false, "", fmt.Errorf("Refusing to compare <nil> to <nil>.")
	}
	if reflect.DeepEqual(actual, matcher.Expected) {
		return true, format.Message(actual, "not to equal", matcher.Expected), nil
	} else {
		return false, format.Message(actual, "to equal", matcher.Expected), nil
	}
}
