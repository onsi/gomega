package matchers

import (
	"fmt"
	"github.com/onsi/gomega/format"
	"reflect"
)

type MatchErrorMatcher struct {
	Expected interface{}
}

func (matcher *MatchErrorMatcher) Match(actual interface{}) (success bool, err error) {
	if isNil(actual) {
		return false, fmt.Errorf("Expected an error, got nil")
	}

	if !isError(actual) {
		return false, fmt.Errorf("Expected an error.  Got:\n%s", format.Object(actual, 1))
	}

	actualErr := actual.(error)

	if isString(matcher.Expected) {
		return reflect.DeepEqual(actualErr.Error(), matcher.Expected), nil
	}

	if isError(matcher.Expected) {
		return reflect.DeepEqual(actualErr, matcher.Expected), nil
	}

	return false, fmt.Errorf("MatchError must be passed an error or string.  Got:\n%s", format.Object(matcher.Expected, 1))
}

func (matcher *MatchErrorMatcher) FailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "to match error", matcher.Expected)
}

func (matcher *MatchErrorMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "not to match error", matcher.Expected)
}
