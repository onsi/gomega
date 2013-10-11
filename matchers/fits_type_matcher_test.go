package matchers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/matchers"
)

type ToEmbed struct {
}

type Embeder struct {
	ToEmbed
}

var _ = Describe("FitsType", func() {
	Context("When asserting equality between types", func() {
		It("should do the right thing", func() {
			Ω(0).Should(FitTypeOf(0))         // Same values
			Ω(5).Should(FitTypeOf(-1))        // different values same type
			Ω("foo").Should(FitTypeOf("bar")) // different values same type
		})
	})

	Context("When asserting nil values", func() {
		It("should error", func() {
			success, _, err := (&FitsTypeMatcher{Expected: nil}).Match(nil)

			Ω(success).Should(BeFalse())
			Ω(err).Should(HaveOccured())
		})
	})
})
