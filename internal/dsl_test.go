package internal_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Gomega DSL", func() {
	Describe("InterceptGomegaFailures", func() {
		Context("when no failures occur", func() {
			It("returns an empty array", func() {
				Expect(InterceptGomegaFailures(func() {
					Expect("hi").To(Equal("hi"))
				})).To(BeEmpty())
			})
		})

		Context("when failures occur", func() {
			It("does not stop execution and returns all the failures as strings", func() {
				Expect(InterceptGomegaFailures(func() {
					Expect("hi").To(Equal("bye"))
					Expect(3).To(Equal(2))
				})).To(Equal([]string{
					"Expected\n    <string>: hi\nto equal\n    <string>: bye",
					"Expected\n    <int>: 3\nto equal\n    <int>: 2",
				}))

			})
		})
	})

	Describe("InterceptGomegaFailure", func() {
		Context("when no failures occur", func() {
			It("returns nil", func() {
				Expect(InterceptGomegaFailure(func() {
					Expect("hi").To(Equal("hi"))
				})).To(BeNil())
			})
		})

		Context("when failures occur", func() {
			It("returns the first failure and stops execution", func() {
				gotThere := false
				Expect(InterceptGomegaFailure(func() {
					Expect("hi").To(Equal("bye"))
					gotThere = true
					Expect(3).To(Equal(2))
				})).To(Equal(errors.New("Expected\n    <string>: hi\nto equal\n    <string>: bye")))
				Expect(gotThere).To(BeFalse())
			})
		})

		Context("when the function panics", func() {
			It("panics", func() {
				Expect(func() {
					InterceptGomegaFailure(func() {
						panic("boom")
					})
				}).To(PanicWith("boom"))
			})
		})
	})
})
