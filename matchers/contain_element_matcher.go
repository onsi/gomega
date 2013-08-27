package matchers

import (
	"fmt"
	"reflect"
)

type ContainElementMatcher struct {
	Element interface{}
}

func (matcher *ContainElementMatcher) Match(actual interface{}) (success bool, message string, err error) {
	if !isArrayOrSlice(actual) {
		return false, "", fmt.Errorf("ContainElement matcher expects an array/slice/string.  Got:%s", formatObject(actual))
	}

	elemMatcher, elementIsMatcher := matcher.Element.(omegaMatcher)
	if !elementIsMatcher {
		elemMatcher = &EqualMatcher{Expected: matcher.Element}
	}

	value := reflect.ValueOf(actual)
	for i := 0; i < value.Len(); i++ {
		success, _, err := elemMatcher.Match(value.Index(i).Interface())
		if err != nil {
			return false, "", fmt.Errorf("ContainElement's element matcher failed with:\n\t%s", err.Error())
		}
		if success {
			if elementIsMatcher {
				return true, formatMessage(actual, "not to contain element matching", matcher.Element), nil
			} else {
				return true, formatMessage(actual, "not to contain element", matcher.Element), nil
			}
		}
	}
	if elementIsMatcher {
		return false, formatMessage(actual, "to contain element matching", matcher.Element), nil
	} else {
		return false, formatMessage(actual, "to contain element", matcher.Element), nil
	}
}
