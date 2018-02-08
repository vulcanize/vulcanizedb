package core

type WatchedEvent struct {
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
