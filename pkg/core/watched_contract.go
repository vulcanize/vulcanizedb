package core

type WatchedContract struct {
	Abi          string
	Hash         string
	Transactions []Transaction
}
