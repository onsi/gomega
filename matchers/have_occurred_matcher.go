package matchers

import (
	"fmt"
)

type HaveOccurredMatcher struct {
}

func (matcher *HaveOccurredMatcher) Match(actual interface{}) (success bool, message string, err error) {
	if actual == nil {
		return false, formatMessage(actual, "to have occurred"), nil
	} else {
		if isError(actual) {
			return true, fmt.Sprintf("Expected error:%s\n\tMessage: %s\n%s", formatObject(actual), actual.(error).Error(), "not to have occurred"), nil
		} else {
			return false, "", fmt.Errorf("Expected an error, got%s", formatObject(actual))
		}
	}
}
