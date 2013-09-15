---
layout: default
title: Gomega
sidebar: Gomega
---

#Gomega: Ginkgo's Preferred Matcher Library

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

If you use the `ginkgo` CLI to `ginkgo bootstrap` a test suite, thie command will be automatically generated for you.

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

On OS X the `Ω` character is easy to type.  Just hit option-z: ⌥z

On the left hand side, you can pass anything you want in to `Ω` and `Expect` for `ACTUAL`.  On the right hand side you must pass an object that satisfies the `OmegaMatcher` interface.  Omega's matchers (e.g. `Equal(EXPECTED)`) are simply functions that create an appropriate `OmegaMatcher` object, passing it the `EXPECTED`.

> The `OmegaMatcher` interface is very simple and is discussed [below](#adding-your-own-matchers).

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

Gomega has support for making *asynchronous* assertions.  You do this by passing a function to `Eventually`:

    Eventually(func() []int {
        return thing.SliceImMonitoring
    }).Should(HaveLen(2))

    Eventually(func() string {
        return thing.Status
    }).ShouldNot(Equal("Stuck Waiting"))

`Eventually` will poll the passed in function (which must have zero-arguments and one return value) repeatedly and check the return value against the `OmegaMatcher`.  `Eventually` then blocks until the match succeeds or until a timeout interval has elapsed.

The default value for the timeout is 5 seconds and the default value for the polling interval is 100 milliseconds.  You can change these values by passing in float64s (in seconds) just after your function:

    Eventually(func() []int {
        return thing.SliceImMonitoring
    }, TIMEOUT_IN_SECONDS, POLLING_INTERVAL_IN_SECONDS).Should(HaveLen(2))

`Eventually` is especially handy when writing integration tests against asynchronous services or components:
    
    externalProcess.DoSomethingAmazing()
    Eventually(func() bool {
        return somethingAmazingHappened()
    }).Should(BeTrue())

> As with synchronous assertions, you can annotate asynchronous assertions by passing a format string and optional inputs after the `OmegaMatcher`.

---

## Provided Matchers

Gomega comes with a bunch of `OmegaMatcher`s.  They're all documented here.  If there's one you'd like to see written either [send a pull request or open an issue](http://github.com/onsi/ginkgo).

These docs only go over the positive assertion case (`Should`), the negative case (`ShouldNot`) is simply the negation of the positive case.  They also use the `Ω` notation, but - as mentioned above - the `Expect` notation is equivalent.

### Equal(expected interface{})

    Ω(ACTUAL).Should(Equal(EXPECTED))

uses [`reflect.DeepEqual`](http://golang.org/pkg/reflect#deepequal) to compare `ACTUAL` with `EXPECTED`.

`reflect.DeepEqual` is awesome.  It will use `==` when appropriate (e.g. when comparing primitives) but will recursively dig into maps, slices, arrays, and even your own structs to ensure deep equality.

It is an error for both `ACTUAL` and `EXPECTED` to be nil, you should use `BeNil()` instead.

> For asserting equality between numbers of different types, you'll want to use the [`BeNumerically()`](#benumerically) matcher

### BeNil()

    Ω(ACTUAL).Should(BeNil())

succeeds if `ACTUAL` is, in fact, `nil`.

### BeZero()

    Ω(ACTUAL).Should(BeZero())

succeeds if `ACTUAL` is the zero value for its type *or* if `ACTUAL` is `nil`.

### BeTrue()

    Ω(ACTUAL).Should(BeTrue())

succeeds if `ACTUAL` is `bool` typed and has the value `true`.  It is an error for `ACTUAL` to not be a `bool`.

> Some matcher libraries have a notion of `truthiness` to assert that an object is present.  You can use `Ω(ACTUAL).ShouldNot(BeZero())` or `Ω(ACTUAL).ShouldNot(BeNil())` to achieve a similar effect.

### BeFalse()

    Ω(ACTUAL).Should(BeTrue())

succeeds if `ACTUAL` is `bool` typed and has the value `false`.  It is an error for `ACTUAL` to not be a `bool`.

### HaveOccured()

    Ω(ERROR).Should(HaveOccured())

succeeds if `ACTUAL` is a non-nil `error`.  Thus, the typical Go error checking pattern looks like:

    err := SomethingThatMightFail()
    Ω(err).ShouldNot(HaveOccured())

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

`ACTUAL` must either be a `string` or a `Stringer` (a type implementing the `String()` method).  Any other input is an error.

> Note, of course, that the `ARGS...` are not required.  They are simply a convenience to allow you to build up strings programmatically inline in the matcher.

### MatchRegexp(regexp string, args ...interface{})

    Ω(ACTUAL).Should(ContainSubstring(STRING, ARGS...))

succeeds if `ACTUAL` is matched by the regular expression string generated by:
    
    fmt.Sprintf(STRING, ARGS...)

`ACTUAL` must either be a `string` or a `Stringer` (a type implementing the `String()` method).  Any other input is an error.  It is also an error for the regular expression to fail to compile.

> Note, of course, that the `ARGS...` are not required.  They are simply a convenience to allow you to build up strings programmatically inline in the matcher.

### ContainElement(element interface{})

    Ω(ACTUAL).Should(ContainElement(ELEMENT))

succeeds if `ACTUAL` contains an element that equals `ELEMENT`.  `ACTUAL` must be an `array` or `slice` -- anything else is an error.

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

### Panic()

    Ω(ACTUAL).Should(Panic())

succeeds if `ACTUAL` is a function that, when invoked, panics.  `ACTUAL` must be a function that takes no arguments and returns no result -- any other type for `ACTUAL` is an error.

---

## Adding Your Own Matchers

A matcher, in Gomega, is any type that satisfies the `OmegaMatcher` interface:

    type OmegaMatcher interface {
        Match(actual interface{}) (success bool, message string, err error)
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

    func (matcher *representJSONMatcher) Match(actual interface{}) (success bool, message string, err error) {
        response, ok := actual.(*http.Response)
        if !ok {
            return false, "", fmt.Errorf("RepresentJSONifiedObject matcher expects an http.Response")
        }

        pointerToObjectOfExpectedType := reflect.New(reflect.TypeOf(matcher.expected)).Interface()
        err = json.NewDecoder(response.Body).Decode(pointerToObjectOfExpectedType)

        if err != nil {
            return false, "", fmt.Errorf("Failed to decode JSON: %s", err.Error())
        }

        decodedObject := reflect.ValueOf(pointerToObjectOfExpectedType).Elem().Interface()

        if reflect.DeepEqual(decodedObject, matcher.expected) {
            return true, fmt.Sprintf("Expected\n\t%#v\nnot to contain the JSON representation of\n\t%#v", actual, matcher.expected), nil
        } else {
            return false, fmt.Sprintf("Expected\n\t%#v\nto contain the JSON representation of\n\t%#v", actual, matcher.expected), nil
        }
    }

Let's break this down:

- Most matchers have a constructor function that returns an instance of the matcher.  In this case we've created `RepresentJSONifiedObject`.  Where possible, your constructor function should take explicit types or interfaces.  For our usecase, however, we need to accept any possible expected type so `RepresentJSONifiedObject` takes an argument with the generic `interface{}` type.
- The constructor function then initializes and returns an instance of our matcher: the `representJSONMatcher`.  These rarely need to be exported outside of your matcher package.
- The `representJSONMatcher` must satisfy the `OmegaMatcher` interface.  It does this by implementing the `Match` method:
    - If the `OmegaMatcher` receives invalid inputs it returns a non-Nil error explaining the problems with the input.  This allows Gomega to fail the assertion whether the assertion is for the positive or negative case.
    - If the `actual` and `expected` values match the matcher should return `true`.  **The matcher should *also* return a failure message appropriate for the *negative* assertion case**.  This is important.  The matcher is not told whether the assertion is positive (`Should`) or negative (`ShouldNot`).  Instead, if it identifies a positive match (by returning `true`) it must return a message that `ShouldNot` can print.
    - Similarly, if the `actual` and `expected` values do not match, the matcher should return `false` along-with a failure message appropriate for the positive assertion (`Should`).
- Finally, it is common for matchers to make extensive use of the `reflect` library to interpret the generic inputs they receive.  In this case, the `representJSONMatcher` goes through some `reflect` gymnastics to create a pointer to a new object with the same type as the `expected` object, read and decode JSON from `actual` into that pointer, and then deference the pointer and compare the result to the `expected` object.

You might testdrive this matcher while writing it using Ginkgo.  Your test might look like:

    package custom_matcher

    import (
        . "github.com/onsi/ginkgo"
        . "github.com/onsi/gomega"

        "bytes"
        "encoding/json"
        "io/ioutil"
        "net/http"
        "strings"

        "testing"
    )

    func TestBootstrap(t *testing.T) {
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
            Ω(err).ShouldNot(HaveOccured())
        })

        Context("when actual is not an http response", func() {
            It("should error", func() {
                _, _, err := RepresentJSONifiedObject(book).Match("not a response")
                Ω(err).Should(HaveOccured())
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
                    _, _, err := RepresentJSONifiedObject(book).Match(response)
                    Ω(err).Should(HaveOccured())
                })
            })
        })
    })

This also offers an example of what using the matcher would look like in your tests.  Note that testing the cases when the matcher returns an error involves creating the matcher and invoking `Match` manually (instead of using an `Ω` or `Expect` assertion).

### Contibuting to Gomega

Contributions are more than welcome.  Either [open an issue](http://github.com/onsi/gomega/issues) for a matcher you'd like to see or, better yet, test drive the matcher and [send a pull request](https://github.com/onsi/gomega/pulls).

When adding a new matcher please mimic the style use in Gomega's current matchers (you should use the tools in `formatSupport.go` and put your tests in `matcher_tests`).  Also, be sure to update the github-pages documentation and include those changes in a separate pull request.