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

type CIDRetriever interface {
	RetrieveCIDs(streamFilters StreamFilters) ([]cidWrapper, error)
}

type EthCIDRetriever struct {
	db *postgres.DB
}

func NewCIDRetriever(db *postgres.DB) *EthCIDRetriever {
	return &EthCIDRetriever{
		db: db,
	}
}

func (ecr *EthCIDRetriever) GetLastBlockNumber() (int64, error) {
	var blockNumber int64
	err := ecr.db.Get(&blockNumber, "SELECT block_number FROM header_cids ORDER BY block_number DESC LIMIT 1 ")
	return blockNumber, err
}
func (ecr *EthCIDRetriever) RetrieveCIDs(streamFilters StreamFilters) ([]cidWrapper, error) {
	var endingBlock int64
	var err error
	if streamFilters.EndingBlock <= 0 || streamFilters.EndingBlock <= streamFilters.StartingBlock {
		endingBlock, err = ecr.GetLastBlockNumber()
		if err != nil {
			return nil, err
		}
	}
	cids := make([]cidWrapper, 0, endingBlock+1-streamFilters.StartingBlock)
	tx, err := ecr.db.Beginx()
	if err != nil {
		return nil, err
	}
	for i := streamFilters.StartingBlock; i <= endingBlock; i++ {
		cw := &cidWrapper{
			BlockNumber:  i,
			Headers:      make([]string, 0),
			Transactions: make([]string, 0),
			Receipts:     make([]string, 0),
			StateNodes:   make([]StateNodeCID, 0),
			StorageNodes: make([]StorageNodeCID, 0),
		}
		if !streamFilters.HeaderFilter.Off {
			err = ecr.retrieveHeaderCIDs(tx, streamFilters, cw, i)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
		var trxIds []int64
		if !streamFilters.TrxFilter.Off {
			trxIds, err = ecr.retrieveTrxCIDs(tx, streamFilters, cw, i)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
		if !streamFilters.ReceiptFilter.Off {
			err = ecr.retrieveRctCIDs(tx, streamFilters, cw, i, trxIds)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
		if !streamFilters.StateFilter.Off {
			err = ecr.retrieveStateCIDs(tx, streamFilters, cw, i)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
		if !streamFilters.StorageFilter.Off {
			err = ecr.retrieveStorageCIDs(tx, streamFilters, cw, i)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
		cids = append(cids, *cw)
	}

	return cids, err
}

func (ecr *EthCIDRetriever) retrieveHeaderCIDs(tx *sqlx.Tx, streamFilters StreamFilters, cids *cidWrapper, blockNumber int64) error {
	var pgStr string
	if streamFilters.HeaderFilter.FinalOnly {
		pgStr = `SELECT cid FROM header_cids
				WHERE block_number = $1
				AND final IS TRUE`
	} else {
		pgStr = `SELECT cid FROM header_cids
				WHERE block_number = $1`
	}
	return tx.Select(cids.Headers, pgStr, blockNumber)
}

func (ecr *EthCIDRetriever) retrieveTrxCIDs(tx *sqlx.Tx, streamFilters StreamFilters, cids *cidWrapper, blockNumber int64) ([]int64, error) {
	args := make([]interface{}, 0, 3)
	type result struct {
		Id  int64  `db:"id"`
		Cid string `db:"cid"`
	}
	results := make([]result, 0)
	pgStr := `SELECT transaction_cids.id, transaction_cids.cid FROM transaction_cids INNER JOIN header_cids ON (transaction_cids.header_id = header_cids.id)
			WHERE header_cids.block_number = $1`
	args = append(args, blockNumber)
	if len(streamFilters.TrxFilter.Dst) > 0 {
		pgStr += ` AND transaction_cids.dst = ANY($2::VARCHAR(66)[])`
		args = append(args, pq.Array(streamFilters.TrxFilter.Dst))
	}
	if len(streamFilters.TrxFilter.Src) > 0 {
		pgStr += ` AND transaction_cids.src = ANY($3::VARCHAR(66)[])`
		args = append(args, pq.Array(streamFilters.TrxFilter.Src))
	}
	err := tx.Select(results, pgStr, args...)
	if err != nil {
		return nil, err
	}
	ids := make([]int64, 0)
	for _, res := range results {
		cids.Transactions = append(cids.Transactions, res.Cid)
		ids = append(ids, res.Id)
	}
	return ids, nil
}

func (ecr *EthCIDRetriever) retrieveRctCIDs(tx *sqlx.Tx, streamFilters StreamFilters, cids *cidWrapper, blockNumber int64, trxIds []int64) error {
	args := make([]interface{}, 0, 2)
	pgStr := `SELECT receipt_cids.cid FROM receipt_cids, transaction_cids, header_cids
			WHERE receipt_cids.tx_id = transaction_cids.id 
			AND transaction_cids.header_id = header_cids.id
			AND header_cids.block_number = $1`
	args = append(args, blockNumber)
	if len(streamFilters.ReceiptFilter.Topic0s) > 0 {
		pgStr += ` AND (receipt_cids.topic0s && $2::VARCHAR(66)[]`
		args = append(args, pq.Array(streamFilters.ReceiptFilter.Topic0s))
	}
	if len(trxIds) > 0 {
		pgStr += ` OR receipt_cids.tx_id = ANY($3::INTEGER[]))`
		args = append(args, pq.Array(trxIds))
	} else {
		pgStr += `)`
	}
	return tx.Select(cids.Receipts, pgStr, args...)
}

func (ecr *EthCIDRetriever) retrieveStateCIDs(tx *sqlx.Tx, streamFilters StreamFilters, cids *cidWrapper, blockNumber int64) error {
	args := make([]interface{}, 0, 2)
	pgStr := `SELECT state_cids.cid, state_cids.state_key FROM state_cids INNER JOIN header_cids ON (state_cids.header_id = header_cids.id)
			WHERE header_cids.block_number = $1`
	args = append(args, blockNumber)
	addrLen := len(streamFilters.StateFilter.Addresses)
	if addrLen > 0 {
		keys := make([]string, 0, addrLen)
		for _, addr := range streamFilters.StateFilter.Addresses {
			keys = append(keys, HexToKey(addr).Hex())
		}
		pgStr += ` AND state_cids.state_key = ANY($2::VARCHAR(66)[])`
		args = append(args, pq.Array(keys))
	}
	return tx.Select(cids.StateNodes, pgStr, args...)
}

func (ecr *EthCIDRetriever) retrieveStorageCIDs(tx *sqlx.Tx, streamFilters StreamFilters, cids *cidWrapper, blockNumber int64) error {
	args := make([]interface{}, 0, 3)
	pgStr := `SELECT storage_cids.cid, state_cids.state_key, storage_cids.storage_key FROM storage_cids, state_cids, header_cids
			WHERE storage_cids.state_id = state_cids.id 
			AND state_cids.header_id = header_cids.id
			AND header_cids.block_number = $1`
	args = append(args, blockNumber)
	addrLen := len(streamFilters.StorageFilter.Addresses)
	if addrLen > 0 {
		keys := make([]string, 0, addrLen)
		for _, addr := range streamFilters.StorageFilter.Addresses {
			keys = append(keys, HexToKey(addr).Hex())
		}
		pgStr += ` AND state_cids.state_key = ANY($2::VARCHAR(66)[])`
		args = append(args, pq.Array(keys))
	}
	if len(streamFilters.StorageFilter.StorageKeys) > 0 {
		pgStr += ` AND storage_cids.storage_key = ANY($3::VARCHAR(66)[])`
		args = append(args, pq.Array(streamFilters.StorageFilter.StorageKeys))
	}
	return tx.Select(cids.StorageNodes, pgStr, args...)
}

// ADD IF LEAF ONLY!!
