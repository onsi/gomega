/*
package gomockadaptor lets you use Gomega matchers as gomock argument matchers.

gomock (both github.com/golang/mock/gomock and go.uber.org/mock/gomock) accepts
any value implementing its Matcher interface:

	type Matcher interface {
		Matches(x any) bool
		String() string
	}

Match wraps a Gomega matcher in a value that satisfies that interface, so the
full set of Gomega matchers becomes available when setting up gomock expectations:

	mock.EXPECT().Send(gomockadaptor.Match(gomega.HaveLen(3)))
	mock.EXPECT().Store(gomockadaptor.Match(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
		"Name": gomega.Equal("Jane"),
	})))

The adaptor deliberately lives outside Gomega's core DSL and does not import
gomock, so it adds no dependencies to Gomega.
*/
package gomockadaptor

import "github.com/onsi/gomega/types"

// Adaptor wraps a Gomega matcher so it satisfies gomock's Matcher interface.
// Use Match to build one.
type Adaptor struct {
	matcher types.GomegaMatcher
	// failureMessage holds the message from the most recent failed match so
	// String can explain why an argument did not match.
	failureMessage string
}

// Match returns an Adaptor that wraps the given Gomega matcher. The returned
// value implements gomock's Matcher interface and can be passed directly to a
// gomock expectation.
func Match(matcher types.GomegaMatcher) *Adaptor {
	return &Adaptor{matcher: matcher}
}

// Matches reports whether x satisfies the wrapped Gomega matcher. A matcher
// that errors is treated as a non-match, with the error surfaced through
// String. This satisfies gomock's Matcher interface.
func (a *Adaptor) Matches(x any) bool {
	a.failureMessage = ""
	success, err := a.matcher.Match(x)
	if err != nil {
		a.failureMessage = err.Error()
		return false
	}
	if !success {
		a.failureMessage = a.matcher.FailureMessage(x)
	}
	return success
}

// String returns the failure message from the most recent failed match, so
// gomock can include it in its report. This satisfies gomock's Matcher
// interface.
func (a *Adaptor) String() string {
	return a.failureMessage
}
