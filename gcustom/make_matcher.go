package gcustom

import (
	"fmt"
	"reflect"
	"strings"
	"text/template"

	"github.com/onsi/gomega/format"
)

var interfaceType = reflect.TypeOf((*interface{})(nil)).Elem()
var errInterface = reflect.TypeOf((*error)(nil)).Elem()

var defaultTemplate = template.Must(ParseTemplate("{{if .Failure}}Custom matcher failed for:{{else}}Custom matcher succeeded (but was expected to fail) for:{{end}}\n{{.FormattedActual}}"))

func formatObject(object any, indent ...uint) string {
	indentation := uint(0)
	if len(indent) > 0 {
		indentation = indent[0]
	}
	return format.Object(object, indentation)
}

func ParseTemplate(templ string) (*template.Template, error) {
	return template.New("template").Funcs(template.FuncMap{
		"format": formatObject,
	}).Parse(templ)
}

func MakeMatcher(matchFunc any, args ...any) CustomGomegaMatcher {
	t := reflect.TypeOf(matchFunc)
	if !(t.Kind() == reflect.Func && t.NumIn() == 1 && t.NumOut() == 2 && t.Out(0).Kind() == reflect.Bool && t.Out(1).Implements(errInterface)) {
		panic("MakeMatcher must be passed a function that takes one argument and returns (bool, error)")
	}
	var finalMatchFunc func(actual any) (bool, error)
	if t.In(0) == interfaceType {
		finalMatchFunc = matchFunc.(func(actual any) (bool, error))
	} else {
		matchFuncValue := reflect.ValueOf(matchFunc)
		finalMatchFunc = reflect.MakeFunc(reflect.TypeOf(finalMatchFunc),
			func(args []reflect.Value) []reflect.Value {
				actual := args[0].Interface()
				if reflect.TypeOf(actual).AssignableTo(t.In(0)) {
					return matchFuncValue.Call([]reflect.Value{reflect.ValueOf(actual)})
				} else {
					return []reflect.Value{
						reflect.ValueOf(false),
						reflect.ValueOf(fmt.Errorf("Matcher expected actual of type <%s>.  Got:\n%s", t.In(0), format.Object(actual, 1))),
					}
				}
			}).Interface().(func(actual any) (bool, error))
	}

	matcher := CustomGomegaMatcher{
		matchFunc:       finalMatchFunc,
		templateMessage: defaultTemplate,
	}

	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			matcher = matcher.WithMessage(v)
		case *template.Template:
			matcher = matcher.WithPrecompiledTemplate(v)
		}
	}

	return matcher
}

type CustomGomegaMatcher struct {
	matchFunc                   func(actual any) (bool, error)
	templateMessage             *template.Template
	templateData                any
	customFailureMessage        func(actual any) string
	customNegatedFailureMessage func(actual any) string
}

func (c CustomGomegaMatcher) WithMessage(message string) CustomGomegaMatcher {
	return c.WithTemplate("Expected:\n{{.FormattedActual}}\n{{.To}} " + message)
}

func (c CustomGomegaMatcher) WithTemplate(templ string, data ...any) CustomGomegaMatcher {
	return c.WithPrecompiledTemplate(template.Must(ParseTemplate(templ)), data...)
}

func (c CustomGomegaMatcher) WithPrecompiledTemplate(templ *template.Template, data ...any) CustomGomegaMatcher {
	c.templateMessage = templ
	c.templateData = nil
	if len(data) > 0 {
		c.templateData = data[0]
	}
	return c
}

func (c CustomGomegaMatcher) WithTemplateData(data any) CustomGomegaMatcher {
	c.templateData = data
	return c
}

func (c CustomGomegaMatcher) Match(actual any) (bool, error) {
	return c.matchFunc(actual)
}

func (c CustomGomegaMatcher) FailureMessage(actual any) string {
	return c.renderTemplateMessage(actual, true)
}

func (c CustomGomegaMatcher) NegatedFailureMessage(actual any) string {
	return c.renderTemplateMessage(actual, false)
}

type templateData struct {
	Failure         bool
	NegatedFailure  bool
	To              string
	FormattedActual string
	Actual          any
	Data            any
}

func (c CustomGomegaMatcher) renderTemplateMessage(actual any, isFailure bool) string {
	var data templateData
	formattedActual := format.Object(actual, 1)
	if isFailure {
		data = templateData{
			Failure:         true,
			NegatedFailure:  false,
			To:              "to",
			FormattedActual: formattedActual,
			Actual:          actual,
			Data:            c.templateData,
		}
	} else {
		data = templateData{
			Failure:         false,
			NegatedFailure:  true,
			To:              "not to",
			FormattedActual: formattedActual,
			Actual:          actual,
			Data:            c.templateData,
		}
	}
	b := &strings.Builder{}
	err := c.templateMessage.Execute(b, data)
	if err != nil {
		return fmt.Sprintf("Failed to render failure message template: %s", err.Error())
	}
	return b.String()
}
