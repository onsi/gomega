package gomega

import (
	"fmt"
)

type actual struct {
	actualInput interface{}
	fail        OmegaFailHandler
}

func newActual(actualInput interface{}, fail OmegaFailHandler) *actual {
	return &actual{
		actualInput: actualInput,
		fail:        fail,
	}
}

func (actual *actual) Should(matcher OmegaMatcher, optionalDescription ...interface{}) {
	actual.runPositiveMatch(matcher, optionalDescription...)
}

func (actual *actual) ShouldNot(matcher OmegaMatcher, optionalDescription ...interface{}) {
	actual.runNegativeMatch(matcher, optionalDescription...)
}

func (actual *actual) To(matcher OmegaMatcher, optionalDescription ...interface{}) {
	actual.runPositiveMatch(matcher, optionalDescription...)
}

func (actual *actual) ToNot(matcher OmegaMatcher, optionalDescription ...interface{}) {
	actual.runNegativeMatch(matcher, optionalDescription...)
}

func (actual *actual) NotTo(matcher OmegaMatcher, optionalDescription ...interface{}) {
	actual.runNegativeMatch(matcher, optionalDescription...)
}

func (actual *actual) buildDescription(optionalDescription ...interface{}) string {
	switch len(optionalDescription) {
	case 0:
		return ""
	default:
		return fmt.Sprintf(optionalDescription[0].(string), optionalDescription[1:]...) + "\n"
	}
}

func (actual *actual) runPositiveMatch(matcher OmegaMatcher, optionalDescription ...interface{}) {
	matches, message, err := matcher.Match(actual.actualInput)
	description := actual.buildDescription(optionalDescription...)
	if err != nil {
		actual.fail(description+err.Error(), 2)
		return
	}
	if !matches {
		actual.fail(description+message, 2)
		return
	}
}

func (actual *actual) runNegativeMatch(matcher OmegaMatcher, optionalDescription ...interface{}) {
	matches, message, err := matcher.Match(actual.actualInput)
	description := actual.buildDescription(optionalDescription...)
	if err != nil {
		actual.fail(description+err.Error(), 2)
		return
	}
	if matches {
		actual.fail(description+message, 2)
		return
	}
}
