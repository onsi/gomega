package matchers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ConsistOf", func() {
	Context("with a slice", func() {
		It("should do the right thing", func() {
			Expect([]string{"foo", "bar", "baz"}).Should(ConsistOf("foo", "bar", "baz"))
			Expect([]string{"foo", "bar", "baz"}).Should(ConsistOf("foo", "bar", "baz"))
			Expect([]string{"foo", "bar", "baz"}).Should(ConsistOf("baz", "bar", "foo"))
			Expect([]string{"foo", "bar", "baz"}).ShouldNot(ConsistOf("baz", "bar", "foo", "foo"))
			Expect([]string{"foo", "bar", "baz"}).ShouldNot(ConsistOf("baz", "foo"))
		})
	})

	Context("with an array", func() {
		It("should do the right thing", func() {
			Expect([3]string{"foo", "bar", "baz"}).Should(ConsistOf("foo", "bar", "baz"))
			Expect([3]string{"foo", "bar", "baz"}).Should(ConsistOf("baz", "bar", "foo"))
			Expect([3]string{"foo", "bar", "baz"}).ShouldNot(ConsistOf("baz", "bar", "foo", "foo"))
			Expect([3]string{"foo", "bar", "baz"}).ShouldNot(ConsistOf("baz", "foo"))
		})
	})

	Context("with a map", func() {
		It("should apply to the values", func() {
			Expect(map[int]string{1: "foo", 2: "bar", 3: "baz"}).Should(ConsistOf("foo", "bar", "baz"))
			Expect(map[int]string{1: "foo", 2: "bar", 3: "baz"}).Should(ConsistOf("baz", "bar", "foo"))
			Expect(map[int]string{1: "foo", 2: "bar", 3: "baz"}).ShouldNot(ConsistOf("baz", "bar", "foo", "foo"))
			Expect(map[int]string{1: "foo", 2: "bar", 3: "baz"}).ShouldNot(ConsistOf("baz", "foo"))
		})

	})

	Context("with anything else", func() {
		It("should error", func() {
			failures := InterceptGomegaFailures(func() {
				Expect("foo").Should(ConsistOf("f", "o", "o"))
			})

			Expect(failures).Should(HaveLen(1))
		})
	})

	When("passed matchers", func() {
		It("should pass if the matchers pass", func() {
			Expect([]string{"foo", "bar", "baz"}).Should(ConsistOf("foo", MatchRegexp("^ba"), "baz"))
			Expect([]string{"foo", "bar", "baz"}).ShouldNot(ConsistOf("foo", MatchRegexp("^ba")))
			Expect([]string{"foo", "bar", "baz"}).ShouldNot(ConsistOf("foo", MatchRegexp("^ba"), MatchRegexp("foo")))
			Expect([]string{"foo", "bar", "baz"}).Should(ConsistOf("foo", MatchRegexp("^ba"), MatchRegexp("^ba")))
			Expect([]string{"foo", "bar", "baz"}).ShouldNot(ConsistOf("foo", MatchRegexp("^ba"), MatchRegexp("turducken")))
		})

		It("should not depend on the order of the matchers", func() {
			Expect([][]int{{1, 2}, {2}}).Should(ConsistOf(ContainElement(1), ContainElement(2)))
			Expect([][]int{{1, 2}, {2}}).Should(ConsistOf(ContainElement(2), ContainElement(1)))
		})

		When("a matcher errors", func() {
			It("should soldier on", func() {
				Expect([]string{"foo", "bar", "baz"}).ShouldNot(ConsistOf(BeFalse(), "foo", "bar"))
				Expect([]interface{}{"foo", "bar", false}).Should(ConsistOf(BeFalse(), ContainSubstring("foo"), "bar"))
			})
		})
	})

	When("passed exactly one argument, and that argument is a slice", func() {
		It("should match against the elements of that argument", func() {
			Expect([]string{"foo", "bar", "baz"}).Should(ConsistOf([]string{"foo", "bar", "baz"}))
			Expect([]string{"foo", "bar", "baz"}).ShouldNot(ConsistOf([]string{"foo", "bar"}))
		})
	})

	Describe("FailureMessage", func() {
		When("actual contains an extra element", func() {
			It("prints the extra element", func() {
				failures := InterceptGomegaFailures(func() {
					Expect([]int{1, 2}).Should(ConsistOf(2))
				})

				expected := "Expected\n.*\\[1, 2\\]\nto consist of\n.*\\[2\\]\nthe extra elements were\n.*\\[1\\]"
				Expect(failures).To(ConsistOf(MatchRegexp(expected)))
			})
		})

		When("actual misses an element", func() {
			It("prints the missing element", func() {
				failures := InterceptGomegaFailures(func() {
					Expect([]int{2}).Should(ConsistOf(1, 2))
				})

				expected := "Expected\n.*\\[2\\]\nto consist of\n.*\\[1, 2\\]\nthe missing elements were\n.*\\[1\\]"
				Expect(failures).To(ConsistOf(MatchRegexp(expected)))
			})
		})

		When("actual contains an extra element and misses an element", func() {
			It("prints both the extra and missing element", func() {
				failures := InterceptGomegaFailures(func() {
					Expect([]int{1, 2}).Should(ConsistOf(2, 3))
				})

				expected := "Expected\n.*\\[1, 2\\]\nto consist of\n.*\\[2, 3\\]\nthe missing elements were\n.*\\[3\\]\nthe extra elements were\n.*\\[1\\]"
				Expect(failures).To(ConsistOf(MatchRegexp(expected)))
			})
		})

		When("expected was specified as an array", func() {
			It("flattens the array in the expectation message", func() {
				failures := InterceptGomegaFailures(func() {
					Expect([]string{"A", "B", "C"}).To(ConsistOf([]string{"A", "B"}))
				})

				expected := `Expected\n.*\["A", "B", "C"\]\nto consist of\n.*: \["A", "B"\]\nthe extra elements were\n.*\["C"\]`
				Expect(failures).To(ConsistOf(MatchRegexp(expected)))
			})

			It("flattens the array in the negated expectation message", func() {
				failures := InterceptGomegaFailures(func() {
					Expect([]string{"A", "B"}).NotTo(ConsistOf([]string{"A", "B"}))
				})

				expected := `Expected\n.*\["A", "B"\]\nnot to consist of\n.*: \["A", "B"\]`
				Expect(failures).To(ConsistOf(MatchRegexp(expected)))
			})
		})

		When("the expected values are the same type", func() {
			It("uses that type for the expectation slice", func() {
				failures := InterceptGomegaFailures(func() {
					Expect([]string{"A", "B"}).To(ConsistOf("A", "C"))
				})

				expected := `to consist of
\s*<\[\]string \| len:2, cap:2>: \["A", "C"\]
the missing elements were
\s*<\[\]string \| len:1, cap:1>: \["C"\]
the extra elements were
\s*<\[\]string \| len:1, cap:1>: \["B"\]`
				Expect(failures).To(ConsistOf(MatchRegexp(expected)))
			})

			It("uses that type for the negated expectation slice", func() {
				failures := InterceptGomegaFailures(func() {
					Expect([]uint64{1, 2}).NotTo(ConsistOf(uint64(1), uint64(2)))
				})

				expected := `not to consist of\n\s*<\[\]uint64 \| len:2, cap:2>: \[1, 2\]`
				Expect(failures).To(ConsistOf(MatchRegexp(expected)))
			})
		})

		When("the expected values are different types", func() {
			It("uses interface{} for the expectation slice", func() {
				failures := InterceptGomegaFailures(func() {
					Expect([]interface{}{1, true}).To(ConsistOf(1, "C"))
				})

				expected := `to consist of
\s*<\[\]interface {} \| len:2, cap:2>: \[1, "C"\]
the missing elements were
\s*<\[\]string \| len:1, cap:1>: \["C"\]
the extra elements were
\s*<\[\]bool \| len:1, cap:1>: \[true\]`
				Expect(failures).To(ConsistOf(MatchRegexp(expected)))
			})

			It("uses interface{} for the negated expectation slice", func() {
				failures := InterceptGomegaFailures(func() {
					Expect([]interface{}{1, "B"}).NotTo(ConsistOf(1, "B"))
				})

				expected := `not to consist of\n\s*<\[\]interface {} \| len:2, cap:2>: \[1, "B"\]`
				Expect(failures).To(ConsistOf(MatchRegexp(expected)))
			})
		})
	})
})
