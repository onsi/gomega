/*

gexec provides support for testing external processes.

Documentation coming soon!

*/

package gexec

import (
	"io"
	"os/exec"
	"reflect"
	"sync"
	"syscall"

	"github.com/onsi/gomega/gbytes"
)

type Session struct {
	command  *exec.Cmd
	Out      *gbytes.Buffer
	Err      *gbytes.Buffer
	lock     *sync.Mutex
	exitCode int
}

func Start(command *exec.Cmd, outWriter io.Writer, errWriter io.Writer) (*Session, error) {
	session := &Session{
		command:  command,
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

func (s *Session) monitorForExit() {
	s.command.Wait()
	s.lock.Lock()
	s.exitCode = s.command.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
	s.lock.Unlock()
}

func (s *Session) getExitCode() int {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.exitCode
}
