package gomega

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	"time"
)

type deprecatedBeTrueMatcher struct{}

func (d *deprecatedBeTrueMatcher) Match(actual interface{}) (success bool, message string, err error) {
	boolValue, isBool := actual.(bool)
	if !isBool {
		return false, "", fmt.Errorf("Expected a boolean.  Got: %v", actual)
	}

	if boolValue {
		return true, "should not be true", nil
	} else {
		return false, "should be true", nil
	}
}

var _ = Describe("DeprecationShim", func() {
	var failureMessage string

	fakeFailHandler := func(message string, callerSkip ...int) {
		failureMessage = message
	}

	BeforeEach(func() {
		failureMessage = ""
	})

	Describe("actual", func() {
		It("should pass through correctly", func() {
			failureMessage = ""
			newActual(true, fakeFailHandler, 1).Should(&deprecatedBeTrueMatcher{})
			Ω(failureMessage).Should(BeEmpty())

			failureMessage = ""
			newActual(false, fakeFailHandler, 1).Should(&deprecatedBeTrueMatcher{})
			Ω(failureMessage).Should(Equal("should be true"))

			failureMessage = ""
			newActual(true, fakeFailHandler, 1).ShouldNot(&deprecatedBeTrueMatcher{})
			Ω(failureMessage).Should(Equal("should not be true"))

			failureMessage = ""
			newActual(false, fakeFailHandler, 1).ShouldNot(&deprecatedBeTrueMatcher{})
			Ω(failureMessage).Should(BeEmpty())

			failureMessage = ""
			newActual(2, fakeFailHandler, 1).ShouldNot(&deprecatedBeTrueMatcher{})
			Ω(failureMessage).Should(Equal("Expected a boolean.  Got: 2"))
		})
	})

	Describe("async actual", func() {
		It("should pass through correctly", func() {
			t := time.Millisecond * 10
			dt := time.Millisecond

			failureMessage = ""
			newAsyncActual(asyncActualTypeEventually, func() bool { return true }, fakeFailHandler, t, dt, 1).Should(&deprecatedBeTrueMatcher{})
			Ω(failureMessage).Should(BeEmpty())

			failureMessage = ""
			newAsyncActual(asyncActualTypeEventually, func() bool { return false }, fakeFailHandler, t, dt, 1).Should(&deprecatedBeTrueMatcher{})
			Ω(failureMessage).Should(ContainSubstring("Timed out"))
			Ω(failureMessage).Should(ContainSubstring("should be true"))

			failureMessage = ""
			newAsyncActual(asyncActualTypeEventually, func() bool { return true }, fakeFailHandler, t, dt, 1).ShouldNot(&deprecatedBeTrueMatcher{})
			Ω(failureMessage).Should(ContainSubstring("Timed out"))
			Ω(failureMessage).Should(ContainSubstring("should not be true"))

			failureMessage = ""
			newAsyncActual(asyncActualTypeEventually, func() bool { return false }, fakeFailHandler, t, dt, 1).ShouldNot(&deprecatedBeTrueMatcher{})
			Ω(failureMessage).Should(BeEmpty())

			failureMessage = ""
			newAsyncActual(asyncActualTypeEventually, func() int { return 2 }, fakeFailHandler, t, dt, 1).ShouldNot(&deprecatedBeTrueMatcher{})
			Ω(failureMessage).Should(ContainSubstring("Timed out"))
			Ω(failureMessage).Should(ContainSubstring("Expected a boolean.  Got: 2"))
		})
	})
})
