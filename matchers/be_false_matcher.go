package matchers

import (
	"fmt"
	"github.com/onsi/gomega/format"
)

type BeFalseMatcher struct {
}

func (matcher *BeFalseMatcher) Match(actual interface{}) (success bool, message string, err error) {
	if !isBool(actual) {
		return false, "", fmt.Errorf("Expected a boolean.  Got:\n%s", format.Object(actual, 1))
	}
	if actual == false {
		return true, format.Message(actual, "not to be false"), nil
	} else {
		return false, format.Message(actual, "to be false"), nil
	}
}
