---
name: assertions
description: Write correct synchronous Gomega assertions — Expect/Ω notation, the To/NotTo/ToNot/Should/ShouldNot equivalences, the multi-return error idiom, Succeed vs HaveOccurred, the .Error() chaining form, annotating assertions (format-string and func()string), tuning failure output via the format subpackage (MaxLength/MaxDepth/UseStringerRepresentation/GomegaStringer/TruncatedDiff/RegisterCustomFormatter/format.Object), and asserting inside helper functions with GinkgoHelper/WithOffset/ExpectWithOffset, NewWithT(t) for plain testing, and the g Gomega callback. Use when writing or reviewing synchronous (non-polling) Gomega assertions.
---

# Synchronous assertions

The core mechanics of asserting *now* (vs. polling — see `gomega:async`). Assumes a dot-import (`. "github.com/onsi/gomega"`). Narrative docs: <https://onsi.github.io/gomega/#making-assertions>. For the big picture and skill map see `gomega:overview`.

## The two notations are identical

```go
Expect(ACTUAL).To(Equal(EXPECTED))      // Expect notation
Expect(ACTUAL).NotTo(Equal(EXPECTED))
Expect(ACTUAL).ToNot(Equal(EXPECTED))   // ToNot == NotTo

Ω(ACTUAL).Should(Equal(EXPECTED))       // Ω notation (⌥z on macOS)
Ω(ACTUAL).ShouldNot(Equal(EXPECTED))
```

`To` == `Should`, and `NotTo` == `ToNot` == `ShouldNot` — pure syntactic sugar. `ACTUAL` (left) is anything; the matcher (right) must satisfy `GomegaMatcher`. Pick one notation per codebase and stay consistent. Matchers are values — see `gomega:matchers`.

## The multi-return error idiom

`Expect`/`Ω` accept **multiple arguments**. The matcher runs against the **first**; the assertion **fails unless every trailing argument is nil or zero-valued**. This collapses Go's `(value, error)` pattern:

```go
// instead of:
result, err := DoSomethingHard()
Expect(err).NotTo(HaveOccurred())
Expect(result).To(Equal("foo"))

// write:
Expect(DoSomethingHard()).To(Equal("foo"))   // passes only if return is ("foo", nil)
```

## Succeed and HaveOccurred

For a function returning **only** `error`, use `Succeed`:

```go
Expect(DoSomethingSimple()).To(Succeed())     // func() error
Expect(DoSomethingSimple()).NotTo(Succeed())
```

For an `error` value you already hold, use `HaveOccurred`:

```go
err := DoSomethingSimple()
Expect(err).NotTo(HaveOccurred())
Expect(err).To(HaveOccurred())
```

Use `Succeed` when calling the function inline; use `HaveOccurred` when you have an `err` variable. To assert on the error's *content* (not just its presence), reach for `MatchError` and friends in `gomega:matchers`.

**Don't** use `Succeed` with a multi-return function. The matcher only sees the *first* return value; the rest are consumed by the multi-return idiom above. `Expect(DoSomethingHard()).To(Succeed())` matches against the `string`, not the `error`, and `Expect(DoSomethingHard()).NotTo(Succeed())` can **never pass**.

## The .Error() chaining form

To assert on the trailing error of a **multi-return** function while ignoring the other returns, chain `.Error()`:

```go
Expect(MultipleReturnValuesFunc()).Error().To(HaveOccurred())
Expect(MultipleReturnValuesFunc()).Error().NotTo(HaveOccurred())
```

`.Error()` retargets the assertion at the **last** return value (the error). The `To(HaveOccurred())` form *additionally* asserts that all the *other* returns are zero-valued; the `NotTo(HaveOccurred())` form lets the other returns be anything. The plain alternative is to capture and assert manually:

```go
_, _, _, err := MultipleReturnValuesFunc()
Expect(err).To(HaveOccurred())
```

## Annotating assertions

Pass a format string (with `fmt.Sprintf` args) or a `func() string` **after** the matcher. It's printed alongside the failure message to add context:

```go
Expect(ACTUAL).To(Equal(EXPECTED), "my annotation %d", foo)
Expect(ACTUAL).To(Equal(EXPECTED), func() string { return expensive() })
```

The `func() string` form is **lazily evaluated** — only called if the assertion fails — so use it for any annotation that's expensive to build.

## Adjusting failure output

On failure Gomega prints a recursive rendering of the objects involved, produced by the `format` subpackage. Tune it via package-level globals (set once, e.g. in a test helper or `TestMain`):

```go
import "github.com/onsi/gomega/format"

format.MaxLength = 4000   // truncate rendered output to N chars; 0 disables truncation
format.MaxDepth = 10      // max recursion depth into nested structures
format.TruncatedDiff = true   // for long strings, show only where they differ (false = full strings)
format.UseStringerRepresentation = false  // true => call String()/GoString() on Stringer/GoStringer types
format.PrintContextObjects = false        // true => print contents of context.Context objects
```

`UseStringerRepresentation` defaults to `false` on purpose: a `String()` method often hides fields you need to diagnose a failure. Leave it off unless you know the stringer is complete.

**`GomegaStringer` interface.** Implement `GomegaString() string` on a type and Gomega always uses it for that type's rendering, regardless of `UseStringerRepresentation`. Best practice: define it in a `_test.go` helper file so you don't leak it into your package's exported API.

```go
func (w Widget) GomegaString() string { return fmt.Sprintf("Widget<%s>", w.id) }
```

**Custom formatters.** Register a `format.CustomFormatter` (`func(value any) (string, bool)`) — return `(rendered, true)` to handle a value, or `("", false)` to pass. Custom formatters take precedence over `GomegaStringer` and `UseStringerRepresentation`, and their output is not truncated:

```go
key := format.RegisterCustomFormatter(myFormatter)
defer format.UnregisterCustomFormatter(key)
```

**`format.Object`.** Reuse Gomega's renderer directly (the int is indentation depth):

```go
fmt.Println(format.Object(theThing, 1))
```

## Making assertions inside helper functions

Factoring assertions into a helper reduces boilerplate, but a naive helper reports failures at the line *inside the helper*, not the caller. Fix the reported line number with one of these.

**With Ginkgo — `GinkgoHelper()` (recommended).** Call it first thing in the helper; Ginkgo then attributes failures to the call site. It nests cleanly across multiple helper layers:

```go
func expectComponents(te TurboEncabulator, components ...string) {
    GinkgoHelper()
    got, err := te.GetComponents()
    Expect(err).NotTo(HaveOccurred())
    Expect(got).To(ConsistOf(components))
}
```

**Without Ginkgo — `WithOffset` / `ExpectWithOffset`.** The offset is how many stack frames to skip; `1` points at the helper's caller. These are equivalent:

```go
Expect(err).WithOffset(1).NotTo(HaveOccurred())
ExpectWithOffset(1, err).NotTo(HaveOccurred())
```

(`Eventually`/`Consistently` have the same `WithOffset` / `…WithOffset` forms — see `gomega:async`.)

**Plain `testing` — `NewWithT(t)`.** Wrap a `*testing.T` to get a `*WithT` carrying `Expect`, `Eventually`, and `Consistently`. Create a fresh one per test (no global fail handler needed):

```go
func TestFarmHasCow(t *testing.T) {
    g := NewWithT(t)
    f := farm.New([]string{"Cow", "Horse"})
    g.Expect(f.HasCow()).To(BeTrue(), "farm should have a cow")
}
```

**The `g Gomega` callback pattern.** A helper can accept a `Gomega` and assert through it, decoupling the helper from how failures are reported:

```go
func expectValidWidget(g Gomega, w Widget) {
    g.Expect(w.ID).NotTo(BeEmpty())
    g.Expect(w.Size).To(BeNumerically(">", 0))
}

g := NewWithT(t)        // or the g passed into an Eventually callback
expectValidWidget(g, w)
```

This is the **same `g Gomega`** that `Eventually(func(g Gomega) {...})` passes into polled callbacks (`gomega:async`) — so a `g Gomega` helper works in both synchronous and polled contexts. To build a `Gomega` wired to your own fail handler, use `NewGomega(failHandler)`.

**Don't** call the global `Expect`/`Ω` inside an `Eventually`/`Consistently` callback — those failures won't be intercepted by the poller. Always assert through the passed-in `g Gomega` there. See `gomega:async`.

For richer, reusable assertions consider a real matcher instead of a helper — see `gomega:custom-matchers`.
