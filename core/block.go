package core

type Block struct {
	Number       int64
	GasLimit     int64
	GasUsed      int64
	Time         int64
	Transactions []Transaction
}
