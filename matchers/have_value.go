package matchers

import (
	"errors"
	"reflect"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

const maxIndirections = 31

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
	Matcher        types.GomegaMatcher // the matcher to apply to the "resolved" actual value.
	resolvedActual interface{}         // the ("resolved") value.
}

func (m *HaveValueMatcher) Match(actual interface{}) (bool, error) {
	val := reflect.ValueOf(actual)
	for allowedIndirs := maxIndirections; allowedIndirs > 0; allowedIndirs-- {
		// return an error if value isn't valid. Please note that we cannot
		// check for nil here, as we might not deal with a pointer or interface
		// at this point.
		if !val.IsValid() {
			return false, errors.New(format.Message(
				actual, "not to be <nil>"))
		}
		switch val.Kind() {
		case reflect.Ptr, reflect.Interface:
			// resolve pointers and interfaces to their values, then rinse and
			// repeat.
			if val.IsNil() {
				return false, errors.New(format.Message(
					actual, "not to be <nil>"))
			}
			val = val.Elem()
			continue
		default:
			// forward the final value to the specified matcher.
			m.resolvedActual = val.Interface()
			return m.Matcher.Match(m.resolvedActual)
		}
	}
	// too many indirections: extreme star gazing, indeed...?
	return false, errors.New(format.Message(actual, "too many indirections"))
}

func (m *HaveValueMatcher) FailureMessage(_ interface{}) (message string) {
	return m.Matcher.FailureMessage(m.resolvedActual)
}

func (m *HaveValueMatcher) NegatedFailureMessage(_ interface{}) (message string) {
	return m.Matcher.NegatedFailureMessage(m.resolvedActual)
}
