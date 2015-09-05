package matchers_test

import (
	"errors"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ImplementsMatcherMatcher", func() {

	var errorInterface = reflect.TypeOf((*error)(nil)).Elem()

	Context("when actual implements the given interface", func() {
		It("matches", func() {
			Expect(errors.New("hi")).To(Implements(errorInterface))
		})
	})

	Context("when actual does not implement the given interface", func() {
		It("does not match", func() {
			Expect(nil).ToNot(Implements(errorInterface))
			Expect("hi").ToNot(Implements(errorInterface))
		})
	})

	Context("failure messages", func() {
		var m = Implements(errorInterface)

		Context("when match fails", func() {
			It("gives a descriptive message", func() {
				actual := "hi"
				Expect(m.Match(actual)).To(BeFalse())
				Expect(m.FailureMessage(actual)).To(Equal("Expected\n    <string>: hi\nto implement error"))
			})
		})

		Context("when match succeeds, but expected it to fail", func() {
			It("gives a descriptive message", func() {
				actual := errors.New("hi")
				Expect(m.Match(actual)).To(BeTrue())
				Expect(m.NegatedFailureMessage(actual)).To(MatchRegexp(
					"Expected\n    <*errors.errorString | 0x[^>]+>: {s: \"hi\"}\nto not implement error"))
			})
		})
	})

	Context("Invalid interface provided", func() {
		It("panics", func() {
			Expect(func() { Implements(nil) }).To(Panic(), "nil interface type")
			Expect(func() { Implements(reflect.TypeOf(123)) }).To(Panic(), "not an interface type")
		})
	})
})
