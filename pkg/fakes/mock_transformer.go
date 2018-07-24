package fakes

import (
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type MockTransformer struct {
	passedHeader   core.Header
	passedHeaderID int64
	executeErr     error
}

func NewMockTransformer() *MockTransformer {
	return &MockTransformer{
		passedHeader:   core.Header{},
		passedHeaderID: 0,
		executeErr:     nil,
	}
}

func (transformer *MockTransformer) SetExecuteErr(err error) {
	transformer.executeErr = err
}

func (transformer *MockTransformer) Execute(header core.Header, headerID int64) error {
	transformer.passedHeader = header
	transformer.passedHeaderID = headerID
	return transformer.executeErr
}

func (transformer *MockTransformer) AssertExecuteCalledWith(header core.Header, headerID int64) {
	Expect(header).To(Equal(transformer.passedHeader))
	Expect(headerID).To(Equal(transformer.passedHeaderID))
}
