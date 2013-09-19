package matchers_test

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/matchers"
)

var _ = Describe("HaveOccured", func() {
	It("should succeed if matching an error", func() {
		Ω(errors.New("Foo")).Should(HaveOccured())
	})

	It("should not succeed with nil", func() {
		Ω(nil).ShouldNot(HaveOccured())
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
