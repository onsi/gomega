package matchers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/matchers/internal/miter"
)

var _ = Describe("ContainMap", func() {
	var (
		mapA    map[string]string
		mapB    map[string]string
		customA *myCustomType
		customB *myCustomType
		objMap  map[*myCustomType]*myCustomType
	)

	BeforeEach(func() {
		mapA = map[string]string{
			"keyA": "valA",
			"key2": "val2",
			"keyB": "valB",
			"key1": "val1",
		}
		mapB = map[string]string{
			"key1": "val1",
			"key2": "val2",
		}

		customA = &myCustomType{s: "a", n: 2, f: 2.3, arr: []string{"ice", "cream"}}
		customB = &myCustomType{s: "b", n: 4, f: 3.1, arr: []string{"cake"}}
		objMap = map[*myCustomType]*myCustomType{
			customA: customA,
			customB: customA,
		}
	})

	When("passed a map", func() {
		It("succeeds when actual contains all key-value pairs of expected", func() {
			Expect(mapA).Should(ContainMap(mapB))
			Expect(mapA).Should(ContainMap(mapA))
			Expect(mapA).Should(ContainMap(map[string]string{"keyA": "valA"}))
		})

		It("fails when actual is missing a key or has a differing value", func() {
			Expect(mapB).ShouldNot(ContainMap(mapA))
			Expect(mapA).ShouldNot(ContainMap(map[string]string{"key1": "val_wrong"}))
			Expect(mapA).ShouldNot(ContainMap(map[string]string{"missingKey": "val1"}))
			Expect(mapA).ShouldNot(ContainMap(map[string]string{"key1": "val1", "missing": "val"}))
		})

		It("handles empty maps correctly", func() {
			emptyMap := map[string]string{}
			Expect(mapA).Should(ContainMap(emptyMap))
			Expect(emptyMap).Should(ContainMap(emptyMap))
			Expect(emptyMap).ShouldNot(ContainMap(mapB))
		})

		It("operates on pointer keys and values", func() {
			Expect(objMap).Should(ContainMap(map[*myCustomType]*myCustomType{customA: customA}))
			Expect(objMap).Should(ContainMap(map[*myCustomType]*myCustomType{
				&myCustomType{s: "b", n: 4, f: 3.1, arr: []string{"cake"}}: &myCustomType{s: "a", n: 2, f: 2.3, arr: []string{"ice", "cream"}},
			}))
			Expect(objMap).ShouldNot(ContainMap(map[*myCustomType]*myCustomType{
				customA: customB,
			}))
		})

		It("operates across different map types when types are compatible", func() {
			heteroActual := map[string]any{
				"a": 1,
				"b": "text",
				"c": true,
			}
			Expect(heteroActual).Should(ContainMap(map[string]int{"a": 1}))
			Expect(heteroActual).Should(ContainMap(map[string]string{"b": "text"}))
			Expect(heteroActual).ShouldNot(ContainMap(map[string]int{"b": 1}))
			Expect(heteroActual).Should(ContainMap(map[any]any{"c": true}))
		})
	})

	When("passed a correctly typed nil map", func() {
		It("behaves like an empty map", func() {
			var nilMap map[string]string
			Expect(mapA).Should(ContainMap(nilMap))
			Expect(nilMap).Should(ContainMap(nilMap))
			Expect(nilMap).Should(ContainMap(map[string]string{}))
			Expect(nilMap).ShouldNot(ContainMap(mapB))
		})
	})

	When("the expected map contains matchers as values", func() {
		It("evaluates values using matchers", func() {
			numMap := map[string]int{
				"first":  10,
				"second": 20,
				"third":  30,
			}
			Expect(numMap).Should(ContainMap(map[string]any{
				"first":  BeNumerically(">", 5),
				"second": BeNumerically("==", 20),
			}))
			Expect(numMap).ShouldNot(ContainMap(map[string]any{
				"first": BeNumerically("<", 5),
			}))
			Expect(numMap).ShouldNot(ContainMap(map[string]any{
				"fourth": BeNumerically(">", 0),
			}))
		})

		It("matches nil values with BeNil", func() {
			mapWithNil := map[string]any{
				"key1": nil,
				"key2": "val",
			}
			Expect(mapWithNil).Should(ContainMap(map[string]any{
				"key1": nil,
			}))
			Expect(mapWithNil).Should(ContainMap(map[string]any{
				"key1": BeNil(),
			}))
			Expect(mapWithNil).ShouldNot(ContainMap(map[string]any{
				"key2": nil,
			}))
		})
	})

	When("the expected map contains matchers as keys", func() {
		It("evaluates keys using matchers", func() {
			Expect(mapA).Should(ContainMap(map[any]string{
				MatchRegexp(`^key[0-9]$`): "val1",
			}))
			Expect(mapA).ShouldNot(ContainMap(map[any]string{
				MatchRegexp(`^key[0-9]$`): "nonexistent",
			}))
		})

		It("matches distinct actual entries for distinct expected matchers", func() {
			numMap := map[string]int{"a": 1, "b": 1}
			Expect(numMap).Should(ContainMap(map[any]int{
				MatchRegexp(`^a$`): 1,
				MatchRegexp(`^b$`): 1,
			}))
			// One entry cannot satisfy two expected entries
			singleEntryMap := map[string]int{"a": 1}
			Expect(singleEntryMap).ShouldNot(ContainMap(map[any]int{
				MatchRegexp(`^a$`): 1,
				ContainSubstring("a"): 1,
			}))
		})
	})

	When("matcher errors occur during evaluation", func() {
		It("continues evaluating remaining entries without failing prematurely", func() {
			// A key matcher that fails on string cannot evaluate numerical comparison, but does not abort
			testMap := map[any]any{
				"strKey": 1,
				2:        2,
			}
			Expect(testMap).Should(ContainMap(map[any]any{
				2: 2,
			}))
		})
	})

	if miter.HasIterators() {
		When("passed an iter.Seq2", func() {
			It("succeeds when actual is an iter.Seq2", func() {
				Expect(universalMapIter2).Should(ContainMap(map[string]int{
					"foo": 0,
					"bar": 42,
				}))
				Expect(universalMapIter2).ShouldNot(ContainMap(map[string]int{
					"foo": 0,
					"bar": 999,
				}))
				Expect(universalMapIter2).ShouldNot(ContainMap(map[string]int{
					"unknown": 0,
				}))
			})

			It("succeeds when expected is an iter.Seq2", func() {
				subIter := func(yield func(string, int) bool) {
					_ = yield("foo", 0) && yield("bar", 42)
				}
				Expect(universalMap).Should(ContainMap(subIter))

				mismatchIter := func(yield func(string, int) bool) {
					_ = yield("foo", 999)
				}
				Expect(universalMap).ShouldNot(ContainMap(mismatchIter))
			})

			It("succeeds when both actual and expected are iter.Seq2", func() {
				subIter := func(yield func(string, int) bool) {
					_ = yield("foo", 0) && yield("baz", 666)
				}
				Expect(universalMapIter2).Should(ContainMap(subIter))
			})

			It("handles empty iter.Seq2", func() {
				Expect(universalMapIter2).Should(ContainMap(emptyIter2))
				Expect(emptyIter2).Should(ContainMap(emptyIter2))
				Expect(emptyIter2).ShouldNot(ContainMap(universalMapIter2))
			})
		})
	}

	When("passed invalid actual or expected arguments", func() {
		It("fails when actual is not a map or iter.Seq2", func() {
			failures := InterceptGomegaFailures(func() {
				Expect("not a map").Should(ContainMap(mapB))
			})
			Expect(failures).Should(HaveLen(1))
			Expect(failures[0]).Should(ContainSubstring("ContainMap matcher expects a map/iter.Seq2"))

			failures = InterceptGomegaFailures(func() {
				Expect([]string{"a", "b"}).Should(ContainMap(mapB))
			})
			Expect(failures).Should(HaveLen(1))
			Expect(failures[0]).Should(ContainSubstring("ContainMap matcher expects a map/iter.Seq2"))

			failures = InterceptGomegaFailures(func() {
				Expect(nil).Should(ContainMap(mapB))
			})
			Expect(failures).Should(HaveLen(1))
			Expect(failures[0]).Should(ContainSubstring("ContainMap matcher expects a map/iter.Seq2"))
		})

		It("fails when expected is not a map or iter.Seq2", func() {
			failures := InterceptGomegaFailures(func() {
				Expect(mapA).Should(ContainMap("not a map"))
			})
			Expect(failures).Should(HaveLen(1))
			Expect(failures[0]).Should(ContainSubstring("ContainMap matcher expects a map/iter.Seq2"))

			failures = InterceptGomegaFailures(func() {
				Expect(mapA).Should(ContainMap([]string{"a"}))
			})
			Expect(failures).Should(HaveLen(1))
			Expect(failures[0]).Should(ContainSubstring("ContainMap matcher expects a map/iter.Seq2"))

			failures = InterceptGomegaFailures(func() {
				Expect(mapA).Should(ContainMap(nil))
			})
			Expect(failures).Should(HaveLen(1))
			Expect(failures[0]).Should(ContainSubstring("ContainMap matcher expects a map/iter.Seq2"))
		})
	})

	Describe("FailureMessage", func() {
		It("formats the failure message with missing entries", func() {
			failures := InterceptGomegaFailures(func() {
				Expect(map[string]int{"a": 1}).Should(ContainMap(map[string]int{"a": 1, "b": 2}))
			})
			Expect(failures).Should(HaveLen(1))
			Expect(failures[0]).Should(ContainSubstring("to contain map"))
			Expect(failures[0]).Should(ContainSubstring("the missing entries were"))
			Expect(failures[0]).Should(ContainSubstring(`"b": 2`))
		})
	})

	Describe("NegatedFailureMessage", func() {
		It("formats the negated failure message", func() {
			failures := InterceptGomegaFailures(func() {
				Expect(map[string]int{"a": 1}).ShouldNot(ContainMap(map[string]int{"a": 1}))
			})
			Expect(failures).Should(HaveLen(1))
			Expect(failures[0]).Should(ContainSubstring("not to contain map"))
		})
	})
})
