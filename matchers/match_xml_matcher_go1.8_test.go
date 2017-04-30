// +build go1.8

package matchers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MatchXMLMatcher Go 1.8", func() {

	var (
		sample_09 = readFileContents("test_data/xml/sample_09.xml")
		sample_10 = readFileContents("test_data/xml/sample_10.xml")
	)

	Context("When passed stringifiables", func() {
		It("should succeed if the XML matches", func() {
			Ω(sample_09).ShouldNot(MatchXML(sample_10)) // same structures with different attribute values
		})
	})
})
