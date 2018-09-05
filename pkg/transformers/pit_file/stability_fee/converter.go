package stability_fee

import (
	"bytes"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Converter interface {
	ToModel(contractAddress string, contractAbi string, ethLog types.Log) (PitFileStabilityFeeModel, error)
}

type PitFileStabilityFeeConverter struct{}

func (PitFileStabilityFeeConverter) ToModel(contractAddress string, contractAbi string, ethLog types.Log) (PitFileStabilityFeeModel, error) {
	what := string(bytes.Trim(ethLog.Topics[2].Bytes(), "\x00"))
	data := common.HexToAddress(ethLog.Topics[1].String()).Hex()

	raw, err := json.Marshal(ethLog)
	return PitFileStabilityFeeModel{
		What:             what,
		Data:             data,
		TransactionIndex: ethLog.TxIndex,
		Raw:              raw,
	}, err
}
