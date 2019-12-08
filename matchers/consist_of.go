// untested sections: 3

package matchers

import (
	"fmt"
	"reflect"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/matchers/support/goraph/bipartitegraph"
)

type ConsistOfMatcher struct {
	Elements        []interface{}
	missingElements []interface{}
	extraElements   []interface{}
}

func (matcher *ConsistOfMatcher) Match(actual interface{}) (success bool, err error) {
	if !isArrayOrSlice(actual) && !isMap(actual) {
		return false, fmt.Errorf("ConsistOf matcher expects an array/slice/map.  Got:\n%s", format.Object(actual, 1))
	}

	elements := matcher.Elements
	if len(matcher.Elements) == 1 && isArrayOrSlice(matcher.Elements[0]) {
		elements = []interface{}{}
		value := reflect.ValueOf(matcher.Elements[0])
		for i := 0; i < value.Len(); i++ {
			elements = append(elements, value.Index(i).Interface())
		}
	}

	matchers := []interface{}{}
	for _, element := range elements {
		matcher, isMatcher := element.(omegaMatcher)
		if !isMatcher {
			matcher = &EqualMatcher{Expected: element}
		}
		matchers = append(matchers, matcher)
	}

	values := matcher.valuesOf(actual)

	neighbours := func(v, m interface{}) (bool, error) {
		match, err := m.(omegaMatcher).Match(v)
		return match && err == nil, nil
	}

	bipartiteGraph, err := bipartitegraph.NewBipartiteGraph(values, matchers, neighbours)
	if err != nil {
		return false, err
	}

	edges := bipartiteGraph.LargestMatching()
	if len(edges) == len(values) && len(edges) == len(matchers) {
		return true, nil
	}

	var missingMatchers []interface{}
	matcher.extraElements, missingMatchers = bipartiteGraph.FreeLeftRight(edges)

	for _, missing := range missingMatchers {
		equalMatcher, ok := missing.(*EqualMatcher)
		if ok {
			missing = equalMatcher.Expected
		}
		matcher.missingElements = append(matcher.missingElements, missing)
	}

	return false, nil
}

func (matcher *ConsistOfMatcher) valuesOf(actual interface{}) []interface{} {
	value := reflect.ValueOf(actual)
	values := []interface{}{}
	if isMap(actual) {
		keys := value.MapKeys()
		for i := 0; i < value.Len(); i++ {
			values = append(values, value.MapIndex(keys[i]).Interface())
		}
	} else {
		for i := 0; i < value.Len(); i++ {
			values = append(values, value.Index(i).Interface())
		}
	}

	return values
}

func (matcher *ConsistOfMatcher) FailureMessage(actual interface{}) (message string) {
	message = format.Message(actual, "to consist of", matcher.Elements)
	if len(matcher.missingElements) > 0 {
		message = fmt.Sprintf("%s\nthe missing elements were\n%s", message,
			format.Object(matcher.missingElements, 1))
	}
	if len(matcher.extraElements) > 0 {
		message = fmt.Sprintf("%s\nthe extra elements were\n%s", message,
			format.Object(matcher.extraElements, 1))
	}
	return
}

func (matcher *ConsistOfMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "not to consist of", matcher.Elements)
}
