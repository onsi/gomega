package matchers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("WithTransformMatcher", func() {

	var plus1 = func(i interface{}) interface{} { return i.(int) + 1 }

	It("works with positive cases", func() {
		Expect(1).To(WithTransform(plus1, Equal(2)))
		Expect(1).To(WithTransform(plus1, WithTransform(plus1, Equal(3))))
		Expect(1).To(WithTransform(plus1, And(Equal(2), BeNumerically(">", 1))))
	})

	It("works with negative cases", func() {
		Expect(1).ToNot(WithTransform(plus1, Equal(3)))
		Expect(1).ToNot(WithTransform(plus1, WithTransform(plus1, Equal(2))))
	})

	Context("failure messages", func() {
		Context("when match fails", func() {
			It("gives a descriptive message", func() {
				m := WithTransform(plus1, Equal(3))
				Expect(m.Match(1)).To(BeFalse())
				Expect(m.FailureMessage(input)).To(Equal("Expected\n    <int>: 2\nto equal\n    <int>: 3"))
			})
		})

		Context("when match succeeds, but expected it to fail", func() {
			It("gives a descriptive message", func() {
				m := Not(WithTransform(plus1, Equal(3)))
				Expect(m.Match(2)).To(BeFalse())
				Expect(m.FailureMessage(input)).To(Equal("Expected\n    <int>: 3\nnot to equal\n    <int>: 3"))
			})
		})
	})
})
