package matchers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

func WithOrderedKeys(maybeMatcher types.GomegaMatcher, keys ...string) types.GomegaMatcher{
	matcher, ok := maybeMatcher.(*MatchUnorderedJSONMatcher)
	if !ok{
		panic("Matcher must be of type MatchUnorderedJSONMatcher")
	}
	matcher.OrderedKeys = make(map[string]bool)
	for _, v  := range keys {
		matcher.OrderedKeys[v] = true
	}
	return matcher

}

type MatchUnorderedJSONMatcher struct {
	JSONToMatch      interface{}
	firstFailurePath []interface{}
	OrderedKeys      map[string]bool
}

func (matcher *MatchUnorderedJSONMatcher) Match(actual interface{}) (success bool, err error) {
	actualString, expectedString, err := matcher.prettyPrint(actual)
	if err != nil {
		return false, err
	}

	var aval interface{}
	var eval interface{}

	// this is guarded by prettyPrint
	json.Unmarshal([]byte(actualString), &aval)
	json.Unmarshal([]byte(expectedString), &eval)
	var equal bool


	equal, matcher.firstFailurePath = matcher.deepEqual(aval, eval, false)
	return equal, nil
}

func (matcher *MatchUnorderedJSONMatcher) FailureMessage(actual interface{}) (message string) {
	actualString, expectedString, _ := matcher.prettyPrint(actual)
	return formattedMessage(format.Message(actualString, "to match JSON of", expectedString), matcher.firstFailurePath)
}

func (matcher *MatchUnorderedJSONMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	actualString, expectedString, _ := matcher.prettyPrint(actual)
	return formattedMessage(format.Message(actualString, "not to match JSON of", expectedString), matcher.firstFailurePath)
}

func (matcher *MatchUnorderedJSONMatcher) prettyPrint(actual interface{}) (actualFormatted, expectedFormatted string, err error) {
	actualString, ok := toString(actual)
	if !ok {
		return "", "", fmt.Errorf("MatchUnorderedJSONMatcher matcher requires a string, stringer, or []byte.  Got actual:\n%s", format.Object(actual, 1))
	}
	expectedString, ok := toString(matcher.JSONToMatch)
	if !ok {
		return "", "", fmt.Errorf("MatchUnorderedJSONMatcher matcher requires a string, stringer, or []byte.  Got expected:\n%s", format.Object(matcher.JSONToMatch, 1))
	}

	abuf := new(bytes.Buffer)
	ebuf := new(bytes.Buffer)

	if err := json.Indent(abuf, []byte(actualString), "", "  "); err != nil {
		return "", "", fmt.Errorf("Actual '%s' should be valid JSON, but it is not.\nUnderlying error:%s", actualString, err)
	}

	if err := json.Indent(ebuf, []byte(expectedString), "", "  "); err != nil {
		return "", "", fmt.Errorf("Expected '%s' should be valid JSON, but it is not.\nUnderlying error:%s", expectedString, err)
	}

	return abuf.String(), ebuf.String(), nil
}

func (matcher *MatchUnorderedJSONMatcher) deepEqual(a interface{}, b interface{}, ordered bool) (bool, []interface{}) {
	var errorPath []interface{}
	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return false, errorPath
	}

	switch a.(type) {
	case []interface{}:
		if ordered {
			return matcher.deepEqualOrderedList(a, b, errorPath)
		} else {
			return matcher.deepEqualUnorderedList(a, b, errorPath)
		}

	case map[string]interface{}:
		if len(a.(map[string]interface{})) != len(b.(map[string]interface{})) {
			return false, errorPath
		}

		for k, v1 := range a.(map[string]interface{}) {
			v2, ok := b.(map[string]interface{})[k]
			if !ok {
				return false, errorPath
			}

			elementEqual, keyPath := matcher.deepEqual(v1, v2, matcher.OrderedKeys[k])
			if !elementEqual {
				return false, append(keyPath, k)
			}
		}
		return true, errorPath

	default:
		return a == b, errorPath
	}
}

func (matcher *MatchUnorderedJSONMatcher) deepEqualUnorderedList(a interface{}, b interface{}, errorPath []interface{}) (bool, []interface{}) {
	if len(a.([]interface{})) != len(b.([]interface{})) {
		return false, errorPath
	}
	matched := make([]bool, len(b.([]interface{})))

	for _, v1 := range a.([]interface{}) {
		foundMatch := false
		for j, v2 := range b.([]interface{}) {
			if matched[j] {
				continue
			}
			elementEqual, _ := matcher.deepEqual(v1, v2, false)
			if elementEqual {
				foundMatch = true
				matched[j] = true
				break
			}
		}
		if !foundMatch {
			return false, errorPath
		}
	}

	return true, errorPath
}

func (matcher *MatchUnorderedJSONMatcher) deepEqualOrderedList(a interface{}, b interface{}, errorPath []interface{}) (bool, []interface{}) {
	if len(a.([]interface{})) != len(b.([]interface{})) {
		return false, errorPath
	}

	for i, v := range a.([]interface{}) {
		elementEqual, keyPath := matcher.deepEqual(v, b.([]interface{})[i], false)
		if !elementEqual {
			return false, append(keyPath, i)
		}
	}
	return true, errorPath
}
