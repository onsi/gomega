package gbytes_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGbytes(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gbytes Suite")
}

func interceptFailures(f func()) []string {
	failures := []string{}
	RegisterFailHandler(func(message string, callerSkip ...int) {
		failures = append(failures, message)
	})
	f()
	RegisterFailHandler(Fail)
	return failures
}
