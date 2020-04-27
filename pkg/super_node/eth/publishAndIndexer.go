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

	"github.com/ipfs/go-ipfs-blockstore"

	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/statediff"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipfs-ds-help"
	"github.com/jmoiron/sqlx"

	common2 "github.com/vulcanize/vulcanizedb/pkg/eth/converters/common"
	"github.com/vulcanize/vulcanizedb/pkg/ipfs/ipld"
	"github.com/vulcanize/vulcanizedb/pkg/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
)

// IPLDPublisherAndIndexer satisfies the IPLDPublisher interface for ethereum
// It interfaces directly with the public.blocks table of PG-IPFS rather than going through an ipfs intermediary
// It publishes and indexes IPLDs together in a single sqlx.Tx
type IPLDPublisherAndIndexer struct {
	db *postgres.DB
}

// NewIPLDPublisherAndIndexer creates a pointer to a new IPLDPublisherAndIndexer which satisfies the IPLDPublisher interface
func NewIPLDPublisherAndIndexer(db *postgres.DB) *IPLDPublisherAndIndexer {
	return &IPLDPublisherAndIndexer{
		db: db,
	}
}

// Publish publishes an IPLDPayload to IPFS and returns the corresponding CIDPayload
func (pub *IPLDPublisherAndIndexer) Publish(payload shared.ConvertedData) (shared.CIDsForIndexing, error) {
	ipldPayload, ok := payload.(ConvertedPayload)
	if !ok {
		return nil, fmt.Errorf("eth publisher expected payload type %T got %T", ConvertedPayload{}, payload)
	}
	// Generate the iplds
	headerNode, uncleNodes, txNodes, txTrieNodes, rctNodes, rctTrieNodes, err := ipld.FromBlockAndReceipts(ipldPayload.Block, ipldPayload.Receipts)
	if err != nil {
		return nil, err
	}

	// Begin new db tx
	tx, err := pub.db.Beginx()
	if err != nil {
		return nil, err
	}

	// Publish trie nodes
	for _, node := range txTrieNodes {
		if err := pub.publishIPLD(tx, node); err != nil {
			shared.Rollback(tx)
			return nil, err
		}
	}
	for _, node := range rctTrieNodes {
		if err := pub.publishIPLD(tx, node); err != nil {
			shared.Rollback(tx)
			return nil, err
		}
	}

	// Publish and index header
	if err := pub.publishIPLD(tx, headerNode); err != nil {
		shared.Rollback(tx)
		return nil, err
	}
	reward := common2.CalcEthBlockReward(ipldPayload.Block.Header(), ipldPayload.Block.Uncles(), ipldPayload.Block.Transactions(), ipldPayload.Receipts)
	header := HeaderModel{
		CID:             headerNode.Cid().String(),
		ParentHash:      ipldPayload.Block.ParentHash().String(),
		BlockNumber:     ipldPayload.Block.Number().String(),
		BlockHash:       ipldPayload.Block.Hash().String(),
		TotalDifficulty: ipldPayload.TotalDifficulty.String(),
		Reward:          reward.String(),
		Bloom:           ipldPayload.Block.Bloom().Bytes(),
		StateRoot:       ipldPayload.Block.Root().String(),
		RctRoot:         ipldPayload.Block.ReceiptHash().String(),
		TxRoot:          ipldPayload.Block.TxHash().String(),
		UncleRoot:       ipldPayload.Block.UncleHash().String(),
		Timestamp:       ipldPayload.Block.Time(),
	}
	headerID, err := pub.indexHeader(tx, header)
	if err != nil {
		shared.Rollback(tx)
		return nil, err
	}

	// Publish and index uncles
	for _, uncleNode := range uncleNodes {
		if err := pub.publishIPLD(tx, uncleNode); err != nil {
			shared.Rollback(tx)
			return nil, err
		}
		uncleReward := common2.CalcUncleMinerReward(ipldPayload.Block.Number().Int64(), uncleNode.Number.Int64())
		uncle := UncleModel{
			CID:        uncleNode.Cid().String(),
			ParentHash: uncleNode.ParentHash.String(),
			BlockHash:  uncleNode.Hash().String(),
			Reward:     uncleReward.String(),
		}
		if err := pub.indexUncle(tx, uncle, headerID); err != nil {
			shared.Rollback(tx)
			return nil, err
		}
	}

	// Publish and index txs and receipts
	for i, txNode := range txNodes {
		if err := pub.publishIPLD(tx, txNode); err != nil {
			shared.Rollback(tx)
			return nil, err
		}
		rctNode := rctNodes[i]
		if err := pub.publishIPLD(tx, rctNode); err != nil {
			shared.Rollback(tx)
			return nil, err
		}
		txModel := ipldPayload.TxMetaData[i]
		txModel.CID = txNode.Cid().String()
		txID, err := pub.indexTx(tx, txModel, headerID)
		if err != nil {
			shared.Rollback(tx)
			return nil, err
		}
		rctModel := ipldPayload.ReceiptMetaData[i]
		rctModel.CID = rctNode.Cid().String()
		if err := pub.indexRct(tx, rctModel, txID); err != nil {
			shared.Rollback(tx)
			return nil, err
		}
	}

	// Publish and index state and storage
	if err := pub.publishAndIndexStateAndStorage(tx, ipldPayload, headerID); err != nil {
		shared.Rollback(tx)
		return nil, err
	}

	// This IPLDPublisher does both publishing and indexing, we do not need to pass anything forward to the indexer
	return nil, tx.Commit()
}

func (pub *IPLDPublisherAndIndexer) publishAndIndexStateAndStorage(tx *sqlx.Tx, ipldPayload ConvertedPayload, headerID int64) error {
	// Publish and index state and storage
	for _, stateNode := range ipldPayload.StateNodes {
		stateIPLD, err := ipld.FromStateTrieRLP(stateNode.Value)
		if err != nil {
			return err
		}
		if err := pub.publishIPLD(tx, stateIPLD); err != nil {
			shared.Rollback(tx)
			return err
		}
		stateModel := StateNodeModel{
			Path:     stateNode.Path,
			StateKey: stateNode.LeafKey.String(),
			CID:      stateIPLD.Cid().String(),
			NodeType: ResolveFromNodeType(stateNode.Type),
		}
		stateID, err := pub.indexState(tx, stateModel, headerID)
		if err != nil {
			return err
		}
		// If we have a leaf, decode and index the account data and publish and index any associated storage diffs
		if stateNode.Type == statediff.Leaf {
			var i []interface{}
			if err := rlp.DecodeBytes(stateNode.Value, &i); err != nil {
				return err
			}
			if len(i) != 2 {
				return fmt.Errorf("IPLDPublisherAndIndexer expected state leaf node rlp to decode into two elements")
			}
			var account state.Account
			if err := rlp.DecodeBytes(i[1].([]byte), &account); err != nil {
				return err
			}
			accountModel := StateAccountModel{
				Balance:     account.Balance.String(),
				Nonce:       account.Nonce,
				CodeHash:    account.CodeHash,
				StorageRoot: account.Root.String(),
			}
			if err := pub.indexAccount(tx, accountModel, stateID); err != nil {
				return err
			}
			statePathHash := crypto.Keccak256Hash(stateNode.Path)
			for _, storageNode := range ipldPayload.StorageNodes[statePathHash] {
				storageIPLD, err := ipld.FromStorageTrieRLP(storageNode.Value)
				if err != nil {
					return err
				}
				if err := pub.publishIPLD(tx, storageIPLD); err != nil {
					return err
				}
				storageModel := StorageNodeModel{
					Path:       storageNode.Path,
					StorageKey: storageNode.LeafKey.Hex(),
					CID:        storageIPLD.Cid().String(),
					NodeType:   ResolveFromNodeType(storageNode.Type),
				}
				if err := pub.indexStorage(tx, storageModel, stateID); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// Index satisfies the shared.CIDIndexer interface
func (pub *IPLDPublisherAndIndexer) Index(cids shared.CIDsForIndexing) error {
	return nil
}

type ipldBase interface {
	Cid() cid.Cid
	RawData() []byte
}

func (pub *IPLDPublisherAndIndexer) publishIPLD(tx *sqlx.Tx, i ipldBase) error {
	dbKey := dshelp.CidToDsKey(i.Cid())
	prefixedKey := blockstore.BlockPrefix.String() + dbKey.String()
	raw := i.RawData()
	_, err := tx.Exec(`INSERT INTO public.blocks (key, data) VALUES ($1, $2) ON CONFLICT (key) DO NOTHING`, prefixedKey, raw)
	return err
}

func (pub *IPLDPublisherAndIndexer) generateAndPublishBlockIPLDs(tx *sqlx.Tx, body *types.Block, receipts types.Receipts) (*ipld.EthHeader,
	[]*ipld.EthHeader, []*ipld.EthTx, []*ipld.EthTxTrie, []*ipld.EthReceipt, []*ipld.EthRctTrie, error) {
	return ipld.FromBlockAndReceipts(body, receipts)
}

func (pub *IPLDPublisherAndIndexer) indexHeader(tx *sqlx.Tx, header HeaderModel) (int64, error) {
	var headerID int64
	err := tx.QueryRowx(`INSERT INTO eth.header_cids (block_number, block_hash, parent_hash, cid, td, node_id, reward, state_root, tx_root, receipt_root, uncle_root, bloom, timestamp, times_validated)
								VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
								ON CONFLICT (block_number, block_hash) DO UPDATE SET (parent_hash, cid, td, node_id, reward, state_root, tx_root, receipt_root, uncle_root, bloom, timestamp, times_validated) = ($3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, eth.header_cids.times_validated + 1)
								RETURNING id`,
		header.BlockNumber, header.BlockHash, header.ParentHash, header.CID, header.TotalDifficulty, pub.db.NodeID, header.Reward, header.StateRoot, header.TxRoot,
		header.RctRoot, header.UncleRoot, header.Bloom, header.Timestamp, 1).Scan(&headerID)
	return headerID, err
}

func (pub *IPLDPublisherAndIndexer) indexUncle(tx *sqlx.Tx, uncle UncleModel, headerID int64) error {
	_, err := tx.Exec(`INSERT INTO eth.uncle_cids (block_hash, header_id, parent_hash, cid, reward) VALUES ($1, $2, $3, $4, $5)
								ON CONFLICT (header_id, block_hash) DO UPDATE SET (parent_hash, cid, reward) = ($3, $4, $5)`,
		uncle.BlockHash, headerID, uncle.ParentHash, uncle.CID, uncle.Reward)
	return err
}

func (pub *IPLDPublisherAndIndexer) indexTx(tx *sqlx.Tx, transaction TxModel, headerID int64) (int64, error) {
	var txID int64
	err := tx.QueryRowx(`INSERT INTO eth.transaction_cids (header_id, tx_hash, cid, dst, src, index) VALUES ($1, $2, $3, $4, $5, $6)
									ON CONFLICT (header_id, tx_hash) DO UPDATE SET (cid, dst, src, index) = ($3, $4, $5, $6)
									RETURNING id`,
		headerID, transaction.TxHash, transaction.CID, transaction.Dst, transaction.Src, transaction.Index).Scan(&txID)
	return txID, err
}

func (pub *IPLDPublisherAndIndexer) indexRct(tx *sqlx.Tx, receipt ReceiptModel, txID int64) error {
	_, err := tx.Exec(`INSERT INTO eth.receipt_cids (tx_id, cid, contract, contract_hash, topic0s, topic1s, topic2s, topic3s, log_contracts) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
							  ON CONFLICT (tx_id) DO UPDATE SET (cid, contract, contract_hash, topic0s, topic1s, topic2s, topic3s, log_contracts) = ($2, $3, $4, $5, $6, $7, $8, $9)`,
		txID, receipt.CID, receipt.Contract, receipt.ContractHash, receipt.Topic0s, receipt.Topic1s, receipt.Topic2s, receipt.Topic3s, receipt.LogContracts)
	return err
}

func (pub *IPLDPublisherAndIndexer) indexState(tx *sqlx.Tx, stateNode StateNodeModel, headerID int64) (int64, error) {
	var stateID int64
	var stateKey string
	if stateNode.StateKey != nullHash.String() {
		stateKey = stateNode.StateKey
	}
	err := tx.QueryRowx(`INSERT INTO eth.state_cids (header_id, state_leaf_key, cid, state_path, node_type) VALUES ($1, $2, $3, $4, $5)
									ON CONFLICT (header_id, state_path) DO UPDATE SET (state_leaf_key, cid, node_type) = ($2, $3, $5)
									RETURNING id`,
		headerID, stateKey, stateNode.CID, stateNode.Path, stateNode.NodeType).Scan(&stateID)
	return stateID, err
}

func (pub *IPLDPublisherAndIndexer) indexStorage(tx *sqlx.Tx, storageNode StorageNodeModel, stateID int64) error {
	var storageKey string
	if storageNode.StorageKey != nullHash.String() {
		storageKey = storageNode.StorageKey
	}
	_, err := tx.Exec(`INSERT INTO eth.storage_cids (state_id, storage_leaf_key, cid, storage_path, node_type) VALUES ($1, $2, $3, $4, $5) 
							  ON CONFLICT (state_id, storage_path) DO UPDATE SET (storage_leaf_key, cid, node_type) = ($2, $3, $5)`,
		stateID, storageKey, storageNode.CID, storageNode.Path, storageNode.NodeType)
	return err
}

func (pub *IPLDPublisherAndIndexer) indexAccount(tx *sqlx.Tx, account StateAccountModel, stateID int64) error {
	_, err := tx.Exec(`INSERT INTO eth.state_accounts (state_id, balance, nonce, code_hash, storage_root) VALUES ($1, $2, $3, $4, $5)
							  ON CONFLICT (state_id) DO UPDATE SET (balance, nonce, code_hash, storage_root) = ($2, $3, $4, $5)`,
		stateID, account.Balance, account.Nonce, account.CodeHash, account.StorageRoot)
	return err
}
