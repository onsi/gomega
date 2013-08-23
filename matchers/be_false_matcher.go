package matchers

import (
	"fmt"
)

type BeFalseMatcher struct {
}

func (matcher *BeFalseMatcher) Match(actual interface{}) (success bool, message string, err error) {
	if !isBool(actual) {
		return false, "", fmt.Errorf("Expected a boolean, got%s", formatObject(actual))
	}
	if actual == false {
		return true, formatMessage(actual, "not to be false"), nil
	} else {
		return false, formatMessage(actual, "to be false"), nil
	}
}
