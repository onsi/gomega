package gomega

import (
	"regexp"
	"runtime/debug"
	"strings"
)

//A simple *testing.T interface wrapper
type GomegaTestingT interface {
	Errorf(format string, args ...interface{})
}

//RegisterTestingT connects Gomega to Golang's XUnit style
//Testing.T tests.  You'll need to call this at the top of each XUnit style test:
//
// func TestFarmHasCow(t *testing.T) {
//     RegisterTestingT(t)
//
//	   f := farm.New([]string{"Cow", "Horse"})
//     Expect(f.HasCow()).To(BeTrue(), "Farm should have cow")
// }
//
// Note that this *testing.T is registered *globally* by Gomega (this is why you don't have to
// pass `t` down to the matcher itself).  This means that you cannot run the XUnit style tests
// in parallel as the global fail handler cannot point to more than one testing.T at a time.
//
// (As an aside: Ginkgo gets around this limitation by running parallel tests in different *processes*).
func RegisterTestingT(t GomegaTestingT) {
	globalFailHandler = buildTestingTOmegaFailHandler(t)
}

func buildTestingTOmegaFailHandler(t GomegaTestingT) OmegaFailHandler {
	return func(message string, callerSkip ...int) {
		skip := 1
		if len(callerSkip) > 0 {
			skip = callerSkip[0]
		}
		stackTrace := pruneStack(string(debug.Stack()), skip)
		t.Errorf("\n%s\n%s", stackTrace, message)
	}
}

func pruneStack(fullStackTrace string, skip int) string {
	stack := strings.Split(fullStackTrace, "\n")
	if len(stack) > 2*(skip+1) {
		stack = stack[2*(skip+1):]
	}
	prunedStack := []string{}
	re := regexp.MustCompile(`\/ginkgo\/|\/pkg\/testing\/|\/pkg\/runtime\/`)
	for i := 0; i < len(stack)/2; i++ {
		if !re.Match([]byte(stack[i*2])) {
			prunedStack = append(prunedStack, stack[i*2])
			prunedStack = append(prunedStack, stack[i*2+1])
		}
	}
	return strings.Join(prunedStack, "\n")
}
