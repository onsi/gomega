package matchers

import (
	"fmt"
	"regexp"
)

type MatchRegexpMatcher struct {
	Regexp string
	Args   []interface{}
}

func (matcher *MatchRegexpMatcher) Match(actual interface{}) (success bool, message string, err error) {
	actualString, ok := toString(actual)
	if ok {
		re := matcher.Regexp
		if len(matcher.Args) > 0 {
			re = fmt.Sprintf(matcher.Regexp, matcher.Args...)
		}

		match, err := regexp.Match(re, []byte(actualString))
		if err != nil {
			return false, "", fmt.Errorf("RegExp match failed to compile with error:\n\t%s", err.Error())
		}
		if match {
			return true, formatMessage(actual, "not to match regular expression", re), nil
		} else {
			return false, formatMessage(actual, "to match regular expression", re), nil
		}
	} else {
		return false, "", fmt.Errorf("RegExp matcher requires a string or stringer.  Got:%s", formatObject(actual))
	}
}
