package vat_grab

import (
	"encoding/json"
	"errors"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"math/big"
)

type VatGrabConverter struct{}

func (VatGrabConverter) ToModels(ethLogs []types.Log) ([]interface{}, error) {
	var models []interface{}
	for _, ethLog := range ethLogs {
		err := verifyLog(ethLog)
		if err != nil {
			return nil, err
		}
		ilk := shared.GetHexWithoutPrefix(ethLog.Topics[1].Bytes())
		urn := shared.GetHexWithoutPrefix(ethLog.Topics[2].Bytes())
		v := shared.GetHexWithoutPrefix(ethLog.Topics[3].Bytes())
		wBytes := shared.GetDataBytesAtIndex(-3, ethLog.Data)
		w := shared.GetHexWithoutPrefix(wBytes)
		dinkBytes := shared.GetDataBytesAtIndex(-2, ethLog.Data)
		dink := big.NewInt(0).SetBytes(dinkBytes).String()
		dartBytes := shared.GetDataBytesAtIndex(-1, ethLog.Data)
		dart := big.NewInt(0).SetBytes(dartBytes).String()

		raw, err := json.Marshal(ethLog)
		if err != nil {
			return nil, err
		}
		model := VatGrabModel{
			Ilk:              ilk,
			Urn:              urn,
			V:                v,
			W:                w,
			Dink:             dink,
			Dart:             dart,
			LogIndex:         ethLog.Index,
			TransactionIndex: ethLog.TxIndex,
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
	if len(log.Data) < constants.DataItemLength {
		return errors.New("log missing data")
	}
	return nil
}
