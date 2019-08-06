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

package repositories

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

const insertHeaderSyncLogQuery = `INSERT INTO header_sync_logs
		(header_id, address, topics, data, block_number, block_hash, tx_index, tx_hash, log_index, raw)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) ON CONFLICT DO NOTHING`

type HeaderSyncLogRepository struct {
	db *postgres.DB
}

func NewHeaderSyncLogRepository(db *postgres.DB) HeaderSyncLogRepository {
	return HeaderSyncLogRepository{db: db}
}

type headerSyncLog struct {
	ID          int64
	HeaderID    int64 `db:"header_id"`
	Address     string
	Topics      pq.ByteaArray
	Data        []byte
	BlockNumber uint64 `db:"block_number"`
	BlockHash   string `db:"block_hash"`
	TxHash      string `db:"tx_hash"`
	TxIndex     uint   `db:"tx_index"`
	LogIndex    uint   `db:"log_index"`
	Transformed bool
	Raw         []byte
}

func (repository HeaderSyncLogRepository) GetUntransformedHeaderSyncLogs() ([]core.HeaderSyncLog, error) {
	rows, queryErr := repository.db.Queryx(`SELECT * FROM public.header_sync_logs WHERE transformed = false`)
	if queryErr != nil {
		return nil, queryErr
	}

	var results []core.HeaderSyncLog
	for rows.Next() {
		var rawLog headerSyncLog
		scanErr := rows.StructScan(&rawLog)
		if scanErr != nil {
			return nil, scanErr
		}
		var logTopics []common.Hash
		for _, topic := range rawLog.Topics {
			logTopics = append(logTopics, common.BytesToHash(topic))
		}
		reconstructedLog := types.Log{
			Address:     common.HexToAddress(rawLog.Address),
			Topics:      logTopics,
			Data:        rawLog.Data,
			BlockNumber: rawLog.BlockNumber,
			TxHash:      common.HexToHash(rawLog.TxHash),
			TxIndex:     rawLog.TxIndex,
			BlockHash:   common.HexToHash(rawLog.BlockHash),
			Index:       rawLog.LogIndex,
			// TODO: revisit if not cascade deleting logs when header removed
			// currently, fetched logs are cascade deleted if removed
			Removed: false,
		}
		result := core.HeaderSyncLog{
			ID:          rawLog.ID,
			HeaderID:    rawLog.HeaderID,
			Log:         reconstructedLog,
			Transformed: rawLog.Transformed,
		}
		// TODO: Consider returning each result async to avoid keeping large result sets in memory
		results = append(results, result)
	}

	return results, nil
}

func (repository HeaderSyncLogRepository) CreateHeaderSyncLogs(headerID int64, logs []types.Log) error {
	tx, txErr := repository.db.Beginx()
	if txErr != nil {
		return txErr
	}
	for _, log := range logs {
		err := insertLog(headerID, log, tx)
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				logrus.Errorf("failed to rollback header sync log insert: %s", rollbackErr.Error())
			}
			return err
		}
	}
	return tx.Commit()
}

func insertLog(headerID int64, log types.Log, tx *sqlx.Tx) error {
	topics := buildTopics(log)
	raw, jsonErr := log.MarshalJSON()
	if jsonErr != nil {
		return jsonErr
	}
	_, insertErr := tx.Exec(insertHeaderSyncLogQuery, headerID, log.Address.Hex(), topics, log.Data, log.BlockNumber,
		log.BlockHash.Hex(), log.TxIndex, log.TxHash.Hex(), log.Index, raw)
	return insertErr
}

func buildTopics(log types.Log) pq.ByteaArray {
	var topics pq.ByteaArray
	for _, topic := range log.Topics {
		topics = append(topics, topic.Bytes())
	}
	return topics
}
