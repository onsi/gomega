package matchers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/matchers"
	"time"
)

var _ = Describe("BeAfter", func() {
	Context("When comparing non-times", func() {
		It("should error", func() {
			success, err := (&BeAfterMatcher{Expected: 1}).Match(1)

			Ω(success).Should(BeFalse())
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("When comparing times", func() {
		It("should do the right thing", func() {
			t0 := time.Now()
			t1 := t0.Add(time.Second)
			Ω(t1).Should(BeAfter(t0))
			Ω(t0).ShouldNot(BeAfter(t1))
			Ω(t0).ShouldNot(BeAfter(t0))
		})
	})
})
