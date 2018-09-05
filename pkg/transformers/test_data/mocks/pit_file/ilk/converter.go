package ilk

import (
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file/ilk"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
)

type MockPitFileIlkConverter struct {
	PassedContractAddress string
	PassedContractABI     string
	PassedLog             types.Log
	converterError        error
}

func (converter *MockPitFileIlkConverter) SetConverterError(err error) {
	converter.converterError = err
}

func (converter *MockPitFileIlkConverter) ToModel(contractAddress string, contractAbi string, ethLog types.Log) (ilk.PitFileIlkModel, error) {
	converter.PassedContractAddress = contractAddress
	converter.PassedContractABI = contractAbi
	converter.PassedLog = ethLog
	return test_data.PitFileIlkModel, converter.converterError
}
