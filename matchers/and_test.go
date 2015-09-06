package matchers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

// sample data
var (
	// example input
	input = "hi"
	// some matchers that succeed against the input
	true1 = HaveLen(2)
	true2 = Equal("hi")
	true3 = MatchRegexp("hi")
	// some matchers that fail against the input.
	false1 = HaveLen(1)
	false2 = Equal("hip")
	false3 = MatchRegexp("hope")
)

// verifyFailureMessage expects the matcher to fail with the given input, and verifies the failure message.
func verifyFailureMessage(m types.GomegaMatcher, input string, expectedFailureMsgFragment string) {
	Expect(m.Match(input)).To(BeFalse())
	Expect(m.FailureMessage(input)).To(Equal(
		"Expected\n    <string>: " + input + "\n" + expectedFailureMsgFragment))
}

var _ = Describe("AndMatcher", func() {
	It("works with positive cases", func() {
		Expect(input).To(And())
		Expect(input).To(And(true1))
		Expect(input).To(And(true1, true2))
		Expect(input).To(And(true1, true2, true3))
	})

	It("works with negative cases", func() {
		Expect(input).ToNot(And(false1, false2))
		Expect(input).ToNot(And(true1, true2, false3))
		Expect(input).ToNot(And(true1, false2, false3))
		Expect(input).ToNot(And(false1, true1, true2))
	})

	Context("failure messages", func() {
		Context("when match fails", func() {
			It("gives a descriptive message", func() {
				verifyFailureMessage(And(false1, true1), input, "to have length 1")
				verifyFailureMessage(And(true1, false2), input, "to equal\n    <string>: hip")
				verifyFailureMessage(And(true1, true2, false3), input, "to match regular expression\n    <string>: hope")
			})
		})

		Context("when match succeeds, but expected it to fail", func() {
			It("gives a descriptive message", func() {
				verifyFailureMessage(Not(And(true1, true2)), input,
					`To not satisfy all of these matchers: [%!s(*matchers.HaveLenMatcher=&{2}) %!s(*matchers.EqualMatcher=&{hi})]`)
			})
		})
	})
})
