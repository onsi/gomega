package internal

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

var errInterface = reflect.TypeOf((*error)(nil)).Elem()
var gomegaType = reflect.TypeOf((*types.Gomega)(nil)).Elem()
var contextType = reflect.TypeOf(new(context.Context)).Elem()

type contextWithAttachProgressReporter interface {
	AttachProgressReporter(func() string) func()
}

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

	actualIsFunc  bool
	actual        interface{}
	argsToForward []interface{}

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

func (assertion *AsyncAssertion) WithArguments(argsToForward ...interface{}) types.AsyncAssertion {
	assertion.argsToForward = argsToForward
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

func (assertion *AsyncAssertion) processReturnValues(values []reflect.Value) (interface{}, error, *AsyncSignalError) {
	var err error
	var asyncSignal *AsyncSignalError = nil

	if len(values) == 0 {
		return nil, fmt.Errorf("No values were returned by the function passed to Gomega"), asyncSignal
	}
	actual := values[0].Interface()
	if asyncSignalErr, ok := AsAsyncSignalError(actual); ok {
		asyncSignal = asyncSignalErr
	}
	for i, extraValue := range values[1:] {
		extra := extraValue.Interface()
		if extra == nil {
			continue
		}
		if asyncSignalErr, ok := AsAsyncSignalError(extra); ok {
			asyncSignal = asyncSignalErr
			continue
		}
		extraType := reflect.TypeOf(extra)
		zero := reflect.Zero(extraType).Interface()
		if reflect.DeepEqual(extra, zero) {
			continue
		}
		if i == len(values)-2 && extraType.Implements(errInterface) {
			err = fmt.Errorf("function returned error: %s\n%s", extra, format.Object(extra, 1))
			continue
		}
		if err == nil {
			err = fmt.Errorf("Unexpected non-nil/non-zero return value at index %d:\n\t<%T>: %#v", i+1, extra, extra)
		}
	}
	if err == nil {
		err = errors.Unwrap(asyncSignal)
	}

	return actual, err, asyncSignal
}

func (assertion *AsyncAssertion) invalidFunctionError(t reflect.Type) error {
	return fmt.Errorf(`The function passed to %s had an invalid signature of %s.  Functions passed to %s must either:

	(a) have return values or
	(b) take a Gomega interface as their first argument and use that Gomega instance to make assertions.

You can learn more at https://onsi.github.io/gomega/#eventually
`, assertion.asyncType, t, assertion.asyncType)
}

func (assertion *AsyncAssertion) noConfiguredContextForFunctionError() error {
	return fmt.Errorf(`The function passed to %s requested a context.Context, but no context has been provided.  Please pass one in using %s().WithContext().

You can learn more at https://onsi.github.io/gomega/#eventually
`, assertion.asyncType, assertion.asyncType)
}

func (assertion *AsyncAssertion) argumentMismatchError(t reflect.Type, numProvided int) error {
	have := "have"
	if numProvided == 1 {
		have = "has"
	}
	return fmt.Errorf(`The function passed to %s has signature %s takes %d arguments but %d %s been provided.  Please use %s().WithArguments() to pass the corect set of arguments.

You can learn more at https://onsi.github.io/gomega/#eventually
`, assertion.asyncType, t, t.NumIn(), numProvided, have, assertion.asyncType)
}

func (assertion *AsyncAssertion) buildActualPoller() (func() (interface{}, error, *AsyncSignalError), error) {
	if !assertion.actualIsFunc {
		return func() (interface{}, error, *AsyncSignalError) { return assertion.actual, nil, nil }, nil
	}
	actualValue := reflect.ValueOf(assertion.actual)
	actualType := reflect.TypeOf(assertion.actual)
	numIn, numOut, isVariadic := actualType.NumIn(), actualType.NumOut(), actualType.IsVariadic()

	if numIn == 0 && numOut == 0 {
		return nil, assertion.invalidFunctionError(actualType)
	}
	takesGomega, takesContext := false, false
	if numIn > 0 {
		takesGomega, takesContext = actualType.In(0).Implements(gomegaType), actualType.In(0).Implements(contextType)
	}
	if takesGomega && numIn > 1 && actualType.In(1).Implements(contextType) {
		takesContext = true
	}
	if takesContext && len(assertion.argsToForward) > 0 && reflect.TypeOf(assertion.argsToForward[0]).Implements(contextType) {
		takesContext = false
	}
	if !takesGomega && numOut == 0 {
		return nil, assertion.invalidFunctionError(actualType)
	}
	if takesContext && assertion.ctx == nil {
		return nil, assertion.noConfiguredContextForFunctionError()
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
	for _, arg := range assertion.argsToForward {
		inValues = append(inValues, reflect.ValueOf(arg))
	}

	if !isVariadic && numIn != len(inValues) {
		return nil, assertion.argumentMismatchError(actualType, len(inValues))
	} else if isVariadic && len(inValues) < numIn-1 {
		return nil, assertion.argumentMismatchError(actualType, len(inValues))
	}

	return func() (actual interface{}, err error, asyncSignal *AsyncSignalError) {
		var values []reflect.Value
		assertionFailure = nil
		defer func() {
			if numOut == 0 && takesGomega {
				actual = assertionFailure
			} else {
				actual, err, asyncSignal = assertion.processReturnValues(values)
				if assertionFailure != nil {
					err = assertionFailure
				}
			}
			if e := recover(); e != nil {
				if asyncSignalErr, ok := AsAsyncSignalError(e); ok {
					asyncSignal = asyncSignalErr
					if err == nil {
						err = errors.Unwrap(asyncSignal)
					}
				} else if assertionFailure == nil {
					panic(e)
				}
			}
		}()
		values = actualValue.Call(inValues)
		return
	}, nil
}

func (assertion *AsyncAssertion) afterTimeout() <-chan time.Time {
	if assertion.timeoutInterval >= 0 {
		return time.After(assertion.timeoutInterval)
	}

	if assertion.asyncType == AsyncAssertionTypeConsistently {
		return time.After(assertion.g.DurationBundle.ConsistentlyDuration)
	} else {
		if assertion.ctx == nil {
			return time.After(assertion.g.DurationBundle.EventuallyTimeout)
		} else {
			return nil
		}
	}
}

func (assertion *AsyncAssertion) afterPolling() <-chan time.Time {
	if assertion.pollingInterval >= 0 {
		return time.After(assertion.pollingInterval)
	}
	if assertion.asyncType == AsyncAssertionTypeConsistently {
		return time.After(assertion.g.DurationBundle.ConsistentlyPollingInterval)
	} else {
		return time.After(assertion.g.DurationBundle.EventuallyPollingInterval)
	}
}

func (assertion *AsyncAssertion) matcherSaysStopTrying(matcher types.GomegaMatcher, value interface{}) *AsyncSignalError {
	if assertion.actualIsFunc || types.MatchMayChangeInTheFuture(matcher, value) {
		return nil
	}
	return StopTrying("No future change is possible.  Bailing out early").(*AsyncSignalError)
}

func (assertion *AsyncAssertion) pollMatcher(matcher types.GomegaMatcher, value interface{}, currentAsyncSignal *AsyncSignalError) (matches bool, err error, asyncSignal *AsyncSignalError) {
	// we pass through the current StopTrying error and only overwrite it with what the matcher says if it is nil
	asyncSignal = currentAsyncSignal

	if currentAsyncSignal == nil || !currentAsyncSignal.StopTrying() {
		asyncSignal = assertion.matcherSaysStopTrying(matcher, value)
	}

	defer func() {
		if e := recover(); e != nil {
			if asyncSignalErr, ok := AsAsyncSignalError(e); ok {
				if asyncSignal == nil {
					asyncSignal = asyncSignalErr
				}
				err = asyncSignalErr
			} else {
				panic(e)
			}
		}
	}()

	matches, err = matcher.Match(value)
	if asyncSignalErr, ok := AsAsyncSignalError(err); ok {
		err = errors.Unwrap(asyncSignalErr)
		if asyncSignal == nil {
			asyncSignal = asyncSignalErr
		}
	}

	return
}

func (assertion *AsyncAssertion) match(matcher types.GomegaMatcher, desiredMatch bool, optionalDescription ...interface{}) bool {
	timer := time.Now()
	timeout := assertion.afterTimeout()
	lock := sync.Mutex{}

	var matches bool
	var err error

	assertion.g.THelper()

	pollActual, err := assertion.buildActualPoller()
	if err != nil {
		assertion.g.Fail(err.Error(), 2+assertion.offset)
		return false
	}

	value, err, asyncSignal := pollActual()

	if err == nil {
		matches, err, asyncSignal = assertion.pollMatcher(matcher, value, asyncSignal)
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

			if asyncSignal != nil && asyncSignal.StopTrying() {
				fail(asyncSignal.Error() + " -")
				return false
			}

			select {
			case <-assertion.afterPolling():
				v, e, as := pollActual()
				if as != nil && as.WasViaPanic() && as.StopTrying() {
					// we were told to stop trying via panic - which means we dont' have reasonable new values
					// we should simply use the old values and exit now
					fail(as.Error() + " -")
					return false
				}
				lock.Lock()
				value, err, asyncSignal = v, e, as
				lock.Unlock()
				if err == nil {
					m, e, as := assertion.pollMatcher(matcher, value, asyncSignal)
					lock.Lock()
					matches, err, asyncSignal = m, e, as
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

			if asyncSignal != nil && asyncSignal.StopTrying() {
				return true
			}

			select {
			case <-assertion.afterPolling():
				v, e, as := pollActual()
				if as != nil && as.WasViaPanic() && as.StopTrying() {
					// we were told to stop trying via panic - which means we made it this far and should return successfully
					return true
				}
				lock.Lock()
				value, err, asyncSignal = v, e, as
				lock.Unlock()
				if err == nil {
					m, e, as := assertion.pollMatcher(matcher, value, asyncSignal)
					lock.Lock()
					matches, err, asyncSignal = m, e, as
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
