---
layout: default
title: Gomega
---
{% raw  %}
![Gomega](./images/gomega.png)

[Gomega](http://github.com/onsi/gomega) is a matcher/assertion library.  It is best paired with the [Ginkgo](http://github.com/onsi/ginkgo) BDD test framework, but can be adapted for use in other contexts too.

## Support Policy

Gomega provides support for versions of Go that are noted by the [Go release policy](https://golang.org/doc/devel/release.html#policy) i.e. N and N-1 major versions.

## Getting Gomega

Just `go get` it:

```bash
$ go get github.com/onsi/gomega/...
```

## Getting Gomega as needed

Instead of getting all of Gomega and it's dependency tree, you can use the go command to get the dependencies as needed.

For example, import gomega in your test code:

```go
import "github.com/onsi/gomega"
```

Use `go get -t` to retrieve the packages referenced in your test code:

```bash
$ cd /path/to/my/app
$ go get -t ./...
```

## Using Gomega with Ginkgo

When a Gomega assertion fails, Gomega calls a `GomegaFailHandler`.  This is a function that you must provide using `gomega.RegisterFailHandler()`.

If you're using Ginkgo, all you need to do is:

```go
gomega.RegisterFailHandler(ginkgo.Fail)
```

before you start your test suite.

If you use the `ginkgo` CLI to `ginkgo bootstrap` a test suite, this hookup will be automatically generated for you.

> `GomegaFailHandler` is defined in the `types` subpackage.

## Using Gomega with Golang's XUnit-style Tests

Though Gomega is tailored to work best with Ginkgo it is easy to use Gomega with Golang's XUnit style tests.  Here's how:

To use Gomega with Golang's XUnit style tests:

```go
func TestFarmHasCow(t *testing.T) {
    g := NewGomegaWithT(t)

    f := farm.New([]string{"Cow", "Horse"})
    g.Expect(f.HasCow()).To(BeTrue(), "Farm should have cow")
}
```

`NewGomegaWithT(t)` wraps a `*testing.T` and returns a struct that supports `Expect`, `Eventually`, and `Consistently`.

## Making Assertions

Gomega provides two notations for making assertions.  These notations are functionally equivalent and their differences are purely aesthetic.

- When you use the `Ω` notation, your assertions look like this:

```go
Ω(ACTUAL).Should(Equal(EXPECTED))
Ω(ACTUAL).ShouldNot(Equal(EXPECTED))
```

- When you use the `Expect` notation, your assertions look like this:

```go
Expect(ACTUAL).To(Equal(EXPECTED))
Expect(ACTUAL).NotTo(Equal(EXPECTED))
Expect(ACTUAL).ToNot(Equal(EXPECTED))
```

On OS X the `Ω` character should be easy to type, it is usually just option-z: `⌥z`.

On the left hand side, you can pass anything you want in to `Ω` and `Expect` for `ACTUAL`.  On the right hand side you must pass an object that satisfies the `GomegaMatcher` interface.  Gomega's matchers (e.g. `Equal(EXPECTED)`) are simply functions that create and initialize an appropriate `GomegaMatcher` object.

> Note that `Should` and `To` are just syntactic sugar and are functionally identical. Same is the case for `ToNot` and `NotTo`.

> The `GomegaMatcher` interface is pretty simple and is discussed in the [custom matchers](#adding-your-own-matchers) section.  It is defined in the `types` subpackage.

### Handling Errors

It is a common pattern, in Golang, for functions and methods to return two things - a value and an error.  For example:

```go
func DoSomethingHard() (string, error) {
    ...
}
```

To assert on the return value of such a method you might write a test that looks like this:

```go
result, err := DoSomethingHard()
Ω(err).ShouldNot(HaveOccurred())
Ω(result).Should(Equal("foo"))
```

Gomega streamlines this very common use case.  Both `Ω` and `Expect` accept *multiple* arguments.  The first argument is passed to the matcher, and the match only succeeds if *all* subsequent arguments are `nil` or zero-valued.  With this, we can rewrite the above example as:

```go
Ω(DoSomethingHard()).Should(Equal("foo"))
```

This will only pass if the return value of `DoSomethingHard()` is `("foo", nil)`.

Additionally, if you call a function with a single `error` return value you can use the `Succeed` matcher to assert the function has returned without error.  So for a function of the form:

```go
func DoSomethingSimple() error {
    ...
}
```

You can either write:

```go
err := DoSomethingSimple()
Ω(err).ShouldNot(HaveOccurred())
```

Or you can write:

```go
Ω(DoSomethingSimple()).Should(Succeed())
```

> You should not use a function with multiple return values (like `DoSomethingHard`) with `Succeed`.  Matchers are only passed the *first* value provided to `Ω`/`Expect`, the subsequent arguments are handled by `Ω` and `Expect` as outlined above.  As a result of this behavior `Ω(DoSomethingHard()).ShouldNot(Succeed())` would never pass.

Assertions about errors on functions with multiple return values can be made as follows (and in a lazy way when not asserting that all other return values are zero values):

```go
_, _, _, err := MUltipleReturnValuesFunc()
Ω(err).Should(HaveOccurred())
```

Alternatively, such error assertions on multi return value functions can be simplified by chaining `Error` to `Ω` and `Expect`. Doing so will additionally automatically assert that all return values, except for the trailing error return value, are in fact zero values:

```go
Ω(MultipleReturnValuesFunc()).Error().Should(HaveOccurred())
```

Similar, asserting that no error occured is supported, too (where the other return values are allowed to take on any value):

```go
Ω(MultipleReturnValuesFunc()).Error().ShouldNot(HaveOccured())
```

### Annotating Assertions

You can annotate any assertion by passing either a format string (and optional inputs to format) or a function of type `func() string` after the `GomegaMatcher`:

```go
Ω(ACTUAL).Should(Equal(EXPECTED), "My annotation %d", foo)
Ω(ACTUAL).ShouldNot(Equal(EXPECTED), "My annotation %d", foo)
Expect(ACTUAL).To(Equal(EXPECTED), "My annotation %d", foo)
Expect(ACTUAL).NotTo(Equal(EXPECTED), "My annotation %d", foo)
Expect(ACTUAL).ToNot(Equal(EXPECTED), "My annotation %d", foo)
Expect(ACTUAL).To(Equal(EXPECTED), func() string { return "My annotation" })
```

If you pass a format string, the format string and inputs will be passed to `fmt.Sprintf(...)`.
If you instead pass a function, the function will be lazily evaluated if the assertion fails.
In both cases, if the assertion fails, Gomega will print your annotation alongside its standard failure message.

This is useful in cases where the standard failure message lacks context.  For example, if the following assertion fails:

```go
Ω(SprocketsAreLeaky()).Should(BeFalse())
```

Gomega will output:

```
Expected
 <bool>: true
to be false
```

But this assertion:

```go
Ω(SprocketsAreLeaky()).Should(BeFalse(), "Sprockets shouldn't leak")
```

Will offer the more helpful output:

```
Sprockets shouldn't leak
Expected
  <bool>: true
to be false
```


### Adjusting Output

When a failure occurs, Gomega prints out a recursive description of the objects involved in the failed assertion.  This output can be very verbose, but Gomega's philosophy is to give as much output as possible to aid in identifying the root cause of a test failure.

These recursive object renditions are performed by the `format` subpackage.  `format` provides some globally adjustable settings to tune Gomega's output:

- `format.MaxLength = 4000`: Gomega will recursively traverse nested data structures as it produces output. If the length of this string representation is more than MaxLength, it will be truncated to MaxLength. To disable this behavior, set the MaxLength to `0`.
- `format.MaxDepth = 10`: Gomega will recursively traverse nested data structures as it produces output. By default the maximum depth of this recursion is set to `10` you can adjust this to see deeper or shallower representations of objects.
- Implementing `format.GomegaStringer`: If `GomegaStringer` interface is implemented on an object, Gomega will call `GomegaString` for an object's string representation. This is regardless of the `format.UseStringerRepresentation` value. Best practice to implement this interface is to implement it in a helper test file (e.g. `helper_test.go`) to avoid leaking it to your package's exported API.
- `format.UseStringerRepresentation = false`: Gomega does *not* call `String` or `GoString` on objects that satisfy the `Stringer` and `GoStringer` interfaces.  Oftentimes such representations, while more human readable, do not contain all the relevant information associated with an object thereby making it harder to understand why a test might be failing.  If you'd rather see the output of `String` or `GoString` set this property to `true`.


> For a tricky example of why `format.UseStringerRepresentation = false` is your friend, check out issue [#37](https://github.com/onsi/gomega/issues/37).

- `format.PrintContextObjects = false`: Gomega by default will not print the content of objects satisfying the context.Context interface, due to too much output. If you want to enable displaying that content, set this property to `true`.

If you want to use Gomega's recursive object description in your own code you can call into the `format` package directly:

```go
fmt.Println(format.Object(theThingYouWantToPrint, 1))
```

- `format.TruncatedDiff = true`: Gomega will truncate long strings and only show where they differ. You can set this to `false` if
you want to see the full strings.

## Making Asynchronous Assertions

Gomega has support for making *asynchronous* assertions.  There are two functions that provide this support: `Eventually` and `Consistently`.

### Eventually

`Eventually` checks that an assertion *eventually* passes.  `Eventually` blocks when called and attempts an assertion periodically until it passes or a timeout occurs.  Both the timeout and polling interval are configurable as optional arguments:

```go
Eventually(ACTUAL, (TIMEOUT), (POLLING_INTERVAL)).Should(MATCHER)
```

The first optional argument is the timeout (which defaults to 1s), the second is the polling interval (which defaults to 10ms).  Both intervals can be specified as time.Duration, parsable duration strings (e.g. "100ms") or `float64` (in which case they are interpreted as seconds).

> As with synchronous assertions, you can annotate asynchronous assertions by passing either a format string and optional inputs or a function of type `func() string` after the `GomegaMatcher`.

Alternatively, the timeout and polling interval can also be specified by chaining `WithTimeout` and `WithPolling` to `Eventually`:

```go
Eventually(ACTUAL).WithTimeout(TIMEOUT).WithPolling(POLLING_INTERVAL).Should(MATCHER)
```

Eventually works with any Gomega compatible matcher and supports making assertions against three categories of `ACTUAL` value:

#### Category 1: Making `Eventually` assertions on values

There are several examples of values that can change over time.  These can be passed in to `Eventually` and will be passed to the matcher repeatedly until a match occurs.  For example:

```go
c := make(chan bool)
go DoStuff(c)
Eventually(c, "50ms").Should(BeClosed())
```

will poll the channel repeatedly until it is closed.  In this example `Eventually` will block until either the specified timeout of 50ms has elapsed or the channel is closed, whichever comes first.

Several Gomega libraries allow you to use Eventually in this way.  For example, the `gomega/gexec` package allows you to block until a `*gexec.Session` exits successfuly via:

```go
Eventually(session).Should(gexec.Exit(0))
```

And the `gomega/gbytes` package allows you to monitor a streaming `*gbytes.Buffer` until a given string is seen:

```go
Eventually(buffer).Should(gbytes.Say("hello there"))
```

In these examples, both `session` and `buffer` are designed to be thread-safe when polled by the `Exit` and `Say` matchers.  This is not true in general of most raw values, so while it is tempting to do something like:

```go
/* === INVALID === */
var s *string
go mutateStringEventually(s)
Eventually(s).Should(Equal("I've changed"))
```

this will trigger Go's race detector as the goroutine polling via Eventually will race over the value of `s` with the goroutine mutating the string.

Similarly, something like `Eventually(slice).Should(HaveLen(N))` probably won't do what you think it should -- `Eventually` will be passed a pointer to the slice, yes, but if the slice is being `append`ed to (as in: `slice = append(slice, ...)`) Go will generate a new pointer and the pointer passed to `Eventually` will not contain the new elements.

In both cases you should always pass `Eventually` a function that, when polled, returns the latest value of the object in question in a thread-safe way.

#### Category 2: Making `Eventually` assertions on functions

`Eventually` can be passed functions that **take no arguments** and **return at least one value**.  When configured this way, `Eventually` will poll the function repeatedly and pass the first returned value to the matcher.

For example:

```go
Eventually(func() int {
   return client.FetchCount()
}).Should(BeNumerically(">=", 17))
```

will repeatedly poll `client.FetchCount` until the `BeNumerically` matcher is satisfied.

> Note that this example could have been written as `Eventually(client.FetchCount).Should(BeNumerically(">=", 17))`

If multple values are returned by the function, `Eventually` will pass the first value to the matcher and require that all others are zero-valued.  This allows you to pass `Eventually` a function that returns a value and an error - a common pattern in Go.

For example, consider a method that returns a value and an error:

```go
func FetchFromDB() (string, error)
```

Then

```go
Eventually(FetchFromDB).Should(Equal("got it"))
```

will pass only if and when the returned error is `nil` *and* the returned string satisfies the matcher.

It is important to note that the function passed into Eventually is invoked **synchronously** when polled.  `Eventually` does not (in fact, it cannot) kill the function if it takes longer to return than `Eventually`'s configured timeout.  You should design your functions with this in mind.

#### Category 3: Making assertions _in_ the function passed into `Eventually`

When testing complex systems it can be valuable to assert that a *set* of assertions passes `Eventually`.  `Eventually` supports this by accepting functions that take **a single `Gomega` argument** and **return zero or more values**.

Here's an example that makes some asssertions and returns a value and error:

```go
Eventually(func(g Gomega) (Widget, error) {
    ids, err := client.FetchIDs()
    g.Expect(err).NotTo(HaveOccurred())
    g.Expect(ids).To(ContainElement(1138))
    return client.FetchWidget(1138)
}).Should(Equal(expectedWidget))
```

will pass only if all the assertions in the polled function pass and the return value satisfied the matcher.  Note that the assertions in the body of the polled function must be performed using the passed-in `g Gomega` object.  If you use the global DSL expectations, `Eventually` will not intercept any failures and the test will fail.

`Eventually` also supports a special case polling function that takes a single `Gomega` argument and returns no values.  `Eventually` assumes such a function is making assertions and is designed to work with the `Succeed` matcher to validate that all assertions have passed.

For example:

```go
Eventually(func(g Gomega) {
   model, err := client.Find(1138)
   g.Expect(err).NotTo(HaveOccurred())
   g.Expect(model.Reticulate()).To(Succeed())
   g.Expect(model.IsReticulated()).To(BeTrue())
   g.Expect(model.Save()).To(Succeed())
}).Should(Succeed())
```

will rerun the function until all assertions pass.


### Consistently

`Consistently` checks that an assertion passes for a period of time.  It does this by polling its argument repeatedly during the period. It fails if the matcher ever fails during that period.

For example:

```go
Consistently(func() []int {
    return thing.MemoryUsage()
}).Should(BeNumerically("<", 10))
```

`Consistently` will poll the passed in function repeatedly and check the return value against the `GomegaMatcher`.  `Consistently` blocks and only returns when the desired duration has elapsed or if the matcher fails.  The default value for the wait-duration is 100 milliseconds.  The default polling interval is 10 milliseconds.  Like `Eventually`, you can change these values by passing them in just after your function:

```go
Consistently(ACTUAL, DURATION, POLLING_INTERVAL).Should(MATCHER)
```

As with `Eventually`, these can be `time.Duration`s, string representations of a `time.Duration` (e.g. `"200ms"`) or `float64`s that are interpreted as seconds.

Also as with `Eventually`, `Consistently` supports chaining `WithTimeout` and `WithPolling` in the form of:

```go
Consistently(ACTUAL).WithTimeout(DURATION).WithPolling(POLLING_INTERVAL).Should(MATCHER)
```

`Consistently` tries to capture the notion that something "does not eventually" happen.  A common use-case is to assert that no goroutine writes to a channel for a period of time.  If you pass `Consistently` an argument that is not a function, it simply passes that argument to the matcher.  So we can assert that:

```go
Consistently(channel).ShouldNot(Receive())
```

To assert that nothing gets sent to a channel.

As with `Eventually`, you can also pass `Consistently` a function.  In fact, `Consistently` works with the three categories of `ACTUAL` value outlined for `Eventually` in the section above.

> Developers often try to use `runtime.Gosched()` to nudge background goroutines to run.  This can lead to flaky tests as it is not deterministic that a given goroutine will run during the `Gosched`.  `Consistently` is particularly handy in these cases: it polls for 100ms which is typically more than enough time for all your Goroutines to run.  Yes, this is basically like putting a time.Sleep() in your tests... Sometimes, when making negative assertions in a concurrent world, that's the best you can do!

### Modifying Default Intervals

By default, `Eventually` will poll every 10 milliseconds for up to 1 second and `Consistently` will monitor every 10 milliseconds for up to 100 milliseconds.  You can modify these defaults across your test suite with:

```go
SetDefaultEventuallyTimeout(t time.Duration)
SetDefaultEventuallyPollingInterval(t time.Duration)
SetDefaultConsistentlyDuration(t time.Duration)
SetDefaultConsistentlyPollingInterval(t time.Duration)
```

You can also adjust these global timeouts by setting the `GOMEGA_DEFAULT_EVENTUALLY_TIMEOUT`, `GOMEGA_DEFAULT_EVENTUALLY_POLLING_INTERVAL`, `GOMEGA_DEFAULT_CONSISTENTLY_DURATION`, and `GOMEGA_DEFAULT_CONSISTENTLY_POLLING_INTERVAL` environment variables to a parseable duration string.

## Making Assertions in Helper Functions

While writing [custom matchers](#adding-your-own-matchers) is an expressive way to make assertions against your code, it is often more convenient to write one-off helper functions like so:

```go
var _ = Describe("Turbo-encabulator", func() {
    ...
    func assertTurboEncabulatorContains(components ...string) {
        teComponents, err := turboEncabulator.GetComponents()
        Expect(err).NotTo(HaveOccurred())

        Expect(teComponents).To(HaveLen(components))
        for _, component := range components {
            Expect(teComponents).To(ContainElement(component))
        }
    }

    It("should have components", func() {
        assertTurboEncabulatorContains("semi-boloid slots", "grammeters")
    })
})
```

This makes your tests more expressive and reduces boilerplate.  However, when an assertion in the helper fails the line numbers provided by Gomega are unhelpful.  Instead of pointing you to the line in your test that failed, they point you the line in the helper.

To get around this, Gomega provides versions of `Expect`, `Eventually` and `Consistently` named `ExpectWithOffset`, `EventuallyWithOffset` and `ConsistentlyWithOffset` that allow you to specify an *offset* in the call stack.  The offset is the first argument to these functions.

With this, we can rewrite our helper as:

```go
func assertTurboEncabulatorContains(components ...string) {
    teComponents, err := turboEncabulator.GetComponents()
    ExpectWithOffset(1, err).NotTo(HaveOccurred())

    ExpectWithOffset(1, teComponents).To(HaveLen(components))
    for _, component := range components {
      ExpectWithOffset(1, teComponents).To(ContainElement(component))
    }
}
```

Now, failed assertions will point to the correct call to the helper in the test.

Alternatively, you can just use the baseline versions of `Expect`, `Eventually` and `Consistently` and combine them with `WithOffset`:

```go
func assertTurboEncabulatorContains(components ...string) {
    teComponents, err := turboEncabulator.GetComponents()
    Expect(err).WithOffset(1).NotTo(HaveOccurred())

    Expect(teComponents).WithOffset(1).To(HaveLen(components))
    for _, component := range components {
      Expect(teComponents).WithOffset(1).To(ContainElement(component))
    }
}
```

## Provided Matchers

Gomega comes with a bunch of `GomegaMatcher`s.  They're all documented here.  If there's one you'd like to see written either [send a pull request or open an issue](http://github.com/onsi/gomega).

A number of community-supported matchers have appeared as well.  A list is maintained on the Gomega [wiki](https://github.com/onsi/gomega/wiki).

These docs only go over the positive assertion case (`Should`), the negative case (`ShouldNot`) is simply the negation of the positive case.  They also use the `Ω` notation, but - as mentioned above - the `Expect` notation is equivalent.

### Asserting Equivalence

#### Equal(expected interface{})

```go
Ω(ACTUAL).Should(Equal(EXPECTED))
```

uses [`reflect.DeepEqual`](http://golang.org/pkg/reflect#deepequal) to compare `ACTUAL` with `EXPECTED`.

`reflect.DeepEqual` is awesome.  It will use `==` when appropriate (e.g. when comparing primitives) but will recursively dig into maps, slices, arrays, and even your own structs to ensure deep equality.  `reflect.DeepEqual`, however, is strict about comparing types.  Both `ACTUAL` and `EXPECTED` *must* have the same type.  If you want to compare across different types (e.g. if you've defined a type alias) you should use `BeEquivalentTo`.

It is an error for both `ACTUAL` and `EXPECTED` to be nil, you should use `BeNil()` instead.

When both `ACTUAL` and `EXPECTED` are a very long strings, it will attempt to pretty-print the diff and display exactly where they differ.

> For asserting equality between numbers of different types, you'll want to use the [`BeNumerically()`](#benumericallycomparator-string-compareto-interface) matcher.

#### BeEquivalentTo(expected interface{})

```go
Ω(ACTUAL).Should(BeEquivalentTo(EXPECTED))
```

Like `Equal`, `BeEquivalentTo` uses `reflect.DeepEqual` to compare `ACTUAL` with `EXPECTED`.  Unlike `Equal`, however, `BeEquivalentTo` will first convert `ACTUAL`'s type to that of `EXPECTED` before making the comparison with `reflect.DeepEqual`.

This means that `BeEquivalentTo` will successfully match equivalent values of different types.  This is particularly useful, for example, with type aliases:

```go
type FoodSrce string

Ω(FoodSrce("Cheeseboard Pizza")
 ).Should(Equal("Cheeseboard Pizza")) //will fail
Ω(FoodSrce("Cheeseboard Pizza")
 ).Should(BeEquivalentTo("Cheeseboard Pizza")) //will pass
```

As with `Equal` it is an error for both `ACTUAL` and `EXPECTED` to be nil, you should use `BeNil()` instead.

As a rule, you **should not** use `BeEquivalentTo` with numbers.  Both of the following assertions are true:

```go
Ω(5.1).Should(BeEquivalentTo(5))
Ω(5).ShouldNot(BeEquivalentTo(5.1))
```

the first assertion passes because 5.1 will be cast to an integer and will get rounded down!  Such false positives are terrible and should be avoided.  Use [`BeNumerically()`](#benumericallycomparator-string-compareto-interface) to compare numbers instead.

#### BeIdenticalTo(expected interface{})

```go
Ω(ACTUAL).Should(BeIdenticalTo(EXPECTED))
```

Like `Equal`, `BeIdenticalTo` compares `ACTUAL` to `EXPECTED` for equality. Unlike `Equal`, however, it uses `==` to compare values. In practice, this means that primitive values like strings, integers and floats are identical to, as well as pointers to values.

`BeIdenticalTo` is most useful when you want to assert that two pointers point to the exact same location in memory.

As with `Equal` it is an error for both `ACTUAL` and `EXPECTED` to be nil, you should use `BeNil()` instead.

#### BeAssignableToTypeOf(expected interface)

```go
Ω(ACTUAL).Should(BeAssignableToTypeOf(EXPECTED interface))
```

succeeds if `ACTUAL` is a type that can be assigned to a variable with the same type as `EXPECTED`.  It is an error for either `ACTUAL` or `EXPECTED` to be `nil`.

### Asserting Presence

#### BeNil()

```go
Ω(ACTUAL).Should(BeNil())
```

succeeds if `ACTUAL` is, in fact, `nil`.

#### BeZero()

```go
Ω(ACTUAL).Should(BeZero())
```

succeeds if `ACTUAL` is the zero value for its type *or* if `ACTUAL` is `nil`.

### Asserting Truthiness

#### BeTrue()

```go
Ω(ACTUAL).Should(BeTrue())
```

succeeds if `ACTUAL` is `bool` typed and has the value `true`.  It is an error for `ACTUAL` to not be a `bool`.

> Some matcher libraries have a notion of "truthiness" to assert that an object is present.  Gomega is strict, and `BeTrue()` only works with `bool`s.  You can use `Ω(ACTUAL).ShouldNot(BeZero())` or `Ω(ACTUAL).ShouldNot(BeNil())` to verify object presence.

#### BeFalse()

```go
Ω(ACTUAL).Should(BeFalse())
```

succeeds if `ACTUAL` is `bool` typed and has the value `false`.  It is an error for `ACTUAL` to not be a `bool`.

### Asserting on Errors

#### HaveOccurred()

```go
Ω(ACTUAL).Should(HaveOccurred())
```

succeeds if `ACTUAL` is a non-nil `error`.  Thus, the typical Go error checking pattern looks like:

```go
err := SomethingThatMightFail()
Ω(err).ShouldNot(HaveOccurred())
```

#### Succeed()

```go
Ω(ACTUAL).Should(Succeed())
```

succeeds if `ACTUAL` is `nil`.  The intended usage is

```go
Ω(FUNCTION()).Should(Succeed())
```

where `FUNCTION()` is a function call that returns an error-type as its *first or only* return value.  See [Handling Errors](#handling-errors) for a more detailed discussion.

#### MatchError(expected interface{})

```go
Ω(ACTUAL).Should(MatchError(EXPECTED))
```

succeeds if `ACTUAL` is a non-nil `error` that matches `EXPECTED`. `EXPECTED` must be one of the following:

- A string, in which case `ACTUAL.Error()` will be compared against `EXPECTED`.
- A matcher, in which case `ACTUAL.Error()` is tested against the matcher.
- An error, in which case `ACTUAL` and `EXPECTED` are compared via `reflect.DeepEqual()`. If they are not deeply equal, they are tested by `errors.Is(ACTUAL, EXPECTED)`. (The latter allows to test whether `ACTUAL` wraps an `EXPECTED` error.)

Any other type for `EXPECTED` is an error.

### Working with Channels

#### BeClosed()

```go
Ω(ACTUAL).Should(BeClosed())
```

succeeds if `ACTUAL` is a closed channel. It is an error to pass a non-channel to `BeClosed`, it is also an error to pass `nil`.

In order to check whether or not the channel is closed, Gomega must try to read from the channel (even in the `ShouldNot(BeClosed())` case).  You should keep this in mind if you wish to make subsequent assertions about values coming down the channel.

Also, if you are testing that a *buffered* channel is closed you must first read all values out of the channel before asserting that it is closed (it is not possible to detect that a buffered-channel has been closed until all its buffered values are read).

Finally, as a corollary: it is an error to check whether or not a send-only channel is closed.

#### Receive()

```go
Ω(ACTUAL).Should(Receive(<optionalPointer>))
```

succeeds if there is a message to be received on actual. Actual must be a channel (and cannot be a send-only channel) -- anything else is an error.

`Receive` returns *immediately*.  It *never* blocks:

- If there is nothing on the channel `c` then `Ω(c).Should(Receive())` will fail and `Ω(c).ShouldNot(Receive())` will pass.
- If there is something on the channel `c` ready to be read, then `Ω(c).Should(Receive())` will pass and `Ω(c).ShouldNot(Receive())` will fail.
- If the channel `c` is closed then `Ω(c).Should(Receive())` will fail and `Ω(c).ShouldNot(Receive())` will pass.

If you have a go-routine running in the background that will write to channel `c`, for example:

```go
go func() {
    time.Sleep(100 * time.Millisecond)
    c <- true
}()
```

you can assert that `c` receives something (anything!) eventually:

```go
Eventually(c).Should(Receive())
```

This will timeout if nothing gets sent to `c` (you can modify the timeout interval as you normally do with `Eventually`).

A similar use-case is to assert that no go-routine writes to a channel (for a period of time).  You can do this with `Consistently`:

```go
Consistently(c).ShouldNot(Receive())
```

`Receive` also allows you to make assertions on the received object.  You do this by passing `Receive` a matcher:

```go
Eventually(c).Should(Receive(Equal("foo")))
```

This assertion will only succeed if `c` receives an object *and* that object satisfies `Equal("foo")`.  Note that `Eventually` will continually poll `c` until this condition is met.  If there are objects coming down the channel that do not satisfy the passed in matcher, they will be pulled off and discarded until an object that *does* satisfy the matcher is received.

In addition, there are occasions when you need to grab the object sent down the channel (e.g. to make several assertions against the object).  To do this, you can ask the `Receive` matcher for the value passed to the channel by passing it a pointer to a variable of the appropriate type:

```go
var receivedBagel Bagel
Eventually(bagelChan).Should(Receive(&receivedBagel))
Ω(receivedBagel.Contents()).Should(ContainElement("cream cheese"))
Ω(receivedBagel.Kind()).Should(Equal("sesame"))
```

Of course, this could have been written as `receivedBagel := <-bagelChan` - however using `Receive` makes it easy to avoid hanging the test suite should nothing ever come down the channel. The pointer can point to any variable whose type is assignable from the channel element type, or if the channel type is an interface and the underlying type is assignable to the pointer.

Finally, `Receive` *never* blocks.  `Eventually(c).Should(Receive())` repeatedly polls `c` in a non-blocking fashion.  That means that you cannot use this pattern to verify that a *non-blocking send* has occurred on the channel - [more details at this GitHub issue](https://github.com/onsi/gomega/issues/82).

#### BeSent(value interface{})

```go
Ω(ACTUAL).Should(BeSent(VALUE))
```

attempts to send `VALUE` to the channel `ACTUAL` without blocking.  It succeeds if this is possible.

`ACTUAL` must be a channel (and cannot be a receive-only channel) that can be sent the type of the `VALUE` passed into `BeSent` -- anything else is an error. In addition, `ACTUAL` must not be closed.

`BeSent` never blocks:

- If the channel `c` is not ready to receive then `Ω(c).Should(BeSent("foo"))` will fail immediately.
- If the channel `c` is eventually ready to receive then `Eventually(c).Should(BeSent("foo"))` will succeed... presuming the channel becomes ready to receive before `Eventually`'s timeout.
- If the channel `c` is closed then `Ω(c).Should(BeSent("foo"))` and `Ω(c).ShouldNot(BeSent("foo"))` will both fail immediately.

Of course, `VALUE` is actually sent to the channel.  The point of `BeSent` is less to make an assertion about the availability of the channel (which is typically an implementation detail that your test should not be concerned with). Rather, the point of `BeSent` is to make it possible to easily and expressively write tests that can timeout on blocked channel sends.

### Working with files

#### BeAnExistingFile

```go
Ω(ACTUAL).Should(BeAnExistingFile())
```

succeeds if a file located at `ACTUAL` exists.

`ACTUAL` must be a string representing the filepath.

#### BeARegularFile

```go
Ω(ACTUAL).Should(BeARegularFile())
```

succeeds IFF a file located at `ACTUAL` exists and is a regular file.

`ACTUAL` must be a string representing the filepath.

#### BeADirectory

```go
Ω(ACTUAL).Should(BeADirectory())
```

succeeds IFF a file is located at `ACTUAL` exists and is a directory.

`ACTUAL` must be a string representing the filepath.

### Working with Strings, JSON and YAML

#### ContainSubstring(substr string, args ...interface{})

```go
Ω(ACTUAL).Should(ContainSubstring(STRING, ARGS...))
```

succeeds if `ACTUAL` contains the substring generated by:

```go
fmt.Sprintf(STRING, ARGS...)
```

`ACTUAL` must either be a `string`, `[]byte` or a `Stringer` (a type implementing the `String()` method).  Any other input is an error.

> Note, of course, that the `ARGS...` are not required.  They are simply a convenience to allow you to build up strings programmatically inline in the matcher.

#### HavePrefix(prefix string, args ...interface{})

```go
Ω(ACTUAL).Should(HavePrefix(STRING, ARGS...))
```

succeeds if `ACTUAL` has the string prefix generated by:

```go
fmt.Sprintf(STRING, ARGS...)
```

`ACTUAL` must either be a `string`, `[]byte` or a `Stringer` (a type implementing the `String()` method).  Any other input is an error.

> Note, of course, that the `ARGS...` are not required.  They are simply a convenience to allow you to build up strings programmatically inline in the matcher.

#### HaveSuffix(suffix string, args ...interface{})

```go
Ω(ACTUAL).Should(HaveSuffix(STRING, ARGS...))
```

succeeds if `ACTUAL` has the string suffix generated by:

```go
fmt.Sprintf(STRING, ARGS...)
```

`ACTUAL` must either be a `string`, `[]byte` or a `Stringer` (a type implementing the `String()` method).  Any other input is an error.

> Note, of course, that the `ARGS...` are not required.  They are simply a convenience to allow you to build up strings programmatically inline in the matcher.

#### MatchRegexp(regexp string, args ...interface{})

```go
Ω(ACTUAL).Should(MatchRegexp(STRING, ARGS...))
```

succeeds if `ACTUAL` is matched by the regular expression string generated by:

```go
fmt.Sprintf(STRING, ARGS...)
```

`ACTUAL` must either be a `string`, `[]byte` or a `Stringer` (a type implementing the `String()` method).  Any other input is an error.  It is also an error for the regular expression to fail to compile.

> Note, of course, that the `ARGS...` are not required.  They are simply a convenience to allow you to build up strings programmatically inline in the matcher.

#### MatchJSON(json interface{})

```go
Ω(ACTUAL).Should(MatchJSON(EXPECTED))
```

Both `ACTUAL` and `EXPECTED` must be a `string`, `[]byte` or a `Stringer`.  `MatchJSON` succeeds if both `ACTUAL` and `EXPECTED` are JSON representations of the same object.  This is verified by parsing both `ACTUAL` and `EXPECTED` and then asserting equality on the resulting objects with `reflect.DeepEqual`.  By doing this `MatchJSON` avoids any issues related to white space, formatting, and key-ordering.

It is an error for either `ACTUAL` or `EXPECTED` to be invalid JSON.

In some cases it is useful to match two JSON strings while ignoring list order.  For this you can use the community maintained [MatchUnorderedJSON](https://github.com/Benjamintf1/Expanded-Unmarshalled-Matchers) matcher.

#### MatchXML(xml interface{})

```go
Ω(ACTUAL).Should(MatchXML(EXPECTED))
```

Both `ACTUAL` and `EXPECTED` must be a `string`, `[]byte` or a `Stringer`.  `MatchXML` succeeds if both `ACTUAL` and `EXPECTED` are XML representations of the same object.  This is verified by parsing both `ACTUAL` and `EXPECTED` and then asserting equality on the resulting objects with `reflect.DeepEqual`.  By doing this `MatchXML` avoids any issues related to white space or formatting.

It is an error for either `ACTUAL` or `EXPECTED` to be invalid XML.

#### MatchYAML(yaml interface{})

```go
Ω(ACTUAL).Should(MatchYAML(EXPECTED))
```

Both `ACTUAL` and `EXPECTED` must be a `string`, `[]byte` or a `Stringer`.  `MatchYAML` succeeds if both `ACTUAL` and `EXPECTED` are YAML representations of the same object.  This is verified by parsing both `ACTUAL` and `EXPECTED` and then asserting equality on the resulting objects with `reflect.DeepEqual`.  By doing this `MatchYAML` avoids any issues related to white space, formatting, and key-ordering.

It is an error for either `ACTUAL` or `EXPECTED` to be invalid YAML.

### Working with Collections

#### BeEmpty()

```go
Ω(ACTUAL).Should(BeEmpty())
```

succeeds if `ACTUAL` is, in fact, empty. `ACTUAL` must be of type `string`, `array`, `map`, `chan`, or `slice`.  It is an error for it to have any other type.

#### HaveLen(count int)

```go
Ω(ACTUAL).Should(HaveLen(INT))
```

succeeds if the length of `ACTUAL` is `INT`. `ACTUAL` must be of type `string`, `array`, `map`, `chan`, or `slice`.  It is an error for it to have any other type.

#### HaveCap(count int)

```go
Ω(ACTUAL).Should(HaveCap(INT))
```

succeeds if the capacity of `ACTUAL` is `INT`. `ACTUAL` must be of type `array`, `chan`, or `slice`.  It is an error for it to have any other type.

#### ContainElement(element interface{})

```go
Ω(ACTUAL).Should(ContainElement(ELEMENT))
```

succeeds if `ACTUAL` contains an element that equals `ELEMENT`.  `ACTUAL` must be an `array`, `slice`, or `map` -- anything else is an error.  For `map`s `ContainElement` searches through the map's values (not keys!).

By default `ContainElement()` uses the `Equal()` matcher under the hood to assert equality between `ACTUAL`'s elements and `ELEMENT`.  You can change this, however, by passing `ContainElement` a `GomegaMatcher`. For example, to check that a slice of strings has an element that matches a substring:

```go
Ω([]string{"Foo", "FooBar"}).Should(ContainElement(ContainSubstring("Bar")))
```

#### ContainElements(element ...interface{})

```go
Ω(ACTUAL).Should(ContainElements(ELEMENT1, ELEMENT2, ELEMENT3, ...))
```

or

```go
Ω(ACTUAL).Should(ContainElements([]SOME_TYPE{ELEMENT1, ELEMENT2, ELEMENT3, ...}))
```

succeeds if `ACTUAL` contains the elements passed into the matcher. The ordering of the elements does not matter.

By default `ContainElements()` uses `Equal()` to match the elements, however custom matchers can be passed in instead.  Here are some examples:

```go
Ω([]string{"Foo", "FooBar"}).Should(ContainElements("FooBar"))
Ω([]string{"Foo", "FooBar"}).Should(ContainElements(ContainSubstring("Bar"), "Foo"))
```

Actual must be an `array`, `slice` or `map`.  For maps, `ContainElements` matches against the `map`'s values.

You typically pass variadic arguments to `ContainElements` (as in the examples above).  However, if you need to pass in a slice you can provided that it
is the only element passed in to `ContainElements`:

```go
Ω([]string{"Foo", "FooBar"}).Should(ContainElements([]string{"FooBar", "Foo"}))
```

Note that Go's type system does not allow you to write this as `ContainElements([]string{"FooBar", "Foo"}...)` as `[]string` and `[]interface{}` are different types - hence the need for this special rule.

The difference between the `ContainElements` and `ConsistOf` matchers is that the latter is more restrictive because the `ConsistOf` matcher checks additionally that the `ACTUAL` elements and the elements passed into the matcher have the same length.

#### BeElementOf(elements ...interface{})

```go
Ω(ACTUAL).Should(BeElementOf(ELEMENT1, ELEMENT2, ELEMENT3, ...))
```

succeeds if `ACTUAL` equals one of the elements passed into the matcher. When a single element `ELEMENT` of type `array` or `slice` is passed into the matcher, `BeElementOf` succeeds if `ELEMENT` contains an element that equals `ACTUAL` (reverse of `ContainElement`). `BeElementOf` always uses the `Equal()` matcher under the hood to assert equality.

#### ConsistOf(element ...interface{})

```go
Ω(ACTUAL).Should(ConsistOf(ELEMENT1, ELEMENT2, ELEMENT3, ...))
```

or

```go
Ω(ACTUAL).Should(ConsistOf([]SOME_TYPE{ELEMENT1, ELEMENT2, ELEMENT3, ...}))
```

succeeds if `ACTUAL` contains precisely the elements passed into the matcher. The ordering of the elements does not matter.

By default `ConsistOf()` uses `Equal()` to match the elements, however custom matchers can be passed in instead.  Here are some examples:

```go
Ω([]string{"Foo", "FooBar"}).Should(ConsistOf("FooBar", "Foo"))
Ω([]string{"Foo", "FooBar"}).Should(ConsistOf(ContainSubstring("Bar"), "Foo"))
Ω([]string{"Foo", "FooBar"}).Should(ConsistOf(ContainSubstring("Foo"), ContainSubstring("Foo")))
```

Actual must be an `array`, `slice` or `map`.  For maps, `ConsistOf` matches against the `map`'s values.

You typically pass variadic arguments to `ConsistOf` (as in the examples above).  However, if you need to pass in a slice you can provided that it
is the only element passed in to `ConsistOf`:

```go
Ω([]string{"Foo", "FooBar"}).Should(ConsistOf([]string{"FooBar", "Foo"}))
```

Note that Go's type system does not allow you to write this as `ConsistOf([]string{"FooBar", "Foo"}...)` as `[]string` and `[]interface{}` are different types - hence the need for this special rule.


#### HaveKey(key interface{})

```go
Ω(ACTUAL).Should(HaveKey(KEY))
```

succeeds if `ACTUAL` is a map with a key that equals `KEY`.  It is an error for `ACTUAL` to not be a `map`.

By default `HaveKey()` uses the `Equal()` matcher under the hood to assert equality between `ACTUAL`'s keys and `KEY`.  You can change this, however, by passing `HaveKey` a `GomegaMatcher`. For example, to check that a map has a key that matches a regular expression:

```go
Ω(map[string]string{"Foo": "Bar", "BazFoo": "Duck"}).Should(HaveKey(MatchRegexp(`.+Foo$`)))
```

#### HaveKeyWithValue(key interface{}, value interface{})

```go
Ω(ACTUAL).Should(HaveKeyWithValue(KEY, VALUE))
```

succeeds if `ACTUAL` is a map with a key that equals `KEY` mapping to a value that equals `VALUE`.  It is an error for `ACTUAL` to not be a `map`.

By default `HaveKeyWithValue()` uses the `Equal()` matcher under the hood to assert equality between `ACTUAL`'s keys and `KEY` and between the associated value and `VALUE`.  You can change this, however, by passing `HaveKeyWithValue` a `GomegaMatcher` for either parameter. For example, to check that a map has a key that matches a regular expression and which is also associated with a value that passes some numerical threshold:

```go
Ω(map[string]int{"Foo": 3, "BazFoo": 4}).Should(HaveKeyWithValue(MatchRegexp(`.+Foo$`), BeNumerically(">", 3)))
```

#### HaveField(field interface{}, value interface{})

```go
Ω(ACTUAL).Should(HaveField(FIELD, VALUE))
```

succeeds if `ACTUAL` is a struct with a value that can be traversed via `FIELD` that equals `VALUE`.  It is an error for `ACTUAL` to not be a `struct`.

By default `HaveField()` uses the `Equal()` matcher under the hood to assert equality between the extracted value and `VALUE`.  You can change this, however, by passing `HaveField` a `GomegaMatcher` for `VALUE`.

`FIELD` allows you to access fields within the `ACTUAL` struct.  Nested structs can be accessed using the `.` delimiter.  `HaveField()` also allows you to invoke methods on the struct by adding a `()` suffix to the `FIELD` - these methods must take no arguments and return exactly one value.  For example consider the following types:

```go
type Book struct {
    Title string
    Author Person
}

type Person struct {
    Name string
    DOB time.Time
}
```

and an instance book `var book = Book{...}` - you can use `HaveField` to make assertions like:

```go
Ω(book).Should(HaveField("Title", "Les Miserables"))
Ω(book).Should(HaveField("Title", ContainSubstring("Les Mis")))
Ω(book).Should(HaveField("Author.Name", "Victor Hugo"))
Ω(book).Should(HaveField("Author.DOB.Year()", BeNumerically("<", 1900)))
```

`HaveField` can pair powerfully with a collection matcher like `ContainElement`.  To assert that a list of books as at least one element with an author born in February you could write:

```go
Ω(books).Should(ContainElement(HaveField("Author.DOB.Month()", Equal(2))))
```

If you want to make lots of complex assertions against the fields of a struct take a look at the [`gstruct`package](#gstruct-testing-complex-data-types) package documented below.  

### Working with Numbers and Times

#### BeNumerically(comparator string, compareTo ...interface{})

```go
Ω(ACTUAL).Should(BeNumerically(COMPARATOR_STRING, EXPECTED, <THRESHOLD>))
```

performs numerical assertions in a type-agnostic way.  `ACTUAL` and `EXPECTED` should be numbers, though the specific type of number is irrelevant (`float32`, `float64`, `uint8`, etc...).  It is an error for `ACTUAL` or `EXPECTED` to not be a number.

There are six supported comparators:

- `Ω(ACTUAL).Should(BeNumerically("==", EXPECTED))`:
    asserts that `ACTUAL` and `EXPECTED` are numerically equal.

- `Ω(ACTUAL).Should(BeNumerically("~", EXPECTED, <THRESHOLD>))`:
    asserts that `ACTUAL` and `EXPECTED` are within `<THRESHOLD>` of one another.  By default `<THRESHOLD>` is `1e-8` but you can specify a custom value.

- `Ω(ACTUAL).Should(BeNumerically(">", EXPECTED))`:
    asserts that `ACTUAL` is greater than `EXPECTED`.

- `Ω(ACTUAL).Should(BeNumerically(">=", EXPECTED))`:
    asserts that `ACTUAL` is greater than or equal to  `EXPECTED`.

- `Ω(ACTUAL).Should(BeNumerically("<", EXPECTED))`:
    asserts that `ACTUAL` is less than `EXPECTED`.

- `Ω(ACTUAL).Should(BeNumerically("<=", EXPECTED))`:
    asserts that `ACTUAL` is less than or equal to `EXPECTED`.

Any other comparator is an error.

#### BeTemporally(comparator string, compareTo time.Time, threshold ...time.Duration)

```go
Ω(ACTUAL).Should(BeTemporally(COMPARATOR_STRING, EXPECTED_TIME, <THRESHOLD_DURATION>))
```

performs time-related assertions.  `ACTUAL` must be a `time.Time`.

There are six supported comparators:

- `Ω(ACTUAL).Should(BeTemporally("==", EXPECTED_TIME))`:
    asserts that `ACTUAL` and `EXPECTED_TIME` are identical `time.Time`s.

- `Ω(ACTUAL).Should(BeTemporally("~", EXPECTED_TIME, <THRESHOLD_DURATION>))`:
    asserts that `ACTUAL` and `EXPECTED_TIME` are within `<THRESHOLD_DURATION>` of one another.  By default `<THRESHOLD_DURATION>` is `time.Millisecond` but you can specify a custom value.

- `Ω(ACTUAL).Should(BeTemporally(">", EXPECTED_TIME))`:
    asserts that `ACTUAL` is after `EXPECTED_TIME`.

- `Ω(ACTUAL).Should(BeTemporally(">=", EXPECTED_TIME))`:
    asserts that `ACTUAL` is after or at `EXPECTED_TIME`.

- `Ω(ACTUAL).Should(BeTemporally("<", EXPECTED_TIME))`:
    asserts that `ACTUAL` is before `EXPECTED_TIME`.

- `Ω(ACTUAL).Should(BeTemporally("<=", EXPECTED_TIME))`:
    asserts that `ACTUAL` is before or at `EXPECTED_TIME`.

Any other comparator is an error.

### Working with Values

#### HaveValue(matcher types.GomegaMatcher)

`HaveValue` applies `MATCHER` to the value that results from dereferencing `ACTUAL` in case of a pointer or an interface, or otherwise `ACTUAL` itself. Pointers and interfaces are dereferenced multiple times as necessary, with a limit of at most 31 dereferences.

```go
Expect(ACTUAL).To(HaveValue(MATCHER))
```

For instance:

```go
i := 42
Expect(&i).To(HaveValue(Equal(42)))
Expect(i).To(HaveValue(Equal(42)))
```

`HaveValue` can be used, for instance, in tests and custom matchers where the it doesn't matter (as opposed to `PointTo`) if a value first needs to be dereferenced or not. This is especially useful to custom matchers that are to be used in mixed contexts of pointers as well as non-pointers.

### Working with HTTP responses

#### HaveHTTPStatus(expected interface{})

```go
  Expect(ACTUAL).To(HaveHTTPStatus(EXPECTED, ...))
```

succeeds if the `Status` or `StatusCode` field of an HTTP response matches.

`ACTUAL` must be either a `*http.Response` or `*httptest.ResponseRecorder`.

`EXPECTED` must be one or more `int` or `string` types. An `int` is compared
to `StatusCode` and a `string` is compared to `Status`.
The matcher succeeds if any of the `EXPECTED` values match.

Here are some examples:

- `Expect(resp).To(HaveHTTPStatus(http.StatusOK, http.StatusNoContext))`:
    asserts that `resp.StatusCode == 200` or `resp.StatusCode == 204`

- `Expect(resp).To(HaveHTTPStatus("404 Not Found"))`:
    asserts that `resp.Status == "404 Not Found"`.

#### HaveHTTPBody(expected interface{})

```go
Expect(ACTUAL).To(HaveHTTPBody(EXPECTED))
```

Succeeds if the body of an HTTP Response matches.

`ACTUAL` must be either a `*http.Response` or `*httptest.ResponseRecorder`.

`EXPECTED` must be one of the following:
- A `string`
- A `[]byte`
- A matcher, in which case the matcher will be called with the body as a `[]byte`.

Here are some examples:

- `Expect(resp).To(HaveHTTPBody("bar"))`:
    asserts that when `resp.Body` is read, it will equal `bar`.

- `Expect(resp).To(HaveHTTPBody(MatchJSON("{\"some\":\"json\""))`:
    asserts that when `resp.Body` is read, the `MatchJSON` matches it to `{"some":"json"}`.

Note that the body is an `io.ReadCloser` and the `HaveHTTPBody()` will read it and the close it.
This means that subsequent attempts to read the body may have unexpected results.

#### HaveHTTPHeaderWithValue(key string, value interface{})

```go
Expect(ACTUAL).To(HaveHTTPHeaderWithValue(KEY, VALUE))
```

Succeeds if the HTTP Response has a matching header and value.

`ACTUAL` must be either a `*http.Response` or `*httptest.ResponseRecorder`.

`KEY` must be a `string`. It is passed to
[`http.Header.Get(key string)`](https://pkg.go.dev/net/http#Header.Get),
and will have the same behaviors regarding order of headers and capitalization.

`VALUE` must be one of the following:
- A `string`
- A matcher, in which case the matcher will be called to match the value.

Here are some examples:

- `Expect(resp).To(HaveHTTPHeaderWithValue("Content-Type", "application/json"))`:
    asserts that the `Content-Type` header has exactly the value `application/json`.

- `Expect(resp).To(HaveHTTPHeaderWithValue(ContainsSubstring("json")))`:
    asserts that the `Content-Type` header contains the substring `json`.

### Asserting on Panics

#### Panic()

```go
Ω(ACTUAL).Should(Panic())
```

succeeds if `ACTUAL` is a function that, when invoked, panics.  `ACTUAL` must be a function that takes no arguments and returns no result -- any other type for `ACTUAL` is an error.

#### PanicWith()

```go
Ω(ACTUAL).Should(PanicWith(VALUE))
```

succeeds if `ACTUAL` is a function that, when invoked, panics with a value of `VALUE`.  `ACTUAL` must be a function that takes no arguments and returns no result -- any other type for `ACTUAL` is an error.

By default `PanicWith()` uses the `Equal()` matcher under the hood to assert equality between `ACTUAL`'s panic value and `VALUE`.  You can change this, however, by passing `PanicWith` a `GomegaMatcher`. For example, to check that the panic value matches a regular expression:

```go
Ω(func() { panic("FooBarBaz") }).Should(PanicWith(MatchRegexp(`.+Baz$`)))
```

### Composing Matchers

You may form larger matcher expressions using the following operators: `And()`, `Or()`, `Not()` and `WithTransform()`.

Note: `And()` and `Or()` can also be referred to as `SatisfyAll()` and `SatisfyAny()`, respectively.

With these operators you can express multiple requirements in a single `Expect()` or `Eventually()` statement. For example:

```go
Expect(number).To(SatisfyAll(
            BeNumerically(">", 0),
            BeNumerically("<", 10)))

Expect(msg).To(SatisfyAny(
            Equal("Success"),
            MatchRegexp(`^Error .+$`)))
```

It can also provide a lightweight syntax to create new matcher types from existing ones. For example:

```go
func BeBetween(min, max int) GomegaMatcher {
    return SatisfyAll(
            BeNumerically(">", min),
            BeNumerically("<", max))
}

Ω(number).Should(BeBetween(0, 10))
```

#### And(matchers ...GomegaMatcher)

#### SatisfyAll(matchers ...GomegaMatcher)

```go
Ω(ACTUAL).Should(And(MATCHER1, MATCHER2, ...))
```

or

```go
Ω(ACTUAL).Should(SatisfyAll(MATCHER1, MATCHER2, ...))
```

succeeds if `ACTUAL` satisfies all of the specified matchers (similar to a logical AND).

Tests the given matchers in order, returning immediately if one fails, without needing to test the remaining matchers.

#### Or(matchers ...GomegaMatcher)

#### SatisfyAny(matchers ...GomegaMatcher)

```go
Ω(ACTUAL).Should(Or(MATCHER1, MATCHER2, ...))
```

or

```go
Ω(ACTUAL).Should(SatisfyAny(MATCHER1, MATCHER2, ...))
```

succeeds if `ACTUAL` satisfies any of the specified matchers (similar to a logical OR).

Tests the given matchers in order, returning immediately if one succeeds, without needing to test the remaining matchers.

#### Not(matcher GomegaMatcher)

```go
Ω(ACTUAL).Should(Not(MATCHER))
```

succeeds if `ACTUAL` does **not** satisfy the specified matcher (similar to a logical NOT).

#### WithTransform(transform interface{}, matcher GomegaMatcher)

```go
Ω(ACTUAL).Should(WithTransform(TRANSFORM, MATCHER))
```

succeeds if applying the `TRANSFORM` function to `ACTUAL` (i.e. the value of `TRANSFORM(ACTUAL)`) will satisfy the given `MATCHER`. For example:

```go
GetColor := func(e Element) Color { return e.Color }

Ω(element).Should(WithTransform(GetColor, Equal(BLUE)))
```

Or the same thing expressed by introducing a new, lightweight matcher:

```go
// HaveColor returns a matcher that expects the element to have the given color.
func HaveColor(c Color) GomegaMatcher {
    return WithTransform(func(e Element) Color {
        return e.Color
    }, Equal(c))
}

Ω(element).Should(HaveColor(BLUE)))
```

`TRANSFORM` functions optionally can return an additional error value in case a transformation is not possible, avoiding the need to `panic`. Returning errors can be useful when using `WithTransform` to build lightweight matchers that accept different value types and that can gracefully fail when presented the wrong value type.

As before, such a `TRANSFORM` expects a single actual value. But now it returns the transformed value together with an error value. This follows the common Go idiom to communicate errors via an explicit, separate return value.

The following lightweight matcher expects to be used either on a `Sprocket` value or `*Sprocket` pointer. It gracefully fails when the actual value is something else.

```go
// HaveSprocketName returns a matcher that expects the actual value to be
// either a Sprocket or a *Sprocket, having the specified name.
func HaveSprocketName(name string) GomegaMatcher {
    return WithTransform(
        func(actual interface{}) (string, error) {
            switch sprocket := actual.(type) {
            case *Sprocket:
                return Sprocket.Name, nil
            case Sprocket:
                return Sprocket.Name, nil
            default:
                return "", fmt.Errorf("HaveSprocketName expects a Sprocket or *Sprocket, but got %T", actual)
            }
        }, Equal(name))
}

Ω(element).Should(HaveSprocketName("gomega")))
```

#### Satisfy(predicate interface{})

```go
Ω(ACTUAL).Should(Satisfy(PREDICATE))
```

succeeds if applying the `PREDICATE` function to `ACTUAL` (i.e. the value of `PREDICATE(ACTUAL)`) will return `true`. For example:

```go
IsEven := func(i int) bool { return i%2==0 }

Ω(number).Should(Satisfy(IsEven))
```

## Adding Your Own Matchers

A matcher, in Gomega, is any type that satisfies the `GomegaMatcher` interface:

```go
type GomegaMatcher interface {
    Match(actual interface{}) (success bool, err error)
    FailureMessage(actual interface{}) (message string)
    NegatedFailureMessage(actual interface{}) (message string)
}
```

For the simplest cases, new matchers can be [created by composition](#composing-matchers).  In addition to this chapter, please take a look at the [Building Custom Matchers](https://onsi.github.io/ginkgo/#building-custom-matchers) section of the Ginkgo and Gomega patterns chapter in the Ginkgo docs.  Gomega's building blocks have evolved since the Gomega docs were written and while this section remains valid - the [Building Custom Matchers](https://onsi.github.io/ginkgo/#building-custom-matchers) docs present a modern way to more quickly construct custom matchers.

But writing domain-specific custom matchers is also trivial and highly encouraged.  Let's work through an example.

> The `GomegaMatcher` interface is defined in the `types` subpackage.

### A Custom Matcher: RepresentJSONifiedObject(EXPECTED interface{})

Say you're working on a JSON API and you want to assert that your server returns the correct JSON representation.  Rather than marshal/unmarshal JSON in your tests, you want to write an expressive matcher that checks that the received response is a JSON representation for the object in question.  This is what the `RepresentJSONifiedObject` matcher could look like:

```go
package json_response_matcher

import (
    "github.com/onsi/gomega/types"

    "encoding/json"
    "fmt"
    "net/http"
    "reflect"
)

func RepresentJSONifiedObject(expected interface{}) types.GomegaMatcher {
    return &representJSONMatcher{
        expected: expected,
    }
}

type representJSONMatcher struct {
    expected interface{}
}

func (matcher *representJSONMatcher) Match(actual interface{}) (success bool, err error) {
    response, ok := actual.(*http.Response)
    if !ok {
        return false, fmt.Errorf("RepresentJSONifiedObject matcher expects an http.Response")
    }

    pointerToObjectOfExpectedType := reflect.New(reflect.TypeOf(matcher.expected)).Interface()
    err = json.NewDecoder(response.Body).Decode(pointerToObjectOfExpectedType)

    if err != nil {
        return false, fmt.Errorf("Failed to decode JSON: %s", err.Error())
    }

    decodedObject := reflect.ValueOf(pointerToObjectOfExpectedType).Elem().Interface()

    return reflect.DeepEqual(decodedObject, matcher.expected), nil
}

func (matcher *representJSONMatcher) FailureMessage(actual interface{}) (message string) {
    return fmt.Sprintf("Expected\n\t%#v\nto contain the JSON representation of\n\t%#v", actual, matcher.expected)
}

func (matcher *representJSONMatcher) NegatedFailureMessage(actual interface{}) (message string) {
    return fmt.Sprintf("Expected\n\t%#v\nnot to contain the JSON representation of\n\t%#v", actual, matcher.expected)
}
```

Let's break this down:

- Most matchers have a constructor function that returns an instance of the matcher.  In this case we've created `RepresentJSONifiedObject`.  Where possible, your constructor function should take explicit types or interfaces.  For our use case, however, we need to accept any possible expected type so `RepresentJSONifiedObject` takes an argument with the generic `interface{}` type.
- The constructor function then initializes and returns an instance of our matcher: the `representJSONMatcher`.  These rarely need to be exported outside of your matcher package.
- The `representJSONMatcher` must satisfy the `GomegaMatcher` interface.  It does this by implementing the `Match`, `FailureMessage`, and `NegatedFailureMessage` method:
    - If the `GomegaMatcher` receives invalid inputs `Match` returns a non-nil error explaining the problems with the input.  This allows Gomega to fail the assertion whether the assertion is for the positive or negative case.
    - If the `actual` and `expected` values match, `Match` should return `true`.
    - Similarly, if the `actual` and `expected` values do not match, `Match` should return `false`.
    - If the `GomegaMatcher` was testing the `Should` case, and `Match` returned `false`, `FailureMessage` will be called to print a message explaining the failure.
    - Likewise, if the `GomegaMatcher` was testing the `ShouldNot` case, and `Match` returned `true`, `NegatedFailureMessage` will be called.
    - It is guaranteed that `FailureMessage` and `NegatedFailureMessage` will only be called *after* `Match`, so you can save off any state you need to compute the messages in `Match`.
- Finally, it is common for matchers to make extensive use of the `reflect` library to interpret the generic inputs they receive.  In this case, the `representJSONMatcher` goes through some `reflect` gymnastics to create a pointer to a new object with the same type as the `expected` object, read and decode JSON from `actual` into that pointer, and then deference the pointer and compare the result to the `expected` object.

You might test drive this matcher while writing it using Ginkgo.  Your test might look like:

```go
package json_response_matcher_test

import (
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
    . "jsonresponsematcher"

    "bytes"
    "encoding/json"
    "io/ioutil"
    "net/http"
    "strings"

    "testing"
)

func TestCustomMatcher(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Custom Matcher Suite")
}

type Book struct {
    Title  string `json:"title"`
    Author string `json:"author"`
}

var _ = Describe("RepresentJSONified Object", func() {
    var (
        book     Book
        bookJSON []byte
        response *http.Response
    )

    BeforeEach(func() {
        book = Book{
            Title:  "Les Miserables",
            Author: "Victor Hugo",
        }

        var err error
        bookJSON, err = json.Marshal(book)
        Ω(err).ShouldNot(HaveOccurred())
    })

    Context("when actual is not an http response", func() {
        It("should error", func() {
            _, err := RepresentJSONifiedObject(book).Match("not a response")
            Ω(err).Should(HaveOccurred())
        })
    })

    Context("when actual is an http response", func() {
        BeforeEach(func() {
            response = &http.Response{}
        })

        Context("with a body containing the JSON representation of actual", func() {
            BeforeEach(func() {
                response.ContentLength = int64(len(bookJSON))
                response.Body = ioutil.NopCloser(bytes.NewBuffer(bookJSON))
            })

            It("should succeed", func() {
                Ω(response).Should(RepresentJSONifiedObject(book))
            })
        })

        Context("with a body containing the JSON representation of something else", func() {
            BeforeEach(func() {
                reader := strings.NewReader(`{}`)
                response.ContentLength = int64(reader.Len())
                response.Body = ioutil.NopCloser(reader)
            })

            It("should fail", func() {
                Ω(response).ShouldNot(RepresentJSONifiedObject(book))
            })
        })

        Context("with a body containing invalid JSON", func() {
            BeforeEach(func() {
                reader := strings.NewReader(`floop`)
                response.ContentLength = int64(reader.Len())
                response.Body = ioutil.NopCloser(reader)
            })

            It("should error", func() {
                _, err := RepresentJSONifiedObject(book).Match(response)
                Ω(err).Should(HaveOccurred())
            })
        })
    })
})
```

This also offers an example of what using the matcher would look like in your tests.  Note that testing the cases when the matcher returns an error involves creating the matcher and invoking `Match` manually (instead of using an `Ω` or `Expect` assertion).

### Aborting Eventually/Consistently

There are sometimes instances where `Eventually` or `Consistently` should stop polling a matcher because the result of the match simply cannot change.

For example, consider a test that looks like:

```go
Eventually(myChannel).Should(Receive(Equal("bar")))
```

`Eventually` will repeatedly invoke the `Receive` matcher against `myChannel` until the match succeeds.  However, if the channel becomes *closed* there is *no way* for the match to ever succeed.  Allowing `Eventually` to continue polling is inefficient and slows the test suite down.

To get around this, a matcher can optionally implement:

```go
MatchMayChangeInTheFuture(actual interface{}) bool
```

This is not part of the `GomegaMatcher` interface and, in general, most matchers do not need to implement `MatchMayChangeInTheFuture`.

If implemented, however, `MatchMayChangeInTheFuture` will be called with the appropriate `actual` value by `Eventually` and `Consistently` *after* the call to `Match` during every polling interval.  If `MatchMayChangeInTheFuture` returns `true`, `Eventually` and `Consistently` will continue polling.  If, however, `MatchMayChangeInTheFuture` returns `false`, `Eventually` and `Consistently` will stop polling and either fail or pass as appropriate.

If you'd like to look at a simple example of `MatchMayChangeInTheFuture` check out [`gexec`'s `Exit` matcher](https://github.com/onsi/gomega/tree/master/gexec/exit_matcher.go).  Here, `MatchMayChangeInTheFuture` returns true if the `gexec.Session` under test has not exited yet, but returns false if it has.  Because of this, if a process exits with status code 3, but an assertion is made of the form:

```go
Eventually(session, 30).Should(gexec.Exit(0))
```

`Eventually` will not block for 30 seconds but will return (and fail, correctly) as soon as the mismatched exit code arrives!

> Note: `Eventually` and `Consistently` only exercise the `MatchMayChangeInTheFuture` method *if* they are passed a bare value.  If they are passed functions to be polled it is not possible to guarantee that the return value of the function will not change between polling intervals.  In this case, `MatchMayChangeInTheFuture` is not called and the polling continues until either a match is found or the timeout elapses.

### Contributing to Gomega

Contributions are more than welcome.  Either [open an issue](http://github.com/onsi/gomega/issues) for a matcher you'd like to see or, better yet, test drive the matcher and [send a pull request](https://github.com/onsi/gomega/pulls).

When adding a new matcher please mimic the style use in Gomega's current matchers: you should use the `format` package to format your output, put the matcher and its tests in the `matchers` package, and the constructor in the `matchers.go` file in the top-level package.

## `ghttp`: Testing HTTP Clients
The `ghttp` package provides support for testing http *clients*.  The typical pattern in Go for testing http clients entails spinning up an `httptest.Server` using the `net/http/httptest` package and attaching test-specific handlers that perform assertions.

`ghttp` provides `ghttp.Server` - a wrapper around `httptest.Server` that allows you to easily build up a stack of test handlers.  These handlers make assertions against the incoming request and return a pre-fabricated response.  `ghttp` provides a number of prebuilt handlers that cover the most common assertions.  You can combine these handlers to build out full-fledged assertions that test multiple aspects of the incoming requests.

The goal of this documentation is to provide you with an adequate mental model to use `ghttp` correctly.  For a full reference of all the available handlers and the various methods on `ghttp.Server` look at the [godoc](https://godoc.org/github.com/onsi/gomega/ghttp) documentation.

### Making assertions against an incoming request

Let's start with a simple example.  Say you are building an API client that provides a `FetchSprockets(category string)` method that makes an http request to a remote server to fetch sprockets of a given category.

For now, let's not worry about the values returned by `FetchSprockets` but simply assert that the correct request was made.  Here's the setup for our `ghttp`-based Ginkgo test:

```go
Describe("The sprockets client", func() {
    var server *ghttp.Server
    var client *sprockets.Client

    BeforeEach(func() {
        server = ghttp.NewServer()
        client = sprockets.NewClient(server.URL())
    })

    AfterEach(func() {
        //shut down the server between tests
        server.Close()
    })
})
```

Note that the server's URL is auto-generated and varies between test runs.  Because of this, you must always inject the server URL into your client.  Let's add a simple test that asserts that `FetchSprockets` hits the correct endpoint with the correct HTTP verb:

```go
Describe("The sprockets client", func() {
    //...see above

    Describe("fetching sprockets", func() {
        BeforeEach(func() {
            server.AppendHandlers(
                ghttp.VerifyRequest("GET", "/sprockets"),
            )
        })

        It("should make a request to fetch sprockets", func() {
            client.FetchSprockets("")
            Ω(server.ReceivedRequests()).Should(HaveLen(1))
        })
    })
})
```

Here we append a `VerifyRequest` handler to the `server` and call `client.FetchSprockets`.  This call (assuming it's a blocking call) will make a round-trip to the test `server` before returning.  The test `server` receives the request and passes it through the `VerifyRequest` handler which will validate that the request is a `GET` request hitting the `/sprockets` endpoint.  If it's not, the test will fail.

Note that the test can pass trivially if `client.FetchSprockets()` doesn't actually make a request.  To guard against this you can assert that the `server` has actually received a request.  All the requests received by the server are saved off and made available via `server.ReceivedRequests()`.  We use this to assert that there should have been exactly one received requests.

> Guarding against the trivial "false positive"  case outlined above isn't really necessary.  Just good practice when test *driving*.

Let's add some more to our example.  Let's make an assertion that `FetchSprockets` can request sprockets filtered by a particular category:

```go
Describe("The sprockets client", func() {
    //...see above

    Describe("fetching sprockets", func() {
        BeforeEach(func() {
            server.AppendHandlers(
                ghttp.VerifyRequest("GET", "/sprockets", "category=encabulators"),
            )
        })

        It("should make a request to fetch sprockets", func() {
            client.FetchSprockets("encabulators")
            Ω(server.ReceivedRequests()).Should(HaveLen(1))
        })
    })
})
```

`ghttp.VerifyRequest` takes an optional third parameter that is matched against the request `URL`'s `RawQuery`.

Let's extend the example some more.  In addition to asserting that the request is a `GET` request to the correct endpoint with the correct query params, let's also assert that it includes the correct `BasicAuth` information and a correct custom header.  Here's the complete example:

```go
Describe("The sprockets client", func() {
    var (
        server *ghttp.Server
        client *sprockets.Client
        username, password string
    )

    BeforeEach(func() {
        username, password = "gopher", "tacoshell"
        server = ghttp.NewServer()
        client = sprockets.NewClient(server.URL(), username, password)
    })

    AfterEach(func() {
        server.Close()
    })

    Describe("fetching sprockets", func() {
        BeforeEach(func() {
            server.AppendHandlers(
                ghttp.CombineHandlers(
                    ghttp.VerifyRequest("GET", "/sprockets", "category=encabulators"),
                    ghttp.VerifyBasicAuth(username, password),
                    ghttp.VerifyHeader(http.Header{
                        "X-Sprocket-API-Version": []string{"1.0"},
                    }),
                )
            )
        })

        It("should make a request to fetch sprockets", func() {
            client.FetchSprockets("encabulators")
            Ω(server.ReceivedRequests()).Should(HaveLen(1))
        })
    })
})
```

This example *combines* multiple `ghttp` verify handlers using `ghttp.CombineHandlers`.  Under the hood, this returns a new handler that wraps and invokes the three passed in verify handlers.  The request sent by the client will pass through each of these verify handlers and must pass them all for the test to pass.

Note that you can easily add your own verify handler into the mix.  Just pass in a regular `http.HandlerFunc` and make assertions against the received request.

> It's important to understand that you must pass `AppendHandlers` **one** handler *per* incoming request (see [below](#handling-multiple-requests)).  In order to apply multiple handlers to a single request we must first combine them with `ghttp.CombineHandlers` and then pass that *one* wrapper handler in to `AppendHandlers`.

### Providing responses

So far, we've only made assertions about the outgoing request.  Clients are also responsible for parsing responses and returning valid data.  Let's say that `FetchSprockets()` returns two things: a slice `[]Sprocket` and an `error`.  Here's what a happy path test that asserts the correct data is returned might look like:

```go
Describe("The sprockets client", func() {
    //...
    Describe("fetching sprockets", func() {
        BeforeEach(func() {
            server.AppendHandlers(
                ghttp.CombineHandlers(
                    ghttp.VerifyRequest("GET", "/sprockets", "category=encabulators"),
                    ghttp.VerifyBasicAuth(username, password),
                    ghttp.VerifyHeader(http.Header{
                        "X-Sprocket-API-Version": []string{"1.0"},
                    }),
                    ghttp.RespondWith(http.StatusOK, `[
                        {"name": "entropic decoupler", "color": "red"},
                        {"name": "defragmenting ramjet", "color": "yellow"}
                    ]`),
                )
            )
        })

        It("should make a request to fetch sprockets", func() {
            sprockets, err := client.FetchSprockets("encabulators")
            Ω(err).ShouldNot(HaveOccurred())
            Ω(sprockets).Should(Equal([]Sprocket{
                sprockets.Sprocket{Name: "entropic decoupler", Color: "red"},
                sprockets.Sprocket{Name: "defragmenting ramjet", Color: "yellow"},
            }))
        })
    })
})
```

We use `ghttp.RespondWith` to specify the response return by the server.  In this case we're passing back a status code of `200` (`http.StatusOK`) and a pile of JSON.  We then assert, in the test, that the client succeeds and returns the correct set of sprockets.

The fact that details of the JSON encoding are bleeding into this test is somewhat unfortunate, and there's a lot of repetition going on.  `ghttp` provides a `RespondWithJSONEncoded` handler that accepts an arbitrary object and JSON encodes it for you.  Here's a cleaner test:

```go
Describe("The sprockets client", func() {
    //...
    Describe("fetching sprockets", func() {
        var returnedSprockets []Sprocket
        BeforeEach(func() {
            returnedSprockets = []Sprocket{
                sprockets.Sprocket{Name: "entropic decoupler", Color: "red"},
                sprockets.Sprocket{Name: "defragmenting ramjet", Color: "yellow"},
            }

            server.AppendHandlers(
                ghttp.CombineHandlers(
                    ghttp.VerifyRequest("GET", "/sprockets", "category=encabulators"),
                    ghttp.VerifyBasicAuth(username, password),
                    ghttp.VerifyHeader(http.Header{
                        "X-Sprocket-API-Version": []string{"1.0"},
                    }),
                    ghttp.RespondWithJSONEncoded(http.StatusOK, returnedSprockets),
                )
            )
        })

        It("should make a request to fetch sprockets", func() {
            sprockets, err := client.FetchSprockets("encabulators")
            Ω(err).ShouldNot(HaveOccurred())
            Ω(sprockets).Should(Equal(returnedSprockets))
        })
    })
})
```

### Testing different response scenarios

Our test currently only handles the happy path where the server returns a `200`.  We should also test a handful of sad paths.  In particular, we'd like to return a `SprocketsErrorNotFound` error when the server `404`s and a `SprocketsErrorUnauthorized` error when the server returns a `401`.  But how to do this without redefining our server handler three times?

`ghttp` provides `RespondWithPtr` and `RespondWithJSONEncodedPtr` for just this use case.  Both take *pointers* to status codes and respond bodies (objects for the `JSON` case).  Here's the more complete test:

```go
Describe("The sprockets client", func() {
    //...
    Describe("fetching sprockets", func() {
        var returnedSprockets []Sprocket
        var statusCode int

        BeforeEach(func() {
            returnedSprockets = []Sprocket{
                sprockets.Sprocket{Name: "entropic decoupler", Color: "red"},
                sprockets.Sprocket{Name: "defragmenting ramjet", Color: "yellow"},
            }

            server.AppendHandlers(
                ghttp.CombineHandlers(
                    ghttp.VerifyRequest("GET", "/sprockets", "category=encabulators"),
                    ghttp.VerifyBasicAuth(username, password),
                    ghttp.VerifyHeader(http.Header{
                        "X-Sprocket-API-Version": []string{"1.0"},
                    }),
                    ghttp.RespondWithJSONEncodedPtr(&statusCode, &returnedSprockets),
                )
            )
        })

        Context("when the request succeeds", func() {
            BeforeEach(func() {
                statusCode = http.StatusOK
            })

            It("should return the fetched sprockets without erroring", func() {
                sprockets, err := client.FetchSprockets("encabulators")
                Ω(err).ShouldNot(HaveOccurred())
                Ω(sprockets).Should(Equal(returnedSprockets))
            })
        })

        Context("when the response is unauthorized", func() {
            BeforeEach(func() {
                statusCode = http.StatusUnauthorized
            })

            It("should return the SprocketsErrorUnauthorized error", func() {
                sprockets, err := client.FetchSprockets("encabulators")
                Ω(sprockets).Should(BeEmpty())
                Ω(err).Should(MatchError(SprocketsErrorUnauthorized))
            })
        })

        Context("when the response is not found", func() {
            BeforeEach(func() {
                statusCode = http.StatusNotFound
            })

            It("should return the SprocketsErrorNotFound error", func() {
                sprockets, err := client.FetchSprockets("encabulators")
                Ω(sprockets).Should(BeEmpty())
                Ω(err).Should(MatchError(SprocketsErrorNotFound))
            })
        })
    })
})
```

In this way, the status code and returned value (not shown here) can be changed in sub-contexts without having to modify the original test setup.

### Handling multiple requests

So far, we've only seen examples where one request is made per test.  `ghttp` supports handling *multiple* requests too.  `server.AppendHandlers` can be passed multiple handlers and these handlers are evaluated in order as requests come in.

This can be helpful in cases where it is not possible (or desirable) to have calls to the client under test only generate *one* request.  A common example is pagination.  If the sprockets API is paginated it may be desirable for `FetchSprockets` to provide a simpler interface that simply fetches all available sprockets.

Here's what a test might look like:

```go
Describe("fetching sprockets from a paginated endpoint", func() {
    var returnedSprockets []Sprocket
    var firstResponse, secondResponse PaginatedResponse

    BeforeEach(func() {
        returnedSprockets = []Sprocket{
            sprockets.Sprocket{Name: "entropic decoupler", Color: "red"},
            sprockets.Sprocket{Name: "defragmenting ramjet", Color: "yellow"},
            sprockets.Sprocket{Name: "parametric demuxer", Color: "blue"},
        }

        firstResponse = sprockets.PaginatedResponse{
            Sprockets: returnedSprockets[0:2], //first batch
            PaginationToken: "get-second-batch", //some opaque non-empty token
        }

        secondResponse = sprockets.PaginatedResponse{
            Sprockets: returnedSprockets[2:], //second batch
            PaginationToken: "", //signifies the last batch
        }

        server.AppendHandlers(
            ghttp.CombineHandlers(
                ghttp.VerifyRequest("GET", "/sprockets", "category=encabulators"),
                ghttp.RespondWithJSONEncoded(http.StatusOK, firstResponse),
            ),
            ghttp.CombineHandlers(
                ghttp.VerifyRequest("GET", "/sprockets", "category=encabulators&pagination-token=get-second-batch"),
                ghttp.RespondWithJSONEncoded(http.StatusOK, secondResponse),
            ),
        )
    })

    It("should fetch all the sprockets", func() {
        sprockets, err := client.FetchSprockets("encabulators")
        Ω(err).ShouldNot(HaveOccurred())
        Ω(sprockets).Should(Equal(returnedSprockets))
    })
})
```

By default the `ghttp` server fails the test if the number of requests received exceeds the number of handlers registered, so this test ensures that the `client` stops sending requests after receiving the second (and final) set of paginated data.

### MUXing Routes to Handlers

`AppendHandlers` allows you to make ordered assertions about incoming requests.  This places a strong constraint on all incoming requests: namely that exactly the right requests have to arrive in exactly the right order and that no additional requests are allowed.

One can take a different testing strategy, however.  Instead of asserting that requests come in in a predefined order, you may which to build a test server that can handle arbitrarily many requests to a set of predefined routes.  In fact, there may be some circumstances where you want to make ordered assertions on *some* requests (via `AppendHandlers`) but still support sending particular responses to *other* requests that may interleave the ordered assertions.

`ghttp` supports these sorts of usecases via `server.RouteToHandler(method, path, handler)`.

Let's cook up an example.  Perhaps, instead of authenticating via basic auth our sprockets client logs in and fetches a token from the server when performing requests that require authentication.  We could pepper our `AppendHandlers` calls with a handler that handles these requests (this is not a terrible idea, of course!) *or* we could set up a single route at the top of our tests.

Here's what such a test might look like:

```go
Describe("CRUDing sprockes", func() {
    BeforeEach(func() {
        server.RouteToHandler("POST", "/login", ghttp.CombineHandlers(
            ghttp.VerifyRequest("POST", "/login", "user=bob&password=password"),
            ghttp.RespondWith(http.StatusOK, "your-auth-token"),
        ))
    })
    Context("GETting sprockets", func() {
        var returnedSprockets []Sprocket

        BeforeEach(func() {
            returnedSprockets = []Sprocket{
                sprockets.Sprocket{Name: "entropic decoupler", Color: "red"},
                sprockets.Sprocket{Name: "defragmenting ramjet", Color: "yellow"},
                sprockets.Sprocket{Name: "parametric demuxer", Color: "blue"},
            }

            server.AppendHandlers(
                ghttp.CombineHandlers(
                    ghttp.VerifyRequest("GET", "/sprockets", "category=encabulators"),
                    ghttp.RespondWithJSONEncoded(http.StatusOK, returnedSprockets),
                ),
            )
        })

        It("should fetch all the sprockets", func() {
            sprockets, err := client.FetchSprockes("encabulators")
            Ω(err).ShouldNot(HaveOccurred())
            Ω(sprockets).Should(Equal(returnedSprockets))
        })
    })

    Context("POSTing sprockets", func() {
        var sprocketToSave Sprocket
        BeforeEach(func() {
            sprocketToSave = sprockets.Sprocket{Name: "endothermic penambulator", Color: "purple"}

            server.AppendHandlers(
                ghttp.CombineHandlers(
                    ghttp.VerifyRequest("POST", "/sprocket", "token=your-auth-token"),
                    ghttp.VerifyJSONRepresenting(sprocketToSave)
                    ghttp.RespondWithJSONEncoded(http.StatusOK, nil),
                ),
            )
        })

        It("should save the sprocket", func() {
            err := client.SaveSprocket(sprocketToSave)
            Ω(err).ShouldNot(HaveOccurred())
        })
    })
})
```

Here, saving a sprocket triggers authentication, which is handled by the registered `RouteToHandler` handler whereas fetching the list of sprockets does not.

> `RouteToHandler` can take either a string as a route (as seen in this example) or a `regexp.Regexp`.

### Allowing unhandled requests

By default, `ghttp`'s server marks the test as failed if a request is made for which there is no registered handler.

It is sometimes useful to have a fake server that simply returns a fixed status code for all unhandled incoming requests.  `ghttp` supports this: just call `server.SetAllowUnhandledRequests(true)` and `server.SetUnhandledRequestStatusCode(statusCode)`, passing whatever status code you'd like to return.

In addition to returning the registered status code, `ghttp`'s server will also save all received requests.  These can be accessed by calling `server.ReceivedRequests()`.  This is useful for cases where you may want to make assertions against requests *after* they've been made.

To bring it all together: there are three ways to instruct a `ghttp` server to handle requests: you can map routes to handlers using `RouteToHandler`, you can append handlers via `AppendHandlers`, and you can `SetAllowUnhandledRequests` and specify the status code by calling `SetUnhandledRequestStatusCode`.

When a `ghttp` server receives a request it first checks against the set of handlers registered via `RouteToHandler` if there is no such handler it proceeds to pop an `AppendHandlers` handler off the stack, if the stack of ordered handlers is empty, it will check whether `GetAllowUnhandledRequests` returns `true` or `false`.  If `false` the test fails.  If `true`, a response is sent with whatever `GetUnhandledRequestStatusCode` returns.

## `gbytes`: Testing Streaming Buffers

`gbytes` implements `gbytes.Buffer` - an `io.WriteCloser` that captures all input to an in-memory buffer.

When used in concert with the `gbytes.Say` matcher, the `gbytes.Buffer` allows you make *ordered* assertions against streaming data.

What follows is a contrived example.  `gbytes` is best paired with [`gexec`](#gexec-testing-external-processes).

Say you have an integration test that is streaming output from an external API.  You can feed this stream into a `gbytes.Buffer` and make ordered assertions like so:

```go
Describe("attach to the data stream", func() {
    var (
        client *apiclient.Client
        buffer *gbytes.Buffer
    )
    BeforeEach(func() {
        buffer = gbytes.NewBuffer()
        client := apiclient.New()
        go client.AttachToDataStream(buffer)
    })

    It("should stream data", func() {
        Eventually(buffer).Should(gbytes.Say(`Attached to stream as client \d+`))

        client.ReticulateSplines()
        Eventually(buffer).Should(gbytes.Say(`reticulating splines`))
        client.EncabulateRetros(7)
        Eventually(buffer).Should(gbytes.Say(`encabulating 7 retros`))
    })
})
```

These assertions will only pass if the strings passed to `Say` (which are interpreted as regular expressions - make sure to escape characters appropriately!) appear in the buffer.  An opaque read cursor (that you cannot access or modify) is fast-forwarded as successful assertions are made. So, for example:

```go
Eventually(buffer).Should(gbytes.Say(`reticulating splines`))
Consistently(buffer).ShouldNot(gbytes.Say(`reticulating splines`))
```

will (counterintuitively) pass.  This allows you to write tests like:

```go
client.ReticulateSplines()
Eventually(buffer).Should(gbytes.Say(`reticulating splines`))
client.ReticulateSplines()
Eventually(buffer).Should(gbytes.Say(`reticulating splines`))
```

and ensure that the test is correctly asserting that `reticulating splines` appears *twice*.

At any time, you can access the entire contents written to the buffer via `buffer.Contents()`.  This includes *everything* ever written to the buffer regardless of the current position of the read cursor.

### Handling branches

Sometimes (rarely!) you must write a test that must perform different actions depending on the output streamed to the buffer.  This can be accomplished using `buffer.Detect`. Here's a contrived example:

```go
func LoginIfNecessary() {
    client.Authorize()
    select {
    case <-buffer.Detect("You are not logged in"):
        client.Login()
    case <-buffer.Detect("Success"):
        return
    case <-time.After(time.Second):
        ginkgo.Fail("timed out waiting for output")
    }
    buffer.CancelDetects()
}
```

`buffer.Detect` takes a string (interpreted as a regular expression) and returns a channel that will fire *once* if the requested string is detected.  Upon detection, the buffer's opaque read cursor is fast-forwarded so subsequent uses of `gbytes.Say` will pick up from where the succeeding `Detect` left off.  You *must* call `buffer.CancelDetects()` to clean up afterwards (`buffer` spawns one goroutine per call to `Detect`).

### Testing `io.Reader`s, `io.Writer`s, and `io.Closer`s

Implementations of `io.Reader`, `io.Writer`, and `io.Closer` are expected to be blocking.  This makes the following class of tests unsafe:

```go
It("should read something", func() {
    p := make([]byte, 5)
    _, err := reader.Read(p)  //unsafe! this could block forever
    Ω(err).ShouldNot(HaveOccurred())
    Ω(p).Should(Equal([]byte("abcde")))
})
```

It is safer to wrap `io.Reader`s, `io.Writer`s, and `io.Closer`s with explicit timeouts.  You can do this with `gbytes.TimeoutReader`, `gbytes.TimeoutWriter`, and `gbytes.TimeoutCloser` like so:

```go
It("should read something", func() {
    p := make([]byte, 5)
    _, err := gbytes.TimeoutReader(reader, time.Second).Read(p)
    Ω(err).ShouldNot(HaveOccurred())
    Ω(p).Should(Equal([]byte("abcde")))
})
```

The `gbytes` wrappers will return `gbytes.ErrTimeout` if a timeout occurs.

In the case of `io.Reader`s you can leverage the `Say` matcher and the functionality of `gbytes.Buffer` by building a `gbytes.Buffer` that reads from the `io.Reader` asynchronously.  You can do this with the `gbytes` package like so:

```go
It("should read something", func() {
    Eventually(gbytes.BufferReader(reader)).Should(gbytes.Say("abcde"))
})
```

`gbytes.BufferReader` takes an `io.Reader` and returns a `gbytes.Buffer`.  Under the hood an `io.Copy` goroutine is launched to copy data from the `io.Reader` into the `gbytes.Buffer`.  The `gbytes.Buffer` is closed when the `io.Copy` completes.  Because the `io.Copy` is launched asynchronously you *must* make assertions against the reader using `Eventually`.


## `gexec`: Testing External Processes

`gexec` simplifies testing external processes.  It can help you [compile go binaries](#compiling-external-binaries), [start external processes](#starting-external-processes), [send signals and wait for them to exit](#sending-signals-and-waiting-for-the-process-to-exit), make [assertions against the exit code](#asserting-against-exit-code), and stream output into `gbytes.Buffer`s to allow you [make assertions against output](#making-assertions-against-the-process-output).

### Compiling external binaries

You use `gexec.Build()` to compile Go binaries.  These are built using `go build` and are stored off in a temporary directory.  You'll want to `gexec.CleanupBuildArtifacts()` when you're done with the test.

A common pattern is to compile binaries once at the beginning of the test using `BeforeSuite` and to clean up once at the end of the test using `AfterSuite`:

```go
var pathToSprocketCLI string

BeforeSuite(func() {
    var err error
    pathToSprocketCLI, err = gexec.Build("github.com/spacely/sprockets")
    Ω(err).ShouldNot(HaveOccurred())
})

AfterSuite(func() {
    gexec.CleanupBuildArtifacts()
})
```

> By default, `gexec.Build` uses the GOPATH specified in your environment.  You can also use `gexec.BuildIn(gopath string, packagePath string)` to specify a custom GOPATH for the build command.  This is useful to, for example, build a binary against its vendored Godeps.

> You can specify arbitrary environment variables for the build command – such as GOOS and GOARCH for building on other platforms – using `gexec.BuildWithEnvironment(packagePath string, envs []string)`.

### Starting external processes

`gexec` provides a `Session` that wraps `exec.Cmd`.  `Session` includes a number of features that will be explored in the next few sections.  You create a `Session` by instructing `gexec` to start a command:

```go
command := exec.Command(pathToSprocketCLI, "-api=127.0.0.1:8899")
session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
Ω(err).ShouldNot(HaveOccurred())
```

`gexec.Start` calls `command.Start` for you and forwards the command's `stdout` and `stderr` to `io.Writer`s that you provide. In the code above, we pass in Ginkgo's `GinkgoWriter`.  This makes working with external processes quite convenient: when a test passes no output is printed to screen, however if a test fails then any output generated by the command will be provided.

> If you want to see all your output regardless of test status, just run `ginkgo` in verbose mode (`-v`) - now everything written to `GinkgoWriter` makes it onto the screen.

### Sending signals and waiting for the process to exit

`gexec.Session` makes it easy to send signals to your started command:

```go
session.Kill() //sends SIGKILL
session.Interrupt() //sends SIGINT
session.Terminate() //sends SIGTERM
session.Signal(signal) //sends the passed in os.Signal signal
```

If the process has already exited these signal calls are no-ops.

In addition to starting the wrapped command, `gexec.Session` also *monitors* the command until it exits.  You can ask `gexec.Session` to `Wait` until the process exits:

```go
session.Wait()
```

this will block until the session exits and will *fail* if it does not exit within the default `Eventually` timeout.  You can override this timeout by specifying a custom one:

```go
session.Wait(5 * time.Second)
```

> Though you can access the wrapped command using `session.Command` you should not attempt to `Wait` on it yourself.  `gexec` has already called `Wait` in order to monitor your process for you.

> Under the hood `session.Wait` simply uses `Eventually`.


Since the signalling methods return the session you can chain calls together:

```go
session.Terminate().Wait()
```

will send `SIGTERM` and then wait for the process to exit.

### Asserting against exit code

Once a session has exited you can fetch its exit code with `session.ExitCode()`.  You can subsequently make assertions against the exit code.

A more idiomatic way to assert that a command has exited is to use the `gexec.Exit()` matcher:

```go
Eventually(session).Should(Exit())
```

Will verify that the `session` exits within `Eventually`'s default timeout.  You can assert that the process exits with a specified exit code too:

```go
Eventually(session).Should(Exit(0))
```

> If the process has not exited yet, `session.ExitCode()` returns `-1`

### Making assertions against the process output

In addition to streaming output to the passed in `io.Writer`s (the `GinkgoWriter` in our example above), `gexec.Start` attaches `gbytes.Buffer`s to the command's output streams.  These are available on the `session` object via:

```go
session.Out //a gbytes.Buffer connected to the command's stdout
session.Err //a gbytes.Buffer connected to the command's stderr
```

This allows you to make assertions against the stream of output:

```go
Eventually(session.Out).Should(gbytes.Say("hello [A-Za-z], nice to meet you"))
Eventually(session.Err).Should(gbytes.Say("oops!"))
```

Since `gexec.Session` is a `gbytes.BufferProvider` that provides the `Out` buffer you can write assertions against `stdout` output like so:

```go
Eventually(session).Should(gbytes.Say("hello [A-Za-z], nice to meet you"))
```

Using the `Say` matcher is convenient when making *ordered* assertions against a stream of data generated by a live process.  Sometimes, however, all you need is to
wait for the process to exit and then make assertions against the entire contents of its output.  Since `Wait()` returns `session` you can wait for the process to exit, then grab all its stdout as a `[]byte` buffer with a simple one-liner:

```go
Ω(session.Wait().Out.Contents()).Should(ContainSubstring("finished successfully"))
```

### Signaling all processes
`gexec` provides methods to track and send signals to all processes that it starts.

```go
gexec.Kill() //sends SIGKILL to all processes
gexec.Terminate() //sends SIGTERM to all processes
gexec.Signal(int)  //sends the passed in os.Signal signal to all the processes
gexec.Interrupt() //sends SIGINT to all processes
```

If the any of the processes have already exited these signal calls are no-ops.

`gexec` also provides methods to cleanup and wait for all the processes it started.

```go
gexec.KillAndWait()
gexec.TerminateAndWait()
```

You can specify a custom timeout by:

```go
gexec.KillAndWait(5 * time.Second)
gexec.TerminateAndWait(2 * time.Second)
```

The timeout is applied for each of the processes.

It is considered good practice to ensure all of your processes have been killed before the end of the test suite. If you are using `ginkgo` you can use:

```go
AfterSuite(func(){
    gexec.KillAndWait()
})
```

Due to the global nature of these methods, keep in mind that signaling processes will affect all processes started by `gexec`, in any context. For example if these methods where used in an `AfterEach`, then processes started in `BeforeSuite` would also be signaled.

## `gstruct`: Testing Complex Data Types

`gstruct` simplifies testing large and nested structs and slices. It is used for building up complex matchers that apply different tests to each field or element.

### Testing type `struct`

`gstruct` provides the `FieldsMatcher` through the `MatchAllFields` and `MatchFields` functions for applying a separate matcher to each field of a struct:

```go
actual := struct{
    A int
    B bool
    C string
}{5, true, "foo"}
Expect(actual).To(MatchAllFields(Fields{
    "A": BeNumerically("<", 10),
    "B": BeTrue(),
    "C": Equal("foo"),
}))
```

`MatchAllFields` requires that every field is matched, and each matcher is mapped to a field. To match a subset or superset of a struct, you should use the `MatchFields` function with the `IgnoreExtras` and `IgnoreMissing` options. `IgnoreExtras` will ignore fields that don't map to a matcher, e.g.

```go
Expect(actual).To(MatchFields(IgnoreExtras, Fields{
    "A": BeNumerically("<", 10),
    "B": BeTrue(),
    // Ignore lack of "C" in the matcher.
}))
```

`IgnoreMissing` will ignore matchers that don't map to a field, e.g.

```go
Expect(actual).To(MatchFields(IgnoreMissing, Fields{
    "A": BeNumerically("<", 10),
    "B": BeTrue(),
    "C": Equal("foo"),
    "D": Equal("bar"), // Ignored, since actual.D does not exist.
}))
```

The options can be combined with the binary or: `IgnoreMissing|IgnoreExtras`.

### Testing type slice

`gstruct` provides the `ElementsMatcher` through the `MatchAllElements` and `MatchElements` function for applying a separate matcher to each element, identified by an `Identifier` function:

```go
actual := []string{
    "A: foo bar baz",
    "B: once upon a time",
    "C: the end",
}
id := func(element interface{}) string {
    return string(element.(string)[0])
}
Expect(actual).To(MatchAllElements(id, Elements{
    "A": Not(BeZero()),
    "B": MatchRegexp("[A-Z]: [a-z ]+"),
    "C": ContainSubstring("end"),
}))
```

`MatchAllElements` requires that there is a 1:1 mapping from every element to every matcher. To match a subset or superset of elements, you should use the `MatchElements` function with the `IgnoreExtras` and `IgnoreMissing` options. `IgnoreExtras` will ignore elements that don't map to a matcher, e.g.

```go
Expect(actual).To(MatchElements(id, IgnoreExtras, Elements{
    "A": Not(BeZero()),
    "B": MatchRegexp("[A-Z]: [a-z ]+"),
    // Ignore lack of "C" in the matcher.
}))
```

`IgnoreMissing` will ignore matchers that don't map to an element, e.g.

```go
Expect(actual).To(MatchElements(id, IgnoreMissing, Elements{
    "A": Not(BeZero()),
    "B": MatchRegexp("[A-Z]: [a-z ]+"),
    "C": ContainSubstring("end"),
    "D": Equal("bar"), // Ignored, since actual.D does not exist.
}))
```

You can also use the flag `AllowDuplicates` to permit multiple elements in your slice to map to a single key and matcher in your fields (this flag is not meaningful when applied to structs).

```go
everyElementID := func(element interface{}) string {
    return "a constant" // Every element will map to the same key in this case; you can group them into multiple keys, however.
}
Expect(actual).To(MatchElements(everyElementID, AllowDuplicates, Elements{
    "a constant": ContainSubstring(": "), // Because every element passes this test
}))
Expect(actual).NotTo(MatchElements(everyElementID, AllowDuplicates, Elements{
    "a constant": ContainSubstring("foo bar baz"), // Only the first element passes this test
}))
```

The options can be combined with the binary or: `IgnoreMissing|IgnoreExtras|AllowDuplicates`.

Additionally, `gstruct` provides `MatchAllElementsWithIndex` and `MatchElementsWithIndex` function for applying a matcher with index to each element, identified by an `IdentifierWithIndex` function. A helper function is also included with `gstruct` called `IndexIdentity` that provides the functionality of the just using the index as your identifier as seen below.

```go
actual := []string{
    "A: foo bar baz",
    "B: once upon a time",
    "C: the end",
}
id := func(index int, _ interface{}) string {
    return strconv.Itoa(index)
}
Expect(actual).To(MatchAllElementsWithIndex(id, Elements{
    "0": Not(BeZero()),
    "1": MatchRegexp("[A-Z]: [a-z ]+"),
    "2": ContainSubstring("end"),
}))
// IndexIdentity is a helper function equivalent to id in this example
Expect(actual).To(MatchAllElementsWithIndex(IndexIdentity, Elements{
    "0": Not(BeZero()),
    "1": MatchRegexp("[A-Z]: [a-z ]+"),
    "2": ContainSubstring("end"),
}))
```

 The `WithIndex` variants take the same options as the other functions.

### Testing type `map`

All of the `*Fields` functions and types have a corresponding definitions `*Keys` which can perform analogous tests against map types:

```go
actual := map[string]string{
    "A": "correct",
    "B": "incorrect",
}

// fails, because `actual` includes the key B
Expect(actual).To(MatchAllKeys(Keys{
    "A": Equal("correct"),
}))

// passes
Expect(actual).To(MatchAllKeys(Keys{
    "A": Equal("correct"),
    "B": Equal("incorrect"),
}))

// passes
Expect(actual).To(MatchKeys(IgnoreMissing, Keys{
    "A": Equal("correct"),
    "B": Equal("incorrect"),
    "C": Equal("whatever"), // ignored, because `actual` doesn't have this key
}))
```

### Testing pointer values

`gstruct` provides the `PointTo` function to apply a matcher to the value pointed-to. It will fail if the pointer value is `nil`:

    foo := 5
    Expect(&foo).To(PointTo(Equal(5)))
    var bar *int
    Expect(bar).NotTo(PointTo(BeNil()))

### Putting it all together: testing complex structures

The `gstruct` matchers are intended to be composable, and can be combined to apply fuzzy-matching to large and deeply nested structures. The additional `Ignore()` and `Reject()` matchers are provided for ignoring (always succeed) fields and elements, or rejecting (always fail) fields and elements.

Example:

```go
coreID := func(element interface{}) string {
    return strconv.Itoa(element.(CoreStats).Index)
}
Expect(actual).To(MatchAllFields(Fields{
  "Name":      Ignore(),
  "StartTime": BeTemporally(">=", time.Now().Add(-100 * time.Hour)),
  "CPU": PointTo(MatchAllFields(Fields{
        "Time":                 BeTemporally(">=", time.Now().Add(-time.Hour)),
        "UsageNanoCores":       BeNumerically("~", 1E9, 1E8),
        "UsageCoreNanoSeconds": BeNumerically(">", 1E6),
        "Cores": MatchElements(coreID, IgnoreExtras, Elements{
            "0": MatchAllFields(Fields{
                Index: Ignore(),
              "UsageNanoCores":       BeNumerically("<", 1E9),
              "UsageCoreNanoSeconds": BeNumerically(">", 1E5),
            }),
            "1": MatchAllFields(Fields{
                Index: Ignore(),
              "UsageNanoCores":       BeNumerically("<", 1E9),
              "UsageCoreNanoSeconds": BeNumerically(">", 1E5),
            }),
        }),
    })),
    "Memory": PointTo(MatchAllFields(Fields{
  "Time": BeTemporally(">=", time.Now().Add(-time.Hour)),
  "AvailableBytes":  BeZero(),
  "UsageBytes":      BeNumerically(">", 5E6),
  "WorkingSetBytes": BeNumerically(">", 5E6),
  "RSSBytes":        BeNumerically("<", 1E9),
  "PageFaults":      BeNumerically("~", 1000, 100),
  "MajorPageFaults": BeNumerically("~", 100, 50),
    })),
    "Rootfs":             m.Ignore(),
    "Logs":               m.Ignore(),
}))
```

## `gmeasure`: Benchmarking Code

`gmeasure` provides support for measuring and recording benchmarks of your code and tests.  It can be used as a simple standalone benchmarking framework, or as part of your code's test suite.  `gmeasure` integrates cleanly with Ginkgo V2 to enable rich benchmarking of code alognside your tests.

### A Mental Model for `gmeasure`

`gmeasure` is organized around the metaphor of `Experiment`s that can each record multiple `Measurement`s.  To use `gmeasure` you create a `NewExperiment` and use the resulting `experiment` object to record values and durations.  You can then print out the `experiment` to get a report of all measurements or access specific measurements and their statistical aggregates to perform comparisons and/or make assertions.

An `experiment` can record _multiple_ `Measurement`s.  Each `Measurement` has a `Name`, a `Type` (either `MeasurementTypeDuration` or `MeasurementTypeValue`), and a collection of recorded data points (of type `float64` for Value measurements and `time.Duration` for Duration measurements).  In this way an experiment might describe a system or context being measured and can contain multiple measurements - one for each aspect of the system in question.

`Experiment`s can either record values and durations that the user passes in directly.  Or they can invoke callbacks and accept their return values as Value data points, or measure their runtimes to compute Duration data points.  `Experiment`s can also _sample_ callbacks, calling them repeatedly to get an ensemble of data points.

A `Measurement` is created when its first data point is recorded by an `Experiment`.  Subsequent data points with the same measurement name are appended to the measurement:

```go
experiment := gmeasure.NewExperiment("My Experiment")
experiment.RecordDuration("runtime", 3*time.Second) //creates a new Measurement called "runtime"
experiment.RecordDuration("runtime", 5*time.Second) //appends a data point to "runtime"
```

As described below, Measurements can be decorated with additional information.  This includes information about the `Units` for the measurement, the `Precision` with which to render the measurement, and any `Style` to apply when rendering the measurement.  Individual data points can also be decorated with an `Annotation` - an arbitrary string that is associated with that data point and gives it context.  Decorations are applied as typed variadic arguments:

```go
experiment := gmeasure.NewExperiment("My Experiment")

// The first call to `RecordValue` for a measurement must set up any units, style, or precision decorations
experiment.RecordValue("length", 3.141, gmeasure.Units("inches"), gmeasure.Style("{{blue}}"), gmeasure.Precision(2), gmeasure.Annotation("box A)"))

// Subsequent calls can attach an annotation.  In this call a new data-point of `2.71` is added to the `length` measurement with the annotation `box B`.
experiment.RecordValue("length", 2.71, gmeasure.Annotation("box B"))
```

Once recorded, `Measurements` can be fetched from the `experiment` by name via `experiment.Get("name")`.  The returned `Measurement` object includes all the data points.  To get a statistical summary of the data points (that includes the min, max, median, mean, and standard deviation) call `measurement.Stats()` or `experiment.GetStats("name")`.  These statistical summaries can also be rank-ordered with `RankStats()`.

`gmeasure` is designed to integrate with Ginkgo.  This is done by registering `Experiment`s, `Measurement`s and `Ranking`s as `ReportEntry`s via Ginkgo's `AddReportEntry`.  This will cause Ginkgo to emit nicely formatted and styled summaries of each of these objects when generating the test report.

Finally, `gmeasure` provides a mechanism to cache `Experiment`s to disk with a specified version number.  This enables multiple use-cases.  You can cache expensive experiments to avoid rerunning them while you iterate on other experiments.  You can also compare experiments to cached experiments to explore whether changes in performance have been introduced to the codebase under test.

`gmeasure` includes detailed [godoc documentation](https://pkg.go.dev/github.com/onsi/gomega/gmeasure) - this narrative documentation is intended to help you get started with `gmeasure`.

### Measuring Values

`Experiment`s can record arbitrary `float64` values.  You can do this by directly providing a `float64` via `experiment.RecordValue(measurementName string, value float64, decorators ...interface{})` or by providing a callback that returns a float64 via `experiment.MeasureValue(measurementName string, callback func() float64, decorators ...interface{})`.

You can apply `Units`, `Style`, and `Precision` decorators to control the appearance of the `Measurement` when reports are generated.  These decorators must be applied when the first data point is recorded but can be elided thereafter.  You can also associate an `Annotation` decoration with any recorded data point.

`Experiment`s are thread-safe so you can call `RecordValue` and `MeasureValue` from any goroutine.

### Measuring Durations

`Experiment`s can record arbitrary `time.Duration` durations.  You can do this by directly providing a `time.Duration` via `experiment.RecordDuration(measurementName string, duration time.Duration, decorators ...interface{})` or by providing a callback via `experiment.MeasureDuration(measurementName string, callback func(), decorators ...interface{})`.  `gmeasure` will run the callback and measure how long it takes to complete.

You can apply `Style` and `Precision` decorators to control the appearance of the `Measurement` when reports are generated.  These decorators must be applied when the first data point is recorded but can be elided thereafter.  You can also associate an `Annotation` decoration with any recorded data point.

`Experiment`s are thread-safe so you can call `RecordDuration` and `MeasureDuration` from any goroutine.

### Sampling

`Experiment`s support sampling callback functions repeatedly to build an ensemble of data points.  All the sampling methods are configured by passing in a `SamplingConfig`:

```go
type SamplingConfig struct {
    N int
    Duration time.Duration
    NumParallel int
    MinSamplingInterval time.Duration
}
```

Setting `SamplingConfig.N` limits the total number of samples to perform to `N`.  Setting `SamplingConfig.Duration` limits the total time spent sampling to `Duration`.  At least one of these fields must be set.  If both are set then `gmeasure` will `sample` until the first limiting condition is met.  Setting `SamplingConfig.MinSamplingInterval` causes `gmeasure` to wait until at least `MinSamplingInterval` has elapsed between subsequent samples.

By default, the `Experiment`'s sampling methods will run their callbacks serially within the calling goroutine.  If `NumParallel` greater than `1`, however, the sampling methods will spin up `NumParallel` goroutines and farm the work among them.  You cannot use `NumParallel` with `MinSamplingInterval`.

The basic sampling method is `experiment.Sample(callback func(idx int), samplingConfig SamplingConfig)`.  This will call the callback function repeatedly, passing in an `idx` counter that increments between each call.  The sampling will end based on the conditions provided in `SamplingConfig`.  Note that `experiment.Sample` is not explicitly associated with a measurement.  You can use `experiment.Sample` whenever you want to repeatedly invoke a callback up to a limit of `N` and/or `Duration`.  You can then record arbitrarily many value or duration measurements in the body of the callback.

A common use-case, however, is to invoke a callback repeatedly to measure its duration or record its returned value and thereby generate an ensemble of data-points.  This is supported via the `SampleX` family of methods built on top of `Sample`:

```go
experiment.SampleValue(measurementName string, callback func(idx int) float64, samplingConfig SamplingConfig, decorations ...interface{})
experiment.SampleDuration(measurementName string, callback func(idx int), samplingConfig SamplingConfig, decorations ...interface{})
experiment.SampleAnnotatedValue(measurementName string, callback func(idx int) (float64, Annotation), samplingConfig SamplingConfig, decorations ...interface{})
experiment.SampleAnnotatedDuration(measurementName string, callback func(idx int) Annotation, samplingConfig SamplingConfig, decorations ...interface{})
```

each of these will contribute data points to the `Measurement` with name `measurementName`.  `SampleValue` records the `float64` values returned by its callback.  `SampleDuration` times each invocation of its callback and records the measured duration.  `SampleAnnotatedValue` and `SampleAnnotatedDuration` expect their callbacks to return `Annotation`s.  These are attached to each generated data point.

All these methods take the same decorators as their corresponding `RecordX` methods.

### Measuring Durations with `Stopwatch`

In addition to `RecordDuration` and `MeasureDuration`, `gmeasure` also provides a `Stopwatch`-based abstraction for recording durations.  To motivate `Stopwatch` consider the following example.  Let's say we want to measure the end-to-end performance of a web-server.  Here's the code we'd like to measure:

```go
It("measures the end-to-end performance of the web-server", func() {
    model, err := client.Fetch("model-id-17")
    Expect(err).NotTo(HaveOccurred())

    err = model.ReticulateSpines()
    Expect(err).NotTo(HaveOccurred())

    Expect(client.Save(model)).To(Succeed())

    reticulatedModels, err := client.List("reticulated-models")
    Expect(err).NotTo(HaveOccurred())
    Expect(reticulatedModels).To(ContainElement(model))
})
```

One approach would be to use `MeasureDuration`:

```go
It("measures the end-to-end performance of the web-server", func() {
    experiment := gmeasure.NewExperiment("end-to-end web-server performance")
    AddReportEntry(experiment.Name, experiment)

    var model Model
    var err error
    experiment.MeasureDuration("fetch", func() {
        model, err = client.Fetch("model-id-17")
    })
    Expect(err).NotTo(HaveOccurred())

    err = model.ReticulateSpines()
    Expect(err).NotTo(HaveOccurred())

    experiment.MeasureDuration("save", func() {
        Expect(client.Save(model)).To(Succeed())
    })

    var reticulatedModels []Models
    experiment.MeasureDuration("list", func() {
        reticulatedModels, err = client.List("reticulated-models")
    })
    Expect(err).NotTo(HaveOccurred())
    Expect(reticulatedModels).To(ContainElement(model))
})
```

this... _works_.  But all those closures and local variables make the test a bit harder to read.  We can clean it up with a `Stopwatch`:

```go
It("measures the end-to-end performance of the web-server", func() {
    experiment := gmeasure.NewExperiment("end-to-end web-server performance")
    AddReportEntry(experiment.Name, experiment)

    stopwatch := experiment.NewStopwatch() // start the stopwatch

    model, err := client.Fetch("model-id-17")
    stopwatch.Record("fetch") // record the amount of time elapsed and store it in a Measurement named fetch
    Expect(err).NotTo(HaveOccurred())

    err = model.ReticulateSpines()
    Expect(err).NotTo(HaveOccurred())

    stopwatch.Reset() // reset the stopwatch
    Expect(client.Save(model)).To(Succeed())
    stopwatch.Record("save").Reset() // record the amount of time elapsed since the last Reset and store it in a Measurement named save, then reset the stopwatch 

    reticulatedModels, err := client.List("reticulated-models")
    stopwatch.Record("list")
    Expect(err).NotTo(HaveOccurred())
    Expect(reticulatedModels).To(ContainElement(model))
})
```

that's now much cleaner and easier to reason about.  If we wanted to sample the server's performance concurrently we could now simply wrap the relevant code in an `experiment.Sample`:

```go
It("measures the end-to-end performance of the web-server", func() {
    experiment := gmeasure.NewExperiment("end-to-end web-server performance")
    AddReportEntry(experiment.Name, experiment)

    experiment.Sample(func(idx int) {
        defer GinkgoRecover() // necessary since these will launch as goroutines and contain assertions
        stopwatch := experiment.NewStopwatch() // we make a new stopwatch for each sample.  Experiments are threadsafe, but Stopwatches are not.

        model, err := client.Fetch("model-id-17")
        stopwatch.Record("fetch")
        Expect(err).NotTo(HaveOccurred())

        err = model.ReticulateSpines()
        Expect(err).NotTo(HaveOccurred())

        stopwatch.Reset()
        Expect(client.Save(model)).To(Succeed())
        stopwatch.Record("save").Reset()

        reticulatedModels, err := client.List("reticulated-models")
        stopwatch.Record("list")
        Expect(err).NotTo(HaveOccurred())
        Expect(reticulatedModels).To(ContainElement(model))
    }, gmeasure.SamplingConfig{N:100, Duration:time.Minute, NumParallel:8})
})
```

Check out the [godoc documentation](https://pkg.go.dev/github.com/onsi/gomega/gmeasure#Stopwatch) for more details about `Stopwatch` including support for `Pause`ing and `Resume`ing the stopwatch.

### Stats and Rankings: Comparing Measurements

Once you've recorded a few measurements you'll want to try to understand and interpret them.  `gmeasure` allows you to quickly compute statistics for a given measurement.  Consider the following example.  Let's say we have two different ideas for how to implement a sorting algorithm and want to hone in on the algorithm with the shortest median runtime.  We could run an experiment:

```go
It("identifies the fastest algorithm", func() {
    experiment := gmeasure.NewExperiment("dueling algorithms")
    AddReportEntry(experiment.Name, experiment)

    experiment.SampleDuration("runtime: algorithm 1", func(_ int) {
        RunAlgorithm1()
    }, gmeasure.SamplingConfig{N:1000})

    experiment.SampleDuration("runtime: algorithm 2", func(_ int) {
        RunAlgorithm2()
    }, gmeasure.SamplingConfig{N:1000})
})
```

This will sample the two competing tables and print out a tabular representation of the resulting statistics.  (Note - you don't need to use Ginkgo here, you could just use `gmeasure` in your code directly and then `fmt.Println` the `experiment` to get the tabular report).

We could compare the tables by eye manually - or ask `gmeasure` to pick the winning algorithm for us:

```go
It("identifies the fastest algorithm", func() {
    experiment := gmeasure.NewExperiment("dueling algorithms")
    AddReportEntry(experiment.Name, experiment)

    experiment.SampleDuration("runtime: algorithm 1", func(_ int) {
        RunAlgorithm1()
    }, gmeasure.SamplingConfig{N:1000})

    experiment.SampleDuration("runtime: algorithm 2", func(_ int) {
        RunAlgorithm2()
    }, gmeasure.SamplingConfig{N:1000})

    ranking := gmeasure.RankStats(gmeasure.LowerMedianIsBetter, experiment.GetStats("runtime: algorithm 1"), experiment.GetStats("runtime: algorithm 2"))
    AddReportEntry("Ranking", ranking)
})
```

This will now emit a ranking result that will highlight the winning algorithm (in this case, the algorithm with the lower Median).  `RankStats` supports the following `RankingCriteria`:

- `LowerMeanIsBetter`
- `HigherMeanIsBetter`
- `LowerMedianIsBetter`
- `HigherMedianIsBetter`
- `LowerMinIsBetter`
- `HigherMinIsBetter`
- `LowerMaxIsBetter`
- `HigherMaxIsBetter`

We can also inspect the statistics of the two algorithms programatically.  `experiment.GetStats` returns a `Stats` object that provides access to the following `Stat`s:

- `StatMin` - the data point with the smallest value
- `StatMax` - the data point with the highest values
- `StatMedian` - the median data point
- `StatMean` - the mean of all the data points
- `StatStdDev` - the standard deviation of all the data points

`Stats` can represent either Value Measurements or Duration Measurements.  When inspecting a Value Measurement you can pull out the requested `Stat` (say, `StatMedian`) via `stats.ValueFor(StatMedian)` - this returns a `float64`.  When inspecting Duration Measurements you can fetch `time.Duration` statistics via `stats.DurationFor(StatX)`.  For either type you can fetch an appropriately formatted string representation of the stat via `stats.StringFor(StatX)`.  You can also get a `float64` for either type by calling `stats.FloatFor(StatX)` (this simply returns a `float64(time.Duration)` for Duration Measurements and can be useful when you need to do some math with the stats).

Going back to our dueling algorithms example.  Lets say we find that Algorithm 2 is the winner with a median runtime of around 3 seconds - and we want to be alerted by a failing test should the winner ever change, or the median runtime vary substantially.  We can do that by writing a few assertions:

```go
It("identifies the fastest algorithm", func() {
    experiment := gmeasure.NewExperiment("dueling algorithms")
    AddReportEntry(experiment.Name, experiment)

    experiment.SampleDuration("runtime: algorithm 1", func(_ int) {
        RunAlgorithm1()
    }, gmeasure.SamplingConfig{N:1000})

    experiment.SampleDuration("runtime: algorithm 2", func(_ int) {
        RunAlgorithm2()
    }, gmeasure.SamplingConfig{N:1000})

    ranking := gmeasure.RankStats(gmeasure.LowerMedianIsBetter, experiment.GetStats("runtime: algorithm 1"), experiment.GetStats("runtime: algorithm 2"))
    AddReportEntry("Ranking", ranking)

    //assert that algorithm 2 is the winner
    Expect(ranking.Winner().MeasurementName).To(Equal("runtime: algorithm 2"))

    //assert that algorithm 2's median is within 0.5 seconds of 3 seconds
    Expect(experiment.GetStats("runtime: algorithm 2").DurationFor(gmeasure.StatMedian)).To(BeNumerically("~", 3*time.Second, 500*time.Millisecond))
})
```

### Formatting Experiment and Measurement Output

`gmeasure` can produce formatted tabular output for `Experiment`s, `Measurement`s, and `Ranking`s.  Each of these objects provides a `String()` method and a `ColorableString()` method.  The `String()` method returns a string that does not include any styling tags whereas the `ColorableString()` method returns a string that includes Ginkgo's console styling tags (e.g. Ginkgo will render a string like `{{blue}}{{bold}}hello{{/}} there` as a bold blue "hello" followed by a default-styled " there").  `ColorableString()` is called for you automatically when you register any of these `gmeasure` objects as Ginkgo `ReportEntry`s.

When printing out `Experiment`s, `gmeasure` will produce a table whose columns correspond to the key statistics provided by `gmeasure.Stats` and whose rows are the various `Measurement`s recorded by the `Experiment`.  Users can also record and emit notes - contextual information about the experiment - by calling `experiment.RecordNote(note string)`.  Each note will get its own row in the table.

When printing out `Measurement`s, `gmeasure` will produce a table that includes _all_ the data points and annotations for the `Measurement`.

When printing out `Ranking`s, `gmeasure` will produce a table similar to the `Experiment` table with the `Measurement`s sorted by `RankingCriteria`.

Users can adjust a few aspects of `gmeasure`s output.  This is done by providing decorators to the `Experiment` methods that record data points:

- `Units(string)` - the `Units` decorator allows you to associate a set of units with a measurement.  Subsequent renderings of the measurement's name will include the units in `[]` square brackets.
- `Precision(int or time.Duration)` - the `Precision` decorator controls the rendering of numerical information.  For Value Measurements an `int` is used to express the number of decimal points to print.  For example `Precision(3)` will render values with `fmt.Sprintf("%.3f", value)`.  For Duration Measurements a `time.Duration` is used to round durations before rendering them.  For example `Precision(time.Second)` will render durations via `duration.Round(time.Second).String()`.
- `Style(string)` - the `Style` decorator allows you to associate a Ginkgo console style to a measurement.  The measurement's row will be rendered with this style.  For example `Style("{{green}}")` will emit a green row.

These formatting decorators **must** be applied to the _first_ data point recorded for a given Measurement (this is when the Measurement object is initialized and its style, precision, and units fields are populated).

Just to get concrete here's a fleshed out example that uses all the things:

```go
It("explores a complex object", func() {
    experiment := gmeasure.NewExperiment("exploring the encabulator")
    AddReportEntry(experiment.Name, experiment)

    experiment.RecordNote("Encabulation Properties")
    experiment.Sample(func(idx int) {
        stopwatch := experiment.NewStopwatch()
        encabulator.Encabulate()
        stopwatch.Record("Encabulate Runtime", gmeasure.Style("{{green}}"), gmeasure.Precision(time.Millisecond))

        var m runtime.MemStats
        runtime.ReadMemStats(&m)
        experiment.RecordValue("Encabulate Memory Usage", float64(m.Alloc / 1024 / 1024), gmeasure.Style("{{red}}"), gmeasure.Precision(3), gmeasure.Units("MB"), gmeasure.Annotation(fmt.Sprintf("%d", idx)))
    }, gmeasure.SamplingConfig{N:1000, NumParallel:4})

    experiment.RecordNote("Encabulation Teardown")
    experiment.MeasureDuration("Teardown Runtime", func() {
        encabulator.Teardown()
    }, gmeasure.Style("{{yellow}}"))

    memoryStats := experiment.GetStats("Encabulate Memory Usage")
    minMemory := memoryStats.ValueFor(gmeasure.StatMin)
    maxMemory := memoryStats.ValueFor(gmeasure.StatMax)
    Expect(maxMemory - minMemory).To(BeNumerically("<=", 10), "Should not see memory fluctuations exceeding 10 megabytes")        
})
```

### Ginkgo Integration

The examples throughout this documentation have illustrated how `gmeasure` interoperates with Ginkgo.  In short - you can emit output for `Experiment`, `Measurement`s, and `Ranking`s by registering them as Ginkgo `ReportEntry`s via `AddReportEntry()`.

This simple connection point ensures that the output is appropriately formatted and associated with the spec in question.  It also ensures that Ginkgo's machine readable reports will include appropriately encoded versions of these `gmeasure` objects.  So, for example, `ginkgo --json-report=report.json` will include JSON encoded `Experiment`s in `report.json` if you remember to `AddReportEntry` the `Experiment`s.

### Caching Experiments 

`gmeasure` supports caching experiments to local disk.  Experiments can be stored and retreived from the cache by name and version number.  Caching allows you to skip rerunning expensive experiments and versioned caching allows you to bust the cache by incrementing the version number.  Under the hood, the cache is simply a set of files in a directory.  Each file contains a JSON encoded header with the experiment's name and version number followed by the JSON-encoded experiment.  The various cache methods are documented over at [pkg.go.dev](https://pkg.go.dev/github.com/onsi/gomega/gmeasure#ExperimentCache).

Using an `ExperimentCache` with Ginkgo takes a little bit of wiring.  Here's an example:

```go
const EXPERIMENT_VERSION = 1 //bump this to bust the cache and recompute _all_ experiments

Describe("some experiments", func() {
    var cache gmeasure.ExperimentCache
    var experiment *gmeasure.Experiment

    BeforeEach(func() {
        cache = gmeasure.NewExperimentCache("./gmeasure-cache")
        name := CurrentSpecReport().LeafNodeText // we use the text in each It block as the name of the experiment
        experiment = cache.Load(name, EXPERIMENT_VERSION) // we try to load the experiment from the cache
        if experiment != nil {
            // we have a cache hit - report on the experiment and skip this test.
            AddReportEntry(experiment)
            Skip("cached")
        }
        //we have a cache miss, make a new experiment and proceed with the test.
        experiment = gmeasure.NewExperiment(name)
        AddReportEntry(experiment)
    })

    It("measures foo runtime", func() {
        experiment.SampleDuration("runtime", func() {
            //do stuff
        }, gmeasure.SamplingConfig{N:100})
    })

    It("measures bar runtime", func() {
        experiment.SampleDuration("runtime", func() {
            //do stuff
        }, gmeasure.SamplingConfig{N:100})
    })

    AfterEach(func() {
        // AfterEaches always run, even for tests that call `Skip`.  So we make sure we aren't a skipped test then save the experiment to the cache
        if !CurrentSpecReport().State.Is(types.SpecStateSkipped) {
            cache.Save(experiment.Name, EXPERIMENT_VERSION, experiment)
        }
    })
})
```

this test will load the experiment from the cache if it's available or run the experiment and store it in the cache if it is not.  Incrementing `EXPERIMENT_VERSION` will force all experiments to rerun.

Another usecase for `ExperimentCache` is to cache and commit experiments to source control for use as future baselines.  Your code can assert that measurements are within a certain range of the stored baseline.  For example:

```go
Describe("server performance", func() {
    It("ensures a performance regression has not been introduced", func() {
        // make an experiment
        experiment := gmeasure.NewExperiment("performance regression test")
        AddReportEntry(experiment.Name, experiment)

        // measure the performance of one endpoint
        experiment.SampleDuration("fetching one", func() {
            model, err := client.Get("id-1")
            Expect(err).NotTo(HaveOccurred())
            Expect(model.Id).To(Equal("id-1"))
        }, gmeasure.SamplingConfig{N:100})

        // measure the performance of another endpoint
        experiment.SampleDuration("listing", func() {
            models, err := client.List()
            Expect(err).NotTo(HaveOccurred())
            Expect(models).To(HaveLen(30))
        }, gmeasure.SamplingConfig{N:100})

        cache := gmeasure.NewExperimentCache("./gemasure-cache")
        baseline := cache.Load("performance regression test", 1)
        if baseline == nil {
            // this is the first run, let's store a baseline
            cache.Save("performacne regression test", 1, experiment)
        } else {
            for _, m := range []string{"fetching one", "listing"} {
                baselineStats := baseline.GetStats(m)
                currentStats := experiment.GetStats(m)

                //make sure the mean of the current performance measure is within one standard deviation of the baseline
                Expect(currentStats.DurationFor(gmeasure.StatMean)).To(BeNumerically("~", baselineStats.DurationFor(gmeasure.StatsMean), baselineStats.DurationFor(gmeasure.StatsStdDev)), m)
            }
        }
    })        
})
```

{% endraw  %}
