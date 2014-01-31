package matchers

import (
	"fmt"
)

type HaveLenMatcher struct {
	Count int
}

func (matcher *HaveLenMatcher) Match(actual interface{}) (success bool, message string, err error) {
	length, ok := lengthOf(actual)
	if ok {
		if length == matcher.Count {
			return true, fmt.Sprintf("Expected%s\n (length: %d) not to have length %d", formatObject(actual), length, matcher.Count), nil
		} else {
			return false, fmt.Sprintf("Expected%s\n (length: %d) to have length %d", formatObject(actual), length, matcher.Count), nil
		}
	} else {
		return false, "", fmt.Errorf("HaveLen matcher expects a string/array/map/channel/slice.  Got:%s", formatObject(actual))
	}
}
