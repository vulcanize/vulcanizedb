package postgres

import (
	"database/sql"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/repositories"
)

func (db *DB) AllWatchedEventLogs() ([]*core.WatchedEventLog, error) {
	rows, err := db.DB.Queryx(`SELECT name, block_number, address, tx_hash, index, topic0, topic1, topic2, topic3, data FROM watched_event_logs`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lgs := make([]*core.WatchedEventLog, 0)
	for rows.Next() {
		lg := new(core.WatchedEventLog)
		err := rows.StructScan(lg)
		if err != nil {
			return nil, err
		}
		lgs = append(lgs, lg)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return lgs, nil
}

func (db DB) GetWatchedEvent(name string) (*core.WatchedEventLog, error) {
	watchedEventLog := core.WatchedEventLog{}
	err := db.DB.Get(&watchedEventLog,
		`SELECT name, 
       block_number, 
       address, 
       tx_hash, 
       topic0, 
       topic1, 
       topic2, 
       topic3, 
       data
FROM watched_event_logs
WHERE name = $1`, name)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return &core.WatchedEventLog{}, repositories.ErrFilterDoesNotExist(name)
		default:
			return &core.WatchedEventLog{}, err
		}
	}
	return &watchedEventLog, nil
}
