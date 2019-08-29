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

package seed_node

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
)

// CIDRepository is an interface for indexing ipfs.CIDPayloads
type CIDRepository interface {
	Index(cidPayload *ipfs.CIDPayload) error
}

// Repository is the underlying struct for the CIDRepository interface
type Repository struct {
	db *postgres.DB
}

// NewCIDRepository creates a new pointer to a Repository which satisfies the CIDRepository interface
func NewCIDRepository(db *postgres.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// Index indexes a cidPayload in Postgres
func (repo *Repository) Index(cidPayload *ipfs.CIDPayload) error {
	tx, err := repo.db.Beginx()
	if err != nil {
		return err
	}
	headerID, err := repo.indexHeaderCID(tx, cidPayload.HeaderCID, cidPayload.BlockNumber, cidPayload.BlockHash.Hex())
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			log.Error(rollbackErr)
		}
		return err
	}
	for uncleHash, cid := range cidPayload.UncleCIDs {
		err = repo.indexUncleCID(tx, cid, cidPayload.BlockNumber, uncleHash.Hex())
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Error(rollbackErr)
			}
			return err
		}
	}
	err = repo.indexTransactionAndReceiptCIDs(tx, cidPayload, headerID)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			log.Error(rollbackErr)
		}
		return err
	}
	err = repo.indexStateAndStorageCIDs(tx, cidPayload, headerID)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			log.Error(rollbackErr)
		}
		return err
	}
	return tx.Commit()
}

func (repo *Repository) indexHeaderCID(tx *sqlx.Tx, cid, blockNumber, hash string) (int64, error) {
	var headerID int64
	err := tx.QueryRowx(`INSERT INTO public.header_cids (block_number, block_hash, cid, final) VALUES ($1, $2, $3, $4)
								ON CONFLICT (block_number, block_hash) DO UPDATE SET (cid, final) = ($3, $4)
								RETURNING id`,
		blockNumber, hash, cid, true).Scan(&headerID)
	return headerID, err
}

func (repo *Repository) indexUncleCID(tx *sqlx.Tx, cid, blockNumber, hash string) error {
	_, err := tx.Exec(`INSERT INTO public.header_cids (block_number, block_hash, cid, final) VALUES ($1, $2, $3, $4)
								ON CONFLICT (block_number, block_hash) DO UPDATE SET (cid, final) = ($3, $4)`,
		blockNumber, hash, cid, false)
	return err
}

func (repo *Repository) indexTransactionAndReceiptCIDs(tx *sqlx.Tx, payload *ipfs.CIDPayload, headerID int64) error {
	for hash, trxCidMeta := range payload.TransactionCIDs {
		var txID int64
		err := tx.QueryRowx(`INSERT INTO public.transaction_cids (header_id, tx_hash, cid, dst, src) VALUES ($1, $2, $3, $4, $5) 
									ON CONFLICT (header_id, tx_hash) DO UPDATE SET (cid, dst, src) = ($3, $4, $5)
									RETURNING id`,
			headerID, hash.Hex(), trxCidMeta.CID, trxCidMeta.Dst, trxCidMeta.Src).Scan(&txID)
		if err != nil {
			return err
		}
		receiptCidMeta, ok := payload.ReceiptCIDs[hash]
		if ok {
			err = repo.indexReceiptCID(tx, receiptCidMeta, txID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (repo *Repository) indexReceiptCID(tx *sqlx.Tx, cidMeta *ipfs.ReceiptMetaData, txID int64) error {
	_, err := tx.Exec(`INSERT INTO public.receipt_cids (tx_id, cid, contract, topic0s) VALUES ($1, $2, $3, $4)`,
		txID, cidMeta.CID, cidMeta.ContractAddress, pq.Array(cidMeta.Topic0s))
	return err
}

func (repo *Repository) indexStateAndStorageCIDs(tx *sqlx.Tx, payload *ipfs.CIDPayload, headerID int64) error {
	for accountKey, stateCID := range payload.StateNodeCIDs {
		var stateID int64
		err := tx.QueryRowx(`INSERT INTO public.state_cids (header_id, state_key, cid, leaf) VALUES ($1, $2, $3, $4)
									ON CONFLICT (header_id, state_key) DO UPDATE SET (cid, leaf) = ($3, $4)
									RETURNING id`,
			headerID, accountKey.Hex(), stateCID.CID, stateCID.Leaf).Scan(&stateID)
		if err != nil {
			return err
		}
		for _, storageCID := range payload.StorageNodeCIDs[accountKey] {
			err = repo.indexStorageCID(tx, storageCID, stateID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (repo *Repository) indexStorageCID(tx *sqlx.Tx, storageCID ipfs.StorageNodeCID, stateID int64) error {
	_, err := tx.Exec(`INSERT INTO public.storage_cids (state_id, storage_key, cid, leaf) VALUES ($1, $2, $3, $4) 
								   ON CONFLICT (state_id, storage_key) DO UPDATE SET (cid, leaf) = ($3, $4)`,
		stateID, storageCID.Key, storageCID.CID, storageCID.Leaf)
	return err
}
