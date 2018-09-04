package vat_init

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/ethereum/go-ethereum/core/types"
)

type Converter interface {
	ToModel(ethLog types.Log) (VatInitModel, error)
}

type VatInitConverter struct{}

func (VatInitConverter) ToModel(ethLog types.Log) (VatInitModel, error) {
	err := verifyLog(ethLog)
	if err != nil {
		return VatInitModel{}, err
	}
	ilk := string(bytes.Trim(ethLog.Topics[1].Bytes(), "\x00"))
	raw, err := json.Marshal(ethLog)
	return VatInitModel{
		Ilk:              ilk,
		TransactionIndex: ethLog.TxIndex,
		Raw:              raw,
	}, err
}

func verifyLog(log types.Log) error {
	if len(log.Topics) < 2 {
		return errors.New("log missing topics")
	}
	return nil
}
