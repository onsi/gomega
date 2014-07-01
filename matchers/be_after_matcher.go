package matchers

import (
	"fmt"
	"github.com/onsi/gomega/format"
	"time"
)

type BeAfterMatcher struct {
	Expected interface{}
}

func (m *BeAfterMatcher) Match(actual interface{}) (bool, error) {
	e, ok := m.Expected.(time.Time)
	if !ok {
		return false, fmt.Errorf("Expectation is not a time.Time: %s", format.Object(m.Expected, 1))
	}
	a, ok := actual.(time.Time)
	if !ok {
		return false, fmt.Errorf("Expected a time.Time. Got: %s", format.Object(actual, 1))
	}
	return a.After(e), nil
}

func (m *BeAfterMatcher) FailureMessage(actual interface{}) string {
	return format.Message(actual, "to be after", m.Expected)
}

func (m *BeAfterMatcher) NegatedFailureMessage(actual interface{}) string {
	return format.Message(actual, "not to be after", m.Expected)
}
