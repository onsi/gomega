package gomega

import (
	"fmt"
)

var didWarnAboutMatcherDeprecation map[string]bool

func init() {
	didWarnAboutMatcherDeprecation = map[string]bool{}
}

type shimMatcher struct {
	matcher DeprecatedOmegaMatcher
	message string
}

func (m *shimMatcher) Match(actual interface{}) (success bool, err error) {
	success, m.message, err = m.matcher.Match(actual)
	return success, err
}

func (m *shimMatcher) FailureMessage(actual interface{}) (message string) {
	return m.message
}

func (m *shimMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return m.message
}

func shimIfNecessary(matcher interface{}) OmegaMatcher {
	switch m := matcher.(type) {
	case OmegaMatcher:
		return m
	case DeprecatedOmegaMatcher:
		matcherType := fmt.Sprintf("%T", matcher)
		if !didWarnAboutMatcherDeprecation[matcherType] {
			fmt.Printf("\nGOMEGA DEPRECATION WARNING\nYou are using a custom matcher of type:\n    %s\nthat conforms to a deprecated matcher interface.\nPlease update your matcher, the old style will be removed in Gomega 1.0.\n", matcherType)
			didWarnAboutMatcherDeprecation[matcherType] = true
		}
		return &shimMatcher{matcher: m}
	default:
		panic(fmt.Sprintf("Not a valid matcher:\n%v", matcher))
	}
}
