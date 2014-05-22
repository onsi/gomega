package gomega

import (
	"github.com/onsi/gomega/matchers"
)

//Equal uses reflect.DeepEqual to compare actual with expected.  Equal is strict about
//types when performing comparisons.
//It is an error for both actual and expected to be nil.  Use BeNil() instead.
func Equal(expected interface{}) OmegaMatcher {
	return &matchers.EqualMatcher{
		Expected: expected,
	}
}

//BeEquivalentTo is more lax than Equal, allowing equality between different types.
//This is done by converting actual to have the type of expected before
//attempting equality with reflect.DeepEqual.
//It is an error for actual and expected to be nil.  Use BeNil() instead.
func BeEquivalentTo(expected interface{}) OmegaMatcher {
	return &matchers.BeEquivalentToMatcher{
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

//HaveOccurred succeeds if actual is a non-nil error
//The typical Go error checking pattern looks like:
//    err := SomethingThatMightFail()
//    Ω(err).ShouldNot(HaveOccurred())
func HaveOccurred() OmegaMatcher {
	return &matchers.HaveOccurredMatcher{}
}

//MatchError succeeds if actual is a non-nil error that matches the passed in string/error.
//
//These are valid use-cases:
//  Ω(err).Should(MatchError("an error")) //asserts that err.Error() == "an error"
//  Ω(err).Should(MatchError(SomeError)) //asserts that err == SomeError (via reflect.DeepEqual)
//
//It is an error for err to be nil or an object that does not implement the Error interface
func MatchError(expected interface{}) OmegaMatcher {
	return &matchers.MatchErrorMatcher{
		Expected: expected,
	}
}

//BeClosed succeeds if actual is a closed channel.
//It is an error to pass a non-channel to BeClosed, it is also an error to pass nil
//
//In order to check whether or not the channel is closed, Gomega must try to read from the channel
//(even in the `ShouldNot(BeClosed())` case).  You should keep this in mind if you wish to make subsequent assertions about
//values coming down the channel.
//
//Also, if you are testing that a *buffered* channel is closed you must first read all values out of the channel before
//asserting that it is closed (it is not possible to detect that a buffered-channel has been closed until all its buffered values are read).
//
//Finally, as a corollary: it is an error to check whether or not a send-only channel is closed.
func BeClosed() OmegaMatcher {
	return &matchers.BeClosedMatcher{}
}

//Receive succeeds if there is a message to be received on actual.
//Actual must be a channel (and cannot be a send-only channel) -- anything else is an error.
//
//Receive returns immediately and never blocks:
//
//- If there is nothing on the channel `c` then Ω(c).Should(Receive()) will fail and Ω(c).ShouldNot(Receive()) will pass.
//
//- If the channel `c` is closed then *both* Ω(c).Should(Receive()) and Ω(c).ShouldNot(Receive()) will error.
//
//- If there is something on the channel `c` ready to be read, then Ω(c).Should(Receive()) will pass and Ω(c).ShouldNot(Receive()) will fail.
//
//If you have a go-routine running in the background that will write to channel `c` you can:
//    Eventually(c).Should(Receive())
//
//This will timeout if nothing gets sent to `c` (you can modify the timeout interval as you normally do with `Eventually`)
//
//A similar use-case is to assert that no go-routine writes to a channel (for a period of time).  You can do this with `Consistently`:
//    Consistently(c).ShouldNot(Receive())
//
//You can pass `Receive` a matcher.  If you do so, it will match the received object against the matcher.  For example:
//    Ω(c).Should(Receive(Equal("foo")))
//
//When given a matcher, `Receive` will always fail if there is nothing to be received on the channel.
//
//Passing Receive a matcher is especially useful when paired with Eventually:
//
//    Eventually(c).Should(Receive(ContainSubstring("bar")))
//
//will repeatedly attempt to pull values out of `c` until a value matching "bar" is received.
//
//Finally, if you want to have a reference to the value *sent* to the channel you can pass the `Receive` matcher a pointer to a variable of the appropriate type:
//    var myThing thing
//    Eventually(thingChan).Should(Receive(&myThing))
//    Ω(myThing.Sprocket).Should(Equal("foo"))
//    Ω(myThing.IsValid()).Should(BeTrue())
func Receive(args ...interface{}) OmegaMatcher {
	var arg interface{}
	if len(args) > 0 {
		arg = args[0]
	}

	return &matchers.ReceiveMatcher{
		Arg: arg,
	}
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

//MatchJSON succeeds if actual is a string or stringer of JSON that matches
//the expected JSON.  The JSONs are decoded and the resulting objects are compared via
//reflect.DeepEqual so things like key-ordering and whitespace shouldn't matter.
func MatchJSON(json interface{}) OmegaMatcher {
	return &matchers.MatchJSONMatcher{
		JSONToMatch: json,
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
//    Ω([]string{"Foo", "FooBar"}).Should(ContainElement(ContainSubstring("Bar")))
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
//    Ω(map[string]string{"Foo": "Bar", "BazFoo": "Duck"}).Should(HaveKey(MatchRegexp(`.+Foo$`)))
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
//    Ω(1.0).Should(BeNumerically("==", 1))
//    Ω(1.0).Should(BeNumerically("~", 0.999, 0.01))
//    Ω(1.0).Should(BeNumerically(">", 0.9))
//    Ω(1.0).Should(BeNumerically(">=", 1.0))
//    Ω(1.0).Should(BeNumerically("<", 3))
//    Ω(1.0).Should(BeNumerically("<=", 1.0))
func BeNumerically(comparator string, compareTo ...interface{}) OmegaMatcher {
	return &matchers.BeNumericallyMatcher{
		Comparator: comparator,
		CompareTo:  compareTo,
	}
}

//BeAssignableToTypeOf succeeds if actual is assignable to the type of expected.
//It will return an error when one of the values is nil.
//	  Ω(0).Should(BeAssignableToTypeOf(0))         // Same values
//	  Ω(5).Should(BeAssignableToTypeOf(-1))        // different values same type
//	  Ω("foo").Should(BeAssignableToTypeOf("bar")) // different values same type
//    Ω(struct{ Foo string }{}).Should(BeAssignableToTypeOf(struct{ Foo string }{}))
func BeAssignableToTypeOf(expected interface{}) OmegaMatcher {
	return &matchers.AssignableToTypeOfMatcher{
		Expected: expected,
	}
}

//Panic succeeds if actual is a function that, when invoked, panics.
//Actual must be a function that takes no arguments and returns no results.
func Panic() OmegaMatcher {
	return &matchers.PanicMatcher{}
}
