package stability_fee

import (
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file/stability_fee"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
)

type MockPitFileStabilityFeeConverter struct {
	converterErr          error
	PassedContractAddress string
	PassedContractABI     string
	PassedLog             types.Log
}

func (converter *MockPitFileStabilityFeeConverter) ToModel(contractAddress string, contractAbi string, ethLog types.Log) (stability_fee.PitFileStabilityFeeModel, error) {
	converter.PassedContractAddress = contractAddress
	converter.PassedContractABI = contractAbi
	converter.PassedLog = ethLog
	return test_data.PitFileStabilityFeeModel, converter.converterErr
}
func (converter *MockPitFileStabilityFeeConverter) SetConverterError(e error) {
	converter.converterErr = e
}
