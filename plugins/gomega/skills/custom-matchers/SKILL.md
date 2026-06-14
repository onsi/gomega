---
name: custom-matchers
description: Writing your own Gomega matchers — the GomegaMatcher interface (Match/FailureMessage/NegatedFailureMessage), gcustom.MakeMatcher with message templates and template data, the format package helpers (format.Message/format.Object), MatchMayChangeInTheFuture and StopTrying for Eventually/Consistently, and how to test and package custom matchers. Use when a built-in or composed matcher can't express your domain assertion and you need to build one.
---

# Writing your own matchers

When no built-in matcher fits and composition (`And`/`Or`/`WithTransform`/`SatisfyAll` — see `gomega:composing-matchers`) can't express it cleanly, write a custom matcher. **Reach for `gcustom.MakeMatcher` first** — it's the modern, low-boilerplate path. Drop to a hand-written type only when you need full control. Docs: <https://onsi.github.io/gomega/#adding-your-own-matchers>.

## gcustom.MakeMatcher — the recommended path

`gcustom.MakeMatcher(matchFunc, ...)` builds a full `types.GomegaMatcher` from one function. The match function must be `func(actual T) (bool, error)`:

```go
import (
	"github.com/onsi/gomega/gcustom"
	"github.com/onsi/gomega/types"

	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

func RepresentJSONifiedObject(expected any) types.GomegaMatcher {
	return gcustom.MakeMatcher(func(response *http.Response) (bool, error) {
		ptr := reflect.New(reflect.TypeOf(expected)).Interface()
		if err := json.NewDecoder(response.Body).Decode(ptr); err != nil {
			return false, fmt.Errorf("failed to decode JSON: %w", err)
		}
		decoded := reflect.ValueOf(ptr).Elem().Interface()
		return reflect.DeepEqual(decoded, expected), nil
	}).WithTemplate("Expected:\n{{.FormattedActual}}\n{{.To}} contain the JSON representation of\n{{format .Data 1}}").WithTemplateData(expected)
}
```

Then `Expect(resp).To(RepresentJSONifiedObject(book))` just works.

**Typed match funcs get free type-checking.** Because the func takes `*http.Response`, `gcustom` rejects any other actual with a clear error *before* calling your code. Use `func(actual any) (bool, error)` only if you want to handle multiple types or do your own type checks.

**Return `(false, err)` for bad input, not a panic.** A non-nil error fails the assertion in *both* the `To` and `NotTo` directions — you can't accidentally pass a negated assertion by feeding garbage.

### Messages: WithMessage vs WithTemplate

- `.WithMessage("contain the JSON representation")` — simplest. Renders `Expected:\n<actual>\nto contain the JSON representation` (and `not to` when negated).
- `.WithTemplate(tmpl, optionalData)` — full control. `WithTemplate(tmpl, data)` is shorthand for `.WithTemplate(tmpl).WithTemplateData(data)`.
- Omit both and you get a generic `Custom matcher failed for:\n<actual>` message.

Template variables and functions:

| token | meaning |
|---|---|
| `{{.Actual}}` | the raw actual value |
| `{{.FormattedActual}}` | actual pretty-printed (indent 1) |
| `{{.To}}` | `"to"` on positive failure, `"not to"` on negated failure |
| `{{.Failure}}` / `{{.NegatedFailure}}` | bools for the two directions |
| `{{.Data}}` | whatever you passed to `WithTemplateData` |
| `{{format <obj> <indent>}}` | render any object with Gomega's formatter |

**Avoid recompiling templates in hot paths.** `WithTemplate` parses on every constructor call. Precompile once with `gcustom.ParseTemplate(str)` (use *that*, not raw `text/template`, so the `format` func is registered) and pass it to `MakeMatcher(matchFunc, tmpl)` or `.WithPrecompiledTemplate(tmpl)`.

## The raw GomegaMatcher interface

`MakeMatcher` returns this; implement it directly only when you need behavior `gcustom` doesn't expose:

```go
type GomegaMatcher interface {
	Match(actual any) (success bool, err error)
	FailureMessage(actual any) (message string)
	NegatedFailureMessage(actual any) (message string)
}
```

A hand-rolled version of the matcher above:

```go
func RepresentJSONifiedObject(expected any) types.GomegaMatcher {
	return &representJSONMatcher{expected: expected}
}

type representJSONMatcher struct{ expected any }

func (m *representJSONMatcher) Match(actual any) (bool, error) {
	response, ok := actual.(*http.Response)
	if !ok {
		return false, fmt.Errorf("RepresentJSONifiedObject matcher expects an *http.Response")
	}
	ptr := reflect.New(reflect.TypeOf(m.expected)).Interface()
	if err := json.NewDecoder(response.Body).Decode(ptr); err != nil {
		return false, fmt.Errorf("failed to decode JSON: %w", err)
	}
	decoded := reflect.ValueOf(ptr).Elem().Interface()
	return reflect.DeepEqual(decoded, m.expected), nil
}

func (m *representJSONMatcher) FailureMessage(actual any) string {
	return format.Message(actual, "to contain the JSON representation of", m.expected)
}

func (m *representJSONMatcher) NegatedFailureMessage(actual any) string {
	return format.Message(actual, "not to contain the JSON representation of", m.expected)
}
```

Conventions:

- A constructor function (here `RepresentJSONifiedObject`) returns the matcher; the struct itself stays unexported. Take concrete types in the constructor where you can; use `any` only when you genuinely must.
- **`FailureMessage`/`NegatedFailureMessage` are always called *after* `Match`** — stash anything you computed in `Match` on the struct and reuse it in the messages.
- Lean on `reflect` to interpret generic `actual`/`expected` inputs.

### format package helpers

Use these for messages instead of hand-rolling `fmt.Sprintf` so output matches every other Gomega matcher:

- `format.Message(actual, "to <verb>", expected)` → `Expected\n<actual>\n<message>\n<expected>`. Omit `expected` for a one-sided message.
- `format.Object(obj, indent)` → the pretty-printed block Gomega uses everywhere (`indent` is a `uint`).
- `format.MessageWithDiff(actual, "to equal", expected)` for string-diff highlighting.

## MatchMayChangeInTheFuture — stopping Eventually/Consistently early

`Eventually`/`Consistently` (see `gomega:async`) re-invoke your matcher each poll. If a match result can become permanently fixed (e.g. a closed channel, an exited process), implement the optional `OracleMatcher` method so polling can short-circuit:

```go
func (m *exitMatcher) MatchMayChangeInTheFuture(actual any) bool {
	session := actual.(*gexec.Session)
	return session.ExitCode() == -1 // true while still running; false once exited
}
```

This is **not** part of `GomegaMatcher` — most matchers don't need it. When present, Gomega calls it after each `Match`; returning `false` stops polling and fails/passes as appropriate (so `Eventually(session).Should(Exit(0))` fails the instant a wrong exit code lands instead of waiting out the timeout). `gcustom` matchers don't expose this; hand-write the type if you need it.

**Only consulted for bare values.** If `Eventually` is polling a *function*, Gomega can't assume the result is stable, so `MatchMayChangeInTheFuture` is skipped.

## Aborting from inside a matcher with StopTrying

A more direct, modern alternative: have `Match` return a `StopTrying(...)` error (or call `StopTrying(...).Now()`) to abort `Eventually`/`Consistently` immediately and fail with your message.

```go
func (m *thing) Match(actual any) (bool, error) {
	if irrecoverable {
		return false, StopTrying("the resource was deleted; it will never appear")
	}
	return check(actual), nil
}
```

You can `.Wrap(err)` and `.Attach("desc", obj)` onto the signal. Full polling-signal semantics (including `TryAgainAfter`) live in `gomega:async`.

## Testing your custom matchers

Drive the matcher under test with ordinary assertions, and for the *error* paths call `Match` directly (an erroring matcher would otherwise just fail your assertion):

```go
var _ = Describe("RepresentJSONifiedObject", func() {
	var book Book
	var response *http.Response

	BeforeEach(func() {
		book = Book{Title: "Les Miserables", Author: "Victor Hugo"}
		j, err := json.Marshal(book)
		Expect(err).NotTo(HaveOccurred())
		response = &http.Response{Body: io.NopCloser(bytes.NewBuffer(j))}
	})

	It("matches an http.Response carrying the object's JSON", func() {
		Expect(response).To(RepresentJSONifiedObject(book))
	})

	It("does not match a different payload", func() {
		response.Body = io.NopCloser(strings.NewReader(`{}`))
		Expect(response).NotTo(RepresentJSONifiedObject(book))
	})

	It("errors on the wrong actual type", func() {
		_, err := RepresentJSONifiedObject(book).Match("not a response")
		Expect(err).To(HaveOccurred())
	})
})
```

Assert on `FailureMessage(actual)` / `NegatedFailureMessage(actual)` when the wording matters.

## Packaging & contributing

- Keep matchers in their own package; export only the constructor.
- To contribute a matcher upstream, mimic Gomega's house style: format output via the `format` package, put the matcher and its tests in the `matchers` package, and add the constructor to `matchers.go` in the top-level package. Issues and PRs welcome at <https://github.com/onsi/gomega>.

## Wrapping the whole Gomega (advanced)

You can replace `gomega.Default` with a type implementing both `gomega.Gomega` and an `Inner() Gomega` method that returns the real default, delegating every call through. Useful for cross-suite logging, call counting, or injecting delays to surface timing dependencies — rarely needed, but the seam exists.

---

See also: `gomega:matchers` (catalog — check before building), `gomega:composing-matchers` (build matchers without code), `gomega:async` (polling signals).
