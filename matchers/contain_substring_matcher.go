package matchers

import (
	"fmt"
	"strings"
)

type ContainSubstringMatcher struct {
	Substr string
}

func (matcher *ContainSubstringMatcher) Match(actual interface{}) (success bool, message string, err error) {
	actualString, ok := toString(actual)
	if ok {
		match := strings.Contains(actualString, matcher.Substr)
		if match {
			return true, formatMessage(actual, "not to contain substring", matcher.Substr), nil
		} else {
			return false, formatMessage(actual, "to contain substring", matcher.Substr), nil
		}
	} else {
		return false, "", fmt.Errorf("ContainSubstring matcher requires a string or stringer.  Got:%s", formatObject(actual))
	}
}
