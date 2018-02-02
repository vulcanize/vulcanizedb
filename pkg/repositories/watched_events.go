package repositories

type WatchedEventLog struct {
	Name        string `json:"name"`                           // name
	BlockNumber int64  `json:"block_number" db:"block_number"` // block_number
	Address     string `json:"address"`                        // address
	TxHash      string `json:"tx_hash" db:"tx_hash"`           // tx_hash
	Index       int64  `json:"index"`                          // index
	Topic0      string `json:"topic0"`                         // topic0
	Topic1      string `json:"topic1"`                         // topic1
	Topic2      string `json:"topic2"`                         // topic2
	Topic3      string `json:"topic3"`                         // topic3
	Data        string `json:"data"`                           // data
}

type WatchedEventLogs interface {
	AllWatchedEventLogs() ([]*WatchedEventLog, error)
}

func (pg *DB) AllWatchedEventLogs() ([]*WatchedEventLog, error) {
	rows, err := pg.Db.Queryx("SELECT name, block_number, address, tx_hash, index, topic0, topic1, topic2, topic3, data FROM watched_event_logs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lgs := make([]*WatchedEventLog, 0)
	for rows.Next() {
		lg := new(WatchedEventLog)
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
