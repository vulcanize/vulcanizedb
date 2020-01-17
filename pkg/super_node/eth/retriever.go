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
	"math/big"

	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/config"
)

// CIDRetriever satisfies the CIDRetriever interface for ethereum
type CIDRetriever struct {
	db *postgres.DB
}

// NewCIDRetriever returns a pointer to a new CIDRetriever which supports the CIDRetriever interface
func NewCIDRetriever(db *postgres.DB) *CIDRetriever {
	return &CIDRetriever{
		db: db,
	}
}

// RetrieveFirstBlockNumber is used to retrieve the first block number in the db
func (ecr *CIDRetriever) RetrieveFirstBlockNumber() (int64, error) {
	var blockNumber int64
	err := ecr.db.Get(&blockNumber, "SELECT block_number FROM header_cids ORDER BY block_number ASC LIMIT 1")
	return blockNumber, err
}

// RetrieveLastBlockNumber is used to retrieve the latest block number in the db
func (ecr *CIDRetriever) RetrieveLastBlockNumber() (int64, error) {
	var blockNumber int64
	err := ecr.db.Get(&blockNumber, "SELECT block_number FROM header_cids ORDER BY block_number DESC LIMIT 1 ")
	return blockNumber, err
}

// Retrieve is used to retrieve all of the CIDs which conform to the passed StreamFilters
func (ecr *CIDRetriever) Retrieve(filter interface{}, blockNumber int64) (interface{}, bool, error) {
	streamFilter, ok := filter.(*config.EthSubscription)
	if !ok {
		return nil, true, fmt.Errorf("eth retriever expected filter type %T got %T", &config.EthSubscription{}, filter)
	}
	log.Debug("retrieving cids")
	tx, err := ecr.db.Beginx()
	if err != nil {
		return nil, true, err
	}

	cw := new(CIDWrapper)
	cw.BlockNumber = big.NewInt(blockNumber)
	// Retrieve cached header CIDs
	if !streamFilter.HeaderFilter.Off {
		cw.Headers, err = ecr.RetrieveHeaderCIDs(tx, blockNumber)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				log.Error(err)
			}
			log.Error("header cid retrieval error")
			return nil, true, err
		}
		if streamFilter.HeaderFilter.Uncles {
			cw.Uncles, err = ecr.RetrieveUncleCIDs(tx, blockNumber)
			if err != nil {
				if err := tx.Rollback(); err != nil {
					log.Error(err)
				}
				log.Error("uncle cid retrieval error")
				return nil, true, err
			}
		}
	}
	// Retrieve cached trx CIDs
	if !streamFilter.TxFilter.Off {
		cw.Transactions, err = ecr.RetrieveTrxCIDs(tx, streamFilter.TxFilter, blockNumber)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				log.Error(err)
			}
			log.Error("transaction cid retrieval error")
			return nil, true, err
		}
	}
	trxIds := make([]int64, 0, len(cw.Transactions))
	for _, tx := range cw.Transactions {
		trxIds = append(trxIds, tx.ID)
	}
	// Retrieve cached receipt CIDs
	if !streamFilter.ReceiptFilter.Off {
		cw.Receipts, err = ecr.RetrieveRctCIDs(tx, streamFilter.ReceiptFilter, blockNumber, trxIds)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				log.Error(err)
			}
			log.Error("receipt cid retrieval error")
			return nil, true, err
		}
	}
	// Retrieve cached state CIDs
	if !streamFilter.StateFilter.Off {
		cw.StateNodes, err = ecr.RetrieveStateCIDs(tx, streamFilter.StateFilter, blockNumber)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				log.Error(err)
			}
			log.Error("state cid retrieval error")
			return nil, true, err
		}
	}
	// Retrieve cached storage CIDs
	if !streamFilter.StorageFilter.Off {
		cw.StorageNodes, err = ecr.RetrieveStorageCIDs(tx, streamFilter.StorageFilter, blockNumber)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				log.Error(err)
			}
			log.Error("storage cid retrieval error")
			return nil, true, err
		}
	}
	return cw, empty(cw), tx.Commit()
}

func empty(cidWrapper *CIDWrapper) bool {
	if len(cidWrapper.Transactions) > 0 || len(cidWrapper.Headers) > 0 || len(cidWrapper.Uncles) > 0 || len(cidWrapper.Receipts) > 0 || len(cidWrapper.StateNodes) > 0 || len(cidWrapper.StorageNodes) > 0 {
		return false
	}
	return true
}

// RetrieveHeaderCIDs retrieves and returns all of the header cids at the provided blockheight
func (ecr *CIDRetriever) RetrieveHeaderCIDs(tx *sqlx.Tx, blockNumber int64) ([]HeaderModel, error) {
	log.Debug("retrieving header cids for block ", blockNumber)
	headers := make([]HeaderModel, 0)
	pgStr := `SELECT * FROM header_cids
				WHERE block_number = $1 AND uncle IS FALSE`
	err := tx.Select(&headers, pgStr, blockNumber)
	return headers, err
}

// RetrieveUncleCIDs retrieves and returns all of the uncle cids at the provided blockheight
func (ecr *CIDRetriever) RetrieveUncleCIDs(tx *sqlx.Tx, blockNumber int64) ([]HeaderModel, error) {
	log.Debug("retrieving header cids for block ", blockNumber)
	headers := make([]HeaderModel, 0)
	pgStr := `SELECT * FROM header_cids
				WHERE block_number = $1 AND uncle IS TRUE`
	err := tx.Select(&headers, pgStr, blockNumber)
	return headers, err
}

// RetrieveTrxCIDs retrieves and returns all of the trx cids at the provided blockheight that conform to the provided filter parameters
// also returns the ids for the returned transaction cids
func (ecr *CIDRetriever) RetrieveTrxCIDs(tx *sqlx.Tx, txFilter config.TxFilter, blockNumber int64) ([]TxModel, error) {
	log.Debug("retrieving transaction cids for block ", blockNumber)
	args := make([]interface{}, 0, 3)
	results := make([]TxModel, 0)
	pgStr := `SELECT transaction_cids.id, transaction_cids.header_id,
 			transaction_cids.tx_hash, transaction_cids.cid,
 			transaction_cids.dst, transaction_cids.src
 			FROM transaction_cids INNER JOIN header_cids ON (transaction_cids.header_id = header_cids.id)
			WHERE header_cids.block_number = $1`
	args = append(args, blockNumber)
	if len(txFilter.Dst) > 0 {
		pgStr += ` AND transaction_cids.dst = ANY($2::VARCHAR(66)[])`
		args = append(args, pq.Array(txFilter.Dst))
	}
	if len(txFilter.Src) > 0 {
		pgStr += ` AND transaction_cids.src = ANY($3::VARCHAR(66)[])`
		args = append(args, pq.Array(txFilter.Src))
	}
	err := tx.Select(&results, pgStr, args...)
	if err != nil {
		return nil, err
	}
	return results, nil
}

// RetrieveRctCIDs retrieves and returns all of the rct cids at the provided blockheight that conform to the provided
// filter parameters and correspond to the provided tx ids
func (ecr *CIDRetriever) RetrieveRctCIDs(tx *sqlx.Tx, rctFilter config.ReceiptFilter, blockNumber int64, trxIds []int64) ([]ReceiptModel, error) {
	log.Debug("retrieving receipt cids for block ", blockNumber)
	args := make([]interface{}, 0, 4)
	pgStr := `SELECT receipt_cids.id, receipt_cids.tx_id, receipt_cids.cid,
 			receipt_cids.contract, receipt_cids.topic0s
 			FROM receipt_cids, transaction_cids, header_cids
			WHERE receipt_cids.tx_id = transaction_cids.id 
			AND transaction_cids.header_id = header_cids.id
			AND header_cids.block_number = $1`
	args = append(args, blockNumber)
	if len(rctFilter.Topic0s) > 0 {
		pgStr += ` AND ((receipt_cids.topic0s && $2::VARCHAR(66)[]`
		args = append(args, pq.Array(rctFilter.Topic0s))
		if len(rctFilter.Contracts) > 0 {
			pgStr += ` AND receipt_cids.contract = ANY($3::VARCHAR(66)[]))`
			args = append(args, pq.Array(rctFilter.Contracts))
			if rctFilter.MatchTxs && len(trxIds) > 0 {
				pgStr += ` OR receipt_cids.tx_id = ANY($4::INTEGER[]))`
				args = append(args, pq.Array(trxIds))
			} else {
				pgStr += `)`
			}
		} else {
			pgStr += `)`
			if rctFilter.MatchTxs && len(trxIds) > 0 {
				pgStr += ` OR receipt_cids.tx_id = ANY($3::INTEGER[]))`
				args = append(args, pq.Array(trxIds))
			} else {
				pgStr += `)`
			}
		}
	} else {
		if len(rctFilter.Contracts) > 0 {
			pgStr += ` AND (receipt_cids.contract = ANY($2::VARCHAR(66)[])`
			args = append(args, pq.Array(rctFilter.Contracts))
			if rctFilter.MatchTxs && len(trxIds) > 0 {
				pgStr += ` OR receipt_cids.tx_id = ANY($3::INTEGER[]))`
				args = append(args, pq.Array(trxIds))
			} else {
				pgStr += `)`
			}
		} else if rctFilter.MatchTxs && len(trxIds) > 0 {
			pgStr += ` AND receipt_cids.tx_id = ANY($2::INTEGER[])`
			args = append(args, pq.Array(trxIds))
		}
	}
	receiptCids := make([]ReceiptModel, 0)
	err := tx.Select(&receiptCids, pgStr, args...)
	if err != nil {
		println(pgStr)
		println("FUCK YOU\r\n\r\n\r\n")
	}
	return receiptCids, err
}

// RetrieveStateCIDs retrieves and returns all of the state node cids at the provided blockheight that conform to the provided filter parameters
func (ecr *CIDRetriever) RetrieveStateCIDs(tx *sqlx.Tx, stateFilter config.StateFilter, blockNumber int64) ([]StateNodeModel, error) {
	log.Debug("retrieving state cids for block ", blockNumber)
	args := make([]interface{}, 0, 2)
	pgStr := `SELECT state_cids.id, state_cids.header_id,
			state_cids.state_key, state_cids.leaf, state_cids.cid
			FROM state_cids INNER JOIN header_cids ON (state_cids.header_id = header_cids.id)
			WHERE header_cids.block_number = $1`
	args = append(args, blockNumber)
	addrLen := len(stateFilter.Addresses)
	if addrLen > 0 {
		keys := make([]string, 0, addrLen)
		for _, addr := range stateFilter.Addresses {
			keys = append(keys, HexToKey(addr).Hex())
		}
		pgStr += ` AND state_cids.state_key = ANY($2::VARCHAR(66)[])`
		args = append(args, pq.Array(keys))
	}
	if !stateFilter.IntermediateNodes {
		pgStr += ` AND state_cids.leaf = TRUE`
	}
	stateNodeCIDs := make([]StateNodeModel, 0)
	err := tx.Select(&stateNodeCIDs, pgStr, args...)
	return stateNodeCIDs, err
}

// RetrieveStorageCIDs retrieves and returns all of the storage node cids at the provided blockheight that conform to the provided filter parameters
func (ecr *CIDRetriever) RetrieveStorageCIDs(tx *sqlx.Tx, storageFilter config.StorageFilter, blockNumber int64) ([]StorageNodeWithStateKeyModel, error) {
	log.Debug("retrieving storage cids for block ", blockNumber)
	args := make([]interface{}, 0, 3)
	pgStr := `SELECT storage_cids.id, storage_cids.state_id, storage_cids.storage_key,
 			storage_cids.leaf, storage_cids.cid, state_cids.state_key FROM storage_cids, state_cids, header_cids
			WHERE storage_cids.state_id = state_cids.id 
			AND state_cids.header_id = header_cids.id
			AND header_cids.block_number = $1`
	args = append(args, blockNumber)
	addrLen := len(storageFilter.Addresses)
	if addrLen > 0 {
		keys := make([]string, 0, addrLen)
		for _, addr := range storageFilter.Addresses {
			keys = append(keys, HexToKey(addr).Hex())
		}
		pgStr += ` AND state_cids.state_key = ANY($2::VARCHAR(66)[])`
		args = append(args, pq.Array(keys))
		if len(storageFilter.StorageKeys) > 0 {
			pgStr += ` AND storage_cids.storage_key = ANY($3::VARCHAR(66)[])`
			args = append(args, pq.Array(storageFilter.StorageKeys))
		}
	} else if len(storageFilter.StorageKeys) > 0 {
		pgStr += ` AND storage_cids.storage_key = ANY($2::VARCHAR(66)[])`
		args = append(args, pq.Array(storageFilter.StorageKeys))
	}
	if !storageFilter.IntermediateNodes {
		pgStr += ` AND storage_cids.leaf = TRUE`
	}
	storageNodeCIDs := make([]StorageNodeWithStateKeyModel, 0)
	err := tx.Select(&storageNodeCIDs, pgStr, args...)
	return storageNodeCIDs, err
}

// RetrieveGapsInData is used to find the the block numbers at which we are missing data in the db
func (ecr *CIDRetriever) RetrieveGapsInData() ([]shared.Gap, error) {
	pgStr := `SELECT header_cids.block_number + 1 AS start, min(fr.block_number) - 1 AS stop FROM header_cids
				LEFT JOIN header_cids r on header_cids.block_number = r.block_number - 1
				LEFT JOIN header_cids fr on header_cids.block_number < fr.block_number
				WHERE r.block_number is NULL and fr.block_number IS NOT NULL
				GROUP BY header_cids.block_number, r.block_number`
	results := make([]struct {
		Start uint64 `db:"start"`
		Stop  uint64 `db:"stop"`
	}, 0)
	err := ecr.db.Select(&results, pgStr)
	if err != nil {
		return nil, err
	}
	gaps := make([]shared.Gap, len(results))
	for i, res := range results {
		gaps[i] = shared.Gap{
			Start: res.Start,
			Stop:  res.Stop,
		}
	}
	return gaps, nil
}

func (ecr *CIDRetriever) Database() *postgres.DB {
	return ecr.db
}
