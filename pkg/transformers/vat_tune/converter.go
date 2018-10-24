package vat_tune

import (
	"bytes"
	"encoding/json"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
)

type VatTuneConverter struct{}

func (VatTuneConverter) ToModels(ethLogs []types.Log) ([]interface{}, error) {
	var models []interface{}
	for _, ethLog := range ethLogs {
		err := verifyLog(ethLog)
		if err != nil {
			return nil, err
		}
		ilk := string(bytes.Trim(ethLog.Topics[1].Bytes(), "\x00"))
		urn := common.BytesToAddress(ethLog.Topics[2].Bytes())
		v := common.BytesToAddress(ethLog.Topics[3].Bytes())
		wBytes := shared.GetDataBytesAtIndex(-3, ethLog.Data)
		w := common.BytesToAddress(wBytes)
		dinkBytes := shared.GetDataBytesAtIndex(-2, ethLog.Data)
		dink := big.NewInt(0).SetBytes(dinkBytes)
		dartBytes := shared.GetDataBytesAtIndex(-1, ethLog.Data)
		dart := big.NewInt(0).SetBytes(dartBytes)

		raw, err := json.Marshal(ethLog)
		if err != nil {
			return nil, err
		}
		model := VatTuneModel{
			Ilk:              ilk,
			Urn:              urn.String(),
			V:                v.String(),
			W:                w.String(),
			Dink:             dink.String(),
			Dart:             dart.String(),
			TransactionIndex: ethLog.TxIndex,
			LogIndex:         ethLog.Index,
			Raw:              raw,
		}
		models = append(models, model)
	}
	return models, nil
}

func verifyLog(log types.Log) error {
	if len(log.Topics) < 4 {
		return errors.New("log missing topics")
	}
	if len(log.Data) < shared.DataItemLength {
		return errors.New("log missing data")
	}
	return nil
}
