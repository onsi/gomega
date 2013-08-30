package gomega

import "time"

var globalFailHandler OmegaFailHandler

type OmegaFailHandler func(message string, callerSkip ...int)

func RegisterFailHandler(handler OmegaFailHandler) {
	globalFailHandler = handler
}

func Ω(actual interface{}) Actual {
	return newActual(actual, globalFailHandler)
}

func Expect(actual interface{}) Actual {
	return newActual(actual, globalFailHandler)
}

func Eventually(actual interface{}, intervals ...float64) AsyncActual {
	timeoutInterval := time.Duration(5 * time.Second)
	pollingInterval := time.Duration(100 * time.Millisecond)
	if len(intervals) > 0 {
		timeoutInterval = time.Duration(intervals[0] * float64(time.Second))
	}
	if len(intervals) > 1 {
		pollingInterval = time.Duration(intervals[1] * float64(time.Second))
	}
	return newAsyncActual(actual, globalFailHandler, timeoutInterval, pollingInterval)
}

type AsyncActual interface {
	Should(matcher OmegaMatcher, optionalDescription ...interface{})
	ShouldNot(matcher OmegaMatcher, optionalDescription ...interface{})
}

type Actual interface {
	Should(matcher OmegaMatcher, optionalDescription ...interface{})
	ShouldNot(matcher OmegaMatcher, optionalDescription ...interface{})

	To(matcher OmegaMatcher, optionalDescription ...interface{})
	ToNot(matcher OmegaMatcher, optionalDescription ...interface{})
	NotTo(matcher OmegaMatcher, optionalDescription ...interface{})
}

type OmegaMatcher interface {
	Match(actual interface{}) (success bool, message string, err error)
}
