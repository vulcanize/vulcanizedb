package fakes

type MockSubscription struct {
	Errs chan error
}

func (m *MockSubscription) Err() <-chan error {
	return m.Errs
}

func (m *MockSubscription) Unsubscribe() {
	panic("implement me")
}
