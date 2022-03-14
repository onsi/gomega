package matchers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/matchers"
)

var _ = Describe("ContainElement", func() {
	Describe("matching only", func() {
		When("passed a supported type", func() {
			Context("and expecting a non-matcher", func() {
				It("should do the right thing", func() {
					Expect([2]int{1, 2}).Should(ContainElement(2))
					Expect([2]int{1, 2}).ShouldNot(ContainElement(3))

					Expect([]int{1, 2}).Should(ContainElement(2))
					Expect([]int{1, 2}).ShouldNot(ContainElement(3))

					Expect(map[string]int{"foo": 1, "bar": 2}).Should(ContainElement(2))
					Expect(map[int]int{3: 1, 4: 2}).ShouldNot(ContainElement(3))

					arr := make([]myCustomType, 2)
					arr[0] = myCustomType{s: "foo", n: 3, f: 2.0, arr: []string{"a", "b"}}
					arr[1] = myCustomType{s: "foo", n: 3, f: 2.0, arr: []string{"a", "c"}}
					Expect(arr).Should(ContainElement(myCustomType{s: "foo", n: 3, f: 2.0, arr: []string{"a", "b"}}))
					Expect(arr).ShouldNot(ContainElement(myCustomType{s: "foo", n: 3, f: 2.0, arr: []string{"b", "c"}}))
				})
			})

			Context("and expecting a matcher", func() {
				It("should pass each element through the matcher", func() {
					Expect([]int{1, 2, 3}).Should(ContainElement(BeNumerically(">=", 3)))
					Expect([]int{1, 2, 3}).ShouldNot(ContainElement(BeNumerically(">", 3)))
					Expect(map[string]int{"foo": 1, "bar": 2}).Should(ContainElement(BeNumerically(">=", 2)))
					Expect(map[string]int{"foo": 1, "bar": 2}).ShouldNot(ContainElement(BeNumerically(">", 2)))
				})

				It("should power through even if the matcher ever fails", func() {
					Expect([]interface{}{1, 2, "3", 4}).Should(ContainElement(BeNumerically(">=", 3)))
				})

				It("should fail if the matcher fails", func() {
					actual := []interface{}{1, 2, "3", "4"}
					success, err := (&ContainElementMatcher{Element: BeNumerically(">=", 3)}).Match(actual)
					Expect(success).Should(BeFalse())
					Expect(err).Should(HaveOccurred())
				})
			})
		})

		When("passed a correctly typed nil", func() {
			It("should operate succesfully on the passed in value", func() {
				var nilSlice []int
				Expect(nilSlice).ShouldNot(ContainElement(1))

				var nilMap map[int]string
				Expect(nilMap).ShouldNot(ContainElement("foo"))
			})
		})

		When("passed an unsupported type", func() {
			It("should error", func() {
				success, err := (&ContainElementMatcher{Element: 0}).Match(0)
				Expect(success).Should(BeFalse())
				Expect(err).Should(HaveOccurred())

				success, err = (&ContainElementMatcher{Element: 0}).Match("abc")
				Expect(success).Should(BeFalse())
				Expect(err).Should(HaveOccurred())

				success, err = (&ContainElementMatcher{Element: 0}).Match(nil)
				Expect(success).Should(BeFalse())
				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Describe("returning findings", func() {
		It("rejects a nil result reference", func() {
			Expect(ContainElement("foo", nil).Match([]string{"foo"})).Error().To(
				MatchError(MatchRegexp(`expects a non-nil pointer.+ Got\n +<nil>: nil`)))
		})

		Context("with match(es)", func() {
			When("passed an assignable result reference", func() {
				It("should assign a single finding to a scalar result reference", func() {
					actual := []string{"bar", "foo"}
					var stash string
					Expect(actual).To(ContainElement("foo", &stash))
					Expect(stash).To(Equal("foo"))

					actualmap := map[int]string{
						1: "bar",
						2: "foo",
					}
					Expect(actualmap).To(ContainElement("foo", &stash))
					Expect(stash).To(Equal("foo"))
				})

				It("should assign a single finding to a slice return reference", func() {
					actual := []string{"bar", "foo", "baz"}
					var stash []string
					Expect(actual).To(ContainElement("foo", &stash))
					Expect(stash).To(HaveLen(1))
					Expect(stash).To(ContainElement("foo"))
				})

				It("should assign multiple findings to a slice return reference", func() {
					actual := []string{"bar", "foo", "bar", "foo"}
					var stash []string
					Expect(actual).To(ContainElement("foo", &stash))
					Expect(stash).To(HaveLen(2))
					Expect(stash).To(HaveEach("foo"))
				})

				It("should assign map findings to a map return reference", func() {
					actual := map[string]string{
						"foo": "foo",
						"bar": "bar",
						"baz": "baz",
					}
					var stash map[string]string
					Expect(actual).To(ContainElement(ContainSubstring("ba"), &stash))
					Expect(stash).To(HaveLen(2))
					Expect(stash).To(ConsistOf("bar", "baz"))
				})
			})

			When("passed a scalar return reference for multiple matches", func() {
				It("should error", func() {
					actual := []string{"foo", "foo"}
					var stash string
					Expect(ContainElement("foo", &stash).Match(actual)).Error().To(
						MatchError(MatchRegexp(`cannot return multiple findings\.  Need \*\[\]string, got \*string`)))
				})
			})

			When("passed an unassignable return reference for matches", func() {
				It("should error for actual []T1, return reference T2", func() {
					actual := []string{"bar", "foo"}
					var stash int
					Expect(ContainElement("foo", &stash).Match(actual)).Error().To(HaveOccurred())
				})
				It("should error for actual []T, return reference [...]T", func() {
					actual := []string{"bar", "foo"}
					var arrstash [2]string
					Expect(ContainElement("foo", &arrstash).Match(actual)).Error().To(HaveOccurred())
				})
				It("should error for actual []interface{}, return reference T", func() {
					actual := []interface{}{"foo", 42}
					var stash int
					Expect(ContainElement(Not(BeZero()), &stash).Match(actual)).Error().To(
						MatchError(MatchRegexp(`cannot return findings\.  Need \*interface.+, got \*int`)))
				})
				It("should error for actual []interface{}, return reference []T", func() {
					actual := []interface{}{"foo", 42}
					var stash []string
					Expect(ContainElement(Not(BeZero()), &stash).Match(actual)).Error().To(
						MatchError(MatchRegexp(`cannot return findings\.  Need \*\[\]interface.+, got \*\[\]string`)))
				})
				It("should error for actual map[T]T, return reference map[T]interface{}", func() {
					actual := map[string]string{
						"foo": "foo",
						"bar": "bar",
						"baz": "baz",
					}
					var stash map[string]interface{}
					Expect(ContainElement(Not(BeZero()), &stash).Match(actual)).Error().To(
						MatchError(MatchRegexp(`cannot return findings\.  Need \*map\[string\]string, got \*map\[string\]interface`)))
				})
				It("should error for actual map[T]T, return reference []T", func() {
					actual := map[string]string{
						"foo": "foo",
						"bar": "bar",
						"baz": "baz",
					}
					var stash []string
					Expect(ContainElement(Not(BeZero()), &stash).Match(actual)).Error().To(
						MatchError(MatchRegexp(`cannot return findings\.  Need \*map\[string\]string, got \*\[\]string`)))
				})

				It("should return a descriptive return reference-type error", func() {
					actual := []string{"bar", "foo"}
					var stash map[string]struct{}
					Expect(ContainElement("foo", &stash).Match(actual)).Error().To(
						MatchError(MatchRegexp(`cannot return findings\.  Need \*\[\]string, got \*map`)))
				})
			})
		})

		Context("without any matches", func() {
			When("the matcher did not error", func() {
				It("should report non-match", func() {
					actual := []string{"bar", "foo"}
					var stash string
					rem := ContainElement("baz", &stash)
					m, err := rem.Match(actual)
					Expect(m).To(BeFalse())
					Expect(err).NotTo(HaveOccurred())
					Expect(rem.FailureMessage(actual)).To(MatchRegexp(`Expected\n.+\nto contain element matching\n.+: baz`))

					var stashslice []string
					rem = ContainElement("baz", &stashslice)
					m, err = rem.Match(actual)
					Expect(m).To(BeFalse())
					Expect(err).NotTo(HaveOccurred())
					Expect(rem.FailureMessage(actual)).To(MatchRegexp(`Expected\n.+\nto contain element matching\n.+: baz`))
				})
			})

			When("the matcher errors", func() {
				It("should report last matcher error", func() {
					actual := []interface{}{"bar", 42}
					var stash []interface{}
					Expect(ContainElement(HaveField("yeehaw", 42), &stash).Match(actual)).Error().To(MatchError(MatchRegexp(`HaveField encountered:\n.*<int>: 42\nWhich is not a struct`)))
				})
			})
		})
	})

})
