package matchers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("OrMatcher", func() {
	It("works with positive cases", func() {
		Expect(input).To(Or(true1))
		Expect(input).To(Or(true1, true2))
		Expect(input).To(Or(true1, false1))
		Expect(input).To(Or(false1, true2))
		Expect(input).To(Or(true1, true2, true3))
		Expect(input).To(Or(true1, true2, false3))
		Expect(input).To(Or(true1, false2, true3))
		Expect(input).To(Or(false1, true2, true3))
		Expect(input).To(Or(true1, false2, false3))
		Expect(input).To(Or(false1, false2, true3))
	})

	It("works with negative cases", func() {
		Expect(input).ToNot(Or())
		Expect(input).ToNot(Or(false1))
		Expect(input).ToNot(Or(false1, false2))
		Expect(input).ToNot(Or(false1, false2, false3))
	})

	Context("failure messages", func() {
		Context("when match fails", func() {
			It("gives a descriptive message", func() {
				verifyFailureMessage(Or(false1, false2), input,
					"To satisfy at least one of these matchers: [%!s(*matchers.HaveLenMatcher=&{1}) %!s(*matchers.EqualMatcher=&{hip})]")
			})
		})

		Context("when match succeeds, but expected it to fail", func() {
			It("gives a descriptive message", func() {
				verifyFailureMessage(Not(Or(true1, true2)), input, `not to have length 2`)
			})
		})
	})
})
