/*
Gomega is the Ginkgo BDD-style testing framework's preferred matcher library.

The godoc documentation describes Gomega's API.  More comprehensive documentation (with examples!) is available at http://onsi.github.io/gomega/

Gomega on Github: http://github.com/onsi/gomega

Learn more about Ginkgo online: http://onsi.github.io/ginkgo

Ginkgo on Github: http://github.com/onsi/ginkgo

Gomega is MIT-Licensed
*/
package gomega

import (
	"errors"
	"fmt"
	"time"

	"github.com/onsi/gomega/internal"
	"github.com/onsi/gomega/types"
)

const GOMEGA_VERSION = "1.15.0"

const nilGomegaPanic = `You are trying to make an assertion, but haven't registered Gomega's fail handler.
If you're using Ginkgo then you probably forgot to put your assertion in an It().
Alternatively, you may have forgotten to register a fail handler with RegisterFailHandler() or RegisterTestingT().
Depending on your vendoring solution you may be inadvertently importing gomega and subpackages (e.g. ghhtp, gexec,...) from different locations.
`

// Gomega describes the essential Gomega DSL. This interface allows libraries
// to abstract between the standard package-level function implementations
// and alternatives like *WithT.
//
// The types in the top-level DSL have gotten a bit messy due to earlier depracations that avoid stuttering
// and due to an accidental use of a concrete type (*WithT) in an earlier release.
//
// As of 1.15 both the WithT and Ginkgo variants of Gomega are implemented by the same underlying object
// however one (the Ginkgo variant) is exported as an interface (types.Gomega) whereas the other (the withT variant)
// is shared as a concrete type (*WithT, which is aliased to *internal.Gomega).  1.15 did not clean this mess up to ensure
// that declarations of *WithT in existing code are not broken by the upgrade to 1.15.
type Gomega = types.Gomega

// DefaultGomega supplies the standard package-level implementation
var Default = Gomega(internal.NewGomega(internal.FetchDefaultDurationBundle()))

// NewGomega returns an instance of Gomega wired into the passed-in fail handler.
// You generally don't need to use this when using Ginkgo - RegisterFailHandler will wire up the global gomega
// However creating a NewGomega with a custom fail handler can be useful in contexts where you want to use Gomega's
// rich ecosystem of matchers without causing a test to fail.  For example, to aggregate a series of potential failures
// or for use in a non-test setting.
func NewGomega(fail types.GomegaFailHandler) Gomega {
	return internal.NewGomega(Default.(*internal.Gomega).DurationBundle).ConfigureWithFailHandler(fail)
}

// WithT wraps a *testing.T and provides `Expect`, `Eventually`, and `Consistently` methods.  This allows you to leverage
// Gomega's rich ecosystem of matchers in standard `testing` test suites.
//
// Use `NewWithT` to instantiate a `WithT`
//
// As of 1.15 both the WithT and Ginkgo variants of Gomega are implemented by the same underlying object
// however one (the Ginkgo variant) is exported as an interface (types.Gomega) whereas the other (the withT variant)
// is shared as a concrete type (*WithT, which is aliased to *internal.Gomega).  1.15 did not clean this mess up to ensure
// that declarations of *WithT in existing code are not broken by the upgrade to 1.15.
type WithT = internal.Gomega

// GomegaWithT is deprecated in favor of gomega.WithT, which does not stutter.
type GomegaWithT = WithT

// NewWithT takes a *testing.T and returngs a `gomega.WithT` allowing you to use `Expect`, `Eventually`, and `Consistently` along with
// Gomega's rich ecosystem of matchers in standard `testing` test suits.
//
//    func TestFarmHasCow(t *testing.T) {
//        g := gomega.NewWithT(t)
//
//        f := farm.New([]string{"Cow", "Horse"})
//        g.Expect(f.HasCow()).To(BeTrue(), "Farm should have cow")
//     }
func NewWithT(t types.GomegaTestingT) *WithT {
	return internal.NewGomega(Default.(*internal.Gomega).DurationBundle).ConfigureWithT(t)
}

// NewGomegaWithT is deprecated in favor of gomega.NewWithT, which does not stutter.
var NewGomegaWithT = NewWithT

// RegisterFailHandler connects Ginkgo to Gomega. When a matcher fails
// the fail handler passed into RegisterFailHandler is called.
func RegisterFailHandler(fail types.GomegaFailHandler) {
	Default.(*internal.Gomega).ConfigureWithFailHandler(fail)
}

// RegisterFailHandlerWithT is deprecated and will be removed in a future release.
// users should use RegisterFailHandler, or RegisterTestingT
func RegisterFailHandlerWithT(_ types.GomegaTestingT, fail types.GomegaFailHandler) {
	fmt.Println("RegisterFailHandlerWithT is deprecated.  Please use RegisterFailHandler or RegisterTestingT instead.")
	Default.(*internal.Gomega).ConfigureWithFailHandler(fail)
}

// RegisterTestingT connects Gomega to Golang's XUnit style
// Testing.T tests.  It is now deprecated and you should use NewWithT() instead to get a fresh instance of Gomega for each test.
func RegisterTestingT(t types.GomegaTestingT) {
	Default.(*internal.Gomega).ConfigureWithT(t)
}

// InterceptGomegaFailures runs a given callback and returns an array of
// failure messages generated by any Gomega assertions within the callback.
// Exeuction continues after the first failure allowing users to collect all failures
// in the callback.
//
// This is most useful when testing custom matchers, but can also be used to check
// on a value using a Gomega assertion without causing a test failure.
func InterceptGomegaFailures(f func()) []string {
	originalHandler := Default.(*internal.Gomega).Fail
	failures := []string{}
	Default.(*internal.Gomega).Fail = func(message string, callerSkip ...int) {
		failures = append(failures, message)
	}
	defer func() {
		Default.(*internal.Gomega).Fail = originalHandler
	}()
	f()
	return failures
}

// InterceptGomegaFailure runs a given callback and returns the first
// failure message generated by any Gomega assertions within the callback, wrapped in an error.
//
// The callback ceases execution as soon as the first failed assertion occurs, however Gomega
// does not register a failure with the FailHandler registered via RegisterFailHandler - it is up
// to the user to decide what to do with the returned error
func InterceptGomegaFailure(f func()) (err error) {
	originalHandler := Default.(*internal.Gomega).Fail
	Default.(*internal.Gomega).Fail = func(message string, callerSkip ...int) {
		err = errors.New(message)
		panic("stop execution")
	}

	defer func() {
		Default.(*internal.Gomega).Fail = originalHandler
		if e := recover(); e != nil {
			if err == nil {
				panic(e)
			}
		}
	}()

	f()
	return err
}

func ensureDefaultGomegaIsConfigured() {
	if !Default.(*internal.Gomega).IsConfigured() {
		panic(nilGomegaPanic)
	}
}

// Ω wraps an actual value allowing assertions to be made on it:
//    Ω("foo").Should(Equal("foo"))
//
// If Ω is passed more than one argument it will pass the *first* argument to the matcher.
// All subsequent arguments will be required to be nil/zero.
//
// This is convenient if you want to make an assertion on a method/function that returns
// a value and an error - a common patter in Go.
//
// For example, given a function with signature:
//    func MyAmazingThing() (int, error)
//
// Then:
//    Ω(MyAmazingThing()).Should(Equal(3))
// Will succeed only if `MyAmazingThing()` returns `(3, nil)`
//
// Ω and Expect are identical
func Ω(actual interface{}, extra ...interface{}) Assertion {
	ensureDefaultGomegaIsConfigured()
	return Default.Ω(actual, extra...)
}

// Expect wraps an actual value allowing assertions to be made on it:
//    Expect("foo").To(Equal("foo"))
//
// If Expect is passed more than one argument it will pass the *first* argument to the matcher.
// All subsequent arguments will be required to be nil/zero.
//
// This is convenient if you want to make an assertion on a method/function that returns
// a value and an error - a common patter in Go.
//
// For example, given a function with signature:
//    func MyAmazingThing() (int, error)
//
// Then:
//    Expect(MyAmazingThing()).Should(Equal(3))
// Will succeed only if `MyAmazingThing()` returns `(3, nil)`
//
// Expect and Ω are identical
func Expect(actual interface{}, extra ...interface{}) Assertion {
	ensureDefaultGomegaIsConfigured()
	return Default.Expect(actual, extra...)
}

// ExpectWithOffset wraps an actual value allowing assertions to be made on it:
//    ExpectWithOffset(1, "foo").To(Equal("foo"))
//
// Unlike `Expect` and `Ω`, `ExpectWithOffset` takes an additional integer argument
// that is used to modify the call-stack offset when computing line numbers.
//
// This is most useful in helper functions that make assertions.  If you want Gomega's
// error message to refer to the calling line in the test (as opposed to the line in the helper function)
// set the first argument of `ExpectWithOffset` appropriately.
func ExpectWithOffset(offset int, actual interface{}, extra ...interface{}) Assertion {
	ensureDefaultGomegaIsConfigured()
	return Default.ExpectWithOffset(offset, actual, extra...)
}

// Eventually wraps an actual value allowing assertions to be made on it.
// The assertion is tried periodically until it passes or a timeout occurs.
//
// Both the timeout and polling interval are configurable as optional arguments:
// The first optional argument is the timeout
// The second optional argument is the polling interval
//
// Both intervals can either be specified as time.Duration, parsable duration strings or as floats/integers.  In the
// last case they are interpreted as seconds.
//
// If Eventually is passed an actual that is a function taking no arguments,
// then Eventually will call the function periodically and try the matcher against the function's first return value.
//
// Example:
//
//    Eventually(func() int {
//        return thingImPolling.Count()
//    }).Should(BeNumerically(">=", 17))
//
// Note that this example could be rewritten:
//
//    Eventually(thingImPolling.Count).Should(BeNumerically(">=", 17))
//
// If the function returns more than one value, then Eventually will pass the first value to the matcher and
// assert that all other values are nil/zero.
// This allows you to pass Eventually a function that returns a value and an error - a common pattern in Go.
//
// For example, consider a method that returns a value and an error:
//    func FetchFromDB() (string, error)
//
// Then
//    Eventually(FetchFromDB).Should(Equal("hasselhoff"))
//
// Will pass only if the the returned error is nil and the returned string passes the matcher.
//
// Eventually allows you to make assertions in the pased-in function.  The function is assumed to have failed and will be retried if any assertion in the function fails.
// For example:
//
//     Eventually(func() Widget {
//	     resp, err := http.Get(url)
//       Expect(err).NotTo(HaveOccurred())
//       defer resp.Body.Close()
//       Expect(resp.SatusCode).To(Equal(http.StatusOK))
//       var widget Widget
//       Expect(json.NewDecoder(resp.Body).Decode(&widget)).To(Succeed())
//       return widget
//     }).Should(Equal(expectedWidget))
//
// will keep trying the passed-in function until all its assertsions pass (i.e. the http request succeeds) _and_ the returned object satisfies the passed-in matcher.
//
// Functions passed to Eventually typically have a return value.  However you are allowed to pass in a function with no return value.  Eventually assumes such a function
// is making assertions and will turn it into a function that returns an error if any assertion fails, or nil if no assertion fails.  This allows you to use the Succeed() matcher
// to express that a complex operation should eventually succeed.  For example:
//
//    Eventually(func() {
//        model, err := db.Find("foo")
//        Expect(err).NotTo(HaveOccurred())
//        Expect(model.Reticulated()).To(BeTrue())
//        Expect(model.Save()).To(Succeed())
//    }).Should(Succeed())
//
// will rerun the function until all its assertions pass.
//
// Eventually's default timeout is 1 second, and its default polling interval is 10ms
func Eventually(actual interface{}, intervals ...interface{}) AsyncAssertion {
	ensureDefaultGomegaIsConfigured()
	return Default.Eventually(actual, intervals...)
}

// EventuallyWithOffset operates like Eventually but takes an additional
// initial argument to indicate an offset in the call stack.  This is useful when building helper
// functions that contain matchers.  To learn more, read about `ExpectWithOffset`.
func EventuallyWithOffset(offset int, actual interface{}, intervals ...interface{}) AsyncAssertion {
	ensureDefaultGomegaIsConfigured()
	return Default.EventuallyWithOffset(offset, actual, intervals...)
}

// Consistently wraps an actual value allowing assertions to be made on it.
// The assertion is tried periodically and is required to pass for a period of time.
//
// Both the total time and polling interval are configurable as optional arguments:
// The first optional argument is the duration that Consistently will run for
// The second optional argument is the polling interval
//
// Both intervals can either be specified as time.Duration, parsable duration strings or as floats/integers.  In the
// last case they are interpreted as seconds.
//
// If Consistently is passed an actual that is a function taking no arguments.
//
// If the function returns one value, then Consistently will call the function periodically and try the matcher against the function's first return value.
//
// If the function returns more than one value, then Consistently will pass the first value to the matcher and
// assert that all other values are nil/zero.
// This allows you to pass Consistently a function that returns a value and an error - a common pattern in Go.
//
// Like Eventually, Consistently allows you to make assertions in the function.  If any assertion fails Consistently will fail.  In addition,
// Consistently also allows you to pass in a function with no return value.  In this case Consistently can be paired with the Succeed() matcher to assert
// that no assertions in the function fail.
//
// Consistently is useful in cases where you want to assert that something *does not happen* over a period of time.
// For example, you want to assert that a goroutine does *not* send data down a channel.  In this case, you could:
//
//   Consistently(channel).ShouldNot(Receive())
//
// Consistently's default duration is 100ms, and its default polling interval is 10ms
func Consistently(actual interface{}, intervals ...interface{}) AsyncAssertion {
	ensureDefaultGomegaIsConfigured()
	return Default.Consistently(actual, intervals...)
}

// ConsistentlyWithOffset operates like Consistently but takes an additional
// initial argument to indicate an offset in the call stack. This is useful when building helper
// functions that contain matchers. To learn more, read about `ExpectWithOffset`.
func ConsistentlyWithOffset(offset int, actual interface{}, intervals ...interface{}) AsyncAssertion {
	ensureDefaultGomegaIsConfigured()
	return Default.ConsistentlyWithOffset(offset, actual, intervals...)
}

// SetDefaultEventuallyTimeout sets the default timeout duration for Eventually. Eventually will repeatedly poll your condition until it succeeds, or until this timeout elapses.
func SetDefaultEventuallyTimeout(t time.Duration) {
	Default.SetDefaultEventuallyTimeout(t)
}

// SetDefaultEventuallyPollingInterval sets the default polling interval for Eventually.
func SetDefaultEventuallyPollingInterval(t time.Duration) {
	Default.SetDefaultEventuallyPollingInterval(t)
}

// SetDefaultConsistentlyDuration sets  the default duration for Consistently. Consistently will verify that your condition is satisfied for this long.
func SetDefaultConsistentlyDuration(t time.Duration) {
	Default.SetDefaultConsistentlyDuration(t)
}

// SetDefaultConsistentlyPollingInterval sets the default polling interval for Consistently.
func SetDefaultConsistentlyPollingInterval(t time.Duration) {
	Default.SetDefaultConsistentlyPollingInterval(t)
}

// AsyncAssertion is returned by Eventually and Consistently and polls the actual value passed into Eventually against
// the matcher passed to the Should and ShouldNot methods.
//
// Both Should and ShouldNot take a variadic optionalDescription argument.
// This argument allows you to make your failure messages more descriptive.
// If a single argument of type `func() string` is passed, this function will be lazily evaluated if a failure occurs
// and the returned string is used to annotate the failure message.
// Otherwise, this argument is passed on to fmt.Sprintf() and then used to annotate the failure message.
//
// Both Should and ShouldNot return a boolean that is true if the assertion passed and false if it failed.
//
// Example:
//
//   Eventually(myChannel).Should(Receive(), "Something should have come down the pipe.")
//   Consistently(myChannel).ShouldNot(Receive(), func() string { return "Nothing should have come down the pipe." })
type AsyncAssertion = types.AsyncAssertion

// GomegaAsyncAssertion is deprecated in favor of AsyncAssertion, which does not stutter.
type GomegaAsyncAssertion = types.AsyncAssertion

// Assertion is returned by Ω and Expect and compares the actual value to the matcher
// passed to the Should/ShouldNot and To/ToNot/NotTo methods.
//
// Typically Should/ShouldNot are used with Ω and To/ToNot/NotTo are used with Expect
// though this is not enforced.
//
// All methods take a variadic optionalDescription argument.
// This argument allows you to make your failure messages more descriptive.
// If a single argument of type `func() string` is passed, this function will be lazily evaluated if a failure occurs
// and the returned string is used to annotate the failure message.
// Otherwise, this argument is passed on to fmt.Sprintf() and then used to annotate the failure message.
//
// All methods return a bool that is true if the assertion passed and false if it failed.
//
// Example:
//
//    Ω(farm.HasCow()).Should(BeTrue(), "Farm %v should have a cow", farm)
type Assertion = types.Assertion

// GomegaAssertion is deprecated in favor of Assertion, which does not stutter.
type GomegaAssertion = types.Assertion

// OmegaMatcher is deprecated in favor of the better-named and better-organized types.GomegaMatcher but sticks around to support existing code that uses it
type OmegaMatcher = types.GomegaMatcher
