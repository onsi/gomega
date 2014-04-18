package gbytes

import (
	"fmt"
	"regexp"

	"github.com/onsi/gomega/format"
)

func Say(expected string, args ...interface{}) *sayMatcher {
	formattedRegexp := expected
	if len(args) > 0 {
		formattedRegexp = fmt.Sprintf(expected, args...)
	}
	return &sayMatcher{
		re: regexp.MustCompile(formattedRegexp),
	}
}

type sayMatcher struct {
	re              *regexp.Regexp
	receivedSayings []byte
}

func (m *sayMatcher) Match(actual interface{}) (success bool, err error) {
	buffer, ok := actual.(*Buffer)
	if !ok {
		return false, fmt.Errorf("Say must be passed a gbytes Buffer.  Got:\n%s", format.Object(actual, 1))
	}

	didSay, sayings := buffer.didSay(m.re)
	m.receivedSayings = sayings

	return didSay, nil
}

func (m *sayMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf(
		"Got stuck at:\n%s\nWaiting for:\n%s",
		format.IndentString(string(m.receivedSayings), 1),
		format.IndentString(m.re.String(), 1),
	)
}

func (m *sayMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf(
		"Saw:\n%s\nWhich matches the unexpected:\n%s",
		format.IndentString(string(m.receivedSayings), 1),
		format.IndentString(m.re.String(), 1),
	)
}
