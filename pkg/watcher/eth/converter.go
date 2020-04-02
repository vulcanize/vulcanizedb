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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"

	common2 "github.com/vulcanize/vulcanizedb/pkg/eth/converters/common"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/eth"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
)

// WatcherConverter converts watched data into models for the trigger tables
type WatcherConverter struct {
	chainConfig *params.ChainConfig
}

// NewWatcherConverter creates a pointer to a new WatcherConverter
func NewWatcherConverter(chainConfig *params.ChainConfig) *WatcherConverter {
	return &WatcherConverter{
		chainConfig: chainConfig,
	}
}

// Convert method is used to convert eth iplds to an cid payload
// Satisfies the shared.PayloadConverter interface
func (pc *WatcherConverter) Convert(ethIPLDs eth.IPLDs) (*eth.CIDPayload, error) {
	numTxs := len(ethIPLDs.Transactions)
	numRcts := len(ethIPLDs.Receipts)
	if numTxs != numRcts {
		return nil, fmt.Errorf("eth converter needs same numbe of receipts and transactions, have %d transactions and %d receipts", numTxs, numRcts)
	}
	// Initialize the payload struct and its fields
	cids := new(eth.CIDPayload)
	cids.UncleCIDs = make([]eth.UncleModel, len(ethIPLDs.Uncles))
	cids.TransactionCIDs = make([]eth.TxModel, numTxs)
	cids.ReceiptCIDs = make(map[common.Hash]eth.ReceiptModel, numTxs)
	cids.StateNodeCIDs = make([]eth.StateNodeModel, len(ethIPLDs.StateNodes))
	cids.StorageNodeCIDs = make(map[common.Hash][]eth.StorageNodeModel, len(ethIPLDs.StateNodes))

	// Unpack header
	var header types.Header
	if err := rlp.DecodeBytes(ethIPLDs.Header.Data, &header); err != nil {
		return nil, err
	}
	// Collect uncles so we can derive miner reward
	uncles := make([]*types.Header, len(ethIPLDs.Uncles))
	for i, uncleIPLD := range ethIPLDs.Uncles {
		var uncle types.Header
		if err := rlp.DecodeBytes(uncleIPLD.Data, &uncle); err != nil {
			return nil, err
		}
		uncleReward := common2.CalcUncleMinerReward(header.Number.Int64(), uncle.Number.Int64())
		uncles[i] = &uncle
		// Uncle data
		cids.UncleCIDs[i] = eth.UncleModel{
			CID:        uncleIPLD.CID,
			BlockHash:  uncle.Hash().String(),
			ParentHash: uncle.ParentHash.String(),
			Reward:     uncleReward.String(),
		}
	}
	// Collect transactions so we can derive receipt fields and miner reward
	signer := types.MakeSigner(pc.chainConfig, header.Number)
	transactions := make(types.Transactions, len(ethIPLDs.Transactions))
	for i, txIPLD := range ethIPLDs.Transactions {
		var tx types.Transaction
		if err := rlp.DecodeBytes(txIPLD.Data, &tx); err != nil {
			return nil, err
		}
		transactions[i] = &tx
		from, err := types.Sender(signer, &tx)
		if err != nil {
			return nil, err
		}
		// Tx data
		cids.TransactionCIDs[i] = eth.TxModel{
			Dst:    shared.HandleNullAddrPointer(tx.To()),
			Src:    shared.HandleNullAddr(from),
			TxHash: tx.Hash().String(),
			Index:  int64(i),
			CID:    txIPLD.CID,
		}
	}
	// Collect receipts so that we can derive the rest of their fields and miner reward
	receipts := make(types.Receipts, len(ethIPLDs.Receipts))
	for i, rctIPLD := range ethIPLDs.Receipts {
		var rct types.Receipt
		if err := rlp.DecodeBytes(rctIPLD.Data, &rct); err != nil {
			return nil, err
		}
		receipts[i] = &rct
	}
	if err := receipts.DeriveFields(pc.chainConfig, header.Hash(), header.Number.Uint64(), transactions); err != nil {
		return nil, err
	}
	for i, receipt := range receipts {
		matchedTx := transactions[i]
		topicSets := make([][]string, 4)
		mappedContracts := make(map[string]bool) // use map to avoid duplicate addresses
		for _, log := range receipt.Logs {
			for i, topic := range log.Topics {
				topicSets[i] = append(topicSets[i], topic.Hex())
			}
			mappedContracts[log.Address.String()] = true
		}
		logContracts := make([]string, 0, len(mappedContracts))
		for addr := range mappedContracts {
			logContracts = append(logContracts, addr)
		}
		// Rct data
		cids.ReceiptCIDs[matchedTx.Hash()] = eth.ReceiptModel{
			CID:          ethIPLDs.Receipts[i].CID,
			Topic0s:      topicSets[0],
			Topic1s:      topicSets[1],
			Topic2s:      topicSets[2],
			Topic3s:      topicSets[3],
			Contract:     receipt.ContractAddress.Hex(),
			LogContracts: logContracts,
		}
	}
	minerReward := common2.CalcEthBlockReward(&header, uncles, transactions, receipts)
	// Header data
	cids.HeaderCID = eth.HeaderModel{
		CID:             ethIPLDs.Header.CID,
		ParentHash:      header.ParentHash.String(),
		BlockHash:       header.Hash().String(),
		BlockNumber:     header.Number.String(),
		TotalDifficulty: ethIPLDs.TotalDifficulty.String(),
		Reward:          minerReward.String(),
	}
	// State data
	for i, stateIPLD := range ethIPLDs.StateNodes {
		cids.StateNodeCIDs[i] = eth.StateNodeModel{
			CID:      stateIPLD.IPLD.CID,
			NodeType: eth.ResolveFromNodeType(stateIPLD.Type),
			StateKey: stateIPLD.StateLeafKey.String(),
		}
	}
	// Storage data
	for _, storageIPLD := range ethIPLDs.StorageNodes {
		cids.StorageNodeCIDs[storageIPLD.StateLeafKey] = append(cids.StorageNodeCIDs[storageIPLD.StateLeafKey], eth.StorageNodeModel{
			CID:        storageIPLD.IPLD.CID,
			NodeType:   eth.ResolveFromNodeType(storageIPLD.Type),
			StorageKey: storageIPLD.StorageLeafKey.String(),
		})
	}
	return cids, nil
}
