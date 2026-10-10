package matchers

import (
	"fmt"

	"github.com/onsi/gomega/format"
)

type HaveLenMatcher struct {
	Count any

	length int
}

func (matcher *HaveLenMatcher) Match(actual any) (success bool, err error) {
	length, ok := lengthOf(actual)
	if !ok {
		return false, fmt.Errorf("HaveLen matcher expects a string/array/map/channel/slice/iterator.  Got:\n%s", format.Object(actual, 1))
	}
	matcher.length = length

	if m, ok := matcher.Count.(omegaMatcher); ok {
		return m.Match(length)
	}
	if count, ok := toInt(matcher.Count); ok {
		return length == count, nil
	}

	return false, fmt.Errorf("HaveLen matcher expects an integer or a GomegaMatcher.  Got:\n%s", format.Object(matcher.Count, 1))
}

func (matcher *HaveLenMatcher) FailureMessage(actual any) (message string) {
	if m, ok := matcher.Count.(omegaMatcher); ok {
		return fmt.Sprintf("Length %d of\n%s\nfailed to satisfy matcher:\n%s", matcher.length, format.Object(actual, 1), format.IndentString(m.FailureMessage(matcher.length), 1))
	}
	count, _ := toInt(matcher.Count)
	return fmt.Sprintf("Expected\n%s\nto have length %d", format.Object(actual, 1), count)
}

func (matcher *HaveLenMatcher) NegatedFailureMessage(actual any) (message string) {
	if m, ok := matcher.Count.(omegaMatcher); ok {
		return fmt.Sprintf("Length %d of\n%s\nsatisfied matcher, but should not have:\n%s", matcher.length, format.Object(actual, 1), format.IndentString(m.NegatedFailureMessage(matcher.length), 1))
	}
	count, _ := toInt(matcher.Count)
	return fmt.Sprintf("Expected\n%s\nnot to have length %d", format.Object(actual, 1), count)
}
