package postgres

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type WatchedEventRepository struct {
	*DB
}

func (watchedEventRepository WatchedEventRepository) GetWatchedEvents(name string) ([]*core.WatchedEvent, error) {
	rows, err := watchedEventRepository.DB.Queryx(`SELECT name, block_number, address, tx_hash, index, topic0, topic1, topic2, topic3, data FROM watched_event_logs where name=$1`, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lgs := make([]*core.WatchedEvent, 0)
	for rows.Next() {
		lg := new(core.WatchedEvent)
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
