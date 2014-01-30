package matchers

import (
	"fmt"
	"reflect"
)

type BeClosedMatcher struct {
}

func (matcher *BeClosedMatcher) Match(actual interface{}) (success bool, message string, err error) {
	if !isChan(actual) {
		return false, "", fmt.Errorf("BeClosed matcher expects a channel.  Got: %s", formatObject(actual))
	}

	channelType := reflect.TypeOf(actual)
	channelValue := reflect.ValueOf(actual)

	var closed bool

	if channelType.ChanDir() == reflect.SendDir {
		return false, "", fmt.Errorf("BeClosed matcher cannot determine if a send-only channel is closed or open.  Got: %s", formatObject(actual))
	} else {
		winnerIndex, _, open := reflect.Select([]reflect.SelectCase{
			reflect.SelectCase{Dir: reflect.SelectRecv, Chan: channelValue},
			reflect.SelectCase{Dir: reflect.SelectDefault},
		})

		if winnerIndex == 0 {
			closed = !open
		} else if winnerIndex == 1 {
			closed = false
		}
	}

	if closed {
		return true, formatMessage(actual, "to be open"), nil
	} else {
		return false, formatMessage(actual, "to be closed"), nil
	}
}
