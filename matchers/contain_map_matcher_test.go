package matchers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
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
				MatchRegexp(`^a$`):    1,
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
		It("errors when actual is not a map or iter.Seq2", func() {
			success, err := ContainMap(mapB).Match("not a map")
			Expect(success).Should(BeFalse())
			Expect(err).Should(MatchError("ContainMap matcher expects a map/iter.Seq2.  Got:\n    <string>: not a map"))

			_, err = ContainMap(mapB).Match([]string{"a", "b"})
			Expect(err).Should(MatchError(HavePrefix("ContainMap matcher expects a map/iter.Seq2.  Got:\n")))

			_, err = ContainMap(mapB).Match(nil)
			Expect(err).Should(MatchError("ContainMap matcher expects a map/iter.Seq2.  Got:\n    <nil>: nil"))
		})

		It("errors when actual is a pointer to a map", func() {
			success, err := ContainMap(mapB).Match(&mapA)
			Expect(success).Should(BeFalse())
			Expect(err).Should(MatchError(HavePrefix("ContainMap matcher expects a map/iter.Seq2.  Got:\n    <*map[string]string | 0x")))
		})

		It("errors, with a different message, when expected is not a map or iter.Seq2", func() {
			success, err := ContainMap("not a map").Match(mapA)
			Expect(success).Should(BeFalse())
			Expect(err).Should(MatchError("ContainMap matcher expects the expected value to be a map/iter.Seq2.  Got:\n    <string>: not a map"))

			_, err = ContainMap([]string{"a"}).Match(mapA)
			Expect(err).Should(MatchError(HavePrefix("ContainMap matcher expects the expected value to be a map/iter.Seq2.  Got:\n")))

			_, err = ContainMap(nil).Match(mapA)
			Expect(err).Should(MatchError("ContainMap matcher expects the expected value to be a map/iter.Seq2.  Got:\n    <nil>: nil"))

			_, err = ContainMap(&mapB).Match(mapA)
			Expect(err).Should(MatchError(HavePrefix("ContainMap matcher expects the expected value to be a map/iter.Seq2.  Got:\n    <*map[string]string | 0x")))
		})
	})

	When("a value matcher for a plain key errors", func() {
		It("returns the error", func() {
			success, err := ContainMap(map[string]any{"a": BeNumerically("~~", 1)}).Match(map[string]int{"a": 1})
			Expect(success).Should(BeFalse())
			Expect(err).Should(MatchError("ContainMap's value matcher for key\n    <string>: a\nfailed with:\n    Unknown comparator: ~~"))
		})

		It("reports the same error regardless of map iteration order", func() {
			for range 50 {
				_, err := ContainMap(map[string]any{
					"c": BeNumerically("~~", 3),
					"a": BeNumerically("~~", 1),
					"b": BeNumerically("~~", 2),
				}).Match(map[string]int{"a": 1, "b": 2, "c": 3})
				Expect(err).Should(MatchError(HavePrefix("ContainMap's value matcher for key\n    <string>: a\n")))
			}
		})

		It("does not run the value matcher, or error, when the key is missing", func() {
			success, err := ContainMap(map[string]any{"b": BeNumerically("~~", 1)}).Match(map[string]int{"a": 1})
			Expect(success).Should(BeFalse())
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Describe("matching each expected entry to a distinct actual entry", func() {
		It("does not let a plain key and a matcher key claim the same actual entry", func() {
			expected := map[any]any{
				"a":                  1,
				MatchRegexp(`^a`):    1,
				HavePrefix("b"):      BeNumerically(">", 0),
				ContainSubstring(""): 3,
			}
			for range 100 {
				Expect(map[string]int{"a": 1, "b": 2, "c": 3}).ShouldNot(ContainMap(expected))
				Expect(map[string]int{"a": 1, "ab": 1, "b": 2, "c": 3}).Should(ContainMap(expected))
				Expect(map[string]int{"a": 1, "ab": 2, "b": 2, "c": 3}).ShouldNot(ContainMap(expected))
			}
		})

		It("finds an assignment regardless of map iteration order when matcher keys overlap", func() {
			actual := map[string]int{"a": 1, "ab": 1, "abc": 1}
			expected := map[any]any{
				MatchRegexp(`^a`):    1,
				MatchRegexp(`^ab`):   1,
				MatchRegexp(`^abc$`): 1,
				MatchRegexp(`^abcd`): BeNil(),
			}
			for range 100 {
				Expect(actual).Should(ContainMap(map[any]any{
					MatchRegexp(`^a`):    1,
					MatchRegexp(`^ab`):   1,
					MatchRegexp(`^abc$`): 1,
				}))
				Expect(actual).ShouldNot(ContainMap(expected))
			}
		})

		It("compares plain keys with Equal semantics, including for pointer and named-type keys", func() {
			type Pair [2]int
			Expect(map[Pair]int{{1, 2}: 3}).Should(ContainMap(map[Pair]int{{1, 2}: 3}))
			Expect(map[Pair]int{{1, 2}: 3}).ShouldNot(ContainMap(map[[2]int]int{{1, 2}: 3}))
			Expect(map[any]int{1: 1}).Should(ContainMap(map[any]int{1: 1}))
			Expect(map[any]int{1: 1}).ShouldNot(ContainMap(map[any]int{int64(1): 1}))
			Expect(map[string]int64{"a": 1}).ShouldNot(ContainMap(map[string]int{"a": 1}))
			Expect(map[string]int64{"a": 1}).Should(ContainMap(map[string]any{"a": BeNumerically("==", 1)}))
		})

		if miter.HasIterators() {
			It("handles repeated keys in iter.Seq2s", func() {
				actual := func(yield func(string, int) bool) {
					_ = yield("a", 1) && yield("a", 2) && yield("b", 3)
				}
				expected := func(yield func(string, any) bool) {
					_ = yield("a", BeNumerically(">=", 1)) && yield("a", 1)
				}
				Expect(actual).Should(ContainMap(expected))

				expected = func(yield func(string, any) bool) {
					_ = yield("a", 1) && yield("a", 1)
				}
				Expect(actual).ShouldNot(ContainMap(expected))
				Expect(map[string]int{"a": 1}).ShouldNot(ContainMap(expected))

				expected = func(yield func(string, any) bool) {
					_ = yield("a", 2) && yield("b", 3)
				}
				Expect(actual).Should(ContainMap(expected))
			})
		}
	})

	When("used as a value matcher inside ContainMap", func() {
		It("matches nested maps partially", func() {
			config := map[string]any{
				"name": "gomega",
				"server": map[string]any{
					"port": 8080,
					"tls":  true,
				},
			}
			Expect(config).Should(ContainMap(map[string]any{
				"server": ContainMap(map[string]any{"port": 8080}),
			}))
			Expect(config).ShouldNot(ContainMap(map[string]any{
				"server": ContainMap(map[string]any{"port": 9090}),
			}))
		})

		It("nests the inner failure message", func() {
			actual := map[string]any{"server": map[string]int{"port": 8080}}
			m := ContainMap(map[string]any{"server": ContainMap(map[string]int{"port": 9090})})
			Expect(m.Match(actual)).Should(BeFalse())
			Expect(m.FailureMessage(actual)).Should(HaveSuffix(`the mismatched values were
    <string>: server
        Expected
            <map[string]int | len:1>: {"port": 8080}
        to contain map
            <map[string]int | len:1>: {"port": 9090}
        the mismatched values were
            <string>: port
                Expected
                    <int>: 8080
                to equal
                    <int>: 9090`))
		})
	})

	Describe("FailureMessage", func() {
		It("reports missing keys", func() {
			m := ContainMap(map[string]int{"b": 2})
			Expect(m.Match(map[string]int{"a": 1})).Should(BeFalse())
			Expect(m.FailureMessage(map[string]int{"a": 1})).Should(Equal(`Expected
    <map[string]int | len:1>: {"a": 1}
to contain map
    <map[string]int | len:1>: {"b": 2}
the missing keys were
    <[]string | len:1, cap:1>: ["b"]`))
		})

		It("reports keys that are present with a mismatched value, with the value matcher's failure message", func() {
			m := ContainMap(map[string]int{"a": 2})
			Expect(m.Match(map[string]int{"a": 1})).Should(BeFalse())
			Expect(m.FailureMessage(map[string]int{"a": 1})).Should(Equal(`Expected
    <map[string]int | len:1>: {"a": 1}
to contain map
    <map[string]int | len:1>: {"a": 2}
the mismatched values were
    <string>: a
        Expected
            <int>: 1
        to equal
            <int>: 2`))

			m = ContainMap(map[string]any{"a": BeNumerically(">", 5)})
			Expect(m.Match(map[string]int{"a": 1})).Should(BeFalse())
			Expect(m.FailureMessage(map[string]int{"a": 1})).Should(HaveSuffix(`
the mismatched values were
    <string>: a
        Expected
            <int>: 1
        to be >
            <int>: 5`))
		})

		It("reports matcher keys that could not be matched", func() {
			keyMatcher := MatchRegexp(`^b`)
			expected := map[any]int{keyMatcher: 1}
			m := ContainMap(expected)
			Expect(m.Match(map[string]int{"a": 1})).Should(BeFalse())
			Expect(m.FailureMessage(map[string]int{"a": 1})).Should(Equal(`Expected
    <map[string]int | len:1>: {"a": 1}
to contain map
` + format.Object(expected, 1) + `
the entries with no matching (distinct) entry in actual were
` + format.Object(keyMatcher, 1) + `
        <int>: 1`))
		})

		It("reports every kind of failure, sorted by key", func() {
			actual := map[string]int{"a": 1, "b": 2, "c": 3}
			m := ContainMap(map[any]any{
				"a":             1,
				"z":             1,
				"y":             2,
				"c":             4,
				"b":             BeNumerically("<", 0),
				HavePrefix("x"): 3,
			})
			Expect(m.Match(actual)).Should(BeFalse())
			Expect(m.FailureMessage(actual)).Should(MatchRegexp(`(?s)
the missing keys were
    <\[\]interface {} \| len:2, cap:2>: \[<string>"y", <string>"z"\]
the mismatched values were
    <string>: b
        Expected
            <int>: 2
        to be <
            <int>: 0
    <string>: c
        Expected
            <int>: 3
        to equal
            <int>: 4
the entries with no matching \(distinct\) entry in actual were
    <\*matchers.HavePrefixMatcher \| 0x[0-9a-f]+>: {Prefix: "x", Args: nil}
        <int>: 3$`))
		})
	})

	Describe("NegatedFailureMessage", func() {
		It("formats the negated failure message", func() {
			m := ContainMap(map[string]int{"a": 1})
			Expect(m.Match(map[string]int{"a": 1, "b": 2})).Should(BeTrue())
			Expect(m.NegatedFailureMessage(map[string]int{"a": 1})).Should(Equal(`Expected
    <map[string]int | len:1>: {"a": 1}
not to contain map
    <map[string]int | len:1>: {"a": 1}`))
		})
	})
})
