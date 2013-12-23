package matchers

import "fmt"

func formatMessage(actual interface{}, message string, expected ...interface{}) string {
	if len(expected) == 0 {
		return fmt.Sprintf("Expected%s\n%s", formatObject(actual), message)
	} else {
		return fmt.Sprintf("Expected%s\n%s%s", formatObject(actual), message, formatObject(expected[0]))
	}
}

func formatObject(object interface{}) string {
	hasLength := false
	length := 0
	if !isString(object) {
		length, hasLength = lengthOf(object)
	}

	if hasLength {
		return fmt.Sprintf("\n\t<%T> of length %d: %#v", object, length, object)
	} else {
		return fmt.Sprintf("\n\t<%T>: %#v", object, object)
	}
}
