package matchers

import "github.com/onsi/gomega/types"

type WithTransformMatcher struct {
	// input
	Transform func(interface{}) interface{}
	Matcher   types.GomegaMatcher

	// state
	transformedValue interface{}
}

func (m *WithTransformMatcher) Match(actual interface{}) (bool, error) {
	m.transformedValue = m.Transform(actual)
	return m.Matcher.Match(m.transformedValue)
}

func (m *WithTransformMatcher) FailureMessage(actual interface{}) (message string) {
	return m.Matcher.FailureMessage(m.transformedValue)
}

func (m *WithTransformMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return m.Matcher.NegatedFailureMessage(m.transformedValue)
}
