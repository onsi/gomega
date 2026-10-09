// untested sections: 2

package matchers

import (
	"fmt"

	"github.com/onsi/gomega/format"
)

type HaveCapMatcher struct {
	Count any
}

func (matcher *HaveCapMatcher) Match(actual any) (success bool, err error) {
	length, ok := capOf(actual)
	if !ok {
		return false, fmt.Errorf("HaveCap matcher expects a array/channel/slice.  Got:\n%s", format.Object(actual, 1))
	}

	if c, ok := matcher.Count.(int); ok {
		return length == c, nil
	}
	if m, ok := matcher.Count.(omegaMatcher); ok {
		return m.Match(length)
	}
	if c, ok := toInt(matcher.Count); ok {
		return length == c, nil
	}

	return false, fmt.Errorf("HaveCap matcher expects an int or GomegaMatcher.  Got:\n%s", format.Object(matcher.Count, 1))
}

func (matcher *HaveCapMatcher) FailureMessage(actual any) (message string) {
	switch m := matcher.Count.(type) {
	case omegaMatcher:
		return format.Message(actual, "to have capacity matching", m)
	default:
		return fmt.Sprintf("Expected\n%s\nto have capacity %v", format.Object(actual, 1), matcher.Count)
	}
}

func (matcher *HaveCapMatcher) NegatedFailureMessage(actual any) (message string) {
	switch m := matcher.Count.(type) {
	case omegaMatcher:
		return format.Message(actual, "not to have capacity matching", m)
	default:
		return fmt.Sprintf("Expected\n%s\nnot to have capacity %v", format.Object(actual, 1), matcher.Count)
	}
}
