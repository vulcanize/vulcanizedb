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
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type Repository interface {
	IndexCIDs(cidPayload *CIDPayload) error
}

type CIDRepository struct {
	db *postgres.DB
}

func NewCIDRepository(db *postgres.DB) *CIDRepository {
	return &CIDRepository{
		db: db,
	}
}

func (repo *CIDRepository) IndexCIDs(cidPayload *CIDPayload) error {
	headerID, err := repo.indexHeaderCID(cidPayload.HeaderCID, cidPayload.BlockNumber, cidPayload.BlockHash)
	if err != nil {
		return err
	}
	err = repo.indexTransactionAndReceiptCIDs(cidPayload, headerID)
	if err != nil {
		return err
	}

}


func (repo *CIDRepository) indexHeaderCID(cid, blockNumber, hash string) (int64, error) {
	var headerID int64
	err := repo.db.QueryRowx(`INSERT INTO public.header_cids (block_number, block_hash, cid) VALUES ($1, $2, $3) 
									RETURNING id
									ON CONFLICT DO UPDATE SET cid = $3`, blockNumber, hash, cid).Scan(&headerID)
	return headerID, err
}

func (repo *CIDRepository) indexTransactionAndReceiptCIDs(payload *CIDPayload, headerID int64) error {
	for hash, trxCid := range payload.TransactionCIDs {
		var txID int64
		err := repo.db.QueryRowx(`INSERT INTO public.transaction_cids (header_id, tx_hash, cid) VALUES ($1, $2, $3) 
										RETURNING id
										ON CONFLICT DO UPDATE SET cid = $3`, headerID, hash.Hex(), trxCid).Scan(&txID)
		if err != nil {
			return err
		}
		receiptCid, ok := payload.ReceiptCIDs[hash]
		if ok {
			err = repo.indexReceiptCID(receiptCid, txID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (repo *CIDRepository) indexReceiptCID(cid string, txId int64) error {
	_, err := repo.db.Exec(`INSERT INTO public.receipt_cids (tx_id, cid) VALUES ($1, $2) 
										ON CONFLICT DO UPDATE SET cid = $2`, txId, cid)
	return err
}

func (repo *CIDRepository) indexStateAndStorageCIDs(payload *CIDPayload, headerID int64) error {
	for accountKey, stateCID := range payload.StateLeafCIDs {
		var stateID int64
		err := repo.db.QueryRowx(`INSERT INTO public.state_cids (header_id, account_key, cid) VALUES ($1, $2, $3) 
										RETURNING id
										ON CONFLICT DO UPDATE SET cid = $3`, headerID, accountKey.Hex(), stateCID).Scan(&stateID)
		if err != nil {
			return err
		}
		for storageKey, storageCID := range payload.StorageLeafCIDs[accountKey] {
			err = repo.indexStorageCID(storageKey.Hex(), storageCID, stateID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (repo *CIDRepository) indexStorageCID(key, cid string, stateId int64) error {
	_, err := repo.db.Exec(`INSERT INTO public.storage_cids (state_id, storage_key, cid) VALUES ($1, $2, $3) 
										ON CONFLICT DO UPDATE SET cid = $3`, stateId, key, cid)
	return err
}