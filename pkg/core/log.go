package core

type Log struct {
	BlockNumber int64
	TxHash      string
	Address     string
	Topics      map[int]string
	Index       int64
	Data        string
}
