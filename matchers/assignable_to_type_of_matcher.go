package matchers

import (
	"fmt"
	"reflect"
	"github.com/onsi/gomega/format"
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
		return true, format.Message(actual, fmt.Sprintf("not to be assignable to the type: %T", matcher.Expected)), nil
	} else {
		return false, format.Message(actual, fmt.Sprintf("to be assignable to the type: %T", matcher.Expected)), nil
	}
}
