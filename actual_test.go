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
					Ω(failureMessage).Should(Equal("negative: The thing I'm testing"))
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
				matcher.errToReturn = nil
			})

			Context("and a positive assertion is being made", func() {
				It("should call the failure callback", func() {
					a.Should(matcher)
					Ω(failureMessage).Should(Equal("positive: The thing I'm testing"))
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
				matcher.errToReturn = nil
			})

			Context("and there is an optional description", func() {
				It("should append the description to the failure message", func() {
					a.Should(matcher, "A description")
					Ω(failureMessage).Should(Equal("A description\npositive: The thing I'm testing"))
					Ω(failureCallerSkip).Should(Equal(3))
				})
			})

			Context("and there are multiple arguments to the optional description", func() {
				It("should append the formatted description to the failure message", func() {
					a.Should(matcher, "A description of [%d]", 3)
					Ω(failureMessage).Should(Equal("A description of [3]\npositive: The thing I'm testing"))
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
					a.Should(matcher)
					Ω(failureMessage).Should(Equal("Kaboom!"))
					Ω(failureCallerSkip).Should(Equal(3))
				})
			})

			Context("and a negative assertion is being made", func() {
				It("should call the failure callback", func() {
					matcher.matchesToReturn = false
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

		Context("when there are extra parameters", func() {
			It("(a simple example)", func() {
				Ω(func() (string, int, error) {
					return "foo", 0, nil
				}()).Should(Equal("foo"))
			})

			Context("when the parameters are all nil or zero", func() {
				It("should invoke the matcher", func() {
					matcher.matchesToReturn = true
					matcher.errToReturn = nil

					var typedNil []string
					a = newActual(input, fakeFailHandler, 1, 0, nil, typedNil)

					result := a.Should(matcher)
					Ω(result).Should(BeTrue())
					Ω(matcher.receivedActual).Should(Equal(input))

					Ω(failureMessage).Should(BeZero())
				})
			})

			Context("when any of the parameters are not nil or zero", func() {
				It("should call the failure callback", func() {
					matcher.matchesToReturn = false
					matcher.errToReturn = nil

					a = newActual(input, fakeFailHandler, 1, errors.New("foo"))
					result := a.Should(matcher)
					Ω(result).Should(BeFalse())
					Ω(matcher.receivedActual).Should(BeZero(), "The matcher doesn't even get called")
					Ω(failureMessage).Should(ContainSubstring("foo"))
					failureMessage = ""

					a = newActual(input, fakeFailHandler, 1, nil, 1)
					result = a.ShouldNot(matcher)
					Ω(result).Should(BeFalse())
					Ω(failureMessage).Should(ContainSubstring("1"))
					failureMessage = ""

					a = newActual(input, fakeFailHandler, 1, nil, 0, []string{"foo"})
					result = a.To(matcher)
					Ω(result).Should(BeFalse())
					Ω(failureMessage).Should(ContainSubstring("foo"))
					failureMessage = ""

					a = newActual(input, fakeFailHandler, 1, nil, 0, []string{"foo"})
					result = a.ToNot(matcher)
					Ω(result).Should(BeFalse())
					Ω(failureMessage).Should(ContainSubstring("foo"))
					failureMessage = ""

					a = newActual(input, fakeFailHandler, 1, nil, 0, []string{"foo"})
					result = a.NotTo(matcher)
					Ω(result).Should(BeFalse())
					Ω(failureMessage).Should(ContainSubstring("foo"))
					Ω(failureCallerSkip).Should(Equal(3))
				})
			})
		})
	})
}
