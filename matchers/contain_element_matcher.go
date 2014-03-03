package matchers

import (
	"fmt"
	"github.com/onsi/gomega/format"
	"reflect"
)

type ContainElementMatcher struct {
	Element interface{}
}

func (matcher *ContainElementMatcher) Match(actual interface{}) (success bool, message string, err error) {
	if !isArrayOrSlice(actual) && !isMap(actual) {
		return false, "", fmt.Errorf("ContainElement matcher expects an array/slice/map.  Got:\n%s", format.Object(actual, 1))
	}

	elemMatcher, elementIsMatcher := matcher.Element.(omegaMatcher)
	matchingString := " matching"
	if !elementIsMatcher {
		elemMatcher = &EqualMatcher{Expected: matcher.Element}
		matchingString = ""
	}

	value := reflect.ValueOf(actual)
	var keys []reflect.Value
	if isMap(actual) {
		keys = value.MapKeys()
	}
	for i := 0; i < value.Len(); i++ {
		var success bool
		var err error
		if isMap(actual) {
			success, _, err = elemMatcher.Match(value.MapIndex(keys[i]).Interface())
		} else {
			success, _, err = elemMatcher.Match(value.Index(i).Interface())
		}
		if err != nil {
			return false, "", fmt.Errorf("ContainElement's element matcher failed with:\n\t%s", err.Error())
		}
		if success {
			return true, format.Message(actual, "not to contain element"+matchingString, matcher.Element), nil
		}
	}

	return false, format.Message(actual, "to contain element"+matchingString, matcher.Element), nil
}
