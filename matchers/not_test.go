package matchers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NotMatcher", func() {
	Context("basic examples", func() {
		It("works", func() {
			Expect(input).To(Not(false1))
			Expect(input).To(Not(Not(true2)))
			Expect(input).ToNot(Not(true3))
			Expect(input).ToNot(Not(Not(false1)))
			Expect(input).To(Not(Not(Not(false2))))
		})
	})

	Context("De Morgan's laws", func() {
		It("~(A && B) == ~A || ~B", func() {
			Expect(input).To(Not(And(false1, false2)))
			Expect(input).To(Or(Not(false1), Not(false2)))
		})
		It("~(A || B) == ~A && ~B", func() {
			Expect(input).To(Not(Or(false1, false2)))
			Expect(input).To(And(Not(false1), Not(false2)))
		})
	})

	Context("failure messages are opposite of original matchers' failure messages", func() {
		Context("when match fails", func() {
			It("gives a descriptive message", func() {
				verifyFailureMessage(Not(HaveLen(2)), input, "not to have length 2")
			})
		})

		Context("when match succeeds, but expected it to fail", func() {
			It("gives a descriptive message", func() {
				verifyFailureMessage(Not(Not(HaveLen(3))), input, "to have length 3")
			})
		})
	})
})
