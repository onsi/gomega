package internal_test

import (
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/internal"
)

var _ = Describe("AsyncSignalError", func() {
	Describe("building StopTrying errors", func() {
		It("can build a formatted message", func() {
			st := StopTrying("I've tried %d times - give up!", 17)
			Ω(st.Error()).Should(Equal("I've tried 17 times - give up!"))
			Ω(errors.Unwrap(st)).Should(BeNil())
		})

		It("can wrap other errors", func() {
			st := StopTrying("Welp!  Server said: %w", fmt.Errorf("ERR_GIVE_UP"))
			Ω(st.Error()).Should(Equal("Welp!  Server said: ERR_GIVE_UP"))
			Ω(errors.Unwrap(st)).Should(Equal(fmt.Errorf("ERR_GIVE_UP")))
		})
	})

	Describe("when invoking Now()", func() {
		It("should not a panic occurred and panic with itself", func() {
			st := StopTrying("bam").(*internal.AsyncSignalError)
			Ω(st.WasViaPanic()).Should(BeFalse())
			Ω(st.Now).Should(PanicWith(st))
			Ω(st.WasViaPanic()).Should(BeTrue())
		})
	})

	Describe("AsAsyncSignalError", func() {
		It("should return false for nils", func() {
			st, ok := internal.AsAsyncSignalError(nil)
			Ω(st).Should(BeNil())
			Ω(ok).Should(BeFalse())
		})

		It("should work when passed a StopTrying error", func() {
			st, ok := internal.AsAsyncSignalError(StopTrying("bam"))
			Ω(st).Should(Equal(StopTrying("bam")))
			Ω(ok).Should(BeTrue())
		})

		It("should work when passed a wrapped error", func() {
			st, ok := internal.AsAsyncSignalError(fmt.Errorf("STOP TRYING %w", StopTrying("bam")))
			Ω(st).Should(Equal(StopTrying("bam")))
			Ω(ok).Should(BeTrue())
		})
	})
})
