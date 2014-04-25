package gomega

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	"testing"
)

func Test(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gomega")
}

func TestTestingT(t *testing.T) {
	RegisterTestingT(t)
	Ω(true).Should(BeTrue())
}

type fakeMatcher struct {
	receivedActual  interface{}
	matchesToReturn bool
	errToReturn     error
}

func (matcher *fakeMatcher) Match(actual interface{}) (bool, error) {
	matcher.receivedActual = actual

	return matcher.matchesToReturn, matcher.errToReturn
}

func (matcher *fakeMatcher) FailureMessage(actual interface{}) string {
	return fmt.Sprintf("positive: %v", actual)
}

func (matcher *fakeMatcher) NegatedFailureMessage(actual interface{}) string {
	return fmt.Sprintf("negative: %v", actual)
}

func interceptFailures(f func()) []string {
	failures := []string{}
	RegisterFailHandler(func(message string, callerSkip ...int) {
		failures = append(failures, message)
	})
	f()
	RegisterFailHandler(Fail)
	return failures
}
