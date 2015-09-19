package matchers

import (
	"fmt"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

type AndMatcher struct {
	Matchers []types.GomegaMatcher

	// state
	firstFailedMatchErrMsg string
}

func (m *AndMatcher) Match(actual interface{}) (success bool, err error) {
	for _, matcher := range m.Matchers {
		success, err := matcher.Match(actual)
		if !success || err != nil {
			m.firstFailedMatchErrMsg = matcher.FailureMessage(actual)
			return false, err
		}
	}
	return true, nil
}

func (m *AndMatcher) FailureMessage(_ interface{}) (message string) {
	return m.firstFailedMatchErrMsg
}

func (m *AndMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	// not the most beautiful list of matchers, but not bad either...
	return format.Message(actual, fmt.Sprintf("To not satisfy all of these matchers: %s", m.Matchers))
}
