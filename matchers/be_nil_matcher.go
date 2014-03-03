package matchers

import "github.com/onsi/gomega/format"

type BeNilMatcher struct {
}

func (matcher *BeNilMatcher) Match(actual interface{}) (success bool, message string, err error) {
	if isNil(actual) {
		return true, format.Message(actual, "not to be nil"), nil
	} else {
		return false, format.Message(actual, "to be nil"), nil
	}
}
