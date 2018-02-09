package testingtsupport

import "github.com/onsi/gomega/types"

type EmptyTWithHelper struct{}

func (e EmptyTWithHelper) Helper() {}

type gomegaTestingT interface {
	Fatalf(format string, args ...interface{})
}

func BuildTestingTGomegaFailWrapper(t gomegaTestingT) *types.GomegaFailWrapper {
	tWithHelper, ok := t.(types.TWithHelper)
	if !ok {
		tWithHelper = EmptyTWithHelper{}
	}

	fail := func(message string, callerSkip ...int) {
		tWithHelper.Helper()
		t.Fatalf("\n%s", message)
	}

	return &types.GomegaFailWrapper{
		Fail:        fail,
		TWithHelper: tWithHelper,
	}
}
