package matchers

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/onsi/gomega/format"
	"gopkg.in/yaml.v2"
)

type MatchYAMLMatcher struct {
	YAMLToMatch interface{}
	firstFailurePath []interface{}

}

func (matcher *MatchYAMLMatcher) Match(actual interface{}) (success bool, err error) {
	actualString, expectedString, err := matcher.toStrings(actual)
	if err != nil {
		return false, err
	}

	var aval interface{}
	var eval interface{}

	if err := yaml.Unmarshal([]byte(actualString), &aval); err != nil {
		return false, fmt.Errorf("Actual '%s' should be valid YAML, but it is not.\nUnderlying error:%s", actualString, err)
	}
	if err := yaml.Unmarshal([]byte(expectedString), &eval); err != nil {
		return false, fmt.Errorf("Expected '%s' should be valid YAML, but it is not.\nUnderlying error:%s", expectedString, err)
	}

	var equal bool
	equal, matcher.firstFailurePath = deepEqualYAML(aval, eval)
	return equal, nil
}

func (matcher *MatchYAMLMatcher) FailureMessage(actual interface{}) (message string) {
	actualString, expectedString, _ := matcher.toNormalisedStrings(actual)
	return formattedMessage(format.Message(actualString, "to match YAML of", expectedString), matcher.firstFailurePath)
}

func (matcher *MatchYAMLMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	actualString, expectedString, _ := matcher.toNormalisedStrings(actual)
	return formattedMessage(format.Message(actualString, "not to match YAML of", expectedString), matcher.firstFailurePath)
}

func (matcher *MatchYAMLMatcher) toNormalisedStrings(actual interface{}) (actualFormatted, expectedFormatted string, err error) {
	actualString, expectedString, err := matcher.toStrings(actual)
	return normalise(actualString), normalise(expectedString), err
}

func normalise(input string) string {
	var val interface{}
	err := yaml.Unmarshal([]byte(input), &val)
	if err != nil {
		panic(err) // guarded by Match
	}
	output, err := yaml.Marshal(val)
	if err != nil {
		panic(err) // guarded by Unmarshal
	}
	return strings.TrimSpace(string(output))
}

func (matcher *MatchYAMLMatcher) toStrings(actual interface{}) (actualFormatted, expectedFormatted string, err error) {
	actualString, ok := toString(actual)
	if !ok {
		return "", "", fmt.Errorf("MatchYAMLMatcher matcher requires a string, stringer, or []byte.  Got actual:\n%s", format.Object(actual, 1))
	}
	expectedString, ok := toString(matcher.YAMLToMatch)
	if !ok {
		return "", "", fmt.Errorf("MatchYAMLMatcher matcher requires a string, stringer, or []byte.  Got expected:\n%s", format.Object(matcher.YAMLToMatch, 1))
	}

	return actualString, expectedString, nil
}

func deepEqualYAML(a interface{}, b interface{}) (bool, []interface{}) {
	var errorPath []interface{}
	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return false, errorPath
	}

	switch a.(type) {
	case []interface{}:
		if len(a.([]interface{})) != len(b.([]interface{})) {
			return false, errorPath
		}

		for i, v := range a.([]interface{}) {
			elementEqual, keyPath := deepEqualYAML(v, b.([]interface{})[i])
			if !elementEqual {
				return false, append(keyPath, i)
			}
		}
		return true, errorPath

	case map[interface{}]interface{}:
		if len(a.(map[interface{}]interface{})) != len(b.(map[interface{}]interface{})) {
			return false, errorPath
		}

		for k, v1 := range a.(map[interface{}]interface{}) {
			v2, ok := b.(map[interface{}]interface{})[k]
			if !ok {
				return false, errorPath
			}
			elementEqual, keyPath := deepEqualYAML(v1, v2)
			if !elementEqual {
				return false, append(keyPath, k)
			}
		}
		return true, errorPath

	default:
		return a == b, errorPath
	}
}
