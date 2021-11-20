package matchers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/matchers"
)

type I interface {
	M()
}

type S struct {
	V int
}

func (s S) M() {}

var _ = Describe("HaveValue", func() {

	It("should fail when passed nil", func() {
		var p *struct{}
		m := HaveValue(BeNil())
		Expect(m.Match(p)).To(BeFalse())
		Expect(m.FailureMessage(p)).To(MatchRegexp("not to be <nil>$"))
	})

	It("should fail when passed nil indirectly", func() {
		var p *struct{}
		m := HaveValue(BeNil())
		Expect(m.Match(&p)).To(BeFalse())
		Expect(m.FailureMessage(&p)).To(MatchRegexp("not to be <nil>$"))
	})

	It("should unwrap the value pointed to", func() {
		i := 1
		Expect(&i).To(HaveValue(Equal(1)))
		Expect(&i).NotTo(HaveValue(Equal(2)))

		p := &i
		Expect(&p).To(HaveValue(Equal(1)))
		Expect(&p).NotTo(HaveValue(Equal(2)))
	})

	It("should unwrap the value inside an interface", func() {
		var i I = &S{V: 42}
		Expect(i).To(HaveValue(Equal(S{V: 42})))
		Expect(i).NotTo(HaveValue(Equal(S{})))
	})

})
