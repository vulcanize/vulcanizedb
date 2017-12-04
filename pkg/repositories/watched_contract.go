package repositories

import "github.com/8thlight/vulcanizedb/pkg/core"

type WatchedContract struct {
	Abi          string
	Hash         string
	Transactions []core.Transaction
}
