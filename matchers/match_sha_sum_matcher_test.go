package matchers_test

import (
	"io/ioutil"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/matchers"
)

var _ = Describe("MatchSHASumOfMatcher", func() {
	Describe("Match", func() {
		Context("when the files have the same contents", func() {
			It("returns true", func() {
				tempDir, err := ioutil.TempDir("", "")
				Expect(err).NotTo(HaveOccurred())

				file1 := filepath.Join(tempDir, "file-1")
				file2 := filepath.Join(tempDir, "file-2")

				err = ioutil.WriteFile(file1, []byte("file contents"), 0644)
				Expect(err).NotTo(HaveOccurred())

				err = ioutil.WriteFile(file2, []byte("file contents"), 0644)
				Expect(err).NotTo(HaveOccurred())

				matcher := MatchSHASumOf(file1)
				success, err := matcher.Match(file2)
				Expect(err).NotTo(HaveOccurred())
				Expect(success).To(BeTrue())
			})
		})

		Context("when the files do not have the same contents", func() {
			It("returns true", func() {
				tempDir, err := ioutil.TempDir("", "")
				Expect(err).NotTo(HaveOccurred())

				file1 := filepath.Join(tempDir, "file-1")
				file2 := filepath.Join(tempDir, "file-2")

				err = ioutil.WriteFile(file1, []byte("some contents"), 0644)
				Expect(err).NotTo(HaveOccurred())

				err = ioutil.WriteFile(file2, []byte("other contents"), 0644)
				Expect(err).NotTo(HaveOccurred())

				matcher := MatchSHASumOf(file1)
				success, err := matcher.Match(file2)
				Expect(err).NotTo(HaveOccurred())
				Expect(success).To(BeFalse())
			})
		})
	})
})
