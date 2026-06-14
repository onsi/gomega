---
name: gleak
description: gleak goroutine leak detection — capture a Goroutines() snapshot before a test, then Eventually(Goroutines).ShouldNot(HaveLeaked(snapshot)) to assert none leaked, with the BeforeEach/AfterEach/DeferCleanup pattern, ignoring matchers IgnoringTopFunction/IgnoringInBacktrace/IgnoringGoroutines/IgnoringCreator, well-known non-leaky goroutines, goroutine IDs, ReportFilenameWithPath, and the Ginkgo -p IgnoreGinkgoParallelClient gotcha. Use when a test must verify goroutines started during the test have all wound down and nothing leaked.
---

# gleak: detecting leaked goroutines

`gleak` discovers all running goroutines and fails a test if any "leaked" — i.e. are still running after the test that aren't well-known framework/runtime goroutines and aren't on your ignore list. Docs: <https://onsi.github.io/gomega/#gleak-finding-leaked-goroutines>. Cross-refs: `gomega:async`, `gomega:matchers`.

> gleak is an experimental Gomega package.

## Import

```go
import (
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gleak" // often dot-imported too, for HaveLeaked/Goroutines/Ignoring*
)
```

## The canonical pattern

Capture a baseline snapshot **before** the test; assert no leak **after**. The cleanup runs `Eventually(Goroutines).ShouldNot(HaveLeaked(snapshot))`:

```go
BeforeEach(func() {
	goods := gleak.Goroutines()          // snapshot of "good" goroutines, taken now
	DeferCleanup(func() {
		Eventually(gleak.Goroutines).ShouldNot(gleak.HaveLeaked(goods))
	})
})
```

Equivalent BeforeEach/AfterEach form (carry the snapshot in a closure variable):

```go
var goods []gleak.Goroutine
BeforeEach(func() { goods = gleak.Goroutines() })
AfterEach(func()  { Eventually(gleak.Goroutines).ShouldNot(gleak.HaveLeaked(goods)) })
```

Simplest possible form — no snapshot, relies only on the built-in well-known list:

```go
AfterEach(func() {
	Eventually(gleak.Goroutines).ShouldNot(gleak.HaveLeaked())
})
```

## Gotchas

**ALWAYS use `Eventually`, never `Expect`.** Goroutines wind down *asynchronously*; a synchronous `Expect` races their shutdown and produces false positives. `Eventually` re-polls until they terminate or it times out (default 1s / 10ms poll — see `gomega:async` to tune).

**Pass `Goroutines`, not `Goroutines()`.** Note the missing `()`. `Eventually` must call it *repeatedly* on each poll. Writing `Eventually(gleak.Goroutines())` snapshots once and defeats the retry. (Calling `gleak.Goroutines()` *with* parens is correct only when *taking a baseline snapshot* to pass into `HaveLeaked`.)

**Capture the baseline at the right time** — in `BeforeEach`, before the test spins anything up, so genuinely pre-existing goroutines are filtered out.

**`HaveLeaked` succeeds when goroutines leaked.** A *success* is a *failure* of your test — so it's almost always used with `ShouldNot`/`NotTo`. The built-in well-known list is always applied and cannot be disabled.

## `HaveLeaked([ignoring...])`

```go
gleak.HaveLeaked(NONLEAKY1, NONLEAKY2, ...)
```

After filtering out well-known goroutines and everything matched by the `ignoring` args, it matches if any goroutines remain. Each arg is either a goroutine matcher or a shorthand:

- `"foo.bar"` → `IgnoringTopFunction("foo.bar")` (exact topmost-function name)
- `"foo.bar..."` → top function name starts with prefix `foo.bar.`
- `"foo.bar [chan receive]"` → exact top function **and** goroutine state begins with `chan receive`
- `[]Goroutine` (a snapshot) → `IgnoringGoroutines(snapshot)`, filtered by goroutine ID
- any `GomegaMatcher` that operates on an actual of type `gleak.Goroutine` (e.g. `HaveField`, `WithTransform`)

## Goroutine matchers (for the ignoring list)

- **`IgnoringTopFunction(name)`** — matches a goroutine whose topmost stack function is `name`; supports the `"name"`, `"prefix..."`, and `"name [state]"` forms above.
- **`IgnoringInBacktrace(substr)`** — matches if `substr` appears *anywhere* in the backtrace (lazy `strings.Contains`).
- **`IgnoringGoroutines(snapshot)`** — matches goroutines that are elements of `snapshot`, compared by goroutine ID.
- **`IgnoringCreator(name)`** — matches by the name of the function that *created* the goroutine; supports `"name"` and `"prefix..."` forms.

```go
Eventually(gleak.Goroutines).ShouldNot(gleak.HaveLeaked(
	gleak.IgnoringTopFunction("github.com/me/pkg.worker"),
	gleak.IgnoringInBacktrace("github.com/some/dep.background"),
	goods, // the snapshot, as IgnoringGoroutines shorthand
))
```

## Ginkgo `-p` (parallel)

**Running with `ginkgo -p` adds a background goroutine** for Ginkgo↔package communication that will look leaked. Call `gleak.IgnoreGinkgoParallelClient()` at the start of each package's `BeforeSuite` so gleak adds it to the ignore list:

```go
var _ = BeforeSuite(func() {
	gleak.IgnoreGinkgoParallelClient()
})
```

## Well-known non-leaky goroutines (ignored by default)

Always filtered, by topmost function name (or backtrace), so you don't list them:

- signal handling: `os/signal.signal_recv`, `os/signal.loop`, `runtime.ensureSigM`
- Go `testing`: `testing.RunTests`, `testing.(*T).Run`, `testing.(*T).Parallel` (all `[chan receive]`)
- Ginkgo internals: `(*Suite).runNode`, the interrupt-handler/progress-signal goroutines, output-interceptor `ResumeIntercepting`, and the v1 spec-runner interrupt handler
- anything with `runtime.ReadTrace` in its backtrace

## Reporting

On a leak, gleak prints only the *leaked* goroutines (not all, unlike a panic), in a compact one-line-per-frame form. This output is **not** subject to `format.MaxLength`. By default locations show package + file + line (`main.foo.func1() at foo/bar.go:123`). Set `gleak.ReportFilenameWithPath = true` for full absolute paths.

## Goroutine IDs

`gleak.Goroutine` values carry a goroutine ID (`goid`) parsed from the runtime stack; `IgnoringGoroutines` matches snapshots by these IDs. IDs are not reused (barring 64-bit counter wraparound) but are not densely sequential. Use them for testing/debugging only — never for program logic.
