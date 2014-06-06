package gexec

import (
	"bytes"
	"io"
	"sync"
)

/*
PrefixedWriter wraps an io.Writer, emiting the passed in prefix at the beginning of each new line.
This can be useful when running multiple gexec.Sessions concurrently - you can prefix the log output of each
session by passing in a PrefixedWriter:

gexec.Start(cmd, NewPrefixedWriter("[my-cmd] ", GinkgoWriter), NewPrefixedWriter("[my-cmd] ", GinkgoWriter))
*/
type PrefixedWriter struct {
	prefix       []byte
	writer       io.Writer
	lock         *sync.Mutex
	isNewLine    bool
	isFirstWrite bool
}

func NewPrefixedWriter(prefix string, writer io.Writer) *PrefixedWriter {
	return &PrefixedWriter{
		prefix:       []byte(prefix),
		writer:       writer,
		lock:         &sync.Mutex{},
		isFirstWrite: true,
	}
}

func (w *PrefixedWriter) Write(b []byte) (int, error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	newLine := []byte("\n")
	segments := bytes.Split(b, newLine)

	if len(segments) != 0 {
		toWrite := []byte{}
		if w.isFirstWrite {
			toWrite = append(toWrite, w.prefix...)
			toWrite = append(toWrite, segments[0]...)
			w.isFirstWrite = false
		} else if w.isNewLine {
			toWrite = append(toWrite, newLine...)
			toWrite = append(toWrite, w.prefix...)
			toWrite = append(toWrite, segments[0]...)
		} else {
			toWrite = append(toWrite, segments[0]...)
		}

		for i := 1; i < len(segments)-1; i++ {
			toWrite = append(toWrite, newLine...)
			toWrite = append(toWrite, w.prefix...)
			toWrite = append(toWrite, segments[i]...)
		}

		if len(segments) > 1 {
			lastSegment := segments[len(segments)-1]

			if len(lastSegment) == 0 {
				w.isNewLine = true
			} else {
				toWrite = append(toWrite, newLine...)
				toWrite = append(toWrite, w.prefix...)
				toWrite = append(toWrite, lastSegment...)
				w.isNewLine = false
			}
		}

		_, err := w.writer.Write(toWrite)
		if err != nil {
			return 0, err
		}
	}

	return len(b), nil
}
