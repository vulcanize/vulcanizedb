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

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

// Repository is an interface for indexing CIDPayloads
type Repository interface {
	Index(cidPayload *CIDPayload) error
}

// CIDRepository is the underlying struct for the Repository interface
type CIDRepository struct {
	db *postgres.DB
}

// NewCIDRepository creates a new pointer to a CIDRepository which satisfies the Repository interface
func NewCIDRepository(db *postgres.DB) *CIDRepository {
	return &CIDRepository{
		db: db,
	}
}

// IndexCIDs indexes a cidPayload in Postgres
func (repo *CIDRepository) Index(cidPayload *CIDPayload) error {
	tx, _ := repo.db.Beginx()
	headerID, err := repo.indexHeaderCID(tx, cidPayload.HeaderCID, cidPayload.BlockNumber, cidPayload.BlockHash)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = repo.indexTransactionAndReceiptCIDs(tx, cidPayload, headerID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = repo.indexStateAndStorageCIDs(tx, cidPayload, headerID)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (repo *CIDRepository) indexHeaderCID(tx *sqlx.Tx, cid, blockNumber, hash string) (int64, error) {
	var headerID int64
	err := tx.QueryRowx(`INSERT INTO public.header_cids (block_number, block_hash, cid) VALUES ($1, $2, $3)
								ON CONFLICT DO UPDATE SET cid = $3
								RETURNING id`,
		blockNumber, hash, cid).Scan(&headerID)
	return headerID, err
}

func (repo *CIDRepository) indexTransactionAndReceiptCIDs(tx *sqlx.Tx, payload *CIDPayload, headerID int64) error {
	for hash, trxCid := range payload.TransactionCIDs {
		var txID int64
		err := tx.QueryRowx(`INSERT INTO public.transaction_cids (header_id, tx_hash, cid) VALUES ($1, $2, $3) 
										ON CONFLICT DO UPDATE SET cid = $3
										RETURNING id`,
			headerID, hash.Hex(), trxCid).Scan(&txID)
		if err != nil {
			return err
		}
		receiptCid, ok := payload.ReceiptCIDs[hash]
		if ok {
			err = repo.indexReceiptCID(tx, receiptCid, txID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (repo *CIDRepository) indexReceiptCID(tx *sqlx.Tx, cid string, txId int64) error {
	_, err := tx.Exec(`INSERT INTO public.receipt_cids (tx_id, cid) VALUES ($1, $2) 
										ON CONFLICT DO UPDATE SET cid = $2`, txId, cid)
	return err
}

func (repo *CIDRepository) indexStateAndStorageCIDs(tx *sqlx.Tx, payload *CIDPayload, headerID int64) error {
	for accountKey, stateCID := range payload.StateLeafCIDs {
		var stateID int64
		err := tx.QueryRowx(`INSERT INTO public.state_cids (header_id, account_key, cid) VALUES ($1, $2, $3) 
										ON CONFLICT DO UPDATE SET cid = $3
										RETURNING id`,
			headerID, accountKey.Hex(), stateCID).Scan(&stateID)
		if err != nil {
			return err
		}
		for storageKey, storageCID := range payload.StorageLeafCIDs[accountKey] {
			err = repo.indexStorageCID(tx, storageKey.Hex(), storageCID, stateID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (repo *CIDRepository) indexStorageCID(tx *sqlx.Tx, key, cid string, stateId int64) error {
	_, err := repo.db.Exec(`INSERT INTO public.storage_cids (state_id, storage_key, cid) VALUES ($1, $2, $3) 
										ON CONFLICT DO UPDATE SET cid = $3`, stateId, key, cid)
	return err
}
