---
name: gstruct
description: Deep, partial matching of nested structs, slices, maps, and pointers with gstruct — MatchAllFields/MatchFields/Fields, MatchAllElements/MatchElements/Elements (idFn), MatchAllKeys/MatchKeys/Keys, PointTo, and the IgnoreExtras/IgnoreMissing/IgnoreUnexportedExtras/AllowDuplicates options, plus Ignore()/Reject(). Use when asserting against large or deeply nested data structures where you want to apply a different matcher to each field, element, or key.
---

# gstruct: matching complex data types

`gstruct` builds composite matchers that apply a *separate* matcher to each field of a struct, each element of a slice, each key of a map, or the target of a pointer. It is the tool for fuzzy-matching large, deeply nested values. Docs: <https://onsi.github.io/gomega/#gstruct-testing-complex-data-types>.

Dot-import gomega; import gstruct normally (or dot-import it too):

```go
import (
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
)
```

**The values in `Fields`/`Elements`/`Keys` are themselves matchers.** That is the whole point: compose with the core catalog (`gomega:matchers`) and combinators like `And`/`Or`/`Not`/`WithTransform` (`gomega:composing-matchers`). Examples below assume `gstruct` is dot-imported for brevity.

If you find yourself building an overly complex `gstruct` matcher, consider if you should be building a custom matcher (`gomega:custom-matchers`) instead.

## Structs — `Fields`

`Fields` is `map[string]types.GomegaMatcher` keyed by field name.

```go
actual := struct {
	A int
	B bool
	C string
}{5, true, "foo"}

Expect(actual).To(MatchAllFields(Fields{
	"A": BeNumerically("<", 10),
	"B": BeTrue(),
	"C": Equal("foo"),
}))
```

`MatchAllFields` requires an exact 1:1 mapping — every field must have a matcher and every matcher must map to a field. This is great for maintainability: adding or removing a struct field breaks the test until you update it.

To match a subset/superset, use `MatchFields(options, Fields{...})`:

```go
// IgnoreExtras: ignore struct fields with no matcher
Expect(actual).To(MatchFields(IgnoreExtras, Fields{
	"A": BeNumerically("<", 10),
	"B": BeTrue(), // no entry for C — fine
}))

// IgnoreMissing: ignore matchers with no corresponding field
Expect(actual).To(MatchFields(IgnoreMissing, Fields{
	"A": BeNumerically("<", 10),
	"B": BeTrue(),
	"C": Equal("foo"),
	"D": Equal("bar"), // ignored — actual has no field D
}))
```

`IgnoreUnexportedExtras` is a middle ground: it ignores only *unexported* extra fields (gstruct can't read them via reflect anyway) while still requiring all exported fields to be matched.

## Slices — `Elements` + an id function

`Elements` is `map[string]types.GomegaMatcher` keyed by a string id. You supply an `Identifier` (`func(element any) string`) that maps each element to its key.

```go
actual := []string{"A: foo bar baz", "B: once upon a time", "C: the end"}
id := func(element any) string { return string(element.(string)[0]) }

Expect(actual).To(MatchAllElements(id, Elements{
	"A": Not(BeZero()),
	"B": MatchRegexp("[A-Z]: [a-z ]+"),
	"C": ContainSubstring("end"),
}))
```

`MatchAllElements` requires a 1:1 mapping. Use `MatchElements(id, options, Elements{...})` with `IgnoreExtras`/`IgnoreMissing` to relax, exactly as with fields. `AllowDuplicates` lets several elements share one key/matcher (all of them must pass).

Index-based variants `MatchAllElementsWithIndex`/`MatchElementsWithIndex` take an `IdentifierWithIndex` (`func(index int, element any) string`); the built-in `IndexIdentity` just uses the index as the key.

## Maps — `Keys`

The `*Fields` API has a `*Keys` mirror for maps: `MatchAllKeys(Keys{...})` and `MatchKeys(options, Keys{...})`.

```go
actual := map[string]string{"A": "correct", "B": "incorrect"}

// MatchAllKeys requires every key be matched (this would fail without B):
Expect(actual).To(MatchAllKeys(Keys{
	"A": Equal("correct"),
	"B": Equal("incorrect"),
}))

// IgnoreMissing tolerates matchers for absent keys:
Expect(actual).To(MatchKeys(IgnoreMissing, Keys{
	"A": Equal("correct"),
	"B": Equal("incorrect"),
	"C": Equal("whatever"), // ignored — actual has no key C
}))
```

## Pointers — `PointTo`

`PointTo(matcher)` dereferences a pointer and applies `matcher` to the pointed-to value. It fails if the pointer is `nil`.

```go
foo := 5
Expect(&foo).To(PointTo(Equal(5)))
var bar *int
Expect(bar).NotTo(PointTo(BeNil())) // nil pointer fails PointTo
```

## `Ignore()` and `Reject()`

`Ignore()` always succeeds, `Reject()` always fails — use them as entries to skip a field/element or to assert one must never appear.

## Putting it all together

The matchers nest arbitrarily, mixing struct/slice/map/pointer matchers with the core catalog:

```go
coreID := func(element any) string { return strconv.Itoa(element.(CoreStats).Index) }

Expect(actual).To(MatchAllFields(Fields{
	"Name":      Ignore(),
	"StartTime": BeTemporally(">=", time.Now().Add(-100*time.Hour)),
	"CPU": PointTo(MatchAllFields(Fields{
		"Time":                 BeTemporally(">=", time.Now().Add(-time.Hour)),
		"UsageNanoCores":       BeNumerically("~", 1e9, 1e8),
		"UsageCoreNanoSeconds": BeNumerically(">", 1e6),
		"Cores": MatchElements(coreID, IgnoreExtras, Elements{
			"0": MatchAllFields(Fields{
				"Index":                Ignore(),
				"UsageNanoCores":       BeNumerically("<", 1e9),
				"UsageCoreNanoSeconds": BeNumerically(">", 1e5),
			}),
		}),
	})),
	"Memory": PointTo(MatchAllFields(Fields{
		"Time":            BeTemporally(">=", time.Now().Add(-time.Hour)),
		"AvailableBytes":  BeZero(),
		"UsageBytes":      BeNumerically(">", 5e6),
		"WorkingSetBytes": BeNumerically(">", 5e6),
	})),
	"Rootfs": Ignore(),
	"Logs":   Ignore(),
}))
```

## Gotchas

- **`MatchAllFields`/`MatchAllElements`/`MatchAllKeys` fail on missing OR extra entries.** They demand an exact 1:1 mapping. The instant you want a partial match, switch to `MatchFields`/`MatchElements`/`MatchKeys` and pass `IgnoreExtras`, `IgnoreMissing`, or both (`IgnoreExtras|IgnoreMissing`).
- **The id function must produce unique, stable keys.** Two elements mapping to the same key collide and fail unless you pass `AllowDuplicates`. Don't derive keys from data that changes between runs.
- **`Elements` keys are strings.** The id function returns a `string`, so numeric ids must be stringified (e.g. `strconv.Itoa(i)`), and the keys in your `Elements{}` literal must match those strings exactly.
- **`PointTo` rejects `nil` pointers** — a `nil` actual fails the match outright, before the inner matcher ever runs.
- Field/element/key matcher values are ordinary matchers, so compose freely with `gomega:matchers` and `gomega:composing-matchers`.
