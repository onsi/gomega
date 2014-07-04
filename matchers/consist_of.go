package matchers

import (
	"fmt"
	"reflect"

	"github.com/onsi/gomega/format"
)

type ConsistOfMatcher struct {
	Elements []interface{}
}

func (matcher *ConsistOfMatcher) Match(actual interface{}) (success bool, err error) {
	if !isArrayOrSlice(actual) && !isMap(actual) {
		return false, fmt.Errorf("ConsistOf matcher expects an array/slice/map.  Got:\n%s", format.Object(actual, 1))
	}

	elements := matcher.Elements
	if len(matcher.Elements) == 1 && isArrayOrSlice(matcher.Elements[0]) {
		elements = []interface{}{}
		value := reflect.ValueOf(matcher.Elements[0])
		for i := 0; i < value.Len(); i++ {
			elements = append(elements, value.Index(i).Interface())
		}
	}

	matchers := map[int]omegaMatcher{}
	for i, element := range elements {
		matcher, isMatcher := element.(omegaMatcher)
		if !isMatcher {
			matcher = &EqualMatcher{Expected: element}
		}
		matchers[i] = matcher
	}

	values := matcher.valuesOf(actual)

	if len(values) != len(matchers) {
		return false, nil
	}

	for _, value := range values {
		found := false
		for key, matcher := range matchers {
			success, err := matcher.Match(value)
			if err != nil {
				continue
			}
			if success {
				found = true
				delete(matchers, key)
				break
			}
		}

		if !found {
			return false, nil
		}
	}

	return true, nil
}

func (matcher *ConsistOfMatcher) valuesOf(actual interface{}) []interface{} {
	value := reflect.ValueOf(actual)
	values := []interface{}{}
	if isMap(actual) {
		keys := value.MapKeys()
		for i := 0; i < value.Len(); i++ {
			values = append(values, value.MapIndex(keys[i]).Interface())
		}
	} else {
		for i := 0; i < value.Len(); i++ {
			values = append(values, value.Index(i).Interface())
		}
	}

	return values
}

func (matcher *ConsistOfMatcher) FailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "to consist of", matcher.Elements)
}

func (matcher *ConsistOfMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "not to consist of", matcher.Elements)
}
