package matchers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BeZero", func() {
	It("should succeed if the passed in object is the zero value for its type", func() {

		Ω("").Should(BeZero())
		Ω(" ").ShouldNot(BeZero())

		Ω(0).Should(BeZero())
		Ω(1).ShouldNot(BeZero())

		Ω(0.0).Should(BeZero())
		Ω(0.1).ShouldNot(BeZero())

		Ω(myCustomType{}).Should(BeZero())
		Ω(myCustomType{s: "a"}).ShouldNot(BeZero())
	})

	It("should succeed when passed nil or a typed nil", func() {
		Ω(nil).Should(BeZero())

		var nilArray []int
		Ω(nilArray).Should(BeZero())
		Ω([]int{1}).ShouldNot(BeZero())

		var nilHash map[string]int
		Ω(nilHash).Should(BeZero())
		Ω(map[string]int{"a": 1}).ShouldNot(BeZero())
	})
})
