package matchers

import (
	"fmt"
	"reflect"
)

func isBool(a interface{}) bool {
	return reflect.TypeOf(a).Kind() == reflect.Bool
}

func isNumber(a interface{}) bool {
	kind := reflect.TypeOf(a).Kind()
	return reflect.Int <= kind && kind <= reflect.Float64
}

func isError(a interface{}) bool {
	_, ok := a.(error)
	return ok
}

func isArrayOrSlice(a interface{}) bool {
	if a == nil {
		return false
	}
	switch reflect.TypeOf(a).Kind() {
	case reflect.Array, reflect.Slice:
		return true
	default:
		return false
	}
}

func toString(a interface{}) (string, bool) {
	aString, isString := a.(string)
	if isString {
		return aString, true
	}
	aStringer, isStringer := a.(fmt.Stringer)
	if isStringer {
		return aStringer.String(), true
	}

	return "", false
}

func lengthOf(a interface{}) (int, bool) {
	if a == nil {
		return 0, false
	}
	switch reflect.TypeOf(a).Kind() {
	case reflect.Map, reflect.Array, reflect.String, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(a).Len(), true
	default:
		return 0, false
	}
}
