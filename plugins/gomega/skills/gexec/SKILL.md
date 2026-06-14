---
name: gexec
description: Testing external processes with gexec — compile binaries with Build/BuildWithEnvironment/BuildIn and CleanupBuildArtifacts, start them with Start returning a *Session, await exit with the Exit matcher (Eventually(session).Should(Exit(0))), Wait/ExitCode, signal via Kill/Terminate/Interrupt/Signal and package-level KillAndWait/TerminateAndWait, and assert on session.Out/Err which are gbytes buffers (Say, Contents). Use when building, running, signaling, or asserting on subprocesses in Go tests.
---

# gexec: testing external processes

`gexec` compiles Go binaries, starts them as subprocesses, sends them signals, and exposes their stdout/stderr as `gbytes.Buffer`s so you can assert on streaming output and exit codes. Import it normally (Gomega is dot-imported):

```go
import (
	. "github.com/onsi/gomega"
	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega/gexec"
)
```

Docs: <https://onsi.github.io/gomega/#gexec-testing-external-processes>.

## Compiling binaries

`gexec.Build(packagePath, ...args)` runs `go build` and returns the path to a temp binary.

```go
var pathToCLI string

BeforeSuite(func() {
	var err error
	pathToCLI, err = gexec.Build("github.com/spacely/sprockets")
	Expect(err).NotTo(HaveOccurred())
})

AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})
```

- `gexec.BuildWithEnvironment(packagePath, env []string, ...args)` — set env vars for the build (e.g. `GOOS`/`GOARCH` for cross-compilation).
- `gexec.BuildIn(gopath, packagePath, ...args)` — build with a custom `GOPATH` (e.g. against vendored deps).

**Always `defer`/`AfterSuite` `gexec.CleanupBuildArtifacts()`.** It deletes the temp binaries; skipping it leaks files into your temp dir across runs.

## Starting a process

`gexec.Start(cmd, outWriter, errWriter)` calls `cmd.Start()` and returns a `*gexec.Session` that wraps and monitors the process, forwarding stdout/stderr to the writers you pass.

```go
command := exec.Command(pathToCLI, "-api=127.0.0.1:8899")
session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
Expect(err).NotTo(HaveOccurred())
```

**Pass `GinkgoWriter` for both writers.** Output is then silent on passing tests but printed on failure (and always under `ginkgo -v`) — invaluable for debugging.

## Asserting on exit

The `gexec.Exit()` matcher is pollable and only works on a `*Session` — pair it with `Eventually` (→ `gomega:async`):

```go
Eventually(session).Should(gexec.Exit())   // exited, any code
Eventually(session).Should(gexec.Exit(0))  // exited with code 0
```

`session.ExitCode()` returns the raw code, or `-1` if the process hasn't exited yet.

`session.Wait([timeout])` blocks until exit (failing if it doesn't exit within the default `Eventually` timeout) and returns the session for chaining:

```go
session.Wait(5 * time.Second)
```

**Don't `Wait` on the wrapped `exec.Cmd` yourself** — `gexec` already calls `Wait` to monitor the process. `session.Wait` is just `Eventually` under the hood.

## Sending signals

```go
session.Kill()         // SIGKILL
session.Interrupt()    // SIGINT
session.Terminate()    // SIGTERM
session.Signal(sig)    // arbitrary os.Signal
```

Each returns the session, so chain with `Wait`: `session.Terminate().Wait()`. Signaling an already-exited process is a no-op.

### Signaling every started session

Package-level helpers signal **all** sessions `gexec` has started, in any context — ideal for cleanup:

```go
gexec.Kill()              // SIGKILL to all
gexec.Terminate()         // SIGTERM to all
gexec.Interrupt()         // SIGINT to all
gexec.Signal(sig)         // os.Signal to all
gexec.KillAndWait()       // signal all, then wait
gexec.TerminateAndWait(2 * time.Second)  // per-process timeout

AfterSuite(func() { gexec.KillAndWait() })
```

**These are global.** Calling them in an `AfterEach` will also signal processes started in `BeforeSuite`. Good practice is to ensure all processes are killed before the suite ends.

## Asserting on output

`session.Out` and `session.Err` are `gbytes.Buffer`s connected to stdout/stderr, so use the `gbytes.Say` matcher for ordered streaming assertions (→ `gomega:gbytes`). The session itself is a `BufferProvider` for `Out`:

```go
Eventually(session.Out).Should(gbytes.Say("hello [A-Za-z], nice to meet you"))
Eventually(session.Err).Should(gbytes.Say("oops!"))
Eventually(session).Should(gbytes.Say("hello"))  // shorthand for session.Out
```

To grab the whole output after exit, `Wait()` returns the session, so `.Out.Contents()` gives a `[]byte` (→ `gomega:matchers`):

```go
Expect(session.Wait().Out.Contents()).To(ContainSubstring("finished successfully"))
```

## End-to-end

```go
var pathToCLI string

var _ = BeforeSuite(func() {
	var err error
	pathToCLI, err = gexec.Build("github.com/spacely/sprockets")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

var _ = It("greets and exits cleanly", func() {
	command := exec.Command(pathToCLI, "--name=George")
	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())

	Eventually(session).Should(gbytes.Say("hello George"))
	Eventually(session).Should(gexec.Exit(0))
	Expect(session.Out.Contents()).To(ContainSubstring("done"))
})
```
