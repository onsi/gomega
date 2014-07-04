package matchers

import (
	"fmt"
	"github.com/onsi/gomega/format"
	"reflect"
)

type MatchArrayOrSliceMatcher struct {
	ExpectedArrayOrSlice interface{}
}

func (matcher *MatchArrayOrSliceMatcher) Match(actual interface{}) (success bool, err error) {
	expected := matcher.ExpectedArrayOrSlice

	if !isArrayOrSlice(expected) {
		return false, fmt.Errorf("MatchArrayOrSlice matcher expects an array or slice.  Got:\n%s", format.Object(expected, 1))
	}

	if !isArrayOrSlice(actual) {
		return false, nil
	}

	expectedValue := reflect.ValueOf(expected)
	actualValue := reflect.ValueOf(actual)
	if expectedValue.Len() != actualValue.Len() {
		return false, nil
	}

	return matchArrayOrSlice(expectedValue, actualValue), nil
}

func (matcher *MatchArrayOrSliceMatcher) FailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "to match array or slice", matcher.ExpectedArrayOrSlice)
}

func (matcher *MatchArrayOrSliceMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "not to match array or slice", matcher.ExpectedArrayOrSlice)
}

func matchArrayOrSlice(arrayOrSlice1 reflect.Value, arrayOrSlice2 reflect.Value) bool {
	netFrequencyMap := make(map[interface{}]int)

	for index := 0; index < arrayOrSlice1.Len(); index++ {
		element := arrayOrSlice1.Index(index).Interface()
		netFrequencyMap[element]++
	}

	for index := 0; index < arrayOrSlice2.Len(); index++ {
		element := arrayOrSlice2.Index(index).Interface()
		netFrequencyMap[element]--

		if netFrequencyMap[element] == 0 {
			delete(netFrequencyMap, element)
		}
	}

	return len(netFrequencyMap) == 0
}
