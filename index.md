---
layout: default
title: Gomega
---

[Gomega](http://github.com/onsi/gomega) is a matcher/assertion library.  It is best paired with the [Ginkgo](http://github.com/onsi/ginkgo) BDD test framework, but can be adapted for use in other contexts too.

---

##Getting Gomega

Just `go get` it:

    $ go get github.com/onsi/gomega

---

##Using Gomega with Ginkgo

When a Gomega assertion fails, Gomega calls an `OmegaFailHandler`.  This is a function that you must provide using `gomega.RegisterFailHandler()`.

If you're using Ginkgo, all you need to do is:

    gomega.RegisterFailHandler(ginkgo.Fail)

before you start your test suite.

If you use the `ginkgo` CLI to `ginkgo bootstrap` a test suite, this hookup will be automatically generated for you.

---

## Using Gomega with Golang's XUnit-style Tests

Though Gomega is tailored to work best with Ginkgo it is easy to use Gomega with Golang's XUnit style tests.  Here's how: me know.

To use Gomega with Golang's XUnit style tests:

    func TestFarmHasCow(t *testing.T) {
        RegisterTestingT(t)

        f := farm.New([]string{"Cow", "Horse"})
        Expect(f.HasCow()).To(BeTrue(), "Farm should have cow")
    }

There are two caveats:

- You **must** register the `t *testing.T` passed to your test with Gomega before you make any assertions associated with that test.  So every `Test...` function in your suite should have the `RegisterTestingT(t)` line.
- Gomega uses a global (singleton) fail handler.  This has the benefit that you don't need to pass the fail handler down to each test, but does mean that *you cannot run your XUnit style tests in parallel with Gomega*.  If you find this odious, open an issue on Github and let

> Gomega tests written with Ginkgo *can* be run in parallel using the `ginkgo` CLI.  This is because Ginkgo runs it's parallel specs in different *processes* whereas the default Golang test runner runs parallel tests in the *same* process.  The latter approach makes your test suite susceptible to test pollution and is avoided by Ginkgo.

---

##Making Assertions

Gomega provides two notations for making assertions.  These notations are functionally equivalent and their differences are purely aesthetic.

- When you use the `Ω` notation, your assertions look like this:

        Ω(ACTUAL).Should(Equal(EXPECTED))
        Ω(ACTUAL).ShouldNot(Equal(EXPECTED))

- When you use the `Expect` notation, your assertions look like this:

        Expect(ACTUAL).To(Equal(EXPECTED))
        Expect(ACTUAL).NotTo(Equal(EXPECTED))
        Expect(ACTUAL).ToNot(Equal(EXPECTED))

On OS X the `Ω` character is easy to type.  Just hit option-z: `⌥z`

On the left hand side, you can pass anything you want in to `Ω` and `Expect` for `ACTUAL`.  On the right hand side you must pass an object that satisfies the `OmegaMatcher` interface.  Gomega's matchers (e.g. `Equal(EXPECTED)`) are simply functions that create and initialize an appropriate `OmegaMatcher` object.

> The `OmegaMatcher` interface is pretty simple and is discussed in the [custom matchers](#adding_your_own_matchers) section.

Each assertion returns a `bool` denoting whether or not the assertion passed.  This is useful for bailing out of a test early if an assertion fails:

        goodToGo := Ω(WeAreSetUp()).Should(BeTrue())
        if !goodToGo {
            return
        }
        doSomethingExpensive()

> With Ginkgo, a failed assertion does *not* bail out of the current test.  It is generally unnecessary to do so, but in cases where bailing out is necessary, use the `bool` return value and pattern outlined above.

### Handling Errors

It is a common pattern, in Golang, for functions and methods to return two things - a value and an error.  For example:

    func DoSomethingHard() (string, error) {
        ...
    }

To assert on the return value of such a method you might write a test that looks like this:

    result, err := DoSomethingHard()
    Ω(err).ShouldNot(HaveOccurred())
    Ω(result).Should(Equal("foo"))

This is a very common use case so Gomega streamlines it for you.  Both `Ω` and `Expect` accept *multiple* arguments.  The first argument is passed to the matcher, and the match only succeeds if *all* subsequent arguments are required to be `nil` or zero-valued.  With this, we can rewrite the above example as:

    Ω(DoSomethingHard()).Should(Equal("foo"))

This will only pass if the return value of `DoSomethingHard()` is `("foo", nil)`.

### Annotating Assertions

You can annotate any assertion by passing a format string (and optional inputs to format) after the `OmegaMatcher`:

    Ω(ACTUAL).Should(Equal(EXPECTED), "My annotation %d", foo)
    Ω(ACTUAL).ShouldNot(Equal(EXPECTED), "My annotation %d", foo)
    Expect(ACTUAL).To(Equal(EXPECTED), "My annotation %d", foo)
    Expect(ACTUAL).NotTo(Equal(EXPECTED), "My annotation %d", foo)
    Expect(ACTUAL).ToNot(Equal(EXPECTED), "My annotation %d", foo)

The format string and inputs will be passed to `fmt.Sprintf(...)`.  If the assertion fails, Gomega will print your annotation alongside its standard failure message.

This is useful in cases where the standard failure message lacks context.  For example, if the following assertion fails:

    Ω(SprocketsAreLeaky()).Should(BeFalse())

Gomega will output:

    Expected
      <bool>: true
    to be false

But this assertion:

    Ω(SprocketsAreLeaky()).Should(BeFalse(), "Sprockets shouldn't leak")

Will offer the more helpful output:

    Sprockets shouldn't leak
    Expected
      <bool>: true
    to be false

---

##Making Asynchronous Assertions

Gomega has support for making *asynchronous* assertions.  There are two functions that provide this support `Eventually` and `Consistently`.

###Eventually

`Eventually` checks that an assertion *eventually* passes.  It does this by polling its argument until the matcher succeeds.

For example:

    Eventually(func() []int {
        return thing.SliceImMonitoring
    }).Should(HaveLen(2))

    Eventually(func() string {
        return thing.Status
    }).ShouldNot(Equal("Stuck Waiting"))

`Eventually` will poll the passed in function (which must have zero-arguments and at least one return value) repeatedly and check the return value against the `OmegaMatcher`.  `Eventually` then blocks until the match succeeds or until a timeout interval has elapsed.

The default value for the timeout is 1 second and the default value for the polling interval is 10 milliseconds.  You can change these values by passing them in just after your function:

    Eventually(func() []int {
        return thing.SliceImMonitoring
    }, TIMEOUT, POLLING_INTERVAL).Should(HaveLen(2))

These can be passed in as `time.Duration`s, string representations of a `time.Duration` (e.g. `"2s"`) or `float64` values (in which case they are interpreted as seconds).

`Eventually` is especially handy when writing integration tests against asynchronous services or components:

    externalProcess.DoSomethingAmazing()
    Eventually(func() bool {
        return somethingAmazingHappened()
    }).Should(BeTrue())

The function that you pass to `Eventually` can have more than one return value.  In that case, `Eventually` passes the first return value to the matcher and asserts that all other return values are `nil` or zero-valued.  This allows you to use `Eventually` with functions that return a value and an error -- a common pattern in Go.  For example, say you have a method on an object named `FetchNameFromNetwork()` that returns a string value and an error.  Given an instance then you could simply write:

    Eventually(myInstance.FetchNameFromNetwork).Should(Equal("archibald"))

If the argument to `Eventually` is *not* a function, `Eventually` will simply run the matcher against the argument.  This works really well with the Gomega matchers geared towards working with channels:

    Eventually(channel).Should(BeClosed())
    Eventually(channel).Should(Receive())

This also pairs well with `gexec`'s `Session` command wrappers nad `gbyte`'s `Buffer`s:

    Eventually(session).Should(gexec.Exit(0)) //the wrapped command should exit with status 0, eventually
    Eventually(buffer).Should(Say("something matching this regexp"))
    Eventually(session.Out).Should(Say("Splines reticulated"))

> Note that `Eventually(slice).Should(HaveLen(N))` probably won't do what you think it should -- eventually will be passed a pointer to the slice, yes, but if the slice is being `append`ed to (as in: `slice := append(slice, ...)`) Go will generate a new pointer and the pointer passed to `Eventually` will not contain the new elements.  In such cases you should always pass `Eventually` a function that, when polled, returns the slice.
> As with synchronous assertions, you can annotate asynchronous assertions by passing a format string and optional inputs after the `OmegaMatcher`.

###Consistently

`Consistently` checks that an assertion passes or a period of time.  It does this by polling its argument for the fixed period of time and fails if the matcher ever fails during that period of time.

For example:

    Consistently(func() []int {
        return thing.MemoryUsage()
    }).Should(BeNumerically("<", 10))

`Consistently` will poll the passed in function (which must have zero-arguments and at least one return value) repeatedly and check the return value against the `OmegaMatcher`.  `Consitently` blocks and only returns when the desired duration has elapsed or if the matcher fails.  The default value for the wait-duration is 100 milliseconds.  The default polling interval is 10 milliseconds.  Like `Eventually`, you can change these values by passing them in just after your function:


    Consistently(func() []int {
        return thing.MemoryUsage()
    }, DURATION, POLLING_INTERVAL).Should(BeNumerically("<", 10))

As with `Eventually`, these can be `time.Duration`s, string representations of a `time.Duration` (e.g. `"200ms"`) or `float64`s that are interpreted as seconds.

`Consistently` tries to capture the notion that something "does not eventually" happen.  A common use-case is to assert that no goroutine writes to a channel for a period of time.  If you pass `Consistently` an argument that is not a function, it simply passes that argument to the matcher.  So we can asser that:

    Consistently(channel).ShouldNot(Receive())

To assert that nothing gets sent to a channel.

As with `Eventually`, if you pass `Consistently` a function that returns more than one value, it will pass the first value to the matcher and assert that all other values are `nil` or zero-valued.

> Developers often try to use `runtime.Gosched()` to nudge background goroutines to run.  This can lead to flaky tests as it is not deterministic that a given goroutine will run during the `Gosched`.  `Consistently` is particularly handy in these cases: it polls for 100ms which is typically more than enough time for all your Goroutines to run.  Yes, this is basically like putting a time.Sleep() in your tests....  Sometimes, when making negative assertions in a concurrent world, that's the best you can do!

###Modifying Default Intervals

By default, `Eventually` will poll every 10 milliseconds for up to 1 second and `Consistently` will monitor every 10 milliseconds for up to 100 milliseconds.  You can modify these defaults across your test suite with:

    SetDefaultEventuallyTimeout(t time.Duration)
    SetDefaultEventuallyPollingInterval(t time.Duration)
    SetDefaultConsistentlyDuration(t time.Duration)
    SetDefaultConsistentlyPollingInterval(t time.Duration)

---

##Making Assertions in Helper Functions

While writing [custom matchers](#adding_your_own_matchers) is an expressive way to make assertions against your code, it is often more convenient to write one-off helper functions like so:

    var _ = Describe("Turboencabulator", func() {
        ...
        assertTurboencabulatorContains(components ...string) {
            components, err := turboEncabulator.GetComponents()
            Expect(err).NotTo(HaveOccurred())

            Expect(components).To(HaveLen(components))
            for _, component := range components {
                Expect(components).To(ContainElement(component))
            }
        }

        It("should have components", func() {
            assertTurboEncabulatorContains("semi-boloid slots", "grammeters")
        })
    })

This makes your tests more expressive and reduces boilerplate.  However, when an assertion in the helper fails the line numbers provided by Gomega are unhelpful.  Instead of pointing you to the line in your test that failed, they point you the line in the helper.

To get around this, Gomega provides versions of `Expect`, `Eventually` and `Consistently` named `ExpectWithOffset`, `EventuallyWithOffset` and `ConsistentlyWithOffset` that allow you to specify an *offset* in the callstack.  The offset is the first argument to these functions.

With this, we can rewrite our helper as:

    assertTurboencabulatorContains(components ...string) {
        components, err := turboEncabulator.GetComponents()
        ExpectWithOffset(1, err).NotTo(HaveOccurred())

        ExpectWithOffset(1, components).To(HaveLen(components))
        for _, component := range components {
            ExpectWithOffset(1, components).To(ContainElement(component))
        }
    }

now, failed assertions will point to the correct call to the helper in the test.

---

## Provided Matchers

Gomega comes with a bunch of `OmegaMatcher`s.  They're all documented here.  If there's one you'd like to see written either [send a pull request or open an issue](http://github.com/onsi/ginkgo).

These docs only go over the positive assertion case (`Should`), the negative case (`ShouldNot`) is simply the negation of the positive case.  They also use the `Ω` notation, but - as mentioned above - the `Expect` notation is equivalent.

### Equal(expected interface{})

    Ω(ACTUAL).Should(Equal(EXPECTED))

uses [`reflect.DeepEqual`](http://golang.org/pkg/reflect#deepequal) to compare `ACTUAL` with `EXPECTED`.

`reflect.DeepEqual` is awesome.  It will use `==` when appropriate (e.g. when comparing primitives) but will recursively dig into maps, slices, arrays, and even your own structs to ensure deep equality.  `reflect.DeepEqual`, however, is strict about comparing types.  Both `ACTUAL` and `EXPECTED` *must* have the same type.  If you want to compare across different types (e.g. if you've defined a type alias) you should use `BeEquivalentTo`

It is an error for both `ACTUAL` and `EXPECTED` to be nil, you should use `BeNil()` instead.

> For asserting equality between numbers of different types, you'll want to use the [`BeNumerically()`](#benumericallycomparator_string_compareto_interface) matcher

### BeEquivalentTo(expected interface{})

    Ω(ACTUAL).Should(BeEquivalentTo(EXPECTED))

Like `Equal`, `BeEquivalentTo` uses `reflect.DeepEqual` to compare `ACTUAL` with `EXPECTED`.  Unlike `Equal`, however, `BeEquivalentTo` will first convert `ACTUAL`s type to that of `EXPECTED` before making the comparison with `reflect.DeepEqual`.

This means that `BeEquivalentTo` will succesfully match equivalent values of different type.  This is particularly useful, for example, with type aliases:

    type FoodSource string

    Ω(FoodSource("Cheeseboard Pizza")).Should(Equal("Cheeseboard Pizza")) //will fail
    Ω(FoodSource("Cheeseboard Pizza")).Should(BeEquivalentTo("Cheeseboard Pizza")) //will pass

As with `Equal` it is an error for both `ACTUAL` and `EXPECTED` to be nil, you should use `BeNil()` instead.

As a rule, you **should not** use `BeEquivalentTo` with numbers.  Both of the following assertions are true:

    Ω(5.1).Should(BeEquivalentTo(5))
    Ω(5).ShouldNot(BeEquivalentTo(5.1))

the first assertion passes because 5.1 will be cast to an integer and will get rounded down!  Such false positives are terrible and should be avoided.  Use [`BeNumerically()`](#benumericallycomparator_string_compareto_interface) to compare numbers instead.

### BeNil()

    Ω(ACTUAL).Should(BeNil())

succeeds if `ACTUAL` is, in fact, `nil`.

### BeZero()

    Ω(ACTUAL).Should(BeZero())

succeeds if `ACTUAL` is the zero value for its type *or* if `ACTUAL` is `nil`.

### BeTrue()

    Ω(ACTUAL).Should(BeTrue())

succeeds if `ACTUAL` is `bool` typed and has the value `true`.  It is an error for `ACTUAL` to not be a `bool`.

> Some matcher libraries have a notion of `truthiness` to assert that an object is present.  Gomega is strict, and `BeTrue()` only works with `bool`s.  You can use `Ω(ACTUAL).ShouldNot(BeZero())` or `Ω(ACTUAL).ShouldNot(BeNil())` to verify object presence.

### BeFalse()

    Ω(ACTUAL).Should(BeFalse())

succeeds if `ACTUAL` is `bool` typed and has the value `false`.  It is an error for `ACTUAL` to not be a `bool`.

### HaveOccurred()

    Ω(ACTUAL).Should(HaveOccurred())

succeeds if `ACTUAL` is a non-nil `error`.  Thus, the typical Go error checking pattern looks like:

    err := SomethingThatMightFail()
    Ω(err).ShouldNot(HaveOccurred())

### MatchError(expected interface{})

    Ω(ACTUAL).Should(MatchError(EXPECTED))

succeeds if `ACTUAL` is a non-nil `error` that matches `EXPECTED`.  `EXPECTED` can be a string, in which case `ACTUAL.Error()` will be compared against `EXPECTED`.  Alternatively, `EXPECTED` can be an error, in which case `ACTUAL` and `ERROR` are compared via `reflect.DeepEqual`.  Any other type for `EXPECTED` is an error.

### BeClosed()

    Ω(ACTUAL).Should(BeClosed())

succeeds if actual is a closed channel. It is an error to pass a non-channel to `BeClosed`, it is also an error to pass `nil`.

In order to check whether or not the channel is closed, Gomega must try to read from the channel (even in the `ShouldNot(BeClosed())` case).  You should keep this in mind if you wish to make subsequent assertions about values coming down the channel.

Also, if you are testing that a *buffered* channel is closed you must first read all values out of the channel before asserting that it is closed (it is not possible to detect that a buffered-channel has been closed until all its buffered values are read).

Finally, as a corollary: it is an error to check whether or not a send-only channel is closed.

### Receive()

    Ω(ACTUAL).Should(Receive(<optionalPointer>))

succeeds if there is a message to be received on actual. Actual must be a channel (and cannot be a send-only channel) -- anything else is an error.

`Receive` returns *immediately*.  It *never* blocks:

- If there is nothing on the channel `c` then `Ω(c).Should(Receive())` will fail and `Ω(c).ShouldNot(Receive())` will pass.
- If there is something on the channel `c` ready to be read, then `Ω(c).Should(Receive())` will pass and `Ω(c).ShouldNot(Receive())` will fail.
- If the channel `c` is closed then *both* `Ω(c).Should(Receive())` and `Ω(c).ShouldNot(Receive())` will error.

If you have a go-routine running in the background that will write to channel `c`, for example:

    go func() {
        time.Sleep(100 * time.Millisecond)
        c <- true
    }()

you can assert that `c` receives something eventually:

    Eventually(c).Should(Receive())

This will timeout if nothing gets sent to `c` (you can modify the timeout interval as you normally do with `Eventually`).

A similar use-case is to assert that no go-routine writes to a channel (for a period of time).  You can do this with `Consistently`:

    Consistently(c).ShouldNot(Receive())

Finally, you often want to make assertions on the value *sent* to the channel.  You can ask the `Receive` matcher for the value passed
to the channel by passing it a pointer to a variable of the appropriate type:

    var receivedString string
    Eventually(stringChan).Should(Receive(&receivedString))
    Ω(receivedString).Shoudl(Equal("foo"))

### BeEmpty()

    Ω(ACTUAL).Should(BeEmpty())

succeeds if `ACTUAL` is, in fact, empty. `ACTUAL` must be of type `string`, `array`, `map`, `chan`, or `slice`.  It is an error for it to have any other type.

### HaveLen(count int)

    Ω(ACTUAL).Should(HaveLen(INT))

succeeds if the length of `ACTUAL` is `INT`. `ACTUAL` must be of type `string`, `array`, `map`, `chan`, or `slice`.  It is an error for it to have any other type.

### ContainSubstring(substr string, args ...interface{})

    Ω(ACTUAL).Should(ContainSubstring(STRING, ARGS...))

succeeds if `ACTUAL` contains the substring generated by:

    fmt.Sprintf(STRING, ARGS...)

`ACTUAL` must either be a `string`, `[]byte` or a `Stringer` (a type implementing the `String()` method).  Any other input is an error.

> Note, of course, that the `ARGS...` are not required.  They are simply a convenience to allow you to build up strings programmatically inline in the matcher.

### MatchRegexp(regexp string, args ...interface{})

    Ω(ACTUAL).Should(MatchRegexp(STRING, ARGS...))

succeeds if `ACTUAL` is matched by the regular expression string generated by:

    fmt.Sprintf(STRING, ARGS...)

`ACTUAL` must either be a `string`, `[]byte` or a `Stringer` (a type implementing the `String()` method).  Any other input is an error.  It is also an error for the regular expression to fail to compile.

> Note, of course, that the `ARGS...` are not required.  They are simply a convenience to allow you to build up strings programmatically inline in the matcher.

### MatchJSON(json interface{})

    Ω(ACTUAL).Should(MatchJSON(EXPECTED))

Both `ACTUAL` and `EXPECTED` must be a `string`, `[]byte` or a `Stringer`.  `MatchJSON` succeeds if bth `ACTUAL` and `EXPECTED` are JSON representations of the same object.  This is verified by parsing both `ACTUAL` and `EXPECTED` and then asserting equality on the resulting objects with `reflect.DeepEqual`.  By doing this `MatchJSON` avoids any issues related to white space, formatting, and key-ordering.

It is an error for either `ACTUAL` or `EXPECTED` to be invalid JSON.

### ContainElement(element interface{})

    Ω(ACTUAL).Should(ContainElement(ELEMENT))

succeeds if `ACTUAL` contains an element that equals `ELEMENT`.  `ACTUAL` must be an `array`, `slice`, or `map` -- anything else is an error.  For `map`s `ContainElement` searches through the map's values (not keys!).

By default `ContainElement()` uses the `Equal()` matcher under the hood to assert equality between `ACTUAL`'s elements and `ELEMENT`.  You can change this, however, by passing `ContainElement` an `OmegaMatcher`. For example, to check that a slice of strings has an element that matches a substring:

    Ω([]string{"Foo", "FooBar"}).Should(ContainElement(ContainSubstring("Bar")))

### HaveKey(key interface{})

    Ω(ACTUAL).Should(HaveKey(KEY))

succeeds if `ACTUAL` is a map with a key that equals `KEY`.  It is an error for `ACTUAL` to not be a `map`.

By default `HaveKey()` uses the `Equal()` matcher under the hood to assert equality between `ACTUAL`'s keys and `KEY`.  You can change this, however, by passing `HaveKey` an `OmegaMatcher`. For example, to check that a map has a key that matches a regular expression:

    Ω(map[string]string{"Foo": "Bar", "BazFoo": "Duck"}).Should(HaveKey(MatchRegexp(`.+Foo$`)))

### BeNumerically(comparator string, compareTo ...interface{})

    Ω(ACTUAL).Should(BeNumerically(COMPARATOR_STRING, EXPECTED, <THRESHOLD>))

performs numerical assertions in a type-agnostic way.  `ACTUAL` and `EXPECTED` should be numebers, though the specific type of number is irrelevant (`float32`, `float64`, `uint8`, etc...).  It is an error for `ACTUAL` or `EXPECTED` to not be a number.

There are six supported comparators:

- `Ω(ACTUAL).Should(BeNumerically("==", EXPECTED))`:
    Asserts that `ACTUAL` and `EXPECTED` are numerically equal

- `Ω(ACTUAL).Should(BeNumerically("~", EXPECTED, <THRESHOLD>))`:
    Asserts that `ACTUAL` and `EXPECTED` are within `<THRESHOLD>` of one another.  By default `<THRESHOLD>` is `1e-8` but you can specify a custom value.

- `Ω(ACTUAL).Should(BeNumerically(">", EXPECTED))`:
    Asserts that `ACTUAL` is greater than `EXPECTED`

- `Ω(ACTUAL).Should(BeNumerically(">=", EXPECTED))`:
    Asserts that `ACTUAL` is greater than or equal to  `EXPECTED`

- `Ω(ACTUAL).Should(BeNumerically("<", EXPECTED))`:
    Asserts that `ACTUAL` is less than `EXPECTED`

- `Ω(ACTUAL).Should(BeNumerically("<=", EXPECTED))`:
    Asserts that `ACTUAL` is less than or equal to `EXPECTED`

Any other comparator is an error.

### BeAssignableToTypeOf(expected interface)

    Ω(ACTUAL).Should(BeAssignableToTypeOf(EXPECTED interface))

succeeds if `ACTUAL` is a type that can be assigned to a variable with the same type as `EXPECTED`.  It is an error for either `ACTUAL` or `EXPECTED` to be `nil`.

### Panic()

    Ω(ACTUAL).Should(Panic())

succeeds if `ACTUAL` is a function that, when invoked, panics.  `ACTUAL` must be a function that takes no arguments and returns no result -- any other type for `ACTUAL` is an error.

---

## Adding Your Own Matchers

A matcher, in Gomega, is any type that satisfies the `OmegaMatcher` interface:

    type OmegaMatcher interface {
        Match(actual interface{}) (success bool, message string, err error)
        FailureMessage(actual interface{}) (message string)
        NegatedFailureMessage(actual interface{}) (message string)
    }

Writing domain-specific custom matchers is trivial and highly encouraged.  Let's work through an example.

### A Custom Matcher: RepresentJSONifiedObject(EXPECTED interface{})

Say you're working on a JSON API and you want to assert that your server returns the correct JSON representation.  Rather than marshal/unmarshal JSON in your tests, you want to write an expressive matcher that checks that the received response is a JSON representation for the object in question.  This is what the `RepresentJSONifiedObject` matcher could look like:

    package json_response_matcher

    import (
        "github.com/onsi/gomega"

        "encoding/json"
        "fmt"
        "net/http"
        "reflect"
    )

    func RepresentJSONifiedObject(expected interface{}) gomega.OmegaMatcher {
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

Let's break this down:

- Most matchers have a constructor function that returns an instance of the matcher.  In this case we've created `RepresentJSONifiedObject`.  Where possible, your constructor function should take explicit types or interfaces.  For our usecase, however, we need to accept any possible expected type so `RepresentJSONifiedObject` takes an argument with the generic `interface{}` type.
- The constructor function then initializes and returns an instance of our matcher: the `representJSONMatcher`.  These rarely need to be exported outside of your matcher package.
- The `representJSONMatcher` must satisfy the `OmegaMatcher` interface.  It does this by implementing the `Match`, `FailureMessage`, and `NegatedFailureMessage` method:
    - If the `OmegaMatcher` receives invalid inputs `Match` returns a non-Nil error explaining the problems with the input.  This allows Gomega to fail the assertion whether the assertion is for the positive or negative case.
    - If the `actual` and `expected` values match, `Match` should return `true`.
    - Similarly, if the `actual` and `expected` values do not match, `Match` should return `false`.
    - If the `OmegaMatcher` was testing the `Should` case, and `Match` returned false, `FailureMessage` will be called to print a message explaining the failure.
    - Likewise, if the `OmegaMatcher` was testing the `ShouldNot` case, and `Match` returned false, `NegatedFailureMessage` will be called.
    - It is guaranteed that `FailureMessage` and `NegatedFailureMessage` will only be called *after* `Match`, so you can save off any state you need to compute the messages in `Match`.
- Finally, it is common for matchers to make extensive use of the `reflect` library to interpret the generic inputs they receive.  In this case, the `representJSONMatcher` goes through some `reflect` gymnastics to create a pointer to a new object with the same type as the `expected` object, read and decode JSON from `actual` into that pointer, and then deference the pointer and compare the result to the `expected` object.

You might testdrive this matcher while writing it using Ginkgo.  Your test might look like:

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

This also offers an example of what using the matcher would look like in your tests.  Note that testing the cases when the matcher returns an error involves creating the matcher and invoking `Match` manually (instead of using an `Ω` or `Expect` assertion).

### Aborting Eventually/Consistently

There are sometimes instances where a `Eventually` and `Consistently` should stop polling a matcher because the result of the match simply cannot change.

For example, consider a test that looks like:
    
    Eventually(myChannel).Should(Receive(Equal("bar")))

`Eventually` will repeatedly invoke the `Receive` matcher against `myChannel` until the match succeeds.  However, if the channel becomes *closed* there is *no way* for the match to ever succeed.  Allowing `Eventually` to conitnue polling is inefficient and slows the test suite down.

To get around this, a matcher can optionally implement:

    MatchMayChangeInTheFuture(actual interface{}) bool

This is not part of the `OmegaMatcher` interface and, in general, most matchers do not need to implement `MatchMayChangeInTheFuture`.

If implemented, however, `MatchMayChangeInTheFuture` will be called with the appropriate `actual` value by `Eventually` and `Consistently` *after* the call to `Match` during every polling interval.  If `MatchMayChangeInTheFuture` returns `true`, `Eventually` and `Consistently` will continue polling.  If, however, `MatchMayChangeInTheFuture` returns `false`, `Eventually` and `Consistently` will abort and either fail or pass as appropriate.

If you'd like to look at a simple example `MatchMayChangeInTheFuture` check out [`gexec`'s `Exit` matcher](https://github.com/onsi/gomega/tree/master/gexec/exit_matcher.go).  Here, `MatchMayChangeInTheFuture` returns true if the `gexec.Session` under test has not exited yet, but returns false if it has.  Because of this, if a process exits with status code 3, but an assertion is made of the form:

    Eventually(session, 30).Should(gexec.Exit(0))

`Eventually` will not block for 30 seconds but will return (and fail, correctly) as soon as the mismatched exit code arrives!

> Note: `Eventually` and `Consistently` only excercise the `MatchMayChangeInTheFuture` method *if* they are passed a bare value.  If they are passed functions to be polled it is not possible to guarantee that the return value of the function will not change between polling intervals.  In this case, `MatchMayChangeInTheFuture` is not called and the polling continues until either a match is found or the timeout elapses.

### Contibuting to Gomega

Contributions are more than welcome.  Either [open an issue](http://github.com/onsi/gomega/issues) for a matcher you'd like to see or, better yet, test drive the matcher and [send a pull request](https://github.com/onsi/gomega/pulls).

When adding a new matcher please mimic the style use in Gomega's current matchers: you should use the tools in `formatSupport.go`, put the matcher and its tests in the `matchers` package, the constructor in the `matchers.go` file in the `omega` package.  Also, be sure to update the github-pages documentation (both `index.md` and `gomega_sidebar.html`) and include those changes in a separate pull request.
