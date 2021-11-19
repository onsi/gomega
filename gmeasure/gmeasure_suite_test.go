package gmeasure_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGmeasure(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gmeasure Suite")
}
