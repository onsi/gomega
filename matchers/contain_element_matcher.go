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
	value := reflect.ValueOf(actual)
	for i := 0; i < value.Len(); i++ {
		if reflect.DeepEqual(value.Index(i).Interface(), matcher.Element) {
			return true, formatMessage(actual, "not to contain element", matcher.Element), nil
		}
	}
	return false, formatMessage(actual, "to contain element", matcher.Element), nil
}
