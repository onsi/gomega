package matchers

import (
	"fmt"
	"github.com/onsi/gomega/format"
)

type BeEmptyMatcher struct {
}

func (matcher *BeEmptyMatcher) Match(actual interface{}) (success bool, message string, err error) {
	length, ok := lengthOf(actual)
	if ok {
		if length == 0 {
			return true, format.Message(actual, "not to be empty"), nil
		} else {
			return false, format.Message(actual, "to be empty"), nil
		}
	} else {
		return false, "", fmt.Errorf("BeEmpty matcher expects a string/array/map/channel/slice.  Got:\n%s", format.Object(actual, 1))
	}
}
