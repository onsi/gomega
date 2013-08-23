package matchers

import (
	"fmt"
)

type HaveOccuredMatcher struct {
}

func (matcher *HaveOccuredMatcher) Match(actual interface{}) (success bool, message string, err error) {
	if actual == nil {
		return false, formatMessage(actual, "to have occured"), nil
	} else {
		if isError(actual) {
			return true, fmt.Sprintf("Expected error:%s\n\tMessage: %s\n%s", formatObject(actual), actual.(error).Error(), "not to have occured"), nil
		} else {
			return false, "", fmt.Errorf("Expected an error, got%s", formatObject(actual))
		}
	}
}
