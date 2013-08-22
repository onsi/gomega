package matchers

import "fmt"

type EqualMatcher struct {
	Expected interface{}
}

func (matcher *EqualMatcher) Match(actual interface{}) (success bool, message string, err error) {
	if actual == matcher.Expected {
		return true, formatMessage(actual, "not to equal", matcher.Expected), nil
	} else {
		return false, formatMessage(actual, "to equal", matcher.Expected), nil
	}
}

type TrueMatcher struct {
}

func (matcher *TrueMatcher) Match(actual interface{}) (success bool, message string, err error) {
	if !isBool(actual) {
		return false, "", fmt.Errorf("Expected a boolean, got%s", formatObject(actual))
	}
	if actual == true {
		return true, formatMessage(actual, "not to be true"), nil
	} else {
		return false, formatMessage(actual, "to be true"), nil
	}
}

type FalseMatcher struct {
}

func (matcher *FalseMatcher) Match(actual interface{}) (success bool, message string, err error) {
	if !isBool(actual) {
		return false, "", fmt.Errorf("Expected a boolean, got%s", formatObject(actual))
	}
	if actual == false {
		return true, formatMessage(actual, "not to be false"), nil
	} else {
		return false, formatMessage(actual, "to be false"), nil
	}
}

type HaveOccuredMatcher struct {
}

func (matcher *HaveOccuredMatcher) Match(actual interface{}) (success bool, message string, err error) {
	if actual == nil {
		return false, formatMessage(actual, "to have occured"), nil
	} else {
		if isError(actual) {
			return true, fmt.Sprintf("Expected error:%s\n\tMessage: %s\n%s", formatObject(actual), actual.(error).Error(), "not to have occured"), nil
		} else {
			return false, "", fmt.Errorf("Expected an error, got%s", formatObject(actual))
		}
	}
}
