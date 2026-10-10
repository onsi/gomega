package matchers

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/matchers/internal/miter"
	"github.com/onsi/gomega/matchers/support/goraph/bipartitegraph"
)

type ContainMapMatcher struct {
	Expected any

	missingKeys       []containMapEntry
	mismatchedEntries []containMapMismatch
	unmatchedEntries  []containMapEntry
}

// containMapEntry is a key:value pair in the expected map
type containMapEntry struct {
	key, value               reflect.Value
	keyMatcher, valueMatcher omegaMatcher

	// plainKey is true when the key is a (non-nil) value rather than a matcher
	plainKey bool
	// lookup is true for plain keys that can be found in actual with a direct (==)
	// lookup rather than by testing every entry in actual
	lookup bool
	// for lookup entries: the actual entries with an equal key, and those whose value
	// also satisfies valueMatcher
	candidates, matches []containMapActualEntry
}

// containMapActualEntry is a key:value pair in actual.  token identifies the entry:
// the key itself when actual is a map, the entry's position when actual is an iter.Seq2
type containMapActualEntry struct {
	token      any
	key, value any
}

type containMapMismatch struct {
	key          reflect.Value
	actualValue  any
	valueMatcher omegaMatcher
}

func containMapEntryMatcher(v reflect.Value) omegaMatcher {
	if v.Interface() == nil {
		return &BeNilMatcher{}
	}
	if m, ok := v.Interface().(omegaMatcher); ok {
		return m
	}
	return &EqualMatcher{Expected: v.Interface()}
}

// containMapLookupable reports whether == on values of type t agrees with reflect.DeepEqual (and
// so with Equal), which makes it safe to find a key of type t with a map lookup.
func containMapLookupable(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Bool, reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return true
	case reflect.Array:
		return containMapLookupable(t.Elem())
	case reflect.Struct:
		for i := range t.NumField() {
			f := t.Field(i)
			if f.Name == "_" || !containMapLookupable(f.Type) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func (matcher *ContainMapMatcher) Match(actual any) (success bool, err error) {
	if !isMap(actual) && !miter.IsSeq2(actual) {
		return false, fmt.Errorf("ContainMap matcher expects a map/iter.Seq2.  Got:\n%s", format.Object(actual, 1))
	}
	if !isMap(matcher.Expected) && !miter.IsSeq2(matcher.Expected) {
		return false, fmt.Errorf("ContainMap matcher expects the expected value to be a map/iter.Seq2.  Got:\n%s", format.Object(matcher.Expected, 1))
	}
	matcher.missingKeys, matcher.mismatchedEntries, matcher.unmatchedEntries = nil, nil, nil

	var expected []*containMapEntry
	addExpected := func(k, v reflect.Value) {
		entry := &containMapEntry{key: k, value: v, keyMatcher: containMapEntryMatcher(k), valueMatcher: containMapEntryMatcher(v)}
		if _, isMatcher := k.Interface().(omegaMatcher); !isMatcher && k.Interface() != nil {
			entry.plainKey = true
			entry.lookup = containMapLookupable(reflect.TypeOf(k.Interface()))
		}
		expected = append(expected, entry)
	}
	if isMap(matcher.Expected) {
		expectedValue := reflect.ValueOf(matcher.Expected)
		for _, k := range expectedValue.MapKeys() {
			addExpected(k, expectedValue.MapIndex(k))
		}
	} else {
		miter.IterateKV(matcher.Expected, func(k, v reflect.Value) bool {
			addExpected(k, v)
			return true
		})
	}
	if len(expected) == 0 {
		return true, nil
	}

	// Entries in actual are only enumerated when we have to: if actual is a map and every expected
	// key can be looked up directly, matching is O(len(expected)).
	var actualEntries []containMapActualEntry
	var actualIndex map[any][]int
	actualIsMap := isMap(actual)
	actualValue := reflect.ValueOf(actual)
	enumerated := false
	enumerateActual := func() {
		if enumerated {
			return
		}
		enumerated = true
		add := func(token any, k, v reflect.Value) {
			key := k.Interface()
			if !actualIsMap && key != nil && containMapLookupable(reflect.TypeOf(key)) {
				if actualIndex == nil {
					actualIndex = map[any][]int{}
				}
				actualIndex[key] = append(actualIndex[key], len(actualEntries))
			}
			actualEntries = append(actualEntries, containMapActualEntry{token: token, key: key, value: v.Interface()})
		}
		if actualIsMap {
			for _, k := range actualValue.MapKeys() {
				add(k.Interface(), k, actualValue.MapIndex(k))
			}
		} else {
			i := 0
			miter.IterateKV(actual, func(k, v reflect.Value) bool {
				add(i, k, v)
				i++
				return true
			})
		}
	}
	if !actualIsMap || slices.ContainsFunc(expected, func(e *containMapEntry) bool { return !e.lookup }) {
		enumerateActual()
	}

	// Phase 1: look up the plain keys and test the values we find.
	var erroredEntries []*containMapEntry
	var valueErrors []error
	claims := map[any]int{}
	unique := true
	for _, e := range expected {
		if !e.lookup {
			continue
		}
		key := e.key.Interface()
		if actualIsMap {
			keyType := actualValue.Type().Key()
			t := reflect.TypeOf(key)
			if t == keyType || keyType.Kind() == reflect.Interface && t.Implements(keyType) {
				if v := actualValue.MapIndex(reflect.ValueOf(key)); v.IsValid() {
					e.candidates = []containMapActualEntry{{token: key, key: key, value: v.Interface()}}
				}
			}
		} else {
			for _, i := range actualIndex[key] {
				e.candidates = append(e.candidates, actualEntries[i])
			}
		}
		var lastErr error
		for _, candidate := range e.candidates {
			success, err := e.valueMatcher.Match(candidate.value)
			if err != nil {
				lastErr = err
				continue
			}
			if success {
				e.matches = append(e.matches, candidate)
				claims[candidate.token]++
			}
		}
		if len(e.matches) == 0 && lastErr != nil {
			erroredEntries = append(erroredEntries, e)
			valueErrors = append(valueErrors, lastErr)
		}
		if len(e.matches) > 1 {
			unique = false
		}
	}
	if len(erroredEntries) > 0 {
		// report the same error regardless of map iteration order
		i := 0
		for j := range erroredEntries {
			if format.Object(erroredEntries[j].key.Interface(), 0) < format.Object(erroredEntries[i].key.Interface(), 0) {
				i = j
			}
		}
		return false, fmt.Errorf("ContainMap's value matcher for key\n%s\nfailed with:\n%s", format.Object(erroredEntries[i].key.Interface(), 1), format.IndentString(valueErrors[i].Error(), 1))
	}
	for _, n := range claims {
		if n > 1 {
			unique = false
		}
	}

	// Phase 2: match everything else.  In the usual case every plain key matched at most one entry in actual
	// and no entry was matched by two plain keys: those entries are then the only possible assignment, so we
	// set them aside and only need to find distinct entries for the remaining (matcher or non-lookup) keys.
	// Otherwise (iter.Seq2s with repeated keys) we fall back to matching every expected entry at once.
	var unmatched []*containMapEntry
	var pending []*containMapEntry
	var available []any
	if unique {
		for _, e := range expected {
			if !e.lookup {
				pending = append(pending, e)
			} else if len(e.matches) == 0 {
				unmatched = append(unmatched, e)
			}
		}
		for _, a := range actualEntries {
			if claims[a.token] == 0 {
				available = append(available, a)
			}
		}
	} else {
		enumerateActual()
		pending = expected
		for _, a := range actualEntries {
			available = append(available, a)
		}
	}
	if len(pending) > 0 {
		pendingValues := make([]any, len(pending))
		for i, e := range pending {
			pendingValues[i] = e
		}
		graph, err := bipartitegraph.NewBipartiteGraph(available, pendingValues, containMapNeighbours)
		if err != nil {
			return false, err
		}
		edges := graph.LargestMatching()
		if len(edges) < len(pending) {
			_, free := graph.FreeLeftRight(edges)
			for _, e := range free {
				unmatched = append(unmatched, e.(*containMapEntry))
			}
		}
	}
	if len(unmatched) == 0 {
		return true, nil
	}

	matcher.classifyUnmatched(unmatched, actualEntries, enumerateActual)
	return false, nil
}

func containMapNeighbours(a, e any) (bool, error) {
	actualEntry, expectedEntry := a.(containMapActualEntry), e.(*containMapEntry)
	if expectedEntry.lookup {
		return slices.ContainsFunc(expectedEntry.matches, func(m containMapActualEntry) bool { return m.token == actualEntry.token }), nil
	}
	// as with ContainElements, a matcher that errors simply does not match
	if success, err := expectedEntry.keyMatcher.Match(actualEntry.key); err != nil || !success {
		return false, nil
	}
	if success, err := expectedEntry.valueMatcher.Match(actualEntry.value); err != nil || !success {
		return false, nil
	}
	return true, nil
}

// classifyUnmatched sorts unmatched expected entries into missing keys, keys whose value did
// not match, and entries that could not be given a distinct entry in actual (including every
// unmatched entry with a matcher as its key)
func (matcher *ContainMapMatcher) classifyUnmatched(unmatched []*containMapEntry, actualEntries []containMapActualEntry, enumerateActual func()) {
	for _, e := range unmatched {
		if !e.plainKey {
			matcher.unmatchedEntries = append(matcher.unmatchedEntries, *e)
			continue
		}
		candidates, matched := e.candidates, len(e.matches) > 0
		if !e.lookup {
			enumerateActual()
			candidates = nil
			for _, a := range actualEntries {
				if success, err := e.keyMatcher.Match(a.key); err == nil && success {
					candidates = append(candidates, a)
					if success, err := e.valueMatcher.Match(a.value); err == nil && success {
						matched = true
					}
				}
			}
		}
		switch {
		case len(candidates) == 0:
			matcher.missingKeys = append(matcher.missingKeys, *e)
		case !matched:
			matcher.mismatchedEntries = append(matcher.mismatchedEntries, containMapMismatch{key: e.key, actualValue: candidates[0].value, valueMatcher: e.valueMatcher})
		default:
			matcher.unmatchedEntries = append(matcher.unmatchedEntries, *e)
		}
	}
}

func (matcher *ContainMapMatcher) FailureMessage(actual any) (message string) {
	message = format.Message(actual, "to contain map", matcher.Expected)

	if len(matcher.missingKeys) > 0 {
		keys := reflect.MakeSlice(reflect.SliceOf(matcher.missingKeys[0].key.Type()), 0, len(matcher.missingKeys))
		for _, e := range sortedByKey(matcher.missingKeys, func(e containMapEntry) reflect.Value { return e.key }) {
			keys = reflect.Append(keys, e.key)
		}
		message += "\nthe missing keys were\n" + format.Object(keys.Interface(), 1)
	}

	if len(matcher.mismatchedEntries) > 0 {
		message += "\nthe mismatched values were"
		for _, m := range sortedByKey(matcher.mismatchedEntries, func(m containMapMismatch) reflect.Value { return m.key }) {
			message += "\n" + format.Object(m.key.Interface(), 1) + "\n" + format.IndentString(m.valueMatcher.FailureMessage(m.actualValue), 2)
		}
	}

	if len(matcher.unmatchedEntries) > 0 {
		entries := make([]string, len(matcher.unmatchedEntries))
		for i, e := range matcher.unmatchedEntries {
			entries[i] = format.Object(e.key.Interface(), 1) + "\n" + format.Object(e.value.Interface(), 2)
		}
		slices.Sort(entries)
		message += "\nthe entries with no matching (distinct) entry in actual were\n" + strings.Join(entries, "\n")
	}

	return message
}

func (matcher *ContainMapMatcher) NegatedFailureMessage(actual any) (message string) {
	return format.Message(actual, "not to contain map", matcher.Expected)
}

// sortedByKey sorts entries by the formatted representation of their key so that failure messages
// do not depend on map iteration order
func sortedByKey[T any](entries []T, key func(T) reflect.Value) []T {
	sorted := slices.Clone(entries)
	slices.SortStableFunc(sorted, func(a, b T) int {
		return strings.Compare(format.Object(key(a).Interface(), 0), format.Object(key(b).Interface(), 0))
	})
	return sorted
}
