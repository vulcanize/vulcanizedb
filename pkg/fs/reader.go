package fs

import "io/ioutil"

type Reader interface {
	Read(path string) ([]byte, error)
}

type FsReader struct {
}

func (FsReader) Read(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}
