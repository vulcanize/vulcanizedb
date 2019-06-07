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
	log "github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

// CIDRetriever is the interface for retrieving CIDs from the Postgres cache
type CIDRetriever interface {
	RetrieveCIDs(streamFilters config.Subscription) ([]CidWrapper, error)
}

// EthCIDRetriever is the underlying struct supporting the CIDRetriever interface
type EthCIDRetriever struct {
	db *postgres.DB
}

// NewCIDRetriever returns a pointer to a new EthCIDRetriever which supports the CIDRetriever interface
func NewCIDRetriever(db *postgres.DB) *EthCIDRetriever {
	return &EthCIDRetriever{
		db: db,
	}
}

// GetLastBlockNumber is used to retrieve the latest block number in the cache
func (ecr *EthCIDRetriever) GetLastBlockNumber() (int64, error) {
	var blockNumber int64
	err := ecr.db.Get(&blockNumber, "SELECT block_number FROM header_cids ORDER BY block_number DESC LIMIT 1 ")
	return blockNumber, err
}

// RetrieveCIDs is used to retrieve all of the CIDs which conform to the passed StreamFilters
func (ecr *EthCIDRetriever) RetrieveCIDs(streamFilters config.Subscription) ([]CidWrapper, error) {
	log.Debug("retrieving cids")
	var endingBlock int64
	var err error
	if streamFilters.EndingBlock <= 0 || streamFilters.EndingBlock <= streamFilters.StartingBlock {
		endingBlock, err = ecr.GetLastBlockNumber()
		if err != nil {
			return nil, err
		}
	}
	cids := make([]CidWrapper, 0, endingBlock+1-streamFilters.StartingBlock)
	tx, err := ecr.db.Beginx()
	if err != nil {
		return nil, err
	}
	for i := streamFilters.StartingBlock; i <= endingBlock; i++ {
		cw := CidWrapper{}
		if !streamFilters.HeaderFilter.Off {
			cw.Headers, err = ecr.retrieveHeaderCIDs(tx, streamFilters, i)
			if err != nil {
				tx.Rollback()
				log.Error("header cid retrieval error")
				return nil, err
			}
		}
		var trxIds []int64
		if !streamFilters.TrxFilter.Off {
			cw.Transactions, trxIds, err = ecr.retrieveTrxCIDs(tx, streamFilters, i)
			if err != nil {
				tx.Rollback()
				log.Error("transaction cid retrieval error")
				return nil, err
			}
		}
		if !streamFilters.ReceiptFilter.Off {
			cw.Receipts, err = ecr.retrieveRctCIDs(tx, streamFilters, i, trxIds)
			if err != nil {
				tx.Rollback()
				log.Error("receipt cid retrieval error")
				return nil, err
			}
		}
		if !streamFilters.StateFilter.Off {
			cw.StateNodes, err = ecr.retrieveStateCIDs(tx, streamFilters, i)
			if err != nil {
				tx.Rollback()
				log.Error("state cid retrieval error")
				return nil, err
			}
		}
		if !streamFilters.StorageFilter.Off {
			cw.StorageNodes, err = ecr.retrieveStorageCIDs(tx, streamFilters, i)
			if err != nil {
				tx.Rollback()
				log.Error("storage cid retrieval error")
				return nil, err
			}
		}
		cids = append(cids, cw)
	}

	return cids, tx.Commit()
}

func (ecr *EthCIDRetriever) retrieveHeaderCIDs(tx *sqlx.Tx, streamFilters config.Subscription, blockNumber int64) ([]string, error) {
	log.Debug("retrieving header cids")
	headers := make([]string, 0)
	pgStr := `SELECT cid FROM header_cids
				WHERE block_number = $1`
	if streamFilters.HeaderFilter.FinalOnly {
		pgStr += ` AND final IS TRUE`
	}
	err := tx.Select(&headers, pgStr, blockNumber)
	return headers, err
}

func (ecr *EthCIDRetriever) retrieveTrxCIDs(tx *sqlx.Tx, streamFilters config.Subscription, blockNumber int64) ([]string, []int64, error) {
	log.Debug("retrieving transaction cids")
	args := make([]interface{}, 0, 3)
	type result struct {
		ID  int64  `db:"id"`
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
	err := tx.Select(&results, pgStr, args...)
	if err != nil {
		return nil, nil, err
	}
	ids := make([]int64, 0, len(results))
	cids := make([]string, 0, len(results))
	for _, res := range results {
		cids = append(cids, res.Cid)
		ids = append(ids, res.ID)
	}
	return cids, ids, nil
}

func (ecr *EthCIDRetriever) retrieveRctCIDs(tx *sqlx.Tx, streamFilters config.Subscription, blockNumber int64, trxIds []int64) ([]string, error) {
	log.Debug("retrieving receipt cids")
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
	receiptCids := make([]string, 0)
	err := tx.Select(&receiptCids, pgStr, args...)
	return receiptCids, err
}

func (ecr *EthCIDRetriever) retrieveStateCIDs(tx *sqlx.Tx, streamFilters config.Subscription, blockNumber int64) ([]StateNodeCID, error) {
	log.Debug("retrieving state cids")
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
	if !streamFilters.StorageFilter.IntermediateNodes {
		pgStr += ` AND state_cids.leaf = TRUE`
	}
	stateNodeCIDs := make([]StateNodeCID, 0)
	err := tx.Select(&stateNodeCIDs, pgStr, args...)
	return stateNodeCIDs, err
}

func (ecr *EthCIDRetriever) retrieveStorageCIDs(tx *sqlx.Tx, streamFilters config.Subscription, blockNumber int64) ([]StorageNodeCID, error) {
	log.Debug("retrieving storage cids")
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
	if !streamFilters.StorageFilter.IntermediateNodes {
		pgStr += ` AND storage_cids.leaf = TRUE`
	}
	storageNodeCIDs := make([]StorageNodeCID, 0)
	err := tx.Select(&storageNodeCIDs, pgStr, args...)
	return storageNodeCIDs, err
}
