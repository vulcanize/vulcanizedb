package core

type Log struct {
	BlockNumber int64
	TxHash      string
	Address     string
	Topics
	Index   int64
	Data    string
	Removed bool
}
