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
	"database/sql"
	"github.com/sirupsen/logrus"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type FullSyncLogRepository struct {
	*postgres.DB
}

func (repository FullSyncLogRepository) CreateLogs(lgs []core.FullSyncLog, receiptID int64) error {
	tx, _ := repository.DB.Beginx()
	for _, tlog := range lgs {
		_, insertLogErr := tx.Exec(
			`INSERT INTO full_sync_logs (block_number, address, tx_hash, index, topic0, topic1, topic2, topic3, data, receipt_id)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
                `,
			tlog.BlockNumber, tlog.Address, tlog.TxHash, tlog.Index, tlog.Topics[0], tlog.Topics[1], tlog.Topics[2], tlog.Topics[3], tlog.Data, receiptID,
		)
		if insertLogErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				logrus.Error("CreateLogs: could not perform rollback: ", rollbackErr)
			}
			return postgres.ErrDBInsertFailed(insertLogErr)
		}
	}
	err := tx.Commit()
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			logrus.Error("CreateLogs: could not perform rollback: ", err)
		}
		return postgres.ErrDBInsertFailed(err)
	}
	return nil
}

func (repository FullSyncLogRepository) GetLogs(address string, blockNumber int64) ([]core.FullSyncLog, error) {
	logRows, err := repository.DB.Query(
		`SELECT block_number,
					  address,
					  tx_hash,
					  index,
					  topic0,
					  topic1,
					  topic2,
					  topic3,
					  data
				FROM full_sync_logs
				WHERE address = $1 AND block_number = $2
				ORDER BY block_number DESC`, address, blockNumber)
	if err != nil {
		return []core.FullSyncLog{}, err
	}
	return repository.loadLogs(logRows)
}

func (repository FullSyncLogRepository) loadLogs(logsRows *sql.Rows) ([]core.FullSyncLog, error) {
	var lgs []core.FullSyncLog
	for logsRows.Next() {
		var blockNumber int64
		var address string
		var txHash string
		var index int64
		var data string
		var topics core.Topics
		err := logsRows.Scan(&blockNumber, &address, &txHash, &index, &topics[0], &topics[1], &topics[2], &topics[3], &data)
		if err != nil {
			logrus.Error("loadLogs: Error scanning a row in logRows: ", err)
			return []core.FullSyncLog{}, err
		}
		lg := core.FullSyncLog{
			BlockNumber: blockNumber,
			TxHash:      txHash,
			Address:     address,
			Index:       index,
			Data:        data,
		}
		for i, topic := range topics {
			lg.Topics[i] = topic
		}
		lgs = append(lgs, lg)
	}
	return lgs, nil
}
