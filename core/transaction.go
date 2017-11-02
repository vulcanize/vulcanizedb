package core

type Transaction struct {
	Hash     string
	Data     []byte
	Nonce    uint64
	To       string
	GasLimit int64
	GasPrice int64
	Value    int64
}
