package repositories

import (
	"context"

	"database/sql"

	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type LogsRepository interface {
	FindLogs(address string, blockNumber int64) []core.Log
	CreateLogs(logs []core.Log) error
}

func (repository Postgres) CreateLogs(logs []core.Log) error {
	tx, _ := repository.Db.BeginTx(context.Background(), nil)
	for _, tlog := range logs {
		_, err := tx.Exec(
			`INSERT INTO logs (block_number, address, tx_hash, index, topic0, topic1, topic2, topic3, data)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
                `,
			tlog.BlockNumber, tlog.Address, tlog.TxHash, tlog.Index, tlog.Topics[0], tlog.Topics[1], tlog.Topics[2], tlog.Topics[3], tlog.Data,
		)
		if err != nil {
			tx.Rollback()
			return ErrDBInsertFailed
		}
	}
	tx.Commit()
	return nil
}

func (repository Postgres) FindLogs(address string, blockNumber int64) []core.Log {
	logRows, _ := repository.Db.Query(
		`SELECT block_number,
					  address,
					  tx_hash,
					  index,
					  topic0,
					  topic1,
					  topic2,
					  topic3,
					  data
				FROM logs
				WHERE address = $1 AND block_number = $2
				ORDER BY block_number DESC`, address, blockNumber)
	return repository.loadLogs(logRows)
}

func (repository Postgres) loadLogs(logsRows *sql.Rows) []core.Log {
	var logs []core.Log
	for logsRows.Next() {
		var blockNumber int64
		var address string
		var txHash string
		var index int64
		var data string
		var topics core.Topics
		logsRows.Scan(&blockNumber, &address, &txHash, &index, &topics[0], &topics[1], &topics[2], &topics[3], &data)
		log := core.Log{
			BlockNumber: blockNumber,
			TxHash:      txHash,
			Address:     address,
			Index:       index,
			Data:        data,
		}
		for i, topic := range topics {
			log.Topics[i] = topic
		}
		logs = append(logs, log)
	}
	return logs
}
