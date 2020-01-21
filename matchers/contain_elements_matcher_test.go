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
	})
})
