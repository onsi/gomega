package matchers

import (
	"fmt"

	"github.com/onsi/gomega/format"
)

type HaveLenMatcher struct {
	Count any
}

func (matcher *HaveLenMatcher) Match(actual any) (success bool, err error) {
	length, ok := lengthOf(actual)
	if !ok {
		return false, fmt.Errorf("HaveLen matcher expects a string/array/map/channel/slice/iterator.  Got:\n%s", format.Object(actual, 1))
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

	return false, fmt.Errorf("HaveLen matcher expects an int or GomegaMatcher.  Got:\n%s", format.Object(matcher.Count, 1))
}

func (matcher *HaveLenMatcher) FailureMessage(actual any) (message string) {
	switch m := matcher.Count.(type) {
	case omegaMatcher:
		return format.Message(actual, "to have length matching", m)
	default:
		return fmt.Sprintf("Expected\n%s\nto have length %v", format.Object(actual, 1), matcher.Count)
	}
}

func (matcher *HaveLenMatcher) NegatedFailureMessage(actual any) (message string) {
	switch m := matcher.Count.(type) {
	case omegaMatcher:
		return format.Message(actual, "not to have length matching", m)
	default:
		return fmt.Sprintf("Expected\n%s\nnot to have length %v", format.Object(actual, 1), matcher.Count)
	}
}
