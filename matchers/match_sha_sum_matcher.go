package matchers

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/onsi/gomega/types"
)

func MatchSHASumOf(expected interface{}) types.GomegaMatcher {
	return &matchSHASumOfMatcher{
		expected: expected,
	}
}

type matchSHASumOfMatcher struct {
	expected interface{}
}

func (matcher *matchSHASumOfMatcher) Match(actual interface{}) (success bool, err error) {
	expectedFile, err := os.Open(matcher.expected.(string))
	if err != nil {
		return false, err
	}

	actualFile, err := os.Open(actual.(string))
	if err != nil {
		return false, err
	}

	expectedHash := sha1.New()
	actualHash := sha1.New()

	_, err = io.Copy(expectedHash, expectedFile)
	if err != nil {
		return false, err
	}

	_, err = io.Copy(actualHash, actualFile)
	if err != nil {
		return false, err
	}

	return reflect.DeepEqual(actualHash.Sum(nil), expectedHash.Sum(nil)), nil
}

func (matcher *matchSHASumOfMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\nto contain the same contents as\n\t%#v", actual, matcher.expected)
}

func (matcher *matchSHASumOfMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\nnot to contain the same contents as\n\t%#v", actual, matcher.expected)
}
