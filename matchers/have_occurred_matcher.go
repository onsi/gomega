package matchers

import (
	"fmt"
	"github.com/onsi/gomega/format"
)

type HaveOccurredMatcher struct {
}

func (matcher *HaveOccurredMatcher) Match(actual interface{}) (success bool, message string, err error) {
	if actual == nil {
		return false, format.Message(actual, "to have occurred"), nil
	} else {
		if isError(actual) {
			return true, fmt.Sprintf("Expected error:%s\n    %s\n%s", format.Object(actual), actual.(error).Error(), "not to have occurred"), nil
		} else {
			return false, "", fmt.Errorf("Expected an error, got%s", format.Object(actual))
		}
	}
}
