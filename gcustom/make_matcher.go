/*
package gcustom provides a simple mechanism for creating custom Gomega matchers
*/
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

/*
ParseTemplate allows you to precompile templates for MakeMatcher's custom matchers.

Use ParseTemplate if you are concerned about performance and would like to avoid repeatedly parsing failure message templates.  The data made available to the template is documented in the WithTemplate() method of CustomGomegaMatcher.

Once parsed you can pass the template in either as an argument to MakeMatcher(matchFunc, <template>) or using MakeMatcher(matchFunc).WithPrecompiledTemplate(template)
*/
func ParseTemplate(templ string) (*template.Template, error) {
	return template.New("template").Funcs(template.FuncMap{
		"format": formatObject,
	}).Parse(templ)
}

/*
MakeMatcher builds a Gomega-compatible matcher from a function (the matchFunc).

matchFunc must return (bool, error) and take a single argument.  If you want to perform type-checking yourself pass in a matchFunc of type `func(actual any) (bool, error)`.  If you want to only operate on a specific type, pass in `func(actual DesiredType) (bool, error)`; MakeMatcher will take care of checking types for you and notifying the user if they use the matcher with an invalid type.

MakeMatcher(matchFunc) builds a matcher with generic failure messages that look like:

	Custom matcher failed for:
	    <formatted actual>

for the positive failure case (i.e. when Expect(actual).To(match) fails) and

	Custom matcher succeeded (but was expected to fail) for:
	     <formatted actual>

for the negative case (i.e. when Expect(actual).NotTo(match) fails).

There are two ways to provide a different message.  You can either provide a simple message string:

	matcher := MakeMatcher(matchFunc, message)
	matcher := MakeMatcher(matchFunc).WithMessage(message)

(where message is of type string) or a template:

	matcher := MakeMatcher(matchFunc).WithTemplate(templateString)

where templateString is a string that is compiled by WithTemplate into a matcher.  Alternatively you can provide a precompiled template like this:

	template, err = gcustom.ParseTemplate(templateString) //use gcustom's ParseTemplate to get some additional functions mixed in
	matcher := MakeMatcher(matchFunc, template)
	matcher := MakeMatcher(matchFunc).WithPrecompiled(template)

When a simple message string is provided the positive failure message will look like:

	Expected:
		<formatted actual>
	to <message>

and the negative failure message will look like:

	Expected:
		<formatted actual>
	not to <message>

A template allows you to have greater control over the message.  For more details see the docs for WithTemplate
*/
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
				if actual == nil && reflect.TypeOf(actual) == reflect.TypeOf(nil) {
					return matchFuncValue.Call([]reflect.Value{reflect.New(t.In(0)).Elem()})
				} else if reflect.TypeOf(actual).AssignableTo(t.In(0)) {
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

// CustomGomegaMatcher is generated by MakeMatcher - you should always use MakeMatcher to construct custom matchers
type CustomGomegaMatcher struct {
	matchFunc                   func(actual any) (bool, error)
	templateMessage             *template.Template
	templateData                any
	customFailureMessage        func(actual any) string
	customNegatedFailureMessage func(actual any) string
}

/*
WithMessage returns a CustomGomegaMatcher configured with a message to display when failure occurs.  Matchers configured this way produce a positive failure message that looks like:

	Expected:
		<formatted actual>
	to <message>

and a negative failure message that looks like:

	Expected:
		<formatted actual>
	not to <message>
*/
func (c CustomGomegaMatcher) WithMessage(message string) CustomGomegaMatcher {
	return c.WithTemplate("Expected:\n{{.FormattedActual}}\n{{.To}} " + message)
}

/*
WithTemplate compiles the passed-in template and returns a CustomGomegaMatcher configured to use that template to generate failure messages.

Templates are provided the following variables and functions:

{{.Failure}} - a bool that, if true, indicates this should be a positive failure message, otherwise this should be a negated failure message
{{.NegatedFailure}} - a bool that, if true, indicates this should be a negated failure message, otherwise this should be a positive failure message
{{.To}} - is set to "to" if this is a positive failure message and "not to" if this is a negated failure message
{{.Actual}} - the actual passed in to the matcher
{{.FormattedActual}} - a string representing the formatted actual.  This can be multiple lines and is always generated with an indentation of 1
{{format <object> <optional-indentation}} - a function that allows you to use Gomega's default formatting from within the template.  The passed-in <object> is formatted and <optional-indentation> can be set to an integer to control indentation.

In addition, you can provide custom data to the template by calling WithTemplate(templateString, data) (where data can be anything).  This is provided to the template as {{.Data}}.

Here's a simple example of all these pieces working together:

	func HaveWidget(widget Widget) OmegaMatcher {
		return MakeMatcher(func(machine Machine) (bool, error) {
			return machine.HasWidget(widget), nil
		}).WithTemplate("Expected:\n{{.FormattedActual}}\n{{.To}} have widget named {{.Data.Name}}:\n{{format .Data 1}}", widget)
	}

	Expect(machine).To(HaveWidget(Widget{Name: "sprocket", Version: 2}))

Would generate a failure message that looks like:

	Expected:
		<formatted machine>
	to have widget named sprocket:
		<formatted sprocket>
*/
func (c CustomGomegaMatcher) WithTemplate(templ string, data ...any) CustomGomegaMatcher {
	return c.WithPrecompiledTemplate(template.Must(ParseTemplate(templ)), data...)
}

/*
WithPrecompiledTemplate returns a CustomGomegaMatcher configured to use the passed-in template.  The template should be precompiled with gcustom.ParseTemplate().

As with WithTemplate() you can provide a single piece of additional data as an optional argument.  This is accessed in the template via {{.Data}}
*/
func (c CustomGomegaMatcher) WithPrecompiledTemplate(templ *template.Template, data ...any) CustomGomegaMatcher {
	c.templateMessage = templ
	c.templateData = nil
	if len(data) > 0 {
		c.templateData = data[0]
	}
	return c
}

/*
WithTemplateData() returns a CustomGomegaMatcher configured to provide it's template with the passed-in data.  The following are equivalent:

MakeMatcher(matchFunc).WithTemplate(templateString, data)
MakeMatcher(matchFunc).WithTemplate(templateString).WithTemplateData(data)
*/
func (c CustomGomegaMatcher) WithTemplateData(data any) CustomGomegaMatcher {
	c.templateData = data
	return c
}

// Match runs the passed-in match func and satisfies the GomegaMatcher interface
func (c CustomGomegaMatcher) Match(actual any) (bool, error) {
	return c.matchFunc(actual)
}

// FailureMessage generates the positive failure message configured via WithMessage or WithTemplate/WithPrecompiledTemplate
// i.e. this is the failure message when Expect(actual).To(match) fails
func (c CustomGomegaMatcher) FailureMessage(actual any) string {
	return c.renderTemplateMessage(actual, true)
}

// NegatedFailureMessage generates the negative failure message configured via WithMessage or WithTemplate/WithPrecompiledTemplate
// i.e. this is the failure message when Expect(actual).NotTo(match) fails
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
