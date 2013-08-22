package matcher_tests

import (
	"errors"
	. "github.com/onsi/godescribe"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/matchers"
)

func init() {
	Describe("Equal", func() {
		Context("when asserting that nil equals nil", func() {
			It("should error", func() {
				matcher := &EqualMatcher{Expected: nil}
				success, _, _ := matcher.Match(nil)

				Ω(success).Should(BeFalse())
				Ω(nil).ShouldNot(Equal(nil))
			})
		})

		Context("When asserting equality of numbers", func() {
			It("should do the right thing", func() {
				Ω(5).Should(Equal(5))
				Ω(5.0).Should(Equal(5))
				Ω(5).Should(Equal(5.0))

				Ω(5).ShouldNot(Equal(3.0))
				Ω(3).ShouldNot(Equal(5))

				Ω(0).Should(Equal(0))
			})
		})
	})

	Describe("BeTrue", func() {
		It("should handle true and false correctly", func() {
			Ω(true).Should(BeTrue())
			Ω(false).ShouldNot(BeTrue())
		})

		It("should only support booleans", func() {
			success, _, err := (&TrueMatcher{}).Match("foo")
			Ω(success).Should(BeFalse())
			Ω(err).Should(HaveOccured())
		})
	})

	Describe("BeFalse", func() {
		It("should handle true and false correctly", func() {
			Ω(true).ShouldNot(BeFalse())
			Ω(false).Should(BeFalse())
		})

		It("should only support booleans", func() {
			success, _, err := (&FalseMatcher{}).Match("foo")
			Ω(success).Should(BeFalse())
			Ω(err).Should(HaveOccured())
		})
	})

	Describe("HaveOccured", func() {
		It("should succeed if matching an error", func() {
			Ω(errors.New("Foo")).ShouldNot(HaveOccured())
		})

		It("should not succed with nil", func() {
			Ω(nil).Should(HaveOccured())
		})

		It("should only support errors and nil", func() {
			success, _, err := (&HaveOccuredMatcher{}).Match("foo")
			Ω(success).Should(BeFalse())
			Ω(err).Should(HaveOccured())
			success, _, err = (&HaveOccuredMatcher{}).Match("")
			Ω(success).Should(BeFalse())
			Ω(err).Should(HaveOccured())
		})
	})
}
