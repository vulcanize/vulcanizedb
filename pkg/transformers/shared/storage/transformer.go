package storage

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/shared"
)

type Transformer interface {
	Execute(row shared.StorageDiffRow) error
	ContractAddress() common.Address
}

type TransformerInitializer func(db *postgres.DB) Transformer
