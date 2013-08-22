package gomega

import "fmt"

var globalFailHandler OmegaFailHandler

type OmegaFailHandler func(message string, callerSkip ...int)

func RegisterFailHandler(handler OmegaFailHandler) {
	fmt.Println("HI!", handler)
	globalFailHandler = handler
}

func Ω(actual interface{}) Actual {
	return newActual(actual, globalFailHandler)
}

func Expect(actual interface{}) Actual {
	return newActual(actual, globalFailHandler)
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
