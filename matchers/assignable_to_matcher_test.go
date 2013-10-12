package matchers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/matchers"
)

var _ = Describe("AssignableTo", func() {
	Context("When asserting equality between types", func() {
		It("should do the right thing", func() {
			Ω(0).Should(FitTypeOf(0))
			Ω(5).Should(FitTypeOf(-1))
			Ω("foo").Should(FitTypeOf("bar"))
		})
	})

	Context("When asserting nil values", func() {
		It("should error", func() {
			success, _, err := (&AssignableToMatcher{Expected: nil}).Match(nil)

			Ω(success).Should(BeFalse())
			Ω(err).Should(HaveOccured())
		})
	})
})
