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

	keyMatcher, keyIsMatcher := matcher.Key.(omegaMatcher)
	if !keyIsMatcher {
		keyMatcher = &EqualMatcher{Expected: matcher.Key}
	}

	keys := reflect.ValueOf(actual).MapKeys()
	for i := 0; i < len(keys); i++ {
		success, _, err := keyMatcher.Match(keys[i].Interface())

		if err != nil {
			return false, "", fmt.Errorf("HaveKey's key matcher failed with:\n\t%s", err.Error())
		}
		if success {
			if keyIsMatcher {
				return true, formatMessage(actual, "not to have key matching", matcher.Key), nil
			} else {
				return true, formatMessage(actual, "not to have key", matcher.Key), nil
			}
		}
	}
	if keyIsMatcher {
		return false, formatMessage(actual, "to have key matching", matcher.Key), nil
	} else {
		return false, formatMessage(actual, "to have key", matcher.Key), nil
	}
}
