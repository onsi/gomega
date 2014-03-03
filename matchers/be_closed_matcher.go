package matchers

import (
	"fmt"
	"github.com/onsi/gomega/format"
	"reflect"
)

type BeClosedMatcher struct {
}

func (matcher *BeClosedMatcher) Match(actual interface{}) (success bool, message string, err error) {
	if !isChan(actual) {
		return false, "", fmt.Errorf("BeClosed matcher expects a channel.  Got:\n%s", format.Object(actual, 1))
	}

	channelType := reflect.TypeOf(actual)
	channelValue := reflect.ValueOf(actual)

	var closed bool

	if channelType.ChanDir() == reflect.SendDir {
		return false, "", fmt.Errorf("BeClosed matcher cannot determine if a send-only channel is closed or open.  Got:\n%s", format.Object(actual, 1))
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
		return true, format.Message(actual, "to be open"), nil
	} else {
		return false, format.Message(actual, "to be closed"), nil
	}
}
