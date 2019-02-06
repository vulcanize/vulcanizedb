package fakes

import (
	"github.com/hpcloud/tail"
	"gopkg.in/tomb.v1"
)

type MockTailer struct {
	Lines      chan *tail.Line
	TailCalled bool
}

func NewMockTailer() *MockTailer {
	return &MockTailer{
		Lines:      make(chan *tail.Line, 1),
		TailCalled: false,
	}
}

func (mock *MockTailer) Tail() (*tail.Tail, error) {
	mock.TailCalled = true
	fakeTail := &tail.Tail{
		Filename: "",
		Lines:    mock.Lines,
		Config:   tail.Config{},
		Tomb:     tomb.Tomb{},
	}
	return fakeTail, nil
}
