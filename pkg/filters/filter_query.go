package filters

import (
	"encoding/json"

	"errors"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type LogFilters []LogFilter

type LogFilter struct {
	Name        string `json:"name"`
	FromBlock   int64  `json:"fromBlock"`
	ToBlock     int64  `json:"toBlock"`
	Address     string `json:"address"`
	core.Topics `json:"topics"`
}

func (filterQuery *LogFilter) UnmarshalJSON(input []byte) error {
	type Alias LogFilter

	var err error
	aux := &struct {
		ToBlock   string `json:"toBlock"`
		FromBlock string `json:"fromBlock"`
		*Alias
	}{
		Alias: (*Alias)(filterQuery),
	}
	if err := json.Unmarshal(input, &aux); err != nil {
		return err
	}
	if filterQuery.Name == "" {
		return errors.New("filters: must provide name for logfilter")
	}
	filterQuery.ToBlock, err = filterQuery.unmarshalFromToBlock(aux.ToBlock)
	if err != nil {
		return errors.New("filters: invalid fromBlock")
	}
	filterQuery.FromBlock, err = filterQuery.unmarshalFromToBlock(aux.FromBlock)
	if err != nil {
		return errors.New("filters: invalid fromBlock")
	}
	if !common.IsHexAddress(filterQuery.Address) {
		return errors.New("filters: invalid address")
	}

	return nil
}

func (filterQuery *LogFilter) unmarshalFromToBlock(auxBlock string) (int64, error) {
	if auxBlock == "" {
		return -1, nil
	}
	block, err := hexutil.DecodeUint64(auxBlock)
	if err != nil {
		return 0, errors.New("filters: invalid block arg")
	}
	return int64(block), nil
}
