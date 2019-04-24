package fakes

import (
	"github.com/hpcloud/tail"
	"gopkg.in/tomb.v1"
)

type MockTailer struct {
	Lines   chan *tail.Line
	TailErr error
}

func NewMockTailer() *MockTailer {
	return &MockTailer{
		Lines: make(chan *tail.Line, 1),
	}
}

func (mock *MockTailer) Tail() (*tail.Tail, error) {
	fakeTail := &tail.Tail{
		Filename: "",
		Lines:    mock.Lines,
		Config:   tail.Config{},
		Tomb:     tomb.Tomb{},
	}
	return fakeTail, mock.TailErr
}
