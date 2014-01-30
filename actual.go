package gomega

import (
	"fmt"
)

type actual struct {
	actualInput interface{}
	fail        OmegaFailHandler
	offset      int
}

func newActual(actualInput interface{}, fail OmegaFailHandler, offset int) *actual {
	return &actual{
		actualInput: actualInput,
		fail:        fail,
		offset:      offset,
	}
}

func (actual *actual) Should(matcher OmegaMatcher, optionalDescription ...interface{}) bool {
	return actual.match(matcher, true, optionalDescription...)
}

func (actual *actual) ShouldNot(matcher OmegaMatcher, optionalDescription ...interface{}) bool {
	return actual.match(matcher, false, optionalDescription...)
}

func (actual *actual) To(matcher OmegaMatcher, optionalDescription ...interface{}) bool {
	return actual.match(matcher, true, optionalDescription...)
}

func (actual *actual) ToNot(matcher OmegaMatcher, optionalDescription ...interface{}) bool {
	return actual.match(matcher, false, optionalDescription...)
}

func (actual *actual) NotTo(matcher OmegaMatcher, optionalDescription ...interface{}) bool {
	return actual.match(matcher, false, optionalDescription...)
}

func (actual *actual) buildDescription(optionalDescription ...interface{}) string {
	switch len(optionalDescription) {
	case 0:
		return ""
	default:
		return fmt.Sprintf(optionalDescription[0].(string), optionalDescription[1:]...) + "\n"
	}
}

func (actual *actual) match(matcher OmegaMatcher, desiredMatch bool, optionalDescription ...interface{}) bool {
	matches, message, err := matcher.Match(actual.actualInput)
	description := actual.buildDescription(optionalDescription...)
	if err != nil {
		actual.fail(description+err.Error(), 2+actual.offset)
		return false
	}
	if matches != desiredMatch {
		actual.fail(description+message, 2+actual.offset)
		return false
	}

	return true
}
