package matchers

import (
	"fmt"
	"reflect"
	"github.com/onsi/gomega/format"
)

type ReceiveMatcher struct {
	Arg interface{}
}

func (matcher *ReceiveMatcher) Match(actual interface{}) (success bool, message string, err error) {
	if !isChan(actual) {
		return false, "", fmt.Errorf("ReceiveMatcher expects a channel.  Got: %s", format.Object(actual))
	}

	channelType := reflect.TypeOf(actual)
	channelValue := reflect.ValueOf(actual)

	if channelType.ChanDir() == reflect.SendDir {
		return false, "", fmt.Errorf("ReceiveMatcher matcher cannot be passed a send-only channel.  Got: %s", format.Object(actual))
	}

	if matcher.Arg != nil {
		argType := reflect.TypeOf(matcher.Arg)
		if argType.Kind() != reflect.Ptr {
			return false, "", fmt.Errorf("Cannot assign a value from the channel:\n\t%s\nTo:\n\t%s\nYou need to pass a pointer!", format.Object(actual), format.Object(matcher.Arg))
		}

		assignable := channelType.Elem().AssignableTo(argType.Elem())
		if !assignable {
			return false, "", fmt.Errorf("Cannot assign a value from the channel:\n\t%s\nTo:\n\t%s", format.Object(actual), format.Object(matcher.Arg))
		}
	}

	var closed bool
	var didReceive bool

	winnerIndex, value, open := reflect.Select([]reflect.SelectCase{
		reflect.SelectCase{Dir: reflect.SelectRecv, Chan: channelValue},
		reflect.SelectCase{Dir: reflect.SelectDefault},
	})

	if winnerIndex == 0 {
		closed = !open
		didReceive = open
	} else if winnerIndex == 1 {
		closed = false
		didReceive = false
	}

	if closed {
		return false, "", fmt.Errorf("ReceiveMatcher was given a closed channel: %s", format.Object(actual))
	}

	if didReceive {
		if matcher.Arg != nil {
			outValue := reflect.ValueOf(matcher.Arg)
			reflect.Indirect(outValue).Set(value)
		}
		return true, format.Message(actual, "not to receive anything"), nil
	} else {
		return false, format.Message(actual, "to receive something"), nil
	}
}
