package matchers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ContainElements", func() {
	Context("with a slice", func() {
		It("should do the right thing", func() {
			Expect([]string{"foo", "bar", "baz"}).Should(ContainElements("foo", "bar", "baz"))
			Expect([]string{"foo", "bar", "baz"}).Should(ContainElements("bar"))
			Expect([]string{"foo", "bar", "baz"}).Should(ContainElements())
			Expect([]string{"foo", "bar", "baz"}).ShouldNot(ContainElements("baz", "bar", "foo", "foo"))
		})
	})

	Context("with an array", func() {
		It("should do the right thing", func() {
			Expect([3]string{"foo", "bar", "baz"}).Should(ContainElements("foo", "bar", "baz"))
			Expect([3]string{"foo", "bar", "baz"}).Should(ContainElements("bar"))
			Expect([3]string{"foo", "bar", "baz"}).Should(ContainElements())
			Expect([3]string{"foo", "bar", "baz"}).ShouldNot(ContainElements("baz", "bar", "foo", "foo"))
		})
	})

	Context("with a map", func() {
		It("should apply to the values", func() {
			Expect(map[int]string{1: "foo", 2: "bar", 3: "baz"}).Should(ContainElements("foo", "bar", "baz"))
			Expect(map[int]string{1: "foo", 2: "bar", 3: "baz"}).Should(ContainElements("bar"))
			Expect(map[int]string{1: "foo", 2: "bar", 3: "baz"}).Should(ContainElements())
			Expect(map[int]string{1: "foo", 2: "bar", 3: "baz"}).ShouldNot(ContainElements("baz", "bar", "foo", "foo"))
		})

	})

	Context("with anything else", func() {
		It("should error", func() {
			failures := InterceptGomegaFailures(func() {
				Expect("foo").Should(ContainElements("f", "o", "o"))
			})

			Expect(failures).Should(HaveLen(1))
		})
	})

	Context("when passed matchers", func() {
		It("should pass if the matchers pass", func() {
			Expect([]string{"foo", "bar", "baz"}).Should(ContainElements("foo", MatchRegexp("^ba"), "baz"))
			Expect([]string{"foo", "bar", "baz"}).Should(ContainElements("foo", MatchRegexp("^ba")))
			Expect([]string{"foo", "bar", "baz"}).ShouldNot(ContainElements("foo", MatchRegexp("^ba"), MatchRegexp("foo")))
			Expect([]string{"foo", "bar", "baz"}).Should(ContainElements("foo", MatchRegexp("^ba"), MatchRegexp("^ba")))
			Expect([]string{"foo", "bar", "baz"}).ShouldNot(ContainElements("foo", MatchRegexp("^ba"), MatchRegexp("turducken")))
		})

		It("should not depend on the order of the matchers", func() {
			Expect([][]int{{1, 2}, {2}}).Should(ContainElements(ContainElement(1), ContainElement(2)))
			Expect([][]int{{1, 2}, {2}}).Should(ContainElements(ContainElement(2), ContainElement(1)))
		})

		Context("when a matcher errors", func() {
			It("should soldier on", func() {
				Expect([]string{"foo", "bar", "baz"}).ShouldNot(ContainElements(BeFalse(), "foo", "bar"))
				Expect([]interface{}{"foo", "bar", false}).Should(ContainElements(BeFalse(), ContainSubstring("foo"), "bar"))
			})
		})
	})

	Context("when passed exactly one argument, and that argument is a slice", func() {
		It("should match against the elements of that argument", func() {
			Expect([]string{"foo", "bar", "baz"}).Should(ContainElements([]string{"foo", "baz"}))
			Expect([]string{"foo", "bar", "baz"}).ShouldNot(ContainElements([]string{"foo", "nope"}))
		})
	})

	Describe("FailureMessage", func() {
		It("prints missing elements", func() {
			failures := InterceptGomegaFailures(func() {
				Expect([]int{2}).Should(ContainElements(1, 2, 3))
			})

			expected := "Expected\n.*\\[2\\]\nto contain elements\n.*\\[1, 2, 3\\]\nthe missing elements were\n.*\\[1, 3\\]"
			Expect(failures).To(ContainElements(MatchRegexp(expected)))
		})

		When("expected was specified as an array", func() {
			It("flattens the array in the expectation message", func() {
				failures := InterceptGomegaFailures(func() {
					Expect([]string{"A", "B", "C"}).To(ContainElements([]string{"A", "D"}))
				})

				expected := `Expected\n.*\["A", "B", "C"\]\nto contain elements\n.*: \["A", "D"\]\nthe missing elements were\n.*\["D"\]`
				Expect(failures).To(ConsistOf(MatchRegexp(expected)))
			})

			It("flattens the array in the negated expectation message", func() {
				failures := InterceptGomegaFailures(func() {
					Expect([]string{"A", "B"}).NotTo(ContainElements([]string{"A", "B"}))
				})

				expected := `Expected\n.*\["A", "B"\]\nnot to contain elements\n.*: \["A", "B"\]`
				Expect(failures).To(ConsistOf(MatchRegexp(expected)))
			})
		})

		When("the expected values are the same type", func() {
			It("uses that type for the expectation slice", func() {
				failures := InterceptGomegaFailures(func() {
					Expect([]string{"A", "B"}).To(ContainElements("A", "B", "C"))
				})

				expected := `to contain elements
\s*<\[\]string \| len:3, cap:3>: \["A", "B", "C"\]
the missing elements were
\s*<\[\]string \| len:1, cap:1>: \["C"\]`
				Expect(failures).To(ConsistOf(MatchRegexp(expected)))
			})

			It("uses that type for the negated expectation slice", func() {
				failures := InterceptGomegaFailures(func() {
					Expect([]uint64{1, 2}).NotTo(ContainElements(uint64(1), uint64(2)))
				})

				expected := `not to contain elements\n\s*<\[\]uint64 \| len:2, cap:2>: \[1, 2\]`
				Expect(failures).To(ConsistOf(MatchRegexp(expected)))
			})
		})

		When("the expected values are different types", func() {
			It("uses interface{} for the expectation slice", func() {
				failures := InterceptGomegaFailures(func() {
					Expect([]interface{}{1, true}).To(ContainElements(1, "C"))
				})

				expected := `to contain elements
\s*<\[\]interface {} \| len:2, cap:2>: \[1, "C"\]
the missing elements were
\s*<\[\]string \| len:1, cap:1>: \["C"\]`
				Expect(failures).To(ConsistOf(MatchRegexp(expected)))
			})

			It("uses interface{} for the negated expectation slice", func() {
				failures := InterceptGomegaFailures(func() {
					Expect([]interface{}{1, "B"}).NotTo(ContainElements(1, "B"))
				})

				expected := `not to contain elements\n\s*<\[\]interface {} \| len:2, cap:2>: \[1, "B"\]`
				Expect(failures).To(ConsistOf(MatchRegexp(expected)))
			})
		})
	})
})
