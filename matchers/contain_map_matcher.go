package matchers

import (
	"fmt"
	"reflect"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/matchers/internal/miter"
	"github.com/onsi/gomega/matchers/support/goraph/bipartitegraph"
)

type ContainMapMatcher struct {
	Expected       any
	missingEntries []expectedMapEntry
}

type mapEntry struct {
	key any
	val any
}

type expectedMapEntry struct {
	keyMatcher   omegaMatcher
	valueMatcher omegaMatcher
	originalKey  any
	originalVal  any
}

func toEntryMatcher(val any) omegaMatcher {
	if val == nil {
		return &BeNilMatcher{}
	}
	if m, ok := val.(omegaMatcher); ok {
		return m
	}
	return &EqualMatcher{Expected: val}
}

func newExpectedMapEntry(key, val any) expectedMapEntry {
	return expectedMapEntry{
		keyMatcher:   toEntryMatcher(key),
		valueMatcher: toEntryMatcher(val),
		originalKey:  key,
		originalVal:  val,
	}
}

func mapEntryNeighbours(actual, expected any) (bool, error) {
	actEntry := actual.(mapEntry)
	expEntry := expected.(expectedMapEntry)

	match, err := expEntry.keyMatcher.Match(actEntry.key)
	if err != nil || !match {
		return false, nil
	}

	match, err = expEntry.valueMatcher.Match(actEntry.val)
	if err != nil || !match {
		return false, nil
	}

	return true, nil
}

func (matcher *ContainMapMatcher) Match(actual any) (success bool, err error) {
	if !isMap(actual) && !miter.IsSeq2(actual) {
		return false, fmt.Errorf("ContainMap matcher expects a map/iter.Seq2.  Got:%s", format.Object(actual, 1))
	}
	if !isMap(matcher.Expected) && !miter.IsSeq2(matcher.Expected) {
		return false, fmt.Errorf("ContainMap matcher expects a map/iter.Seq2.  Got:%s", format.Object(matcher.Expected, 1))
	}

	matcher.missingEntries = nil

	var expectedEntries []any
	if miter.IsSeq2(matcher.Expected) {
		miter.IterateKV(matcher.Expected, func(k, v reflect.Value) bool {
			expectedEntries = append(expectedEntries, newExpectedMapEntry(k.Interface(), v.Interface()))
			return true
		})
	} else {
		v := reflect.ValueOf(matcher.Expected)
		for _, k := range v.MapKeys() {
			expectedEntries = append(expectedEntries, newExpectedMapEntry(k.Interface(), v.MapIndex(k).Interface()))
		}
	}

	if len(expectedEntries) == 0 {
		return true, nil
	}

	var actualEntries []any
	if miter.IsSeq2(actual) {
		miter.IterateKV(actual, func(k, v reflect.Value) bool {
			actualEntries = append(actualEntries, mapEntry{key: k.Interface(), val: v.Interface()})
			return true
		})
	} else {
		v := reflect.ValueOf(actual)
		for _, k := range v.MapKeys() {
			actualEntries = append(actualEntries, mapEntry{key: k.Interface(), val: v.MapIndex(k).Interface()})
		}
	}

	bipartiteGraph, err := bipartitegraph.NewBipartiteGraph(actualEntries, expectedEntries, mapEntryNeighbours)
	if err != nil {
		return false, err
	}

	edges := bipartiteGraph.LargestMatching()
	if len(edges) == len(expectedEntries) {
		return true, nil
	}

	_, missing := bipartiteGraph.FreeLeftRight(edges)
	for _, m := range missing {
		matcher.missingEntries = append(matcher.missingEntries, m.(expectedMapEntry))
	}

	return false, nil
}

func (matcher *ContainMapMatcher) FailureMessage(actual any) (message string) {
	message = format.Message(actual, "to contain map", matcher.Expected)
	return appendMissingMapEntries(message, matcher.missingEntriesPresentable())
}

func (matcher *ContainMapMatcher) NegatedFailureMessage(actual any) (message string) {
	return format.Message(actual, "not to contain map", matcher.Expected)
}

func appendMissingMapEntries(message string, missing any) string {
	if missing == nil {
		return message
	}
	return fmt.Sprintf("%s\nthe missing entries were\n%s", message,
		format.Object(missing, 1))
}

func (matcher *ContainMapMatcher) missingEntriesPresentable() any {
	if len(matcher.missingEntries) == 0 {
		return nil
	}

	if isMap(matcher.Expected) {
		expType := reflect.TypeOf(matcher.Expected)
		keyType := expType.Key()
		valType := expType.Elem()

		canSet := true
		for _, entry := range matcher.missingEntries {
			kVal := reflect.ValueOf(entry.originalKey)
			if !kVal.IsValid() || !kVal.Type().AssignableTo(keyType) {
				canSet = false
				break
			}
			if entry.originalVal == nil {
				if valType.Kind() != reflect.Interface && valType.Kind() != reflect.Pointer && valType.Kind() != reflect.Map && valType.Kind() != reflect.Slice && valType.Kind() != reflect.Func && valType.Kind() != reflect.Chan {
					canSet = false
					break
				}
			} else {
				vVal := reflect.ValueOf(entry.originalVal)
				if !vVal.Type().AssignableTo(valType) {
					canSet = false
					break
				}
			}
		}

		if canSet {
			res := reflect.MakeMapWithSize(expType, len(matcher.missingEntries))
			for _, entry := range matcher.missingEntries {
				kVal := reflect.ValueOf(entry.originalKey)
				var vVal reflect.Value
				if entry.originalVal == nil {
					vVal = reflect.Zero(valType)
				} else {
					vVal = reflect.ValueOf(entry.originalVal)
				}
				res.SetMapIndex(kVal, vVal)
			}
			return res.Interface()
		}
	}

	canBeMap := true
	for _, entry := range matcher.missingEntries {
		if entry.originalKey == nil {
			canBeMap = false
			break
		}
		if !reflect.TypeOf(entry.originalKey).Comparable() {
			canBeMap = false
			break
		}
	}
	if canBeMap {
		res := make(map[any]any, len(matcher.missingEntries))
		for _, entry := range matcher.missingEntries {
			res[entry.originalKey] = entry.originalVal
		}
		return res
	}

	type entryStruct struct {
		Key   any
		Value any
	}
	entries := make([]entryStruct, len(matcher.missingEntries))
	for i, entry := range matcher.missingEntries {
		entries[i] = entryStruct{Key: entry.originalKey, Value: entry.originalVal}
	}
	return entries
}
