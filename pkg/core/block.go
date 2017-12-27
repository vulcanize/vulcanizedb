package core

type Block struct {
	Difficulty   int64
	ExtraData    string
	GasLimit     int64
	GasUsed      int64
	Hash         string
	IsFinal      bool
	Miner        string
	Nonce        string
	Number       int64
	ParentHash   string
	Size         int64
	Time         int64
	Transactions []Transaction
	UncleHash    string
}
