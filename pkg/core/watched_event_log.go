package core

type WatchedEvent struct {
	LogID       int64  `json:"log_id" db:"id"`
	Name        string `json:"name"`
	BlockNumber int64  `json:"block_number" db:"block_number"`
	Address     string `json:"address"`
	TxHash      string `json:"tx_hash" db:"tx_hash"`
	Index       int64  `json:"index"`
	Topic0      string `json:"topic0"`
	Topic1      string `json:"topic1"`
	Topic2      string `json:"topic2"`
	Topic3      string `json:"topic3"`
	Data        string `json:"data"`
}
