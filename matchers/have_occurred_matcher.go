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
			return true, fmt.Sprintf("Expected error:\n%s\n%s\n%s", format.Object(actual, 1), format.IndentString(actual.(error).Error(), 1), "not to have occurred"), nil
		} else {
			return false, "", fmt.Errorf("Expected an error.  Got:\n%s", format.Object(actual, 1))
		}
	}
}
