package shared

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type Watcher struct {
	Transformers []Transformer
	DB           postgres.DB
	Blockchain   core.Blockchain
}

func (watcher *Watcher) AddTransformers(us []TransformerInitializer) {
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
