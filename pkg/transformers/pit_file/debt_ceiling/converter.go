package debt_ceiling

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

type Converter interface {
	ToModel(contractAddress string, contractAbi string, ethLog types.Log) (PitFileDebtCeilingModel, error)
}

type PitFileDebtCeilingConverter struct{}

func (PitFileDebtCeilingConverter) ToModel(contractAddress string, contractAbi string, ethLog types.Log) (PitFileDebtCeilingModel, error) {
	what := common.HexToAddress(ethLog.Topics[1].String()).String()
	itemByteLength := 32
	riskBytes := ethLog.Data[len(ethLog.Data)-itemByteLength:]
	data := big.NewInt(0).SetBytes(riskBytes).String()

	raw, err := json.Marshal(ethLog)
	return PitFileDebtCeilingModel{
		What:             what,
		Data:             data,
		TransactionIndex: ethLog.TxIndex,
		Raw:              raw,
	}, err
}
