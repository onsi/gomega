package matchers

import (
	"fmt"
)

type BeTrueMatcher struct {
}

func (matcher *BeTrueMatcher) Match(actual interface{}) (success bool, message string, err error) {
	if !isBool(actual) {
		return false, "", fmt.Errorf("Expected a boolean, got%s", formatObject(actual))
	}
	if actual == true {
		return true, formatMessage(actual, "not to be true"), nil
	} else {
		return false, formatMessage(actual, "to be true"), nil
	}
}
