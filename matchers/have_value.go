package matchers

import (
	"reflect"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

// HavePoint applies the given matcher to the resulting value after optionally
// resolving actual to the value it points to or its interface value, repeatedly
// as necessary. It fails if a pointer or interface is nil. In contrast to
// gstruct.PointTo, HaveValue does not expect actual to be a pointer but instead
// also accepts non-pointer and interface values.
//
//   actual := 42
//   Expect(actual).To(HaveValue(42))
//   Expect(&actual).To(HaveValue(42))
func HaveValue(matcher types.GomegaMatcher) types.GomegaMatcher {
	return &HaveValueMatcher{
		Matcher: matcher,
	}
}

type HaveValueMatcher struct {
	Matcher types.GomegaMatcher // given matcher to apply to resolved value.
	failure string              // failure message, if any.
}

func (m *HaveValueMatcher) Match(actual interface{}) (bool, error) {
	val := reflect.ValueOf(actual)
	for allowedIndirs := 32; allowedIndirs > 0; allowedIndirs-- {
		// return an error if value isn't valid.
		if !val.IsValid() {
			m.failure = format.Message(
				actual, "not to be <nil>")
			return false, nil
		}
		switch val.Kind() {
		case reflect.Ptr, reflect.Interface:
			// resolve pointers and interfaces to their values, then rinse and
			// repeat.
			if val.IsNil() {
				m.failure = format.Message(
					actual, "not to be <nil>")
				return false, nil
			}
			val = val.Elem()
			continue
		default:
			// forward the final value to the specified matcher.
			elem := val.Interface()
			match, err := m.Matcher.Match(elem)
			if !match {
				m.failure = m.Matcher.FailureMessage(elem)
			}
			return match, err
		}
	}
	// too many indirections: extreme star gazing, indeed...?
	m.failure = format.Message(actual, "indirecting too many times")
	return false, nil
}

func (m *HaveValueMatcher) FailureMessage(_ interface{}) (message string) {
	return m.failure
}

func (m *HaveValueMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return m.Matcher.NegatedFailureMessage(actual)
}
