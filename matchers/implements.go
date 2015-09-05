package matchers

import (
	"fmt"
	"github.com/onsi/gomega/format"
	"reflect"
)

type ImplementsMatcher struct {
	InterfaceType reflect.Type
}

func (m *ImplementsMatcher) Match(actual interface{}) (success bool, err error) {
	if actual == nil {
		return false, nil
	}
	a := reflect.ValueOf(actual)
	return a.Type().Implements(m.InterfaceType), nil
}

func (m *ImplementsMatcher) FailureMessage(actual interface{}) (message string) {
	return format.Message(actual, fmt.Sprintf("to implement %s", m.InterfaceType))
}

func (m *ImplementsMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, fmt.Sprintf("to not implement %s", m.InterfaceType))
}
