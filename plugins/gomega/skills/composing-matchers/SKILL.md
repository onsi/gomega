---
name: composing-matchers
description: Build compound Gomega assertions by combining matchers — And/SatisfyAll (all pass), Or/SatisfyAny (any pass), Not (negate), WithTransform to map the actual before matching, Satisfy for an ad-hoc predicate, HaveValue to dereference pointers/interfaces, HaveField for struct fields and method results, HaveEach for every element, plus the matchers-as-arguments idiom that lets you nest matchers inside ContainElement/ConsistOf/HaveKeyWithValue/Receive. Use when one Expect needs several requirements at once, or you want to assert deep into a value without writing a custom matcher.
---

# Composing and transforming matchers

Gomega matchers are values, so you can wrap and nest them to express several requirements in a single `Expect`/`Eventually`. Start with `gomega:matchers` for the catalog; reach here when one assertion needs to do more. Docs: <https://onsi.github.io/gomega/#composing-matchers>.

## Logical combinators

```go
Expect(n).To(And(BeNumerically(">", 0), BeNumerically("<", 10)))   // all must pass
Expect(msg).To(Or(Equal("Success"), MatchRegexp(`^Error .+$`)))    // any may pass
Expect(s).To(Not(BeEmpty()))                                        // negate
```

- `And(ms ...GomegaMatcher)` / `SatisfyAll(...)` — identical; succeeds only if every matcher passes. Short-circuits on the first failure.
- `Or(ms ...GomegaMatcher)` / `SatisfyAny(...)` — identical; succeeds if any matcher passes. Short-circuits on the first success.
- `Not(matcher GomegaMatcher)` — logical negation; usually clearer than the `.NotTo(...)` form when nested.

These are the building blocks for lightweight named matchers:

```go
func BeBetween(min, max int) GomegaMatcher {
    return SatisfyAll(BeNumerically(">", min), BeNumerically("<", max))
}
Expect(n).To(BeBetween(0, 10))
```

## WithTransform — map the actual, then match

`WithTransform(transform any, matcher GomegaMatcher)` applies `transform` to `ACTUAL` and matches the result.

```go
Expect(element).To(WithTransform(func(e Element) Color { return e.Color }, Equal(BLUE)))
```

**Transform-function rules:** it must take exactly one argument (the actual) and return either one value, or `(value, error)`. Returning an error fails the assertion gracefully instead of panicking — useful for transforms that accept several input types:

```go
func HaveSprocketName(name string) GomegaMatcher {
    return WithTransform(func(actual any) (string, error) {
        switch s := actual.(type) {
        case *Sprocket: return s.Name, nil
        case Sprocket:  return s.Name, nil
        default:        return "", fmt.Errorf("expected a Sprocket, got %T", actual)
        }
    }, Equal(name))
}
```

This is the standard pattern for distilling a `WithTransform` into a reusable named matcher.

## Satisfy — an ad-hoc boolean predicate

`Satisfy(predicate any)` passes when `predicate(ACTUAL)` returns `true`. The predicate takes one argument and returns one `bool`.

```go
Expect(n).To(Satisfy(func(i int) bool { return i%2 == 0 }))
```

Reach for `Satisfy` for a one-off check; `WithTransform` when you want the inner matcher's richer failure message. For anything reused, write a real matcher → `gomega:custom-matchers`.

## HaveValue — dereference before matching

`HaveValue(matcher GomegaMatcher)` dereferences `ACTUAL` through pointers/interfaces (up to 31 levels) and applies `matcher` to the underlying value, so the same matcher works whether the actual is a value or a pointer. Fails on `nil`.

```go
i := 42
Expect(&i).To(HaveValue(Equal(42)))
Expect(i).To(HaveValue(Equal(42)))
```

**`Not(HaveValue(...))` does not suppress the nil error.** To accept `nil` while rejecting a specific value, use `Or(BeNil(), Not(HaveValue(...)))`.

## HaveField — assert on a struct field or method result

`HaveField(field string, value any)` extracts a field (or method result) from a struct and matches it. If `value` is a matcher it is used directly; otherwise it is wrapped in `Equal`.

```go
Expect(book).To(HaveField("Title", "Les Miserables"))
Expect(book).To(HaveField("Title", ContainSubstring("Les Mis")))   // value can be a matcher
Expect(book).To(HaveField("Author.Name", "Victor Hugo"))           // nested with "."
Expect(book).To(HaveField("Author.DOB.Year()", BeNumerically("<", 1900))) // "()" calls a method
```

The `()` suffix invokes a method that takes no arguments and returns exactly one value. **A missing field is an error, not a failed match.** To guard, pair with `HaveExistingField(field)`: `And(HaveExistingField("X"), HaveField("X", VALUE))` — handy when reusing `HaveField` as a `ContainElement` filter.

## HaveEach — every element must match

`HaveEach(element any)` passes when every element of an array, slice, map (values), or `iter.Seq`/`iter.Seq2` matches. A plain value is wrapped in `Equal`; pass a matcher for richer checks.

```go
Expect([]string{"Foo", "FooBar"}).To(HaveEach(ContainSubstring("Foo")))
```

**An empty (or nil) collection is an error** — it is ambiguous whether "each" holds. If empty is acceptable, use `Or(BeEmpty(), HaveEach(...))`.

## The big idea: matchers as arguments

Nearly every collection/channel matcher accepts matchers where it accepts values, so you compose by nesting rather than learning new operators:

```go
Expect(xs).To(ContainElement(HaveField("Name", "gomega")))
Expect(xs).To(ContainElement(HaveField("Name", BeKeyOf(names)), &found)) // matcher + capture pointer
Expect(books).To(ConsistOf(
    HaveField("Author.Name", "Victor Hugo"),
    And(HaveField("Pages", BeNumerically(">", 100)), HaveField("InPrint", BeTrue())),
))
Expect(m).To(HaveKeyWithValue("user", HaveField("Active", BeTrue())))
Eventually(ch).Should(Receive(HaveField("Kind", "sesame")))
```

This composes cleanly with `Eventually`/`Consistently` too → `gomega:async`.

## When to switch to gstruct

For deep, partial matching across nested structs, slices, and maps, hand-rolled `HaveField`/`And`/`HaveEach` trees get noisy. `gstruct`'s `MatchFields`/`MatchAllElements`/`MatchKeys` are purpose-built for that shape → `gomega:gstruct`.
