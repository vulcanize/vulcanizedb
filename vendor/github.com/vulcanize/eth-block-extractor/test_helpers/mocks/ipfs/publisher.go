package ipfs

import (
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/eth-block-extractor/test_helpers"
)

type MockPublisher struct {
	err              error
	passedBlockDatas []interface{}
	returnStrings    [][]string
}

func NewMockPublisher() *MockPublisher {
	return &MockPublisher{
		err:              nil,
		passedBlockDatas: []interface{}{},
		returnStrings:    nil,
	}
}

func (publisher *MockPublisher) SetReturnStrings(returnBytes [][]string) {
	publisher.returnStrings = returnBytes
}

func (publisher *MockPublisher) SetError(err error) {
	publisher.err = err
}

func (publisher *MockPublisher) Write(input interface{}) ([]string, error) {
	publisher.passedBlockDatas = append(publisher.passedBlockDatas, input)
	if publisher.err != nil {
		return nil, publisher.err
	}
	var stringsToReturn []string
	if len(publisher.returnStrings) > 0 {
		stringsToReturn = publisher.returnStrings[0]
		if len(publisher.returnStrings) > 1 {
			publisher.returnStrings = publisher.returnStrings[1:]
		} else {
			publisher.returnStrings = [][]string{{test_helpers.FakeString}}
		}
	} else {
		stringsToReturn = []string{test_helpers.FakeString}
	}
	return stringsToReturn, nil
}

func (publisher *MockPublisher) AssertWriteCalledWithBytes(inputs [][]byte) {
	for i := 0; i < len(inputs); i++ {
		Expect(publisher.passedBlockDatas).To(ContainElement(inputs[i]))
	}
	for i := 0; i < len(publisher.passedBlockDatas); i++ {
		Expect(inputs).To(ContainElement(publisher.passedBlockDatas[i]))
	}
}

func (publisher *MockPublisher) AssertWriteCalledWithInterfaces(interfaces []interface{}) {
	for i := 0; i < len(interfaces); i++ {
		Expect(publisher.passedBlockDatas).To(ContainElement(interfaces[i]))
	}
	for i := 0; i < len(publisher.passedBlockDatas); i++ {
		Expect(interfaces).To(ContainElement(publisher.passedBlockDatas[i]))
	}
}

func (publisher *MockPublisher) AssertWriteCalledWithBodies(bodies []*types.Body) {
	var expected []*types.Body
	for i := 0; i < len(publisher.passedBlockDatas); i++ {
		expected = append(expected, publisher.passedBlockDatas[i].(*types.Body))
	}
	Expect(expected).To(Equal(bodies))
}
