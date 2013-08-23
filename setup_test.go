package gomega

import (
	. "github.com/onsi/ginkgo"
	"testing"
)

func Test(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gomega")
}
