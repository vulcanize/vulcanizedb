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

package eth

import (
	"fmt"

	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/statediff"
)

// PayloadConverter satisfies the PayloadConverter interface for ethereum
type PayloadConverter struct {
	chainConfig *params.ChainConfig
}

// NewPayloadConverter creates a pointer to a new PayloadConverter which satisfies the PayloadConverter interface
func NewPayloadConverter(chainConfig *params.ChainConfig) *PayloadConverter {
	return &PayloadConverter{
		chainConfig: chainConfig,
	}
}

// Convert method is used to convert a eth statediff.Payload to an IPLDPayload
// Satisfies the shared.PayloadConverter interface
func (pc *PayloadConverter) Convert(payload shared.RawChainData) (shared.StreamedIPLDs, error) {
	stateDiffPayload, ok := payload.(statediff.Payload)
	if !ok {
		return nil, fmt.Errorf("eth converter: expected payload type %T got %T", statediff.Payload{}, payload)
	}
	// Unpack block rlp to access fields
	block := new(types.Block)
	if err := rlp.DecodeBytes(stateDiffPayload.BlockRlp, block); err != nil {
		return nil, err
	}
	trxLen := len(block.Transactions())
	convertedPayload := IPLDPayload{
		TotalDifficulty: stateDiffPayload.TotalDifficulty,
		Block:           block,
		TxMetaData:      make([]TxModel, 0, trxLen),
		Receipts:        make(types.Receipts, 0, trxLen),
		ReceiptMetaData: make([]ReceiptModel, 0, trxLen),
		StateNodes:      make([]TrieNode, 0),
		StorageNodes:    make(map[common.Hash][]TrieNode),
	}
	signer := types.MakeSigner(pc.chainConfig, block.Number())
	transactions := block.Transactions()
	for i, trx := range transactions {
		// Extract to and from data from the the transactions for indexing
		from, err := types.Sender(signer, trx)
		if err != nil {
			return nil, err
		}
		txMeta := TxModel{
			Dst:    handleNullAddr(trx.To()),
			Src:    handleNullAddr(&from),
			TxHash: trx.Hash().String(),
			Index:  int64(i),
		}
		// txMeta will have same index as its corresponding trx in the convertedPayload.BlockBody
		convertedPayload.TxMetaData = append(convertedPayload.TxMetaData, txMeta)
	}

	// Decode receipts for this block
	receipts := make(types.Receipts, 0)
	if err := rlp.DecodeBytes(stateDiffPayload.ReceiptsRlp, &receipts); err != nil {
		return nil, err
	}
	// Derive any missing fields
	if err := receipts.DeriveFields(pc.chainConfig, block.Hash(), block.NumberU64(), block.Transactions()); err != nil {
		return nil, err
	}
	for i, receipt := range receipts {
		// If the transaction for this receipt has a "to" address, the above DeriveFields() fails to assign it to the receipt's ContractAddress
		// If it doesn't have a "to" address, it correctly derives it and assigns it to to the receipt's ContractAddress
		// Weird, right?
		if transactions[i].To() != nil {
			receipt.ContractAddress = *transactions[i].To()
		}
		// Extract topic and contract data from the receipt for indexing
		topicSets := make([][]string, 4)
		for _, log := range receipt.Logs {
			for i := range topicSets {
				if i < len(log.Topics) {
					topicSets[i] = append(topicSets[i], log.Topics[i].Hex())
				}
			}
		}
		rctMeta := ReceiptModel{
			Topic0s:  topicSets[0],
			Topic1s:  topicSets[1],
			Topic2s:  topicSets[2],
			Topic3s:  topicSets[3],
			Contract: receipt.ContractAddress.Hex(),
		}
		// receipt and rctMeta will have same indexes
		convertedPayload.Receipts = append(convertedPayload.Receipts, receipt)
		convertedPayload.ReceiptMetaData = append(convertedPayload.ReceiptMetaData, rctMeta)
	}

	// Unpack state diff rlp to access fields
	stateDiff := new(statediff.StateDiff)
	if err := rlp.DecodeBytes(stateDiffPayload.StateDiffRlp, stateDiff); err != nil {
		return nil, err
	}
	for _, createdAccount := range stateDiff.CreatedAccounts {
		hashKey := common.BytesToHash(createdAccount.Key)
		convertedPayload.StateNodes = append(convertedPayload.StateNodes, TrieNode{
			Key:   hashKey,
			Value: createdAccount.Value,
			Leaf:  createdAccount.Leaf,
		})
		for _, storageDiff := range createdAccount.Storage {
			convertedPayload.StorageNodes[hashKey] = append(convertedPayload.StorageNodes[hashKey], TrieNode{
				Key:   common.BytesToHash(storageDiff.Key),
				Value: storageDiff.Value,
				Leaf:  storageDiff.Leaf,
			})
		}
	}
	for _, deletedAccount := range stateDiff.DeletedAccounts {
		hashKey := common.BytesToHash(deletedAccount.Key)
		convertedPayload.StateNodes = append(convertedPayload.StateNodes, TrieNode{
			Key:   hashKey,
			Value: deletedAccount.Value,
			Leaf:  deletedAccount.Leaf,
		})
		for _, storageDiff := range deletedAccount.Storage {
			convertedPayload.StorageNodes[hashKey] = append(convertedPayload.StorageNodes[hashKey], TrieNode{
				Key:   common.BytesToHash(storageDiff.Key),
				Value: storageDiff.Value,
				Leaf:  storageDiff.Leaf,
			})
		}
	}
	for _, updatedAccount := range stateDiff.UpdatedAccounts {
		hashKey := common.BytesToHash(updatedAccount.Key)
		convertedPayload.StateNodes = append(convertedPayload.StateNodes, TrieNode{
			Key:   hashKey,
			Value: updatedAccount.Value,
			Leaf:  updatedAccount.Leaf,
		})
		for _, storageDiff := range updatedAccount.Storage {
			convertedPayload.StorageNodes[hashKey] = append(convertedPayload.StorageNodes[hashKey], TrieNode{
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
