package gomega

import (
	. "github.com/onsi/ginkgo"
	"time"
)

func init() {
	Describe("AsyncActual", func() {
		var (
			failureMessage string
			callerSkip     int
		)

		var fakeFailHandler = func(message string, skip ...int) {
			failureMessage = message
			callerSkip = skip[0]
		}

		BeforeEach(func() {
			failureMessage = ""
			callerSkip = 0
		})

		Context("when passed a function", func() {
			Context("the positive case", func() {
				It("should poll the function and matcher", func() {
					arr := []int{}
					a := newAsyncActual(func() []int {
						arr = append(arr, 1)
						return arr
					}, fakeFailHandler, time.Duration(0.2*float64(time.Second)), time.Duration(0.02*float64(time.Second)), 1)

					a.Should(HaveLen(10))

					Ω(arr).Should(HaveLen(10))
					Ω(failureMessage).Should(BeZero())
				})

				It("should continue when the matcher errors", func() {
					var arr = []int{}
					a := newAsyncActual(func() interface{} {
						arr = append(arr, 1)
						if len(arr) == 4 {
							return 0 //this should cause the matcher to error
						}
						return arr
					}, fakeFailHandler, time.Duration(0.2*float64(time.Second)), time.Duration(0.02*float64(time.Second)), 1)

					a.Should(HaveLen(4), "My description %d", 2)

					Ω(failureMessage).Should(ContainSubstring("Timed out after"))
					Ω(failureMessage).Should(ContainSubstring("My description 2"))
					Ω(callerSkip).Should(Equal(3))
				})

				It("should be able to timeout", func() {
					arr := []int{}
					a := newAsyncActual(func() []int {
						arr = append(arr, 1)
						return arr
					}, fakeFailHandler, time.Duration(0.2*float64(time.Second)), time.Duration(0.02*float64(time.Second)), 1)

					a.Should(HaveLen(11), "My description %d", 2)

					Ω(arr).Should(HaveLen(10))
					Ω(failureMessage).Should(ContainSubstring("Timed out after"))
					Ω(failureMessage).Should(ContainSubstring("My description 2"))
					Ω(callerSkip).Should(Equal(3))
				})
			})

			Context("the negative case", func() {
				It("should poll the function and matcher", func() {
					counter := 0
					arr := []int{}
					a := newAsyncActual(func() []int {
						counter += 1
						if counter >= 10 {
							arr = append(arr, 1)
						}
						return arr
					}, fakeFailHandler, time.Duration(0.2*float64(time.Second)), time.Duration(0.02*float64(time.Second)), 1)

					a.ShouldNot(HaveLen(0))

					Ω(arr).Should(HaveLen(1))
					Ω(failureMessage).Should(BeZero())
				})

				It("should timeout when the matcher errors", func() {
					a := newAsyncActual(func() interface{} {
						return 0 //this should cause the matcher to error
					}, fakeFailHandler, time.Duration(0.2*float64(time.Second)), time.Duration(0.02*float64(time.Second)), 1)

					a.ShouldNot(HaveLen(0), "My description %d", 2)

					Ω(failureMessage).Should(ContainSubstring("Timed out after"))
					Ω(failureMessage).Should(ContainSubstring("Error:"))
					Ω(failureMessage).Should(ContainSubstring("My description 2"))
					Ω(callerSkip).Should(Equal(3))
				})

				It("should be able to timeout", func() {
					a := newAsyncActual(func() []int {
						return []int{}
					}, fakeFailHandler, time.Duration(0.2*float64(time.Second)), time.Duration(0.02*float64(time.Second)), 1)

					a.ShouldNot(HaveLen(0), "My description %d", 2)

					Ω(failureMessage).Should(ContainSubstring("Timed out after"))
					Ω(failureMessage).Should(ContainSubstring("My description 2"))
					Ω(callerSkip).Should(Equal(3))
				})
			})
		})

		Context("when passed a function with the wrong # or arguments & returns", func() {
			It("should panic", func() {
				Ω(func() {
					newAsyncActual(func() {}, fakeFailHandler, 0, 0, 1)
				}).Should(Panic())

				Ω(func() {
					newAsyncActual(func(a string) int { return 0 }, fakeFailHandler, 0, 0, 1)
				}).Should(Panic())

				Ω(func() {
					newAsyncActual(func() int { return 0 }, fakeFailHandler, 0, 0, 1)
				}).ShouldNot(Panic())
			})
		})
	})
}
