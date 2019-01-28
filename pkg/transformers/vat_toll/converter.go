package vat_toll

import (
	"encoding/json"
	"errors"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type VatTollConverter struct{}

func (VatTollConverter) ToModels(ethLogs []types.Log) ([]interface{}, error) {
	var models []interface{}
	for _, ethLog := range ethLogs {
		err := verifyLog(ethLog)
		if err != nil {
			return nil, err
		}
		ilk := shared.GetHexWithoutPrefix(ethLog.Topics[1].Bytes())
		urn := common.BytesToAddress(ethLog.Topics[2].Bytes()[:common.AddressLength])
		take := ethLog.Topics[3].Big()

		raw, err := json.Marshal(ethLog)
		if err != nil {
			return nil, err
		}
		model := VatTollModel{
			Ilk:              ilk,
			Urn:              urn.String(),
			Take:             take.String(),
			TransactionIndex: ethLog.TxIndex,
			LogIndex:         ethLog.Index,
			Raw:              raw,
		}
		models = append(models, model)
	}
	return models, nil
}

func verifyLog(log types.Log) error {
	numTopicInValidLog := 4
	if len(log.Topics) < numTopicInValidLog {
		return errors.New("log missing topics")
	}
	return nil
}
