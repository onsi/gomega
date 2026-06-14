---
name: matchers
description: The complete catalog of Gomega's built-in matchers, grouped by category — equivalence (Equal/BeEquivalentTo/BeComparableTo/BeIdenticalTo/BeAssignableToTypeOf), presence (BeNil/BeZero/BeEmpty), truthiness (BeTrue/BeFalse/BeTrueBecause), errors (HaveOccurred/Succeed/MatchError), channels (Receive/BeClosed/BeSent), files, strings/JSON/XML/YAML, collections (ContainElement/ConsistOf/HaveExactElements/HaveKey), structs (HaveField), numbers/times (BeNumerically/BeTemporally), values (HaveValue), HTTP responses, and panics. Use when you need to find or choose the right matcher for an assertion instead of defaulting to Equal.
---

# Gomega matcher catalog

Full reference: <https://onsi.github.io/gomega/#provided-matchers>. Assumes dot-import.

Three rules: **prefer the most specific matcher** — it produces far better failure messages. **Every matcher is negatable** (`NotTo`/`ShouldNot`). **Many matchers accept other matchers as arguments** (e.g. `ContainElement(ContainSubstring("x"))`) — compose freely. Anything taking a `format string, args ...any` runs `fmt.Sprintf` on it.

## Asserting Equivalence

- `Equal(expected)` — deep equality via `reflect.DeepEqual`; **type-strict** (actual and expected must be the same type). The default; reach for something more specific first.
- `BeEquivalentTo(expected)` — like `Equal` but **converts actual's type to expected's first**. Laxer and **risky** — `5.1` matches `BeEquivalentTo(5)` via truncation. Never use with numbers (use `BeNumerically`); fine for type aliases.
- `BeComparableTo(expected, options ...cmp.Option)` — deep equality via go-cmp (`github.com/google/go-cmp`); pass `cmp.Options` to ignore fields, compare unexported, set tolerances, etc.
- `BeIdenticalTo(expected)` — `==` identity; for primitives, or to assert two **pointers point to the same memory**.
- `BeAssignableToTypeOf(expected)` — succeeds if actual is assignable to a variable of expected's type. Asserts *type*, not value.

## Asserting Presence

- `BeNil()` — actual is `nil`. **Gotcha:** a non-nil interface holding a nil concrete pointer is *not* `nil`.
- `BeZero()` — actual is the zero value for its type (or `nil`).
- `BeEmpty()` — actual (`string`/array/map/chan/slice/iterator) has zero length.

## Asserting Truthiness

`BeTrue`/`BeFalse` require a `bool` (no "truthiness" — use `ShouldNot(BeNil())`/`ShouldNot(BeZero())` for presence).

- `BeTrue()` — actual is `true`. Weak failure message; prefer `BeTrueBecause`.
- `BeFalse()` — actual is `false`. Prefer `BeFalseBecause`.
- `BeTrueBecause(reason, args...)` — `BeTrue` with an explanatory message. **Best practice.**
- `BeFalseBecause(reason, args...)` — `BeFalse` with an explanatory message. **Best practice.**

## Asserting on Errors

Also surfaced in `gomega:assertions`, which covers the multi-return error idiom.

- `HaveOccurred()` — actual is a non-nil `error`. Idiom: `Expect(err).NotTo(HaveOccurred())`.
- `Succeed()` — actual error is `nil`. Idiom: `Expect(fn()).To(Succeed())` for funcs returning error first/only.
- `MatchError(expected, [funcDescription])` — **polymorphic**: `string` → `err.Error() == s`; `error` → `errors.Is` then `reflect.DeepEqual` against wrapped errors; matcher → applied to `err.Error()`; `func(error) bool` → predicate (**requires** the description second arg).
- `MatchErrorStrictly(expected)` — succeeds only if both non-nil and `errors.Is(actual, expected)`; no string fallback.

## Working with Channels

- `Receive([&val], [matcher])` — non-blocking: a value is ready to receive. `Receive(&val)` captures it into a pointer; `Receive(matcher)` asserts on the received value; `Receive(&val, matcher)` does both. Pairs with `Eventually`/`Consistently` → `gomega:async`.
- `BeClosed()` — actual is a closed channel. Reads from the channel to check; drain buffered channels first.
- `BeSent(value)` — non-blocking send of `value` onto actual succeeds (and actually sends).

## Working with files

Actual must be a filepath `string`.

- `BeAnExistingFile()` — a file exists at the path.
- `BeARegularFile()` — exists and is a regular file.
- `BeADirectory()` — exists and is a directory.

## Working with Strings, JSON and YAML

String matchers accept `string`/`[]byte`/`Stringer`. The `args...` forms run `fmt.Sprintf`.

- `ContainSubstring(substr, args...)` — actual contains the substring.
- `HavePrefix(prefix, args...)` — actual starts with the string.
- `HaveSuffix(suffix, args...)` — actual ends with the string.
- `MatchRegexp(regexp, args...)` — actual matches the regular expression.
- `MatchJSON(expected)` — actual and expected are the same JSON object (ignores whitespace/formatting/key order).
- `MatchXML(expected)` — actual and expected are the same XML object (ignores whitespace/formatting).
- `MatchYAML(expected)` — actual and expected are the same YAML object (ignores whitespace/formatting/key order).
- `HaveLen(count)` — string (or collection) has the given length.

## Working with Collections

Actual may be array/slice/map (and, on Go 1.23+, `iter.Seq`/`iter.Seq2` iterators). For maps, element matchers search **values**. Elements passed in may themselves be matchers.

- `HaveLen(count)` — length is `count`.
- `HaveCap(count)` — capacity is `count` (array/chan/slice).
- `BeEmpty()` — length zero.
- `ContainElement(element, [&pointer])` — contains a matching element. With a pointer second arg, **extracts** matches into it (scalar for one, slice/map for many).
- `ContainElements(elements...)` — contains all of the given elements (order-independent; extras allowed). Pass a single slice arg if needed.
- `ConsistOf(elements...)` — contains **precisely** these elements, order-independent (same length, no extras). vs `ContainElement(s)`: `ConsistOf` also checks length.
- `HaveExactElements(elements...)` — contains precisely these elements **in order** (array/slice).
- `BeElementOf(elements...)` — actual equals one of the given elements (always uses `Equal`).
- `BeKeyOf(map)` — actual equals one of the map's keys (always uses `Equal`).
- `HaveEach(element)` — every element matches (errors on empty collection).
- `HaveKey(key)` — map has a matching key.
- `HaveKeyWithValue(key, value)` — map has a matching key mapped to a matching value.

For deep/nested collection matching see `gomega:gstruct`; to compose element matchers see `gomega:composing-matchers`.

## Working with Structs

- `HaveField(field, value)` — struct's `field` matches `value`. `field` supports nested traversal (`"A.B.C"`) and zero-arg method calls (`"Method()"`, `"A.DOB.Year()"`). `value` may be a matcher. Missing field is an error.
- `HaveExistingField(field)` — struct has `field` regardless of value; combine with `And(HaveExistingField(f), HaveField(f, v))` or use as a filter.

For rich nested struct/slice/map matching see `gomega:gstruct`.

## Working with Numbers and Times

- `BeNumerically(comparator, expected, [threshold])` — type-agnostic numeric compare. Comparators: `"=="`, `">"`, `">="`, `"<"`, `"<="`, and `"~"` (**approximate** — within `threshold`, default `1e-8`). Use this for cross-type number equality.
- `BeTemporally(comparator, time, [threshold])` — `time.Time` compare. Same six comparators; `"~"` is within `threshold` (default `time.Millisecond`).

## Working with Values

- `HaveValue(matcher)` — dereferences pointers/interfaces (up to 31 levels) and applies `matcher` to the value; fails on nil. Lets one matcher work for both pointer and non-pointer actuals. Also in `gomega:composing-matchers`.
- `BeIdenticalTo(expected)` — `==` identity (see Equivalence above).

## Working with HTTP responses

Actual must be `*http.Response` or `*httptest.ResponseRecorder`.

- `HaveHTTPStatus(expected...)` — matches `StatusCode` (int) or `Status` (string); succeeds if any expected value matches.
- `HaveHTTPHeaderWithValue(key, value)` — header `key` matches `value` (string or matcher).
- `HaveHTTPBody(expected)` — response body matches (string, `[]byte`, or matcher called with `[]byte`). Reads and closes the body.

## Asserting on Panics

Actual must be a `func()` (no args, no returns).

- `Panic()` — invoking actual panics.
- `PanicWith(value)` — invoking actual panics with a matching value (`value` may be a matcher).

## Going further

- Composing/transforming matchers (`And`/`Or`/`Not`, `SatisfyAll`/`SatisfyAny`, `WithTransform`) → `gomega:composing-matchers`
- Writing your own matcher (`GomegaMatcher`, `gcustom`) → `gomega:custom-matchers`
- Deep, partial matching of nested structs/slices/maps → `gomega:gstruct`
