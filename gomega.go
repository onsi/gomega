/*
Gomega is the Ginkgo BDD-style testing framework's preferred matcher library.

The godoc documentation describes Gomega's API.  More comprehensive documentation (with examples!) is available at http://onsi.github.io/gomega/

Gomega on Github: http://github.com/onsi/gomega

Learn more about Ginkgo online: http://onsi.github.io/ginkgo

Ginkgo on Github: http://github.com/onsi/ginkgo

Gomega is MIT-Licensed
*/
package gomega

import "time"

const GOMEGA_VERSION = "0.9"

var globalFailHandler OmegaFailHandler

type OmegaFailHandler func(message string, callerSkip ...int)

//RegisterFailHandler connects Ginkgo to Gomega.  When a matcher fails
//the fail handler passed into RegisterFailHandler is called.
func RegisterFailHandler(handler OmegaFailHandler) {
	globalFailHandler = handler
}

//Ω wraps an actual value allowing assertions to be made on it:
//	Ω("foo").Should(Equal("foo"))
//
//Ω and Expect are identical
func Ω(actual interface{}) Actual {
	return newActual(actual, globalFailHandler)
}

//Expect wraps an actual value allowing assertions to be made on it:
//	Expect("foo").To(Equal("foo"))
//
//Expect and Ω are identical
func Expect(actual interface{}) Actual {
	return newActual(actual, globalFailHandler)
}

//Eventually wraps an actual value allowing assertions to be made on it.
//The assertion is tried periodically until it passes or a timeout occurs.
//
//Both the timeout and polling interval are configurable as optional arguments:
//The first optional argument is the timeout in seconds expressed as a float64.
//The second optional argument is the polling interval in seconds expressd as a float64.
//
//If Eventually is passed an actual that is a function taking no arguments and returning one value,
//then Eventually will call the function periodically and try the matcher against the function's return value.
//
//Example:
//
//  Eventually(func() int {
//    return thingImPolling.Count()
//  }).Should(BeNumerically(">=", 17))
func Eventually(actual interface{}, intervals ...float64) AsyncActual {
	timeoutInterval := time.Duration(1 * time.Second)
	pollingInterval := time.Duration(10 * time.Millisecond)
	if len(intervals) > 0 {
		timeoutInterval = time.Duration(intervals[0] * float64(time.Second))
	}
	if len(intervals) > 1 {
		pollingInterval = time.Duration(intervals[1] * float64(time.Second))
	}
	return newAsyncActual(actual, globalFailHandler, timeoutInterval, pollingInterval)
}

//AsyncActual is returned by Eventually and polls the actual value passed into Eventually against
//the matcher passed to the Should and ShouldNot methods.
//
//Both Should and ShouldNot take a variadic optionalDescription argument.  This is passed on to
//fmt.Sprintf() and is used to annotate failure messages.  This allows you to make your failure messages more
//descriptive
//
//Both Should and ShouldNot return a boolean that is true if the assertion passed and false if it failed.
//
//Example:
//
//  Eventually(myChannel).Should(HaveLen(1), "Something should have come down the pipe.")
type AsyncActual interface {
	Should(matcher OmegaMatcher, optionalDescription ...interface{}) bool
	ShouldNot(matcher OmegaMatcher, optionalDescription ...interface{}) bool
}

//Actual is returned by Ω and Expect and compares the actual value to the matcher
//passed to the Should/ShouldNot and To/ToNot/NotTo methods.
//
//Typically Should/ShouldNot are used with Ω and To/ToNot/NotTo are used with Expect
//though this is not enforced.
//
//All methods take a variadic optionalDescription argument.  This is passed on to fmt.Sprintf()
//and is used to annotate failure messages.
//
//All methods return a bool that is true if hte assertion passed and false if it failed.
//
//Example:
//
//   Ω(farm.HasCow()).Should(BeTrue(), "Farm %v should have a cow", farm)
type Actual interface {
	Should(matcher OmegaMatcher, optionalDescription ...interface{}) bool
	ShouldNot(matcher OmegaMatcher, optionalDescription ...interface{}) bool

	To(matcher OmegaMatcher, optionalDescription ...interface{}) bool
	ToNot(matcher OmegaMatcher, optionalDescription ...interface{}) bool
	NotTo(matcher OmegaMatcher, optionalDescription ...interface{}) bool
}

//All Gomega matchers must implement the OmegaMatcher interface
//
//For details on writing custom matchers, check out: http://onsi.github.io/gomega/#adding_your_own_matchers
type OmegaMatcher interface {
	Match(actual interface{}) (success bool, message string, err error)
}
