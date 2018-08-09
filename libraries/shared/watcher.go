package shared

import (
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
		transformer := transformerInitializer(&watcher.DB, watcher.Blockchain)
		watcher.Transformers = append(watcher.Transformers, transformer)
	}
}

func (watcher *Watcher) Execute() error {
	var err error
	for _, transformer := range watcher.Transformers {
		err = transformer.Execute()
	}
	return err
}
