---
name: async
description: Polling assertions in Gomega — Eventually (poll until it passes) and Consistently (must keep passing), the func(g Gomega) callback idiom, WithTimeout/WithPolling/Within/ProbeEvery, WithContext and Ginkgo SpecContext, StopTrying/TryAgainAfter bail-outs, MustPassRepeatedly, and default-interval tuning. Use when an assertion can't be true synchronously — anything involving goroutines, channels, network calls, eventual consistency, or "wait until / stays true".
---

# Asynchronous assertions: Eventually & Consistently

`Expect` asserts *now*. When the thing under test settles over time, poll it. See `gomega:overview` for the mental model and `gomega:assertions` for the synchronous side. Dot-import assumed; with plain `testing` use `g := NewWithT(t)` and call `g.Eventually(...)`. Docs: <https://onsi.github.io/gomega/#making-asynchronous-assertions>.

- **`Eventually`** — polls until the matcher passes or the timeout elapses. Default: poll every 10ms for up to 1s.
- **`Consistently`** — polls for the whole duration and fails if the matcher *ever* fails. Default: poll every 10ms for 100ms. Use it to assert something "does not eventually" happen.

```go
Eventually(client.FetchCount).Should(BeNumerically(">=", 17))
Consistently(channel).ShouldNot(Receive())   // nothing ever arrives
```

## The three things you can poll

**1. A bare value — rarely what you want.** Gomega passes the value to the matcher every poll, but nothing re-reads it.

```go
Eventually(c, "50ms").Should(BeClosed())   // works: matcher re-inspects the channel
```

This is only correct for values whose *matcher* re-inspects live state (channels, `*gexec.Session`, `*gbytes.Buffer`). **Gotcha:** `Eventually(s).Should(Equal("changed"))` or `Eventually(slice).Should(HaveLen(3))` will **not** see updates — the value is captured once, and racing a goroutine over it trips the race detector. Pass a **function** instead.

**2. A function returning at least one value.** Polled repeatedly; the first return value goes to the matcher.

```go
Eventually(func() int { return client.FetchCount() }).Should(BeNumerically(">=", 17))
```

The **multi-return error idiom** applies: extra return values must all be zero, so a `(value, error)` function fails the poll while `err != nil` and matches `value` once it's `nil`.

```go
func FetchFromDB() (string, error)
Eventually(FetchFromDB).Should(Equal("got it"))   // passes when err==nil AND string matches
```

Pass arguments with `.WithArguments(...)` (supports variadic):

```go
Eventually(FetchFullName).WithArguments(1138).Should(Equal("Wookie"))
```

**3. A `func(g Gomega)` callback — the recommended powerful form.** The function makes its *own* assertions against the passed-in `g`; the whole block is retried until every assertion passes. Return zero values and match with `Succeed`, or return values to also match the result.

```go
Eventually(func(g Gomega) {
    model, err := client.Find(1138)
    g.Expect(err).NotTo(HaveOccurred())
    g.Expect(model.Reticulate()).To(Succeed())
    g.Expect(model.IsReticulated()).To(BeTrue())
}).Should(Succeed())

Eventually(func(g Gomega) (Widget, error) {
    ids, err := client.FetchIDs()
    g.Expect(err).NotTo(HaveOccurred())
    g.Expect(ids).To(ContainElement(1138))
    return client.FetchWidget(1138)
}).Should(Equal(expectedWidget))
```

**Gotcha:** inside the callback you **must** use the passed-in `g`, not the global `Expect`. Global assertions aren't intercepted by `Eventually` and will fail the test on the first miss instead of being retried.

## Configuring timeout & polling

Optional positional args (after `ACTUAL`): timeout, then polling interval, then `context.Context`. Each duration may be a `time.Duration`, a parseable string (`"100ms"`), or a `float64` (interpreted as seconds).

```go
Eventually(ACTUAL, "2s", "50ms").Should(MATCHER)
Consistently(ACTUAL, 500*time.Millisecond).Should(MATCHER)
```

Equivalently — and more readable — chain. `Within`/`ProbeEvery` are aliases for `WithTimeout`/`WithPolling`:

```go
Eventually(ACTUAL).WithTimeout(2 * time.Second).WithPolling(50 * time.Millisecond).Should(MATCHER)
Eventually(ACTUAL).Within(2 * time.Second).ProbeEvery(50 * time.Millisecond).Should(MATCHER)
```

For `Consistently` the first duration is the **window it must hold for**, not a timeout.

## Contexts & cancellation

`.WithContext(ctx)` (or passing `ctx` positionally, even first: `Eventually(ctx, ACTUAL)`) lets a cancelled context stop the poll.

```go
Eventually(ACTUAL).WithTimeout(t).WithPolling(p).WithContext(ctx).Should(MATCHER)
```

**Timeout × context interaction (`Eventually`):**
- Context **and** explicit timeout → stops at whichever comes first.
- Context **and no** explicit timeout → **no timeout applied**; polls until the context is cancelled. This is intentional so one `ctx` can govern a batch of `Eventually`s. Opt out with `EnforceDefaultTimeoutsWhenUsingContexts()`.
- `Consistently` always uses its `DURATION`; a cancelled context makes it **bail out early as a failure**, never as the duration controller.

**Ginkgo `SpecContext`** plugs straight in, and Gomega auto-injects the context as the polled function's first arg (after `g Gomega`, if present):

```go
It("fetches the count", func(ctx SpecContext) {
    Eventually(client.FetchCount).WithContext(ctx).WithArguments("/users").Should(BeNumerically(">=", 17))
}, SpecTimeout(time.Second))
```

The polled function runs **synchronously** — `Eventually` cannot kill a slow function, so thread the context *into* it (e.g. `client.FetchCount(ctx, ...)`) to make it interruptible. With `SpecContext`, Ginkgo's Progress Reports also surface which `Eventually` was running and its latest failure on timeout.

## Bailing out early: StopTrying / TryAgainAfter

Both signals work whether **returned as an `error`** or thrown via **`.Now()`** (a panic the poller catches). They also work from inside `func(g Gomega)` callbacks and from matchers.

**`StopTrying(msg)`** — stop polling immediately. **Always a failure for `Eventually`.**

```go
Eventually(func() (string, error) {
    if playerIndex == numPlayers {
        return "", StopTrying("no more players left")   // as error
    }
    name := client.FetchPlayer(playerIndex); playerIndex++
    return name, nil
}).Should(Equal("Patrick Mahomes"))

Eventually(func() []string {
    names, err := client.FetchAllPlayers()
    if err == client.IRRECOVERABLE_ERROR {
        StopTrying("irrecoverable error").Now()          // panic form
    }
    return names
}).Should(ContainElement("Patrick Mahomes"))
```

For `Consistently` only, `StopTrying(msg).Successfully()` ends the window early **without** failing (e.g. the monitored goroutine has finished). Not valid with `Eventually`, which always treats `StopTrying` as failure.

Enrich the message with `.Wrap(err)` (renders `<message>: <err>`) and `.Attach(description, obj)` (formats `obj` through Gomega's formatter; repeatable).

**`TryAgainAfter(duration)`** — same return/`.Now()` mechanics, but instead of stopping it waits `duration` before the next poll. Use it to back off when a service reports "unavailable, retry later." If the overall timeout elapses during that wait, both `Eventually` and `Consistently` fail and print the message (also supports `.Wrap`/`.Attach`).

## MustPassRepeatedly

Require N consecutive passes before `Eventually` succeeds — useful when a single transient pass isn't enough:

```go
Eventually(ACTUAL).MustPassRepeatedly(3).Should(MATCHER)
```

## Matchers can stop the poll

A matcher's `Match` returning `StopTrying` (or calling `.Now()`) halts polling and **fails** — same always-a-failure rule. Separately, a matcher implementing `MatchMayChangeInTheFuture(actual) bool` returning `false` signals out-of-band that no further change is possible, so `Eventually`/`Consistently` can stop without that being inherently a failure (e.g. `Receive` on a closed channel). See `gomega:custom-matchers` for both mechanisms.

## Modifying default intervals

Suite-wide, set in a `BeforeSuite`/`TestMain`:

```go
SetDefaultEventuallyTimeout(t time.Duration)
SetDefaultEventuallyPollingInterval(t time.Duration)
SetDefaultConsistentlyDuration(t time.Duration)
SetDefaultConsistentlyPollingInterval(t time.Duration)
```

Or via env vars (lower precedence than the `SetDefault…` calls), all parseable duration strings: `GOMEGA_DEFAULT_EVENTUALLY_TIMEOUT`, `GOMEGA_DEFAULT_EVENTUALLY_POLLING_INTERVAL`, `GOMEGA_DEFAULT_CONSISTENTLY_DURATION`, `GOMEGA_DEFAULT_CONSISTENTLY_POLLING_INTERVAL`, and `GOMEGA_ENFORCE_DEFAULT_TIMEOUTS_WHEN_USING_CONTEXTS`.
