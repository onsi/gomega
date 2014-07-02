package matchers

import (
	"fmt"
	"github.com/onsi/gomega/format"
	"time"
)

type BeTemporallyMatcher struct {
	Comparator string
	CompareTo  []interface{}
}

func (matcher *BeTemporallyMatcher) FailureMessage(actual interface{}) (message string) {
	return format.Message(actual, fmt.Sprintf("to be %s", matcher.Comparator), matcher.CompareTo[0])
}

func (matcher *BeTemporallyMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, fmt.Sprintf("not to be %s", matcher.Comparator), matcher.CompareTo[0])
}

func (matcher *BeTemporallyMatcher) Match(actual interface{}) (bool, error) {
	// predicate to test for time.Time type
	isTime := func(t interface{}) bool {
		_, ok := t.(time.Time)
		return ok
	}

	if len(matcher.CompareTo) == 0 || len(matcher.CompareTo) > 2 {
		return false, fmt.Errorf("BeTemporally requires 1 or 2 CompareTo arguments.  Got:\n%s", format.Object(matcher.CompareTo, 1))
	}
	if !isTime(actual) {
		return false, fmt.Errorf("Expected a time.Time.  Got:\n%s", format.Object(actual, 1))
	}
	if !isTime(matcher.CompareTo[0]) {
		return false, fmt.Errorf("Expected a time.Time.  Got:\n%s", format.Object(matcher.CompareTo[0], 1))
	}

	switch matcher.Comparator {
	case "==", "~", ">", ">=", "<", "<=":
	default:
		return false, fmt.Errorf("Unknown comparator: %s", matcher.Comparator)
	}

	var secondOperand = time.Millisecond
	if len(matcher.CompareTo) == 2 {
		var ok bool
		secondOperand, ok = matcher.CompareTo[1].(time.Duration)
		if !ok {
			return false, fmt.Errorf("Expected a time.Duration.  Got:\n%s", format.Object(matcher.CompareTo[1], 1))
		}
	}

	return matcher.matchTimes(actual.(time.Time), matcher.CompareTo[0].(time.Time), secondOperand), nil
}

func (matcher *BeTemporallyMatcher) matchTimes(actual, compareTo time.Time, threshold time.Duration) (success bool) {
	switch matcher.Comparator {
	case "==":
		return actual.Equal(compareTo)
	case "~":
		diff := actual.Sub(compareTo)
		return -threshold <= diff && diff <= threshold
	case ">":
		return actual.After(compareTo)
	case ">=":
		return !actual.Before(compareTo)
	case "<":
		return actual.Before(compareTo)
	case "<=":
		return !actual.After(compareTo)
	}
	return false
}
