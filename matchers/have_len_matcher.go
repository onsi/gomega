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
			return true, fmt.Sprintf("Expected%s\n not to have length %d", format.Object(actual), matcher.Count), nil
		} else {
			return false, fmt.Sprintf("Expected%s\n to have length %d", format.Object(actual), matcher.Count), nil
		}
	} else {
		return false, "", fmt.Errorf("HaveLen matcher expects a string/array/map/channel/slice.  Got:%s", format.Object(actual))
	}
}
