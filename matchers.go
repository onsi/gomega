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

func MatchRegexp(regexp string) OmegaMatcher {
	return &matchers.MatchRegexpMatcher{
		Regexp: regexp,
	}
}

func ContainSubstring(substr string) OmegaMatcher {
	return &matchers.ContainSubstringMatcher{
		Substr: substr,
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

func BeNumerically(comparator string, compareTo ...interface{}) OmegaMatcher {
	return &matchers.BeNumericallyMatcher{
		Comparator: comparator,
		CompareTo:  compareTo,
	}
}

//HaveKey
//Panic

//HaveSameElementsAs? (//order independent)
//BeSameInstanceAs //identical!

//HaveSameTypeAs
//ImplementSameInterfaceAs
