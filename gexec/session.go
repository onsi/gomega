/*
Package gexec provides support for testing external processes.
*/
package gexec

import (
	"io"
	"os/exec"
	"reflect"
	"sync"
	"syscall"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

type Session struct {
	//The wrapped command
	Command *exec.Cmd

	//A *gbytes.Buffer connected to the command's stdout
	Out *gbytes.Buffer

	//A *gbytes.Buffer connected to the command's stderr
	Err *gbytes.Buffer

	lock     *sync.Mutex
	exitCode int
}

/*
Start starts the passed-in *exec.Cmd command.  It wraps the command in a *gexec.Session.

The session pipes the command's stdout and stderr to two *gbytes.Buffers available as properties on the session: session.Out and session.Err.
These buffers can be used with the gbytes.Say matcher to match against unread output:

	Ω(session.Out).Should(gbytes.Say("foo-out"))
	Ω(session.Err).Should(gbytes.Say("foo-err"))

In addition, Session satisfies the gbytes.BufferProvider interface and provides the stdout *gbytes.Buffer.  This allows you to replace the first line, above, with:

	Ω(session).Should(gbytes.Say("foo-out"))

When outWriter and/or errWriter are non-nil, the session will pipe stdout and/or stderr output both into the session *gybtes.Buffers and to the passed-in outWriter/errWriter.
This is useful for capturing the process's output or logging it to screen.  In particular, when using Ginkgo it can be convenient to direct output to the GinkgoWriter:

	session, err := Start(command, GinkgoWriter, GinkgoWriter)

This will log output when running tests in verbose mode, but - otherwise - will only log output when a test fails.

The session wrapper is responsible for waiting on the *exec.Cmd command.  You *should not* call command.Wait() yourself.
Instead, to assert that the command has exited you can use the gexec.Exit matcher:

	Ω(session).Should(gexec.Exit())

When the session exits it closes the stdout and stderr gbytes buffers.  This will short circuit any
Eventuallys waiting fo the buffers to Say something.
*/
func Start(command *exec.Cmd, outWriter io.Writer, errWriter io.Writer) (*Session, error) {
	session := &Session{
		Command:  command,
		Out:      gbytes.NewBuffer(),
		Err:      gbytes.NewBuffer(),
		lock:     &sync.Mutex{},
		exitCode: -1,
	}

	var commandOut, commandErr io.Writer

	commandOut, commandErr = session.Out, session.Err

	if outWriter != nil && !reflect.ValueOf(outWriter).IsNil() {
		commandOut = io.MultiWriter(commandOut, outWriter)
	}

	if errWriter != nil && !reflect.ValueOf(errWriter).IsNil() {
		commandErr = io.MultiWriter(commandErr, errWriter)
	}

	command.Stdout = commandOut
	command.Stderr = commandErr

	err := command.Start()
	if err == nil {
		go session.monitorForExit()
	}

	return session, err
}

/*
Buffer implements the gbytes.BufferProvider interface and returns s.Out
This allows you to make gbytes.Say matcher assertions against stdout without having to reference .Out:

	Eventually(session).Should(gbytes.Say("foo"))
*/
func (s *Session) Buffer() *gbytes.Buffer {
	return s.Out
}

/*
ExitCode returns the wrapped command's exit code.  If the command hasn't exited yet, ExitCode returns -1.

To assert that the command has exited it is more convenient to use the Exit matcher:

	Eventually(s).Should(gexec.Exit())
*/
func (s *Session) ExitCode() int {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.exitCode
}

/*
Wait waits until the wrapped command exits.  It can be passed an optional timeout.
If the command does not exit within the timeout, Wait will trigger a test failure.

Wait returns the session, making it possible to chain:

	session.Wait().Out.Contents()

will wait for the command to exit then return the entirety of Out's contents.

Wait uses eventually under the hood and accepts the same timeout/polling intervals that eventually does.
*/
func (s *Session) Wait(timeout ...interface{}) *Session {
	Eventually(s, timeout...).Should(Exit())
	return s
}

func (s *Session) monitorForExit() {
	s.Command.Wait()
	s.lock.Lock()
	s.exitCode = s.Command.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
	s.Out.Close()
	s.Err.Close()
	s.lock.Unlock()
}
