package matchers

import (
	"fmt"
	"github.com/onsi/gomega/format"
	"reflect"
)

type BeEquivalentToMatcher struct {
	Expected interface{}
}

func (matcher *BeEquivalentToMatcher) Match(actual interface{}) (success bool, message string, err error) {
	if actual == nil && matcher.Expected == nil {
		return false, "", fmt.Errorf("Both actual and expected must not be nil.")
	}

	convertedActual := actual

	if actual != nil && matcher.Expected != nil && reflect.TypeOf(actual).ConvertibleTo(reflect.TypeOf(matcher.Expected)) {
		convertedActual = reflect.ValueOf(actual).Convert(reflect.TypeOf(matcher.Expected)).Interface()
	}

	if reflect.DeepEqual(convertedActual, matcher.Expected) {
		return true, format.Message(actual, "not to be equivalent to", matcher.Expected), nil
	} else {
		return false, format.Message(actual, "to be equivalent to", matcher.Expected), nil
	}
}
