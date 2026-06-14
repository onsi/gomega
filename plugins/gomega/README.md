# Gomega plugin for Claude Code

Skills that help an AI agent (and you) write expressive, correct assertions with [Gomega](https://onsi.github.io/gomega/) — the matcher/assertion library for Go, best paired with [Ginkgo](https://github.com/onsi/ginkgo).

## Install

The Gomega repo doubles as the marketplace:

```
/plugin marketplace add onsi/gomega
/plugin install gomega@gomega
```

## What you get

All skills are namespaced under `gomega:` and activate when you're writing or reviewing Go tests that use Gomega.

| Skill | Use it when |
|---|---|
| `gomega:overview` | You want the mental model — `Expect`/`Ω`, matchers-are-values, sync vs async, and a map of the whole library (read me first). |
| `gomega:assertions` | You're writing synchronous assertions: the multi-return error idiom, `Succeed`/`HaveOccurred`, annotations, output tuning, and asserting inside helpers. |
| `gomega:async` | You need `Eventually`/`Consistently` — polling functions, the `g Gomega` callback, contexts, timeouts, and `StopTrying`/`TryAgainAfter`. |
| `gomega:matchers` | You want to find or choose the right matcher — the full catalog, grouped, with the gotchas flagged (start here to stop over-using `Equal`). |
| `gomega:composing-matchers` | You're combining matchers: `And`/`Or`/`Not`, `WithTransform`, `HaveField`, `HaveValue`, and nesting matchers as arguments. |
| `gomega:custom-matchers` | A built-in or composed matcher can't say what you mean and you need to write your own (`GomegaMatcher`, `gcustom`). |
| `gomega:gstruct` | You're asserting against large or deeply nested structs, slices, and maps. |
| `gomega:ghttp` | You're testing code that makes outbound HTTP requests. |
| `gomega:gexec` | You're building, running, signaling, or asserting on external processes. |
| `gomega:gbytes` | You're asserting on streaming or incremental output rather than a complete value. |
| `gomega:gleak` | You want to verify a test left no goroutines leaked. |
| `gomega:gmeasure` | You need human-readable benchmarks, performance reports, or regression baselines. |

## Versioning

These skills track the Gomega library. The narrative docs at <https://onsi.github.io/gomega/> are the source of truth; pin to the Gomega version you've `go get`'d.
