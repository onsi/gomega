package gomega

import (
	"github.com/onsi/gomega/matchers"
)

func Equal(expected interface{}) OmegaMatcher {
	return &matchers.EqualMatcher{
		Expected: expected,
	}
}

func DeepEqual(expected interface{}) OmegaMatcher {
	return nil
}

func BeTrue() OmegaMatcher {
	return &matchers.TrueMatcher{}
}

func BeFalse() OmegaMatcher {
	return &matchers.FalseMatcher{}
}

func HaveOccured() OmegaMatcher {
	return &matchers.HaveOccuredMatcher{}
}

// func Be(comparator string, expected interface{}) *OmegaMatcher {
// 	return nil
// }

// func BeNil() *OmegaMatcher {
// 	return nil
// }

// func BeEmpty() *OmegaMatcher {
// 	return nil
// }

// func MatchRegex(matchingRegex string) *OmegaMatcher {
// 	return nil
// }

// func HaveLen(length int) *OmegaMatcher {
// 	return nil
// }

// func HaveKey(key interface{}) *OmegaMatcher {
// 	return nil
// }

// func ContainElement(element interface{}) *OmegaMatcher {
// 	return nil
// }

// func Panic() *OmegaMatcher {
// 	return nil
// }

// func HaveType() *OmegaMatcher {
// 	return nil
// }

// func ImplementInterface() *OmegaMatcher {
// 	return nil
// }
