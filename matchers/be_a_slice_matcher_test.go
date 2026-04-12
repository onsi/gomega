package matchers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/matchers"
)

var _ = Describe("BeASliceMatcher", func() {
	When("passed a slice", func() {
		It("should succeed", func() {
			Expect([]int{1, 2, 3}).Should(BeASlice())
			Expect([]string{}).Should(BeASlice())
			Expect([]any{1, "two", 3.0}).Should(BeASlice())
			Expect([]byte("hello")).Should(BeASlice())
		})
	})

	When("passed a nil slice", func() {
		It("should succeed", func() {
			var s []int
			Expect(s).Should(BeASlice())
		})
	})

	When("passed a non-slice", func() {
		It("should fail", func() {
			Expect("hello").ShouldNot(BeASlice())
			Expect(42).ShouldNot(BeASlice())
			Expect(true).ShouldNot(BeASlice())
			Expect(map[string]int{"a": 1}).ShouldNot(BeASlice())
			Expect([3]int{1, 2, 3}).ShouldNot(BeASlice())
			Expect(struct{}{}).ShouldNot(BeASlice())
		})
	})

	When("passed nil", func() {
		It("should error", func() {
			success, err := (&BeASliceMatcher{}).Match(nil)
			Expect(success).Should(BeFalse())
			Expect(err).Should(HaveOccurred())
		})
	})

	It("should produce appropriate failure messages", func() {
		matcher := &BeASliceMatcher{}

		matcher.Match("hello")
		Expect(matcher.FailureMessage("hello")).Should(ContainSubstring("to be a slice"))
		Expect(matcher.NegatedFailureMessage("hello")).Should(ContainSubstring("not to be a slice"))
	})
})
