// untested sections: 2

package matchers

import (
	"fmt"

	"github.com/onsi/gomega/format"
)

type HaveCapMatcher struct {
	Count any

	capacity int
}

func (matcher *HaveCapMatcher) Match(actual any) (success bool, err error) {
	capacity, ok := capOf(actual)
	if !ok {
		return false, fmt.Errorf("HaveCap matcher expects a array/channel/slice.  Got:\n%s", format.Object(actual, 1))
	}
	matcher.capacity = capacity

	if m, ok := matcher.Count.(omegaMatcher); ok {
		return m.Match(capacity)
	}
	if count, ok := toInt(matcher.Count); ok {
		return capacity == count, nil
	}

	return false, fmt.Errorf("HaveCap matcher expects an integer or a GomegaMatcher.  Got:\n%s", format.Object(matcher.Count, 1))
}

func (matcher *HaveCapMatcher) FailureMessage(actual any) (message string) {
	if m, ok := matcher.Count.(omegaMatcher); ok {
		return fmt.Sprintf("Capacity %d of\n%s\nfailed to satisfy matcher:\n%s", matcher.capacity, format.Object(actual, 1), format.IndentString(m.FailureMessage(matcher.capacity), 1))
	}
	count, _ := toInt(matcher.Count)
	return fmt.Sprintf("Expected\n%s\nto have capacity %d", format.Object(actual, 1), count)
}

func (matcher *HaveCapMatcher) NegatedFailureMessage(actual any) (message string) {
	if m, ok := matcher.Count.(omegaMatcher); ok {
		return fmt.Sprintf("Capacity %d of\n%s\nsatisfied matcher, but should not have:\n%s", matcher.capacity, format.Object(actual, 1), format.IndentString(m.NegatedFailureMessage(matcher.capacity), 1))
	}
	count, _ := toInt(matcher.Count)
	return fmt.Sprintf("Expected\n%s\nnot to have capacity %d", format.Object(actual, 1), count)
}
