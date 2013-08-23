package matchers

import (
	"fmt"
	"regexp"
)

type MatchRegexpMatcher struct {
	Regexp string
}

func (matcher *MatchRegexpMatcher) Match(actual interface{}) (success bool, message string, err error) {
	actualString, ok := toString(actual)
	if ok {
		match, err := regexp.Match(matcher.Regexp, []byte(actualString))
		if err != nil {
			return false, "", fmt.Errorf("RegExp match failed to compile with error:\n\t%s", err.Error())
		}
		if match {
			return true, formatMessage(actual, "not to match regular expression", matcher.Regexp), nil
		} else {
			return false, formatMessage(actual, "to match regular expression", matcher.Regexp), nil
		}
	} else {
		return false, "", fmt.Errorf("RegExp matcher requires a string or stringer.  Got:%s", formatObject(actual))
	}
}
