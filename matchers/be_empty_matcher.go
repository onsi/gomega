package matchers

import (
	"fmt"
)

type BeEmptyMatcher struct {
}

func (matcher *BeEmptyMatcher) Match(actual interface{}) (success bool, message string, err error) {
	length, ok := lengthOf(actual)
	if ok {
		if length == 0 {
			return true, formatMessage(actual, "not to be empty"), nil
		} else {
			return false, formatMessage(actual, "to be empty"), nil
		}
	} else {
		return false, "", fmt.Errorf("BeEmpty matcher expects a string/array/map/channel/slice.  Got:%s", formatObject(actual))
	}
}
