package core

type Block struct {
	Reward       float64 `db:"reward"`
	Difficulty   int64   `db:"difficulty"`
	ExtraData    string  `db:"extra_data"`
	GasLimit     int64   `db:"gaslimit"`
	GasUsed      int64   `db:"gasused"`
	Hash         string  `db:"hash"`
	IsFinal      bool    `db:"is_final"`
	Miner        string  `db:"miner"`
	Nonce        string  `db:"nonce"`
	Number       int64   `db:"number"`
	ParentHash   string  `db:"parenthash"`
	Size         int64   `db:"size"`
	Time         int64   `db:"time"`
	Transactions []Transaction
	UncleHash    string  `db:"uncle_hash"`
	UnclesReward float64 `db:"uncles_reward"`
}
