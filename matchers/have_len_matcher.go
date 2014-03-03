package matchers

import (
	"fmt"
	"github.com/onsi/gomega/format"
)

type HaveLenMatcher struct {
	Count int
}

func (matcher *HaveLenMatcher) Match(actual interface{}) (success bool, message string, err error) {
	length, ok := lengthOf(actual)
	if ok {
		if length == matcher.Count {
			return true, fmt.Sprintf("Expected\n%s\nnot to have length %d", format.Object(actual, 1), matcher.Count), nil
		} else {
			return false, fmt.Sprintf("Expected\n%s\nto have length %d", format.Object(actual, 1), matcher.Count), nil
		}
	} else {
		return false, "", fmt.Errorf("HaveLen matcher expects a string/array/map/channel/slice.  Got:\n%s", format.Object(actual, 1))
	}
}
