package gomega

import (
	"github.com/onsi/gomega/matchers"
)

//Equal uses reflect.DeepEqual to compare actual with expected
//It is an error for both actual and expected to be nil.  Use BeNil() instead.
func Equal(expected interface{}) OmegaMatcher {
	return &matchers.EqualMatcher{
		Expected: expected,
	}
}

//BeNil succeeds if actual is nil
func BeNil() OmegaMatcher {
	return &matchers.BeNilMatcher{}
}

//BeTrue succeeds if actual is true
func BeTrue() OmegaMatcher {
	return &matchers.BeTrueMatcher{}
}

//BeFalse succeeds if actual is false
func BeFalse() OmegaMatcher {
	return &matchers.BeFalseMatcher{}
}

//HaveOccured succeeds if actual is a non-nil error
//The typical Go error checking pattern looks like:
//
//  err := SomethingThatMightFail()
//  Ω(err).ShouldNot(HaveOccured())
func HaveOccured() OmegaMatcher {
	return &matchers.HaveOccuredMatcher{}
}

//MatchRegexp succeeds if actual is a string or stringer that matches the
//passed-in regexp.  Optional arguments can be provided to construct a regexp
//via fmt.Sprintf().
func MatchRegexp(regexp string, args ...interface{}) OmegaMatcher {
	return &matchers.MatchRegexpMatcher{
		Regexp: regexp,
		Args:   args,
	}
}

//ContainSubstring succeeds if actual is a string or stringer that contains the
//passed-in regexp.  Optional arguments can be provided to construct the substring
//via fmt.Sprintf().
func ContainSubstring(substr string, args ...interface{}) OmegaMatcher {
	return &matchers.ContainSubstringMatcher{
		Substr: substr,
		Args:   args,
	}
}

//BeEmpty succeeds if actual is empty.  Actual must be of type string, array, map, chan, or slice.
func BeEmpty() OmegaMatcher {
	return &matchers.BeEmptyMatcher{}
}

//HaveLen succeeds if actual has the passed-in length.  Actual must be of type string, array, map, chan, or slice.
func HaveLen(count int) OmegaMatcher {
	return &matchers.HaveLenMatcher{
		Count: count,
	}
}

//BeZero succeeds if actual is the zero value for its type or if actual is nil.
func BeZero() OmegaMatcher {
	return &matchers.BeZeroMatcher{}
}

//ContainElement succeeds if actual contains the passed in element.
//By default ContainElement() uses Equal() to perform the match, however a
//matcher can be passed in instead:
//
//  Ω([]string{"Foo", "FooBar"}).Should(ContainElement(ContainSubstring("Bar")))
//
//Actual must be an array, slice or map.
//For maps, containElement searches through the map's values.
func ContainElement(element interface{}) OmegaMatcher {
	return &matchers.ContainElementMatcher{
		Element: element,
	}
}

//HaveKey succeeds if actual is a map with the passed in key.
//By default HaveKey uses Equal() to perform the match, however a
//matcher can be passed in instead:
//
//  Ω(map[string]string{"Foo": "Bar", "BazFoo": "Duck"}).Should(HaveKey(MatchRegexp(`.+Foo$`)))
func HaveKey(key interface{}) OmegaMatcher {
	return &matchers.HaveKeyMatcher{
		Key: key,
	}
}

//BeNumerically performs numerical assertions in a type-agnostic way.
//Actual and expected should be numbers, though the specific type of
//number is irrelevant (floa32, float64, uint8, etc...).
//
//There are six, self-explanatory, supported comparators:
//
//  Ω(1.0).Should(BeNumerically("==", 1))
//  Ω(1.0).Should(BeNumerically("~", 0.999, 0.01))
//  Ω(1.0).Should(BeNumerically(">", 0.9))
//  Ω(1.0).Should(BeNumerically(">=", 1.0))
//  Ω(1.0).Should(BeNumerically("<", 3))
//  Ω(1.0).Should(BeNumerically("<=", 1.0))
func BeNumerically(comparator string, compareTo ...interface{}) OmegaMatcher {
	return &matchers.BeNumericallyMatcher{
		Comparator: comparator,
		CompareTo:  compareTo,
	}
}

//BeAssignableTo succeeds if actual is assignable to the type of actual.
//It will return an error when one of the values is nil.
//
//	Ω(0).Should(BeAssignableTo(0))         // Same values
//	Ω(5).Should(BeAssignableTo(-1))        // different values same type
//	Ω("foo").Should(BeAssignableTo("bar")) // different values same type
func BeAssignableTo(expected interface{}) OmegaMatcher {
	return &matchers.AssignableToMatcher{
		Expected: expected,
	}
}

//Panic succeeds if actual is a function that, when invoked, panics.
//Actual must be a function that takes no arguments and returns no results.
func Panic() OmegaMatcher {
	return &matchers.PanicMatcher{}
}
