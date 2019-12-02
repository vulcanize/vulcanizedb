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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/statediff"
)

// PayloadConverter interface is used to convert a geth statediff.Payload to our IPLDPayload type
type PayloadConverter interface {
	Convert(payload statediff.Payload) (*IPLDPayload, error)
}

// Converter is the underlying struct for the PayloadConverter interface
type Converter struct {
	chainConfig *params.ChainConfig
}

// NewPayloadConverter creates a pointer to a new Converter which satisfies the PayloadConverter interface
func NewPayloadConverter(chainConfig *params.ChainConfig) *Converter {
	return &Converter{
		chainConfig: chainConfig,
	}
}

// Convert method is used to convert a geth statediff.Payload to a IPLDPayload
func (pc *Converter) Convert(payload statediff.Payload) (*IPLDPayload, error) {
	// Unpack block rlp to access fields
	block := new(types.Block)
	decodeErr := rlp.DecodeBytes(payload.BlockRlp, block)
	if decodeErr != nil {
		return nil, decodeErr
	}
	header := block.Header()
	headerRlp, encodeErr := rlp.EncodeToBytes(header)
	if encodeErr != nil {
		return nil, encodeErr
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
	signer := types.MakeSigner(pc.chainConfig, block.Number())
	transactions := block.Transactions()
	for _, trx := range transactions {
		// Extract to and from data from the the transactions for indexing
		from, senderErr := types.Sender(signer, trx)
		if senderErr != nil {
			return nil, senderErr
		}
		txMeta := &TrxMetaData{
			Dst: handleNullAddr(trx.To()),
			Src: handleNullAddr(&from),
		}
		// txMeta will have same index as its corresponding trx in the convertedPayload.BlockBody
		convertedPayload.TrxMetaData = append(convertedPayload.TrxMetaData, txMeta)
	}

	// Decode receipts for this block
	receipts := make(types.Receipts, 0)
	decodeErr = rlp.DecodeBytes(payload.ReceiptsRlp, &receipts)
	if decodeErr != nil {
		return nil, decodeErr
	}
	// Derive any missing fields
	deriveErr := receipts.DeriveFields(pc.chainConfig, block.Hash(), block.NumberU64(), block.Transactions())
	if deriveErr != nil {
		return nil, deriveErr
	}
	for i, receipt := range receipts {
		// If the transaction for this receipt has a "to" address, the above DeriveFields() fails to assign it to the receipt's ContractAddress
		// If it doesn't have a "to" address, it correctly derives it and assigns it to to the receipt's ContractAddress
		// Weird, right?
		if transactions[i].To() != nil {
			receipt.ContractAddress = *transactions[i].To()
		}
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
	decodeErr = rlp.DecodeBytes(payload.StateDiffRlp, stateDiff)
	if decodeErr != nil {
		return nil, decodeErr
	}
	for _, createdAccount := range stateDiff.CreatedAccounts {
		hashKey := common.BytesToHash(createdAccount.Key)
		convertedPayload.StateNodes[hashKey] = StateNode{
			Value: createdAccount.Value,
			Leaf:  createdAccount.Leaf,
		}
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
