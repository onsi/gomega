package matchers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FitsType", func() {
	Context("When asserting equality between types", func() {
		It("should do the right thing", func() {
			Ω(0).Should(FitTypeOf(0))
			Ω(5).Should(FitTypeOf(-1))
			Ω("foo").Should(FitTypeOf("bar"))
		})
	})
})
