package matcher_tests

import (
	"errors"
	. "github.com/onsi/godescribe"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/matchers"
)

type myStringer struct {
	a string
}

func (s *myStringer) String() string {
	return s.a
}

type myCustomType struct {
	s   string
	n   int
	f   float32
	arr []string
}

func init() {
	Describe("Equal", func() {
		Context("when asserting that nil equals nil", func() {
			It("should error", func() {
				success, _, err := (&EqualMatcher{Expected: nil}).Match(nil)

				Ω(success).Should(BeFalse())
				Ω(err).Should(HaveOccured())
			})
		})

		Context("When asserting equality between objects", func() {
			It("should do the right thing", func() {
				Ω(5).Should(Equal(5))
				Ω(5.0).Should(Equal(5.0))

				Ω(5).ShouldNot(Equal("5"))
				Ω(5).ShouldNot(Equal(5.0))
				Ω(5).ShouldNot(Equal(3))

				Ω("5").Should(Equal("5"))
				Ω([]int{1, 2}).Should(Equal([]int{1, 2}))
				Ω([]int{1, 2}).ShouldNot(Equal([]int{2, 1}))
				Ω(map[string]string{"a": "b", "c": "d"}).Should(Equal(map[string]string{"a": "b", "c": "d"}))
				Ω(map[string]string{"a": "b", "c": "d"}).ShouldNot(Equal(map[string]string{"a": "b", "c": "e"}))
				Ω(errors.New("foo")).Should(Equal(errors.New("foo")))
				Ω(errors.New("foo")).ShouldNot(Equal(errors.New("bar")))

				Ω(myCustomType{s: "foo", n: 3, f: 2.0, arr: []string{"a", "b"}}).Should(Equal(myCustomType{s: "foo", n: 3, f: 2.0, arr: []string{"a", "b"}}))
				Ω(myCustomType{s: "foo", n: 3, f: 2.0, arr: []string{"a", "b"}}).ShouldNot(Equal(myCustomType{s: "bar", n: 3, f: 2.0, arr: []string{"a", "b"}}))
				Ω(myCustomType{s: "foo", n: 3, f: 2.0, arr: []string{"a", "b"}}).ShouldNot(Equal(myCustomType{s: "foo", n: 2, f: 2.0, arr: []string{"a", "b"}}))
				Ω(myCustomType{s: "foo", n: 3, f: 2.0, arr: []string{"a", "b"}}).ShouldNot(Equal(myCustomType{s: "foo", n: 3, f: 3.0, arr: []string{"a", "b"}}))
				Ω(myCustomType{s: "foo", n: 3, f: 2.0, arr: []string{"a", "b"}}).ShouldNot(Equal(myCustomType{s: "foo", n: 3, f: 2.0, arr: []string{"a", "b", "c"}}))
			})
		})
	})

	Describe("BeNil", func() {
		It("should succeed when passed nil", func() {
			Ω(nil).Should(BeNil())
		})

		It("should not succeed when not passed nil", func() {
			Ω(0).ShouldNot(BeNil())
			Ω(false).ShouldNot(BeNil())
			Ω("").ShouldNot(BeNil())
		})
	})

	Describe("BeTrue", func() {
		It("should handle true and false correctly", func() {
			Ω(true).Should(BeTrue())
			Ω(false).ShouldNot(BeTrue())
		})

		It("should only support booleans", func() {
			success, _, err := (&BeTrueMatcher{}).Match("foo")
			Ω(success).Should(BeFalse())
			Ω(err).Should(HaveOccured())
		})
	})

	Describe("BeFalse", func() {
		It("should handle true and false correctly", func() {
			Ω(true).ShouldNot(BeFalse())
			Ω(false).Should(BeFalse())
		})

		It("should only support booleans", func() {
			success, _, err := (&BeFalseMatcher{}).Match("foo")
			Ω(success).Should(BeFalse())
			Ω(err).Should(HaveOccured())
		})
	})

	Describe("HaveOccured", func() {
		It("should succeed if matching an error", func() {
			Ω(errors.New("Foo")).Should(HaveOccured())
		})

		It("should not succed with nil", func() {
			Ω(nil).ShouldNot(HaveOccured())
		})

		It("should only support errors and nil", func() {
			success, _, err := (&HaveOccuredMatcher{}).Match("foo")
			Ω(success).Should(BeFalse())
			Ω(err).Should(HaveOccured())

			success, _, err = (&HaveOccuredMatcher{}).Match("")
			Ω(success).Should(BeFalse())
			Ω(err).Should(HaveOccured())
		})
	})

	Describe("MatchRegexp", func() {
		Context("when actual is a string", func() {
			It("should match against the string", func() {
				Ω(" a2!bla").Should(MatchRegexp(`\d!`))
				Ω(" a2!bla").ShouldNot(MatchRegexp(`[A-Z]`))
			})
		})

		Context("when actual is a stringer", func() {
			It("should call the stringer and match agains the returned string", func() {
				Ω(&myStringer{a: "Abc3"}).Should(MatchRegexp(`[A-Z][a-z]+\d`))
			})
		})

		Context("when actual is neither a string nor a stringer", func() {
			It("should error", func() {
				success, _, err := (&MatchRegexpMatcher{Regexp: `\d`}).Match(2)
				Ω(success).Should(BeFalse())
				Ω(err).Should(HaveOccured())
			})
		})

		Context("when the passed in regexp fails to compile", func() {
			It("should error", func() {
				success, _, err := (&MatchRegexpMatcher{Regexp: "("}).Match("Foo")
				Ω(success).Should(BeFalse())
				Ω(err).Should(HaveOccured())
			})
		})
	})

	Describe("ContainSubstringMatcher", func() {
		Context("when actual is a string", func() {
			It("should match against the string", func() {
				Ω("Marvelous").Should(ContainSubstring("rve"))
				Ω("Marvelous").ShouldNot(ContainSubstring("boo"))
			})
		})

		Context("when actual is a stringer", func() {
			It("should call the stringer and match agains the returned string", func() {
				Ω(&myStringer{a: "Abc3"}).Should(ContainSubstring("bc3"))
			})
		})

		Context("when actual is neither a string nor a stringer", func() {
			It("should error", func() {
				success, _, err := (&ContainSubstringMatcher{Substr: "2"}).Match(2)
				Ω(success).Should(BeFalse())
				Ω(err).Should(HaveOccured())
			})
		})
	})

	Describe("BeEmpty", func() {
		Context("when passed a supported type", func() {
			It("should do the right thing", func() {
				Ω("").Should(BeEmpty())
				Ω(" ").ShouldNot(BeEmpty())

				Ω([0]int{}).Should(BeEmpty())
				Ω([1]int{1}).ShouldNot(BeEmpty())

				Ω([]int{}).Should(BeEmpty())
				Ω([]int{1}).ShouldNot(BeEmpty())

				Ω(map[string]int{}).Should(BeEmpty())
				Ω(map[string]int{"a": 1}).ShouldNot(BeEmpty())

				c := make(chan bool, 1)
				Ω(c).Should(BeEmpty())
				c <- true
				Ω(c).ShouldNot(BeEmpty())
			})
		})

		Context("when passed an unsupported type", func() {
			It("should error", func() {
				success, _, err := (&BeEmptyMatcher{}).Match(0)
				Ω(success).Should(BeFalse())
				Ω(err).Should(HaveOccured())

				success, _, err = (&BeEmptyMatcher{}).Match(nil)
				Ω(success).Should(BeFalse())
				Ω(err).Should(HaveOccured())
			})
		})
	})

	Describe("HaveLen", func() {
		Context("when passed a supported type", func() {
			It("should do the right thing", func() {
				Ω("").Should(HaveLen(0))
				Ω("AA").Should(HaveLen(2))

				Ω([0]int{}).Should(HaveLen(0))
				Ω([2]int{1, 2}).Should(HaveLen(2))

				Ω([]int{}).Should(HaveLen(0))
				Ω([]int{1, 2, 3}).Should(HaveLen(3))

				Ω(map[string]int{}).Should(HaveLen(0))
				Ω(map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}).Should(HaveLen(4))

				c := make(chan bool, 3)
				Ω(c).Should(HaveLen(0))
				c <- true
				c <- true
				Ω(c).Should(HaveLen(2))
			})
		})

		Context("when passed an unsupported type", func() {
			It("should error", func() {
				success, _, err := (&HaveLenMatcher{Count: 0}).Match(0)
				Ω(success).Should(BeFalse())
				Ω(err).Should(HaveOccured())

				success, _, err = (&HaveLenMatcher{Count: 0}).Match(nil)
				Ω(success).Should(BeFalse())
				Ω(err).Should(HaveOccured())
			})
		})
	})

	Describe("BeZero", func() {
		It("should succeed if the passed in object is the zero value for its type", func() {
			Ω("").Should(BeZero())
			Ω(" ").ShouldNot(BeZero())

			Ω(0).Should(BeZero())
			Ω(1).ShouldNot(BeZero())

			Ω(0.0).Should(BeZero())
			Ω(0.1).ShouldNot(BeZero())

			// Ω([]int{}).Should(BeZero())
			Ω([]int{1}).ShouldNot(BeZero())

			// Ω(map[string]int{}).Should(BeZero())
			Ω(map[string]int{"a": 1}).ShouldNot(BeZero())

			Ω(myCustomType{}).Should(BeZero())
			Ω(myCustomType{s: "a"}).ShouldNot(BeZero())
		})
	})

	Describe("ContainElement", func() {
		Context("when passed a supported type", func() {
			It("should do the right thing", func() {
				Ω([2]int{1, 2}).Should(ContainElement(2))
				Ω([2]int{1, 2}).ShouldNot(ContainElement(3))

				Ω([]int{1, 2}).Should(ContainElement(2))
				Ω([]int{1, 2}).ShouldNot(ContainElement(3))

				arr := make([]myCustomType, 2)
				arr[0] = myCustomType{s: "foo", n: 3, f: 2.0, arr: []string{"a", "b"}}
				arr[1] = myCustomType{s: "foo", n: 3, f: 2.0, arr: []string{"a", "c"}}
				Ω(arr).Should(ContainElement(myCustomType{s: "foo", n: 3, f: 2.0, arr: []string{"a", "b"}}))
				Ω(arr).ShouldNot(ContainElement(myCustomType{s: "foo", n: 3, f: 2.0, arr: []string{"b", "c"}}))
			})
		})

		Context("when passed an unsupported type", func() {
			It("should error", func() {
				success, _, err := (&ContainElementMatcher{Element: 0}).Match(0)
				Ω(success).Should(BeFalse())
				Ω(err).Should(HaveOccured())

				success, _, err = (&ContainElementMatcher{Element: 0}).Match("abc")
				Ω(success).Should(BeFalse())
				Ω(err).Should(HaveOccured())

				success, _, err = (&ContainElementMatcher{Element: 0}).Match(map[string]int{"a": 1, "b": 2})
				Ω(success).Should(BeFalse())
				Ω(err).Should(HaveOccured())

				success, _, err = (&ContainElementMatcher{Element: 0}).Match(nil)
				Ω(success).Should(BeFalse())
				Ω(err).Should(HaveOccured())
			})
		})
	})
}
