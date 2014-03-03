package matchers

import (
	"fmt"
	"github.com/onsi/gomega/format"
	"reflect"
)

type PanicMatcher struct{}

func (matcher *PanicMatcher) Match(actual interface{}) (success bool, message string, err error) {
	if actual == nil {
		return false, "", fmt.Errorf("PanicMatcher expects a non-nil actual.")
	}
	actualType := reflect.TypeOf(actual)
	if actualType.Kind() != reflect.Func {
		return false, "", fmt.Errorf("PanicMatcher expects a function.  Got:\n%s", format.Object(actual, 1))
	}
	if !(actualType.NumIn() == 0 && actualType.NumOut() == 0) {
		return false, "", fmt.Errorf("PanicMatcher expects a function with no arguments and no return value.  Got:\n%s", format.Object(actual, 1))
	}

	success = false
	message = format.Message(actual, "to panic")
	err = nil

	defer func() {
		if e := recover(); e != nil {
			success = true
			message = format.Message(actual, "not to panic")
			err = nil
		}
	}()

	reflect.ValueOf(actual).Call([]reflect.Value{})

	return
}
