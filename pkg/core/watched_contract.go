package core

type WatchedContract struct {
	Hash         string
	Transactions []Transaction
}

type ContractAttribute struct {
	Name string
	Type string
}
