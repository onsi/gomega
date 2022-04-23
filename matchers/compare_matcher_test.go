package matchers_test

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/matchers"
)

type wrapError struct {
	msg string
	err error
}

func (e wrapError) Error() string {
	return e.msg
}

func (e wrapError) Unwrap() error {
	return e.err
}

var _ = Describe("Compare", func() {
	When("asserting that nil compare equal nil", func() {
		It("should error", func() {
			success, err := (&CompareMatcher{Expected: nil}).Match(nil)

			Expect(success).Should(BeFalse())
			Expect(err).Should(HaveOccurred())
		})
	})

	Context("When asserting on nil", func() {
		It("should do the right thing", func() {
			Expect("foo").ShouldNot(Compare(nil))
			Expect(nil).ShouldNot(Compare(3))
			Expect([]int{1, 2}).ShouldNot(Compare(nil))
		})
	})

	Context("When asserting time with different location ", func() {
		var t1, t2 time.Time

		BeforeEach(func() {
			t1 = time.Time{}
			t2 = time.Time{}.Local()
		})

		It("should do the right thing", func() {
			Expect(t1).ShouldNot(Equal(t2))
			Expect(t1).Should(Compare(t2))
		})
	})

	Context("When struct contain unexported fields", func() {
		type structWithUnexportedFields struct {
			unexported string
			Exported   string
		}

		var s1, s2 structWithUnexportedFields

		BeforeEach(func() {
			s1 = structWithUnexportedFields{unexported: "unexported", Exported: "Exported"}
			s2 = structWithUnexportedFields{unexported: "unexported", Exported: "Exported"}
		})

		It("should panic with unexported field", func() {
			defer func() {
				if e := recover(); e != nil {
					Expect(e).Should(HavePrefix("cannot handle unexported field at"))
				}
			}()

			matcher := &CompareMatcher{
				Expected: s1,
			}
			_, _ = matcher.Match(s2)
		})

		It("should do the right thing", func() {
			Expect(s1).Should(Compare(s2, cmpopts.IgnoreUnexported(structWithUnexportedFields{})))
		})
	})

	Context("When compare error", func() {
		var err1, err2 error

		It("not equal", func() {
			err1 = errors.New("error")
			err2 = errors.New("error")
			Expect(err1).ShouldNot(Compare(err2, cmpopts.EquateErrors()))
		})

		It("equal if err1 is err2", func() {
			err1 = errors.New("error")
			err2 = &wrapError{
				msg: "some error",
				err: err1,
			}

			Expect(err1).Should(Compare(err2, cmpopts.EquateErrors()))
		})
	})

	Context("When asserting equal between objects", func() {
		It("should do the right thing", func() {
			Expect(5).Should(Compare(5))
			Expect(5.0).Should(Compare(5.0))

			Expect(5).ShouldNot(Compare("5"))
			Expect(5).ShouldNot(Compare(5.0))
			Expect(5).ShouldNot(Compare(3))

			Expect("5").Should(Compare("5"))
			Expect([]int{1, 2}).Should(Compare([]int{1, 2}))
			Expect([]int{1, 2}).ShouldNot(Compare([]int{2, 1}))
			Expect([]byte{'f', 'o', 'o'}).Should(Compare([]byte{'f', 'o', 'o'}))
			Expect([]byte{'f', 'o', 'o'}).ShouldNot(Compare([]byte{'b', 'a', 'r'}))
			Expect(map[string]string{"a": "b", "c": "d"}).Should(Compare(map[string]string{"a": "b", "c": "d"}))
			Expect(map[string]string{"a": "b", "c": "d"}).ShouldNot(Compare(map[string]string{"a": "b", "c": "e"}))

			Expect(myCustomType{s: "abc", n: 3, f: 2.0, arr: []string{"a", "b"}}).Should(Compare(myCustomType{s: "foo", n: 3, f: 2.0, arr: []string{"a", "b"}}, cmpopts.IgnoreUnexported(myCustomType{})))

			Expect(myCustomType{s: "foo", n: 3, f: 2.0, arr: []string{"a", "b"}}).Should(Compare(myCustomType{s: "foo", n: 3, f: 2.0, arr: []string{"a", "b"}}, cmp.AllowUnexported(myCustomType{})))
			Expect(myCustomType{s: "foo", n: 3, f: 2.0, arr: []string{"a", "b"}}).ShouldNot(Compare(myCustomType{s: "bar", n: 3, f: 2.0, arr: []string{"a", "b"}}, cmp.AllowUnexported(myCustomType{})))
			Expect(myCustomType{s: "foo", n: 3, f: 2.0, arr: []string{"a", "b"}}).ShouldNot(Compare(myCustomType{s: "foo", n: 2, f: 2.0, arr: []string{"a", "b"}}, cmp.AllowUnexported(myCustomType{})))
			Expect(myCustomType{s: "foo", n: 3, f: 2.0, arr: []string{"a", "b"}}).ShouldNot(Compare(myCustomType{s: "foo", n: 3, f: 3.0, arr: []string{"a", "b"}}, cmp.AllowUnexported(myCustomType{})))
			Expect(myCustomType{s: "foo", n: 3, f: 2.0, arr: []string{"a", "b"}}).ShouldNot(Compare(myCustomType{s: "foo", n: 3, f: 2.0, arr: []string{"a", "b", "c"}}, cmp.AllowUnexported(myCustomType{})))
		})
	})
})
