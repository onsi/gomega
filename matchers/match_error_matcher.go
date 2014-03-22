package matchers

import (
	"fmt"
	"github.com/onsi/gomega/format"
	"reflect"
)

type MatchErrorMatcher struct {
	Expected interface{}
}

func (matcher *MatchErrorMatcher) Match(actual interface{}) (success bool, message string, err error) {
	if isNil(actual) {
		return false, "", fmt.Errorf("Expected an error, got nil")
	}

	if !isError(actual) {
		return false, "", fmt.Errorf("Expected an error.  Got:\n%s", format.Object(actual, 1))
	}

	actualErr := actual.(error)

	if isString(matcher.Expected) {
		if reflect.DeepEqual(actualErr.Error(), matcher.Expected) {
			return true, format.Message(actual, "not to match error", matcher.Expected), nil
		} else {
			return false, format.Message(actual, "to match error", matcher.Expected), nil
		}
	}

	if isError(matcher.Expected) {
		if reflect.DeepEqual(actualErr, matcher.Expected) {
			return true, format.Message(actual, "not to match error", matcher.Expected), nil
		} else {
			return false, format.Message(actual, "to match error", matcher.Expected), nil
		}
	}

	return false, "", fmt.Errorf("MatchError must be passed an error or string.  Got:\n%s", format.Object(matcher.Expected, 1))
}
