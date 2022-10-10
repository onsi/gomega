package internal

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/onsi/gomega/types"
)

type AsyncAssertionType uint

const (
	AsyncAssertionTypeEventually AsyncAssertionType = iota
	AsyncAssertionTypeConsistently
)

func (at AsyncAssertionType) String() string {
	switch at {
	case AsyncAssertionTypeEventually:
		return "Eventually"
	case AsyncAssertionTypeConsistently:
		return "Consistently"
	}
	return "INVALID ASYNC ASSERTION TYPE"
}

type AsyncAssertion struct {
	asyncType AsyncAssertionType

	actualIsFunc bool
	actual       interface{}

	timeoutInterval time.Duration
	pollingInterval time.Duration
	ctx             context.Context
	offset          int
	g               *Gomega
}

func NewAsyncAssertion(asyncType AsyncAssertionType, actualInput interface{}, g *Gomega, timeoutInterval time.Duration, pollingInterval time.Duration, ctx context.Context, offset int) *AsyncAssertion {
	out := &AsyncAssertion{
		asyncType:       asyncType,
		timeoutInterval: timeoutInterval,
		pollingInterval: pollingInterval,
		offset:          offset,
		ctx:             ctx,
		g:               g,
	}

	out.actual = actualInput
	if actualInput != nil && reflect.TypeOf(actualInput).Kind() == reflect.Func {
		out.actualIsFunc = true
	}

	return out
}

func (assertion *AsyncAssertion) WithOffset(offset int) types.AsyncAssertion {
	assertion.offset = offset
	return assertion
}

func (assertion *AsyncAssertion) WithTimeout(interval time.Duration) types.AsyncAssertion {
	assertion.timeoutInterval = interval
	return assertion
}

func (assertion *AsyncAssertion) WithPolling(interval time.Duration) types.AsyncAssertion {
	assertion.pollingInterval = interval
	return assertion
}

func (assertion *AsyncAssertion) Within(timeout time.Duration) types.AsyncAssertion {
	assertion.timeoutInterval = timeout
	return assertion
}

func (assertion *AsyncAssertion) ProbeEvery(interval time.Duration) types.AsyncAssertion {
	assertion.pollingInterval = interval
	return assertion
}

func (assertion *AsyncAssertion) WithContext(ctx context.Context) types.AsyncAssertion {
	assertion.ctx = ctx
	return assertion
}

func (assertion *AsyncAssertion) Should(matcher types.GomegaMatcher, optionalDescription ...interface{}) bool {
	assertion.g.THelper()
	vetOptionalDescription("Asynchronous assertion", optionalDescription...)
	return assertion.match(matcher, true, optionalDescription...)
}

func (assertion *AsyncAssertion) ShouldNot(matcher types.GomegaMatcher, optionalDescription ...interface{}) bool {
	assertion.g.THelper()
	vetOptionalDescription("Asynchronous assertion", optionalDescription...)
	return assertion.match(matcher, false, optionalDescription...)
}

func (assertion *AsyncAssertion) buildDescription(optionalDescription ...interface{}) string {
	switch len(optionalDescription) {
	case 0:
		return ""
	case 1:
		if describe, ok := optionalDescription[0].(func() string); ok {
			return describe() + "\n"
		}
	}
	return fmt.Sprintf(optionalDescription[0].(string), optionalDescription[1:]...) + "\n"
}

func (assertion *AsyncAssertion) processReturnValues(values []reflect.Value) (interface{}, error) {
	if len(values) == 0 {
		return nil, fmt.Errorf("No values were returned by the function passed to Gomega")
	}
	actual := values[0].Interface()
	for i, extraValue := range values[1:] {
		extra := extraValue.Interface()
		if extra == nil {
			continue
		}
		zero := reflect.Zero(extraValue.Type()).Interface()
		if reflect.DeepEqual(extra, zero) {
			continue
		}
		return actual, fmt.Errorf("Unexpected non-nil/non-zero argument at index %d:\n\t<%T>: %#v", i+1, extra, extra)
	}
	return actual, nil
}

var gomegaType = reflect.TypeOf((*types.Gomega)(nil)).Elem()
var contextType = reflect.TypeOf(new(context.Context)).Elem()

func (assertion *AsyncAssertion) invalidFunctionError(t reflect.Type) error {
	return fmt.Errorf(`The function passed to %s had an invalid signature of %s.  Functions passed to %s must either:

	(a) have return values or
	(b) take a Gomega interface as their first argument and use that Gomega instance to make assertions.

You can learn more at https://onsi.github.io/gomega/#eventually
`, assertion.asyncType, t, assertion.asyncType)
}

func (assertion *AsyncAssertion) noConfiguredContextForFunctionError(t reflect.Type) error {
	return fmt.Errorf(`The function passed to %s requested a context.Context, but no context has been provided to %s.  Please pass one in using %s().WithContext().

You can learn more at https://onsi.github.io/gomega/#eventually
`, assertion.asyncType, t, assertion.asyncType)
}

func (assertion *AsyncAssertion) buildActualPoller() (func() (interface{}, error), error) {
	if !assertion.actualIsFunc {
		return func() (interface{}, error) { return assertion.actual, nil }, nil
	}
	actualValue := reflect.ValueOf(assertion.actual)
	actualType := reflect.TypeOf(assertion.actual)
	numIn, numOut := actualType.NumIn(), actualType.NumOut()

	if numIn == 0 && numOut == 0 {
		return nil, assertion.invalidFunctionError(actualType)
	} else if numIn == 0 {
		return func() (interface{}, error) { return assertion.processReturnValues(actualValue.Call([]reflect.Value{})) }, nil
	}
	takesGomega, takesContext := actualType.In(0).Implements(gomegaType), actualType.In(0).Implements(contextType)
	if takesGomega && numIn > 1 && actualType.In(1).Implements(contextType) {
		takesContext = true
	}
	if !takesGomega && numOut == 0 {
		return nil, assertion.invalidFunctionError(actualType)
	}
	if takesContext && assertion.ctx == nil {
		return nil, assertion.noConfiguredContextForFunctionError(actualType)
	}
	remainingIn := numIn
	if takesGomega {
		remainingIn -= 1
	}
	if takesContext {
		remainingIn -= 1
	}
	if remainingIn > 0 {
		return nil, assertion.invalidFunctionError(actualType)
	}

	var assertionFailure error
	inValues := []reflect.Value{}
	if takesGomega {
		inValues = append(inValues, reflect.ValueOf(NewGomega(assertion.g.DurationBundle).ConfigureWithFailHandler(func(message string, callerSkip ...int) {
			skip := 0
			if len(callerSkip) > 0 {
				skip = callerSkip[0]
			}
			_, file, line, _ := runtime.Caller(skip + 1)
			assertionFailure = fmt.Errorf("Assertion in callback at %s:%d failed:\n%s", file, line, message)
			panic("stop execution")
		})))
	}
	if takesContext {
		inValues = append(inValues, reflect.ValueOf(assertion.ctx))
	}

	return func() (actual interface{}, err error) {
		var values []reflect.Value
		assertionFailure = nil
		defer func() {
			if numOut == 0 {
				actual = assertionFailure
			} else {
				actual, err = assertion.processReturnValues(values)
				if assertionFailure != nil {
					err = assertionFailure
				}
			}
			if e := recover(); e != nil && assertionFailure == nil {
				panic(e)
			}
		}()
		values = actualValue.Call(inValues)
		return
	}, nil
}

func (assertion *AsyncAssertion) matcherMayChange(matcher types.GomegaMatcher, value interface{}) bool {
	if assertion.actualIsFunc {
		return true
	}
	return types.MatchMayChangeInTheFuture(matcher, value)
}

type contextWithAttachProgressReporter interface {
	AttachProgressReporter(func() string) func()
}

func (assertion *AsyncAssertion) match(matcher types.GomegaMatcher, desiredMatch bool, optionalDescription ...interface{}) bool {
	timer := time.Now()
	timeout := time.After(assertion.timeoutInterval)
	lock := sync.Mutex{}

	var matches bool
	var err error
	mayChange := true

	assertion.g.THelper()

	pollActual, err := assertion.buildActualPoller()
	if err != nil {
		assertion.g.Fail(err.Error(), 2+assertion.offset)
		return false
	}

	value, err := pollActual()
	if err == nil {
		mayChange = assertion.matcherMayChange(matcher, value)
		matches, err = matcher.Match(value)
	}

	messageGenerator := func() string {
		// can be called out of band by Ginkgo if the user requests a progress report
		lock.Lock()
		defer lock.Unlock()
		errMsg := ""
		message := ""
		if err != nil {
			errMsg = "Error: " + err.Error()
		} else {
			if desiredMatch {
				message = matcher.FailureMessage(value)
			} else {
				message = matcher.NegatedFailureMessage(value)
			}
		}
		description := assertion.buildDescription(optionalDescription...)
		return fmt.Sprintf("%s%s%s", description, message, errMsg)
	}

	fail := func(preamble string) {
		assertion.g.THelper()
		assertion.g.Fail(fmt.Sprintf("%s after %.3fs.\n%s", preamble, time.Since(timer).Seconds(), messageGenerator()), 3+assertion.offset)
	}

	var contextDone <-chan struct{}
	if assertion.ctx != nil {
		contextDone = assertion.ctx.Done()
		if v, ok := assertion.ctx.Value("GINKGO_SPEC_CONTEXT").(contextWithAttachProgressReporter); ok {
			detach := v.AttachProgressReporter(messageGenerator)
			defer detach()
		}
	}

	if assertion.asyncType == AsyncAssertionTypeEventually {
		for {
			if err == nil && matches == desiredMatch {
				return true
			}

			if !mayChange {
				fail("No future change is possible.  Bailing out early")
				return false
			}

			select {
			case <-time.After(assertion.pollingInterval):
				v, e := pollActual()
				lock.Lock()
				value, err = v, e
				lock.Unlock()
				if err == nil {
					mayChange = assertion.matcherMayChange(matcher, value)
					matches, e = matcher.Match(value)
					lock.Lock()
					err = e
					lock.Unlock()
				}
			case <-contextDone:
				fail("Context was cancelled")
				return false
			case <-timeout:
				fail("Timed out")
				return false
			}
		}
	} else if assertion.asyncType == AsyncAssertionTypeConsistently {
		for {
			if !(err == nil && matches == desiredMatch) {
				fail("Failed")
				return false
			}

			if !mayChange {
				return true
			}

			select {
			case <-time.After(assertion.pollingInterval):
				v, e := pollActual()
				lock.Lock()
				value, err = v, e
				lock.Unlock()
				if err == nil {
					mayChange = assertion.matcherMayChange(matcher, value)
					matches, e = matcher.Match(value)
					lock.Lock()
					err = e
					lock.Unlock()
				}
			case <-contextDone:
				fail("Context was cancelled")
				return false
			case <-timeout:
				return true
			}
		}
	}

	return false
}
