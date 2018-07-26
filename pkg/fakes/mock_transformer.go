package fakes

import (
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type MockTransformer struct {
	passedHeader   core.Header
	passedHeaderID int64
	executeCalled  bool
	executeErr     error
}

func NewMockTransformer() *MockTransformer {
	return &MockTransformer{
		passedHeader:   core.Header{},
		passedHeaderID: 0,
		executeCalled:  false,
		executeErr:     nil,
	}
}

func (transformer *MockTransformer) SetExecuteErr(err error) {
	transformer.executeErr = err
}

func (transformer *MockTransformer) Execute(header core.Header, headerID int64) error {
	transformer.executeCalled = true
	transformer.passedHeader = header
	transformer.passedHeaderID = headerID
	return transformer.executeErr
}

func (transformer *MockTransformer) AssertExecuteCalledWith(header core.Header, headerID int64) {
	Expect(transformer.executeCalled).To(BeTrue())
	Expect(header).To(Equal(transformer.passedHeader))
	Expect(headerID).To(Equal(transformer.passedHeaderID))
}

func (tranformer *MockTransformer) AssertExecuteNotCalled() {
	Expect(tranformer.executeCalled).To(BeFalse())
}
