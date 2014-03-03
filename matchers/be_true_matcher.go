package matchers

import (
	"fmt"
	"github.com/onsi/gomega/format"
)

type BeTrueMatcher struct {
}

func (matcher *BeTrueMatcher) Match(actual interface{}) (success bool, message string, err error) {
	if !isBool(actual) {
		return false, "", fmt.Errorf("Expected a boolean.  Got:\n%s", format.Object(actual, 1))
	}
	if actual == true {
		return true, format.Message(actual, "not to be true"), nil
	} else {
		return false, format.Message(actual, "to be true"), nil
	}
}
