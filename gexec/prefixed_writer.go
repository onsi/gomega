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
	prefix    []byte
	writer    io.Writer
	lock      *sync.Mutex
	isNewLine bool
}

func NewPrefixedWriter(prefix string, writer io.Writer) *PrefixedWriter {
	return &PrefixedWriter{
		prefix:    []byte(prefix),
		writer:    writer,
		lock:      &sync.Mutex{},
		isNewLine: true,
	}
}

func (w *PrefixedWriter) Write(b []byte) (int, error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	newLine := []byte("\n")
	segments := bytes.Split(b, newLine)

	for i, segment := range segments {
		if w.isNewLine {
			_, err := w.writer.Write(w.prefix)
			if err != nil {
				return 0, err
			}
			w.isNewLine = false
		}
		_, err := w.writer.Write(segment)
		if err != nil {
			return 0, err
		}
		if i < len(segments)-1 {
			_, err := w.writer.Write(newLine)
			if err != nil {
				return 0, err
			}
			w.isNewLine = true
		}
	}

	return len(b), nil
}
