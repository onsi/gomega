package matcher_tests

import (
	. "github.com/onsi/godescribe"
	. "github.com/onsi/gomega"
	"testing"
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

func Test(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gomega")
}
