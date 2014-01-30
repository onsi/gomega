package gomega

import (
	"errors"
	. "github.com/onsi/ginkgo"
)

func init() {
	Describe("Actual", func() {
		var (
			a                 *actual
			failureMessage    string
			failureCallerSkip int
			matcher           *fakeMatcher
		)

		input := "The thing I'm testing"

		var fakeFailHandler = func(message string, callerSkip ...int) {
			failureMessage = message
			if len(callerSkip) == 1 {
				failureCallerSkip = callerSkip[0]
			}
		}

		BeforeEach(func() {
			matcher = &fakeMatcher{}
			failureMessage = ""
			failureCallerSkip = 0
			a = newActual(input, fakeFailHandler, 1)
		})

		Context("when called", func() {
			It("should pass the provided input value to the matcher", func() {
				a.Should(matcher)

				Ω(matcher.receivedActual).Should(Equal(input))
				matcher.receivedActual = ""

				a.ShouldNot(matcher)

				Ω(matcher.receivedActual).Should(Equal(input))
				matcher.receivedActual = ""

				a.To(matcher)

				Ω(matcher.receivedActual).Should(Equal(input))
				matcher.receivedActual = ""

				a.ToNot(matcher)

				Ω(matcher.receivedActual).Should(Equal(input))
				matcher.receivedActual = ""

				a.NotTo(matcher)

				Ω(matcher.receivedActual).Should(Equal(input))
			})
		})

		Context("when the matcher succeeds", func() {
			BeforeEach(func() {
				matcher.matchesToReturn = true
				matcher.messageToReturn = "The negative failure message"
				matcher.errToReturn = nil
			})

			Context("and a positive assertion is being made", func() {
				It("should not call the failure callback", func() {
					a.Should(matcher)
					Ω(failureMessage).Should(Equal(""))
				})

				It("should be true", func() {
					Ω(a.Should(matcher)).Should(BeTrue())
				})
			})

			Context("and a negative assertion is being made", func() {
				It("should call the failure callback", func() {
					a.ShouldNot(matcher)
					Ω(failureMessage).Should(Equal("The negative failure message"))
					Ω(failureCallerSkip).Should(Equal(3))
				})

				It("should be false", func() {
					Ω(a.ShouldNot(matcher)).Should(BeFalse())
				})
			})
		})

		Context("when the matcher fails", func() {
			BeforeEach(func() {
				matcher.matchesToReturn = false
				matcher.messageToReturn = "The positive failure message"
				matcher.errToReturn = nil
			})

			Context("and a positive assertion is being made", func() {
				It("should call the failure callback", func() {
					a.Should(matcher)
					Ω(failureMessage).Should(Equal("The positive failure message"))
					Ω(failureCallerSkip).Should(Equal(3))
				})

				It("should be false", func() {
					Ω(a.Should(matcher)).Should(BeFalse())
				})
			})

			Context("and a negative assertion is being made", func() {
				It("should not call the failure callback", func() {
					a.ShouldNot(matcher)
					Ω(failureMessage).Should(Equal(""))
				})

				It("should be true", func() {
					Ω(a.ShouldNot(matcher)).Should(BeTrue())
				})
			})
		})

		Context("When reporting a failure", func() {
			BeforeEach(func() {
				matcher.matchesToReturn = false
				matcher.messageToReturn = "The positive failure message"
				matcher.errToReturn = nil
			})

			Context("and there is an optional description", func() {
				It("should append the description to the failure message", func() {
					a.Should(matcher, "A description")
					Ω(failureMessage).Should(Equal("A description\nThe positive failure message"))
					Ω(failureCallerSkip).Should(Equal(3))
				})
			})

			Context("and there are multiple arguments to the optional description", func() {
				It("should append the formatted description to the failure message", func() {
					a.Should(matcher, "A description of [%d]", 3)
					Ω(failureMessage).Should(Equal("A description of [3]\nThe positive failure message"))
					Ω(failureCallerSkip).Should(Equal(3))
				})
			})
		})

		Context("When the matcher returns an error", func() {
			BeforeEach(func() {
				matcher.errToReturn = errors.New("Kaboom!")
			})

			Context("and a positive assertion is being made", func() {
				It("should call the failure callback", func() {
					matcher.matchesToReturn = true
					matcher.messageToReturn = "Ignore me"
					a.Should(matcher)
					Ω(failureMessage).Should(Equal("Kaboom!"))
					Ω(failureCallerSkip).Should(Equal(3))
				})
			})

			Context("and a negative assertion is being made", func() {
				It("should call the failure callback", func() {
					matcher.matchesToReturn = false
					matcher.messageToReturn = "Ignore me"
					a.ShouldNot(matcher)
					Ω(failureMessage).Should(Equal("Kaboom!"))
					Ω(failureCallerSkip).Should(Equal(3))
				})
			})

			It("should always be false", func() {
				Ω(a.Should(matcher)).Should(BeFalse())
				Ω(a.ShouldNot(matcher)).Should(BeFalse())
			})
		})
	})
}
