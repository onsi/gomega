package gomega

import (
	"fmt"
	"reflect"
)

type actual struct {
	actualInput interface{}
	fail        OmegaFailHandler
	offset      int
	extra       []interface{}
}

func newActual(actualInput interface{}, fail OmegaFailHandler, offset int, extra ...interface{}) *actual {
	return &actual{
		actualInput: actualInput,
		fail:        fail,
		offset:      offset,
		extra:       extra,
	}
}

func (actual *actual) Should(matcher OmegaMatcher, optionalDescription ...interface{}) bool {
	return actual.vetExtras(optionalDescription...) && actual.match(matcher, true, optionalDescription...)
}

func (actual *actual) ShouldNot(matcher OmegaMatcher, optionalDescription ...interface{}) bool {
	return actual.vetExtras(optionalDescription...) && actual.match(matcher, false, optionalDescription...)
}

func (actual *actual) To(matcher OmegaMatcher, optionalDescription ...interface{}) bool {
	return actual.vetExtras(optionalDescription...) && actual.match(matcher, true, optionalDescription...)
}

func (actual *actual) ToNot(matcher OmegaMatcher, optionalDescription ...interface{}) bool {
	return actual.vetExtras(optionalDescription...) && actual.match(matcher, false, optionalDescription...)
}

func (actual *actual) NotTo(matcher OmegaMatcher, optionalDescription ...interface{}) bool {
	return actual.vetExtras(optionalDescription...) && actual.match(matcher, false, optionalDescription...)
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

func (actual *actual) vetExtras(optionalDescription ...interface{}) bool {
	success, message := vetExtras(actual.extra)
	if success {
		return true
	}

	description := actual.buildDescription(optionalDescription...)
	actual.fail(description+message, 2+actual.offset)
	return false
}

func vetExtras(extras []interface{}) (bool, string) {
	for i, extra := range extras {
		if extra != nil {
			zeroValue := reflect.Zero(reflect.TypeOf(extra)).Interface()
			if !reflect.DeepEqual(zeroValue, extra) {
				message := fmt.Sprintf("Unexpected non-nil/non-zero extra argument at index %d:\n\t<%T>: %#v", i+1, extra, extra)
				return false, message
			}
		}
	}
	return true, ""
}
