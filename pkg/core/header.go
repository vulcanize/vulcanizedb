package core

type Header struct {
	BlockNumber int64 `db:"block_number"`
	Hash        string
	Raw         []byte
}
