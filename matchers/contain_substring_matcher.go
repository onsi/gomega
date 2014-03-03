package matchers

import (
	"fmt"
	"github.com/onsi/gomega/format"
	"strings"
)

type ContainSubstringMatcher struct {
	Substr string
	Args   []interface{}
}

func (matcher *ContainSubstringMatcher) Match(actual interface{}) (success bool, message string, err error) {
	actualString, ok := toString(actual)
	if ok {
		stringToMatch := matcher.Substr
		if len(matcher.Args) > 0 {
			stringToMatch = fmt.Sprintf(matcher.Substr, matcher.Args...)
		}
		match := strings.Contains(actualString, stringToMatch)
		if match {
			return true, format.Message(actual, "not to contain substring", stringToMatch), nil
		} else {
			return false, format.Message(actual, "to contain substring", stringToMatch), nil
		}
	} else {
		return false, "", fmt.Errorf("ContainSubstring matcher requires a string or stringer.  Got:\n%s", format.Object(actual, 1))
	}
}
