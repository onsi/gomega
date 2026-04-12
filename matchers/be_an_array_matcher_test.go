package matchers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/matchers"
)

var _ = Describe("BeAnArrayMatcher", func() {
	When("passed an array", func() {
		It("should succeed", func() {
			Expect([3]int{1, 2, 3}).Should(BeAnArray())
			Expect([0]string{}).Should(BeAnArray())
			Expect([2]any{1, "two"}).Should(BeAnArray())
			Expect([4]byte{1, 2, 3, 4}).Should(BeAnArray())
		})
	})

	When("passed a non-array", func() {
		It("should fail", func() {
			Expect("hello").ShouldNot(BeAnArray())
			Expect(42).ShouldNot(BeAnArray())
			Expect(true).ShouldNot(BeAnArray())
			Expect(map[string]int{"a": 1}).ShouldNot(BeAnArray())
			Expect([]int{1, 2, 3}).ShouldNot(BeAnArray())
			Expect(struct{}{}).ShouldNot(BeAnArray())
		})
	})

	When("passed nil", func() {
		It("should error", func() {
			success, err := (&BeAnArrayMatcher{}).Match(nil)
			Expect(success).Should(BeFalse())
			Expect(err).Should(HaveOccurred())
		})
	})

	It("should produce appropriate failure messages", func() {
		matcher := &BeAnArrayMatcher{}

		matcher.Match("hello")
		Expect(matcher.FailureMessage("hello")).Should(ContainSubstring("to be an array"))
		Expect(matcher.NegatedFailureMessage("hello")).Should(ContainSubstring("not to be an array"))
	})
})
