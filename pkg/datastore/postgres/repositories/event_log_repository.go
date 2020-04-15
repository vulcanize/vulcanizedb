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
	"github.com/makerdao/vulcanizedb/libraries/shared/repository"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/sirupsen/logrus"
)

const insertEventLogQuery = `INSERT INTO public.event_logs
		(header_id, address, topics, data, block_number, block_hash, tx_index, tx_hash, log_index, raw)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) ON CONFLICT DO NOTHING`

type EventLogRepository struct {
	db *postgres.DB
}

func NewEventLogRepository(db *postgres.DB) EventLogRepository {
	return EventLogRepository{
		db: db,
	}
}

type rawEventLog struct {
	ID          int64
	HeaderID    int64 `db:"header_id"`
	Address     int64
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

func (repo EventLogRepository) GetUntransformedEventLogs() ([]core.EventLog, error) {
	rows, queryErr := repo.db.Queryx(`SELECT * FROM public.event_logs WHERE transformed is false`)
	if queryErr != nil {
		return nil, queryErr
	}

	var results []core.EventLog
	for rows.Next() {
		var rawLog rawEventLog
		scanErr := rows.StructScan(&rawLog)
		if scanErr != nil {
			return nil, scanErr
		}
		var logTopics []common.Hash
		for _, topic := range rawLog.Topics {
			logTopics = append(logTopics, common.BytesToHash(topic))
		}
		address, addrErr := repository.GetAddressById(repo.db, rawLog.Address)
		if addrErr != nil {
			return nil, addrErr
		}
		reconstructedLog := types.Log{
			Address:     common.HexToAddress(address),
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
		result := core.EventLog{
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

func (repo EventLogRepository) CreateEventLogs(headerID int64, logs []types.Log) error {
	tx, txErr := repo.db.Beginx()
	if txErr != nil {
		return txErr
	}
	for _, log := range logs {
		err := repo.insertLog(headerID, log, tx)
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

func (repo EventLogRepository) insertLog(headerID int64, log types.Log, tx *sqlx.Tx) error {
	topics := buildTopics(log)
	raw, jsonErr := log.MarshalJSON()
	if jsonErr != nil {
		return jsonErr
	}
	addressID, addrErr := repository.GetOrCreateAddressInTransaction(tx, log.Address.Hex())
	if addrErr != nil {
		return addrErr
	}
	_, insertErr := tx.Exec(insertEventLogQuery, headerID, addressID, topics, log.Data, log.BlockNumber,
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
