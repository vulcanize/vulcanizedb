package core

type Header struct {
	Id          int64
	BlockNumber int64 `db:"block_number"`
	Hash        string
	Raw         []byte
}
