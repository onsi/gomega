---
name: gbytes
description: Testing streaming io buffers with gbytes — gbytes.NewBuffer() (an io.Writer also returned by gexec sessions), the Say(regexp) matcher that forward-scans from a moving read cursor, the canonical Eventually(buffer).Should(Say(...)) streaming pattern, sequential cursor-advancing Say calls, Contents(), BufferWithBytes/BufferReader, buffer.Detect for branching, and TimeoutReader/Writer/Closer for testing blocking io.Reader/Writer/Closer. Use when asserting on streaming or incremental output (process stdout/stderr, API streams, io.Readers) rather than a complete value.
---

# gbytes: testing streaming buffers

`gbytes` provides `*gbytes.Buffer`, an in-memory `io.Writer` (and `io.Reader`/`io.Closer`) that captures everything written to it, plus the `Say` matcher for making **ordered** assertions against streaming data as it arrives. Pairs naturally with `gomega:async` (poll with `Eventually`) and `gomega:gexec` (a session's `.Out`/`.Err` are `*gbytes.Buffer`s, and the session itself is a `BufferProvider`). Docs: <https://onsi.github.io/gomega/#gbytes-testing-streaming-buffers>.

```go
import (
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)
```

## The buffer

`gbytes.NewBuffer()` returns an empty `*gbytes.Buffer`. Hand it to anything that wants an `io.Writer` and it accumulates the bytes:

```go
buffer := gbytes.NewBuffer()
go client.AttachToDataStream(buffer) // client writes to buffer concurrently
```

- `gbytes.BufferWithBytes(b []byte)` — seed a buffer with already-captured bytes.
- `buffer.Contents() []byte` — **all** bytes ever written, regardless of cursor position. Use this for whole-buffer assertions (`Expect(buffer.Contents()).To(ContainSubstring(...))`).
- `buffer.Close()` / `buffer.Closed()` — mark a buffer done. Writes to a closed buffer error. A closed buffer tells `Eventually` to give up (see below).

## The `Say` matcher

`gbytes.Say(pattern, args...)` takes a **regular expression** (optionally `fmt.Sprintf`-formatted with `args`) and matches against the buffer's *unread* portion:

```go
Expect(buffer).To(gbytes.Say(`hello \w+`))
```

**Say is cursor-based, not whole-buffer.** Each buffer carries an opaque read cursor you cannot access. When `Say` matches, it fast-forwards the cursor to just past the match. The next `Say` only sees bytes *after* that point. This is what makes ordered assertions work — and the gotcha below.

## The canonical streaming pattern: `Eventually` + `Say`

Streaming output arrives over time, so poll the buffer with `Eventually` (`gomega:async`). Each successful `Say` advances the cursor, so successive `Say`s assert on successive output in order:

```go
Eventually(buffer).Should(gbytes.Say(`Attached as client \d+`))

client.ReticulateSplines()
Eventually(buffer).Should(gbytes.Say(`reticulating splines`))

client.EncabulateRetros(7)
Eventually(buffer).Should(gbytes.Say(`encabulating 7 retros`))
```

Because the cursor only moves forward, this counts repeats correctly — these two assertions pass only if `reticulating splines` appears **twice**:

```go
client.ReticulateSplines()
Eventually(buffer).Should(gbytes.Say(`reticulating splines`))
client.ReticulateSplines()
Eventually(buffer).Should(gbytes.Say(`reticulating splines`))
```

And consequently this (counterintuitively) passes — the first `Say` consumed the match:

```go
Eventually(buffer).Should(gbytes.Say(`reticulating splines`))
Consistently(buffer).ShouldNot(gbytes.Say(`reticulating splines`))
```

`Say` works on a `*gbytes.Buffer` or any `BufferProvider` (anything with a `Buffer() *gbytes.Buffer` method, e.g. a `gexec.Session`), so `Eventually(session).Should(gbytes.Say(...))` works directly.

## Handling branches

When the test must react to *whichever* output arrives, use `buffer.Detect(regexp, args...)`, which returns a channel that fires once on match (and fast-forwards the cursor). Always `CancelDetects()` to clean up the spawned goroutines:

```go
client.Authorize()
select {
case <-buffer.Detect("You are not logged in"):
	client.Login()
case <-buffer.Detect("Success"):
	// carry on
case <-time.After(time.Second):
	Fail("timed out waiting for output")
}
buffer.CancelDetects()
```

## Testing `io.Reader`/`io.Writer`/`io.Closer`

These interfaces are expected to block, so calling `Read`/`Write`/`Close` directly in a test risks hanging forever. Wrap them with timeouts; the wrappers return `gbytes.ErrTimeout` if the operation doesn't complete in time:

```go
p := make([]byte, 5)
_, err := gbytes.TimeoutReader(reader, time.Second).Read(p)
Expect(err).NotTo(HaveOccurred())
```

`gbytes.TimeoutReader`, `gbytes.TimeoutWriter`, and `gbytes.TimeoutCloser` each wrap the matching interface. To use `Say` against an `io.Reader`, wrap it with `gbytes.BufferReader(reader)` — it launches an `io.Copy` goroutine into a fresh buffer (closed when the copy completes). Because the copy is async you **must** use `Eventually`:

```go
Eventually(gbytes.BufferReader(reader)).Should(gbytes.Say("abcde"))
```

## Gotchas

- **Say takes a regexp — escape metacharacters.** `(`, `)`, `.`, `[`, `+`, `?`, `\d`, etc. are regex syntax. To match a literal `splines (v2)` write `gbytes.Say(`splines \(v2\)`)`. The pattern is compiled with `regexp.MustCompile`, so a bad pattern panics.
- **The cursor only moves forward.** You cannot re-match output an earlier `Say` already consumed. If you need to assert on the whole buffer, use `buffer.Contents()` with a string matcher (`gomega:matchers`) instead of `Say`.
- **Pair `Say` with `Eventually` for live streams; use plain `Expect` only on an already-complete buffer.** A bare `Expect(buffer).To(Say(...))` checks the buffer once, right now — it will flake if the data hasn't arrived yet.
- **A closed buffer aborts `Eventually` early.** Once `Close()` is called, a pending `Say` can never succeed on new data, so the matcher signals `Eventually`/`Consistently` to stop polling and fail immediately rather than wait out the timeout.
