package matchers

import (
	"fmt"
	"reflect"
)

type HaveKeyMatcher struct {
	Key interface{}
}

func (matcher *HaveKeyMatcher) Match(actual interface{}) (success bool, message string, err error) {
	if !isMap(actual) {
		return false, "", fmt.Errorf("HaveKey matcher expects a map.  Got: %s", formatObject(actual))
	}

	if reflect.TypeOf(actual).Key() != reflect.TypeOf(matcher.Key) {
		return false, "", fmt.Errorf("Mismatch between map key type (%v) and the expected key type (%T)", reflect.TypeOf(actual).Key(), matcher.Key)
	}

	//lame that we can't just look up the key...
	//func (Value) MapIndex(key Value) returns the zero value, not an error, if the key is missing.  boo...
	keys := reflect.ValueOf(actual).MapKeys()
	for i := 0; i < len(keys); i++ {
		if reflect.DeepEqual(keys[i].Interface(), matcher.Key) {
			return true, formatMessage(actual, "not to have key", matcher.Key), nil
		}
	}
	return false, formatMessage(actual, "to have key", matcher.Key), nil
}
