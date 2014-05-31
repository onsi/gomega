package matchers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/matchers"
)

var _ = Describe("MatchArrayOrSlice", func() {
	Context("when expected is not an array or slice", func() {
		It("should error", func() {
			success, err := (&MatchArrayOrSliceMatcher{ExpectedArrayOrSlice: 123}).Match([]int{123})
			Ω(success).Should(BeFalse())
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("when expected is an array or slice", func() {
		Context("when actual is not an array or slice", func() {
			It("should fail the assertion", func() {
				Ω(123).ShouldNot(MatchArrayOrSlice([]int{123}))
			})
		})

		Context("when actual and expected have the same type", func() {
			Context("when they have different length", func() {
				It("should fail the assertion", func() {
					Ω([]int{}).ShouldNot(MatchArrayOrSlice([]int{123}))
				})
			})

			Context("when they have the same elements with the same frequencies in the same order", func() {
				It("should pass the assertion", func() {
					Ω([]int{1, 2, 1}).Should(MatchArrayOrSlice([]int{1, 2, 1}))
				})
			})

			Context("when they have the same elements with the same frequencies in a different order", func() {
				It("should pass the assertion", func() {
					Ω([]int{1, 2, 1}).Should(MatchArrayOrSlice([]int{2, 1, 1}))
				})
			})

			Context("when they have the same elements with different frequencies", func() {
				It("should fail the assertion", func() {
					Ω([]int{1, 2, 1}).ShouldNot(MatchArrayOrSlice([]int{2, 2, 1}))
				})
			})
		})

		Context("when actual and expected have different types", func() {
			Context("when the difference is merely superficial", func() {
				It("should pass the assertion", func() {
					Ω([]int{123}).Should(MatchArrayOrSlice([]interface{}{123}))
					Ω([]int{}).Should(MatchArrayOrSlice([]string{}))
				})
			})

			Context("when the difference is material", func() {
				It("should fail the assertion", func() {
					Ω([]int{123}).ShouldNot(MatchArrayOrSlice([]interface{}{"123"}))
				})
			})
		})
	})

	// This test has entirely negligible effect on build time.
	It("is very fast for arrays or slices of unequal length", func() {
		var veryLongArray [999999]byte
		var evenLongerArray [1000000]byte
		Ω(veryLongArray).ShouldNot(MatchArrayOrSlice(evenLongerArray))
	})

	// This test adds less than 2% to the total time to run `ginkgo -r` in the 'matchers' directory.
	It("is not too slow even for large arrays or slices of the same length", func() {
		var longArray [9999]byte
		var otherLongArray [9999]byte
		longArray[0] = 1
		otherLongArray[9998] = 1
		Ω(longArray).Should(MatchArrayOrSlice(otherLongArray))
	})
})
