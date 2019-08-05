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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/statediff"

	"github.com/vulcanize/vulcanizedb/pkg/core"
)

// PayloadConverter interface is used to convert a geth statediff.Payload to our IPLDPayload type
type PayloadConverter interface {
	Convert(payload statediff.Payload) (*IPLDPayload, error)
}

// Converter is the underlying struct for the PayloadConverter interface
type Converter struct {
	client core.EthClient
}

// NewPayloadConverter creates a pointer to a new Converter which satisfies the PayloadConverter interface
func NewPayloadConverter(client core.EthClient) *Converter {
	return &Converter{
		client: client,
	}
}

// Convert method is used to convert a geth statediff.Payload to a IPLDPayload
func (pc *Converter) Convert(payload statediff.Payload) (*IPLDPayload, error) {
	// Unpack block rlp to access fields
	block := new(types.Block)
	err := rlp.DecodeBytes(payload.BlockRlp, block)
	header := block.Header()
	headerRlp, err := rlp.EncodeToBytes(header)
	if err != nil {
		return nil, err
	}
	trxLen := len(block.Transactions())
	convertedPayload := &IPLDPayload{
		BlockHash:       block.Hash(),
		BlockNumber:     block.Number(),
		HeaderRLP:       headerRlp,
		BlockBody:       block.Body(),
		TrxMetaData:     make([]*TrxMetaData, 0, trxLen),
		Receipts:        make(types.Receipts, 0, trxLen),
		ReceiptMetaData: make([]*ReceiptMetaData, 0, trxLen),
		StateNodes:      make(map[common.Hash]StateNode),
		StorageNodes:    make(map[common.Hash][]StorageNode),
	}
	for gethTransactionIndex, trx := range block.Transactions() {
		// Extract to and from data from the the transactions for indexing
		from, err := pc.client.TransactionSender(context.Background(), trx, block.Hash(), uint(gethTransactionIndex))
		if err != nil {
			return nil, err
		}
		txMeta := &TrxMetaData{
			Dst: handleNullAddr(trx.To()),
			Src: from.Hex(),
		}
		// txMeta will have same index as its corresponding trx in the convertedPayload.BlockBody
		convertedPayload.TrxMetaData = append(convertedPayload.TrxMetaData, txMeta)
	}

	// Decode receipts for this block
	receipts := make(types.Receipts, 0)
	err = rlp.DecodeBytes(payload.ReceiptsRlp, &receipts)
	if err != nil {
		return nil, err
	}
	for _, receipt := range receipts {
		// Extract topic0 data from the receipt's logs for indexing
		rctMeta := &ReceiptMetaData{
			Topic0s:         make([]string, 0, len(receipt.Logs)),
			ContractAddress: receipt.ContractAddress.Hex(),
		}
		for _, log := range receipt.Logs {
			if len(log.Topics) < 1 {
				continue
			}
			rctMeta.Topic0s = append(rctMeta.Topic0s, log.Topics[0].Hex())
		}
		// receipt and rctMeta will have same indexes
		convertedPayload.Receipts = append(convertedPayload.Receipts, receipt)
		convertedPayload.ReceiptMetaData = append(convertedPayload.ReceiptMetaData, rctMeta)
	}

	// Unpack state diff rlp to access fields
	stateDiff := new(statediff.StateDiff)
	err = rlp.DecodeBytes(payload.StateDiffRlp, stateDiff)
	if err != nil {
		return nil, err
	}
	for _, createdAccount := range stateDiff.CreatedAccounts {
		hashKey := common.BytesToHash(createdAccount.Key)
		convertedPayload.StateNodes[hashKey] = StateNode{
			Value: createdAccount.Value,
			Leaf:  createdAccount.Leaf,
		}
		convertedPayload.StorageNodes[hashKey] = make([]StorageNode, 0)
		for _, storageDiff := range createdAccount.Storage {
			convertedPayload.StorageNodes[hashKey] = append(convertedPayload.StorageNodes[hashKey], StorageNode{
				Key:   common.BytesToHash(storageDiff.Key),
				Value: storageDiff.Value,
				Leaf:  storageDiff.Leaf,
			})
		}
	}
	for _, deletedAccount := range stateDiff.DeletedAccounts {
		hashKey := common.BytesToHash(deletedAccount.Key)
		convertedPayload.StateNodes[hashKey] = StateNode{
			Value: deletedAccount.Value,
			Leaf:  deletedAccount.Leaf,
		}
		convertedPayload.StorageNodes[hashKey] = make([]StorageNode, 0)
		for _, storageDiff := range deletedAccount.Storage {
			convertedPayload.StorageNodes[hashKey] = append(convertedPayload.StorageNodes[hashKey], StorageNode{
				Key:   common.BytesToHash(storageDiff.Key),
				Value: storageDiff.Value,
				Leaf:  storageDiff.Leaf,
			})
		}
	}
	for _, updatedAccount := range stateDiff.UpdatedAccounts {
		hashKey := common.BytesToHash(updatedAccount.Key)
		convertedPayload.StateNodes[hashKey] = StateNode{
			Value: updatedAccount.Value,
			Leaf:  updatedAccount.Leaf,
		}
		convertedPayload.StorageNodes[hashKey] = make([]StorageNode, 0)
		for _, storageDiff := range updatedAccount.Storage {
			convertedPayload.StorageNodes[hashKey] = append(convertedPayload.StorageNodes[hashKey], StorageNode{
				Key:   common.BytesToHash(storageDiff.Key),
				Value: storageDiff.Value,
				Leaf:  storageDiff.Leaf,
			})
		}
	}
	return convertedPayload, nil
}

func handleNullAddr(to *common.Address) string {
	if to == nil {
		return "0x0000000000000000000000000000000000000000000000000000000000000000"
	}
	return to.Hex()
}
