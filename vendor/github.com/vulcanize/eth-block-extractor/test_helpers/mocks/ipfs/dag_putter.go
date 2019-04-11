package ipfs

type MockDagPutter struct {
	Called          bool
	PassedInterface interface{}
	Err             error
}

func NewMockDagPutter() *MockDagPutter {
	return &MockDagPutter{
		Called:          false,
		PassedInterface: nil,
		Err:             nil,
	}
}

func (mdp *MockDagPutter) SetError(err error) {
	mdp.Err = err
}

func (mdp *MockDagPutter) DagPut(raw interface{}) ([]string, error) {
	mdp.Called = true
	mdp.PassedInterface = raw
	return nil, mdp.Err
}
