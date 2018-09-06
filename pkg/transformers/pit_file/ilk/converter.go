package ilk

import (
	"bytes"
	"encoding/json"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

type Converter interface {
	ToModel(contractAddress string, contractAbi string, ethLog types.Log) (PitFileIlkModel, error)
}

type PitFileIlkConverter struct{}

func (PitFileIlkConverter) ToModel(contractAddress string, contractAbi string, ethLog types.Log) (entity PitFileIlkModel, err error) {
	ilk := string(bytes.Trim(ethLog.Topics[2].Bytes(), "\x00"))
	what := string(bytes.Trim(ethLog.Topics[3].Bytes(), "\x00"))
	itemByteLength := 32
	riskBytes := ethLog.Data[len(ethLog.Data)-itemByteLength:]
	risk := big.NewInt(0).SetBytes(riskBytes).String()

	raw, err := json.Marshal(ethLog)
	return PitFileIlkModel{
		Ilk:              ilk,
		What:             what,
		Data:             risk,
		TransactionIndex: ethLog.TxIndex,
		Raw:              raw,
	}, err
}
