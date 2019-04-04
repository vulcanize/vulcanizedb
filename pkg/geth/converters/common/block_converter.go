// VulcanizeDB
// Copyright © 2019 Vulcanize

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

package common

import (
	"encoding/json"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type BlockConverter struct {
	transactionConverter TransactionConverter
}

func NewBlockConverter(transactionConverter TransactionConverter) BlockConverter {
	return BlockConverter{transactionConverter: transactionConverter}
}

func (bc BlockConverter) ToCoreBlock(gethBlock *types.Block) (core.Block, error) {
	transactions, err := bc.transactionConverter.ConvertBlockTransactionsToCore(gethBlock)
	if err != nil {
		return core.Block{}, err
	}

	coreBlock := core.Block{
		Difficulty:   gethBlock.Difficulty().Int64(),
		ExtraData:    hexutil.Encode(gethBlock.Extra()),
		GasLimit:     gethBlock.GasLimit(),
		GasUsed:      gethBlock.GasUsed(),
		Hash:         gethBlock.Hash().Hex(),
		Miner:        strings.ToLower(gethBlock.Coinbase().Hex()),
		Nonce:        hexutil.Encode(gethBlock.Header().Nonce[:]),
		Number:       gethBlock.Number().Int64(),
		ParentHash:   gethBlock.ParentHash().Hex(),
		Size:         gethBlock.Size().String(),
		Time:         gethBlock.Time().Int64(),
		Transactions: transactions,
		UncleHash:    gethBlock.UncleHash().Hex(),
	}
	coreBlock.Reward = CalcBlockReward(coreBlock, gethBlock.Uncles()).String()
	totalUncleReward, uncles := bc.ToCoreUncle(coreBlock, gethBlock.Uncles())

	coreBlock.UnclesReward = totalUncleReward.String()
	coreBlock.Uncles = uncles
	return coreBlock, nil
}

// Rewards for the miners of uncles is calculated as (U_n + 8 - B_n) * R / 8
// Where U_n is the uncle block number, B_n is the parent block number and R is the static block reward at B_n
// https://github.com/ethereum/go-ethereum/issues/1591
// https://ethereum.stackexchange.com/questions/27172/different-uncles-reward
// https://github.com/ethereum/homestead-guide/issues/399
// Returns the total uncle reward and the individual processed uncles
func (bc BlockConverter) ToCoreUncle(block core.Block, uncles []*types.Header) (*big.Int, []core.Uncle) {
	totalUncleRewards := new(big.Int)
	coreUncles := make([]core.Uncle, 0, len(uncles))
	for _, uncle := range uncles {
		thisUncleReward := calcUncleMinerReward(block.Number, uncle.Number.Int64())
		raw, _ := json.Marshal(uncle)
		coreUncle := core.Uncle{
			Miner:     uncle.Coinbase.Hex(),
			Hash:      uncle.Hash().Hex(),
			Raw:       raw,
			Reward:    thisUncleReward.String(),
			Timestamp: uncle.Time.String(),
		}
		coreUncles = append(coreUncles, coreUncle)
		totalUncleRewards.Add(totalUncleRewards, thisUncleReward)
	}
	return totalUncleRewards, coreUncles
}
