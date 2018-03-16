package shared

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type Transformer interface {
	Execute() error
}

type TransformerInitializer func(db *postgres.DB, blockchain core.Blockchain) Transformer

func HexToInt64(byteString string) int64 {
	value := common.HexToHash(byteString)
	return value.Big().Int64()
}

func HexToString(byteString string) string {
	value := common.HexToHash(byteString)
	return value.Big().String()
}
