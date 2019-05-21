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
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

// Still seeing some errors from tx and storage indexing processes... due to fk constraints being broken

// CIDRepository is an interface for indexing CIDPayloads
type CIDRepository interface {
	Index(cidPayload *CIDPayload) error
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
func (repo *Repository) Index(cidPayload *CIDPayload) error {
	tx, _ := repo.db.Beginx()
	headerID, err := repo.indexHeaderCID(tx, cidPayload.HeaderCID, cidPayload.BlockNumber, cidPayload.BlockHash.Hex())
	if err != nil {
		tx.Rollback()
		return err
	}
	for uncleHash, cid := range cidPayload.UncleCIDS {
		err = repo.indexUncleCID(tx, cid, cidPayload.BlockNumber, uncleHash.Hex())
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	err = repo.indexTransactionAndReceiptCIDs(cidPayload, headerID)
	if err != nil {
		return err
	}
	err = repo.indexStateAndStorageCIDs(cidPayload, headerID)
	if err != nil {
		return err
	}
	return nil
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
	_, err := tx.Queryx(`INSERT INTO public.header_cids (block_number, block_hash, cid, final) VALUES ($1, $2, $3, $4)
								ON CONFLICT (block_number, block_hash) DO UPDATE SET (cid, final) = ($3, $4)`,
		blockNumber, hash, cid, false)
	return err
}

func (repo *Repository) indexTransactionAndReceiptCIDs(payload *CIDPayload, headerID int64) error {
	tx, _ := repo.db.Beginx()
	for hash, trxCidMeta := range payload.TransactionCIDs {
		var txID int64
		err := tx.QueryRowx(`INSERT INTO public.transaction_cids (header_id, tx_hash, cid, dst, src) VALUES ($1, $2, $3, $4, $5) 
									ON CONFLICT (header_id, tx_hash) DO UPDATE SET (cid, dst, src) = ($3, $4, $5)
									RETURNING id`,
			headerID, hash.Hex(), trxCidMeta.CID, trxCidMeta.Dst, trxCidMeta.Src).Scan(&txID)
		if err != nil {
			tx.Rollback()
			return err
		}
		receiptCidMeta, ok := payload.ReceiptCIDs[hash]
		if ok {
			err = repo.indexReceiptCID(tx, receiptCidMeta, txID)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	return tx.Commit()
}

func (repo *Repository) indexReceiptCID(tx *sqlx.Tx, cidMeta *ReceiptMetaData, txID int64) error {
	_, err := tx.Exec(`INSERT INTO public.receipt_cids (tx_id, cid, topic0s) VALUES ($1, $2, $3)`,
		txID, cidMeta.CID, pq.Array(cidMeta.Topic0s))
	return err
}

func (repo *Repository) indexStateAndStorageCIDs(payload *CIDPayload, headerID int64) error {
	tx, _ := repo.db.Beginx()
	for accountKey, stateCID := range payload.StateNodeCIDs {
		var stateID int64
		err := tx.QueryRowx(`INSERT INTO public.state_cids (header_id, state_key, cid, leaf) VALUES ($1, $2, $3, $4)
									ON CONFLICT (header_id, state_key) DO UPDATE SET (cid, leaf) = ($3, $4)
									RETURNING id`,
			headerID, accountKey.Hex(), stateCID.CID, stateCID.Leaf).Scan(&stateID)
		if err != nil {
			tx.Rollback()
			return err
		}
		for _, storageCID := range payload.StorageNodeCIDs[accountKey] {
			err = repo.indexStorageCID(tx, storageCID, stateID)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	return tx.Commit()
}

func (repo *Repository) indexStorageCID(tx *sqlx.Tx, storageCID StorageNodeCID, stateID int64) error {
	_, err := repo.db.Exec(`INSERT INTO public.storage_cids (state_id, storage_key, cid, leaf) VALUES ($1, $2, $3, $4) 
								   ON CONFLICT (state_id, storage_key) DO UPDATE SET (cid, leaf) = ($3, $4)`,
		stateID, storageCID.Key, storageCID.CID, storageCID.Leaf)
	return err
}
