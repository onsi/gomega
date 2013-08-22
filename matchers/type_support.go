package matchers

import "reflect"

func isBool(a interface{}) bool {
	return reflect.TypeOf(a).Kind() == reflect.Bool
}

func isError(a interface{}) bool {
	_, ok := a.(error)
	return ok
}
