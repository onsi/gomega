## 1.0 (8/2/2014)

No changes. Dropping "beta" from the version number.

## 1.0.0-beta (7/8/2014)
Breaking Changes:

- Changed OmegaMatcher interface.  Instead of having `Match` return failure messages, two new methods `FailureMessage` and `NegatedFailureMessage` are called instead.
- Moved and renamed OmegaFailHandler to types.GomegaFailHandler and OmegaMatcher to types.GomegaMatcher.  Any references to OmegaMatcher in any custom matchers will need to be changed to point to types.GomegaMatcher

New Test-Support Features:

- `ghttp`: supports testing http clients
    - Provides a flexible fake http server
    - Provides a collection of chainable http handlers that perform assertions.
- `gbytes`: supports making ordered assertions against streams of data
    - Provides a `gbytes.Buffer`
    - Provides a `Say` matcher to perform ordered assertions against output data
- `gexec`: supports testing external processes
    - Provides support for building Go binaries
    - Wraps and starts `exec.Cmd` commands
    - Makes it easy to assert against stdout and stderr
    - Makes it easy to send signals and wait for processes to exit
    - Provides an `Exit` matcher to assert against exit code.

DSL Changes:

- `Eventually` and `Consistently` can accept `time.Duration` interval and polling inputs.
- The default timeouts for `Eventually` and `Consistently` are now configurable.

New Matchers:

- `ConsistOf`: order-independent assertion against the elements of an array/slice or keys of a map.
- `BeTemporally`: like `BeNumerically` but for `time.Time`
- `HaveKeyWithValue`: asserts a map has a given key with the given value.

Updated Matchers:

- `Receive` matcher can take a matcher as an argument and passes only if the channel under test receives an objet that satisfies the passed-in matcher.
- Matchers that implement `MatchMayChangeInTheFuture(actual interface{}) bool` can inform `Eventually` and/or `Consistently` when a match has no chance of changing status in the future.  For example, `Receive` returns `false` when a channel is closed.

Misc:

- Start using semantic versioning
- Start maintaining changelog

Major refactor:

- Pull out Gomega's internal to `internal`
