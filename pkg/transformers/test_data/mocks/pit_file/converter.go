package pit_file

import (
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
)

type MockPitFileConverter struct {
	PassedContractAddress string
	PassedContractABI     string
	PassedLog             types.Log
	converterError        error
}

func (converter *MockPitFileConverter) SetConverterError(err error) {
	converter.converterError = err
}

func (converter *MockPitFileConverter) ToModel(contractAddress string, contractAbi string, ethLog types.Log) (pit_file.PitFileModel, error) {
	converter.PassedContractAddress = contractAddress
	converter.PassedContractABI = contractAbi
	converter.PassedLog = ethLog
	return test_data.PitFileModel, converter.converterError
}
