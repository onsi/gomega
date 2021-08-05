package internal_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Making Synchronous Assertions", func() {
	var SHOULD_MATCH = true
	var SHOULD_NOT_MATCH = false
	var IT_PASSES = true
	var IT_FAILS = false

	Extras := func(extras ...interface{}) []interface{} {
		return extras
	}

	OptionalDescription := func(optionalDescription ...interface{}) []interface{} {
		return optionalDescription
	}

	DescribeTable(
		"the various cases",
		func(actual interface{}, extras []interface{}, optionalDescription []interface{}, isPositiveAssertion bool, expectedFailureMessage string, expectedReturnValue bool) {
			if isPositiveAssertion {
				ig := NewInstrumentedGomega()
				returnValue := ig.G.Expect(actual, extras...).To(SpecMatch(), optionalDescription...)
				Expect(returnValue).To(Equal(expectedReturnValue))
				Expect(ig.FailureMessage).To(Equal(expectedFailureMessage))
				if expectedFailureMessage != "" {
					Expect(ig.FailureSkip).To(Equal([]int{2}))
				}
				Expect(ig.RegisteredHelpers).To(ContainElement("(*Assertion).To"))

				ig = NewInstrumentedGomega()
				returnValue = ig.G.ExpectWithOffset(3, actual, extras...).To(SpecMatch(), optionalDescription...)
				Expect(returnValue).To(Equal(expectedReturnValue))
				Expect(ig.FailureMessage).To(Equal(expectedFailureMessage))
				if expectedFailureMessage != "" {
					Expect(ig.FailureSkip).To(Equal([]int{5}))
				}
				Expect(ig.RegisteredHelpers).To(ContainElement("(*Assertion).To"))

				ig = NewInstrumentedGomega()
				returnValue = ig.G.Ω(actual, extras...).Should(SpecMatch(), optionalDescription...)
				Expect(returnValue).To(Equal(expectedReturnValue))
				Expect(ig.FailureMessage).To(Equal(expectedFailureMessage))
				if expectedFailureMessage != "" {
					Expect(ig.FailureSkip).To(Equal([]int{2}))
				}
				Expect(ig.RegisteredHelpers).To(ContainElement("(*Assertion).Should"))
			} else {
				ig := NewInstrumentedGomega()
				returnValue := ig.G.Expect(actual, extras...).ToNot(SpecMatch(), optionalDescription...)
				Expect(returnValue).To(Equal(expectedReturnValue))
				Expect(ig.FailureMessage).To(Equal(expectedFailureMessage))
				if expectedFailureMessage != "" {
					Expect(ig.FailureSkip).To(Equal([]int{2}))
				}
				Expect(ig.RegisteredHelpers).To(ContainElement("(*Assertion).ToNot"))

				ig = NewInstrumentedGomega()
				returnValue = ig.G.Expect(actual, extras...).NotTo(SpecMatch(), optionalDescription...)
				Expect(returnValue).To(Equal(expectedReturnValue))
				Expect(ig.FailureMessage).To(Equal(expectedFailureMessage))
				if expectedFailureMessage != "" {
					Expect(ig.FailureSkip).To(Equal([]int{2}))
				}
				Expect(ig.RegisteredHelpers).To(ContainElement("(*Assertion).NotTo"))

				ig = NewInstrumentedGomega()
				returnValue = ig.G.ExpectWithOffset(3, actual, extras...).NotTo(SpecMatch(), optionalDescription...)
				Expect(returnValue).To(Equal(expectedReturnValue))
				Expect(ig.FailureMessage).To(Equal(expectedFailureMessage))
				if expectedFailureMessage != "" {
					Expect(ig.FailureSkip).To(Equal([]int{5}))
				}
				Expect(ig.RegisteredHelpers).To(ContainElement("(*Assertion).NotTo"))

				ig = NewInstrumentedGomega()
				returnValue = ig.G.Ω(actual, extras...).ShouldNot(SpecMatch(), optionalDescription...)
				Expect(returnValue).To(Equal(expectedReturnValue))
				Expect(ig.FailureMessage).To(Equal(expectedFailureMessage))
				if expectedFailureMessage != "" {
					Expect(ig.FailureSkip).To(Equal([]int{2}))
				}
				Expect(ig.RegisteredHelpers).To(ContainElement("(*Assertion).ShouldNot"))
			}
		},
		Entry(
			"when the matcher matches and a positive assertion is being made",
			MATCH, Extras(), OptionalDescription(),
			SHOULD_MATCH, "", IT_PASSES,
		),
		Entry(
			"when the matcher matches and a negative assertion is being made",
			MATCH, Extras(), OptionalDescription(),
			SHOULD_NOT_MATCH, "negative: match", IT_FAILS,
		),
		Entry(
			"when the matcher does not match and a positive assertion is being made",
			NO_MATCH, Extras(), OptionalDescription(),
			SHOULD_MATCH, "positive: no match", IT_FAILS,
		),
		Entry(
			"when the matcher does not match and a negative assertion is being made",
			NO_MATCH, Extras(), OptionalDescription(),
			SHOULD_NOT_MATCH, "", IT_PASSES,
		),
		Entry(
			"when the matcher returns an error and a positive assertion is being made",
			ERR_MATCH, Extras(), OptionalDescription(),
			SHOULD_MATCH, "spec matcher error", IT_FAILS,
		),
		Entry(
			"when the matcher returns an error and a negative assertion is being made",
			ERR_MATCH, Extras(), OptionalDescription(),
			SHOULD_NOT_MATCH, "spec matcher error", IT_FAILS,
		),
		Entry(
			"when a failure occurs and there is a single optional description",
			NO_MATCH, Extras(), OptionalDescription("a description"),
			SHOULD_MATCH, "a description\npositive: no match", IT_FAILS,
		),
		Entry(
			"when a failure occurs and there are multiple optional descriptions",
			NO_MATCH, Extras(), OptionalDescription("a description of [%d]", 3),
			SHOULD_MATCH, "a description of [3]\npositive: no match", IT_FAILS,
		),
		Entry(
			"when a failure occurs and the optional description is a function",
			NO_MATCH, Extras(), OptionalDescription(func() string { return "a description" }),
			SHOULD_MATCH, "a description\npositive: no match", IT_FAILS,
		),
		Entry(
			"when the matcher matches and zero-valued extra parameters are included, it passes",
			MATCH, Extras(0, "", struct{ Foo string }{}, nil), OptionalDescription(),
			SHOULD_MATCH, "", IT_PASSES,
		),
		Entry(
			"when the matcher matches but a non-zero-valued extra parameter is included, it fails",
			MATCH, Extras(1, "bam", struct{ Foo string }{Foo: "foo"}, nil), OptionalDescription(),
			SHOULD_MATCH, "Unexpected non-nil/non-zero extra argument at index 1:\n\t<int>: 1", IT_FAILS,
		),
	)
})
