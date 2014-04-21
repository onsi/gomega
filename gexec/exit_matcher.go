package gexec

import (
	"fmt"

	"github.com/onsi/gomega/format"
)

func Exit(optionalExitCode ...int) *exitMatcher {
	exitCode := -1
	if len(optionalExitCode) > 0 {
		exitCode = optionalExitCode[0]
	}

	return &exitMatcher{
		exitCode: exitCode,
	}
}

type exitMatcher struct {
	exitCode       int
	didExit        bool
	actualExitCode int
}

func (m *exitMatcher) Match(actual interface{}) (success bool, err error) {
	session, ok := actual.(*Session)
	if !ok {
		return false, fmt.Errorf("Exit must be passed a gexit session.  Got:\n%s", format.Object(actual, 1))
	}

	m.actualExitCode = session.getExitCode()

	if m.actualExitCode == -1 {
		return false, nil
	}

	if m.exitCode == -1 {
		return true, nil
	}
	return m.exitCode == m.actualExitCode, nil
}

func (m *exitMatcher) FailureMessage(actual interface{}) (message string) {
	if m.actualExitCode == -1 {
		return "Expected process to exit.  It did not."
	} else {
		return format.Message(m.actualExitCode, "to match exit code:", m.exitCode)
	}
}

func (m *exitMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	if m.actualExitCode == -1 {
		return "you really shouldn't be able to see this!"
	} else {
		if m.exitCode == -1 {
			return "Expected process not to exit.  It did."
		} else {
			return format.Message(m.actualExitCode, "not to match exit code:", m.exitCode)
		}
	}
}
