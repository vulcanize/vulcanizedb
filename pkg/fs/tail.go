package fs

import "github.com/hpcloud/tail"

type Tailer interface {
	Tail() (*tail.Tail, error)
}

type FileTailer struct {
	Path string
}

func (tailer FileTailer) Tail() (*tail.Tail, error) {
	return tail.TailFile(tailer.Path, tail.Config{Follow: true})
}
