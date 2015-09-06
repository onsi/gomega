package matchers

import (
	"fmt"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

type OrMatcher struct {
	Matchers []types.GomegaMatcher

	// state
	successfulMatcher types.GomegaMatcher
}

func (m *OrMatcher) Match(actual interface{}) (success bool, err error) {
	for _, matcher := range m.Matchers {
		success, err := matcher.Match(actual)
		if err != nil {
			return false, err
		}
		if success {
			m.successfulMatcher = matcher
			return true, nil
		}
	}
	return false, nil
}

func (m *OrMatcher) FailureMessage(actual interface{}) (message string) {
	// not the most beautiful list of matchers, but not bad either...
	return format.Message(actual, fmt.Sprintf("To satisfy at least one of these matchers: %s", m.Matchers))
}

func (m *OrMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return m.successfulMatcher.NegatedFailureMessage(actual)
}
