package matchers_test

import (
	"math"

	"github.com/onsi/gomega/matchers/internal/miter"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/matchers"
)

var _ = Describe("HaveLen", func() {
	When("passed a supported type", func() {
		It("should do the right thing", func() {
			Expect("").Should(HaveLen(0))
			Expect("AA").Should(HaveLen(2))

			Expect([0]int{}).Should(HaveLen(0))
			Expect([2]int{1, 2}).Should(HaveLen(2))

			Expect([]int{}).Should(HaveLen(0))
			Expect([]int{1, 2, 3}).Should(HaveLen(3))

			Expect(map[string]int{}).Should(HaveLen(0))
			Expect(map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}).Should(HaveLen(4))

			c := make(chan bool, 3)
			Expect(c).Should(HaveLen(0))
			c <- true
			c <- true
			Expect(c).Should(HaveLen(2))
		})
	})

	When("passed a correctly typed nil", func() {
		It("should operate successfully on the passed in value", func() {
			var nilSlice []int
			Expect(nilSlice).Should(HaveLen(0))

			var nilMap map[int]string
			Expect(nilMap).Should(HaveLen(0))
		})
	})

	When("passed an unsupported type", func() {
		It("should error", func() {
			success, err := (&HaveLenMatcher{Count: 0}).Match(0)
			Expect(success).Should(BeFalse())
			Expect(err).Should(HaveOccurred())

			success, err = (&HaveLenMatcher{Count: 0}).Match(nil)
			Expect(success).Should(BeFalse())
			Expect(err).Should(HaveOccurred())
		})
	})

	When("passed a matcher", func() {
		It("should evaluate strings, arrays, slices, maps, channels with sub-matchers", func() {
			Expect("").Should(HaveLen(BeNumerically("==", 0)))
			Expect("AA").Should(HaveLen(BeNumerically("==", 2)))
			Expect("AA").Should(HaveLen(BeNumerically(">", 1)))
			Expect("AA").ShouldNot(HaveLen(BeNumerically(">", 2)))

			Expect([0]int{}).Should(HaveLen(BeNumerically("==", 0)))
			Expect([2]int{1, 2}).Should(HaveLen(BeNumerically("==", 2)))
			Expect([2]int{1, 2}).Should(HaveLen(BeNumerically("<", 5)))
			Expect([2]int{1, 2}).ShouldNot(HaveLen(BeNumerically("==", 0)))

			Expect([]int{}).Should(HaveLen(BeNumerically("==", 0)))
			Expect([]int{1, 2, 3}).Should(HaveLen(BeNumerically("==", 3)))
			Expect([]int{1, 2, 3}).Should(HaveLen(BeNumerically(">", 2)))
			Expect([]int{1, 2, 3}).Should(HaveLen(Equal(3)))
			Expect([]int{1, 2, 3}).ShouldNot(HaveLen(BeNumerically("<", 2)))

			Expect(map[string]int{}).Should(HaveLen(BeNumerically("==", 0)))
			Expect(map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}).Should(HaveLen(BeNumerically("==", 4)))
			Expect(map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}).Should(HaveLen(BeNumerically(">=", 4)))
			Expect(map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}).ShouldNot(HaveLen(BeNumerically("==", 0)))

			c := make(chan bool, 3)
			Expect(c).Should(HaveLen(BeNumerically("==", 0)))
			c <- true
			c <- true
			Expect(c).Should(HaveLen(BeNumerically("==", 2)))
			Expect(c).Should(HaveLen(BeNumerically(">", 1)))
			Expect(c).ShouldNot(HaveLen(BeNumerically("==", 0)))
		})

		It("should report the length and the sub-matcher's failure message", func() {
			m := HaveLen(BeNumerically(">", 5))
			Expect(m.Match([]int{1, 2})).To(BeFalse())
			Expect(m.FailureMessage([]int{1, 2})).To(Equal("Length 2 of\n    <[]int | len:2, cap:2>: [1, 2]\nfailed to satisfy matcher:\n    Expected\n        <int>: 2\n    to be >\n        <int>: 5"))

			m = HaveLen(BeNumerically(">", 1))
			Expect(m.Match([]int{1, 2})).To(BeTrue())
			Expect(m.NegatedFailureMessage([]int{1, 2})).To(Equal("Length 2 of\n    <[]int | len:2, cap:2>: [1, 2]\nsatisfied matcher, but should not have:\n    Expected\n        <int>: 2\n    not to be >\n        <int>: 1"))
		})

		It("should propagate sub-matcher errors", func() {
			success, err := HaveLen(BeNumerically("invalid-op", 1)).Match([]int{1})
			Expect(success).Should(BeFalse())
			Expect(err).Should(HaveOccurred())
		})
	})

	When("passed a non-int count", func() {
		It("should support other integer types", func() {
			Expect([]int{1, 2}).Should(HaveLen(int64(2)))
			Expect([]int{1, 2}).Should(HaveLen(uint(2)))
			Expect([]int{1, 2}).Should(HaveLen(int32(2)))
			Expect([]int{1, 2}).Should(HaveLen(uint8(2)))
			Expect([]int{1, 2}).ShouldNot(HaveLen(int64(3)))
		})

		It("should support integral floats, so that untyped float constants keep working", func() {
			const N = 3.0
			Expect([]int{1, 2, 3}).Should(HaveLen(N))
			Expect([]int{1, 2, 3}).Should(HaveLen(3.0))
			Expect([]int{1, 2, 3}).Should(HaveLen(float32(3)))
			Expect([]int{1, 2, 3}).ShouldNot(HaveLen(2.0))
		})

		It("should report the plain-integer failure messages", func() {
			Expect(HaveLen(2).FailureMessage([]int{1})).To(Equal("Expected\n    <[]int | len:1, cap:1>: [1]\nto have length 2"))
			Expect(HaveLen(2).NegatedFailureMessage([]int{1, 2})).To(Equal("Expected\n    <[]int | len:2, cap:2>: [1, 2]\nnot to have length 2"))
			Expect(HaveLen(int64(2)).FailureMessage([]int{1})).To(Equal("Expected\n    <[]int | len:1, cap:1>: [1]\nto have length 2"))
			Expect(HaveLen(2.0).FailureMessage([]int{1})).To(Equal("Expected\n    <[]int | len:1, cap:1>: [1]\nto have length 2"))
		})

		It("should error when passed a non-integral float", func() {
			success, err := HaveLen(2.5).Match([]int{1, 2})
			Expect(success).Should(BeFalse())
			Expect(err).Should(MatchError("HaveLen matcher expects an integer or a GomegaMatcher.  Got:\n    <float64>: 2.5"))
		})

		It("should error when passed an integer that does not fit in an int", func() {
			success, err := HaveLen(^uint64(0)).Match([]int{1})
			Expect(success).Should(BeFalse())
			Expect(err).Should(MatchError(ContainSubstring("HaveLen matcher expects an integer or a GomegaMatcher.")))

			success, err = HaveLen(-float64(math.MinInt)).Match([]int{1})
			Expect(success).Should(BeFalse())
			Expect(err).Should(HaveOccurred())

			success, err = HaveLen(math.Inf(1)).Match([]int{1})
			Expect(success).Should(BeFalse())
			Expect(err).Should(HaveOccurred())

			success, err = HaveLen(math.NaN()).Match([]int{1})
			Expect(success).Should(BeFalse())
			Expect(err).Should(HaveOccurred())
		})

		It("should error when passed anything else", func() {
			success, err := HaveLen("invalid").Match([]int{1})
			Expect(success).Should(BeFalse())
			Expect(err).Should(MatchError("HaveLen matcher expects an integer or a GomegaMatcher.  Got:\n    <string>: invalid"))

			success, err = HaveLen(nil).Match([]int{1})
			Expect(success).Should(BeFalse())
			Expect(err).Should(HaveOccurred())
		})
	})

	Context("iterators", func() {
		BeforeEach(func() {
			if !miter.HasIterators() {
				Skip("iterators not available")
			}
		})

		When("passed an iterator type", func() {
			It("should do the right thing", func() {
				Expect(emptyIter).To(HaveLen(0))
				Expect(emptyIter).To(HaveLen(BeNumerically("==", 0)))
				Expect(emptyIter2).To(HaveLen(0))

				Expect(universalIter).To(HaveLen(len(universalElements)))
				Expect(universalIter).To(HaveLen(BeNumerically("==", len(universalElements))))
				Expect(universalIter).To(HaveLen(BeNumerically(">", 0)))
				Expect(universalIter2).To(HaveLen(len(universalElements)))
			})
		})

		When("passed a correctly typed nil", func() {
			It("should operate successfully on the passed in value", func() {
				var nilIter func(func(string) bool)
				Expect(nilIter).Should(HaveLen(0))

				var nilIter2 func(func(int, string) bool)
				Expect(nilIter2).Should(HaveLen(0))
			})
		})
	})
})
