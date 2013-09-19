package matchers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/matchers"
)

var _ = Describe("Panic", func() {
	Context("when passed something that's not a function that takes zero arguments and returns nothing", func() {
		It("should error", func() {
			success, _, err := (&PanicMatcher{}).Match("foo")
			Ω(success).Should(BeFalse())
			Ω(err).Should(HaveOccured())

			success, _, err = (&PanicMatcher{}).Match(nil)
			Ω(success).Should(BeFalse())
			Ω(err).Should(HaveOccured())

			success, _, err = (&PanicMatcher{}).Match(func(foo string) {})
			Ω(success).Should(BeFalse())
			Ω(err).Should(HaveOccured())

			success, _, err = (&PanicMatcher{}).Match(func() string { return "bar" })
			Ω(success).Should(BeFalse())
			Ω(err).Should(HaveOccured())
		})
	})

	Context("when passed a function of the correct type", func() {
		It("should call the function and pass if the function panics", func() {
			Ω(func() { panic("ack!") }).Should(Panic())
			Ω(func() {}).ShouldNot(Panic())
		})
	})
})
