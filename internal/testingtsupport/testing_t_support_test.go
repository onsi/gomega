package testingtsupport_test

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestTestingT(t *testing.T) {
	RegisterTestingT(t)
	Ω(true).Should(BeTrue())
}

func TestNewGomegaWithT(t *testing.T) {
	g := NewWithT(t)
	g.Expect(true).To(BeTrue())
}
