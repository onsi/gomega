package gleak

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGleak(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gleak Suite")
}
