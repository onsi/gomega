package matchers

type BeNilMatcher struct {
}

func (matcher *BeNilMatcher) Match(actual interface{}) (success bool, message string, err error) {
	if isNil(actual) {
		return true, formatMessage(actual, "not to be nil"), nil
	} else {
		return false, formatMessage(actual, "to be nil"), nil
	}
}
