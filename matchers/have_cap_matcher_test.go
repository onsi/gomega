package matchers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/matchers"
)

var _ = Describe("HaveCap", func() {
	When("passed a supported type", func() {
		It("should do the right thing", func() {
			Expect([0]int{}).Should(HaveCap(0))
			Expect([2]int{1}).Should(HaveCap(2))

			Expect([]int{}).Should(HaveCap(0))
			Expect([]int{1, 2, 3, 4, 5}[:2]).Should(HaveCap(5))
			Expect(make([]int, 0, 5)).Should(HaveCap(5))

			c := make(chan bool, 3)
			Expect(c).Should(HaveCap(3))
			c <- true
			c <- true
			Expect(c).Should(HaveCap(3))

			Expect(make(chan bool)).Should(HaveCap(0))
		})
	})

	When("passed a correctly typed nil", func() {
		It("should operate successfully on the passed in value", func() {
			var nilSlice []int
			Expect(nilSlice).Should(HaveCap(0))

			var nilChan chan int
			Expect(nilChan).Should(HaveCap(0))
		})
	})

	When("passed an unsupported type", func() {
		It("should error", func() {
			success, err := (&HaveCapMatcher{Count: 0}).Match(0)
			Expect(success).Should(BeFalse())
			Expect(err).Should(HaveOccurred())

			success, err = (&HaveCapMatcher{Count: 0}).Match(nil)
			Expect(success).Should(BeFalse())
			Expect(err).Should(HaveOccurred())
		})
	})

	When("passed a matcher", func() {
		It("should evaluate arrays, slices, and channels with sub-matchers", func() {
			Expect([0]int{}).Should(HaveCap(BeNumerically("==", 0)))
			Expect([2]int{1}).Should(HaveCap(BeNumerically("==", 2)))
			Expect([2]int{1}).Should(HaveCap(BeNumerically(">", 1)))
			Expect([2]int{1}).ShouldNot(HaveCap(BeNumerically("<", 2)))

			Expect([]int{}).Should(HaveCap(BeNumerically("==", 0)))
			Expect([]int{1, 2, 3, 4, 5}[:2]).Should(HaveCap(BeNumerically("==", 5)))
			Expect([]int{1, 2, 3, 4, 5}[:2]).Should(HaveCap(BeNumerically(">", 3)))
			Expect([]int{1, 2, 3, 4, 5}[:2]).ShouldNot(HaveCap(BeNumerically("<", 5)))
			Expect(make([]int, 0, 5)).Should(HaveCap(BeNumerically("==", 5)))
			Expect(make([]int, 0, 5)).Should(HaveCap(Equal(5)))
			Expect(make([]int, 0, 5)).ShouldNot(HaveCap(BeNumerically(">", 5)))

			c := make(chan bool, 3)
			Expect(c).Should(HaveCap(BeNumerically("==", 3)))
			Expect(c).Should(HaveCap(BeNumerically(">=", 3)))
			Expect(c).ShouldNot(HaveCap(BeNumerically("==", 0)))
		})

		It("should support non-int integer types", func() {
			Expect(make([]int, 0, 5)).Should(HaveCap(int64(5)))
			Expect(make([]int, 0, 5)).Should(HaveCap(uint(5)))
			Expect(make([]int, 0, 5)).Should(HaveCap(int32(5)))
		})

		It("should report meaningful failure and negated failure messages", func() {
			m := HaveCap(BeNumerically(">", 5))
			Expect(m.FailureMessage(make([]int, 0, 2))).To(ContainSubstring("to have capacity matching"))
			Expect(m.NegatedFailureMessage(make([]int, 0, 2))).To(ContainSubstring("not to have capacity matching"))

			mInt := HaveCap(5)
			Expect(mInt.FailureMessage(make([]int, 0, 2))).To(Equal("Expected\n    <[]int | len:0, cap:2>: []\nto have capacity 5"))
			Expect(mInt.NegatedFailureMessage(make([]int, 0, 5))).To(Equal("Expected\n    <[]int | len:0, cap:5>: []\nnot to have capacity 5"))
		})

		It("should propagate sub-matcher errors", func() {
			success, err := HaveCap(BeNumerically("invalid-op", 1)).Match([]int{1})
			Expect(success).Should(BeFalse())
			Expect(err).Should(HaveOccurred())
		})

		It("should error when passed an invalid count argument", func() {
			success, err := HaveCap("invalid").Match([]int{1})
			Expect(success).Should(BeFalse())
			Expect(err).Should(MatchError(ContainSubstring("HaveCap matcher expects an int or GomegaMatcher")))
		})
	})
})
