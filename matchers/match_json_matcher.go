package matchers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

type MatchJSONMatcher struct {
	JSONToMatch interface{}
}

func (matcher *MatchJSONMatcher) Match(actual interface{}) (success bool, message string, err error) {
	actualString, aok := toString(actual)
	expectedString, eok := toString(matcher.JSONToMatch)

	if aok && eok {
		abuf := new(bytes.Buffer)
		ebuf := new(bytes.Buffer)

		if err := json.Indent(abuf, []byte(actualString), "", ""); err != nil {
			return false, "", err
		}

		if err := json.Indent(ebuf, []byte(expectedString), "", ""); err != nil {
			return false, "", err
		}

		var aval interface{}
		var eval interface{}

		json.Unmarshal([]byte(actualString), &aval)
		json.Unmarshal([]byte(expectedString), &eval)

		if reflect.DeepEqual(aval, eval) {
			return true, formatMessage(abuf.String(), "not to match JSON of", ebuf.String()), nil
		} else {
			return false, formatMessage(abuf.String(), "to match JSON of", ebuf.String()), nil
		}
	} else {
		return false, "", fmt.Errorf("MatchJSONMatcher matcher requires a string or stringer.  Got:%s", formatObject(actual))
	}
	return false, "", nil
}
