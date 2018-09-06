package debt_ceiling

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file/debt_ceiling"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
)

type MockPitFileDebtCeilingConverter struct {
	converterErr          error
	PassedContractAddress string
	PassedContractABI     string
	PassedLog             types.Log
}

func (converter *MockPitFileDebtCeilingConverter) ToModel(contractAddress string, contractAbi string, ethLog types.Log) (debt_ceiling.PitFileDebtCeilingModel, error) {
	converter.PassedContractAddress = contractAddress
	converter.PassedContractABI = contractAbi
	converter.PassedLog = ethLog
	return test_data.PitFileDebtCeilingModel, converter.converterErr
}

func (converter *MockPitFileDebtCeilingConverter) SetConverterError(e error) {
	converter.converterErr = e
}
