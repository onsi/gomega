package matchers_test

import (
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/matchers"
)

type FakeIsError struct {
	isError bool
}

func (f *FakeIsError) Error() string {
	return fmt.Sprintf("is other error: %T", f.isError)
}

func (f *FakeIsError) Is(other error) bool {
	return f.isError
}

var _ = Describe("MatchErrorStrictlyMatcher", func() {
	Context("When asserting against an error", func() {
		When("passed an error", func() {
			It("should succeed when errors.Is returns true", func() {
				err := errors.New("an error")
				fmtErr := fmt.Errorf("an error")
				isError := &FakeIsError{true}

				Expect(err).To(MatchErrorStrictly(err))
				Expect(fmtErr).To(MatchErrorStrictly(fmtErr))
				Expect(isError).To(MatchErrorStrictly(errors.New("any error should match")))
			})

			It("should fail when errors.Is returns false", func() {
				err := errors.New("an error")
				fmtErr := fmt.Errorf("an error")
				isNotError := &FakeIsError{false}

				Expect(err).ToNot(MatchErrorStrictly(errors.New("another error")))
				Expect(fmtErr).ToNot(MatchErrorStrictly(fmt.Errorf("an error")))

				// errors.Is first checks if the values equal via ==, so we must point
				// to different instances of otherwise equal FakeIsError
				Expect(isNotError).ToNot(MatchErrorStrictly(&FakeIsError{false}))
			})

			It("should succeed when any error in the chain matches the passed error", func() {
				innerErr := errors.New("inner error")
				outerErr := fmt.Errorf("outer error wrapping: %w", innerErr)

				Expect(outerErr).To(MatchErrorStrictly(innerErr))
			})
		})
	})

	When("expected is nil", func() {
		It("should fail with an appropriate error", func() {
			_, err := (&MatchErrorStrictlyMatcher{
				Expected: nil,
			}).Match(errors.New("an error"))
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("ToNot(HaveOccurred())"))
		})
	})

	When("passed nil", func() {
		It("should fail", func() {
			_, err := (&MatchErrorStrictlyMatcher{
				Expected: errors.New("an error"),
			}).Match(nil)
			Expect(err).To(HaveOccurred())
		})
	})

	When("passed a non-error", func() {
		It("should fail", func() {
			_, err := (&MatchErrorStrictlyMatcher{
				Expected: errors.New("an error"),
			}).Match("an error")
			Expect(err).To(HaveOccurred())

			_, err = (&MatchErrorStrictlyMatcher{
				Expected: errors.New("an error"),
			}).Match(3)
			Expect(err).To(HaveOccurred())
		})
	})

	It("shows failure message", func() {
		failuresMessages := InterceptGomegaFailures(func() {
			Expect(errors.New("foo")).To(MatchErrorStrictly(errors.New("bar")))
		})
		Expect(failuresMessages[0]).To(ContainSubstring("foo\n    {s: \"foo\"}\nto match error\n    <*errors.errorString"))
	})

	It("shows negated failure message", func() {
		err := errors.New("foo")
		failuresMessages := InterceptGomegaFailures(func() {
			Expect(err).ToNot(MatchErrorStrictly(err))
		})
		Expect(failuresMessages[0]).To(ContainSubstring("foo\n    {s: \"foo\"}\nnot to match error\n    <*errors.errorString"))
	})

})
