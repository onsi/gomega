package gexec_test

import (
	"bytes"
	"os/exec"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Session", func() {
	var command *exec.Cmd
	var session *Session

	var outWriter, errWriter *bytes.Buffer

	BeforeEach(func() {
		outWriter = nil
		errWriter = nil
	})

	JustBeforeEach(func() {
		command = exec.Command(fireflyPath)
		var err error
		session, err = Start(command, outWriter, errWriter)
		Ω(err).ShouldNot(HaveOccurred())
	})

	Context("running a command", func() {
		It("should start the process", func() {
			Ω(command.Process).ShouldNot(BeNil())
		})

		It("should wrap the process's stdout and stderr with gbytes buffers", func(done Done) {
			Eventually(session.Out).Should(Say("We've done the impossible, and that makes us mighty"))
			Eventually(session.Err).Should(Say("Ah, curse your sudden but inevitable betrayal!"))
			defer session.Out.CancelDetects()

			select {
			case <-session.Out.Detect("Can we maybe vote on the whole murdering people issue"):
				Eventually(session).Should(Exit(0))
			case <-session.Out.Detect("I swear by my pretty floral bonnet, I will end you."):
				Eventually(session).Should(Exit(1))
			case <-session.Out.Detect("My work's illegal, but at least it's honest."):
				Eventually(session).Should(Exit(2))
			}

			close(done)
		})
	})

	Context("when wrapping out and err", func() {
		BeforeEach(func() {
			outWriter = &bytes.Buffer{}
			errWriter = &bytes.Buffer{}
		})

		It("should route to both the provided writers and the gbytes buffers", func() {
			Eventually(session.Out).Should(Say("We've done the impossible, and that makes us mighty"))
			Eventually(session.Err).Should(Say("Ah, curse your sudden but inevitable betrayal!"))

			Ω(outWriter.String()).Should(ContainSubstring("We've done the impossible, and that makes us mighty"))
			Ω(errWriter.String()).Should(ContainSubstring("Ah, curse your sudden but inevitable betrayal!"))

			Eventually(session).Should(Exit())

			Ω(outWriter.Bytes()).Should(Equal(session.Out.Contents()))
			Ω(errWriter.Bytes()).Should(Equal(session.Err.Contents()))
		})
	})

	Describe("when the command fails to start", func() {
		It("should return an error", func() {
			_, err := Start(exec.Command("agklsjdfas"), nil, nil)
			Ω(err).Should(HaveOccurred())
		})
	})
})
