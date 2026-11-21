package gexec_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
	"time"
)

func TestGexec(t *testing.T) {
	AfterSuite(func() {
		gexec.CleanupBuildArtifacts()
	})

	// Allow for a slow child process startup or exit on loaded machines.
	SetDefaultEventuallyTimeout(5 * time.Second)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gexec Suite")
}
