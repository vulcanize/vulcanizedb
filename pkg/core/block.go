package core

type Block struct {
	Difficulty   int64
	GasLimit     int64
	GasUsed      int64
	Hash         string
	Nonce        string
	Number       int64
	ParentHash   string
	Size         int64
	Time         int64
	Transactions []Transaction
	UncleHash    string
	IsFinal      bool
}
