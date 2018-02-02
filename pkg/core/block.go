package core

type Block struct {
	Reward       float64 `db:"block_reward"`
	Difficulty   int64   `db:"block_difficulty"`
	ExtraData    string  `db:"block_extra_data"`
	GasLimit     int64   `db:"block_gaslimit"`
	GasUsed      int64   `db:"block_gasused"`
	Hash         string  `db:"block_hash"`
	IsFinal      bool    `db:"is_final"`
	Miner        string  `db:"block_miner"`
	Nonce        string  `db:"block_nonce"`
	Number       int64   `db:"block_number"`
	ParentHash   string  `db:"block_parenthash"`
	Size         int64   `db:"block_size"`
	Time         int64   `db:"block_time"`
	Transactions []Transaction
	UncleHash    string  `db:"uncle_hash"`
	UnclesReward float64 `db:"block_uncles_reward"`
}
