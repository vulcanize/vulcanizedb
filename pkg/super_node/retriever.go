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

package super_node

import (
	"math/big"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
)

// CIDRetriever is the interface for retrieving CIDs from the Postgres cache
type CIDRetriever interface {
	RetrieveCIDs(streamFilters config.Subscription, blockNumber int64) (*ipfs.CIDWrapper, error)
	RetrieveLastBlockNumber() (int64, error)
	RetrieveFirstBlockNumber() (int64, error)
	RetrieveGapsInData() ([][2]int64, error)
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

// RetrieveFirstBlockNumber is used to retrieve the first block number in the db
func (ecr *EthCIDRetriever) RetrieveFirstBlockNumber() (int64, error) {
	var blockNumber int64
	err := ecr.db.Get(&blockNumber, "SELECT block_number FROM header_cids ORDER BY block_number ASC LIMIT 1")
	return blockNumber, err
}

// RetrieveLastBlockNumber is used to retrieve the latest block number in the db
func (ecr *EthCIDRetriever) RetrieveLastBlockNumber() (int64, error) {
	var blockNumber int64
	err := ecr.db.Get(&blockNumber, "SELECT block_number FROM header_cids ORDER BY block_number DESC LIMIT 1 ")
	return blockNumber, err
}

// RetrieveCIDs is used to retrieve all of the CIDs which conform to the passed StreamFilters
func (ecr *EthCIDRetriever) RetrieveCIDs(streamFilters config.Subscription, blockNumber int64) (*ipfs.CIDWrapper, error) {
	log.Debug("retrieving cids")
	tx, beginErr := ecr.db.Beginx()
	if beginErr != nil {
		return nil, beginErr
	}
	// THIS IS SUPER EXPENSIVE HAVING TO CYCLE THROUGH EACH BLOCK, NEED BETTER WAY TO FETCH CIDS
	// WHILE STILL MAINTAINING RELATION INFO ABOUT WHAT BLOCK THE CIDS BELONG TO
	cw := new(ipfs.CIDWrapper)
	cw.BlockNumber = big.NewInt(blockNumber)

	// Retrieve cached header CIDs
	if !streamFilters.HeaderFilter.Off {
		var headersErr error
		cw.Headers, headersErr = ecr.retrieveHeaderCIDs(tx, streamFilters, blockNumber)
		if headersErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Error(rollbackErr)
			}
			log.Error("header cid retrieval error")
			return nil, headersErr
		}
		if !streamFilters.HeaderFilter.FinalOnly {
			var unclesErr error
			cw.Uncles, unclesErr = ecr.retrieveUncleCIDs(tx, streamFilters, blockNumber)
			if unclesErr != nil {
				rollbackErr := tx.Rollback()
				if rollbackErr != nil {
					log.Error(rollbackErr)
				}
				log.Error("uncle cid retrieval error")
				return nil, unclesErr
			}
		}
	}

	// Retrieve cached trx CIDs
	var trxIds []int64
	if !streamFilters.TrxFilter.Off {
		var trxsErr error
		cw.Transactions, trxIds, trxsErr = ecr.retrieveTrxCIDs(tx, streamFilters, blockNumber)
		if trxsErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Error(rollbackErr)
			}
			log.Error("transaction cid retrieval error")
			return nil, trxsErr
		}
	}

	// Retrieve cached receipt CIDs
	if !streamFilters.ReceiptFilter.Off {
		var rctsErr error
		cw.Receipts, rctsErr = ecr.retrieveRctCIDs(tx, streamFilters, blockNumber, trxIds)
		if rctsErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Error(rollbackErr)
			}
			log.Error("receipt cid retrieval error")
			return nil, rctsErr
		}
	}

	// Retrieve cached state CIDs
	if !streamFilters.StateFilter.Off {
		var stateErr error
		cw.StateNodes, stateErr = ecr.retrieveStateCIDs(tx, streamFilters, blockNumber)
		if stateErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Error(rollbackErr)
			}
			log.Error("state cid retrieval error")
			return nil, stateErr
		}
	}

	// Retrieve cached storage CIDs
	if !streamFilters.StorageFilter.Off {
		var storageErr error
		cw.StorageNodes, storageErr = ecr.retrieveStorageCIDs(tx, streamFilters, blockNumber)
		if storageErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Error(rollbackErr)
			}
			log.Error("storage cid retrieval error")
			return nil, storageErr
		}
	}

	return cw, tx.Commit()
}

func (ecr *EthCIDRetriever) retrieveHeaderCIDs(tx *sqlx.Tx, streamFilters config.Subscription, blockNumber int64) ([]string, error) {
	log.Debug("retrieving header cids for block ", blockNumber)
	headers := make([]string, 0)
	pgStr := `SELECT cid FROM header_cids
				WHERE block_number = $1 AND final IS TRUE`
	err := tx.Select(&headers, pgStr, blockNumber)
	return headers, err
}

func (ecr *EthCIDRetriever) retrieveUncleCIDs(tx *sqlx.Tx, streamFilters config.Subscription, blockNumber int64) ([]string, error) {
	log.Debug("retrieving header cids for block ", blockNumber)
	headers := make([]string, 0)
	pgStr := `SELECT cid FROM header_cids
				WHERE block_number = $1 AND final IS FALSE`
	err := tx.Select(&headers, pgStr, blockNumber)
	return headers, err
}

func (ecr *EthCIDRetriever) retrieveTrxCIDs(tx *sqlx.Tx, streamFilters config.Subscription, blockNumber int64) ([]string, []int64, error) {
	log.Debug("retrieving transaction cids for block ", blockNumber)
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
	log.Debug("retrieving receipt cids for block ", blockNumber)
	args := make([]interface{}, 0, 4)
	pgStr := `SELECT receipt_cids.cid FROM receipt_cids, transaction_cids, header_cids
			WHERE receipt_cids.tx_id = transaction_cids.id 
			AND transaction_cids.header_id = header_cids.id
			AND header_cids.block_number = $1`
	args = append(args, blockNumber)
	if len(streamFilters.ReceiptFilter.Topic0s) > 0 {
		pgStr += ` AND ((receipt_cids.topic0s && $2::VARCHAR(66)[]`
		args = append(args, pq.Array(streamFilters.ReceiptFilter.Topic0s))
		if len(streamFilters.ReceiptFilter.Contracts) > 0 {
			pgStr += ` AND receipt_cids.contract = ANY($3::VARCHAR(66)[]))`
			args = append(args, pq.Array(streamFilters.ReceiptFilter.Contracts))
			if len(trxIds) > 0 {
				pgStr += ` OR receipt_cids.tx_id = ANY($4::INTEGER[]))`
				args = append(args, pq.Array(trxIds))
			} else {
				pgStr += `)`
			}
		} else {
			pgStr += `)`
			if len(trxIds) > 0 {
				pgStr += ` OR receipt_cids.tx_id = ANY($3::INTEGER[]))`
				args = append(args, pq.Array(trxIds))
			} else {
				pgStr += `)`
			}
		}
	} else {
		if len(streamFilters.ReceiptFilter.Contracts) > 0 {
			pgStr += ` AND (receipt_cids.contract = ANY($2::VARCHAR(66)[])`
			args = append(args, pq.Array(streamFilters.ReceiptFilter.Contracts))
			if len(trxIds) > 0 {
				pgStr += ` OR receipt_cids.tx_id = ANY($3::INTEGER[]))`
				args = append(args, pq.Array(trxIds))
			} else {
				pgStr += `)`
			}
		} else if len(trxIds) > 0 {
			pgStr += ` AND receipt_cids.tx_id = ANY($2::INTEGER[])`
			args = append(args, pq.Array(trxIds))
		}
	}
	receiptCids := make([]string, 0)
	err := tx.Select(&receiptCids, pgStr, args...)
	return receiptCids, err
}

func (ecr *EthCIDRetriever) retrieveStateCIDs(tx *sqlx.Tx, streamFilters config.Subscription, blockNumber int64) ([]ipfs.StateNodeCID, error) {
	log.Debug("retrieving state cids for block ", blockNumber)
	args := make([]interface{}, 0, 2)
	pgStr := `SELECT state_cids.cid, state_cids.state_key, state_cids.leaf FROM state_cids INNER JOIN header_cids ON (state_cids.header_id = header_cids.id)
			WHERE header_cids.block_number = $1`
	args = append(args, blockNumber)
	addrLen := len(streamFilters.StateFilter.Addresses)
	if addrLen > 0 {
		keys := make([]string, 0, addrLen)
		for _, addr := range streamFilters.StateFilter.Addresses {
			keys = append(keys, ipfs.HexToKey(addr).Hex())
		}
		pgStr += ` AND state_cids.state_key = ANY($2::VARCHAR(66)[])`
		args = append(args, pq.Array(keys))
	}
	if !streamFilters.StorageFilter.IntermediateNodes {
		pgStr += ` AND state_cids.leaf = TRUE`
	}
	stateNodeCIDs := make([]ipfs.StateNodeCID, 0)
	err := tx.Select(&stateNodeCIDs, pgStr, args...)
	return stateNodeCIDs, err
}

func (ecr *EthCIDRetriever) retrieveStorageCIDs(tx *sqlx.Tx, streamFilters config.Subscription, blockNumber int64) ([]ipfs.StorageNodeCID, error) {
	log.Debug("retrieving storage cids for block ", blockNumber)
	args := make([]interface{}, 0, 3)
	pgStr := `SELECT storage_cids.cid, state_cids.state_key, storage_cids.storage_key, storage_cids.leaf FROM storage_cids, state_cids, header_cids
			WHERE storage_cids.state_id = state_cids.id 
			AND state_cids.header_id = header_cids.id
			AND header_cids.block_number = $1`
	args = append(args, blockNumber)
	addrLen := len(streamFilters.StorageFilter.Addresses)
	if addrLen > 0 {
		keys := make([]string, 0, addrLen)
		for _, addr := range streamFilters.StorageFilter.Addresses {
			keys = append(keys, ipfs.HexToKey(addr).Hex())
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
	storageNodeCIDs := make([]ipfs.StorageNodeCID, 0)
	err := tx.Select(&storageNodeCIDs, pgStr, args...)
	return storageNodeCIDs, err
}

type gap struct {
	Start int64 `db:"start"`
	Stop  int64 `db:"stop"`
}

// RetrieveGapsInData is used to find the the block numbers at which we are missing data in the db
func (ecr *EthCIDRetriever) RetrieveGapsInData() ([][2]int64, error) {
	pgStr := `SELECT header_cids.block_number + 1 AS start, min(fr.block_number) - 1 AS stop FROM header_cids
				LEFT JOIN header_cids r on header_cids.block_number = r.block_number - 1
				LEFT JOIN header_cids fr on header_cids.block_number < fr.block_number
				WHERE r.block_number is NULL and fr.block_number IS NOT NULL
				GROUP BY header_cids.block_number, r.block_number`
	gaps := make([]gap, 0)
	err := ecr.db.Select(&gaps, pgStr)
	if err != nil {
		return nil, err
	}
	gapRanges := make([][2]int64, 0)
	for _, gap := range gaps {
		gapRanges = append(gapRanges, [2]int64{gap.Start, gap.Stop})
	}
	return gapRanges, nil
}
