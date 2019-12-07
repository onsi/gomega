package bipartitegraph_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestBipartitegraph(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bipartitegraph Suite")
}
