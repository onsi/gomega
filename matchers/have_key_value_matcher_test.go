package matchers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/matchers"
)

var _ = Describe("HaveKeyValue", func() {
	var (
		stringKeys map[string]int
		intKeys    map[int]string
		objKeys    map[*myCustomType]*myCustomType

		customA *myCustomType
		customB *myCustomType
	)
	BeforeEach(func() {
		stringKeys = map[string]int{"foo": 2, "bar": 3}
		intKeys = map[int]string{2: "foo", 3: "bar"}

		customA = &myCustomType{s: "a", n: 2, f: 2.3, arr: []string{"ice", "cream"}}
		customB = &myCustomType{s: "b", n: 4, f: 3.1, arr: []string{"cake"}}
		objKeys = map[*myCustomType]*myCustomType{customA: customA, customB: customA}
	})

	Context("when passed a map", func() {
		It("should do the right thing", func() {
			Ω(stringKeys).Should(HaveKeyValue("foo", 2))
			Ω(stringKeys).ShouldNot(HaveKeyValue("foo", 1))
			Ω(stringKeys).ShouldNot(HaveKeyValue("baz", 2))
			Ω(stringKeys).ShouldNot(HaveKeyValue("baz", 1))

			Ω(intKeys).Should(HaveKeyValue(2, "foo"))
			Ω(intKeys).ShouldNot(HaveKeyValue(4, "foo"))
			Ω(intKeys).ShouldNot(HaveKeyValue(2, "baz"))

			Ω(objKeys).Should(HaveKeyValue(customA, customA))
			Ω(objKeys).Should(HaveKeyValue(&myCustomType{s: "b", n: 4, f: 3.1, arr: []string{"cake"}}, &myCustomType{s: "a", n: 2, f: 2.3, arr: []string{"ice", "cream"}}))
			Ω(objKeys).ShouldNot(HaveKeyValue(&myCustomType{s: "b", n: 4, f: 3.1, arr: []string{"apple", "pie"}}, customA))
		})
	})

	Context("when passed a correctly typed nil", func() {
		It("should operate succesfully on the passed in value", func() {
			var nilMap map[int]string
			Ω(nilMap).ShouldNot(HaveKeyValue("foo", "bar"))
		})
	})

	Context("when the passed in key is actually a matcher", func() {
		It("should pass each element through the matcher", func() {
			Ω(stringKeys).Should(HaveKeyValue(ContainSubstring("oo"), 2))
			Ω(intKeys).Should(HaveKeyValue(2, ContainSubstring("oo")))
			Ω(stringKeys).ShouldNot(HaveKeyValue(ContainSubstring("foobar"), 2))
		})

		It("should fail if the matcher ever fails", func() {
			actual := map[interface{}]string{"foo": "a", 3: "b", "bar": "c"}
			success, err := (&HaveKeyValueMatcher{Key: ContainSubstring("ar"), Value: 2}).Match(actual)
			Ω(success).Should(BeFalse())
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("when passed something that is not a map", func() {
		It("should error", func() {
			success, err := (&HaveKeyValueMatcher{Key: "foo", Value: "bar"}).Match([]string{"foo"})
			Ω(success).Should(BeFalse())
			Ω(err).Should(HaveOccurred())

			success, err = (&HaveKeyValueMatcher{Key: "foo", Value: "bar"}).Match(nil)
			Ω(success).Should(BeFalse())
			Ω(err).Should(HaveOccurred())
		})
	})
})
