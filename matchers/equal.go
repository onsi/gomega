package matchers

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type EqualMatcher struct {
	Expected interface{}
}

func (matcher *EqualMatcher) Match(actual interface{}) (success bool, message string, err error) {
	if actual == nil && matcher.Expected == nil {
		return false, "", fmt.Errorf("Refusing to compare <nil> to <nil>.")
	}
	if reflect.DeepEqual(actual, matcher.Expected) {
		return true, formatMessage(actual, "not to equal", matcher.Expected), nil
	} else {
		return false, formatMessage(actual, "to equal", matcher.Expected), nil
	}
	return
}

type BeNilMatcher struct {
}

func (matcher *BeNilMatcher) Match(actual interface{}) (success bool, message string, err error) {
	if actual == nil {
		return true, formatMessage(actual, "not to be nil"), nil
	} else {
		return false, formatMessage(actual, "to be nil"), nil
	}
}

type BeTrueMatcher struct {
}

func (matcher *BeTrueMatcher) Match(actual interface{}) (success bool, message string, err error) {
	if !isBool(actual) {
		return false, "", fmt.Errorf("Expected a boolean, got%s", formatObject(actual))
	}
	if actual == true {
		return true, formatMessage(actual, "not to be true"), nil
	} else {
		return false, formatMessage(actual, "to be true"), nil
	}
}

type BeFalseMatcher struct {
}

func (matcher *BeFalseMatcher) Match(actual interface{}) (success bool, message string, err error) {
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

type ContainSubstringMatcher struct {
	Substr string
}

func (matcher *ContainSubstringMatcher) Match(actual interface{}) (success bool, message string, err error) {
	actualString, ok := toString(actual)
	if ok {
		match := strings.Contains(actualString, matcher.Substr)
		if match {
			return true, formatMessage(actual, "not to contain substring", matcher.Substr), nil
		} else {
			return false, formatMessage(actual, "to contain substring", matcher.Substr), nil
		}
	} else {
		return false, "", fmt.Errorf("ContainSubstring matcher requires a string or stringer.  Got:%s", formatObject(actual))
	}
}

type BeEmptyMatcher struct {
}

func (matcher *BeEmptyMatcher) Match(actual interface{}) (success bool, message string, err error) {
	length, ok := lengthOf(actual)
	if ok {
		if length == 0 {
			return true, formatMessage(actual, "not to be empty"), nil
		} else {
			return false, formatMessage(actual, "to be empty"), nil
		}
	} else {
		return false, "", fmt.Errorf("BeEmpty matcher expects a string/array/map/channel/slice.  Got:%s", formatObject(actual))
	}
}

type HaveLenMatcher struct {
	Count int
}

func (matcher *HaveLenMatcher) Match(actual interface{}) (success bool, message string, err error) {
	length, ok := lengthOf(actual)
	if ok {
		if length == matcher.Count {
			return true, fmt.Sprintf("Expected%s\n (length: %d) not to have length %d", formatObject(actual), length, matcher.Count), nil
		} else {
			return false, fmt.Sprintf("Expected%s\n (length: %d) to have length %d", formatObject(actual), length, matcher.Count), nil
		}
	} else {
		return false, "", fmt.Errorf("BeEmpty matcher expects a string/array/map/channel/slice.  Got:%s", formatObject(actual))
	}
}

type BeZeroMatcher struct {
}

func (matcher *BeZeroMatcher) Match(actual interface{}) (success bool, message string, err error) {
	zeroValue := reflect.Zero(reflect.TypeOf(actual)).Interface()
	if reflect.DeepEqual(zeroValue, actual) {
		return true, formatMessage(actual, "not to be zero-valued"), nil
	} else {
		return false, formatMessage(actual, "to be zero-valued"), nil
	}
}

type ContainElementMatcher struct {
	Element interface{}
}

func (matcher *ContainElementMatcher) Match(actual interface{}) (success bool, message string, err error) {
	if !isArrayOrSlice(actual) {
		return false, "", fmt.Errorf("ContainElement matcher expects an array/slice/string.  Got:%s", formatObject(actual))
	}
	value := reflect.ValueOf(actual)
	for i := 0; i < value.Len(); i++ {
		if reflect.DeepEqual(value.Index(i).Interface(), matcher.Element) {
			return true, formatMessage(actual, "not to contain element", matcher.Element), nil
		}
	}
	return false, formatMessage(actual, "to contain element", matcher.Element), nil
}
