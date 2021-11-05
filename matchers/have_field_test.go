package matchers_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

type Book struct {
	Title  string
	Author person
	Pages  int
}

func (book Book) AuthorName() string {
	return fmt.Sprintf("%s %s", book.Author.FirstName, book.Author.LastName)
}

func (book Book) AbbreviatedAuthor() person {
	return person{
		FirstName: book.Author.FirstName[0:3],
		LastName:  book.Author.LastName[0:3],
		DOB:       book.Author.DOB,
	}
}

func (book Book) NoReturn() {
}

func (book Book) TooManyReturn() (string, error) {
	return "", nil
}

func (book Book) HasArg(arg string) string {
	return arg
}

type person struct {
	FirstName string
	LastName  string
	DOB       time.Time
}

var _ = Describe("HaveField", func() {
	var book Book
	BeforeEach(func() {
		book = Book{
			Title: "Les Miserables",
			Author: person{
				FirstName: "Victor",
				LastName:  "Hugo",
				DOB:       time.Date(1802, 2, 26, 0, 0, 0, 0, time.UTC),
			},
			Pages: 2783,
		}
	})

	DescribeTable("traversing the struct works",
		func(field string, expected interface{}) {
			Ω(book).Should(HaveField(field, expected))
		},
		Entry("Top-level field with default submatcher", "Title", "Les Miserables"),
		Entry("Top-level field with custom submatcher", "Title", ContainSubstring("Les Mis")),
		Entry("Nested field", "Author.FirstName", "Victor"),
		Entry("Top-level method", "AuthorName()", "Victor Hugo"),
		Entry("Nested method", "Author.DOB.Year()", BeNumerically("<", 1900)),
		Entry("Traversing past a method", "AbbreviatedAuthor().FirstName", Equal("Vic")),
	)

	DescribeTable("negation works",
		func(field string, expected interface{}) {
			Ω(book).ShouldNot(HaveField(field, expected))
		},
		Entry("Top-level field with default submatcher", "Title", "Les Mis"),
		Entry("Top-level field with custom submatcher", "Title", ContainSubstring("Notre Dame")),
		Entry("Nested field", "Author.FirstName", "Hugo"),
		Entry("Top-level method", "AuthorName()", "Victor M. Hugo"),
		Entry("Nested method", "Author.DOB.Year()", BeNumerically(">", 1900)),
	)

	Describe("when field lookup fails", func() {
		It("errors appropriately", func() {
			success, err := HaveField("BookName", "Les Miserables").Match(book)
			Ω(success).Should(BeFalse())
			Ω(err.Error()).Should(ContainSubstring("HaveField could not find field named '%s' in struct:", "BookName"))

			success, err = HaveField("BookName", "Les Miserables").Match(book)
			Ω(success).Should(BeFalse())
			Ω(err.Error()).Should(ContainSubstring("HaveField could not find field named '%s' in struct:", "BookName"))

			success, err = HaveField("AuthorName", "Victor Hugo").Match(book)
			Ω(success).Should(BeFalse())
			Ω(err.Error()).Should(ContainSubstring("HaveField could not find field named '%s' in struct:", "AuthorName"))

			success, err = HaveField("Title()", "Les Miserables").Match(book)
			Ω(success).Should(BeFalse())
			Ω(err.Error()).Should(ContainSubstring("HaveField could not find method named '%s' in struct of type matchers_test.Book.", "Title()"))

			success, err = HaveField("NoReturn()", "Les Miserables").Match(book)
			Ω(success).Should(BeFalse())
			Ω(err.Error()).Should(ContainSubstring("HaveField found an invalid method named 'NoReturn()' in struct of type matchers_test.Book.\nMethods must take no arguments and return exactly one value."))

			success, err = HaveField("TooManyReturn()", "Les Miserables").Match(book)
			Ω(success).Should(BeFalse())
			Ω(err.Error()).Should(ContainSubstring("HaveField found an invalid method named 'TooManyReturn()' in struct of type matchers_test.Book.\nMethods must take no arguments and return exactly one value."))

			success, err = HaveField("HasArg()", "Les Miserables").Match(book)
			Ω(success).Should(BeFalse())
			Ω(err.Error()).Should(ContainSubstring("HaveField found an invalid method named 'HasArg()' in struct of type matchers_test.Book.\nMethods must take no arguments and return exactly one value."))

			success, err = HaveField("Pages.Count", 2783).Match(book)
			Ω(success).Should(BeFalse())
			Ω(err.Error()).Should(Equal("HaveField encountered:\n    <int>: 2783\nWhich is not a struct."))

			success, err = HaveField("Author.Abbreviation", "Vic").Match(book)
			Ω(success).Should(BeFalse())
			Ω(err.Error()).Should(ContainSubstring("HaveField could not find field named '%s' in struct:", "Abbreviation"))
		})
	})

	Describe("Failure Messages", func() {
		It("renders the underlying matcher failure", func() {
			matcher := HaveField("Title", "Les Mis")
			success, err := matcher.Match(book)
			Ω(success).Should(BeFalse())
			Ω(err).ShouldNot(HaveOccurred())

			msg := matcher.FailureMessage(book)
			Ω(msg).Should(Equal("Value for field 'Title' failed to satisfy matcher.\nExpected\n    <string>: Les Miserables\nto equal\n    <string>: Les Mis"))

			matcher = HaveField("Title", "Les Miserables")
			success, err = matcher.Match(book)
			Ω(success).Should(BeTrue())
			Ω(err).ShouldNot(HaveOccurred())

			msg = matcher.NegatedFailureMessage(book)
			Ω(msg).Should(Equal("Value for field 'Title' satisfied matcher, but should not have.\nExpected\n    <string>: Les Miserables\nnot to equal\n    <string>: Les Miserables"))
		})
	})
})
