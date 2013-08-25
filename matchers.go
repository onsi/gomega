package gomega

import (
	"github.com/onsi/gomega/matchers"
)

func Equal(expected interface{}) OmegaMatcher {
	return &matchers.EqualMatcher{
		Expected: expected,
	}
}

func BeNil() OmegaMatcher {
	return &matchers.BeNilMatcher{}
}

func BeTrue() OmegaMatcher {
	return &matchers.BeTrueMatcher{}
}

func BeFalse() OmegaMatcher {
	return &matchers.BeFalseMatcher{}
}

func HaveOccured() OmegaMatcher {
	return &matchers.HaveOccuredMatcher{}
}

func MatchRegexp(regexp string, args ...interface{}) OmegaMatcher {
	return &matchers.MatchRegexpMatcher{
		Regexp: regexp,
		Args:   args,
	}
}

func ContainSubstring(substr string, args ...interface{}) OmegaMatcher {
	return &matchers.ContainSubstringMatcher{
		Substr: substr,
		Args:   args,
	}
}

func BeEmpty() OmegaMatcher {
	return &matchers.BeEmptyMatcher{}
}

func HaveLen(count int) OmegaMatcher {
	return &matchers.HaveLenMatcher{
		Count: count,
	}
}

func BeZero() OmegaMatcher {
	return &matchers.BeZeroMatcher{}
}

func ContainElement(element interface{}) OmegaMatcher {
	return &matchers.ContainElementMatcher{
		Element: element,
	}
}

func HaveKey(key interface{}) OmegaMatcher {
	return &matchers.HaveKeyMatcher{
		Key: key,
	}
}

func BeNumerically(comparator string, compareTo ...interface{}) OmegaMatcher {
	return &matchers.BeNumericallyMatcher{
		Comparator: comparator,
		CompareTo:  compareTo,
	}
}

// func Panic() OmegaMatcher {
// 	return &matchers.PanicMatcher{}
// }

// func Receive(expected interface{}) OmegaMatcher {
// 	return &matchers.ReceiveMatcher{
// 		Expected: expected,
// 	}
// }

//TODO:
//HaveSameTypeAs
//ImplementSameInterfaceAs

//MAYBE:
//HaveSameElementsAs? (//order independent)
//BeSameInstanceAs? (//pointer identity)
