package matchers

import (
	"reflect"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

// HaveValue applies the given matcher to the value of actual, optionally and
// repeatedly dereferencing pointers or taking the concrete value of interfaces.
// Thus, the matcher will always be applied to non-pointer and non-interface
// values only. HaveValue will fail with an error if a pointer or interface is
// nil. It will also fail for more than 31 pointer or interface dereferences to
// guard against mistakenly applying it to arbitrarily deep linked pointers.
//
// HaveValue differs from gstruct.PointTo in that it does not expect actual to
// be a pointer (as gstruct.PointTo does) but instead also accepts non-pointer
// and even interface values.
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
	Matcher types.GomegaMatcher // the matcher to apply to the "resolved" actual value.
	failure string              // failure message, if any.
}

func (m *HaveValueMatcher) Match(actual interface{}) (bool, error) {
	val := reflect.ValueOf(actual)
	for allowedIndirs := 32; allowedIndirs > 0; allowedIndirs-- {
		// return an error if value isn't valid. Please note that we cannot
		// check for nil here, as we might not deal with a pointer or interface
		// at this point.
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
