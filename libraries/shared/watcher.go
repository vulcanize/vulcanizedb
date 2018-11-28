package shared

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
)

type Watcher struct {
	Transformers []shared.Transformer
	DB           postgres.DB
	Blockchain   core.BlockChain
}

func (watcher *Watcher) AddTransformers(us []shared.TransformerInitializer) {
	for _, transformerInitializer := range us {
		transformer := transformerInitializer(&watcher.DB)
		watcher.Transformers = append(watcher.Transformers, transformer)
	}
}

func (watcher *Watcher) Execute() error {
	// TODO Solve checkedHeadersColumn issue
	// TODO Handle start and end numbers in transformers?
	var missingHeaders []core.Header

	// TODO Get contract addresses and topic0s
	var logs []types.Log

	var err error
	for _, transformer := range watcher.Transformers {
		err = transformer.Execute(logs, missingHeaders)
	}
	return err
}
