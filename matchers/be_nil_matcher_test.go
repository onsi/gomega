package matchers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BeNil", func() {
	It("should succeed when passed nil", func() {
		Ω(nil).Should(BeNil())
	})

	It("should not succeed when not passed nil", func() {
		Ω(0).ShouldNot(BeNil())
		Ω(false).ShouldNot(BeNil())
		Ω("").ShouldNot(BeNil())
	})
})
