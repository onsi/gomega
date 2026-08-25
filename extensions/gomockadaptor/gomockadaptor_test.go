package gomockadaptor_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/extensions/gomockadaptor"
	"github.com/onsi/gomega/types"
)

// gomockMatcher mirrors gomock's Matcher interface so we can assert, without
// depending on gomock, that an Adaptor can be used wherever one is expected.
type gomockMatcher interface {
	Matches(x any) bool
	String() string
}

var _ gomockMatcher = gomockadaptor.Match(Equal(1))

// erroringMatcher is a GomegaMatcher whose Match always returns an error.
type erroringMatcher struct{}

func (erroringMatcher) Match(actual any) (bool, error) {
	return false, errors.New("boom")
}
func (erroringMatcher) FailureMessage(actual any) string        { return "should not be used" }
func (erroringMatcher) NegatedFailureMessage(actual any) string { return "should not be used" }

var _ types.GomegaMatcher = erroringMatcher{}

var _ = Describe("Gomockadaptor", func() {
	It("reports a match when the wrapped matcher succeeds", func() {
		adaptor := gomockadaptor.Match(Equal(3))
		Expect(adaptor.Matches(3)).To(BeTrue())
		Expect(adaptor.String()).To(BeEmpty())
	})

	It("reports no match when the wrapped matcher fails, and surfaces the failure message", func() {
		adaptor := gomockadaptor.Match(Equal(3))
		Expect(adaptor.Matches(4)).To(BeFalse())
		Expect(adaptor.String()).To(ContainSubstring("to equal"))
	})

	It("treats a matcher error as a non-match and surfaces the error", func() {
		adaptor := gomockadaptor.Match(erroringMatcher{})
		Expect(adaptor.Matches("anything")).To(BeFalse())
		Expect(adaptor.String()).To(Equal("boom"))
	})

	It("clears the failure message between calls", func() {
		adaptor := gomockadaptor.Match(Equal(3))
		Expect(adaptor.Matches(4)).To(BeFalse())
		Expect(adaptor.String()).NotTo(BeEmpty())
		Expect(adaptor.Matches(3)).To(BeTrue())
		Expect(adaptor.String()).To(BeEmpty())
	})

	It("works with composed matchers", func() {
		adaptor := gomockadaptor.Match(HaveLen(2))
		Expect(adaptor.Matches([]string{"a", "b"})).To(BeTrue())
		Expect(adaptor.Matches([]string{"a"})).To(BeFalse())
	})
})
