package fakes

type MockStatusWriter struct {
	WriteCalled bool
}

func (w *MockStatusWriter) Write() error {
	w.WriteCalled = true
	return nil
}
