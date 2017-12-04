package repositories

import "github.com/8thlight/vulcanizedb/pkg/core"

type WatchedContract struct {
	Hash         string
	Transactions []core.Transaction
}
