package fs

import (
	"fmt"
	"os"
)

type StatusWriter interface {
	Write() error
}

type statusWriter struct {
	file      string
	message   []byte
	perm      os.FileMode
	fileFlags int
}

func NewStatusWriter(file string, message []byte) statusWriter {
	return statusWriter{
		file:      file,
		message:   message,
		perm:      0644,
		fileFlags: os.O_APPEND | os.O_CREATE | os.O_WRONLY,
	}
}
func (w statusWriter) Write() error {
	file, openErr := os.OpenFile(w.file, w.fileFlags, w.perm)
	if openErr != nil {
		return fmt.Errorf("error opening file %s: %w", w.file, openErr)
	}

	_, writeErr := file.Write(w.message)
	if writeErr != nil {
		closeErr := file.Close()
		if closeErr != nil {
			return fmt.Errorf("error closing file %s: %w", w.file, closeErr)
		}
		return fmt.Errorf("error writing file %s: %w", w.file, writeErr)
	}
	return file.Close()
}
