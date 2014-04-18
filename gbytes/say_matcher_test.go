package gbytes_test

import (
	"time"
	. "github.com/onsi/gomega/gbytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SayMatcher", func() {
	var buffer *Buffer

	BeforeEach(func() {
		buffer = NewBuffer()
		buffer.Write([]byte("abc"))
	})

	Context("when actual is not a gexec Buffer", func() {
		It("should error", func() {
			failures := interceptFailures(func() {
				Ω("foo").Should(Say("foo"))
			})
			Ω(failures[0]).Should(ContainSubstring("gbytes Buffer"))
		})
	})

	Context("when a match is found", func() {
		It("should succeed", func() {
			Ω(buffer).Should(Say("abc"))
		})

		It("should support printf-like formatting", func() {
			Ω(buffer).Should(Say("a%sc", "b"))
		})

		It("should use a regular expression", func() {
			Ω(buffer).Should(Say("a.c"))
		})

		It("should fastforward the buffer", func() {
			buffer.Write([]byte("def"))
			Ω(buffer).Should(Say("abcd"))
			Ω(buffer).Should(Say("ef"))
			Ω(buffer).ShouldNot(Say("[a-z]"))
		})
	})

	Context("when a positive match fails", func() {
		It("should report where it got stuck", func() {
			Ω(buffer).Should(Say("abc"))
			buffer.Write([]byte("def"))
			failures := interceptFailures(func() {
				Ω(buffer).Should(Say("abc"))
			})
			Ω(failures[0]).Should(ContainSubstring("Got stuck at:"))
			Ω(failures[0]).Should(ContainSubstring("def"))
		})
	})

	Context("when a negative match fails", func() {
		It("should report where it got stuck", func() {
			failures := interceptFailures(func() {
				Ω(buffer).ShouldNot(Say("abc"))
			})
			Ω(failures[0]).Should(ContainSubstring("Saw:"))
			Ω(failures[0]).Should(ContainSubstring("Which matches the unexpected:"))
			Ω(failures[0]).Should(ContainSubstring("abc"))
		})
	})

	Context("when a match is not found", func() {
		It("should not fastforward the buffer", func() {
			Ω(buffer).ShouldNot(Say("def"))
			Ω(buffer).Should(Say("abc"))
		})
	})

	Context("a nice real-life example", func() {
		It("should behave well", func() {
			Ω(buffer).Should(Say("abc"))
			go func() {
				time.Sleep(10 * time.Millisecond)
				buffer.Write([]byte("def"))
			}()
			Ω(buffer).ShouldNot(Say("def"))
			Eventually(buffer).Should(Say("def"))
		})
	})
})
