// VulcanizeDB
// Copyright © 2018 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package filters

import (
	"encoding/json"

	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type LogFilters []LogFilter

type LogFilter struct {
	Name        string `json:"name"`
	FromBlock   int64  `json:"fromBlock" db:"from_block"`
	ToBlock     int64  `json:"toBlock" db:"to_block"`
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
	if err = json.Unmarshal(input, &aux); err != nil {
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
