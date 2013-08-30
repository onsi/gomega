package gomega

import (
	. "github.com/onsi/ginkgo"
	"testing"
)

func Test(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gomega")
}

type fakeMatcher struct {
	receivedActual  interface{}
	matchesToReturn bool
	messageToReturn string
	errToReturn     error
}

func (matcher *fakeMatcher) Match(actual interface{}) (bool, string, error) {
	matcher.receivedActual = actual

	return matcher.matchesToReturn, matcher.messageToReturn, matcher.errToReturn
}
