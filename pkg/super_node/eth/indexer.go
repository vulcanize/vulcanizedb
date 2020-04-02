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
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/pkg/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
)

var (
	nullHash = common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")
)

// Indexer satisfies the Indexer interface for ethereum
type CIDIndexer struct {
	db *postgres.DB
}

// NewCIDIndexer creates a new pointer to a Indexer which satisfies the CIDIndexer interface
func NewCIDIndexer(db *postgres.DB) *CIDIndexer {
	return &CIDIndexer{
		db: db,
	}
}

// Index indexes a cidPayload in Postgres
func (in *CIDIndexer) Index(cids shared.CIDsForIndexing) error {
	cidPayload, ok := cids.(*CIDPayload)
	if !ok {
		return fmt.Errorf("eth indexer expected cids type %T got %T", &CIDPayload{}, cids)
	}
	tx, err := in.db.Beginx()
	if err != nil {
		return err
	}
	headerID, err := in.indexHeaderCID(tx, cidPayload.HeaderCID, in.db.NodeID)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			log.Error(err)
		}
		log.Error("eth indexer error when indexing header")
		return err
	}
	for _, uncle := range cidPayload.UncleCIDs {
		if err := in.indexUncleCID(tx, uncle, headerID); err != nil {
			if err := tx.Rollback(); err != nil {
				log.Error(err)
			}
			log.Error("eth indexer error when indexing uncle")
			return err
		}
	}
	if err := in.indexTransactionAndReceiptCIDs(tx, cidPayload, headerID); err != nil {
		if err := tx.Rollback(); err != nil {
			log.Error(err)
		}
		log.Error("eth indexer error when indexing transactions and receipts")
		return err
	}
	if err := in.indexStateAndStorageCIDs(tx, cidPayload, headerID); err != nil {
		if err := tx.Rollback(); err != nil {
			log.Error(err)
		}
		log.Error("eth indexer error when indexing state and storage nodes")
		return err
	}
	return tx.Commit()
}

func (in *CIDIndexer) indexHeaderCID(tx *sqlx.Tx, header HeaderModel, nodeID int64) (int64, error) {
	var headerID int64
	err := tx.QueryRowx(`INSERT INTO eth.header_cids (block_number, block_hash, parent_hash, cid, td, node_id, reward, state_root, tx_root, receipt_root, uncle_root, bloom, timestamp) 
								VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
								ON CONFLICT (block_number, block_hash) DO UPDATE SET (parent_hash, cid, td, node_id, reward, state_root, tx_root, receipt_root, uncle_root, bloom, timestamp) = ($3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
								RETURNING id`,
		header.BlockNumber, header.BlockHash, header.ParentHash, header.CID, header.TotalDifficulty, nodeID, header.Reward, header.StateRoot, header.TxRoot,
		header.RctRoot, header.UncleRoot, header.Bloom, header.Timestamp).Scan(&headerID)
	return headerID, err
}

func (in *CIDIndexer) indexUncleCID(tx *sqlx.Tx, uncle UncleModel, headerID int64) error {
	_, err := tx.Exec(`INSERT INTO eth.uncle_cids (block_hash, header_id, parent_hash, cid, reward) VALUES ($1, $2, $3, $4, $5)
								ON CONFLICT (header_id, block_hash) DO UPDATE SET (parent_hash, cid, reward) = ($3, $4, $5)`,
		uncle.BlockHash, headerID, uncle.ParentHash, uncle.CID, uncle.Reward)
	return err
}

func (in *CIDIndexer) indexTransactionAndReceiptCIDs(tx *sqlx.Tx, payload *CIDPayload, headerID int64) error {
	for _, trxCidMeta := range payload.TransactionCIDs {
		var txID int64
		err := tx.QueryRowx(`INSERT INTO eth.transaction_cids (header_id, tx_hash, cid, dst, src, index) VALUES ($1, $2, $3, $4, $5, $6)
									ON CONFLICT (header_id, tx_hash) DO UPDATE SET (cid, dst, src, index) = ($3, $4, $5, $6)
									RETURNING id`,
			headerID, trxCidMeta.TxHash, trxCidMeta.CID, trxCidMeta.Dst, trxCidMeta.Src, trxCidMeta.Index).Scan(&txID)
		if err != nil {
			return err
		}
		receiptCidMeta, ok := payload.ReceiptCIDs[common.HexToHash(trxCidMeta.TxHash)]
		if ok {
			if err := in.indexReceiptCID(tx, receiptCidMeta, txID); err != nil {
				return err
			}
		}
	}
	return nil
}

func (in *CIDIndexer) indexReceiptCID(tx *sqlx.Tx, cidMeta ReceiptModel, txID int64) error {
	_, err := tx.Exec(`INSERT INTO eth.receipt_cids (tx_id, cid, contract, topic0s, topic1s, topic2s, topic3s, log_contracts) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		txID, cidMeta.CID, cidMeta.Contract, cidMeta.Topic0s, cidMeta.Topic1s, cidMeta.Topic2s, cidMeta.Topic3s, cidMeta.LogContracts)
	return err
}

func (in *CIDIndexer) indexStateAndStorageCIDs(tx *sqlx.Tx, payload *CIDPayload, headerID int64) error {
	for _, stateCID := range payload.StateNodeCIDs {
		var stateID int64
		var stateKey string
		if stateCID.StateKey != nullHash.String() {
			stateKey = stateCID.StateKey
		}
		err := tx.QueryRowx(`INSERT INTO eth.state_cids (header_id, state_leaf_key, cid, state_path, node_type) VALUES ($1, $2, $3, $4, $5)
									ON CONFLICT (header_id, state_path) DO UPDATE SET (state_leaf_key, cid, node_type) = ($2, $3, $5)
									RETURNING id`,
			headerID, stateKey, stateCID.CID, stateCID.Path, stateCID.NodeType).Scan(&stateID)
		if err != nil {
			return err
		}
		// If we have a state leaf node, index the associated account and storage nodes
		if stateCID.NodeType == 2 {
			pathKey := crypto.Keccak256Hash(stateCID.Path)
			for _, storageCID := range payload.StorageNodeCIDs[pathKey] {
				if err := in.indexStorageCID(tx, storageCID, stateID); err != nil {
					return err
				}
			}
			if stateAccount, ok := payload.StateAccounts[pathKey]; ok {
				if err := in.indexStateAccount(tx, stateAccount, stateID); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (in *CIDIndexer) indexStateAccount(tx *sqlx.Tx, stateAccount StateAccountModel, stateID int64) error {
	_, err := tx.Exec(`INSERT INTO eth.state_accounts (state_id, balance, nonce, code_hash, storage_root) VALUES ($1, $2, $3, $4, $5)
								ON CONFLICT (state_id) DO UPDATE SET (balance, nonce, code_hash, storage_root) = ($2, $3, $4, $5)`,
		stateID, stateAccount.Balance, stateAccount.Nonce, stateAccount.CodeHash, stateAccount.StorageRoot)
	return err
}

func (in *CIDIndexer) indexStorageCID(tx *sqlx.Tx, storageCID StorageNodeModel, stateID int64) error {
	var storageKey string
	if storageCID.StorageKey != nullHash.String() {
		storageKey = storageCID.StorageKey
	}
	_, err := tx.Exec(`INSERT INTO eth.storage_cids (state_id, storage_leaf_key, cid, storage_path, node_type) VALUES ($1, $2, $3, $4, $5) 
								   ON CONFLICT (state_id, storage_path) DO UPDATE SET (storage_leaf_key, cid, node_type) = ($2, $3, $5)`,
		stateID, storageKey, storageCID.CID, storageCID.Path, storageCID.NodeType)
	return err
}
