package shared

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type Handler interface {
	Execute() error
}

type HandlerInitializer func(db *postgres.DB, blockchain core.Blockchain) Handler

func HexToInt64(byteString string) int64 {
	intHash := common.HexToHash(byteString)
	return intHash.Big().Int64()
}

func HexToString(byteString string) string {
	value := common.HexToHash(byteString)
	return value.Big().String()
}
