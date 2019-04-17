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

package ipfs

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/statediff"

	"github.com/vulcanize/vulcanizedb/pkg/core"
)

// Converter interface is used to convert a geth statediff.Payload to our IPLDPayload type
type Converter interface {
	Convert(payload statediff.Payload) (*IPLDPayload, error)
}

// PayloadConverter is the underlying struct for the Converter interface
type PayloadConverter struct {
	client core.EthClient
}

// IPLDPayload is a custom type which packages ETH data for the IPFS publisher
type IPLDPayload struct {
	HeaderRLP  []byte
	BlockNumber *big.Int
	BlockHash   common.Hash
	BlockBody  *types.Body
	Receipts   types.Receipts
	StateLeafs  map[common.Hash][]byte
	StorageLeafs map[common.Hash]map[common.Hash][]byte
}

// NewPayloadConverter creates a pointer to a new PayloadConverter which satisfies the Converter interface
func NewPayloadConverter(client core.EthClient) *PayloadConverter {
	return &PayloadConverter{
		client: client,
	}
}

// Convert method is used to convert a geth statediff.Payload to a IPLDPayload
func (pc *PayloadConverter) Convert(payload statediff.Payload) (*IPLDPayload, error) {
	// Unpack block rlp to access fields
	block := new(types.Block)
	err := rlp.DecodeBytes(payload.BlockRlp, block)
	header := block.Header()
	headerRlp, err := rlp.EncodeToBytes(header)
	if err != nil {
		return nil, err
	}
	convertedPayload := &IPLDPayload{
		BlockHash: block.Hash(),
		BlockNumber: block.Number(),
		HeaderRLP: headerRlp,
		BlockBody: block.Body(),
		Receipts:  make(types.Receipts, 0),
		StateLeafs: make(map[common.Hash][]byte),
		StorageLeafs: make(map[common.Hash]map[common.Hash][]byte),
	}
	for _, trx := range block.Transactions() {
		gethReceipt, err := pc.client.TransactionReceipt(context.Background(), trx.Hash())
		if err != nil {
			return nil, err
		}
		convertedPayload.Receipts = append(convertedPayload.Receipts, gethReceipt)
	}

	// Unpack state diff rlp to access fields
	stateDiff := new(statediff.StateDiff)
	err = rlp.DecodeBytes(payload.StateDiffRlp, stateDiff)
	if err != nil {
		return nil, err
	}
	for addr, createdAccount := range stateDiff.CreatedAccounts {
		convertedPayload.StateLeafs[addr] = createdAccount.Value
		convertedPayload.StorageLeafs[addr] = make(map[common.Hash][]byte)
		for _, storageDiff := range createdAccount.Storage {
			convertedPayload.StorageLeafs[addr][common.BytesToHash(storageDiff.Key)] = storageDiff.Value
		}
	}
	for addr, deletedAccount := range stateDiff.DeletedAccounts {
		convertedPayload.StateLeafs[addr] = deletedAccount.Value
		convertedPayload.StorageLeafs[addr] = make(map[common.Hash][]byte)
		for _, storageDiff := range deletedAccount.Storage {
			convertedPayload.StorageLeafs[addr][common.BytesToHash(storageDiff.Key)] = storageDiff.Value
		}
	}
	for addr, updatedAccount := range stateDiff.UpdatedAccounts {
		convertedPayload.StateLeafs[addr] = updatedAccount.Value
		convertedPayload.StorageLeafs[addr] = make(map[common.Hash][]byte)
		for _, storageDiff := range updatedAccount.Storage {
			convertedPayload.StorageLeafs[addr][common.BytesToHash(storageDiff.Key)] = storageDiff.Value
		}
	}
	return convertedPayload, nil
}

