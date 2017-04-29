package matchers

import (
	"encoding/xml"
	"fmt"
	"reflect"

	"github.com/onsi/gomega/format"
)

type MatchXMLMatcher struct {
	XMLToMatch interface{}
}

func (matcher *MatchXMLMatcher) Match(actual interface{}) (success bool, err error) {
	actualString, expectedString, err := matcher.formattedPrint(actual)
	if err != nil {
		return false, err
	}

	aval := &xmlNode{}
	eval := &xmlNode{}

	if err := xml.Unmarshal([]byte(actualString), aval); err != nil {
		return false, fmt.Errorf("Actual '%s' should be valid XML, but it is not.\nUnderlying error:%s", actualString, err)
	}
	if err := xml.Unmarshal([]byte(expectedString), eval); err != nil {
		return false, fmt.Errorf("Expected '%s' should be valid XML, but it is not.\nUnderlying error:%s", expectedString, err)
	}

	aval.Clean()
	eval.Clean()

	return reflect.DeepEqual(aval, eval), nil
}

func (matcher *MatchXMLMatcher) FailureMessage(actual interface{}) (message string) {
	actualString, expectedString, _ := matcher.formattedPrint(actual)
	return fmt.Sprintf("Expected\n%s\nto match XML of\n%s", actualString, expectedString)
}

func (matcher *MatchXMLMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	actualString, expectedString, _ := matcher.formattedPrint(actual)
	return fmt.Sprintf("Expected\n%s\nnot to match XML of\n%s", actualString, expectedString)
}

func (matcher *MatchXMLMatcher) formattedPrint(actual interface{}) (actualString, expectedString string, err error) {
	var ok bool
	actualString, ok = toString(actual)
	if !ok {
		return "", "", fmt.Errorf("MatchXMLMatcher matcher requires a string, stringer, or []byte.  Got actual:\n%s", format.Object(actual, 1))
	}
	expectedString, ok = toString(matcher.XMLToMatch)
	if !ok {
		return "", "", fmt.Errorf("MatchXMLMatcher matcher requires a string, stringer, or []byte.  Got expected:\n%s", format.Object(matcher.XMLToMatch, 1))
	}
	return actualString, expectedString, nil
}
