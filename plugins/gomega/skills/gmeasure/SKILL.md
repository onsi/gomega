---
name: gmeasure
description: Benchmark and measure Go code with gmeasure — an Experiment groups named Measurements, recorded via RecordValue/RecordDuration/MeasureDuration or repeated Sample/SampleValue/SampleDuration with SamplingConfig, timed inline with a Stopwatch, summarized through GetStats/Stats (StatMin/Max/Mean/Median/StdDev, ValueFor/DurationFor) and compared with RankStats; decorate output with Units/Precision/Style/Annotation, render in Ginkgo via AddReportEntry, and persist with ExperimentCache. Use when you need human-readable benchmarks, performance reports, or regression baselines (not pass/fail assertions on their own).
---

# gmeasure: benchmarking and measuring code

`gmeasure` records benchmarks as **Experiments** that hold one or more named **Measurements**. Use it standalone (`fmt.Println(experiment)`) or wire it into Ginkgo for rich report output. Docs: <https://onsi.github.io/gomega/#gmeasure-benchmarking-code>. For the broader library see `gomega:overview`.

```go
import "github.com/onsi/gomega/gmeasure"
```

## Mental model

- An **Experiment** (`gmeasure.NewExperiment(name)`) groups related measurements for one system/context.
- A **Measurement** is a named bag of data points plus a `Type`: `MeasurementTypeValue` (`float64`) or `MeasurementTypeDuration` (`time.Duration`). It is created on its first recorded data point; later records of the same name append.
- **Stats** are statistical aggregates (min/max/mean/median/stddev) computed over a measurement's data points.

**gmeasure does not fail tests by itself.** It is a benchmarking/reporting tool. To gate on results you must pull `Stats` and write your own `Expect(...)` (or use `RankStats(...).Winner()`).

```go
experiment := gmeasure.NewExperiment("My Experiment")
experiment.RecordDuration("runtime", 3*time.Second) // creates the "runtime" measurement
experiment.RecordDuration("runtime", 5*time.Second) // appends a data point
```

## Recording values and durations

```go
// Direct values / durations:
experiment.RecordValue("length", 3.141)
experiment.RecordDuration("runtime", 200*time.Millisecond)

// Callback-driven — gmeasure times MeasureDuration for you:
v := experiment.MeasureValue("length", func() float64 { return computeLength() })
d := experiment.MeasureDuration("save", func() { client.Save(model) })
```

Experiments are **thread-safe** — `RecordX`/`MeasureX` may be called from any goroutine.

## Sampling: ensembles of data points

Run a callback repeatedly to build up many data points. Configure with `SamplingConfig`:

```go
type SamplingConfig struct {
	N                   int           // cap on number of samples
	Duration            time.Duration // cap on total sampling time
	NumParallel         int           // run samples across this many goroutines (>1)
	MinSamplingInterval time.Duration // minimum gap between samples (incompatible with NumParallel)
}
```

**At least one of `N` or `Duration` must be set** — otherwise sampling has no stop condition. With both, sampling stops at whichever limit hits first.

```go
// SampleDuration: time each call, append to "runtime"
experiment.SampleDuration("runtime", func(idx int) {
	RunAlgorithm()
}, gmeasure.SamplingConfig{N: 1000})

// SampleValue: record each returned float64
experiment.SampleValue("alloc-mb", func(idx int) float64 {
	return currentAllocMB()
}, gmeasure.SamplingConfig{Duration: time.Minute, NumParallel: 4})
```

`SampleAnnotatedDuration` / `SampleAnnotatedValue` take callbacks that also return a `gmeasure.Annotation` per data point. The bare `experiment.Sample(func(idx int){...}, cfg)` just drives the loop and isn't tied to a measurement — record whatever you like inside it.

## Stopwatch: timing sections inline

`experiment.NewStopwatch()` starts immediately. `Record(name)` stores elapsed time since the last `Reset` (or since creation) into a duration measurement; it returns the stopwatch so you can chain. **Stopwatches are not thread-safe — make a fresh one per goroutine inside `Sample`.**

```go
It("measures the end-to-end performance of the web-server", func() {
	experiment := gmeasure.NewExperiment("end-to-end performance")
	AddReportEntry(experiment.Name, experiment)

	experiment.Sample(func(idx int) {
		defer GinkgoRecover() // these run as goroutines and contain assertions
		stopwatch := experiment.NewStopwatch()

		model, err := client.Fetch("model-id-17")
		stopwatch.Record("fetch")
		Expect(err).NotTo(HaveOccurred())

		stopwatch.Reset()
		Expect(client.Save(model)).To(Succeed())
		stopwatch.Record("save").Reset()

		_, err = client.List("reticulated-models")
		stopwatch.Record("list")
		Expect(err).NotTo(HaveOccurred())
	}, gmeasure.SamplingConfig{N: 100, Duration: time.Minute, NumParallel: 8})
})
```

`Pause()` / `Resume()` bracket out work you don't want counted.

## Stats and rankings

`experiment.GetStats(name)` returns a `Stats`. Pull individual stats with the `gmeasure.Stat` enum: `StatMin`, `StatMax`, `StatMean`, `StatMedian`, `StatStdDev`.

```go
stats := experiment.GetStats("runtime")
med := stats.DurationFor(gmeasure.StatMedian) // time.Duration (Duration measurements)
mb  := experiment.GetStats("alloc-mb").ValueFor(gmeasure.StatMax) // float64 (Value measurements)

// FloatFor(stat) works for either type (durations come back as float64(ns));
// StringFor(stat) returns a formatted, precision-aware string.
```

Compare measurements with `gmeasure.RankStats(criterion, ...Stats)`, then `.Winner()`:

```go
ranking := gmeasure.RankStats(gmeasure.LowerMedianIsBetter,
	experiment.GetStats("runtime: algorithm 1"),
	experiment.GetStats("runtime: algorithm 2"),
)
AddReportEntry("Ranking", ranking)
Expect(ranking.Winner().MeasurementName).To(Equal("runtime: algorithm 2"))
```

Criteria: `LowerMeanIsBetter`, `HigherMeanIsBetter`, `LowerMedianIsBetter`, `HigherMedianIsBetter`, `LowerMinIsBetter`, `HigherMinIsBetter`, `LowerMaxIsBetter`, `HigherMaxIsBetter`.

## Decorations: units, precision, style, annotations

Pass these as variadic args to any `RecordX`/`MeasureX`/`SampleX` call. **`Units`, `Precision`, and `Style` must be set on the first data point of a measurement** (that's when the measurement is initialized); later they're ignored. `Annotation` can be attached to any individual data point.

```go
experiment.RecordValue("length", 3.141,
	gmeasure.Units("inches"),     // rendered as "length [inches]"
	gmeasure.Precision(2),        // int → %.2f for values
	gmeasure.Style("{{blue}}"),   // Ginkgo console style for the row
	gmeasure.Annotation("box A"),
)
experiment.RecordValue("length", 2.71, gmeasure.Annotation("box B")) // appends w/ annotation

// For Duration measurements, Precision takes a time.Duration to round to:
experiment.MeasureDuration("teardown", teardown, gmeasure.Precision(time.Millisecond))
```

`experiment.RecordNote("...")` adds a contextual row to the rendered table.

## Ginkgo integration

Register the experiment (and any rankings) as report entries so Ginkgo renders styled tables and includes them in machine-readable reports (`ginkgo --json-report`):

```go
experiment := gmeasure.NewExperiment("my benchmark")
AddReportEntry(experiment.Name, experiment) // also works for Measurement and Ranking
```

**Without `AddReportEntry` you'll see no output under Ginkgo.** Outside Ginkgo, just `fmt.Println(experiment)` — `Experiment`/`Measurement`/`Ranking` all implement `String()` (and `ColorableString()` for styled output).

## Caching experiments

`gmeasure.NewExperimentCache(dir)` (returns `(ExperimentCache, error)`) persists experiments to disk keyed by name + version. Use it to skip expensive re-runs or to store committed baselines for regression checks.

```go
cache, err := gmeasure.NewExperimentCache("./gmeasure-cache")
Expect(err).NotTo(HaveOccurred())

const VERSION = 1 // bump to bust the cache and force recomputation
if experiment := cache.Load(name, VERSION); experiment != nil {
	AddReportEntry(experiment.Name, experiment)
	Skip("cached") // reuse cached results, skip re-measuring
} else {
	experiment = gmeasure.NewExperiment(name)
	// ... measure ...
	cache.Save(experiment.Name, VERSION, experiment)
}
```

`cache.Load` returns `nil` on a miss. Other methods: `Save`, `Delete`, `List`, `Clear`. For a regression gate, load a committed baseline and assert the current stats stay within a tolerance:

```go
baseline := cache.Load("perf", 1)
if baseline == nil {
	cache.Save("perf", 1, experiment) // first run establishes the baseline
} else {
	cur := experiment.GetStats("fetch")
	base := baseline.GetStats("fetch")
	Expect(cur.DurationFor(gmeasure.StatMean)).To(BeNumerically("~",
		base.DurationFor(gmeasure.StatMean), base.DurationFor(gmeasure.StatStdDev)))
}
```
