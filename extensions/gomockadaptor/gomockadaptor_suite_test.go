package gomockadaptor_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGomockadaptor(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gomockadaptor Suite")
}
