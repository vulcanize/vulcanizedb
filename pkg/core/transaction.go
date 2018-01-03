package core

type Transaction struct {
	Hash     string
	Data     string
	Nonce    uint64
	To       string
	From     string
	GasLimit int64
	GasPrice int64
	Receipt
	Value int64
}
